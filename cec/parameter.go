package cec

type ParameterType uint32

const (
	ParameterTypeString ParameterType = iota
	ParameterTypeUnknown
)

type Parameter struct {
	Type        ParameterType
	ValueString string
	ValueRaw    []byte
}

type nativeParameter struct {
	paramType ParameterType
	paramData uintptr
}

func (n nativeParameter) Go() Parameter {
	panic("TODO: nativeParameter.Go")
}
