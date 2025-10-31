//go:build linux

package libevdi

const (
	_IOC_NONE  uintptr = 0
	_IOC_WRITE uintptr = 1
	_IOC_READ  uintptr = 2

	_IOC_NRBITS   = 8
	_IOC_TYPEBITS = 8
	_IOC_SIZEBITS = 14

	_IOC_NRSHIFT   int = 0
	_IOC_TYPESHIFT int = (_IOC_NRSHIFT + _IOC_NRBITS)
	_IOC_SIZESHIFT int = (_IOC_TYPESHIFT + _IOC_TYPEBITS)
	_IOC_DIRSHIFT  int = (_IOC_SIZESHIFT + _IOC_SIZEBITS)
)

func IOC(dir, typ, nr, size uintptr) uintptr {
	return (dir << _IOC_DIRSHIFT) |
		(size << _IOC_SIZESHIFT) |
		(typ << _IOC_TYPESHIFT) |
		(nr << _IOC_NRSHIFT)
}

// Remplacement approximatif de la macro _IOW(type, nr, size)
// Les constantes de direction, type et bits sont dans unix/ioctl.go
func IOW(typ, nr, size uintptr) uintptr {
	return IOC(_IOC_WRITE, typ, nr, size)
}

func IOR(typ, nr, size uintptr) uintptr {
	return IOC(_IOC_READ, typ, nr, size)
}

func IOWR(typ, nr, size uintptr) uintptr {
	return IOC(_IOC_WRITE|_IOC_READ, typ, nr, size)
}

func IO(typ, nr uintptr) uintptr {
	return IOC(_IOC_NONE, typ, nr, 0)
}
