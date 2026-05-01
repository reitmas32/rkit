package logger

import "strings"

type Level int

const (
	LevelAll   Level = iota // all messages
	LevelDebug              // debug, info, warn, error, fatal, panic
	LevelInfo               // info, warn, error, fatal, panic
	LevelWarn               // warn, error, fatal, panic
	LevelError              // error, fatal, panic
	LevelFatal              // fatal, panic
	LevelPanic              // panic
)

// ParseLevel converts a string level ("debug", "info", ...) into Level.
// Returns ok=false if the input is unknown.
func ParseLevel(s string) (lvl Level, ok bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "all":
		return LevelAll, true
	case "debug":
		return LevelDebug, true
	case "info":
		return LevelInfo, true
	case "warn", "warning":
		return LevelWarn, true
	case "error":
		return LevelError, true
	case "fatal":
		return LevelFatal, true
	case "panic":
		return LevelPanic, true
	default:
		return LevelAll, false
	}
}

func (l Level) String() string {
	switch l {
	case LevelAll:
		return "all"
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	case LevelPanic:
		return "panic"
	default:
		return "unknown"
	}
}

// IsLoggable returns true if a message at msgLevel should be printed
// when the configured minimum level is minLevel.
func IsLoggable(minLevel, msgLevel Level) bool {
	// "all" prints everything
	if minLevel == LevelAll {
		return true
	}
	// Higher value = more severe. Print if message severity >= configured minimum.
	return msgLevel >= minLevel
}

// IsLoggableLevel is a convenience wrapper for string inputs (env/config).
// Unknown levels return false (fail-closed).
func IsLoggableLevel(minLevelStr, msgLevelStr string) bool {
	min, okMin := ParseLevel(minLevelStr)
	msg, okMsg := ParseLevel(msgLevelStr)
	if !okMin || !okMsg {
		return false
	}
	return IsLoggable(min, msg)
}
