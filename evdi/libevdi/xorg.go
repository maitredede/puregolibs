//go:build linux

package libevdi

import (
	"github.com/prometheus/procfs"
)

func XorgRunning() bool {
	return iterateThroughAllProcessFoldersAndFindXorg()
}

func iterateThroughAllProcessFoldersAndFindXorg() bool {
	procs, err := procfs.AllProcs()
	if err != nil {
		evdiLogDebug("can't open /proc folder: %v", err)
		return false
	}
	for _, p := range procs {
		s, err := p.Stat()
		if err != nil {
			evdiLogDebug("can't stat %v: %v", p, err)
			continue
		}
		if s.Comm == "Xorg" {
			return true
		}
	}
	return false
}

// func isXorgProcessFolder(procEntry os.DirEntry) bool {
// 	folderName := procEntry.Name()
// 	if isNumeric(folderName) && isXorg(folderName) {
// 		return true
// 	}
// 	return false
// }

// func isNumeric(str string) bool {
// 	if len(str) == 0 {
// 		return false
// 	}
// 	for _, c := range str {
// 		if !unicode.IsDigit(c) {
// 			return false
// 		}
// 	}
// 	return true
// }

// func isXorg(pid string) bool {
// 	processFolder := fmt.Sprintf("/proc/%s/stat", pid)
// 	procfs.Self()
// }
