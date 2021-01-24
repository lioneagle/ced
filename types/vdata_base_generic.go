package types

//go:generate genny -in=$GOFILE -out=vfloat64_base.go gen "DataType=float64"

import (
	"github.com/cheekybits/genny/generic"
	"github.com/pkg/errors"
)

type DataType generic.Type

type VDataTypeBase struct {
	name string
	data []DataType
}

func NewVDataTypeBase(name string, data []DataType) *VDataTypeBase {
	ret := &VDataTypeBase{
		name: name,
		data: data,
	}

	if data == nil {
		ret.data = make([]DataType, 0)
	}
	return ret
}

func (this *VDataTypeBase) GetName() string {
	return this.name
}

func (this *VDataTypeBase) Len() int {
	return len(this.data)
}

func (this *VDataTypeBase) Add(vals ...DataType) error {
	this.data = append(this.data, vals...)
	return nil
}

func (this *VDataTypeBase) GetData() interface{} {
	return this.data
}

func (this *VDataTypeBase) GetDataAt(index int) DataType {
	return this.data[index]
}

func (this *VDataTypeBase) Clone(name string) *VDataTypeBase {
	data := make([]DataType, len(this.data))
	copy(data, this.data)
	return NewVDataTypeBase(name, data)
}

func (this *VDataTypeBase) Map(op func(params ...DataType) DataType, rhs ...*VDataTypeBase) error {
	if op == nil {
		return errors.Errorf("op is nil")
	}

	num := len(this.data)

	for _, v := range rhs {
		if num > v.Len() {
			num = v.Len()
		}
	}

	params := make([]DataType, len(rhs), len(rhs))

	for i := 0; i < num; i++ {
		for j := 0; j < len(rhs); j++ {
			params[j] = rhs[j].GetDataAt(i)
		}

		this.data[i] = op(params...)
	}

	return nil
}
