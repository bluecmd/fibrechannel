package encoding

import (
	"fmt"

	"strings"
)

type bitfield struct {
	// 32 j
	name   string
	isbool bool
	masks  []int
	bytes  []int
	shifts []uint
	fshift uint
}

type BitStruct struct {
	Name   string
	offset int
	fields []*bitfield
}

func (t *BitStruct) BoolBit(n string) *bitfield {
	bit := uint(7 - t.offset%8)
	mask := 1 << bit
	shift := bit
	b := t.offset / 8
	bf := &bitfield{n, true, []int{mask}, []int{b}, []uint{0}, shift}
	t.fields = append(t.fields, bf)
	t.offset++
	return bf
}

func (t *BitStruct) SkipBit(n int) {
	t.offset += n
}

func (t *BitStruct) IntField(n string, bits int) *bitfield {
	bs := []int{}
	masks := []int{}
	shifts := []uint{}

	shift := uint((bits / 8) * 8)
	if bits%8 == 0 {
		shift -= 8
	}
	mask := 0
	b := t.offset / 8
	for i := t.offset; i < t.offset+bits; i++ {
		if i%8 == 0 && mask > 0 {
			bs = append(bs, b)
			masks = append(masks, mask)
			shifts = append(shifts, shift)
			mask = 0
			shift -= 8
			b++
		}
		bit := uint(7 - i%8)
		mask |= 1 << bit
	}

	if mask > 0 {
		bs = append(bs, b)
		masks = append(masks, mask)
		shifts = append(shifts, shift)
	}

	fshift := uint(7 - (t.offset+bits-1)%8)
	bf := &bitfield{n, false, masks, bs, shifts, fshift}
	t.fields = append(t.fields, bf)
	t.offset += bits
	return bf
}

func (t *BitStruct) TypeName() string {
	return t.Name
}

func (t *BitStruct) Consts() []NamedConstant {
	return []NamedConstant{}
}

func (t *BitStruct) Functions() []Function {
	return []Function{}
}

func (t *BitStruct) TypeDefs() []TypeDef {
	td := []TypeDef{}
	mine := "type " + t.Name + " struct { "
	for _, f := range t.fields {
		tn := "int"
		if f.isbool {
			tn = "bool"
		}
		mine += fmt.Sprintf("%s %s\n", f.name, tn)
	}
	mine += " }"
	return append(td, TypeDef(mine))
}

func (t *BitStruct) Deser(p Context, m string) ([]Statement, error) {
	stmt := fmt.Sprintf("{ var bs [%d]byte\n", t.offset/8)
	stmt += " _io.Read(bs[:])\n if _io.Error != nil { return _io.Pos, _io.Error }\n"
	for _, f := range t.fields {
		expr := "0"
		for i, _ := range f.masks {
			expr += fmt.Sprintf("| int(bs[%d] & 0x%x) << %d", f.bytes[i], f.masks[i], f.shifts[i])
		}
		if f.isbool {
			stmt += fmt.Sprintf(" %s.%s = (%s) == 0x%x\n", m, f.name, expr, 1<<f.fshift)
		} else {
			stmt += fmt.Sprintf(" %s.%s = ((%s) >> %d)\n", m, f.name, expr, f.fshift)
		}
	}

	stmt += "}"

	stmt = strings.ReplaceAll(stmt, "<< 0", "")
	stmt = strings.ReplaceAll(stmt, ">> 0", "")
	stmt = strings.ReplaceAll(stmt, " 0|", "")

	return []Statement{Statement(stmt)}, nil
}

func (t *BitStruct) PreSer(p Context, m string) ([]Statement, error) {
	return []Statement{}, nil
}

func (t *BitStruct) Ser(p Context, m string) ([]Statement, error) {
	stmt := fmt.Sprintf("{ var bs [%d]byte\n", t.offset/8)
	stmt += "bool2int := func(v bool) int { if v { return 1 } \n return 0 }\n"
	bytes := t.offset / 8
	if t.offset%8 > 0 {
		bytes++
	}
	for i := 0; i < bytes; i++ {
		expr := "0"
		for _, f := range t.fields {
			for j, b := range f.bytes {
				if b != i {
					continue
				}
				conv := "int"
				if f.isbool {
					conv = "bool2int"
				}
				expr += fmt.Sprintf("| ((%s(%s.%s) << %d) >> %d) & 0x%x", conv, m, f.name, f.fshift, f.shifts[j], f.masks[j])
			}
		}
		stmt += fmt.Sprintf("bs[%d] = byte(%s)\n", i, expr)
	}

	stmt += " _io.Write(bs[:])\n"
	stmt += "}"

	stmt = strings.ReplaceAll(stmt, "<< 0", "")
	stmt = strings.ReplaceAll(stmt, ">> 0", "")
	stmt = strings.ReplaceAll(stmt, " 0|", "")
	return []Statement{Statement(stmt)}, nil
}

func (t *BitStruct) FindReference(needle interface{}) string {
	for _, v := range t.fields {
		if v == needle {
			return v.name
		}
	}
	return ""
}

func NewBitStruct(n string) *BitStruct {
	return &BitStruct{
		Name: n,
	}
}
