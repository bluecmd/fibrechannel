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
	rf += "_io := encoding.Reader{R: r}\n"
	for _, r := range t.Deser(t, "o") {
		rf += string(r) + "\n"
		rf += "if _io.Error != nil { return _io.Pos, _io.Error }\n"
	}
	rf += "return _io.Pos, nil }"

	wt := fmt.Sprintf("func (o *%s) WriteTo(w io.Writer) (int64, error) {\n", t.Name)
	wt += "_io := encoding.Writer{W: w}\n"
	for _, r := range t.PreSer(t, "o") {
		wt += string(r) + "\n"
	}
	for _, r := range t.Ser(t, "o") {
		wt += string(r) + "\n"
		wt += "if _io.Error != nil { return _io.Pos, _io.Error }\n"
	}
	wt += "return _io.Pos, nil }"

	fcts := []Function{Function(rf), Function(wt)}
	for _, f := range t.Fields {
		fcts = append(fcts, f.Type.Functions()...)
	}
	return fcts
}

func (t *Struct) Deser(p Type, m string) []Statement {
	o := []Statement{}
	for _, f := range t.Fields {
		o = append(o, f.Type.Deser(t, m+"."+f.Name)...)
	}
	return o
}

func (t *Struct) TypeDefs() []TypeDef {
	td := []TypeDef{}
	mine := "type " + t.Name + " struct { "
	for _, f := range t.Fields {
		td = append(td, f.Type.TypeDefs()...)
		tn := f.Type.TypeName()
		if tn != "" {
			mine += fmt.Sprintf("%s %s\n", f.Name, tn)
		}
	}
	mine += " }"
	return append(td, TypeDef(mine))
}

func (t *Struct) PreSer(p Type, m string) []Statement {
	o := []Statement{}
	for _, f := range t.Fields {
		o = append(o, f.Type.PreSer(t, m+"."+f.Name)...)
	}
	return o
}

func (t *Struct) Ser(p Type, m string) []Statement {
	o := []Statement{}
	for _, f := range t.Fields {
		o = append(o, f.Type.Ser(t, m+"."+f.Name)...)
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
	sort.Slice(n, func(i, j int) bool { return n[i].Value < n[j].Value })
	return n
}

func (t *Enum) Functions() []Function {
	str := fmt.Sprintf("func (o *%s) String() string { switch *o {\n", t.Name)
	for _, v := range t.Consts() {
		str += fmt.Sprintf(" case 0x%x:\n", v.Value)
		str += fmt.Sprintf("   return \"%s <0x%x> (%s)\"\n", v.Name, v.Value, v.Comment)
	}
	str += "default:\n\treturn fmt.Sprintf(\"--Invalid Enum Value-- <0x%x>\", *o)\n}}"
	return []Function{Function(str)}
}

func (t *Enum) TypeDefs() []TypeDef {
	return []TypeDef{TypeDef(fmt.Sprintf("type %s uint%d", t.Name, t.Size))}
}

func (t *Enum) Deser(p Type, m string) []Statement {
	return []Statement{Statement(fmt.Sprintf("_io.ReadObject(&%s)", m))}
}

func (t *Enum) PreSer(p Type, m string) []Statement {
	return []Statement{}
}

func (t *Enum) Ser(p Type, m string) []Statement {
	return []Statement{Statement(fmt.Sprintf("_io.WriteObject(%s)", m))}
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
		stmt += "if n, err := i.ReadFrom(&_io); err != nil { return n, err }\n"
		stmt += fmt.Sprintf("%s = i\n", m)
	}
	stmt += "}\n"

	return []Statement{Statement(stmt)}
}

func (t *SwitchedType) PreSer(p Type, m string) []Statement {
	if t.SwitchedOn.Parent != p {
		panic("Multi-level conditions not supported")
	}
	stmt := fmt.Sprintf("switch %s.(type) {\n", m)
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
		stmt += fmt.Sprintf("case *%s:\n", c.TypeName())
		stmt += "if n, err := i.WriteTo(&_io); err != nil { return n, err }\n"
	}
	stmt += "default:\n"
	stmt += "  return _io.Pos, fmt.Errorf(\"Unsupported type %v\", i)\n"
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
	return []Statement{Statement(fmt.Sprintf("_io.ReadObject(&%s)", m))}
}

func (t *Unsigned) PreSer(p Type, m string) []Statement {
	return []Statement{}
}

func (t *Unsigned) Ser(p Type, m string) []Statement {
	return []Statement{Statement(fmt.Sprintf("_io.WriteObject(%s)", m))}
}

type ByteArray struct {
	Count int
}

func (t *ByteArray) TypeName() string {
	return fmt.Sprintf("[%d]byte", t.Count)
}

func (t *ByteArray) Consts() []NamedConstant {
	return []NamedConstant{}
}

func (t *ByteArray) Functions() []Function {
	return []Function{}
}

func (t *ByteArray) TypeDefs() []TypeDef {
	return []TypeDef{}
}

func (t *ByteArray) Deser(p Type, m string) []Statement {
	return []Statement{Statement(fmt.Sprintf("_io.ReadObject(&%s)", m))}
}

func (t *ByteArray) PreSer(p Type, m string) []Statement {
	return []Statement{}
}

func (t *ByteArray) Ser(p Type, m string) []Statement {
	return []Statement{Statement(fmt.Sprintf("_io.WriteObject(%s)", m))}
}

type Object struct {
	Class string
}

func (t *Object) TypeName() string {
	return t.Class
}

func (t *Object) Consts() []NamedConstant {
	return []NamedConstant{}
}

func (t *Object) Functions() []Function {
	return []Function{}
}

func (t *Object) TypeDefs() []TypeDef {
	return []TypeDef{}
}

func (t *Object) Deser(p Type, m string) []Statement {
	return []Statement{
		Statement(fmt.Sprintf("if n, err := %s.ReadFrom(&_io); err != nil { return n, err }", m))}
}

func (t *Object) PreSer(p Type, m string) []Statement {
	return []Statement{}
}

func (t *Object) Ser(p Type, m string) []Statement {
	return []Statement{
		Statement(fmt.Sprintf("if n, err := %s.WriteTo(&_io); err != nil { return n, err }", m))}
}

type Skip struct {
	Size Size
}

func (t *Skip) TypeName() string {
	return ""
}

func (t *Skip) Consts() []NamedConstant {
	return []NamedConstant{}
}

func (t *Skip) Functions() []Function {
	return []Function{}
}

func (t *Skip) TypeDefs() []TypeDef {
	return []TypeDef{}
}

func (t *Skip) bytes() int {
	odd := t.Size & 0x7
	if odd != 0 {
		return int(t.Size)/8 + 1
	}
	return int(t.Size) / 8
}

func (t *Skip) Deser(p Type, m string) []Statement {
	return []Statement{Statement(fmt.Sprintf("_io.Skip(%d)", t.bytes()))}
}

func (t *Skip) PreSer(p Type, m string) []Statement {
	return []Statement{}
}

func (t *Skip) Ser(p Type, m string) []Statement {
	return []Statement{Statement(fmt.Sprintf("_io.Skip(%d)", t.bytes()))}
}

type Array struct {
	Count int
	Type Type
}

func (t *Array) TypeName() string {
	return fmt.Sprintf("[%d]%s", t.Count, t.Type.TypeName())
}

func (t *Array) Consts() []NamedConstant {
	return []NamedConstant{}
}

func (t *Array) Functions() []Function {
	return []Function{}
}

func (t *Array) TypeDefs() []TypeDef {
	return []TypeDef{}
}

func (t *Array) Deser(p Type, m string) []Statement {
	stmt := []Statement{}
	for i := 0; i < t.Count; i++ {
		stmt = append(stmt, t.Type.Deser(t, m + fmt.Sprintf("[%d]", i))...)
	}
	return stmt
}

func (t *Array) PreSer(p Type, m string) []Statement {
	stmt := []Statement{}
	for i := 0; i < t.Count; i++ {
		stmt = append(stmt, t.Type.PreSer(t, m + fmt.Sprintf("[%d]", i))...)
	}
	return stmt
}

func (t *Array) Ser(p Type, m string) []Statement {
	stmt := []Statement{}
	for i := 0; i < t.Count; i++ {
		stmt = append(stmt, t.Type.Ser(t, m + fmt.Sprintf("[%d]", i))...)
	}
	return stmt
}

func NewStruct(n string) *Struct {
	return &Struct{
		Name: n,
	}
}

func Generate(pkg string, imports []string, ts ...Type) ([]byte, error) {
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
	sort.Slice(typedefs, func(i, j int) bool { return typedefs[i] < typedefs[j] })
	sort.Slice(funcs, func(i, j int) bool { return funcs[i] < funcs[j] })

	doc := new(bytes.Buffer)
	bw := bufio.NewWriter(doc)
	bw.WriteString("// Generated by Fibre Channel protocol generator\n")
	bw.WriteString("// Any manual changes will be lost\n")
	bw.WriteString("\n")
	bw.WriteString(fmt.Sprintf("package %s", pkg))
	bw.WriteString("\n")
	bw.WriteString("import (\n")
	bw.WriteString("\"io\"\n")
	bw.WriteString("\"fmt\"\n")
	bw.WriteString("\n")
	bw.WriteString("\"github.com/bluecmd/fibrechannel/encoding\"\n")
	for _, i := range imports {
		bw.WriteString(fmt.Sprintf("\"%s\"\n", i))
	}
	bw.WriteString(")\n")
	bw.WriteString("\n")
	bw.WriteString("const (\n")

	d := ""
	p := ""
	for _, c := range consts {
		if c.Domain != d && d != "" {
			bw.WriteString("\n")
		}
		d = c.Domain
		if p == c.Name {
			continue
		}
		p = c.Name
		bw.WriteString(fmt.Sprintf("\t%s = 0x%x // %s\n", c.Name, c.Value, c.Comment))
	}
	bw.WriteString(")\n")

	for _, td := range typedefs {
		if p == string(td) {
			continue
		}
		p = string(td)
		bw.WriteString(fmt.Sprintf("%s\n\n", td))
	}

	for _, f := range funcs {
		if p == string(f) {
			continue
		}
		p = string(f)
		bw.WriteString(fmt.Sprintf("%s\n\n", f))
	}

	bw.Flush()

	out, err := format.Source(doc.Bytes())
	if err != nil {
		return doc.Bytes(), err
	}
	return out, nil
}
