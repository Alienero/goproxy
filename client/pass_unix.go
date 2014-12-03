// +build darwin dragonfly freebsd linux netbsd openbsd

package main

import (
	"github.com/Alienero/goproxy/3rd/code.google.com/p/gopass"
)

func GetPass() string {
	psw, err := gopass.GetPass("Please input the password>>")
	if err != nil {
		panic(err)
	}
	return psw
}
