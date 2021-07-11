package main

import (
	"fmt"
	"strings"
)

type termVal struct {
	fieldName   string
	injectPoint []string
	raw         string
}

type pathTerm struct {
	term termVal
}

type pathPattern struct {
	bindFields []string
	base       []string
	terms      []pathTerm
	pathMap    []int
}

func pathParse(path string) (*pathPattern, error) {
	// /svc/{a}/b/{c}
	// /svc{/sd}
	var comp = []rune{}
	var temp = []rune{}
	var inClosure = false
	var termsI = 0
	var pTerms = []string{}
	var terms = []string{}
	var pathMap = []int{}
	for i, c := range path {
		if inClosure && i == len(path)-1 && c != '}' {
			return nil, fmt.Errorf("expected '}' got '%c' at col %d", c, i)
		}
		if !inClosure && i == len(path)-1 {
			comp = append(comp, c)
			pTerms = append(pTerms, string(comp))
		} else if c == '{' {
			if !inClosure {
				pTerms = append(pTerms, string(comp))
				comp = []rune{}
				termsI++
				inClosure = true
				continue
			} else {
				return nil, fmt.Errorf("expected '}' got '{' at col %d", i)
			}
		} else if c == '}' {
			if !inClosure {
				return nil, fmt.Errorf("expected '{' got '%c' at col %d", c, i)
			} else {
				if len(temp) < 1 {
					return nil, fmt.Errorf("expected field name, got '%c' at col %d", c, i)
				}
				inClosure = false
				terms = append(terms, string(temp))
				pathMap = append(pathMap, termsI)
				termsI++
				temp = []rune{}
				continue
			}
		}
		if inClosure {
			temp = append(temp, c)
		} else {
			comp = append(comp, c)
		}
	}

	bindFields := []string{}
	pathTerms := []pathTerm{}
	for i, term := range terms {
		kv := strings.SplitN(term, "=", 2)
		k, v := kv[0], ""
		if len(kv) > 1 {
			v = kv[1]
		}
		if len(k) < 1 {
			return nil, fmt.Errorf("expected field name got %+v at mapping %d (kv %+#v)", k, pathMap[i], kv)
		}
		bindFields = append(bindFields, k)
		if v == "" {
			pathTerms = append(pathTerms, pathTerm{
				term: termVal{
					fieldName: k,
					raw:       term,
				},
			})
		} else {
			sp, err := pathFindPoint(v)
			if err != nil {
				return nil, err
			}
			pathTerms = append(pathTerms, pathTerm{
				term: termVal{
					fieldName:   k,
					injectPoint: sp,
					raw:         v,
				},
			})
		}
	}
	return &pathPattern{
		bindFields: bindFields,
		base:       pTerms,
		terms:      pathTerms,
		pathMap:    pathMap,
	}, nil
}

func pathFindPoint(v string) (sp []string, err error) {
	cnt := strings.Count(v, "**")
	if cnt > 0 {
		err = fmt.Errorf("expected one wildchar got more than one in a row")
		return
	}
	sp = strings.Split(v, "*")
	return
}
