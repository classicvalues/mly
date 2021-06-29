package storable

import (
	"github.com/viant/mly/common"
	"reflect"
)

//Field represents a  default storable field descriptor
type Field struct {
	Name     string
	DataType string
	dataType reflect.Type
}

func (f *Field) Type() reflect.Type {
	return f.dataType
}

func (f *Field) Init() (err error) {
	if f.dataType != nil {
		return nil
	}
	f.dataType, err = common.DataType(f.DataType)
	return err
}

//NewFields create new Fields
func NewFields(name string, dataType string) []*Field {
	field := &Field{Name: name, DataType: dataType}
	_ = field.Init()
	return []*Field{
		field,
	}
}
