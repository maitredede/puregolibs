package libfreefare

import (
	"errors"
	"unsafe"

	"github.com/maitredede/puregolibs/libnfc"
	"github.com/maitredede/puregolibs/strings"
)

type TagInfo struct {
	Type  TagType
	Name  string
	UID   string
	Error string
}

// type nativeNfcTarget struct {
// 	nti nativeNfcTargetInfo
// 	nm  libnfc.Modulation
// }

// type nativeNfcTargetInfo [283]byte

// type nativeFreefareTag struct {
// 	device  uintptr
// 	info    nativeNfcTarget
// 	typ     int16
// 	active  int16
// 	timeout int16
// 	tag     uintptr
// }

func GetTagInfos(device *libnfc.NfcDevice) ([]TagInfo, error) {
	libInit()

	tags, err := GetTags(device)
	if err != nil {
		return nil, err
	}
	defer tags.Close()

	infos := make([]TagInfo, tags.Len())
	for i := 0; i < tags.Len(); i++ {
		tag := tags.Get(i)

		info := TagInfo{
			Type:  tag.TagType(),
			Name:  tag.Name(),
			UID:   tag.UID(),
			Error: tag.Error(),
		}
		infos[i] = info
	}
	return infos, nil
}

func GetTags(device *libnfc.NfcDevice) (MifareTagList, error) {
	libInit()

	lst := MifareTagList{}

	tags := libGetTags(device.Ptr())
	if tags == nil {
		return lst, nil
	}
	lst.tags = tags

	ptr := *(*unsafe.Pointer)(unsafe.Pointer(tags))
	size := int(unsafe.Sizeof((uintptr)(0)))
	var length int
	for {
		if *(*uintptr)(unsafe.Add(ptr, uintptr(length*size))) == 0 {
			break
		}
		length++
	}
	values := unsafe.Slice(tags, length)
	for _, tagPtr := range values {
		if tagPtr == 0 {
			break
		}

		item := MifareTag{
			ptr:         tagPtr,
			shouldClear: false,
		}
		lst.items = append(lst.items, item)
	}
	return lst, nil
}

type MifareTagList struct {
	tags  *uintptr
	items []MifareTag
}

func (l *MifareTagList) Close() {
	if l.tags != nil {
		libFreeTags(l.tags)
		l.tags = nil
	}
}

func (l *MifareTagList) Len() int {
	return len(l.items)
}

func (l *MifareTagList) Get(i int) MifareTag {
	return l.items[i]
}

type MifareTag struct {
	ptr         uintptr
	shouldClear bool
}

func (t *MifareTag) Close() error {
	if t.ptr == 0 {
		return ErrTagClosed
	}
	if t.shouldClear {
		libFreeTag(t.ptr)
		t.ptr = 0
	}
	return nil
}

func (t *MifareTag) TagType() TagType {
	if t.ptr == 0 {
		return TagType(0)
	}
	return libGetTagType(t.ptr)
}

func (t *MifareTag) Name() string {
	if t.ptr == 0 {
		return ""
	}
	return libGetTagFirendlyName(t.ptr)
}

func (t *MifareTag) UID() string {
	if t.ptr == 0 {
		return ""
	}
	uidPtr := libGetTagUID(t.ptr)
	defer libFree(uidPtr)

	return strings.GoString(uidPtr)
}

func (t *MifareTag) Error() string {
	if t.ptr == 0 {
		return ""
	}
	return libStrError(t.ptr)
}

func (t *MifareTag) MifareClassicConnect() error {
	if t.ptr == 0 {
		return ErrTagClosed
	}
	ret := libMifareClassicConnect(t.ptr)
	if ret != 0 {
		return errors.New(libStrError(t.ptr))
	}
	return nil
}

func (t *MifareTag) MifareClassicDisconnect() error {
	if t.ptr == 0 {
		return ErrTagClosed
	}
	ret := libMifareClassicDisconnect(t.ptr)
	if ret != 0 {
		return errors.New(libStrError(t.ptr))
	}
	return nil
}

func (t *MifareTag) ReadMad() (*Mad, error) {
	if t.ptr == 0 {
		return nil, ErrTagClosed
	}
	madPtr := libMadRead(t.ptr)
	if madPtr == 0 {
		return nil, errors.New("no MAD detected")
	}
	mad := &Mad{
		tag: t,
		ptr: madPtr,
	}
	return mad, nil
}
