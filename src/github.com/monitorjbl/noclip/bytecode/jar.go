package bytecode

import (
	"archive/zip"
	"strings"
	"fmt"
	"io"
	"encoding/binary"
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
	class.CanonicalName = classNameFromPath(classfile.Name)
	class.Size = classfile.FileInfo().Size()

	fmt.Printf("Loading class %v\n", class.CanonicalName)
	//load magic value
	magic := read32(reader)
	if magic != MAGIC_NUMBER {
		fmt.Print("Class %v does not look like a normal class file", class.CanonicalName)
	}

	//load version fields
	class.MinorVersion = read16(reader)
	class.MajorVersion = read16(reader)

	//load constant pool
	class.ConstantPool = readConstantPool(reader, read16(reader))

	//load access flags
	index := read16(reader)

	//load reference to this class
	index = read16(reader)
	entry := class.ConstantPool[index - 1]
	if (entry.Type() != cp_class) {
		fmt.Printf("Class was malformed!")
	}
	thisEntry, _ := entry.(ConstantPool_Class)
	this, _ := class.ConstantPool[thisEntry.NameIndex - 1].(ConstantPool_UTF8)
	fmt.Printf("%v\n", this.Value)

	return class
}

func readConstantPool(reader io.ReadCloser, constantPoolSize uint16) ([]ConstantPoolEntry) {
	pool := make([]ConstantPoolEntry, 0)
	for i := 0; i < int(constantPoolSize - 1); i++ {
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
	if err != nil {
		fmt.Printf("Error! %v", err)
	}
	return content
}

func classNameFromPath(filename string) (string) {
	name := strings.Replace(filename, "/", ".", -1)
	return strings.Replace(name, ".class", "", -1)
}