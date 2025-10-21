package sane

import (
	"errors"
	"fmt"
	"log/slog"
	"unsafe"

	"github.com/maitredede/puregolibs/sane/internal"
	"github.com/maitredede/puregolibs/strings"
)

type OptionDescriptor struct {
	handle *Handle

	Number         int
	Name           string
	Title          string
	Description    string
	Type           Type
	Unit           Unit
	BinSize        int
	Cap            Cap
	ConstraintType ConstraintType
	Constraint     any
}

func (h *Handle) GetOptionDescriptors() ([]OptionDescriptor, error) {
	var count uint32
	var flags internal.SANE_Int
	ret := libSaneControlOption(h.h, 0, SANEActionGetValue, unsafe.Pointer(&count), &flags)
	if ret != StatusGood {
		return nil, fmt.Errorf("can't count options: %w", mkError(ret))
	}

	descriptors := make([]OptionDescriptor, count)
	for i := 0; i < int(count); i++ {

		// descPtr := libSaneGetOptionDescriptor(h.h, internal.SANE_Int(i))
		// d := convertOptionDescriptor(descPtr, i)
		nativeDesc := libSaneGetOptionDescriptor(h.h, internal.SANE_Int(i))
		d := convertOptionDescriptor(nativeDesc, i)

		d.handle = h

		descriptors[i] = d
	}
	return descriptors, nil
}

func convertOptionDescriptor(desc *internal.SANE_Option_Descriptor, number int) OptionDescriptor {
	d := OptionDescriptor{
		Number:         number,
		Type:           Type(desc.Type),
		Unit:           Unit(desc.Unit),
		BinSize:        int(desc.Size),
		Cap:            Cap(desc.Cap),
		ConstraintType: ConstraintType(desc.ConstraintType),
	}
	if desc.Name != nil {
		d.Name = strings.GoString((*byte)(desc.Name))
	}
	if desc.Title != nil {
		d.Title = strings.GoString((*byte)(desc.Title))
	}
	if desc.Desc != nil {
		d.Description = strings.GoString((*byte)(desc.Desc))
	}

	switch d.ConstraintType {
	case ConstraintRange:
		c := (*internalSANE_Range)(unsafe.Pointer(desc.Constraint))
		switch d.Type {
		case TypeInt:
			d.Constraint = &Range{
				Min:   intFromSane(c.Min),
				Max:   intFromSane(c.Max),
				Quant: intFromSane(c.Quant),
			}
		case TypeFloat:
			d.Constraint = &Range{
				Min:   floatFromSane(c.Min),
				Max:   floatFromSane(c.Max),
				Quant: floatFromSane(c.Quant),
			}
		}
		// c := (*SANE_Range)(unsafe.Pointer(desc.Constraint))
		// d.Constaint = c
	case ConstraintWordList:
		arr := (*internal.SANE_Word)(unsafe.Pointer(desc.Constraint))
		l := *arr
		c := unsafe.Slice(arr, uintptr(l+1))
		d.Constraint = c[1:]
	case ConstraintStringList:
		//arr := (*uintptr)(unsafe.Pointer(desc.Constraint))
		arr := (*[1 << 30]*uintptr)(unsafe.Pointer(desc.Constraint))
		strs := make([]string, 0, 30)
		ci := 0
		for {
			p := arr[ci]
			if p == nil {
				break
			}
			ci++

			s := strings.GoString((*byte)(unsafe.Pointer(p)))
			strs = append(strs, s)
		}
		d.Constraint = strs
	}
	return d
}

type internalSANE_Range struct {
	Min   internal.SANE_Word
	Max   internal.SANE_Word
	Quant internal.SANE_Word
}

func (h *Handle) GetOptionDescriptor(option int) *OptionDescriptor {
	ptr := libSaneGetOptionDescriptor(h.h, internal.SANE_Int(option))
	if ptr == nil {
		return nil
	}

	desc := convertOptionDescriptor(ptr, option)

	return &desc
}

func (d OptionDescriptor) GetValue() (any, error) {

	switch d.Type {
	case TypeBool:
		if d.BinSize != 4 {
			panic(fmt.Sprintf("FIXME: OptionDescriptor.GetValue bool binsize=%d", d.BinSize))
		}
		return d.GetValueBool()
	case TypeInt:
		if d.BinSize != 4 {
			panic(fmt.Sprintf("FIXME: OptionDescriptor.GetValue int binsize=%d", d.BinSize))
		}
		return d.GetValueInt()
	case TypeFloat:
		if d.BinSize != 4 {
			panic(fmt.Sprintf("FIXME: OptionDescriptor.GetValue float binsize=%d", d.BinSize))
		}
		return d.GetValueFloat()
	case TypeString:
		return d.GetValueString()
	}
	panic(fmt.Sprintf("FIXME: OptionDescriptor.GetValue unknown(%s) binsize=%d", d.Type, d.BinSize))
}

func (d OptionDescriptor) GetValueBool() (bool, error) {
	var val internal.SANE_Bool
	valPtr := unsafe.Pointer(&val)
	var flags internal.SANE_Int
	ret := libSaneControlOption(d.handle.h, internal.SANE_Int(d.Number), SANEActionGetValue, valPtr, &flags)
	if ret != StatusGood {
		return false, mkError(ret)
	}
	if flags != 0 {
		info := infoFromValue(flags)
		slog.Warn(fmt.Sprintf("getValueBool(%d) (%s) info=%+v", d.Number, d.Name, info))
	}
	return val.Go(), nil
}

func (d OptionDescriptor) GetValueInt() (int, error) {
	var val internal.SANE_Int
	valPtr := unsafe.Pointer(&val)
	var flags internal.SANE_Int
	ret := libSaneControlOption(d.handle.h, internal.SANE_Int(d.Number), SANEActionGetValue, valPtr, &flags)
	if ret != StatusGood {
		return 0, mkError(ret)
	}
	if flags != 0 {
		info := infoFromValue(flags)
		slog.Warn(fmt.Sprintf("getValueInt(%d) (%s) info=%+v", d.Number, d.Name, info))
	}
	return int(val), nil
}

func (d OptionDescriptor) GetValueFloat() (float64, error) {
	var val internal.SANE_Word
	valPtr := unsafe.Pointer(&val)
	var flags internal.SANE_Int
	ret := libSaneControlOption(d.handle.h, internal.SANE_Int(d.Number), SANEActionGetValue, valPtr, &flags)
	if ret != StatusGood {
		return 0, mkError(ret)
	}
	if flags != 0 {
		info := infoFromValue(flags)
		slog.Warn(fmt.Sprintf("getValueFloat(%d) (%s) info=%+v", d.Number, d.Name, info))
	}
	return floatFromSane(val), nil
}

func (d OptionDescriptor) GetValueString() (string, error) {
	val := make([]byte, d.BinSize)
	valPtr := unsafe.Pointer(&val)
	var flags internal.SANE_Int
	ret := libSaneControlOption(d.handle.h, internal.SANE_Int(d.Number), SANEActionGetValue, valPtr, &flags)
	if ret != StatusGood {
		return "", mkError(ret)
	}
	if flags != 0 {
		info := infoFromValue(flags)
		slog.Warn(fmt.Sprintf("getValueString(%d) (%s) info=%+v", d.Number, d.Name, info))
	}
	valStr := strings.GoString((*byte)(valPtr))
	return valStr, nil
}

func (d OptionDescriptor) SetValueAuto() (Info, error) {
	var flags internal.SANE_Int
	ret := libSaneControlOption(d.handle.h, internal.SANE_Int(d.Number), SANEActionSetAuto, nil, &flags)
	if ret != StatusGood {
		return Info{}, mkError(ret)
	}
	if flags != 0 {
		info := infoFromValue(flags)
		slog.Warn(fmt.Sprintf("setValueAuto(%d) (%s) info=%+v", d.Number, d.Name, info))
	}
	info := infoFromValue(flags)
	return info, nil
}

func (h *Handle) SetOptionValueAuto(name string) (Info, error) {
	opts, err := h.GetOptionDescriptors()
	if err != nil {
		return Info{}, err
	}

	for _, desc := range opts {
		if desc.Name != name {
			continue
		}

		return desc.SetValueAuto()
	}
	return Info{}, fmt.Errorf("option '%s' not found", name)
}

func (d OptionDescriptor) SetValue(value any) (Info, error) {
	switch d.Type {
	case TypeBool:
		if d.BinSize != 4 {
			panic(fmt.Sprintf("FIXME: OptionDescriptor.SetValue bool binsize=%d", d.BinSize))
		}
		b, ok := value.(bool)
		if !ok {
			return Info{}, errors.New("value not a bool")
		}
		return d.SetValueBool(b)
	case TypeInt:
		if d.BinSize != 4 {
			panic(fmt.Sprintf("FIXME: OptionDescriptor.SetValue int binsize=%d", d.BinSize))
		}
		i, ok := value.(int)
		if !ok {
			return Info{}, errors.New("value not a int")
		}
		return d.SetValueInt(i)
	case TypeFloat:
		if d.BinSize != 4 {
			panic(fmt.Sprintf("FIXME: OptionDescriptor.SetValue float binsize=%d", d.BinSize))
		}
		f, ok := value.(float64)
		if !ok {
			return Info{}, errors.New("value not a float64")
		}
		return d.SetValueFloat(f)
	case TypeString:
		s, ok := value.(string)
		if !ok {
			return Info{}, errors.New("value not a string")
		}
		return d.SetValueString(s)
	}
	panic(fmt.Sprintf("FIXME: OptionDescriptor.SetValue unknown(%s) binsize=%d", d.Type, d.BinSize))
}

func (h *Handle) SetOptionValue(name string, value any) (Info, error) {
	opts, err := h.GetOptionDescriptors()
	if err != nil {
		return Info{}, err
	}

	for _, desc := range opts {
		if desc.Name != name {
			continue
		}

		return desc.SetValue(value)
	}
	return Info{}, fmt.Errorf("option '%s' not found", name)
}

func (d OptionDescriptor) SetValueBool(value bool) (Info, error) {
	var nValue internal.SANE_Bool
	if value {
		nValue = internal.SANE_TRUE
	} else {
		nValue = internal.SANE_FALSE
	}
	var flags internal.SANE_Int
	ret := libSaneControlOption(d.handle.h, internal.SANE_Int(d.Number), SANEActionSetValue, unsafe.Pointer(&nValue), &flags)
	if ret != StatusGood {
		return Info{}, mkError(ret)
	}
	if flags != 0 {
		info := infoFromValue(flags)
		slog.Warn(fmt.Sprintf("setValueBool(%d) (%s) info=%+v", d.Number, d.Name, info))
	}
	info := infoFromValue(flags)
	return info, nil
}

func (d OptionDescriptor) SetValueString(value string) (Info, error) {
	nValue := strings.CString(value)
	var flags internal.SANE_Int
	ret := libSaneControlOption(d.handle.h, internal.SANE_Int(d.Number), SANEActionSetValue, unsafe.Pointer(nValue), &flags)
	if ret != StatusGood {
		return Info{}, mkError(ret)
	}
	if flags != 0 {
		info := infoFromValue(flags)
		slog.Warn(fmt.Sprintf("setValueString(%d) (%s) info=%+v", d.Number, d.Name, info))
	}
	info := infoFromValue(flags)
	return info, nil
}

func (d OptionDescriptor) SetValueInt(value int) (Info, error) {
	var nValue internal.SANE_Int
	nValue = internal.SANE_Int(value)

	var flags internal.SANE_Int
	ret := libSaneControlOption(d.handle.h, internal.SANE_Int(d.Number), SANEActionSetValue, unsafe.Pointer(&nValue), &flags)
	if ret != StatusGood {
		return Info{}, mkError(ret)
	}
	if flags != 0 {
		info := infoFromValue(flags)
		slog.Warn(fmt.Sprintf("setValueInt(%d) (%s) info=%+v", d.Number, d.Name, info))
	}
	info := infoFromValue(flags)
	return info, nil
}

func (d OptionDescriptor) SetValueFloat(value float64) (Info, error) {
	nValue := floatToSane(value)

	var flags internal.SANE_Int
	ret := libSaneControlOption(d.handle.h, internal.SANE_Int(d.Number), SANEActionSetValue, unsafe.Pointer(&nValue), &flags)
	if ret != StatusGood {
		return Info{}, mkError(ret)
	}
	if flags != 0 {
		info := infoFromValue(flags)
		slog.Warn(fmt.Sprintf("setValueInt(%d) (%s) info=%+v", d.Number, d.Name, info))
	}
	info := infoFromValue(flags)
	return info, nil
}
