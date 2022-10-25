package aops

type StmtParams struct {
	FunStmt   []string
	DeferStmt []string
	Packs     []Pack
}

type Pack struct {
	Name string
	Path string
}
