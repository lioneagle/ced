package types

type VDataType uint8

const (
	VDATETIME VDataType = iota
	VINT
	VFLOAT
	VSTRING

	vDataTypeMAX
)

var g_datatype_names = []string{
	"vdatetime",
	"vint",
	"vfloat",
	"vstring",
}

func (v VDataType) String() string {
	if v >= vDataTypeMAX {
		return "unknown datatype"
	}
	return g_datatype_names[v]
}

type VData interface {
	GetType() VDataType
	GetName() string
	Len() int
	AddByStrings(vals ...string) error
	GetData() interface{}
	GetDataStringAt(index int) string
	String() []string
}

/*func NewVData(name string, datatype VDataType) VData {
	switch datatype {
	case VDATETIME:
		return NewVDatetime(name, nil)
	case VINT:
		return NewVInt(name, nil)
	case VFLOAT:
		return NewVFloat(name, nil)
	case VSTRING:
		return NewVString(name, nil)
	default:
		return nil
	}
}

func GetAsVDatetime(data VData) *VDatetime {
	v, ok := data.(*VDatetime)
	if !ok {
		return nil
	}
	return v
}

func GetAsVInt(data VData) *VInt {
	v, ok := data.(*VInt)
	if !ok {
		return nil
	}
	return v
}

func GetAsVFloat(data VData) *VFloat {
	v, ok := data.(*VFloat)
	if !ok {
		return nil
	}
	return v
}

func GetAsVString(data VData) *VString {
	v, ok := data.(*VString)
	if !ok {
		return nil
	}
	return v
}
*/
