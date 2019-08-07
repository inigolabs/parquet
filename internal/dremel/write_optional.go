package dremel

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/parsyl/parquet/internal/fields"
	"github.com/parsyl/parquet/internal/structs"
)

func init() {
	funcs := template.FuncMap{
		"removeStar": func(s string) string {
			return strings.Replace(s, "*", "", 1)
		},
		"plusOne": func(i int) int { return i + 1 },
		"notNil":  func(x *ifElse) bool { return x != nil },
	}

	var err error
	writeTpl, err = template.New("output").Funcs(funcs).Parse(`func write{{.FuncName}}(x *{{.Field.Type}}, vals []{{removeStar .Field.TypeName}}, defs, reps []uint8) (int, int) {
	def := defs[0]
	switch def { {{range $i, $case := .Cases}}{{$def:=plusOne $i}}
	case {{$def}}:
	{{$defIndex := $.Field.DefIndex $def}}{{if $case.UseIf $.Field.Len $defIndex $.Field.NDefs $def $.MaxDef}}{{template "ifelse" $case}}{{else}}{{$case.If.Val}}{{end}}{{if eq $def $.MaxDef}}
	return 1, 1{{end}}{{end}}
	}

	return 0, 1
}`)
	if err != nil {
		log.Fatal(err)
	}

	writeTpl, err = writeTpl.Parse(`{{define "ifelse"}}if {{.If.Cond}} {
	{{.If.Val}}
} {{range $else := .ElseIf}} else if {{$else.Cond}} {
	{{$else.Val}}
}{{end}} {{if notNil .Else}} else {
	{{.Else.Val}}
} {{end}}{{end}}`)
	if err != nil {
		log.Fatal(err)
	}
}

var (
	writeTpl *template.Template
)

type writeInput struct {
	fields.Field
	Cases    []ifElses
	FuncName string
}

type ifElse struct {
	Cond string
	Val  string
}

type ifElses struct {
	If     ifElse
	ElseIf []ifElse
	Else   *ifElse
}

func (i ifElses) UseIf(l, defIndex, nDefs, def, maxDef int) bool {
	if nDefs == 1 && defIndex == l-1 {
		return false
	}
	if def < maxDef {
		return true
	}
	return len(i.ElseIf) > 0 || i.Else != nil
}

func writeOptional(f fields.Field) string {
	i := writeInput{
		Field:    f,
		FuncName: strings.Join(f.FieldNames, ""),
		Cases:    writeOptionalCases(f),
	}

	var buf bytes.Buffer
	err := writeTpl.Execute(&buf, i)
	if err != nil {
		log.Fatal(err) //TODO: return error
	}
	return string(buf.Bytes())
}

func writeOptionalCases(f fields.Field) []ifElses {
	var out []ifElses
	for def := 1; def <= defs(f); def++ {
		out = append(out, ifelse(def, f))
	}
	return out
}

// return an if else block for the definition level
func ifelse(def int, f fields.Field) ifElses {
	opts := optionals(def, f)
	var out ifElses
	for i, o := range opts {
		p := f.Parent(o + 1)
		if i == 0 {
			n := strings.Join(p.FieldNames, ".")
			cond := fmt.Sprintf("x.%s == nil", n)
			out.If.Cond = cond
			out.If.Val = fmt.Sprintf("x.%s = %s", n, removeAssigment(structs.Init(def, 0, 0, f)))
			if len(opts) == 1 && def == f.MaxDef() {
				ch := f.Child(len(f.FieldNames) - 1)
				out.Else = &ifElse{
					Val: fmt.Sprintf("x.%s = %s", strings.Join(f.FieldNames, "."), removeAssigment(structs.Init(def, 0, 0, ch))),
				}
			}
		} else if i+1 == f.MaxDef() {
			ch := f.Child(len(f.FieldNames) - 1)
			out.Else = &ifElse{
				Val: fmt.Sprintf("x.%s = %s", strings.Join(f.FieldNames, "."), removeAssigment(structs.Init(def, 0, 0, ch))),
			}
		} else {
			n := strings.Join(p.FieldNames, ".")
			cond := fmt.Sprintf("x.%s == nil", n)
			ch := f.Child(o)
			out.ElseIf = append(out.ElseIf, ifElse{
				Cond: cond,
				Val:  fmt.Sprintf("x.%s = %s", n, removeAssigment(structs.Init(def-i, 0, 0, ch))),
			})
		}
	}

	return out
}

func removeAssigment(s string) string {
	return strings.Replace(s[strings.Index(s, "= ")+1:], "nVals", "0", 1)
}

func optionals(def int, f fields.Field) []int {
	var out []int
	for i, rt := range f.RepetitionTypes {
		if rt == fields.Optional {
			out = append(out, i)
			if len(out) == def {
				return out
			}
		}
	}
	return out
}
