package libevdi

type AccessMode int

const (
	R_OK AccessMode = 4 /* Test for read permission.  */
	W_OK AccessMode = 2 /* Test for write permission.  */
	X_OK AccessMode = 1 /* Test for execute permission.  */
	F_OK AccessMode = 0 /* Test for existence.  */
)

func access(file string, mode AccessMode) bool {
	panic("TODO")
}
