# goAOP
Wrote by golang, designed for AOP.

## What is goAOP?

`goAOP` is a package for AOP(golang), it designed for use AOP easily.

When `goAOP` executes, it will scan all source code files. If these files contain valid aop id,
then `goAOP` will try to add pre-defined code in origin files.

## What id is a valid AOP id?

A valid AOP id standard is :

1. A word starts with '@'
2. This word can not end with '@'
3. This word is case-sensitive.

So the follow ids are valid:
1. @middleware-a
2. @middleware
3. @MiddleWare-C

also these follow ids are invalid:
1. middleware-a(not starts with '@')
2. @middleware-c@(ends with '@')
3. @(only has a char '@', not has valid word)

Please note that, since AOP id is case-sensitive, so `@middleware-a` not equals with `@Middleware-A`.

## How is goAOP work?

Let's take a demo code. Suppose we have a code snippet like bellow (fully code references unitTest dir):

```golang
package main

// InvokeSecondFunction demo code snippet
// @middleware-b
func InvokeSecondFunction() {}

```

When `goAOP` execute, it scans source code first. It will get a ast struct of source code. Second, goAOP try 
to parser all decls(decls is a slice contains all declare variable), typically, the comment also is a decl. So
goAOP will get all comments of `InvokeSecondFunction`。

Third, goAOP check whether it( comment of `InvokeSecondFunction`) has valid AOP id. If it finds valid AOP id, then try to 
match aop params from configure file. We suppose there has params of '@middleware-b' in configure file like bellow:

```toml
[[middleware]]
    id="@middleware-b"
    funcStmt=["""func(){
                    log.Println("middleware-b install")
                }()
            """]
    deferStmt=["""defer func(){
                    log.Println("middleware-b install")
                }()
            """]
    [[middleware.package]]
        name = "log"
        path = """ "github.com/sirupsen/logrus" """
```
Forth, goAOP will add `funcStmt` , `deferStmt` in function body by order. 

At last, goAOP also will add package in the head of origin file.

So we can get the last code like that:

```golang
package main

import log "github.com/sirupsen/logrus"

// InvokeSecondFunction demo code snippet
// @middleware-b
func InvokeSecondFunction() {
	func(){
		log.Println("middleware-b install")
	}()

	defer func(){
		log.Println("middleware-b install")
	}()
}
```

If you choose replace origin file, goAOP will cover origin file.

## How to build goAOP binary?

In this package, there has a sdk package and a main package. If you want to use goAOP directly, then you can build cli dir. 

Execute the follow command:

```shell
go build -o bin/aop cli/*
```

You will get an execute file like follow:

```golang
Usage of ./bin/aop:
  -config string
    	The runtime config, default is aop.toml (default "aop.toml")
  -debug
    	Enable / Disable debug output
  -dir string
    	The source code file dir path
  -replace
    	Replace source code file or not, default is true (default true)
```

Then execute `./bin/aop -config example/aop.toml -dir ./unitTests`, you will see the effective.

`goAOP` also supply a configure file , named aop.toml, in example dir. 

## How to use goAOP sdk?

`aops` is sdk dir, developer can invoke sdk in there's code. Since hard to understand go parser package, so developer can reference sdk usage from unit test code in aops dir.

## Some use cases.

+ [Use func as a depend condition](doc/case-01.md).