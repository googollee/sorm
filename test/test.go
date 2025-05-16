package main

import (
	"fmt"
)

func f(s *string) {
	fmt.Println(*s)
}

func main() {
	m := map[string]string{"a": "a"}
	f(&m["a"])
}
