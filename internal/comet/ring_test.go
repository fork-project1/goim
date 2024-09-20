package comet

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	//NewRing(17)
	fmt.Printf("[16] %b\n", 16)
	fmt.Printf("[15] %b\n", 15)
	fmt.Printf("[14] %b\n", 14)
	fmt.Printf("[13] %b\n", 13)
	fmt.Printf("[12] %b\n", 12)
	fmt.Printf("[11] %b\n", 11)
	NewRing(15)
}
