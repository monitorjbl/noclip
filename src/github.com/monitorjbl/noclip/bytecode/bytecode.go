package bytecode

import (
	"io"
	"fmt"
)

type LocalVariableTableEntry struct {
	start_pc         uint16;
	length           uint16;
	name_index       uint16;
	descriptor_index uint16;
	index            uint16;
}

func readAttributes(reader io.ReadCloser, class *ClassFile, cp *[]ConstantPoolEntry, method string) ([]string) {
	attrCount := read16(reader)
	attrs := make([]string, 0)
	//fmt.Printf("\tAttrs: %v\n", attrCount)
	for i := 0; i < int(attrCount); i++ {
		attrs = append(attrs, readAttribute(reader, class, cp, method))
	}
	return attrs
}

func readAttribute(reader io.ReadCloser, class *ClassFile, cp *[]ConstantPoolEntry, method string) (string) {
	attrName := lookupUTF8(class, cp, read16(reader))
	attrLength := read32(reader)
	//fmt.Printf("\tAttr %v: %v\n", attrName, attrLength)
	switch attrName {
	case "Code":
		readCodeAttribute(reader, class, cp, method)
		break
	case "LineNumberTable":
		readLineNumberTableAttribute(reader)
		break
	case "LocalVariableTable":
		readLocalVariableTableAttribute(reader, class, cp)
		break
	case "ConstantValue":
		readConstantValueAttribute(reader)
		break
	case "Synthetic":
		break
	case "SourceFile":
		readSourceFileAttribute(reader)
		break
	case "InnerClasses":
		readInnerClassesAttribute(reader)
		break
	case "Deprecated":
		break
	case "Exceptions":
		readExceptionsAttribute(reader)
		break
	default:
		fmt.Printf("Not sure what to do with %v, just reading out %v bytes\n", attrName, attrLength)
		readSimple32(reader, attrLength)
	}
	return attrName
}

func readSourceFileAttribute(reader io.ReadCloser) {
	//sourcefile_index
	read16(reader)
}

func readConstantValueAttribute(reader io.ReadCloser) {
	//constantvalue_index
	read16(reader)
}

func readLineNumberTableAttribute(reader io.ReadCloser) {
	//line_number_table_length
	line_number_length := read16(reader)
	//read out all line_number_table entries
	readSimple(reader, line_number_length * 4)
}

func readInnerClassesAttribute(reader io.ReadCloser) {
	//number_of_classes
	number_of_classes := read16(reader)
	//read out all classes
	readSimple(reader, number_of_classes * 4)
}

func readExceptionsAttribute(reader io.ReadCloser) {
	//number_of_exceptions
	number_of_exceptions := read16(reader)
	//read out all exceptions
	readSimple(reader, number_of_exceptions * 2)
}

func readLocalVariableTableAttribute(reader io.ReadCloser, class *ClassFile, cp *[]ConstantPoolEntry) {
	//line_number_table_length
	local_variable_length := int(read16(reader))
	//read out all line_number_table entries
	for i := 0; i < local_variable_length; i++ {
		entry := new(LocalVariableTableEntry)
		entry.start_pc = read16(reader)
		entry.length = read16(reader)
		entry.name_index = read16(reader)
		entry.descriptor_index = read16(reader)
		entry.index = read16(reader)

		//name := lookupUTF8(class, cp, entry.name_index)
		//varType := lookupUTF8(class, cp, entry.descriptor_index)
		//fmt.Printf("\t\t\tVariable %v [%v]\n", name, varType)
	}
}

func readCodeAttribute(reader io.ReadCloser, class *ClassFile, cp *[]ConstantPoolEntry, method string) {
	//max_stack
	read16(reader)
	//max_locals
	read16(reader)
	//code_length
	code_length := read32(reader)
	//fmt.Printf("\t\tCode: %v\n", code_length)

	//read out all bytecode
	readSimple32(reader, code_length)

	//exception_table_length
	exception_table_length := read16(reader)
	//fmt.Printf("\t\tExceptionTable: %v\n", exception_table_length*8)

	//read out all exception handler info
	readSimple(reader, exception_table_length * 8)

	//read out all sub attributes
	readAttributes(reader, class, cp, method)
}
