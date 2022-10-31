package aops

const (
	AddFuncWithoutDepends = 1 + iota
	AddDeferFuncStmt
	AddDeferFuncWithVarStmt
	AddFuncWithVarStmt
	AddReturnFuncWithoutVarStmt
)

const (
	AddFuncWithoutDependsStr       = "add-func-without-depends"
	AddFuncWithVarStmtStr          = "add-func-with-var-depend"
	AddDeferFuncStmtStr            = "add-defer-func"
	AddDeferFuncWithVarStmtStr     = "add-defer-func-with-var-depend"
	AddReturnFuncWithoutVarStmtStr = "add-return-func-without-var"
)

//var stmtKindExecuteOrder = []int{
//	AddFuncWithoutDepends,
//	AddFuncWithVarStmt,
//	AddDeferFuncStmt,
//	AddDeferFuncWithVarStmt,
//}
