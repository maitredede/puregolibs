package cec

import (
	gostrings "strings"
)

type LogLevel int32

const (
	LogError   LogLevel = 1
	LogWarning LogLevel = 2
	LogNotice  LogLevel = 4
	LogTraffic LogLevel = 8
	LogDebug   LogLevel = 16
	LogAll     LogLevel = 31
)

func (l LogLevel) String() string {
	levels := map[LogLevel]string{
		LogError:   "error",
		LogWarning: "warning",
		LogNotice:  "notice",
		LogTraffic: "traffic",
		LogDebug:   "debug",
	}
	strs := make([]string, 0, len(levels))
	for t, s := range levels {
		if l == t {
			return s
		} else {
			if (l & t) == t {
				strs = append(strs, s)
			}
		}
	}
	return gostrings.Join(strs, ",")
}
