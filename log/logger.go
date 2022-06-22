package log

import "fmt"

type Logger interface {
	Log(string, ...any)
}

type PrefixedLogger struct {
	Prefix string
}

func (l *PrefixedLogger) Log(format string, args ...any) {
	fmt.Printf("%s %s\n", l.Prefix, fmt.Sprintf(format, args...))
}

func NewLogger(prefixFormat string, args ...any) *PrefixedLogger {
	return &PrefixedLogger{fmt.Sprintf(prefixFormat, args...)}
}

func (l *PrefixedLogger) Extend(prefixFormat string, args ...any) *PrefixedLogger {
	return &PrefixedLogger{
		fmt.Sprintf("%s %s", l.Prefix,
			fmt.Sprintf(prefixFormat, args...)),
	}
}
