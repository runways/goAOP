package case_01

type FirstStruct struct {
	name string `json:"name"`
	Age  int    `json:"age"`
}

// invokeFunctionWithInjection
// @middleware-injection(path:"/user/id")
func (fs FirstStruct) invokeFunctionWithInjection() {

}
