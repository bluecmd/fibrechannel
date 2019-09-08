// v2 encoding, this will replace the reflect based one earlier
package encoding

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"sort"
)

type Size int
type TypeDef string
type Function string
type Statement string

const (
	Bits          = Size(1)
	Bytes         = Size(8)
	RemainingData = Size(-1)
)

var (
	Uint8  = &Unsigned{8 * Bits}
	Uint16 = &Unsigned{16 * Bits}
	Uint32 = &Unsigned{32 * Bits}
	Uint64 = &Unsigned{64 * Bits}
)

type Type interface {
	TypeName() string
	Consts() []NamedConstant
	TypeDefs() []TypeDef
	Functions() []Function
	Deser(p Type, m string) []Statement
	PreSer(p Type, m string) []Statement
	Ser(p Type, m string) []Statement
}

type Field struct {
	Name   string
	Type   Type
	Parent *Struct
}

type Struct struct {
	Name   string
	Fields []*Field
}

func (t *Struct) TypeName() string {
	return t.Name
}

func (t *Struct) Field(n string, ty Type) *Field {
	f := &Field{n, ty, t}
	t.Fields = append(t.Fields, f)
	return f
}

func (t *Struct) Consts() []NamedConstant {
	nc := []NamedConstant{}
	for _, c := range t.Fields {
		nc = append(nc, c.Type.Consts()...)
	}
	return nc
}

func (t *Struct) Functions() []Function {
	rf := fmt.Sprintf("func (o *%s) ReadFrom(r io.Reader) (int64, error) {\n", t.Name)
	rf += "_io := encoding.Reader{r}\n"
	for _, r := range t.Deser(t, "o") {
		rf += string(r) + "\n"
		rf += "if _io.Error != nil { return _io.Pos, _io.Error }\n"
	}
	rf += "return _io.Pos, nil }"

	wt := fmt.Sprintf("func (o *%s) WriteTo(w io.Writer) (int64, error) {\n", t.Name)
	wt += "_io := encoding.Writer{r}\n"
	for _, r := range t.PreSer(t, "o") {
		wt += string(r) + "\n"
	}
	for _, r := range t.Ser(t, "o") {
		wt += string(r) + "\n"
		wt += "if _io.Error != nil { return _io.Pos, _io.Error }\n"
	}
	wt += "return _io.Pos, nil }"

	return []Function{Function(rf), Function(wt)}
}

func (t *Struct) Deser(p Type, m string) []Statement {
	o := []Statement{}
	for _, f := range t.Fields {
		if f.Name == "" {
			o = append(o, f.Type.Deser(t, "")...)
		} else {
			o = append(o, f.Type.Deser(t, m+"."+f.Name)...)
		}
	}
	return o
}

func (t *Struct) TypeDefs() []TypeDef {
	td := []TypeDef{}
	mine := "type " + t.Name + " struct { "
	for _, f := range t.Fields {
		td = append(td, f.Type.TypeDefs()...)
		mine += fmt.Sprintf("%s %s\n", f.Name, f.Type.TypeName())
	}
	mine += " }"
	return append(td, TypeDef(mine))
}

func (t *Struct) PreSer(p Type, m string) []Statement {
	o := []Statement{}
	for _, f := range t.Fields {
		if f.Name == "" {
			o = append(o, f.Type.PreSer(t, "")...)
		} else {
			o = append(o, f.Type.PreSer(t, m+"."+f.Name)...)
		}
	}
	return o
}

func (t *Struct) Ser(p Type, m string) []Statement {
	o := []Statement{}
	for _, f := range t.Fields {
		if f.Name == "" {
			o = append(o, f.Type.Ser(t, "")...)
		} else {
			o = append(o, f.Type.Ser(t, m+"."+f.Name)...)
		}
	}
	return o
}

type NamedConstant struct {
	Name   string
	Domain string
	Constant
}

type Constant struct {
	Value   int
	Comment string
}

type Enum struct {
	Name   string
	Size   Size
	Values map[string]Constant
}

func (t *Enum) TypeName() string {
	return t.Name
}

func (t *Enum) Consts() []NamedConstant {
	n := []NamedConstant{}
	for k, v := range t.Values {
		n = append(n, NamedConstant{k, t.Name, v})
	}
	return n
}

func (t *Enum) Functions() []Function {
	return []Function{}
}

func (t *Enum) TypeDefs() []TypeDef {
	return []TypeDef{TypeDef(fmt.Sprintf("type %s uint%d", t.Name, t.Size))}
}

func (t *Enum) Deser(p Type, m string) []Statement {
	return []Statement{Statement(fmt.Sprintf("_io.ReadUint%d(&%s)", t.Size, m))}
}

func (t *Enum) PreSer(p Type, m string) []Statement {
	return []Statement{}
}

func (t *Enum) Ser(p Type, m string) []Statement {
	return []Statement{Statement(fmt.Sprintf("_io.WriteUint%d(%s)", t.Size, m))}
}

type SwitchedType struct {
	Name       string
	SwitchedOn *Field
	Cases      map[string]Type
}

func (t *SwitchedType) TypeName() string {
	return "interface{}"
}

func (t *SwitchedType) Consts() []NamedConstant {
	nc := []NamedConstant{}
	for _, c := range t.Cases {
		nc = append(nc, c.Consts()...)
	}
	return nc
}

func (t *SwitchedType) Functions() []Function {
	return []Function{}
}

func (t *SwitchedType) TypeDefs() []TypeDef {
	return []TypeDef{}
}

func (t *SwitchedType) Deser(p Type, m string) []Statement {
	if t.SwitchedOn.Parent != p {
		panic("Multi-level conditions not supported")
	}
	stmt := fmt.Sprintf("switch o.%s {\n", t.SwitchedOn.Name)
	for k, c := range t.Cases {
		stmt += fmt.Sprintf("case %s:\n", k)
		stmt += fmt.Sprintf("i := &%s{}\n", c.TypeName())
		stmt += "if n, err := i.ReadFrom(_io.NewReader()); err != nil { return n, err }\n"
		stmt += fmt.Sprintf("%s = i\n", m)
	}
	stmt += "}\n"

	return []Statement{Statement(stmt)}
}

func (t *SwitchedType) PreSer(p Type, m string) []Statement {
	if t.SwitchedOn.Parent != p {
		panic("Multi-level conditions not supported")
	}
	stmt := fmt.Sprintf("switch i := %s.(type) {\n", m)
	for k, c := range t.Cases {
		stmt += fmt.Sprintf("case %s:\n", c.TypeName())
		stmt += fmt.Sprintf("o.%s = %s\n", t.SwitchedOn.Name, k)
	}
	stmt += "}\n"

	return []Statement{Statement(stmt)}
}

func (t *SwitchedType) Ser(p Type, m string) []Statement {
	if t.SwitchedOn.Parent != p {
		panic("Multi-level conditions not supported")
	}
	stmt := fmt.Sprintf("switch i := %s.(type) {\n", m)
	for _, c := range t.Cases {
		stmt += fmt.Sprintf("case %s:\n", c.TypeName())
		stmt += "if n, err := i.WriteTo(_io.NewWriter()); err != nil { return n, err }\n"
	}
	stmt += "}\n"

	return []Statement{Statement(stmt)}
}

type Unsigned struct {
	Size Size
}

func (t *Unsigned) TypeName() string {
	return fmt.Sprintf("uint%d", t.bytes()*8)
}

func (t *Unsigned) Consts() []NamedConstant {
	return []NamedConstant{}
}

func (t *Unsigned) Functions() []Function {
	return []Function{}
}

func (t *Unsigned) TypeDefs() []TypeDef {
	return []TypeDef{}
}

func (t *Unsigned) bytes() int {
	odd := t.Size & 0x7
	if odd != 0 {
		return int(t.Size)/8 + 1
	}
	return int(t.Size) / 8
}

func (t *Unsigned) Deser(p Type, m string) []Statement {
	if m == "" {
		return []Statement{Statement(fmt.Sprintf("_io.Skip(%d)", t.bytes()))}
	}
	return []Statement{Statement(fmt.Sprintf("_io.ReadUint%d(&%s)", t.bytes()*8, m))}
}

func (t *Unsigned) PreSer(p Type, m string) []Statement {
	return []Statement{}
}

func (t *Unsigned) Ser(p Type, m string) []Statement {
	if m == "" {
		return []Statement{Statement(fmt.Sprintf("_io.Skip(%d)", t.bytes()))}
	}
	return []Statement{Statement(fmt.Sprintf("_io.WriteUint%d(%s)", t.bytes()*8, m))}
}

func NewStruct(n string) *Struct {
	return &Struct{
		Name: n,
	}
}

func Generate(pkg string, ts ...Type) ([]byte, error) {
	consts := []NamedConstant{}
	typedefs := []TypeDef{}
	funcs := []Function{}

	for _, t := range ts {
		consts = append(consts, t.Consts()...)
		typedefs = append(typedefs, t.TypeDefs()...)
		funcs = append(funcs, t.Functions()...)
	}

	sort.Slice(consts, func(i, j int) bool {
		if consts[i].Domain != consts[j].Domain {
			return consts[i].Domain < consts[j].Domain
		}
		return consts[i].Value < consts[j].Value
	})

	doc := new(bytes.Buffer)
	bw := bufio.NewWriter(doc)
	bw.WriteString("// Generated by Fibre Channel protocol generator\n")
	bw.WriteString("// Any manual changes will be lost\n")
	bw.WriteString("\n")
	bw.WriteString(fmt.Sprintf("package %s", pkg))
	bw.WriteString("\n")
	bw.WriteString("const (\n")

	d := ""
	for _, c := range consts {
		if c.Domain != d && d != "" {
			bw.WriteString("\n")
		}
		d = c.Domain
		bw.WriteString(fmt.Sprintf("\t%s = 0x%x // %s\n", c.Name, c.Value, c.Comment))
	}
	bw.WriteString(")\n")

	for _, td := range typedefs {
		bw.WriteString(fmt.Sprintf("%s\n\n", td))
	}

	for _, f := range funcs {
		bw.WriteString(fmt.Sprintf("%s\n\n", f))
	}

	bw.Flush()

	fmt.Print(string(doc.Bytes()))
	return format.Source(doc.Bytes())
}
