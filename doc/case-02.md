# Insert stmt without variable dependence

Assume we have some source codes as bellow:
```go
func test(){
	fmt.Println("Hello World")
}
```

Someone wants to insert some code and these code not depend on any variable, so we define aop configure like that:
```toml
[[middleware]]
    id="@middleware-without-variable"
    [[middleware.Stmt]]
        kind ="add-func-with-var-depend"
        code = [
            """
            x:=1
            """
        ]
```

At last, we can get the code:
```go
func test(){
    x:=1
	fmt.Println("Hello World")
}
```
