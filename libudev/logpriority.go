package libudev

import "log/syslog"

type LogPriority int32

const (
	LogEmerg    LogPriority = LogPriority(syslog.LOG_EMERG)
	LogAlert    LogPriority = LogPriority(syslog.LOG_ALERT)
	LogCritical LogPriority = LogPriority(syslog.LOG_CRIT)
	LogError    LogPriority = LogPriority(syslog.LOG_ERR)
	LogWarning  LogPriority = LogPriority(syslog.LOG_WARNING)
	LogNotice   LogPriority = LogPriority(syslog.LOG_NOTICE)
	LogInfo     LogPriority = LogPriority(syslog.LOG_INFO)
	LogDebug    LogPriority = LogPriority(syslog.LOG_DEBUG)
)
