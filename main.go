package main

import (
	"fmt"
)

type HSL struct {
	H, S, L int
}

func (hsl *HSL) String() string {
	return fmt.Sprintf("HSL(%d, %d%%, %d%%)", hsl.H, hsl.S, hsl.L)
}

type HSV struct {
	H, S, V int
}

func (hsl *HSV) String() string {
	return fmt.Sprintf("HSV(%d, %d%%, %d%%)", hsl.H, hsl.S, hsl.V)
}

type RGB struct {
	R, G, B uint8
}

func main() {

}
