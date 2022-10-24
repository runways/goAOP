package unitTests

import (
	"fmt"
)

type FirstStruct struct {
	name string
	Age  int
}

// InvokeFirstFunction
// @middleware-a
func (fs FirstStruct) InvokeFirstFunction() string {
	fmt.Println("---> fs.InvokeFirstFunction")
	return ""
}

func (fs FirstStruct) InvokeSecondFunction() string {
	return ""
}
