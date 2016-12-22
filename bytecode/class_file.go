package bytecode

import "fmt"

type ClassFile struct {
	FileName        string
	CanonicalName   string
	Size            int64

	MinorVersion    uint16
	MajorVersion    uint16

	Superclass      string
	Interfaces      []string
	Fields          []ClassField
	Methods         []ClassMethod

	ClassReferences []string
}

type ClassField struct {
	Name       string
	Type       string
	Attributes []string
}

type ClassMethod struct {
	Name        string
	Description string
	Attributes  []string
}

func (c ClassFile) ToString() (string) {
	return fmt.Sprintf("%v", c.CanonicalName)
}

