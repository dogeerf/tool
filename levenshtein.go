package main

import (
	"fmt"
	"github.com/arbovm/levenshtein"
)

func main() {
	s1 := "kitten"
	s2 := "sitting"
	fmt.Printf("The distance between %v and %v is %v\n",
		s1, s2, levenshtein.Distance(s1, s2))
	// -> The distance between kitten and sitting is 3
}
