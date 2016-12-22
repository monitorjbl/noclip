package bytecode

import "fmt"

type ClassFile struct {
	CanonicalName string
	Size          int64

	MinorVersion  uint16
	MajorVersion  uint16
	ConstantPool  []ConstantPoolEntry

	Superclass string
}


func (c ClassFile) ToString() (string) {
	return fmt.Sprintf("%v [%v.%v]: %v", c.CanonicalName, c.MajorVersion, c.MinorVersion, c.Size)
}

