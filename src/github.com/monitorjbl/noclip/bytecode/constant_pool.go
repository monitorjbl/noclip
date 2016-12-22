package bytecode

//All struct definitions come from here: https://docs.oracle.com/javase/specs/jvms/se7/html/jvms-4.html

const (
	cp_utf8 = 1
	cp_integer = 3
	cp_float = 4
	cp_long = 5
	cp_double = 6
	cp_class = 7
	cp_string = 8
	cp_field_ref = 9
	cp_method_ref = 10
	cp_interface_method_ref = 11
	cp_name_and_type = 12
	cp_method_handle = 15
	cp_indy = 16
	cp_method_type = 18
)

type ConstantPoolEntry interface {
	Type() uint8
}

type ConstantPool_Class struct {
	ConstantPoolEntry
	NameIndex uint16
}

func (c ConstantPool_Class) Type() (uint8) {
	return cp_class
}

type ConstantPool_FieldRef struct {
	ConstantPoolEntry
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

func (c ConstantPool_FieldRef) Type() (uint8) {
	return cp_field_ref
}

type ConstantPool_MethodRef struct {
	ConstantPoolEntry
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

func (c ConstantPool_MethodRef) Type() (uint8) {
	return cp_method_ref
}

type ConstantPool_InterfaceMethodRef struct {
	ConstantPoolEntry
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

func (c ConstantPool_InterfaceMethodRef) Type() (uint8) {
	return cp_interface_method_ref
}

type ConstantPool_String struct {
	ConstantPoolEntry
	StringIndex uint16
}

func (c ConstantPool_String) Type() (uint8) {
	return cp_string
}

type ConstantPool_Integer struct {
	ConstantPoolEntry
	Value uint32
}

func (c ConstantPool_Integer) Type() (uint8) {
	return cp_integer
}

type ConstantPool_Float struct {
	ConstantPoolEntry
	Value uint32
}

func (c ConstantPool_Float) Type() (uint8) {
	return cp_float
}

type ConstantPool_Long struct {
	ConstantPoolEntry
	High uint32
	Low  uint32
}

func (c ConstantPool_Long) Type() (uint8) {
	return cp_long
}

type ConstantPool_Double struct {
	ConstantPoolEntry
	High uint32
	Low  uint32
}

func (c ConstantPool_Double) Type() (uint8) {
	return cp_double
}

type ConstantPool_NameAndType struct {
	ConstantPoolEntry
	NameIndex       uint16
	DescriptorIndex uint16
}

func (c ConstantPool_NameAndType) Type() (uint8) {
	return cp_name_and_type
}

type ConstantPool_UTF8 struct {
	ConstantPoolEntry
	Value string
}

func (c ConstantPool_UTF8) Type() (uint8) {
	return cp_utf8
}

type ConstantPool_MethodHandle struct {
	ConstantPoolEntry
	Kind  uint8
	Index uint16
}

func (c ConstantPool_MethodHandle) Type() (uint8) {
	return cp_method_handle
}

type ConstantPool_Indy struct {
	ConstantPoolEntry
	BootstrapMethodIndex uint16
	NameAndTypeIndex     uint16
}

func (c ConstantPool_Indy) Type() (uint8) {
	return cp_indy
}

type ConstantPool_MethodType struct {
	ConstantPoolEntry
	DescriptorIndex uint16
}

func (c ConstantPool_MethodType) Type() (uint8) {
	return cp_method_type
}
