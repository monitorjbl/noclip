package bytecode

import (
	"archive/zip"
	"strings"
	"fmt"
	"io"
	"encoding/binary"
	"log"
)

const MAGIC_NUMBER uint32 = 0xCAFEBABE;

func LoadJar(filename string) ([]*ClassFile) {
	var closer, err = zip.OpenReader(filename)
	if err != nil {
		fmt.Printf("Error! %v", err)
	}

	var files = closer.File
	return loadClasses(&files)
}

func loadClasses(files *[]*zip.File) ([]*ClassFile) {
	arr := make([]*ClassFile, 0)
	for _, f := range *files {
		if (!f.FileInfo().IsDir() && strings.HasSuffix(f.Name, ".class")) {
			arr = append(arr, loadClassFile(f))
		}
	}
	return arr
}

func loadClassFile(classfile *zip.File) (*ClassFile) {
	reader, err := classfile.Open()
	if err != nil {
		fmt.Printf("Error! %v", err)
	}

	//generate class file
	class := new(ClassFile)
	class.FileName = classfile.Name
	class.Size = classfile.FileInfo().Size()

	fmt.Printf("Loading class %v\n", class.FileName)
	//load magic value
	magic := read32(reader)
	if magic != MAGIC_NUMBER {
		log.Fatal(fmt.Sprintf("Class %v does not look like a normal class file", class.CanonicalName))
	}

	//load version fields
	class.MinorVersion = read16(reader)
	class.MajorVersion = read16(reader)

	//load constant pool
	constantPool := readConstantPool(reader)

	//load access flags
	read16(reader)

	//load reference to this class
	class.CanonicalName = classNameFromPath(lookupClassEntry(class, &constantPool, read16(reader)))

	//load reference to superclass
	class.Superclass = classNameFromPath(lookupClassEntry(class, &constantPool, read16(reader)))

	//load interfaces
	class.Interfaces = readInterfaces(reader, class, &constantPool)

	//load fields
	class.Fields = readFields(reader, class, &constantPool)

	//load methods
	class.Methods = readMethods(reader, class, &constantPool)

	//read attributes
	readAttributes(reader, class, &constantPool)

	return class
}

func readConstantPool(reader io.ReadCloser) ([]ConstantPoolEntry) {
	poolSize := read16(reader)
	pool := make([]ConstantPoolEntry, 0)
	for i := 0; i < int(poolSize - 1); i++ {
		pool = append(pool, readConstantPoolEntry(reader, read8(reader)))
	}
	return pool
}

func readConstantPoolEntry(reader io.ReadCloser, tag uint8) (ConstantPoolEntry) {
	switch tag {
	case cp_utf8:
		length := read16(reader)
		str := readSimple(reader, length)
		return ConstantPool_UTF8{Value:string(str)}
	case cp_integer:
		return ConstantPool_Integer{Value:read32(reader)}
	case cp_float:
		return ConstantPool_Float{Value:read32(reader)}
	case cp_long:
		//TODO: deal with the fact that 8-byte constants take up two entries
		return ConstantPool_Long{High:read32(reader), Low:read32(reader)}
	case cp_double:
		//TODO: deal with the fact that 8-byte constants take up two entries
		return ConstantPool_Double{High:read32(reader), Low:read32(reader)}
	case cp_class:
		return ConstantPool_Class{NameIndex:read16(reader)}
	case cp_string:
		return ConstantPool_String{StringIndex:read16(reader)}
	case cp_field_ref:
		return ConstantPool_FieldRef{ClassIndex:read16(reader), NameAndTypeIndex:read16(reader)}
	case cp_method_ref:
		return ConstantPool_MethodRef{ClassIndex:read16(reader), NameAndTypeIndex:read16(reader)}
	case cp_interface_method_ref:
		return ConstantPool_InterfaceMethodRef{ClassIndex:read16(reader), NameAndTypeIndex:read16(reader)}
	case cp_name_and_type:
		return ConstantPool_NameAndType{NameIndex:read16(reader), DescriptorIndex:read16(reader)}
	case cp_method_handle:
		return ConstantPool_MethodHandle{Kind:read8(reader), Index:read16(reader)}
	case cp_indy:
		return ConstantPool_MethodType{DescriptorIndex:read16(reader)}
	case cp_method_type:
		return ConstantPool_Indy{BootstrapMethodIndex:read16(reader), NameAndTypeIndex:read16(reader)}
	default:
		return nil
	}

}

func readInterfaces(reader io.ReadCloser, class *ClassFile, cp *[]ConstantPoolEntry) ([]string) {
	interfaceCount := read16(reader)
	interfaces := make([]string, 0)
	for i := 0; i < int(interfaceCount); i++ {
		interfaces = append(interfaces, lookupClassEntry(class, cp, read16(reader)))
	}
	return interfaces
}
func readFields(reader io.ReadCloser, class *ClassFile, cp *[]ConstantPoolEntry) ([]ClassField) {
	fieldCount := read16(reader)
	fields := make([]ClassField, 0)
	for i := 0; i < int(fieldCount); i++ {
		field := ClassField{}
		//read out access values for the field
		read16(reader)
		field.Name = lookupUTF8(class, cp, read16(reader))
		field.Type = lookupUTF8(class, cp, read16(reader))
		readAttributes(reader, class, cp)
		fields = append(fields, field)
	}
	return fields
}

func readMethods(reader io.ReadCloser, class *ClassFile, cp *[]ConstantPoolEntry) ([]ClassMethod) {
	methodCount := read16(reader)
	methods := make([]ClassMethod, 0)
	for i := 0; i < int(methodCount); i++ {
		method := ClassMethod{}
		//read out access values for the field
		read16(reader)
		method.Name = lookupUTF8(class, cp, read16(reader))
		method.ReturnType = lookupUTF8(class, cp, read16(reader))
		readAttributes(reader, class, cp)
		methods = append(methods, method)
	}
	return methods
}

func readAttributes(reader io.ReadCloser, class *ClassFile, cp *[]ConstantPoolEntry) ([]string) {
	attrCount := read16(reader)
	attrs := make([]string, 0)
	for i := 0; i < int(attrCount); i++ {
		attrs = append(attrs, lookupUTF8(class, cp, read16(reader)))

		//read out the parts that we don't care about
		attrLength := read32(reader)
		readSimple32(reader, attrLength)
	}
	return attrs
}

func lookupClassEntry(class *ClassFile, cp *[]ConstantPoolEntry, index uint16) (string) {
	pool := *cp
	entry := pool[index - 1]
	if (entry.Type() != cp_class) {
		log.Fatal(fmt.Sprintf("Class %v was malformed!", class.CanonicalName))
	}
	thisEntry, _ := entry.(ConstantPool_Class)
	this, _ := pool[thisEntry.NameIndex - 1].(ConstantPool_UTF8)
	return this.Value
}

func lookupUTF8(class *ClassFile, cp *[]ConstantPoolEntry, index uint16) (string) {
	pool := *cp
	entry := pool[index - 1]
	if (entry.Type() != cp_utf8) {
		log.Fatal(fmt.Sprintf("Class %v was malformed!", class.CanonicalName))
	}
	this, _ := entry.(ConstantPool_UTF8)
	return this.Value
}

func read8(reader io.ReadCloser) (uint8) {
	return uint8(readSimple(reader, 1)[0])
}

func read16(reader io.ReadCloser) (uint16) {
	return binary.BigEndian.Uint16(readSimple(reader, 2))
}

func read32(reader io.ReadCloser) (uint32) {
	return binary.BigEndian.Uint32(readSimple(reader, 4))
}

func readSimple(reader io.ReadCloser, length uint16) ([]byte) {
	content := make([]byte, length)
	_, err := reader.Read(content)
	if err != nil && err != io.EOF {
		log.Fatal(fmt.Sprintf("Could not read class file, got error %v", err))
	}
	return content
}

func readSimple32(reader io.ReadCloser, length uint32) ([]byte) {
	content := make([]byte, length)
	_, err := reader.Read(content)
	if err != nil && err != io.EOF {
		log.Fatal(fmt.Sprintf("Could not read class file, got error %v", err))
	}
	return content
}

func classNameFromPath(filename string) (string) {
	name := strings.Replace(filename, "/", ".", -1)
	return strings.Replace(name, ".class", "", -1)
}