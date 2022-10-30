package aops

// StmtParams The stmt will insert into function body.
// FunStmt declares stander func stmt, e.g. func(){xxxx}()
// DeferStmt declares some stmt in defer block. e.g. defer func(){xxxx}()
// FunVarStmt declares some stmt in return function body, e.g.
// If FunVarStmt = []string{" fmt.Println("xxxx") "}, then the result is:
//
//	return func(){
//		 // Below stmts from FunVarStmt
//	  fmt.Println("xxxx")
//	  // insert end
//	  xxxx
//	}
//
// DeclStmt declares some stmt bellow specify variable. like that,
// If we put some string stmt in []string{
// `
//
//		 defer func(err error){
//	    stmt block
//		 }(err)
//
// `
// }
// Since we don't know the err variable position in source code, but we needn't think about it.
// `AddCode` will try to find the err declare position, then insert stmt bellow it. But one thing
// we should notice that DeclStmt now only support bind one variable. That means more stmts need use
// the same variable. For example, the previous stmt bind err is valid. If there has other stmt bind
// x(new variable). We can not find different variable at the same time, so we can not insert stmt right.
type StmtParams struct {
	FunStmt    []string
	DeferStmt  []string
	FunVarStmt []string
	DeclStmt   []DeclParams
	Packs      []Pack
}

type Pack struct {
	Name string
	Path string
}

// DeclParams store stmt insert behind specify variable
type DeclParams struct {
	VarName string
	Stmt    []string
}
