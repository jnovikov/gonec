package ast

import (
	"github.com/covrom/gonec/bincode/binstmt"
	"github.com/covrom/gonec/pos"
)

// Stmt provides all of interfaces for statement.
type Stmt interface {
	pos.Pos
	stmt()
	Simplify()
	BinTo(*binstmt.BinStmts, int, *int)
}

// StmtImpl provide commonly implementations for Stmt..
type StmtImpl struct {
	pos.PosImpl // StmtImpl provide Pos() function.
}

// stmt provide restraint interface.
func (x *StmtImpl) stmt() {}

type Stmts []Stmt

func (x Stmts) BinTo(bins *binstmt.BinStmts, reg int, lid *int) {
	for _, st := range x {
		st.BinTo(bins, reg, lid)
	}
}

// ExprStmt provide expression statement.
type ExprStmt struct {
	StmtImpl
	Expr Expr
}

func (x *ExprStmt) Simplify() {
	x.Expr = x.Expr.Simplify()
}

func (s *ExprStmt) BinTo(bins *binstmt.BinStmts, reg int, lid *int) {
	s.Expr.BinTo(bins, reg, lid, true)
	// *bins = append(*bins, addBinExpr(s.Expr, reg, lid, true)...)
}

// IfStmt provide "if/else" statement.
type IfStmt struct {
	StmtImpl
	If     Expr
	Then   Stmts
	ElseIf Stmts // This is array of IfStmt
	Else   Stmts
}

func (x *IfStmt) Simplify() {
	x.If = x.If.Simplify()
	for _, st := range x.Then {
		st.Simplify()
	}
	for _, st := range x.ElseIf {
		st.Simplify()
	}
	for _, st := range x.Else {
		st.Simplify()
	}
}

func (s *IfStmt) BinTo(bins *binstmt.BinStmts, reg int, lid *int) {
	*lid++
	lend := *lid

	// Если
	s.If.BinTo(bins, reg, lid, false)

	*lid++
	lf := *lid

	bins.Append(binstmt.NewBinJFALSE(reg, lf, s))

	// Тогда
	s.Then.BinTo(bins, reg, lid)

	bins.Append(binstmt.NewBinJMP(lend, s))

	// ИначеЕсли
	bins.Append(binstmt.NewBinLABEL(lf, s))

	for _, elif := range s.ElseIf {
		stmtif := elif.(*IfStmt)

		stmtif.If.BinTo(bins, reg, lid, false)

		// если ложь, то перейдем на следующее условие
		*lid++
		li := *lid

		bins.Append(binstmt.NewBinJFALSE(reg, li, stmtif))

		stmtif.Then.BinTo(bins, reg, lid)

		bins.Append(binstmt.NewBinJMP(lend, stmtif))

		bins.Append(binstmt.NewBinLABEL(li, stmtif))
	}

	// Иначе
	if len(s.Else) > 0 {
		s.Else.BinTo(bins, reg, lid)
	}
	// КонецЕсли
	bins.Append(binstmt.NewBinLABEL(lend, s))

	// освобождаем память
	bins.Append(binstmt.NewBinFREE(reg+1, s))
}

// TryStmt provide "try/catch/finally" statement.
type TryStmt struct {
	StmtImpl
	Try Stmts
	// Var     string
	Catch Stmts
	// Finally Stmts
}

func (x *TryStmt) Simplify() {
	for _, st := range x.Try {
		st.Simplify()
	}
	for _, st := range x.Catch {
		st.Simplify()
	}
}

func (s *TryStmt) BinTo(bins *binstmt.BinStmts, reg int, lid *int) {
	*lid++
	lend := *lid
	*lid++
	li := *lid
	// эта инструкция сообщает, в каком регистре будет отслеживаться ошибка выполнения кода до блока CATCH
	// по-умолчанию, ошибка в регистрах не отслеживается, а передается по уровням исполнения вирт. машины
	bins = appendBin(bins,
		&BinTRY{
			Reg:    reg,
			JumpTo: li,
		}, s)
	bins = append(bins, BinaryCode(s.Try, reg+1, lid).Code...) // чтобы не затереть регистр с ошибкой, увеличиваем номер
	// сюда переходим, если в блоке выше возникла ошибка
	bins = appendBin(bins,
		&BinLABEL{
			Label: li,
		}, s)
	// CATCH работает как JFALSE, и определяет функцию ОписаниеОшибки()
	bins = appendBin(bins,
		&BinCATCH{
			Reg:    reg,
			JumpTo: lend,
		}, s)
	// тело обработки ошибки
	bins = append(bins, BinaryCode(s.Catch, reg, lid).Code...) // регистр с ошибкой больше не нужен, текст определен функцией
	// КонецПопытки
	bins = appendBin(bins,
		&BinLABEL{
			Label: lend,
		}, s)
	// снимаем со стека состояние обработки ошибок, чтобы последующий код не был включен в текущую обработку
	bins = appendBin(bins,
		&BinPOPTRY{
			CatchLabel: li,
		}, s)
	// освобождаем память
	bins = appendBin(bins,
		&BinFREE{
			Reg: reg + 1,
		}, s)

}	
// ForStmt provide "for in" expression statement.
type ForStmt struct {
	StmtImpl
	Var   int //string
	Value Expr
	Stmts Stmts
}

func (x *ForStmt) Simplify() {
	x.Value = x.Value.Simplify()
	for _, st := range x.Stmts {
		st.Simplify()
	}
}

// NumForStmt name = expr1 to expr2
type NumForStmt struct {
	StmtImpl
	Name  int //string
	Expr1 Expr
	Expr2 Expr
	Stmts Stmts
}

func (x *NumForStmt) Simplify() {
	x.Expr1 = x.Expr1.Simplify()
	x.Expr2 = x.Expr2.Simplify()
	for _, st := range x.Stmts {
		st.Simplify()
	}
}

// CForStmt provide C-style "for (;;)" expression statement.
// type CForStmt struct {
// 	StmtImpl
// 	Expr1 Expr
// 	Expr2 Expr
// 	Expr3 Expr
// 	Stmts Stmts
// }

// LoopStmt provide "for expr" expression statement.
type LoopStmt struct {
	StmtImpl
	Expr  Expr
	Stmts Stmts
}

func (x *LoopStmt) Simplify() {
	x.Expr = x.Expr.Simplify()
	for _, st := range x.Stmts {
		st.Simplify()
	}
}

// BreakStmt provide "break" expression statement.
type BreakStmt struct {
	StmtImpl
}

func (x *BreakStmt) Simplify() {}

// ContinueStmt provide "continue" expression statement.
type ContinueStmt struct {
	StmtImpl
}

func (x *ContinueStmt) Simplify() {}

// ForStmt provide "return" expression statement.
type ReturnStmt struct {
	StmtImpl
	Exprs []Expr
}

func (x *ReturnStmt) Simplify() {
	for i := range x.Exprs {
		x.Exprs[i] = x.Exprs[i].Simplify()
	}
}

// ThrowStmt provide "throw" expression statement.
type ThrowStmt struct {
	StmtImpl
	Expr Expr
}

func (x *ThrowStmt) Simplify() {
	x.Expr = x.Expr.Simplify()
}

// ModuleStmt provide "module" expression statement.
type ModuleStmt struct {
	StmtImpl
	Name  int //string
	Stmts Stmts
}

func (x *ModuleStmt) Simplify() {
	for _, st := range x.Stmts {
		st.Simplify()
	}
}

// VarStmt provide statement to let variables in current scope.
type VarStmt struct {
	StmtImpl
	Names []int //string
	Exprs []Expr
}

func (x *VarStmt) Simplify() {
	for i := range x.Exprs {
		x.Exprs[i] = x.Exprs[i].Simplify()
	}
}

// SwitchStmt provide switch statement.
type SwitchStmt struct {
	StmtImpl
	Expr  Expr
	Cases Stmts
}

func (x *SwitchStmt) Simplify() {
	x.Expr = x.Expr.Simplify()
	for _, st := range x.Cases {
		st.Simplify()
	}
}

// SelectStmt provide switch statement.
type SelectStmt struct {
	StmtImpl
	Cases Stmts
}

func (x *SelectStmt) Simplify() {
	for _, st := range x.Cases {
		st.Simplify()
	}
}

// CaseStmt provide switch/case statement.
type CaseStmt struct {
	StmtImpl
	Expr  Expr
	Stmts Stmts
}

func (x *CaseStmt) Simplify() {
	x.Expr = x.Expr.Simplify()
	for _, st := range x.Stmts {
		st.Simplify()
	}
}

// DefaultStmt provide switch/default statement.
type DefaultStmt struct {
	StmtImpl
	Stmts Stmts
}

func (x *DefaultStmt) Simplify() {
	for _, st := range x.Stmts {
		st.Simplify()
	}
}

// LetsStmt provide multiple statement of let.
type LetsStmt struct {
	StmtImpl
	Lhss     []Expr
	Operator string
	Rhss     []Expr
}

func (x *LetsStmt) Simplify() {
	for i := range x.Lhss {
		x.Lhss[i] = x.Lhss[i].Simplify()
	}
	for i := range x.Rhss {
		x.Rhss[i] = x.Rhss[i].Simplify()
	}
}
