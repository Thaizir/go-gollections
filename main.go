package main

import "github.com/thaizir/gollections/set"

func main() {

	s := set.New[string]()
	result := s.Insert("Thaizir")
	println(result)
	println(s.Size())
	println(s.Contains("Thaizir"))
}
