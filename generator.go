package main

import (
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	genQueryOptional = true
)

func generateFieldModifier(g *protogen.GeneratedFile, genTmpVar func() string, kind protoreflect.Kind, queryEscape bool) (mod func(string) string) {
	switch kind {
	case protoreflect.StringKind:
		mod = func(s string) string {
			if queryEscape {
				return g.QualifiedGoIdent(urlPackage.Ident("QueryEscape")) + "(" + s + ")"
			}
			return s
		}
	case protoreflect.BoolKind:
		mod = func(s string) string {
			t := genTmpVar()
			g.P(t, " := ", "\"false\"")
			g.P("if ", s, "{")
			g.P(t, " = \"true\"")
			g.P("}")
			return t
		}
	case protoreflect.Int64Kind, protoreflect.Int32Kind:
		// t := genTmpVar()
		mod = func(s string) string {
			// g.P(t, " := ", strconvPackage.Ident("FormatInt"), "(", s, ", 10)")
			// return t
			if kind == protoreflect.Int32Kind {
				s = "int64(" + s + ")"
			}
			return g.QualifiedGoIdent(strconvPackage.Ident("FormatInt")) + "(" + s + ", 10)"
		}
	case protoreflect.EnumKind:
		mod = func(s string) string {
			return s + ".String()"
		}
	default:
		mod = func(s string) string {
			return s
		}
	}
	return
}

func generateMessageUtil(g *protogen.GeneratedFile, msg *protogen.Message) (err error) {
	generateMessageToURLUtil(g, msg, "", "")
	return
}

func generateMessageFieldsURLSet(g *protogen.GeneratedFile, msg *protogen.Message, genTmpVar func() string, accessorValue, keyPrefix string) {
	fields, fOrder := extractMessageField(msg)
	// wip: ignore fields that already used in url path

	for _, fName := range fOrder {
		f, ok := fields[fName]
		if !ok {
			continue // skip
		}

		// non-repeated message
		// https://cloud.google.com/endpoints/docs/grpc-service-config/reference/rpc/google.api
		if f.Desc.IsList() {
			continue
		}

		if keyPrefix != "" {
			fName = keyPrefix + "." + fName
		}
		var inp = f.GoName
		if accessorValue != "" {
			inp = accessorValue + "." + inp
		}
		tt := f.Desc.Kind()
		switch tt {
		case protoreflect.MessageKind:
			if genQueryOptional {
				g.P("if ", "u."+inp, " != nil {")
			}
			generateMessageFieldsURLSet(g, f.Message, genTmpVar, inp, fName)
			if genQueryOptional {
				g.P("}")
			}
			continue
		}
		mod := generateFieldModifier(g, genTmpVar, tt, false)
		g.P("q.Set(\"", fName, "\", ", mod("u."+inp), ")")
	}
}

func generateMessageToURLUtil(g *protogen.GeneratedFile, msg *protogen.Message, accessorValue, keyPrefix string) {
	g.P()
	g.P("// QueryString returns http url.Values of ", msg.GoIdent.GoName)
	g.P("func (u *", msg.GoIdent.GoName, ") QueryString() ", urlPackage.Ident("Values"), "{")
	g.P("var q = ", urlPackage.Ident("Values"), "{}")

	tmpI := 0
	genTmpVar := func() string {
		defer func() {
			tmpI++
		}()
		return fmt.Sprintf("tm%d", tmpI)
	}

	generateMessageFieldsURLSet(g, msg, genTmpVar, accessorValue, keyPrefix)
	g.P("return q")
	g.P("}")
}
