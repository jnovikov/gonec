package bincode

///////////////////////////////////////////////////////////////
// компиляция в байткод
///////////////////////////////////////////////////////////////

// func BinaryCode(inast ast.Stmts, reg int, lid *int) (bcd binstmt.BinCode) {
// 	bins := bcd.Code
// 	inast.BinTo(&bins, reg, lid)
// 	bcd.Code = bins
// 	bcd.MapLabels()
// }

// for _, st := range inast {

// 	// перебираем все подвыражения и команды, и выстраиваем их в линию
// 	// если в команде есть выражение - определяем новый id регистра, присваиваем ему выражение, а в команду передаем id этого регистра
// 	switch s := st.(type) {
// 	// case *ast.ExprStmt:
// 	// 	bins = append(bins, addBinExpr(s.Expr, reg, lid, true)...)
// 	// case *ast.IfStmt:
// 	// 	*lid++
// 	// 	lend := *lid
// 	// 	// Если
// 	// 	bins = append(bins, addBinExpr(s.If, reg, lid, false)...)
// 	// 	*lid++
// 	// 	lf := *lid
// 	// 	bins = appendBin(bins,
// 	// 		&BinJFALSE{
// 	// 			Reg:    reg,
// 	// 			JumpTo: lf,
// 	// 		}, s)
// 	// 	// Тогда
// 	// 	bins = append(bins, BinaryCode(s.Then, reg, lid).Code...)
// 	// 	bins = appendBin(bins,
// 	// 		&BinJMP{
// 	// 			JumpTo: lend,
// 	// 		}, s)
// 	// 	// ИначеЕсли
// 	// 	bins = appendBin(bins,
// 	// 		&BinLABEL{
// 	// 			Label: lf,
// 	// 		}, s)

// 	// 	for _, elif := range s.ElseIf {
// 	// 		stmtif := elif.(*ast.IfStmt)
// 	// 		bins = append(bins, addBinExpr(stmtif.If, reg, lid, false)...)
// 	// 		// если ложь, то перейдем на следующее условие
// 	// 		*lid++
// 	// 		li := *lid
// 	// 		bins = appendBin(bins,
// 	// 			&BinJFALSE{
// 	// 				Reg:    reg,
// 	// 				JumpTo: li,
// 	// 			}, stmtif)
// 	// 		bins = append(bins, BinaryCode(stmtif.Then, reg, lid).Code...)
// 	// 		bins = appendBin(bins,
// 	// 			&BinJMP{
// 	// 				JumpTo: lend,
// 	// 			}, stmtif)
// 	// 		bins = appendBin(bins,
// 	// 			&BinLABEL{
// 	// 				Label: li,
// 	// 			}, stmtif)
// 	// 	}
// 	// 	// Иначе
// 	// 	if len(s.Else) > 0 {
// 	// 		bins = append(bins, BinaryCode(s.Else, reg, lid).Code...)
// 	// 	}
// 	// 	// КонецЕсли
// 	// 	bins = appendBin(bins,
// 	// 		&BinLABEL{
// 	// 			Label: lend,
// 	// 		}, s)
// 	// 	// освобождаем память
// 	// 	bins = appendBin(bins,
// 	// 		&BinFREE{
// 	// 			Reg: reg + 1,
// 	// 		}, s)

// 	// case *ast.TryStmt:
// 	// 	*lid++
// 	// 	lend := *lid
// 	// 	*lid++
// 	// 	li := *lid
// 	// 	// эта инструкция сообщает, в каком регистре будет отслеживаться ошибка выполнения кода до блока CATCH
// 	// 	// по-умолчанию, ошибка в регистрах не отслеживается, а передается по уровням исполнения вирт. машины
// 	// 	bins = appendBin(bins,
// 	// 		&BinTRY{
// 	// 			Reg:    reg,
// 	// 			JumpTo: li,
// 	// 		}, s)
// 	// 	bins = append(bins, BinaryCode(s.Try, reg+1, lid).Code...) // чтобы не затереть регистр с ошибкой, увеличиваем номер
// 	// 	// сюда переходим, если в блоке выше возникла ошибка
// 	// 	bins = appendBin(bins,
// 	// 		&BinLABEL{
// 	// 			Label: li,
// 	// 		}, s)
// 	// 	// CATCH работает как JFALSE, и определяет функцию ОписаниеОшибки()
// 	// 	bins = appendBin(bins,
// 	// 		&BinCATCH{
// 	// 			Reg:    reg,
// 	// 			JumpTo: lend,
// 	// 		}, s)
// 	// 	// тело обработки ошибки
// 	// 	bins = append(bins, BinaryCode(s.Catch, reg, lid).Code...) // регистр с ошибкой больше не нужен, текст определен функцией
// 	// 	// КонецПопытки
// 	// 	bins = appendBin(bins,
// 	// 		&BinLABEL{
// 	// 			Label: lend,
// 	// 		}, s)
// 	// 	// снимаем со стека состояние обработки ошибок, чтобы последующий код не был включен в текущую обработку
// 	// 	bins = appendBin(bins,
// 	// 		&BinPOPTRY{
// 	// 			CatchLabel: li,
// 	// 		}, s)
// 	// 	// освобождаем память
// 	// 	bins = appendBin(bins,
// 	// 		&BinFREE{
// 	// 			Reg: reg + 1,
// 	// 		}, s)

// 	// case *ast.ForStmt:
// 	// 	// для каждого
// 	// 	bins = append(bins, addBinExpr(s.Value, reg, lid, false)...)

// 	// 	*lid++
// 	// 	lend := *lid
// 	// 	*lid++
// 	// 	li := *lid

// 	// 	regiter := reg + 1
// 	// 	regval := reg + 2
// 	// 	regsub := reg + 3
// 	// 	// инициализируем итератор, параметры цикла и цикл в стеке циклов
// 	// 	bins = appendBin(bins,
// 	// 		&BinFOREACH{
// 	// 			Reg:           reg,
// 	// 			RegIter:       regiter,
// 	// 			BreakLabel:    lend,
// 	// 			ContinueLabel: li,
// 	// 		}, s)
// 	// 	// очередная итерация
// 	// 	// сюда же переходим по Продолжить
// 	// 	bins = appendBin(bins,
// 	// 		&BinLABEL{
// 	// 			Label: li,
// 	// 		}, s)
// 	// 	bins = appendBin(bins,
// 	// 		&BinNEXT{
// 	// 			Reg:     reg,
// 	// 			RegIter: regiter,
// 	// 			RegVal:  regval,
// 	// 			JumpTo:  lend,
// 	// 		}, s)
// 	// 	// устанавливаем переменную-итератор
// 	// 	bins = appendBin(bins,
// 	// 		&BinSET{
// 	// 			Reg: regval,
// 	// 			Id:  s.Var,
// 	// 		}, s)

// 	// 	bins = append(bins, BinaryCode(s.Stmts, regsub, lid).Code...)

// 	// 	// повторяем итерацию
// 	// 	bins = appendBin(bins,
// 	// 		&BinJMP{
// 	// 			JumpTo: li,
// 	// 		}, s)

// 	// 	// КонецЦикла
// 	// 	bins = appendBin(bins,
// 	// 		&BinLABEL{
// 	// 			Label: lend,
// 	// 		}, s)
// 	// 	// снимаем со стека наличие цикла для Прервать и Продолжить
// 	// 	bins = appendBin(bins,
// 	// 		&BinPOPFOR{
// 	// 			ContinueLabel: li,
// 	// 		}, s)
// 	// 	// освобождаем память
// 	// 	bins = appendBin(bins,
// 	// 		&BinFREE{
// 	// 			Reg: reg + 1,
// 	// 		}, s)

// 	// case *ast.NumForStmt:
// 	// 	// для .. по ..
// 	// 	regfrom := reg + 1
// 	// 	regto := reg + 2
// 	// 	regsub := reg + 3

// 	// 	bins = append(bins, addBinExpr(s.Expr1, regfrom, lid, false)...)
// 	// 	bins = append(bins, addBinExpr(s.Expr2, regto, lid, false)...)

// 	// 	*lid++
// 	// 	lend := *lid
// 	// 	*lid++
// 	// 	li := *lid

// 	// 	// инициализируем итератор, параметры цикла и цикл в стеке циклов
// 	// 	bins = appendBin(bins,
// 	// 		&BinFORNUM{
// 	// 			Reg:           reg,
// 	// 			RegFrom:       regfrom,
// 	// 			RegTo:         regto,
// 	// 			BreakLabel:    lend,
// 	// 			ContinueLabel: li,
// 	// 		}, s)
// 	// 	// очередная итерация
// 	// 	// сюда же переходим по Продолжить
// 	// 	bins = appendBin(bins,
// 	// 		&BinLABEL{
// 	// 			Label: li,
// 	// 		}, s)
// 	// 	bins = appendBin(bins,
// 	// 		&BinNEXTNUM{
// 	// 			Reg:     reg,
// 	// 			RegFrom: regfrom,
// 	// 			RegTo:   regto,
// 	// 			JumpTo:  lend, // сюда же переходим по Прервать
// 	// 		}, s)
// 	// 	// устанавливаем переменную-итератор
// 	// 	bins = appendBin(bins,
// 	// 		&BinSET{
// 	// 			Reg: reg,
// 	// 			Id:  s.Name,
// 	// 		}, s)

// 	// 	bins = append(bins, BinaryCode(s.Stmts, regsub, lid).Code...)

// 	// 	// повторяем итерацию
// 	// 	bins = appendBin(bins,
// 	// 		&BinJMP{
// 	// 			JumpTo: li,
// 	// 		}, s)

// 	// 	// КонецЦикла
// 	// 	bins = appendBin(bins,
// 	// 		&BinLABEL{
// 	// 			Label: lend,
// 	// 		}, s)
// 	// 	// снимаем со стека наличие цикла для Прервать и Продолжить
// 	// 	bins = appendBin(bins,
// 	// 		&BinPOPFOR{
// 	// 			ContinueLabel: li,
// 	// 		}, s)
// 	// 	// освобождаем память
// 	// 	bins = appendBin(bins,
// 	// 		&BinFREE{
// 	// 			Reg: reg + 1,
// 	// 		}, s)

// 	// case *ast.LoopStmt:
// 	// 	*lid++
// 	// 	lend := *lid
// 	// 	*lid++
// 	// 	li := *lid
// 	// 	bins = appendBin(bins,
// 	// 		&BinWHILE{
// 	// 			BreakLabel:    lend,
// 	// 			ContinueLabel: li,
// 	// 		}, s)
// 	// 	// очередная итерация
// 	// 	// сюда же переходим по Продолжить
// 	// 	bins = appendBin(bins,
// 	// 		&BinLABEL{
// 	// 			Label: li,
// 	// 		}, s)
// 	// 	bins = append(bins, addBinExpr(s.Expr, reg, lid, false)...)
// 	// 	bins = appendBin(bins,
// 	// 		&BinJFALSE{
// 	// 			Reg:    reg,
// 	// 			JumpTo: lend,
// 	// 		}, s)
// 	// 	// тело цикла
// 	// 	bins = append(bins, BinaryCode(s.Stmts, reg+1, lid).Code...)

// 	// 	// повторяем итерацию
// 	// 	bins = appendBin(bins,
// 	// 		&BinJMP{
// 	// 			JumpTo: li,
// 	// 		}, s)

// 	// 	// КонецЦикла
// 	// 	bins = appendBin(bins,
// 	// 		&BinLABEL{
// 	// 			Label: lend,
// 	// 		}, s)
// 	// 	// снимаем со стека наличие цикла для Прервать и Продолжить
// 	// 	bins = appendBin(bins,
// 	// 		&BinPOPFOR{
// 	// 			ContinueLabel: li,
// 	// 		}, s)
// 	// 	// освобождаем память
// 	// 	bins = appendBin(bins,
// 	// 		&BinFREE{
// 	// 			Reg: reg + 1,
// 	// 		}, s)

// 	// case *ast.BreakStmt:
// 	// 	bins = appendBin(bins,
// 	// 		&BinBREAK{}, s)

// 	// case *ast.ContinueStmt:
// 	// 	bins = appendBin(bins,
// 	// 		&BinCONTINUE{}, s)

// 	// case *ast.ReturnStmt:
// 	// 	if len(s.Exprs) == 0 {
// 	// 		bins = appendBin(bins,
// 	// 			&BinLOAD{
// 	// 				Reg: reg, // основной регистр
// 	// 				Val: nil,
// 	// 			}, s)
// 	// 	}
// 	// 	if len(s.Exprs) == 1 {
// 	// 		// одиночное значение в reg
// 	// 		bins = append(bins, addBinExpr(s.Exprs[0], reg, lid, false)...)
// 	// 	} else {
// 	// 		// создание слайса в reg
// 	// 		bins = appendBin(bins,
// 	// 			&BinMAKESLICE{
// 	// 				Reg: reg,
// 	// 				Len: len(s.Exprs),
// 	// 				Cap: len(s.Exprs),
// 	// 			}, s)

// 	// 		for i, ee := range s.Exprs {
// 	// 			bins = append(bins, addBinExpr(ee, reg+1, lid, false)...)
// 	// 			bins = appendBin(bins,
// 	// 				&BinSETIDX{
// 	// 					Reg:    reg,
// 	// 					Index:  i,
// 	// 					RegVal: reg + 1,
// 	// 				}, ee)
// 	// 		}
// 	// 	}
// 	// 	// в reg имеем значение или структуру возврата
// 	// 	bins = appendBin(bins,
// 	// 		&BinFREE{
// 	// 			Reg: reg + 1,
// 	// 		}, s)
// 	// 	bins = appendBin(bins,
// 	// 		&BinRET{
// 	// 			Reg: reg,
// 	// 		}, s)

// 	// case *ast.ThrowStmt:
// 	// 	bins = append(bins, addBinExpr(s.Expr, reg, lid, false)...)
// 	// 	bins = appendBin(bins,
// 	// 		&BinTHROW{
// 	// 			Reg: reg,
// 	// 		}, s)

// 	// case *ast.ModuleStmt:
// 	// 	if s.Name == env.UniqueNames.Set("_") {
// 	// 		// добавляем все операторы в текущий контекст
// 	// 		bins = append(bins, BinaryCode(s.Stmts, reg, lid).Code...)
// 	// 	} else {
// 	// 		bins = appendBin(bins,
// 	// 			&BinMODULE{
// 	// 				Name: s.Name,
// 	// 				Code: BinaryCode(s.Stmts, 0, lid),
// 	// 			}, s)
// 	// 	}
// 	// case *ast.SwitchStmt:
// 	// 	bins = append(bins, addBinExpr(s.Expr, reg, lid, true)...)
// 	// 	// сравниваем с каждым case
// 	// 	*lid++
// 	// 	lend := *lid
// 	// 	var default_stmt *ast.DefaultStmt
// 	// 	for _, ss := range s.Cases {
// 	// 		if ssd, ok := ss.(*ast.DefaultStmt); ok {
// 	// 			default_stmt = ssd
// 	// 			continue
// 	// 		}
// 	// 		*lid++
// 	// 		li := *lid
// 	// 		case_stmt := ss.(*ast.CaseStmt)
// 	// 		bins = append(bins, addBinExpr(case_stmt.Expr, reg+1, lid, false)...)
// 	// 		bins = appendBin(bins,
// 	// 			&BinEQUAL{
// 	// 				Reg:  reg + 2,
// 	// 				Reg1: reg,
// 	// 				Reg2: reg + 1,
// 	// 			}, case_stmt)

// 	// 		bins = appendBin(bins,
// 	// 			&BinJFALSE{
// 	// 				Reg:    reg + 2,
// 	// 				JumpTo: li,
// 	// 			}, case_stmt)
// 	// 		bins = append(bins, BinaryCode(case_stmt.Stmts, reg, lid).Code...)
// 	// 		bins = appendBin(bins,
// 	// 			&BinJMP{
// 	// 				JumpTo: lend,
// 	// 			}, case_stmt)

// 	// 		bins = appendBin(bins,
// 	// 			&BinLABEL{
// 	// 				Label: li,
// 	// 			}, case_stmt)
// 	// 	}
// 	// 	if default_stmt != nil {
// 	// 		bins = append(bins, BinaryCode(default_stmt.Stmts, reg, lid).Code...)
// 	// 	}
// 	// 	bins = appendBin(bins,
// 	// 		&BinLABEL{
// 	// 			Label: lend,
// 	// 		}, s)
// 	// 	// освобождаем память
// 	// 	bins = appendBin(bins,
// 	// 		&BinFREE{
// 	// 			Reg: reg + 1,
// 	// 		}, s)

// 	// case *ast.SelectStmt:
// 	// 	*lid++
// 	// 	lstart := *lid
// 	// 	bins = appendBin(bins,
// 	// 		&BinLABEL{
// 	// 			Label: lstart,
// 	// 		}, s)

// 	// 	*lid++
// 	// 	lend := *lid
// 	// 	var default_stmt *ast.DefaultStmt
// 	// 	for _, ss := range s.Cases {
// 	// 		if ssd, ok := ss.(*ast.DefaultStmt); ok {
// 	// 			default_stmt = ssd
// 	// 			continue
// 	// 		}
// 	// 		*lid++
// 	// 		li := *lid
// 	// 		case_stmt := ss.(*ast.CaseStmt)
// 	// 		e, ok := case_stmt.Expr.(*ast.ChanExpr)
// 	// 		if !ok {
// 	// 			panic(NewStringError(case_stmt, "При выборе вариантов из каналов допустимы только выражения с каналами"))
// 	// 		}
// 	// 		// определяем значение справа
// 	// 		bins = append(bins, addBinExpr(e.Rhs, reg, lid, false)...)
// 	// 		if e.Lhs == nil {
// 	// 			// слева нет значения - это временное чтение из канала без сохранения значения в переменной
// 	// 			bins = appendBin(bins,
// 	// 				&BinTRYRECV{
// 	// 					Reg:       reg,
// 	// 					RegVal:    reg + 1,
// 	// 					RegOk:     reg + 2,
// 	// 					RegClosed: reg + 3,
// 	// 				}, e.Rhs)
// 	// 			// если канал закрыт или не получено значение - идем в следующую ветку
// 	// 			bins = appendBin(bins,
// 	// 				&BinJFALSE{
// 	// 					Reg:    reg + 2,
// 	// 					JumpTo: li,
// 	// 				}, s)
// 	// 		} else {
// 	// 			// значение слева
// 	// 			bins = append(bins, addBinExpr(e.Lhs, reg+1, lid, false)...)

// 	// 			// проверяем: слева канал?
// 	// 			bins = appendBin(bins,
// 	// 				&BinMV{
// 	// 					RegFrom: reg + 1,
// 	// 					RegTo:   reg + 3,
// 	// 				}, e)
// 	// 			bins = appendBin(bins,
// 	// 				&BinISKIND{
// 	// 					Reg:  reg + 3,
// 	// 					Kind: reflect.Chan,
// 	// 				}, e)

// 	// 			*lid++
// 	// 			li3 := *lid

// 	// 			bins = appendBin(bins,
// 	// 				&BinJFALSE{
// 	// 					Reg:    reg + 3,
// 	// 					JumpTo: li3,
// 	// 				}, e)

// 	// 			// слева канал - пишем в него правое

// 	// 			bins = appendBin(bins,
// 	// 				&BinTRYSEND{
// 	// 					Reg:    reg + 1,
// 	// 					RegVal: reg,
// 	// 					RegOk:  reg + 2,
// 	// 					// RegClosed: reg + 3,
// 	// 				}, e.Lhs)

// 	// 			*lid++
// 	// 			li2 := *lid

// 	// 			// если отправлено значение - выполняем код блока
// 	// 			bins = appendBin(bins,
// 	// 				&BinJTRUE{
// 	// 					Reg:    reg + 2,
// 	// 					JumpTo: li2,
// 	// 				}, s)
// 	// 			// если не отправлено значение - идем в следующую ветку
// 	// 			// если канал закрыт - будет паника
// 	// 			bins = appendBin(bins,
// 	// 				&BinJMP{
// 	// 					JumpTo: li,
// 	// 				}, s)

// 	// 			// иначе справа канал, а слева переменная (установим, если прочитали из канала)
// 	// 			bins = appendBin(bins,
// 	// 				&BinLABEL{
// 	// 					Label: li3,
// 	// 				}, s)

// 	// 			bins = appendBin(bins,
// 	// 				&BinTRYRECV{
// 	// 					Reg:       reg,
// 	// 					RegVal:    reg + 1,
// 	// 					RegOk:     reg + 2,
// 	// 					RegClosed: reg + 3,
// 	// 				}, e.Rhs)

// 	// 			// если канал закрыт или не получено значение - идем в следующую ветку
// 	// 			bins = appendBin(bins,
// 	// 				&BinJFALSE{
// 	// 					Reg:    reg + 2,
// 	// 					JumpTo: li,
// 	// 				}, s)

// 	// 			// устанавливаем переменную прочитанным значением
// 	// 			bins = append(bins, addBinLetExpr(e.Lhs, reg+1, lid)...)

// 	// 			bins = appendBin(bins,
// 	// 				&BinLABEL{
// 	// 					Label: li2,
// 	// 				}, s)

// 	// 		}
// 	// 		// отправили или прочитали - выполняем ветку кода и выходим из цикла
// 	// 		bins = append(bins, BinaryCode(case_stmt.Stmts, reg, lid).Code...)

// 	// 		// выходим из цикла
// 	// 		bins = appendBin(bins,
// 	// 			&BinJMP{
// 	// 				JumpTo: lend,
// 	// 			}, case_stmt)

// 	// 		// к следующему case
// 	// 		bins = appendBin(bins,
// 	// 			&BinLABEL{
// 	// 				Label: li,
// 	// 			}, s)
// 	// 	}
// 	// 	// если ни одна из веток не сработала - проверяем default
// 	// 	if default_stmt != nil {
// 	// 		bins = append(bins, BinaryCode(default_stmt.Stmts, reg, lid).Code...)
// 	// 	} else {
// 	// 		// допускаем обработку других горутин
// 	// 		bins = appendBin(bins,
// 	// 			&BinGOSHED{}, s)
// 	// 		bins = appendBin(bins,
// 	// 			&BinJMP{
// 	// 				JumpTo: lstart,
// 	// 			}, s)
// 	// 	}
// 	// 	bins = appendBin(bins,
// 	// 		&BinLABEL{
// 	// 			Label: lend,
// 	// 		}, s)
// 	// 	// освобождаем память
// 	// 	bins = appendBin(bins,
// 	// 		&BinFREE{
// 	// 			Reg: reg + 1,
// 	// 		}, s)

// 	// case *ast.LetsStmt:
// 	// 	// если справа одно выражение - присваиваем его всем левым
// 	// 	// и если там массив, то по очереди элементы, начиная с 0-го
// 	// 	// иначе с обеих сторон должно быть одинаковое число выражений, они попарно присваиваются
// 	// 	if len(s.Rhss) == 1 && len(s.Lhss) > 1 {
// 	// 		bins = append(bins, addBinExpr(s.Rhss[0], reg, lid, false)...)
// 	// 		// проверяем на массив
// 	// 		*lid++
// 	// 		lend := *lid
// 	// 		*lid++
// 	// 		li := *lid
// 	// 		bins = appendBin(bins,
// 	// 			&BinISSLICE{
// 	// 				Reg:     reg,
// 	// 				RegBool: reg + 1,
// 	// 			}, s)
// 	// 		bins = appendBin(bins,
// 	// 			&BinJFALSE{
// 	// 				Reg:    reg + 1,
// 	// 				JumpTo: li,
// 	// 			}, s)

// 	// 		// присваиваем из слайса
// 	// 		i := 0
// 	// 		for _, e := range s.Lhss {
// 	// 			// в рег+1 сохраним очередной элемент
// 	// 			bins = appendBin(bins,
// 	// 				&BinMV{
// 	// 					RegFrom: reg,
// 	// 					RegTo:   reg + 1,
// 	// 				}, e)
// 	// 			bins = appendBin(bins,
// 	// 				&BinLOAD{
// 	// 					Reg: reg + 2,
// 	// 					Val: i,
// 	// 				}, e)
// 	// 			bins = appendBin(bins,
// 	// 				&BinGETIDX{
// 	// 					Reg:      reg + 1,
// 	// 					RegIndex: reg + 2,
// 	// 				}, e)

// 	// 			bins = append(bins, addBinLetExpr(e, reg+1, lid)...)
// 	// 			i++
// 	// 		}

// 	// 		bins = appendBin(bins,
// 	// 			&BinJMP{
// 	// 				JumpTo: lend,
// 	// 			}, s)

// 	// 		// присваиваем одно и то же значение
// 	// 		bins = appendBin(bins,
// 	// 			&BinLABEL{
// 	// 				Label: li,
// 	// 			}, s)
// 	// 		for _, e := range s.Lhss {
// 	// 			bins = append(bins, addBinLetExpr(e, reg, lid)...)
// 	// 		}

// 	// 		bins = appendBin(bins,
// 	// 			&BinLABEL{
// 	// 				Label: lend,
// 	// 			}, s)
// 	// 	} else {
// 	// 		if len(s.Lhss) == len(s.Rhss) {
// 	// 			// сначала все вычисляем в разные регистры, затем все присваиваем
// 	// 			// так обеспечиваем взаимный обмен
// 	// 			for i := range s.Lhss {
// 	// 				bins = append(bins, addBinExpr(s.Rhss[i], reg+i, lid, false)...)
// 	// 			}
// 	// 			for i, e := range s.Lhss {
// 	// 				bins = append(bins, addBinLetExpr(e, reg+i, lid)...)
// 	// 			}
// 	// 		} else {
// 	// 			// ошибка
// 	// 			panic(NewStringError(s, "Количество переменных и значений должно совпадать или значение должно быть одно"))
// 	// 		}
// 	// 	}
// 	// 	// освобождаем память
// 	// 	// bins = appendBin(bins,
// 	// 	// 	&BinFREE{
// 	// 	// 		Reg: reg + 1,
// 	// 	// 	}, s)

// 	// case *ast.VarStmt:
// 	// 	// если справа одно выражение - присваиваем его всем левым
// 	// 	// иначе с обеих сторон должно быть одинаковое число выражений, они попарно присваиваются
// 	// 	if len(s.Exprs) == 1 {
// 	// 		bins = append(bins, addBinExpr(s.Exprs[0], reg, lid, false)...)
// 	// 		for _, e := range s.Names {
// 	// 			bins = appendBin(bins,
// 	// 				&BinSET{
// 	// 					Reg: reg,
// 	// 					Id:  e,
// 	// 				}, s)
// 	// 		}
// 	// 	} else {
// 	// 		if len(s.Exprs) == len(s.Names) {
// 	// 			for i, e := range s.Exprs {
// 	// 				bins = append(bins, addBinExpr(e, reg, lid, false)...)
// 	// 				bins = appendBin(bins,
// 	// 					&BinSET{
// 	// 						Reg: reg,
// 	// 						Id:  s.Names[i],
// 	// 					}, s)
// 	// 			}
// 	// 		} else {
// 	// 			// ошибка
// 	// 			panic(NewStringError(s, "Количество переменных и значений должно совпадать или значение должно быть одно"))
// 	// 		}
// 	// 	}
// 	// 	// освобождаем память
// 	// 	// bins = appendBin(bins,
// 	// 	// 	&BinFREE{
// 	// 	// 		Reg: reg + 1,
// 	// 	// 	}, s)

// 	}
// }
// return
// }

// func appendBin(bins BinStmts, b BinStmt, e pos.Pos) BinStmts {
// 	b.SetPosition(e.Position())
// 	return append(bins, b)
// }

// func addBinLetExpr(e ast.Expr, reg int, lid *int) (bins BinStmts) {
// 	// присваиваем значению переменной из e значение из регистра reg
// 	switch ee := e.(type) {
// case *ast.IdentExpr:
// 	bins = appendBin(bins,
// 		&BinSET{
// 			Reg: reg,
// 			Id:  ee.Id,
// 		}, e)

// case *ast.MemberExpr:
// 	bins = append(bins, addBinExpr(ee.Expr, reg+1, lid, false)...)
// 	bins = appendBin(bins,
// 		&BinSETMEMBER{
// 			Reg:    reg + 1,
// 			Id:     ee.Name,
// 			RegVal: reg,
// 		}, e)

// case *ast.ItemExpr:
// 	*lid++
// 	lend := *lid
// 	bins = append(bins, addBinExpr(ee.Value, reg+1, lid, false)...)
// 	bins = append(bins, addBinExpr(ee.Index, reg+2, lid, false)...)
// 	bins = appendBin(bins,
// 		&BinSETITEM{
// 			Reg:        reg + 1,
// 			RegIndex:   reg + 2,
// 			RegVal:     reg,
// 			RegNeedLet: reg + 3,
// 		}, e)
// 	bins = appendBin(bins,
// 		&BinJFALSE{
// 			Reg:    reg + 3,
// 			JumpTo: lend,
// 		}, ee)

// 	bins = append(bins, addBinLetExpr(ee.Value, reg+1, lid)...)

// 	bins = appendBin(bins,
// 		&BinLABEL{
// 			Label: lend,
// 		}, ee)

// case *ast.SliceExpr:
// 	*lid++
// 	lend := *lid
// 	bins = append(bins, addBinExpr(ee.Value, reg+1, lid, false)...)
// 	bins = append(bins, addBinExpr(ee.Begin, reg+2, lid, false)...)
// 	bins = append(bins, addBinExpr(ee.End, reg+3, lid, false)...)
// 	bins = appendBin(bins,
// 		&BinSETSLICE{
// 			Reg:        reg + 1,
// 			RegBegin:   reg + 2,
// 			RegEnd:     reg + 3,
// 			RegVal:     reg,
// 			RegNeedLet: reg + 4,
// 		}, e)
// 	bins = appendBin(bins,
// 		&BinJFALSE{
// 			Reg:    reg + 4,
// 			JumpTo: lend,
// 		}, ee)

// 	bins = append(bins, addBinLetExpr(ee.Value, reg+1, lid)...)

// 	bins = appendBin(bins,
// 		&BinLABEL{
// 			Label: lend,
// 		}, ee)

// 	default:
// 		// ошибка
// 		panic(NewStringError(e, "Неверная операция"))
// 	}
// 	return
// }

// func addBinExpr(expr ast.Expr, reg int, lid *int, inStmt bool) (bins BinStmts) {
//inStmt=true - признак запуска выражения как опреатора в блоке кода, иначе это подвыражение
// if expr == nil {
// 	bins = appendBin(bins,
// 		&BinLOAD{
// 			Reg: reg,
// 			Val: nil,
// 		}, &ast.NativeExpr{Value: nil}) // т.к. expr == nil, то у него нет Pos
// 	return
// }
// switch e := expr.(type) {
// case *ast.NativeExpr:
// 	// добавляем команду загрузки значения
// 	bins = appendBin(bins,
// 		&BinLOAD{
// 			Reg: reg, // основной регистр
// 			Val: e.Value,
// 		}, e)
// case *ast.NumberExpr:
// 	// команда на загрузку строки в регистр и ее преобразование в число, в регистр
// 	bins = appendBin(bins,
// 		&BinLOAD{
// 			Reg: reg,
// 			Val: e.Lit,
// 		}, e)

// 	bins = appendBin(bins,
// 		&BinCASTNUM{
// 			Reg: reg,
// 		}, e)
// case *ast.StringExpr:
// 	bins = appendBin(bins,
// 		&BinLOAD{
// 			Reg: reg,
// 			Val: e.Lit,
// 		}, e)
// case *ast.ConstExpr:
// 	b := BinLOAD{
// 		Reg: reg,
// 	}
// 	switch strings.ToLower(e.Value) {
// 	case "истина":
// 		b.Val = true
// 	case "ложь":
// 		b.Val = false
// 	case "null":
// 		b.Val = ast.NullVar
// 	default:
// 		b.Val = nil
// 	}
// 	bins = appendBin(bins, &b, e)
// case *ast.ArrayExpr:
// 	// создание слайса
// 	bins = appendBin(bins,
// 		&BinMAKESLICE{
// 			Reg: reg,
// 			Len: len(e.Exprs),
// 			Cap: len(e.Exprs),
// 		}, e)

// 	for i, ee := range e.Exprs {
// 		// каждое выражение сохраняем в следующем по номеру регистре (относительно регистра слайса)
// 		bins = append(bins, addBinExpr(ee, reg+1, lid, false)...)
// 		bins = appendBin(bins,
// 			&BinSETIDX{
// 				Reg:    reg,
// 				Index:  i,
// 				RegVal: reg + 1,
// 			}, ee)
// 	}
// case *ast.MapExpr:
// 	// создание мапы
// 	bins = appendBin(bins,
// 		&BinMAKEMAP{
// 			Reg: reg,
// 			Len: len(e.MapExpr),
// 		}, e)

// 	for k, ee := range e.MapExpr {
// 		bins = append(bins, addBinExpr(ee, reg+1, lid, false)...)
// 		bins = appendBin(bins,
// 			&BinSETKEY{
// 				Reg:    reg,
// 				Key:    k,
// 				RegVal: reg + 1,
// 			}, ee)
// 	}
// case *ast.IdentExpr:
// 	bins = appendBin(bins,
// 		&BinGET{
// 			Reg: reg,
// 			Id:  e.Id,
// 		}, e)
// case *ast.UnaryExpr:
// 	bins = append(bins, addBinExpr(e.Expr, reg, lid, false)...)
// 	bins = appendBin(bins,
// 		&BinUNARY{
// 			Reg: reg,
// 			Op:  rune(e.Operator[0]),
// 		}, e)
// case *ast.AddrExpr:
// 	switch ee := e.Expr.(type) {
// 	case *ast.IdentExpr:
// 		bins = appendBin(bins,
// 			&BinADDRID{
// 				Reg:  reg,
// 				Name: ee.Id,
// 			}, e)
// 	case *ast.MemberExpr:
// 		bins = append(bins, addBinExpr(ee.Expr, reg, lid, false)...)
// 		bins = appendBin(bins,
// 			&BinADDRMBR{
// 				Reg:  reg,
// 				Name: ee.Name,
// 			}, e)
// 	default:
// 		panic(NewStringError(e, "Неверная операция над значением"))
// 	}
// case *ast.DerefExpr:
// 	switch ee := e.Expr.(type) {
// 	case *ast.IdentExpr:
// 		bins = appendBin(bins,
// 			&BinUNREFID{
// 				Reg:  reg,
// 				Name: ee.Id,
// 			}, e)
// 	case *ast.MemberExpr:
// 		bins = append(bins, addBinExpr(ee.Expr, reg, lid, false)...)
// 		bins = appendBin(bins,
// 			&BinUNREFMBR{
// 				Reg:  reg,
// 				Name: ee.Name,
// 			}, e)
// 	default:
// 		panic(NewStringError(e, "Неверная операция над значением"))
// 	}
// case *ast.ParenExpr:
// 	bins = append(bins, addBinExpr(e.SubExpr, reg, lid, false)...)
// case *ast.BinOpExpr:
// 	oper := OperMap[e.Operator]
// 	// если это равенство в контексте исполнения блока кода, то это присваивание, а не вычисление выражения
// 	if inStmt && oper == EQL {
// 		bins = append(bins, BinaryCode([]ast.Stmt{
// 			&ast.LetsStmt{
// 				Lhss:     e.Lhss,
// 				Operator: "=",
// 				Rhss:     e.Rhss,
// 			},
// 		}, reg, lid).Code...)
// 		return
// 	}
// 	if len(e.Lhss) != 1 || len(e.Rhss) != 1 {
// 		panic(NewStringError(e, "С каждой стороны операции может быть только одно выражение"))
// 	}
// 	// сначала вычисляем левую часть
// 	bins = append(bins, addBinExpr(e.Lhss[0], reg, lid, false)...)
// 	switch oper {
// 	case LOR:
// 		*lid++
// 		lab := *lid
// 		// вставляем проверку на истину слева и возвращаем ее, не вычисляя правую часть, иначе возвращаем правую часть
// 		bins = appendBin(bins,
// 			&BinJTRUE{
// 				Reg:    reg,
// 				JumpTo: lab,
// 			}, e)
// 		bins = append(bins, addBinExpr(e.Rhss[0], reg, lid, false)...)
// 		bins = appendBin(bins,
// 			&BinLABEL{
// 				Label: lab,
// 			}, e)
// 	case LAND:
// 		*lid++
// 		lab := *lid
// 		// вставляем проверку на ложь слева и возвращаем ее, не вычисляя правую часть, иначе возвращаем правую часть
// 		bins = appendBin(bins,
// 			&BinJFALSE{
// 				Reg:    reg,
// 				JumpTo: lab,
// 			}, e)
// 		bins = append(bins, addBinExpr(e.Rhss[0], reg, lid, false)...)
// 		bins = appendBin(bins,
// 			&BinLABEL{
// 				Label: lab,
// 			}, e)
// 	default:
// 		bins = append(bins, addBinExpr(e.Rhss[0], reg+1, lid, false)...)
// 		bins = appendBin(bins,
// 			&BinOPER{
// 				RegL: reg, // сюда же помещается результат
// 				RegR: reg + 1,
// 				Op:   oper,
// 			}, e)
// 	}
// case *ast.TernaryOpExpr:
// 	bins = append(bins, addBinExpr(e.Expr, reg, lid, false)...)
// 	*lid++
// 	lab := *lid
// 	bins = appendBin(bins,
// 		&BinJFALSE{
// 			Reg:    reg,
// 			JumpTo: lab,
// 		}, e)
// 	// если истина - берем левое выражение
// 	bins = append(bins, addBinExpr(e.Lhs, reg, lid, false)...)
// 	// прыгаем в конец
// 	*lid++
// 	lend := *lid
// 	bins = appendBin(bins,
// 		&BinJMP{
// 			JumpTo: lend,
// 		}, e)

// 	// правое выражение
// 	bins = appendBin(bins,
// 		&BinLABEL{
// 			Label: lab,
// 		}, e)
// 	bins = append(bins, addBinExpr(e.Rhs, reg, lid, false)...)
// 	bins = appendBin(bins,
// 		&BinLABEL{
// 			Label: lend,
// 		}, e)

// case *ast.CallExpr:
// 	// если это анонимный вызов, то в reg сама функция, значит, параметры записываем в reg+1, иначе в reg
// 	var regoff int
// 	if e.Name == 0 {
// 		regoff = 1
// 	}

// 	// помещаем аргументы в массив аргументов в reg
// 	bins = appendBin(bins,
// 		&BinMAKESLICE{
// 			Reg: reg + regoff,
// 			Len: len(e.SubExprs),
// 			Cap: len(e.SubExprs),
// 		}, e)

// 	for i, ee := range e.SubExprs {
// 		// каждое выражение сохраняем в следующем по номеру регистре (относительно регистра слайса)
// 		bins = append(bins, addBinExpr(ee, reg+1+regoff, lid, false)...)
// 		bins = appendBin(bins,
// 			&BinSETIDX{
// 				Reg:    reg + regoff,
// 				Index:  i,
// 				RegVal: reg + 1 + regoff,
// 			}, ee)
// 	}

// 	bins = appendBin(bins,
// 		&BinCALL{
// 			Name:    e.Name,
// 			NumArgs: len(e.SubExprs),
// 			RegArgs: reg, // для анонимных (Name==0) - тут будет функция, иначе первый аргумент (см. выше)
// 			VarArg:  e.VarArg,
// 			Go:      e.Go,
// 			RegRets: reg,
// 		}, e)

// case *ast.AnonCallExpr:
// 	// помещаем в регистр значение функции (тип func, или ссылку на него, или интерфейс с ним)
// 	bins = append(bins, addBinExpr(e.Expr, reg, lid, false)...)
// 	// далее аргументы, как при вызове обычной функции
// 	bins = append(bins, addBinExpr(&ast.CallExpr{
// 		Name:     0,
// 		SubExprs: e.SubExprs,
// 		VarArg:   e.VarArg,
// 		Go:       e.Go,
// 	}, reg, lid, false)...) // передаем именно reg, т.к. он для Name==0 означает функцию, которую надо вызвать в BinCALL

// case *ast.MemberExpr:
// 	// здесь идет только вычисление значения свойства
// 	bins = append(bins, addBinExpr(e.Expr, reg, lid, false)...)
// 	bins = appendBin(bins,
// 		&BinGETMEMBER{
// 			Name: e.Name,
// 			Reg:  reg,
// 		}, e)
// case *ast.ItemExpr:
// 	// только вычисление значения по индексу
// 	bins = append(bins, addBinExpr(e.Value, reg, lid, false)...)
// 	bins = append(bins, addBinExpr(e.Index, reg+1, lid, false)...)
// 	bins = appendBin(bins,
// 		&BinGETIDX{
// 			Reg:      reg,
// 			RegIndex: reg + 1,
// 		}, e)
// case *ast.SliceExpr:
// 	// только вычисление субслайса
// 	bins = append(bins, addBinExpr(e.Value, reg, lid, false)...)
// 	bins = append(bins, addBinExpr(e.Begin, reg+1, lid, false)...)
// 	bins = append(bins, addBinExpr(e.End, reg+2, lid, false)...)
// 	bins = appendBin(bins,
// 		&BinGETSUBSLICE{
// 			Reg:      reg,
// 			RegBegin: reg + 1,
// 			RegEnd:   reg + 2,
// 		}, e)
// case *ast.FuncExpr:
// 	// *lid++
// 	// lend := *lid
// 	bins = appendBin(bins,
// 		&BinFUNC{
// 			Reg:    reg,
// 			Name:   e.Name,
// 			Code:   BinaryCode(e.Stmts, 0, lid),
// 			Args:   e.Args,
// 			VarArg: e.VarArg,
// 			// ReturnTo: lend,
// 		}, e)
// 	// КонецФункции
// 	// bins = appendBin(bins,
// 	// 	&BinLABEL{
// 	// 		Label: lend,
// 	// 	}, e)
// 	// // возвращаем значения в регистре reg, установленные функцией
// 	// bins = appendBin(bins,
// 	// 	&BinRET{
// 	// 		Reg: reg,
// 	// 	}, e)

// case *ast.TypeCast:
// 	bins = append(bins, addBinExpr(e.CastExpr, reg, lid, false)...)
// 	if e.TypeExpr == nil {
// 		bins = appendBin(bins,
// 			&BinLOAD{
// 				Reg:  reg + 1,
// 				Val:  e.Type,
// 				IsId: true,
// 			}, e)
// 	} else {
// 		bins = append(bins, addBinExpr(e.TypeExpr, reg+1, lid, false)...)
// 		bins = appendBin(bins,
// 			&BinSETNAME{
// 				Reg: reg + 1,
// 			}, e)
// 	}
// 	bins = appendBin(bins,
// 		&BinCASTTYPE{
// 			Reg:     reg,
// 			TypeReg: reg + 1,
// 		}, e)
// case *ast.MakeExpr:
// 	if e.TypeExpr == nil {
// 		bins = appendBin(bins,
// 			&BinLOAD{
// 				Reg:  reg,
// 				Val:  e.Type,
// 				IsId: true,
// 			}, e)
// 	} else {
// 		bins = append(bins, addBinExpr(e.TypeExpr, reg, lid, false)...)
// 		bins = appendBin(bins,
// 			&BinSETNAME{
// 				Reg: reg,
// 			}, e)
// 	}
// 	bins = appendBin(bins,
// 		&BinMAKE{
// 			Reg: reg,
// 		}, e)
// case *ast.MakeChanExpr:
// 	if e.SizeExpr == nil {
// 		bins = appendBin(bins,
// 			&BinLOAD{
// 				Reg: reg,
// 				Val: int64(0),
// 			}, e)
// 	} else {
// 		bins = append(bins, addBinExpr(e.SizeExpr, reg, lid, false)...)
// 	}
// 	bins = appendBin(bins,
// 		&BinMAKECHAN{
// 			Reg: reg,
// 		}, e)
// case *ast.MakeArrayExpr:
// 	bins = append(bins, addBinExpr(e.LenExpr, reg, lid, false)...)
// 	if e.CapExpr == nil {
// 		bins = appendBin(bins,
// 			&BinMV{
// 				RegFrom: reg,
// 				RegTo:   reg + 1,
// 			}, e)
// 	} else {
// 		bins = append(bins, addBinExpr(e.CapExpr, reg+1, lid, false)...)
// 	}
// 	bins = appendBin(bins,
// 		&BinMAKEARR{
// 			Reg:    reg,
// 			RegCap: reg + 1,
// 		}, e)
// case *ast.ChanExpr:

// 	// определяем значение справа
// 	bins = append(bins, addBinExpr(e.Rhs, reg+1, lid, false)...)
// 	if e.Lhs == nil {
// 		// слева нет значения - это временное чтение из канала без сохранения значения в переменной
// 		bins = appendBin(bins,
// 			&BinCHANRECV{
// 				Reg:    reg + 1,
// 				RegVal: reg,
// 			}, e)
// 	} else {
// 		// значение слева
// 		bins = append(bins, addBinExpr(e.Lhs, reg+2, lid, false)...)
// 		bins = appendBin(bins,
// 			&BinMV{
// 				RegFrom: reg + 2,
// 				RegTo:   reg + 3,
// 			}, e)
// 		// слева канал - пишем в него правое
// 		bins = appendBin(bins,
// 			&BinISKIND{
// 				Reg:  reg + 3,
// 				Kind: reflect.Chan,
// 			}, e)
// 		*lid++
// 		li := *lid
// 		bins = appendBin(bins,
// 			&BinJFALSE{
// 				Reg:    reg + 3,
// 				JumpTo: li,
// 			}, e)

// 		bins = appendBin(bins,
// 			&BinCHANSEND{
// 				Reg:    reg + 2,
// 				RegVal: reg + 1,
// 			}, e)

// 		bins = appendBin(bins,
// 			&BinLOAD{
// 				Reg: reg,
// 				Val: true,
// 			}, e)

// 		*lid++
// 		li2 := *lid

// 		bins = appendBin(bins,
// 			&BinJMP{
// 				JumpTo: li2,
// 			}, e)

// 		// иначе справа канал, а слева переменная (установим, если прочитали из канала)
// 		bins = appendBin(bins,
// 			&BinLABEL{
// 				Label: li,
// 			}, e)
// 		bins = appendBin(bins,
// 			&BinCHANRECV{
// 				Reg:    reg + 1,
// 				RegVal: reg,
// 			}, e)

// 		bins = append(bins, addBinLetExpr(e.Lhs, reg, lid)...)

// 		bins = appendBin(bins,
// 			&BinLABEL{
// 				Label: li2,
// 			}, e)
// 	}

// case *ast.AssocExpr:
// 	switch e.Operator {
// 	case "++":
// 		if alhs, ok := e.Lhs.(*ast.IdentExpr); ok {
// 			bins = appendBin(bins,
// 				&BinGET{
// 					Reg: reg,
// 					Id:  alhs.Id,
// 				}, alhs)
// 			bins = appendBin(bins,
// 				&BinINC{
// 					Reg: reg,
// 				}, alhs)
// 			bins = appendBin(bins,
// 				&BinSET{
// 					Reg: reg,
// 					Id:  alhs.Id,
// 				}, alhs)
// 		}
// 	case "--":
// 		if alhs, ok := e.Lhs.(*ast.IdentExpr); ok {
// 			bins = appendBin(bins,
// 				&BinGET{
// 					Reg: reg,
// 					Id:  alhs.Id,
// 				}, alhs)
// 			bins = appendBin(bins,
// 				&BinDEC{
// 					Reg: reg,
// 				}, alhs)
// 			bins = appendBin(bins,
// 				&BinSET{
// 					Reg: reg,
// 					Id:  alhs.Id,
// 				}, alhs)
// 		}
// 	default:
// 		bins = append(bins, addBinExpr(&ast.BinOpExpr{Lhss: []ast.Expr{e.Lhs}, Operator: e.Operator[0:1], Rhss: []ast.Expr{e.Rhs}}, reg, lid, false)...)
// 		bins = append(bins, addBinLetExpr(e.Lhs, reg, lid)...)
// 	}

// case *ast.LetExpr:
// 	// пока не используется (не распознается парсером), планируется добавить предопределенные значения для функций

// 	}

// 	return
// }
