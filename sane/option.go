package sane

import (
	"fmt"
	"unsafe"

	"github.com/maitredede/puregolibs/strings"
)

// Option represents a scanning option.
type Option struct {
	Name         string        // option name
	Group        string        // option group
	Title        string        // option title
	Desc         string        // option description
	Type         Type          // option type
	Unit         Unit          // units
	Length       int           // vector length for vector-valued options
	ConstrSet    []interface{} // constraint set
	ConstrRange  *Range        // constraint range
	IsActive     bool          // whether option is active
	IsSettable   bool          // whether option can be set
	IsDetectable bool          // whether option value can be detected
	IsAutomatic  bool          // whether option has an auto value
	IsEmulated   bool          // whether option is emulated
	IsAdvanced   bool          // whether option is advanced
	index        int           // internal option index
	size         int           // internal option size in bytes
}

type OptionDescriptor struct {
	Number         int
	Name           string
	Title          string
	Description    string
	Type           Type
	Unit           Unit
	Size           int
	Cap            Cap
	ConstraintType ConstraintType
	Constaint      any
}

type internalOptionDescriptor struct {
	Name           uintptr
	Title          uintptr
	Desc           uintptr
	Type           Type
	Unit           Unit
	Size           SANE_Int
	Cap            Cap
	ConstraintType ConstraintType
	Constraint     uintptr
}

func GetOptionDescriptors(h SANE_Handle) ([]OptionDescriptor, error) {
	var count uint32
	var flags Info
	ret := libSaneControlOption(h, 0, SANEActionGetValue, unsafe.Pointer(&count), &flags)
	if ret != StatusGood {
		return nil, fmt.Errorf("can't count options: %w", mkError(ret))
	}

	descriptors := make([]OptionDescriptor, count)
	for i := 0; i < int(count); i++ {

		descPtr := libSaneGetOptionDescriptor(h, SANE_Int(i))
		desc := (*internalOptionDescriptor)(unsafe.Pointer(descPtr))

		d := OptionDescriptor{
			Number:         i,
			Type:           desc.Type,
			Unit:           desc.Unit,
			Size:           int(desc.Size),
			Cap:            desc.Cap,
			ConstraintType: desc.ConstraintType,
		}
		if desc.Name != 0 {
			d.Name = strings.GoString(desc.Name)
		}
		if desc.Title != 0 {
			d.Title = strings.GoString(desc.Title)
		}
		if desc.Desc != 0 {
			d.Description = strings.GoString(desc.Desc)
		}

		switch desc.ConstraintType {
		case ConstraintRange:
			c := (*SANE_Range)(unsafe.Pointer(desc.Constraint))
			d.Constaint = c
		case ConstraintWordList:
			arr := (*SANE_Word)(unsafe.Pointer(desc.Constraint))
			l := *arr
			c := unsafe.Slice(arr, uintptr(l+1))
			d.Constaint = c[1:]
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

				s := strings.GoString(uintptr(unsafe.Pointer(p)))
				strs = append(strs, s)
			}
			d.Constaint = strs
		}

		descriptors[i] = d
	}
	return descriptors, nil
}

type SANE_Range struct {
	Min   SANE_Word
	Max   SANE_Word
	Quant SANE_Word
}
