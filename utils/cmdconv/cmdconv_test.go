package cmdconv

import (
	"fmt"
	"testing"
)

func TestCmdConv(t *testing.T) {
	args := [][]byte{[]byte("key1"), []byte("val1")}
	line := ToCmdLineArgs("del", args)
	for _, bytes := range line {
		fmt.Println(string(bytes))
	}
}
