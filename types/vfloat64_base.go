// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package types

import "github.com/pkg/errors"

type VFloat64Base struct {
	name string
	data []float64
}

func NewVFloat64Base(name string, data []float64) *VFloat64Base {
	ret := &VFloat64Base{
		name: name,
		data: data,
	}

	if data == nil {
		ret.data = make([]float64, 0)
	}
	return ret
}

func (this *VFloat64Base) GetName() string {
	return this.name
}

func (this *VFloat64Base) Len() int {
	return len(this.data)
}

func (this *VFloat64Base) Add(vals ...float64) error {
	this.data = append(this.data, vals...)
	return nil
}

func (this *VFloat64Base) GetData() interface{} {
	return this.data
}

func (this *VFloat64Base) GetDataAt(index int) float64 {
	return this.data[index]
}

func (this *VFloat64Base) Clone(name string) *VFloat64Base {
	data := make([]float64, len(this.data))
	copy(data, this.data)
	return NewVFloat64Base(name, data)
}

func (this *VFloat64Base) Map(op func(params ...float64) float64, rhs ...*VFloat64Base) error {
	if op == nil {
		return errors.Errorf("op is nil")
	}

	num := len(this.data)

	for _, v := range rhs {
		if num > v.Len() {
			num = v.Len()
		}
	}

	params := make([]float64, len(rhs), len(rhs))

	for i := 0; i < num; i++ {
		for j := 0; j < len(rhs); j++ {
			params[j] = rhs[j].GetDataAt(i)
		}

		this.data[i] = op(params...)
	}

	return nil
}
