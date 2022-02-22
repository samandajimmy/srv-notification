package logger

type ContextKey = uint8

const (
	_ ContextKey = iota + 1
	RequestIdKey
)
