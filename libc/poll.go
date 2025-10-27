package libc

type Pollfd struct {
	Fd      int32
	Events  int16
	Revents int16
}

const (
	POLLIN  int16 = 0x001
	POLLPRI int16 = 0x002
	POLLOUT int16 = 0x004

	POLLRDNORM int16 = 0x040
	POLLRDBAND int16 = 0x080
	POLLWRNORM int16 = 0x100
	POLLWRBAND int16 = 0x200

	POLLMSG    int16 = 0x400
	POLLREMOVE int16 = 0x1000
	POLLRDHUP  int16 = 0x2000

	POLLERR  int16 = 0x008
	POLLHUP  int16 = 0x010
	POLLNVAL int16 = 0x020
)

func Poll(fds []Pollfd, timeout int) int {
	initCLib()

	fdsPtr := &fds[0]
	nfds := len(fds)

	ret := libcPoll(fdsPtr, uint32(nfds), int32(timeout))
	return int(ret)
}
