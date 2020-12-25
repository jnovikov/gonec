// Package core implements core interface for gonec script.
package core

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/covrom/decnum"

	"github.com/covrom/gonec/names"
	"github.com/satori/go.uuid"
)

// LoadAllBuiltins is a convenience function that loads all defineSd builtins.
func LoadAllBuiltins(env *Env) {
	Import(env)

	pkgs := map[string]func(env *Env) *Env{
		// "sort":          gonec_sort.Import,
		// "strings":       gonec_strings.Import,
	}

	env.DefineS("импорт", VMFunc(func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		if len(args) != 1 {
			return VMErrorNeedSinglePacketName
		}
		if s, ok := args[0].(VMString); ok {
			if loader, ok := pkgs[strings.ToLower(string(s))]; ok {
				rets.Append(loader(env)) // возвращает окружение, инициализированное пакетом
				return nil
			}
			return fmt.Errorf("Пакет '%s' не найден", s)
		} else {
			return VMErrorNeedString
		}
	}))

	// успешно загружен глобальный контекст
	env.SetBuiltsIsLoaded()
}

// Import общая стандартная бибилиотека
func Import(env *Env) *Env {

	env.DefineS("длина", VMFuncMustParams(1, func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		if rv, ok := args[0].(VMIndexer); ok {
			rets.Append(rv.Length())
			return nil
		}
		return VMErrorNeedLength
	}))

	env.DefineS("диапазон", VMFunc(func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		if len(args) < 1 {
			return VMErrorNoArgs
		}
		if len(args) > 2 {
			return VMErrorNeedLengthOrBoundary
		}
		var min, max int64
		var arr VMSlice
		if len(args) == 1 {
			min = 0
			maxvm, ok := args[0].(VMInt)
			if !ok {
				return VMErrorNeedInt
			}
			max = maxvm.Int() - 1
		} else {
			minvm, ok := args[0].(VMInt)
			if !ok {
				return VMErrorNeedInt
			}
			min = minvm.Int()
			maxvm, ok := args[1].(VMInt)
			if !ok {
				return VMErrorNeedInt
			}
			max = maxvm.Int()
		}
		if min > max {
			return VMErrorNeedLess
		}
		arr = make(VMSlice, max-min+1)

		for i := min; i <= max; i++ {
			arr[i-min] = VMInt(i)
		}
		rets.Append(arr)
		return nil
	}))

	env.DefineS("текущаядата", VMFuncMustParams(0, func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		rets.Append(Now())
		return nil
	}))

	env.DefineS("прошловременис", VMFuncMustParams(1, func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		if rv, ok := args[0].(VMDateTimer); ok {
			rets.Append(Now().Sub(rv.Time()))
			return nil
		}
		return VMErrorNeedDate
	}))

	env.DefineS("пауза", VMFuncMustParams(1, func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		if v, ok := args[0].(VMNumberer); ok {
			sec1 := NewVMDecNumFromInt64(int64(VMSecond))
			time.Sleep(time.Duration(v.DecNum().Mul(sec1).Int()))
			return nil
		}
		return VMErrorNeedSeconds
	}))

	env.DefineS("длительностьнаносекунды", VMNanosecond)
	env.DefineS("длительностьмикросекунды", VMMicrosecond)
	env.DefineS("длительностьмиллисекунды", VMMillisecond)
	env.DefineS("длительностьсекунды", VMSecond)
	env.DefineS("длительностьминуты", VMMinute)
	env.DefineS("длительностьчаса", VMHour)
	env.DefineS("длительностьдня", VMDay)

	env.DefineS("прочитатьфайл", VMFuncMustParams(1, func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		if v, ok := args[0].(VMString); ok {
			data, err := ioutil.ReadFile(v.String())
			rets.Append(VMString(data), VMBool(err == nil))
			return nil
		}
		return VMErrorNeedString
	}))

	env.DefineS("хэш", VMFuncMustParams(1, func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		if v, ok := args[0].(VMHasher); ok {
			rets.Append(v.Hash())
			return nil
		}
		return VMErrorNeedHash
	}))

	env.DefineS("уникальныйидентификатор", VMFuncMustParams(0, func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		rets.Append(VMString(uuid.NewV1().String()))
		return nil
	}))

	env.DefineS("получитьмассивизпула", VMFuncMustParams(0, func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		rets.Append(GetGlobalVMSlice())
		return nil
	}))

	env.DefineS("вернутьмассиввпул", VMFuncMustParams(1, func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		if v, ok := args[0].(VMSlice); ok {
			PutGlobalVMSlice(v)
			return nil
		}
		return VMErrorNeedMap
	}))

	env.DefineS("случайнаястрока", VMFuncMustParams(1, func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		if v, ok := args[0].(VMInt); ok {
			rets.Append(VMString(MustGenerateRandomString(int(v))))
			return nil
		}
		return VMErrorNeedInt
	}))

	env.DefineS("нрег", VMFuncMustParams(1, func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		if v, ok := args[0].(VMStringer); ok {
			rets.Append(VMString(strings.ToLower(string(v.String()))))
			return nil
		}
		return VMErrorNeedString
	}))

	env.DefineS("врег", VMFuncMustParams(1, func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		if v, ok := args[0].(VMStringer); ok {
			rets.Append(VMString(strings.ToUpper(string(v.String()))))
			return nil
		}
		return VMErrorNeedString
	}))

	env.DefineS("стрсодержит", VMFuncMustParams(2, func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		v1, ok1 := args[0].(VMStringer)
		v2, ok2 := args[1].(VMStringer)
		if ok1 && ok2 {
			rets.Append(VMBool(strings.Contains(string(v1.String()), string(v2.String()))))
			return nil
		}
		return VMErrorNeedString
	}))

	env.DefineS("стрсодержитлюбой", VMFuncMustParams(2, func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		v1, ok1 := args[0].(VMStringer)
		v2, ok2 := args[1].(VMStringer)
		if ok1 && ok2 {
			rets.Append(VMBool(strings.ContainsAny(string(v1.String()), string(v2.String()))))
			return nil
		}
		return VMErrorNeedString
	}))

	env.DefineS("стрколичество", VMFuncMustParams(2, func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		v1, ok1 := args[0].(VMStringer)
		v2, ok2 := args[1].(VMStringer)
		if ok1 && ok2 {
			rets.Append(VMInt(strings.Count(string(v1.String()), string(v2.String()))))
			return nil
		}
		return VMErrorNeedString
	}))

	env.DefineS("стрнайти", VMFuncMustParams(2, func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		v1, ok1 := args[0].(VMStringer)
		v2, ok2 := args[1].(VMStringer)
		if ok1 && ok2 {
			rets.Append(VMInt(strings.Index(string(v1.String()), string(v2.String()))))
			return nil
		}
		return VMErrorNeedString
	}))

	env.DefineS("стрнайтилюбой", VMFuncMustParams(2, func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		v1, ok1 := args[0].(VMStringer)
		v2, ok2 := args[1].(VMStringer)
		if ok1 && ok2 {
			rets.Append(VMInt(strings.IndexAny(string(v1.String()), string(v2.String()))))
			return nil
		}
		return VMErrorNeedString
	}))

	env.DefineS("стрнайтипоследний", VMFuncMustParams(2, func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		v1, ok1 := args[0].(VMStringer)
		v2, ok2 := args[1].(VMStringer)
		if ok1 && ok2 {
			rets.Append(VMInt(strings.LastIndex(string(v1.String()), string(v2.String()))))
			return nil
		}
		return VMErrorNeedString
	}))

	env.DefineS("стрзаменить", VMFuncMustParams(3, func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		v1, ok1 := args[0].(VMStringer)
		v2, ok2 := args[1].(VMStringer)
		v3, ok3 := args[2].(VMStringer)
		if ok1 && ok2 && ok3 {
			rets.Append(VMString(strings.Replace(string(v1.String()), string(v2.String()), string(v3.String()), -1)))
			return nil
		}
		return VMErrorNeedString
	}))

	env.DefineS("стрдекодироватьзапрос", VMFuncMustParams(1, func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		if v, ok := args[0].(VMString); ok {
			dec, err := url.QueryUnescape(v.String())
			rets.Append(VMString(dec), VMBool(err == nil))
			return nil
		}
		return VMErrorNeedString
	}))

	env.DefineS("окр", VMFuncMustParams(2, func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		v1, ok1 := args[0].(VMDecNum)
		if !ok1 {
			return VMErrorNeedDecNum
		}
		v2, ok2 := args[1].(VMInt)
		if !ok2 {
			return VMErrorNeedInt
		}

		rets.Append(VMDecNum{num: v1.num.RoundWithMode(int32(v2), decnum.RoundHalfUp)})
		return nil
	}))

	env.DefineS("формат", VMFunc(func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		if len(args) < 2 {
			return VMErrorNeedFormatAndArgs
		}
		if v, ok := args[0].(VMString); ok {
			as := VMSlice(args[1:]).Args()
			rets.Append(VMString(env.Sprintf(string(v), as...)))
			return nil
		}
		return VMErrorNeedString
	}))

	env.DefineS("кодсимвола", VMFuncMustParams(1, func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		if v, ok := args[0].(VMStringer); ok {
			s := v.String()
			if len(s) == 0 {
				rets.Append(VMInt(0))
			} else {
				rets.Append(VMInt([]rune(s)[0]))
			}
			return nil
		}
		return VMErrorNeedString
	}))

	env.DefineS("типзнч", VMFuncMustParams(1, func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		if args[0] == nil || args[0] == VMNil {
			rets.Append(VMString("Неопределено"))
			return nil
		}
		rets.Append(VMString(names.UniqueNames.Get(env.TypeName(reflect.TypeOf(args[0])))))
		return nil
	}))

	env.DefineS("сообщить", VMFunc(func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		if len(args) == 0 {
			env.Println()
			return nil
		}
		as := args.Args()
		env.Println(as...)
		return nil
	}))

	env.DefineS("сообщитьф", VMFunc(func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		if len(args) < 2 {
			return VMErrorNeedFormatAndArgs
		}
		if v, ok := args[0].(VMString); ok {
			as := VMSlice(args[1:]).Args()
			env.Printf(string(v), as...)
			return nil
		}
		return VMErrorNeedString

	}))

	env.DefineS("обработатьгорутины", VMFuncMustParams(0, func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		runtime.Gosched()
		return nil
	}))

	env.DefineS("переменнаяокружения", VMFuncMustParams(1, func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		if v, ok := args[0].(VMString); ok {
			val, setted := os.LookupEnv(string(v))
			rets.Append(VMString(val))
			rets.Append(VMBool(setted))
			return nil
		}
		return VMErrorNeedSeconds
	}))

	// при изменении состава типов не забывать изменять их и в lexer.go
	env.DefineTypeS("целоечисло", ReflectVMInt)
	env.DefineTypeS("число", ReflectVMDecNum)
	env.DefineTypeS("булево", ReflectVMBool)
	env.DefineTypeS("строка", ReflectVMString)
	env.DefineTypeS("массив", ReflectVMSlice)
	env.DefineTypeS("структура", ReflectVMStringMap)
	env.DefineTypeS("дата", ReflectVMTime)
	env.DefineTypeS("длительность", ReflectVMTimeDuration)

	env.DefineTypeS("группаожидания", ReflectVMWaitGroup)
	env.DefineTypeS("файловаябазаданных", ReflectVMBoltDB)

	env.DefineTypeStruct("сервер", &VMServer{})
	env.DefineTypeStruct("клиент", &VMClient{})

	env.DefineTypeStruct("таблицазначений", &VMTable{})
	env.DefineTypeStruct("колонкатаблицызначений", &VMTableColumn{})
	env.DefineTypeStruct("коллекцияколоноктаблицызначений", &VMTableColumns{})
	env.DefineTypeStruct("строкатаблицызначений", &VMTableLine{})

	//////////////////
	env.DefineTypeStruct("__функциональнаяструктуратест__", &TttStructTest{})

	env.DefineS("__дамп__", VMFunc(func(args VMSlice, rets *VMSlice, envout *(*Env)) error {
		*envout = env
		env.Dump()
		return nil
	}))
	/////////////////////

	return env
}

/////////////////
// TttStructTest - тестовая структура для отладки работы с системными функциональными структурами
type TttStructTest struct {
	VMMetaObj

	ПолеЦелоеЧисло VMInt
	ПолеСтрока     VMString
}

func (tst *TttStructTest) VMRegister() {
	tst.VMRegisterMethod("ВСтроку", tst.ВСтроку)
	tst.VMRegisterField("ПолеЦелоеЧисло", &tst.ПолеЦелоеЧисло)
	tst.VMRegisterField("ПолеСтрока", &tst.ПолеСтрока)
}

// обратите внимание - русскоязычное название метода для структуры и формат для быстрого вызова
func (tst *TttStructTest) ВСтроку(args VMSlice, rets *VMSlice, envout *(*Env)) error {
	rets.Append(VMString(fmt.Sprintf("ПолеЦелоеЧисло=%v, ПолеСтрока=%v", tst.ПолеЦелоеЧисло, tst.ПолеСтрока)))
	return nil
}

/////////////////
