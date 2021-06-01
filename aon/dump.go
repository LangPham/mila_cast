package aon


import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"strings"
)

func Dump(int interface{}, label ...string) {
	if len(label) == 1 {
		fmt.Println(strings.Repeat("=",30) + "| "+ label[0] +" |" + strings.Repeat("=",30) + " BEGIN")
	} else {
		fmt.Println( strings.Repeat("=",70)+ " BEGIN")
	}
	spew.Dump(int)
	if len(label) == 1 {
		fmt.Println(strings.Repeat("=",30) + "| "+ label[0] +" |" + strings.Repeat("=",30)+ " END")
	} else {
		fmt.Println(strings.Repeat("=",70) +  " END")
	}
}
