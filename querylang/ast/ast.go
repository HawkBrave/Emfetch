package ast

type RootNodeType int
type Value interface{}

const (
	ObjectRoot RootNodeType = iota
	ArrayRoot
)

type RootNode struct {
	RootValue *Value
	Type      RootNodeType
}

type Identifier struct {
	Type  string
	Value string
}

type Property struct {
	Type  string
	Key   Identifier
	Value Value
}

type Object struct {
	Type     string
	Children []Property
	Start    int
	End      int
}

type Array struct {
	Type     string
	Children []Value
	Start    int
	End      int
}

type Literal struct {
	Type  string
	Value Value
}

type state int

const (
	// String states
	StringStart state = iota
	StringQuoteOrChar
	Escape
	// Number states
	NumberStart
	NumberMinus
	NumberZero
	NumberDigit
	NumberPoint
	NumberDigitFraction
	NumberExp
	NumberExpDigitOrSign
	// Property states
	PropertyStart
	PropertyKey
	PropertyColon
	// Object states
	ObjStart
	ObjOpen
	ObjProperty
	ObjComma
	// Array states
	ArrayStart
	ArrayOpen
	ArrayValue
	ArrayComma
)
