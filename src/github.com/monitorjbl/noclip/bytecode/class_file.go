package bytecode

import "fmt"

type ClassFile struct {
	FileName      string
	CanonicalName string
	Size          int64

	MinorVersion  uint16
	MajorVersion  uint16

	Superclass    string
	Interfaces    []string
	Fields        []ClassField
	Methods       []ClassMethod
}

type ClassField struct {
	Name       string
	Type       string
	Attributes []string
}

type ClassMethod struct {
	Name       string
	ReturnType string
	Parameters []string
}

func (c ClassFile) ToString() (string) {
	return fmt.Sprintf("%v, %v", c.CanonicalName, c.Methods)
}

