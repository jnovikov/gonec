package bincode

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/covrom/gonec/ast"
	"github.com/covrom/gonec/builtins"
)

type BinStmt interface {
	ast.Pos
	binstmt()
	SwapId(map[int]int)
}

type BinStmtImpl struct {
	ast.PosImpl
	fmt.Stringer
}

func (x *BinStmtImpl) binstmt()           {}
func (x *BinStmtImpl) SwapId(map[int]int) {}

type BinCode []BinStmt

func (v BinCode) String() string {
	s := ""
	for _, e := range v {
		s += fmt.Sprintf("%v\n", e)
	}
	return s
}

func WriteBinCode(w io.Writer, v BinCode) error {
	zw := gzip.NewWriter(w)
	zw.Name = "Gonec binary code"
	zw.Comment = "Created with https://covrom.github.io/gonec/ by Roman TSovanyan rs@tsov.pro"
	zw.ModTime = time.Now()

	enc := gob.NewEncoder(zw)

	// так же сохраняем уникальные имена
	if err := enc.Encode(*ast.UniqueNames); err != nil {
		return err
	}

	if err := enc.Encode(v); err != nil {
		return err
	}

	if err := zw.Close(); err != nil {
		return err
	}
	return nil
}

func ReadBinCode(r io.Reader) (res BinCode, err error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	dec := gob.NewDecoder(zr)

	var gnxNames = ast.NewEnvNames()

	if err := dec.Decode(gnxNames); err != nil {
		return nil, err
	}
	if err := dec.Decode(&res); err != nil {
		return nil, err
	}
	if err := zr.Close(); err != nil {
		return nil, err
	}

	// переносим загруженные имена в текущий контекст
	// и заменяем идентификаторы в загружаемом коде в случае конфликта
	swapIdents := make(map[int]int)

	// log.Println(gnxNames)

	for i, v := range gnxNames.Handlow {

		// log.Printf("Проверяем %d, %q", i, v)

		if vv, ok := ast.UniqueNames.GetLowerCaseOk(i); ok {
			// под тем же идентификатором находится другая строка, без учета регистра
			if v != vv {
				// новый id
				ii := ast.UniqueNames.Set(gnxNames.Handles[i])
				swapIdents[i] = ii

				// log.Printf("Заменяем %d на %d для загружаемого %q, уже есть %q\n", i, ii, v, vv)

			}
		} else {
			// такого идентификатора еще нет - устанавливаем значение на него
			// последующие идентификаторы ast.UniqueNames будут идти после него

			// log.Printf("Устанавливаем %d для загружаемого %q\n", i, gnxNames.Handles[i])

			ast.UniqueNames.SetToId(gnxNames.Handles[i], i)
		}
	}

	// заменяем идентификаторы, если при слиянии были конфликты
	for _, v := range res {
		v.SwapId(swapIdents)
	}

	return res, nil
}

func init() {
	gob.Register(&ast.EnvNames{})
	gob.Register(core.VMTime{})
	gob.Register(core.VMSlice{})
	gob.Register(core.VMStringMap{})

	gob.Register(&BinLOAD{})
	gob.Register(&BinMV{})
	gob.Register(&BinEQUAL{})
	gob.Register(&BinCASTNUM{})
	gob.Register(&BinMAKESLICE{})
	gob.Register(&BinSETIDX{})
	gob.Register(&BinMAKEMAP{})
	gob.Register(&BinSETKEY{})
	gob.Register(&BinGET{})
	gob.Register(&BinSET{})
	gob.Register(&BinSETMEMBER{})
	gob.Register(&BinSETNAME{})
	gob.Register(&BinSETITEM{})
	gob.Register(&BinSETSLICE{})
	gob.Register(&BinUNARY{})
	gob.Register(&BinADDRID{})
	gob.Register(&BinADDRMBR{})
	gob.Register(&BinUNREFID{})
	gob.Register(&BinUNREFMBR{})
	gob.Register(&BinLABEL{})
	gob.Register(&BinJMP{})
	gob.Register(&BinJTRUE{})
	gob.Register(&BinJFALSE{})
	gob.Register(&BinOPER{})
	gob.Register(&BinCALL{})
	gob.Register(&BinGETMEMBER{})
	gob.Register(&BinGETIDX{})
	gob.Register(&BinGETSUBSLICE{})
	gob.Register(&BinFUNC{})
	gob.Register(&BinCASTTYPE{})
	gob.Register(&BinMAKE{})
	gob.Register(&BinMAKECHAN{})
	gob.Register(&BinMAKEARR{})
	gob.Register(&BinCHANRECV{})
	gob.Register(&BinCHANSEND{})
	gob.Register(&BinISKIND{})
	gob.Register(&BinISSLICE{})
	gob.Register(&BinTRY{})
	gob.Register(&BinCATCH{})
	gob.Register(&BinPOPTRY{})
	gob.Register(&BinFOREACH{})
	gob.Register(&BinNEXT{})
	gob.Register(&BinPOPFOR{})
	gob.Register(&BinFORNUM{})
	gob.Register(&BinNEXTNUM{})
	gob.Register(&BinWHILE{})
	gob.Register(&BinBREAK{})
	gob.Register(&BinCONTINUE{})
	gob.Register(&BinRET{})
	gob.Register(&BinTHROW{})
	gob.Register(&BinMODULE{})
	gob.Register(&BinERROR{})
	gob.Register(&BinTRYRECV{})
	gob.Register(&BinTRYSEND{})
	gob.Register(&BinGOSHED{})
	gob.Register(&BinINC{})
	gob.Register(&BinDEC{})
	gob.Register(&BinFREE{})

}

//////////////////////
// команды байткода
//////////////////////

type BinLOAD struct {
	BinStmtImpl

	Reg  int
	Val  interface{}
	IsId bool
}

func (v *BinLOAD) SwapId(m map[int]int) {
	if v.IsId {
		if newid, ok := m[v.Val.(int)]; ok {
			v.Val = newid
		}
	}
}
func (v BinLOAD) String() string {
	if v.IsId {
		return fmt.Sprintf("LOAD r%d, %#v", v.Reg, ast.UniqueNames.Get(v.Val.(int)))
	}
	return fmt.Sprintf("LOAD r%d, %#v", v.Reg, v.Val)
}

type BinMV struct {
	BinStmtImpl

	RegFrom int
	RegTo   int
}

func (v BinMV) String() string {
	return fmt.Sprintf("MV r%d, r%d", v.RegTo, v.RegFrom)
}

type BinEQUAL struct {
	BinStmtImpl

	Reg  int // результат сравнения помещаем сюда, тип "булево"
	Reg1 int
	Reg2 int
}

func (v BinEQUAL) String() string {
	return fmt.Sprintf("EQUAL r%d, r%d == r%d", v.Reg, v.Reg1, v.Reg2)
}

type BinCASTNUM struct {
	BinStmtImpl

	Reg int
}

func (v BinCASTNUM) String() string {
	return fmt.Sprintf("CAST r%d, NUMBER", v.Reg)
}

type BinMAKESLICE struct {
	BinStmtImpl

	Reg int
	Len int
	Cap int
}

func (v BinMAKESLICE) String() string {
	return fmt.Sprintf("MAKESLICE r%d, LEN %d, CAP %d", v.Reg, v.Len, v.Cap)
}

type BinSETIDX struct {
	BinStmtImpl

	Reg    int
	Index  int
	RegVal int
}

func (v BinSETIDX) String() string {
	return fmt.Sprintf("SETIDX r%d[%d], r%d", v.Reg, v.Index, v.RegVal)
}

type BinMAKEMAP struct {
	BinStmtImpl

	Reg int
	Len int
}

func (v BinMAKEMAP) String() string {
	return fmt.Sprintf("MAKEMAP r%d, LEN %d", v.Reg, v.Len)
}

type BinSETKEY struct {
	BinStmtImpl

	Reg    int
	Key    string
	RegVal int
}

func (v BinSETKEY) String() string {
	return fmt.Sprintf("SETKEY r%d[%q], r%d", v.Reg, v.Key, v.RegVal)
}

type BinGET struct {
	BinStmtImpl

	Reg int
	Id  int
}

func (v *BinGET) SwapId(m map[int]int) {
	if newid, ok := m[v.Id]; ok {
		v.Id = newid
		// log.Printf("Замена в %#v %v\n",v, v)
	}
}
func (v BinGET) String() string {
	return fmt.Sprintf("GET r%d, %q", v.Reg, ast.UniqueNames.Get(v.Id))
}

type BinSET struct {
	BinStmtImpl

	Id  int // id переменной
	Reg int // регистр со значением
}

func (v *BinSET) SwapId(m map[int]int) {
	if newid, ok := m[v.Id]; ok {
		v.Id = newid
		// log.Printf("Замена в %#v %v\n",v, v)
	}
}
func (v BinSET) String() string {
	return fmt.Sprintf("SET %q, r%d", ast.UniqueNames.Get(v.Id), v.Reg)
}

type BinSETMEMBER struct {
	BinStmtImpl

	Reg    int // регистр со структтурой или мапой
	Id     int // id поля структуры или мапы
	RegVal int // регистр со значением
}

func (v *BinSETMEMBER) SwapId(m map[int]int) {
	if newid, ok := m[v.Id]; ok {
		v.Id = newid
		// log.Printf("Замена в %#v %v\n",v, v)
	}
}
func (v BinSETMEMBER) String() string {
	return fmt.Sprintf("SETMEMBER r%d.%q, r%d", v.Reg, ast.UniqueNames.Get(v.Id), v.RegVal)
}

type BinSETNAME struct {
	BinStmtImpl

	Reg int // регистр с именем (строкой), сюда же возвращается id имени, записанного в ast.UniqueNames.Set()
}

func (v BinSETNAME) String() string {
	return fmt.Sprintf("SETNAME r%d", v.Reg)
}

type BinSETITEM struct {
	BinStmtImpl

	Reg        int
	RegIndex   int
	RegVal     int
	RegNeedLet int
}

func (v BinSETITEM) String() string {
	return fmt.Sprintf("SETITEM r%d[r%d], r%d", v.Reg, v.RegIndex, v.RegVal)
}

type BinSETSLICE struct {
	BinStmtImpl

	Reg        int
	RegBegin   int
	RegEnd     int
	RegVal     int
	RegNeedLet int
}

func (v BinSETSLICE) String() string {
	return fmt.Sprintf("SETSLICE r%d[r%d:r%d], r%d", v.Reg, v.RegBegin, v.RegEnd, v.RegVal)
}

type BinUNARY struct {
	BinStmtImpl

	Reg int
	Op  rune // - ! ^
}

func (v BinUNARY) String() string {
	return fmt.Sprintf("UNARY %sr%d", string(v.Op), v.Reg)
}

type BinADDRID struct {
	BinStmtImpl

	Reg  int
	Name int
}

func (v *BinADDRID) SwapId(m map[int]int) {
	if newid, ok := m[v.Name]; ok {
		v.Name = newid
	}
}
func (v BinADDRID) String() string {
	return fmt.Sprintf("ADDRID r%d, %q", v.Reg, ast.UniqueNames.Get(v.Name))
}

type BinADDRMBR struct {
	BinStmtImpl

	Reg  int
	Name int
}

func (v *BinADDRMBR) SwapId(m map[int]int) {
	if newid, ok := m[v.Name]; ok {
		v.Name = newid
	}
}
func (v BinADDRMBR) String() string {
	return fmt.Sprintf("ADDRMBR r%d, r%d.%q", v.Reg, v.Reg, ast.UniqueNames.Get(v.Name))
}

type BinUNREFID struct {
	BinStmtImpl

	Reg  int
	Name int
}

func (v *BinUNREFID) SwapId(m map[int]int) {
	if newid, ok := m[v.Name]; ok {
		v.Name = newid
	}
}
func (v BinUNREFID) String() string {
	return fmt.Sprintf("UNREFID r%d, %q", v.Reg, ast.UniqueNames.Get(v.Name))
}

type BinUNREFMBR struct {
	BinStmtImpl

	Reg  int
	Name int
}

func (v *BinUNREFMBR) SwapId(m map[int]int) {
	if newid, ok := m[v.Name]; ok {
		v.Name = newid
	}
}
func (v BinUNREFMBR) String() string {
	return fmt.Sprintf("UNREFMBR r%d, r%d.%q", v.Reg, v.Reg, ast.UniqueNames.Get(v.Name))
}

type BinLABEL struct {
	BinStmtImpl

	Label int
}

func (v BinLABEL) String() string {
	return fmt.Sprintf("L%d:", v.Label)
}

type BinJMP struct {
	BinStmtImpl

	JumpTo int
}

func (v BinJMP) String() string {
	return fmt.Sprintf("JMP L%d", v.JumpTo)
}

type BinJTRUE struct {
	BinStmtImpl

	Reg    int
	JumpTo int
}

func (v BinJTRUE) String() string {
	return fmt.Sprintf("JTRUE r%d, L%d", v.Reg, v.JumpTo)
}

type BinJFALSE struct {
	BinStmtImpl

	Reg    int
	JumpTo int
}

func (v BinJFALSE) String() string {
	return fmt.Sprintf("JFALSE r%d, L%d", v.Reg, v.JumpTo)
}

type BinOPER struct {
	BinStmtImpl

	RegL int // сюда же помещается результат
	RegR int
	Op   int
}

func (v BinOPER) String() string {
	return fmt.Sprintf("OP r%d, %q, r%d", v.RegL, OperMapR[v.Op], v.RegR)
}

type BinCALL struct {
	BinStmtImpl

	Name int // либо вызов по имени из ast.UniqueNames, если Name != 0
	// либо вызов обработчика (Name==0), напр. для анонимной функции
	// (выражение типа func, или ссылка или интерфейс с ним, находится в reg, а параметры начиная с reg+1)
	NumArgs int // число аргументов, которое надо взять на входе из массива (Reg)
	RegArgs int // регистр с массивом аругментов
	RegRets int // массив с возвращаемыми из функции значениями

	// в последнем регистре (в RegArgs) может быть передан
	// массив аргументов переменной длины, и это приемлемо для вызываемой функции (оператор "...")
	// таким массивом будет только последний аргумент
	VarArg bool

	Go bool // признак необходимости запуска в новой горутине
}

func (v *BinCALL) SwapId(m map[int]int) {
	if v.Name == 0 {
		return
	}
	if newid, ok := m[v.Name]; ok {
		v.Name = newid
		// log.Printf("Замена в %#v %v\n",v, v)
	}
}
func (v BinCALL) String() string {
	if v.Name == 0 {
		return fmt.Sprintf("CALL REG r%d, ARGS r%d, ARGS_COUNT %d, VARARG %v, GO %v, RETURN r%d", v.RegArgs, v.RegArgs+1, v.NumArgs, v.VarArg, v.Go, v.RegRets)
	}
	return fmt.Sprintf("CALL %q, ARGS r%d, ARGS_COUNT %d, VARARG %v, GO %v, RETURN r%d", ast.UniqueNames.Get(v.Name), v.RegArgs, v.NumArgs, v.VarArg, v.Go, v.RegRets)
}

type BinGETMEMBER struct {
	BinStmtImpl

	Reg  int
	Name int
}

func (v *BinGETMEMBER) SwapId(m map[int]int) {
	if newid, ok := m[v.Name]; ok {
		v.Name = newid
		// log.Printf("Замена в %#v %v\n",v, v)
	}
}
func (v BinGETMEMBER) String() string {
	return fmt.Sprintf("GETMEMBER r%d, %q", v.Reg, ast.UniqueNames.Get(v.Name))
}

type BinGETIDX struct {
	BinStmtImpl

	Reg      int
	RegIndex int
}

func (v BinGETIDX) String() string {
	return fmt.Sprintf("GETIDX r%d[r%d]", v.Reg, v.RegIndex)
}

type BinGETSUBSLICE struct {
	BinStmtImpl

	Reg      int
	RegBegin int
	RegEnd   int
}

func (v BinGETSUBSLICE) String() string {
	return fmt.Sprintf("SLICE r%d[r%d : r%d]", v.Reg, v.RegBegin, v.RegEnd)
}

type BinFUNC struct {
	BinStmtImpl

	Reg    int // регистр, в который сохраняется значение определяемой функции типа func
	Name   int
	Code   BinCode
	Args   []int // идентификаторы параметров
	VarArg bool
	// ReturnTo int //метка инструкции возврата из функции
}

func (v *BinFUNC) SwapId(m map[int]int) {
	if newid, ok := m[v.Name]; ok && v.Name != 0 {
		v.Name = newid
		// log.Printf("Замена в %#v %v\n",v, v)
	}
	for i := range v.Args {
		if newid, ok := m[v.Args[i]]; ok && v.Args[i] != 0 {
			v.Args[i] = newid
			// log.Printf("Замена в аргументах %#v %v\n",v, v)
		}
	}
}
func (v BinFUNC) String() string {
	s := ""
	for _, a := range v.Args {
		if s != "" {
			s += ", "
		}
		s += ast.UniqueNames.Get(a)
	}
	vrg := ""
	if v.VarArg {
		vrg = "..."
	}
	return fmt.Sprintf("FUNC r%d, %q, (%s %s)\n{\n%v}\n", v.Reg, ast.UniqueNames.Get(v.Name), s, vrg, v.Code)
}

type BinCASTTYPE struct {
	BinStmtImpl

	Reg     int
	TypeReg int
}

func (v BinCASTTYPE) String() string {
	return fmt.Sprintf("CAST r%d AS TYPE r%d", v.Reg, v.TypeReg)
}

type BinMAKE struct {
	BinStmtImpl

	Reg int // здесь id типа, и сюда же пишем новое значение
}

func (v BinMAKE) String() string {
	return fmt.Sprintf("MAKE r%d AS TYPE r%d", v.Reg, v.Reg)
}

type BinMAKECHAN struct {
	BinStmtImpl

	Reg int // тут размер буфера (0=без буфера), сюда же помещается созданный канал
}

func (v BinMAKECHAN) String() string {
	return fmt.Sprintf("MAKECHAN r%d SIZE r%d", v.Reg, v.Reg)
}

type BinMAKEARR struct {
	BinStmtImpl

	Reg    int // тут длина, сюда же помещается слайс
	RegCap int
}

func (v BinMAKEARR) String() string {
	return fmt.Sprintf("MAKEARR r%d, LEN r%d, CAP r%d", v.Reg, v.Reg, v.RegCap)
}

type BinCHANRECV struct {
	BinStmtImpl
	// с ожиданием
	Reg    int // канал
	RegVal int // сюда помещается результат
}

func (v BinCHANRECV) String() string {
	return fmt.Sprintf("<-CHAN r%d, r%d", v.RegVal, v.Reg)
}

type BinCHANSEND struct {
	BinStmtImpl
	// с ожиданием
	Reg    int // канал
	RegVal int // значение
}

func (v BinCHANSEND) String() string {
	return fmt.Sprintf("CHAN<- r%d, r%d", v.Reg, v.RegVal)
}

type BinISKIND struct {
	BinStmtImpl

	Reg  int          // значение для проверки, сюда же возвращается bool
	Kind reflect.Kind // категория типа значения в reg
}

func (v BinISKIND) String() string {
	return fmt.Sprintf("ISKIND r%d, %s", v.Reg, v.Kind)
}

type BinISSLICE struct {
	BinStmtImpl

	Reg     int // значение для проверки
	RegBool int //сюда возвращается bool
}

func (v BinISSLICE) String() string {
	return fmt.Sprintf("ISSLICE r%d, r%d", v.RegBool, v.Reg)
}

type BinTRY struct {
	BinStmtImpl

	Reg    int // регистр, куда будет помещаться error во время выполнения последующего кода
	JumpTo int // метка блока обработки ошибки
}

func (v BinTRY) String() string {
	return fmt.Sprintf("TRY r%d, CATCH L%d", v.Reg, v.JumpTo)
}

type BinCATCH struct {
	BinStmtImpl

	Reg    int
	JumpTo int
}

func (v BinCATCH) String() string {
	return fmt.Sprintf("CATCH r%d, NOERR L%d", v.Reg, v.JumpTo)
}

type BinPOPTRY struct {
	BinStmtImpl

	CatchLabel int // снимаем со стека исключений конструкцию с этим регистром
}

func (v BinPOPTRY) String() string {
	return fmt.Sprintf("POPTRY L%d", v.CatchLabel)
}

type BinFOREACH struct {
	BinStmtImpl

	Reg           int // регистр для итерационного выбора из него значений
	RegIter       int // в этот регистр будет записываться итератор
	BreakLabel    int
	ContinueLabel int
}

func (v BinFOREACH) String() string {
	return fmt.Sprintf("FOREACH r%d, ITER r%d, BREAK TO L%d", v.Reg, v.RegIter, v.BreakLabel)
}

type BinNEXT struct {
	BinStmtImpl

	Reg int // выбираем из этого регистра следующее значение и помещаем в регистр RegVal
	// это может быть очередное значение из слайса или из канала, зависит от типа значения в Reg
	RegVal  int
	RegIter int // регистр с итератором, инициализированным FOREACH
	JumpTo  int // переход в случае, если нет очередного значения (достигнут конец выборки)
	// туда же переходим по Прервать
}

func (v BinNEXT) String() string {
	return fmt.Sprintf("NEXT r%d, FROM r%d, ITER r%d, ENDLOOP L%d", v.RegVal, v.Reg, v.RegIter, v.JumpTo)
}

type BinPOPFOR struct {
	BinStmtImpl

	ContinueLabel int // снимаем со стека циклов конструкцию с этой меткой
}

func (v BinPOPFOR) String() string {
	return fmt.Sprintf("POPFOR L%d", v.ContinueLabel)
}

type BinFORNUM struct {
	BinStmtImpl

	Reg           int // регистр для итерационного значения
	RegFrom       int // регистр с начальным значением
	RegTo         int // регистр с конечным значением
	BreakLabel    int
	ContinueLabel int
}

func (v BinFORNUM) String() string {
	return fmt.Sprintf("FORNUM r%d, FROM r%d, TO r%d, BREAK TO L%d", v.Reg, v.RegFrom, v.RegTo, v.BreakLabel)
}

type BinNEXTNUM struct {
	BinStmtImpl

	Reg     int // следующее значение итератора
	RegFrom int // регистр с начальным значением
	RegTo   int // регистр с конечным значением
	JumpTo  int // переход в случае, если значение после увеличения стало больше, чем ранее определенное в RegTo
	// туда же переходим по Прервать
}

func (v BinNEXTNUM) String() string {
	return fmt.Sprintf("NEXTNUM r%d, ENDLOOP L%d", v.Reg, v.JumpTo)
}

type BinWHILE struct {
	BinStmtImpl

	BreakLabel    int
	ContinueLabel int
}

func (v BinWHILE) String() string {
	return fmt.Sprintf("WHILE BREAK TO L%d", v.BreakLabel)
}

type BinBREAK struct {
	BinStmtImpl
}

func (v BinBREAK) String() string {
	return fmt.Sprintf("BREAK")
}

type BinCONTINUE struct {
	BinStmtImpl
}

func (v BinCONTINUE) String() string {
	return fmt.Sprintf("CONTINUE")
}

type BinRET struct {
	BinStmtImpl

	Reg int
}

func (v BinRET) String() string {
	return fmt.Sprintf("RETURN r%d", v.Reg)
}

type BinTHROW struct {
	BinStmtImpl

	Reg int
}

func (v BinTHROW) String() string {
	return fmt.Sprintf("THROW r%d", v.Reg)
}

type BinMODULE struct {
	BinStmtImpl

	Name int
	Code BinCode
}

func (v *BinMODULE) SwapId(m map[int]int) {
	if newid, ok := m[v.Name]; ok {
		v.Name = newid
		// log.Printf("Замена в %#v %v\n",v, v)
	}
}
func (v BinMODULE) String() string {
	return fmt.Sprintf("MODULE %s\n{\n%v}\n", ast.UniqueNames.Get(v.Name), v.Code)
}

type BinERROR struct {
	BinStmtImpl

	Error string
}

func (v BinERROR) String() string {
	return fmt.Sprintf("ERROR %q", v.Error)
}

type BinTRYRECV struct {
	BinStmtImpl

	Reg       int // на входе канал, на выходе тоже
	RegVal    int // получаемое значение
	RegOk     int // успешное чтение, или не было чтения, или в Reg не канал
	RegClosed int // в этот регистр помещается true если канал закрыт
}

func (v BinTRYRECV) String() string {
	return fmt.Sprintf("TRYRECV r%d, OK r%d", v.Reg, v.RegOk)
}

type BinTRYSEND struct {
	BinStmtImpl

	Reg    int // на входе канал, на выходе тоже
	RegVal int // регистр со значением для отправки
	RegOk  int // успешно передано в канал, или не было передачи, или в Reg не канал
	// RegClosed int // в этот регистр помещается true если канал закрыт
}

func (v BinTRYSEND) String() string {
	return fmt.Sprintf("TRYSEND r%d, r%d, OK r%d", v.Reg, v.RegVal, v.RegOk)
}

type BinGOSHED struct {
	BinStmtImpl
}

func (v BinGOSHED) String() string {
	return fmt.Sprintf("GOSHED")
}

type BinINC struct {
	BinStmtImpl

	Reg int
}

func (v BinINC) String() string {
	return fmt.Sprintf("INC r%d", v.Reg)
}

type BinDEC struct {
	BinStmtImpl

	Reg int
}

func (v BinDEC) String() string {
	return fmt.Sprintf("DEC r%d", v.Reg)
}

type BinFREE struct {
	BinStmtImpl

	Reg int
}

func (v BinFREE) String() string {
	return fmt.Sprintf("FREE FROM r%d", v.Reg)
}
