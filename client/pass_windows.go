package main

import (
	"bufio"
	"log"
	"os"
)

func GetPass() string {
	bio := bufio.NewReader(os.Stdin)
	print("Please input the password>>")
	line, _, err := bio.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	return string(line)
}
