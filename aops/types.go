package aops

// StmtParams The stmt will insert into function body.
// FunStmt declares stander func stmt, e.g. func(){xxxx}()
// DeferStmt declares some stmt in defer block. e.g. defer func(){xxxx}()
// FunVarStmt declares some stmt in return function body, e.g.
// If FunVarStmt = []string{" fmt.Println("xxxx") "}, then the result is:
// return func(){
//	 // Below stmts from FunVarStmt
//   fmt.Println("xxxx")
//   // insert end
//   xxxx
// }
type StmtParams struct {
	FunStmt    []string
	DeferStmt  []string
	FunVarStmt []string
	Packs      []Pack
}

type Pack struct {
	Name string
	Path string
}
