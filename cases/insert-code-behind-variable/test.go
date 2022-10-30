package insert_code_behind_variable

import "fmt"

type FirstStruct struct {
	name string `json:"name"`
	Age  int    `json:"age"`
}

// invokeFirstFunction
// @middleware-err
func (fs FirstStruct) invokeFirstFunction() {

	err := fmt.Errorf("a pre-defind error")
	fmt.Println(err.Error())
}

// invokeSecondFunction
// @middleware-err
func (fs FirstStruct) invokeSecondFunction() {

	err := fmt.Errorf("a pre-defind error")
	err = fmt.Errorf("New error ")
	fmt.Println(err.Error())
}

// invokeThreeFunction
// @middleware-err
func invokeThreeFunction() {
	err := fmt.Errorf("a pre-defind error")
	err = fmt.Errorf("New error ")
	fmt.Println(err.Error())
}

// invokeFourFunction
// @middleware-err
func invokeFourFunction() func(err error) {
	err := fmt.Errorf("a pre-defind error")
	err = fmt.Errorf("New error ")

	if err != nil {
	}
	return func(err error) {
		fmt.Println(err.Error())
	}
}

// invokeFiveFunction
// @middleware-err
// this scene not support right now
func invokeFiveFunction() func() {
	return func() {
		err := fmt.Errorf("a pre-defind error")
		err = fmt.Errorf("New error ")

		if err != nil {
		}
	}
}
