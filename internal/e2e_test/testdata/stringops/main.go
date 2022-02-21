package main

func testindex() {
	println("abc"[0])
	s := "hello"
	println(s[0])
	println(s[len(s)-1])
	// for i := 0; i < len(s); i++ {
	// 	println(s[i])
	// }
}

func main() {
	testindex()
}
