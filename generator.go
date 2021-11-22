package main

import (
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	genBoolUseInt       = false
	genQueryOptional    = true
	genUnmarshalOptions = "DiscardUnknown: true"
	genUseValidator     = true // use go-validator https://github.com/asaskevich/govalidator
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
			bf := "false"
			bt := "true"
			if genBoolUseInt {
				bf = "0"
				bt = "1"
			}
			g.P(t, " := ", "\"", bf, "\"")
			g.P("if ", s, "{")
			g.P(t, " = \"", bt, "\"")
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
	case protoreflect.FloatKind:
		mod = func(s string) string {
			s = "float64(" + s + ")"
			return g.QualifiedGoIdent(strconvPackage.Ident("FormatFloat")) + "(" + s + ", 'f', -1, 64)"
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

func generateMessageFieldForMapVMessage(g *protogen.GeneratedFile, msg *protogen.Message, genTmpVar func() string, accessorValue, addtKey, keyPrefix string) {
	fields, order := extractMessageField(msg)
	for _, fk := range order {
		f, ok := fields[fk]
		if !ok {
			continue
		}
		tt := f.Desc.Kind()
		if tt == protoreflect.MessageKind {
			accessorVal := accessorValue
			if accessorVal != "" {
				accessorVal = accessorVal + "." + f.GoName
			}
			g.P("if ", accessorValue, " != nil {")
			generateMessageFieldForMapVMessage(g, f.Message, genTmpVar, accessorVal, addtKey, keyPrefix)
			g.P("}")
			continue
		}
		mod := generateFieldModifier(g, genTmpVar, tt, false)
		g.P("q.Set(\"", keyPrefix, "[\" + k + \"]", addtKey+"["+string(f.Desc.Name())+"]", "\", ", mod(accessorValue+"."+f.GoName), ") // ", tt)
	}
}
func generateMessageFieldForMap(g *protogen.GeneratedFile, msg *protogen.Message, genTmpVar func() string, accessorValue, keyPrefix string) (cont bool) {
	fs, fo := extractMessageField(msg) // key, value
	fk, fv := fs[fo[0]], fs[fo[1]]
	// cant have map or list as key
	if fk.Desc.IsMap() || fk.Desc.IsList() {
		return true
	}
	if fv.Desc.IsMap() || fv.Desc.IsList() {
		return true
	}
	cInp := accessorValue

	g.P("for k, v := range ", cInp, " {")
	fVal := "v"
	fName := keyPrefix
	addtKey := ""

	tt := fv.Desc.Kind()
	if tt == protoreflect.MessageKind {
		accessorValue = "v"
		generateMessageFieldForMapVMessage(g, fv.Message, genTmpVar, accessorValue, addtKey, fName)
	} else {
		mod := generateFieldModifier(g, genTmpVar, tt, false)
		g.P("q.Set(\"", fName, "[\" + k + \"]", addtKey, "\", ", mod(fVal), ")")
	}
	g.P("}")
	return false
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
		cInp := "u." + inp

		tt := f.Desc.Kind()
		switch tt {
		case protoreflect.MessageKind:
			if genQueryOptional {
				g.P("if ", cInp, " != nil {")
			}
			if f.Desc.IsMap() {
				if generateMessageFieldForMap(g, f.Message, genTmpVar, cInp, fName) {
					continue
				}
			} else {
				generateMessageFieldsURLSet(g, f.Message, genTmpVar, inp, fName)
			}
			if genQueryOptional {
				g.P("}")
			}
			continue
		}
		mod := generateFieldModifier(g, genTmpVar, tt, false)
		g.P("q.Set(\"", fName, "\", ", mod(cInp), ")")
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
