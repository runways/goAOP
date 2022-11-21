# Use func as a depend condition.

We assume the source code is :
```go
func (fs FirstStruct) addWithFuncDependWithInjection() {
	cli := sqlx.NewMysql(fs.dataSource)
	
	fmt.Println(cli)
}
```

Next we define the config like that:
```toml
[[middleware]]
    id="@job-mysql-health-check"
    [[middleware.Stmt]]
        kind = "add-func-with-var-depend"
        code = [
            """
            go _check.Mysql(__varName__)
            """
        ] 
        funDepend = ["sqlx.NewMysql"]
    [[middleware.package]]
        name = "_check"
        path = """ "github.com/mysql/check" """ # a fake mysql check library
```

At last, we can get the code. 
```go
func (fs FirstStruct) addWithFuncDependWithInjection() {
	cli := sqlx.NewMysql(fs.dataSource)
	go _check.Mysql(cli)
	
	fmt.Println(cli)
}
```