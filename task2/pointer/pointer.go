package main

import "fmt"

/*
*
  - 编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，

在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值。
*/
func add(i *int) int {
	return *i + 10
}

/*
*

	实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
*/
func sliceMultiplyTwo(s []int) []int {
	// 移除返回值
	for i := range s {
		s[i] *= 2
	}
	return s
}

func main() {
	i := 12
	fmt.Println("调用add前 i =", i)
	add(&i)
	fmt.Println("调用add后 i =", i)

	s := []int{1, 2, 3, 4, 5}
	fmt.Println("调用sliceMultiplyTwo前 s =", s)
	sliceMultiplyTwo(s)
	fmt.Println("调用sliceMultiplyTwo后 s =", s)
}
