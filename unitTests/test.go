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
// @middleware-b
// @middleware-return
func (fs FirstStruct) InvokeFirstFunction() string {
	x := "a pre-defind variable"
	fmt.Println("---> fs.InvokeFirstFunction", x)
	
	// Since this function return a string, but it will not insert
	// 'return code' although user use @middleware-return
	return ""
}

func (fs FirstStruct) InvokeSecondFunction() string {
	return ""
}

// InvokeThirdFunction
// @middleware-return
func (fs FirstStruct) InvokeThirdFunction() func() {
	return func() {
		fmt.Println("InvokeThirdFunction end")
	}
}

// InvokeFirstFunction a same name function, but not belongs any struct
// @middleware-a
func InvokeFirstFunction() {}

// InvokeSecondFunction a same name function, but not belongs any struct
// @middleware-b
func InvokeSecondFunction() {}
