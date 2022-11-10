package aops

type OperationKind int

// StmtParams The stmt will insert into function body.
// There are three kind of stmt.
//
// FunStmt declares stander func stmt, e.g. func(){xxxx}()
// DeferStmt declares some stmt in defer block. e.g. defer func(){xxxx}()
// FunVarStmt declares some stmt in return function body, e.g.
//
// If FunVarStmt = []string{" fmt.Println("xxxx") "}, then the result is:
//
//	return func(){
//	  // Below stmts from FunVarStmt
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
//	    	stmt block
//		 }(err)
//
// `
// }
// Since we don't know the err variable position in source code, but we needn't care about it.
// `AddCode` will try to find the err declare position, then insert stmt bellow it. But one thing
// we should notice that DeclStmt now only support binding one variable. That means all stmts need use
// the same variable. For example, the previous stmt bind err is valid. If there has other stmts binding
// str(a new string variable). We can't find different variables at the same time, so we can not insert stmt right.
//
// At last, Packs save the import data. Maybe user has import the same package, so named a unique name
// for avoid repeat is a good idea.
type StmtParams struct {
	DeclStmt []DeclParams
	Stmts    []StmtParam
	Packs    []Pack
}

type Pack struct {
	Name string
	Path string
}

// DeclParams store stmt insert behind specify variable
type DeclParams struct {
	VarName  string // VarName is the variable name ,like 'x := 1', the x is var name.
	Stmt     []string
	FuncName string // FuncName is the func name, like 'x := fmt.Sprintf', the `fmt.Sprintf` is func name.
}

// StmtParam store the metadata of stmt.
// Kind decides to how and where to insert stmt.
// Stmt is the string of stmt, use parseStmt before use these.
// Depends are the dependence conditions
type StmtParam struct {
	Kind        OperationKind
	Stmt        []string
	Depends     []string
	FuncDepends []string
}

type StmtDepend interface {
	Depend() []string
}

// StmtVarDepend Variable dependency condition.
// VarName are all the variable string names.
// AOP will find all the variable init complete position
type StmtVarDepend struct {
	VarName []string
}

func (sd StmtVarDepend) Depend() []string {
	return sd.VarName
}
