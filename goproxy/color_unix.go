// +build darwin dragonfly freebsd linux netbsd openbsd

package goproxy

const (
	Red   = "31"
	Green = "32"
	Blue  = "34"

	Start = "\033["
	Mid   = ";49;1m"
	End   = "\033[39;49;0m"
)
