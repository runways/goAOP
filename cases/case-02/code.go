package case_02

import (
	"fmt"
	"math"
)

type FirstStruct struct {
	name string `json:"name"`
	Age  int    `json:"age"`
}

// addWithFuncDependWithInjection
// @middleware-func-depend
func (fs FirstStruct) addWithFuncDependWithInjection() {
	y := math.Round(1)
	
	fmt.Println(y)
}
