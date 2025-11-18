package libudev

import "unsafe"

type ListEntry unsafe.Pointer

func ListEntryForEach(firstEntry ListEntry, fn func(entry ListEntry)) {
	initLib()

	for listEntry := firstEntry; listEntry != nil; listEntry = libudevListEntryGetNext(listEntry) {
		fn(listEntry)
	}
}

func ListEntryGetName(entry ListEntry) string {
	initLib()

	return libudevListEntryGetName(entry)
}
