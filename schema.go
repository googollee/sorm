package sorm

import (
	"reflect"
)

type fieldType string

const (
	fieldText    fieldType = "TEXT"
	fieldInteger           = "INTEGER"
	fieldVarChar           = "VARCHAR"
)

type FieldType struct {
	typ  fieldType
	size int
}

func Text() FieldType {
	return FieldType{
		typ: fieldText,
	}
}
func VarChar(size int) FieldType {
	return FieldType{
		typ:  fieldVarChar,
		size: size,
	}
}

func Integer() FieldType {
	return FieldType{
		typ: fieldInteger,
	}
}

func IntegerByte(size int) FieldType {
	return FieldType{
		typ:  fieldInteger,
		size: size,
	}
}

type tableDefinition struct {
	fields  []Field
	indexes []Index
}

type Schema struct {
	fieldStart uintptr
	fieldSize  uintptr
	fields     map[any]*Field

	definition *tableDefinition
}

func newSchema(model any) *Schema {
	vmodel := reflect.ValueOf(model)
	if vmodel.Kind() != reflect.Pointer {
		panic("should be a pointer")
	}

	return &Schema{
		fieldStart: uintptr(vmodel.UnsafePointer()),
		fieldSize:  vmodel.Type().Elem().Size(),
		fields:     make(map[any]*Field),
		definition: &tableDefinition{},
	}
}

func (t *Schema) Field(f any) *Field {
	vfield := reflect.ValueOf(f)
	if vfield.Kind() != reflect.Pointer {
		panic("should be a pointer")
	}

	addrField := uintptr(vfield.UnsafePointer())
	offset := addrField - t.fieldStart
	if offset < 0 || offset > t.fieldSize {
		panic("should be one of the struct field")
	}

	tField := vfield.Type().Elem()

	t.definition.fields = append(t.definition.fields, newFieldWithType(tField))
	pField := &t.definition.fields[len(t.definition.fields)-1]
	t.fields[f] = pField
	pField.parent = t.definition

	return pField
}

func (t *Schema) UniqueIndex(fields ...any) {
	index := Index{
		fields: make([]*Field, 0, len(fields)),
		unique: true,
	}

	for _, field := range fields {
		schemaField, ok := t.fields[field]
		if !ok {
			panic("should be one of the struct field")
		}

		index.fields = append(index.fields, schemaField)
	}

	t.definition.indexes = append(t.definition.indexes, index)
}

type Field struct {
	parent *tableDefinition

	name          string
	typ           FieldType
	autoIncrement bool
	primaryKey    bool
	nullable      bool
}

func newFieldWithType(fieldType reflect.Type) Field {
	ret := Field{}

	switch fieldType.Kind() {
	case reflect.String:
		ret.typ = Text()

	case reflect.Int:
		ret.typ = Integer()
	case reflect.Int8:
		ret.typ = IntegerByte(1)
	case reflect.Int16:
		ret.typ = IntegerByte(2)
	case reflect.Int32:
		ret.typ = IntegerByte(4)
	case reflect.Int64:
		ret.typ = IntegerByte(8)

	case reflect.Uint:
		ret.typ = Integer()
	case reflect.Uint8:
		ret.typ = IntegerByte(1)
	case reflect.Uint16:
		ret.typ = IntegerByte(2)
	case reflect.Uint32:
		ret.typ = IntegerByte(4)
	case reflect.Uint64:
		ret.typ = IntegerByte(8)
	}

	return ret
}

func (f *Field) Name(name string) *Field {
	f.name = name
	return f
}

func (f *Field) PrimaryKey() *Field {
	f.primaryKey = true
	f.autoIncrement = true
	return f
}

func (f *Field) AutoIncrement(ok bool) *Field {
	f.autoIncrement = ok
	return f
}

func (f *Field) Type(typ FieldType) *Field {
	f.typ = typ
	return f
}

func (f *Field) Nullable(ok bool) *Field {
	f.nullable = ok
	return f
}

func (f *Field) NotNull() *Field {
	f.nullable = false
	return f
}

func (f *Field) UniqueIndex() *Field {
	index := Index{
		fields: []*Field{f},
		unique: true,
	}

	f.parent.indexes = append(f.parent.indexes, index)

	return f
}

type Index struct {
	fields []*Field
	unique bool
}
