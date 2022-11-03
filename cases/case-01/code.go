package case_01

// There are some examples.
// 1. `invokeFunctionWithInjection` use a middleware with a param. When execute finish, the full code will be
//
//// invokeFunctionWithInjection
//// @middleware-injection(path:"/user/id")
//func (fs FirstStruct) invokeFunctionWithInjection() {
//	path := "/user/id"
//	fmt.Println("path: %s",
//
//		path)
//
//} more detail info please reference `internal_test.go`  TestAddCode unittest(A ID with param).
// 2. `invokeSecondFunctionWithInjection` use a default value param. When execute finish, the code will be :
//// invokeSecondFunctionWithInjection
//// @middleware-injection(path:"")
//func (fs FirstStruct) invokeSecondFunctionWithInjection() {
//	path := ""
//	fmt.Println("path: %s",
//
//		path)
//
//}
// 3. `invokeThirdFunctionWithInjection` use an int value param. When execute finish, the code will be:
//func (fs FirstStruct) invokeThirdFunctionWithInjection() {
//	path := 100
//	fmt.Println("path: %v",
//
//		path)
//
//}
// 4. `invokeFourFunctionWithInjection` use two params. When execute finish, the code will be:
//func (fs FirstStruct) invokeFourFunctionWithInjection() {
//	path := 100
//	name := "a string param"
//	fmt.Println("path: %v, name: %v",
//
//		path, name,
//	)
//
//}
// 5. `invokeFiveFunctionWithInjection` use two params and a inject param. When execute finish, the code will be:
//func (fs FirstStruct) invokeFiveFunctionWithInjection() {
//	path := 100
//	name := "a string param"
//	f := inject.
//		AOPLabel{Name: "invokeFiveFunctionWithInjection",
//		Owner: "FirstStruct"}
//	fmt.Println("path: %v, name: %v inject: %v",
//
//		path, name, f)
//
//}
// If user use inject id, please add `"github.com/runways/goAOP/inject"` in Pack block.

type FirstStruct struct {
	name string `json:"name"`
	Age  int    `json:"age"`
}

// invokeFunctionWithInjection
// @middleware-injection(path:"/user/id")
func (fs FirstStruct) invokeFunctionWithInjection() {}

// invokeSecondFunctionWithInjection
// @middleware-injection(path:"")
func (fs FirstStruct) invokeSecondFunctionWithInjection() {}

// invokeThirdFunctionWithInjection
// @middleware-injection(path:100)
func (fs FirstStruct) invokeThirdFunctionWithInjection() {}

// invokeFourFunctionWithInjection
// @middleware-injection(path:100, name:"a string param")
func (fs FirstStruct) invokeFourFunctionWithInjection() {}

// invokeFiveFunctionWithInjection
// @middleware-injection(path:100, name:"a string param", f:@inject)
func (fs FirstStruct) invokeFiveFunctionWithInjection() {}
