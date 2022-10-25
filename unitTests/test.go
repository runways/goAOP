package unitTests

import (
	"fmt"
)

type FirstStruct struct {
	name string `json:"name"`
	Age  int    `json:"age"`
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

// InvokeFirstFunction a same name function, but not belongs any struct
// @middleware-a
func InvokeFirstFunction() {}

// InvokeSecondFunction a same name function, but not belongs any struct
// @middleware-b
func InvokeSecondFunction() {}
