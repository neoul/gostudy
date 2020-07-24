package main

import "fmt"

func main() {

	var ss []int = nil
	for _, x := range ss {
		fmt.Println(x)
	}
	s := make([]int, 0, len(ss))
	fmt.Println(s)
	s = append(s, 100)
	fmt.Println(s)
	for i := 0; i < 10; i++ {
		s = append(s, i)
	}
	fmt.Println(s)

}
