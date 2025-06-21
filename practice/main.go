package main

import "fmt"

//func add(x, y int) int {
//	return x + y
//}
//
//func swap(x string, y string) (string, string) {
//	return y, x
//}
//
//func split(sum int) (x int, y int) {
//	y = sum / 2
//	x = sum - y
//	return
//}
//
//var c, python, java bool
//var i, j int = 1, 2
//
//func main() {
//	fmt.Println("Hello World")
//	fmt.Println(math.MaxInt8)
//	fmt.Println(rand.Intn(2))
//	fmt.Println(math.Pi)
//	fmt.Println(add(2, 3))
//	fmt.Println(swap("world", "hello"))
//	fmt.Println(split(12))
//	var i int
//	fmt.Println(i, c, python, java)
//	var c, python, java = true, false, "ні!"
//	fmt.Println(i, j, c, python, java)
//	var l, j int = 1, 2
//	k := 3
//	x, y, z := true, false, "ні!"
//	fmt.Println(l, j, k, x, y, z)
//}
//
//func main2() {
//	var i int
//	var f float64
//	var b bool
//	var s string
//	fmt.Printf("%v %v %v %q\n", i, f, b, s)
//}
//
//func init() {
//	fmt.Println("Init is executed before ...")
//	main2()
//}

func main() {
	var sum int
	for i := 0; i < 10; i++ {
		sum = sum + i
	}
	fmt.Println(sum)
}
