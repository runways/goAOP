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
// @middleware-return-err
func invokeFourFunction() func(err2 error) {
	err := fmt.Errorf("a pre-defind error")
	err = fmt.Errorf("New error ")

	if err != nil {
	}
	return func(err error) {
		str := "xx"
		fmt.Println(err.Error(), str)
	}
}

// invokeFiveFunction
// @middleware-err
func invokeFiveFunction() func() {
	return func() {
		str := "hello world"
		err := fmt.Errorf("a pre-defind error")
		err = fmt.Errorf("New error ")

		if err != nil {
			fmt.Println(str)
		}
	}
}

// invokeFiveFunction
// @middleware-err
func invokeSixFunction() func(error) {
	return func(err error) {
		str := "hello world"
		err = fmt.Errorf("a pre-defind error")
		err = fmt.Errorf("New error ")

		if err != nil {
			fmt.Println(str)
		}
	}
}

// invokeSevenFunction
// @middleware-func-var-err
func invokeSevenFunction() func(e func(error)) {
	return func(e func(error)) {
		e(fmt.Errorf("a func with error params"))
	}
}
