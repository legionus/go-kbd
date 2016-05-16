package lexer

import (
	"encoding/json"
)

//go:generate stringer -type=Kind

// Kind allows to distinguish between different types of leaves.
type Kind int

// Leaf kinds.
const (
	Unknown Kind = iota
	Space
	NewLine
	Comment
	Equals
	Dash
	Comma
	Plus
	Char
	Number
	Unicode
	Literal
	Escaped
	Quote
	Strval
	AltIsMeta
	Strings
	String
	Charset
	Keymaps
	Keycode
	Plain
	CapsShift
	Compose
	Control
	CtrlL
	CtrlR
	AltGr
	Alt
	ShiftL
	ShiftR
	Shift
	Usual
	For
	As
	On
	To
	Include
)

// Leaf is a basic node type. Represents a piece of the input data.
type Leaf struct {
	Kind Kind
	Data []byte
//	pos  Pos
}

//func (l Leaf) Pos() Pos {
//	return l.pos
//}

//func (l Leaf) End() Pos {
//	return l.pos + Pos(len(l.Data))
//}

func (l Leaf) MarshalText() ([]byte, error) {
	return l.Data, nil
}

func (l Leaf) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string
		Kind string
		Data string
	}{
		"Leaf",
		l.Kind.String(),
		string(l.Data),
	})
}
