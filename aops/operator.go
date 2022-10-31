package aops

import "go/ast"

// There are all support operators.  These operators execute by bellow orders:
//
// 1. addDeferWithoutVarOperator
// 2. addFuncWithoutDependsOperator
// 3. addStmtAsFuncWithVarOperator
// 4. addStmtAsReturnOperator
// 5. addReturnWithBindVarOperator
// 6. addStmtBindVarOperator

// addReturnWithBindVarOperator Find the return function, then insert code in target function.
// The variable depend on contains basic type, like int, string, error etc. Also, support function type,
// like that func(e func()), the e is a func variable.
// Detail usage please reference `cases/insert-return-func-with-var` and `unitTests/test.go`
func addReturnWithBindVarOperator(t *ast.FuncDecl, def []ast.Stmt, depend []string) error {
	for _, _f := range t.Body.List {
		switch t := _f.(type) {
		case *ast.ReturnStmt:
			return _addFuncCodeWithVar(t, def, depend)
			//case
		}

	}

	return nil
}

// addFuncWithoutDependsOperator Insert expr that in the fun list to source code by order.
func addFuncWithoutDependsOperator(t *ast.FuncDecl, fun []ast.Expr) error {
	var stats []ast.Stmt
	for _, e := range fun {
		stats = append(stats, &ast.ExprStmt{
			X: e,
		})
	}

	stats = append(stats, t.Body.List...)
	t.Body.List = stats

	return nil
}

// addDeferWithoutVarOperator Insert stmt ignore any variable depend.
func addDeferWithoutVarOperator(t *ast.FuncDecl, def []ast.Stmt) error {
	var stats []ast.Stmt
	for _, e := range def {
		stats = append(stats, e)
	}

	stats = append(stats, t.Body.List...)
	t.Body.List = stats

	return nil
}

// addStmtAsFuncWithVarOperator Insert stmt with specify variable.
// Now only support specify one variable. If there has no depend on variable, it will do nothing.
func addStmtAsFuncWithVarOperator(t *ast.FuncDecl, def []ast.Stmt, depend []string) error {
	if len(depend) > 0 {
		return addStmtBlockBindVarOperator(t, []DeclParams{
			{
				VarName: depend[0],
				Stmt:    nil,
			},
		}, def)
	}

	return nil
}

// addStmtAsReturnOperator check whether this function has a func variable as return data.
// If it has function as return, then add pre-defined code. Otherwise, do nothing.
func addStmtAsReturnOperator(t *ast.FuncDecl, fun []ast.Stmt) error {
	for _, _f := range t.Body.List {
		rf, ok := _f.(*ast.ReturnStmt)
		if ok {
			//	One function only has one return stmt block. So no need check twice.
			return _addFuncCode(rf, fun)
		}
	}

	return nil
}

// _addFuncCodeWithVar Find the specific variable position. First find variable from params list,
// then find in body.
// If it finds variable in params, then it will insert all exprs in the head of function body.
// If it finds variable in body scope, then it will insert all exprs behind the variable.
func _addFuncCodeWithVar(t *ast.ReturnStmt, exprs []ast.Stmt, depend []string) error {
	for _, returnFunc := range t.Results {
		rf, ok := returnFunc.(*ast.FuncLit)
		if ok {
			findPosition := false
			// First check whether the variable depends on  defined in params
			for _, p := range rf.Type.Params.List {
				for _, pName := range p.Names {
					if pName.Name == depend[0] && findPosition == false {
						// The variable depends on find in params.
						// Insert stmt in body
						findPosition = true
					}
				}
			}

			if !findPosition {
				//	Try to find variable in body
				jump := false
				var _stmt []ast.Stmt
				for _, body := range rf.Body.List {
					as, ok := body.(*ast.AssignStmt)
					if ok {
						// jump == true means stmt block has inserted complete.
						// So ignore surplus stmts.
						if jump {
							_stmt = append(_stmt, body)
							continue
						}

						// check whether is the variable that we are finding.
						// x,y := 1, "ff"
						// lhs    rhs
						for _, lhs := range as.Lhs {
							if ident, ok := lhs.(*ast.Ident); ok {
								if ident.Name == depend[0] &&
									!jump {
									//	I find dp.VarName position. Then insert all stmt behind it.
									_stmt = append(_stmt, body)
									_stmt = append(_stmt, exprs...)
									jump = true
								}
							}
						}

						// The variable is not we are finding, so ignore it.
						if !jump {
							_stmt = append(_stmt, body)
						}
					} else {
						_stmt = append(_stmt, body)
					}
				}
				rf.Body.List = _stmt
				return nil
			}

			// Find variable in params list, so insert code in body.
			stats := make([]ast.Stmt, 0, len(rf.Body.List)+len(exprs))

			for _, e := range exprs {
				stats = append(stats, e)
			}

			stats = append(stats, rf.Body.List...)
			rf.Body.List = stats
		}
	}

	return nil
}

// _addFuncCode Insert all exprs in the head of return function body. The difference from _addFuncCodeWithVar
// is that this function no binding any variable. So _addFuncCode fit the closure scene.
func _addFuncCode(t *ast.ReturnStmt, exprs []ast.Stmt) error {

	for _, returnFunc := range t.Results {
		rf, ok := returnFunc.(*ast.FuncLit)
		if ok {
			// I think if it has a empty body also is OK.
			stats := make([]ast.Stmt, 0, len(rf.Body.List)+len(exprs))
			for _, e := range exprs {
				stats = append(stats, e)
			}

			stats = append(stats, rf.Body.List...)
			rf.Body.List = stats
		}
	}

	return nil
}

// addStmtBlockBindVarOperator Insert all stmts behind the specific variable. This function
// only support binding one variable.
func addStmtBlockBindVarOperator(t *ast.FuncDecl, v []DeclParams, stmt []ast.Stmt) error {
	if len(v) == 0 {
		return nil
	}

	dp := v[0]

	var _stmt []ast.Stmt
	var _stmtBlock []ast.Stmt = stmt

	jump := false

	for _, body := range t.Body.List {
		as, ok := body.(*ast.AssignStmt)
		if ok {
			// jump == true means stmt block has inserted complete.
			// So ignore surplus stmts.
			if jump {
				_stmt = append(_stmt, body)
				continue
			}

			// check whether is the variable that we are finding.
			// x,y := 1, "ff"
			// lhs    rhs
			for _, lhs := range as.Lhs {
				if ident, ok := lhs.(*ast.Ident); ok {
					if ident.Name == dp.VarName &&
						!jump {
						//	I find dp.VarName position. Then insert all stmt behind it.
						_stmt = append(_stmt, body)
						_stmt = append(_stmt, _stmtBlock...)
						jump = true
					}
				}
			}

			// The variable is not we are finding, so ignore it.
			if !jump {
				_stmt = append(_stmt, body)
			}
		} else {
			_stmt = append(_stmt, body)
		}
	}
	if len(_stmt) > 0 {
		t.Body.List = _stmt
	}

	return nil
}

// addStmtBindVar insert declare stmt behind specify variable.
// If v is nil, then do nothing.
// If v is not nil, then try to find the position of variable that
// ident by v[0].Name. Then insert all stmt that stores in v[0].Stmt.
// Return nil if there occur any unexpected error.
func addStmtBindVarOperator(t *ast.FuncDecl, v []DeclParams) error {
	if len(v) == 0 {
		return nil
	}

	dp := v[0]

	var _stmt []ast.Stmt
	var _stmtBlock []ast.Stmt

	// convert string slice to ast.Stmt
	for _, s := range dp.Stmt {
		_s, err := parserStmt(s)
		if err != nil {
			return err
		}

		_stmtBlock = append(_stmtBlock, _s)
	}

	jump := false

	for _, body := range t.Body.List {
		as, ok := body.(*ast.AssignStmt)
		if ok {
			// jump == true means stmt block has inserted complete.
			// So ignore surplus stmts.
			if jump {
				_stmt = append(_stmt, body)
				continue
			}

			// check whether is the variable that we are finding.
			// x,y := 1, "ff"
			// lhs    rhs
			for _, lhs := range as.Lhs {
				if ident, ok := lhs.(*ast.Ident); ok {
					if ident.Name == dp.VarName &&
						!jump {
						//	I find dp.VarName position. Then insert all stmt behind it.
						_stmt = append(_stmt, body)
						_stmt = append(_stmt, _stmtBlock...)
						jump = true
					}
				}
			}

			// The variable is not we are finding, so ignore it.
			if !jump {
				_stmt = append(_stmt, body)
			}
		} else {
			_stmt = append(_stmt, body)
		}
	}
	if len(_stmt) > 0 {
		t.Body.List = _stmt
	}

	return nil
}
