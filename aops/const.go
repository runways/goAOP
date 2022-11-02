package aops

const (
	AddFuncWithoutDepends = 1 + iota
	AddDeferFuncStmt
	AddDeferFuncWithVarStmt
	AddFuncWithVarStmt
	AddReturnFuncWithoutVarStmt
	AddReturnFuncWithVarStmt
	AddFuncWithoutDependsWithInject
)

const (
	AddFuncWithoutDependsStr           = "add-func-without-depends"
	AddFuncWithVarStmtStr              = "add-func-with-var-depend"
	AddFuncWithoutDependsWithInjectStr = "add-func-without-depends-with-injection"
	AddDeferFuncStmtStr                = "add-defer-func"
	AddDeferFuncWithVarStmtStr         = "add-defer-func-with-var-depend"
	AddReturnFuncWithoutVarStmtStr     = "add-return-func-without-var"
	AddReturnFuncWithVarStmtStr        = "add-return-func-with-var"
)
