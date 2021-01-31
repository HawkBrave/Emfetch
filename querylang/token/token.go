package token

import "fmt"

type Type string

const (
	// Token/character we don't know about
	Illegal Type = "ILLEGAL"

	// End of file
	EOF Type = "EOF"

	// Literals
	String Type = "STRING"
	Number Type = "NUMBER"

	// The six structural tokens
	LeftBrace    Type = "{"
	RightBrace   Type = "}"
	LeftBracket  Type = "["
	RightBracket Type = "]"
	Comma        Type = ","
	Colon        Type = ":"
	Loop         Type = "*"

	// Values
	True  Type = "TRUE"
	False Type = "FALSE"
	Null  Type = "NULL"
)

type Token struct {
	Type    Type
	Literal string
	Line    int
	Start   int
	End     int
}

var validIdentifiers = map[string]Type{
	"true":  True,
	"false": False,
	"null":  Null,
}

func LookupIdentifier(identifier string) (Type, error) {
	if token, ok := validIdentifiers[identifier]; ok {
		return token, nil
	}
	return "", fmt.Errorf("Expected a valid JSON identifier. Found: %s", identifier)
}
