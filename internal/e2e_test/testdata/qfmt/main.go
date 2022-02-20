package main

import "fmt"

func main() {
	println(fmt.Sprintf("ok"))
	println(fmt.Sprintf("%s", "ok"))
	println(fmt.Sprintf("%s:%d", "file.go", 1043))

	formatString := "hello, %s!"
	println(fmt.Sprintf(formatString, "world"))
}
