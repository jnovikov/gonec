package core

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/covrom/gonec/names"
)

// Env provides interface to run VM. This mean function scope and blocked-scope.
// If stack goes to blocked-scope, it will make new Env.
type Env struct {
	sync.RWMutex
	name         string
	env          map[int]VMValuer
	typ          map[int]reflect.Type
	parent       *Env
	interrupt    *bool
	stdout       io.Writer
	sid          string
	goRunned     bool
	lastid       int
	lastval      VMValuer
	builtsLoaded bool
}

func (e *Env) vmval() {} // нужно для того, чтобы *Env можно было сохранять в переменные VMValuer

// NewEnv creates new global scope.
// !!!не забывать вызывать core.LoadAllBuiltins(m)!!!
func NewEnv() *Env {
	b := false

	m := &Env{
		env:          make(map[int]VMValuer),
		typ:          make(map[int]reflect.Type),
		parent:       nil,
		interrupt:    &b,
		stdout:       os.Stdout,
		goRunned:     false,
		lastid:       -1,
		builtsLoaded: false,
	}
	return m
}

// NewEnv создает новое окружение под глобальным контекстом переданного в e
func (e *Env) NewEnv() *Env {
	for ee := e; ee != nil; ee = ee.parent {
		if ee.parent == nil {
			return &Env{
				env:          make(map[int]VMValuer),
				typ:          make(map[int]reflect.Type),
				parent:       ee,
				interrupt:    e.interrupt,
				stdout:       e.stdout,
				goRunned:     false,
				lastid:       -1,
				builtsLoaded: ee.builtsLoaded,
			}

		}
	}
	panic("Не найден глобальный контекст!")
}

// NewSubEnv создает новое окружение под e, нужно для замыкания в анонимных функциях
func (e *Env) NewSubEnv() *Env {
	return &Env{
		env:          make(map[int]VMValuer),
		typ:          make(map[int]reflect.Type),
		parent:       e,
		interrupt:    e.interrupt,
		stdout:       e.stdout,
		goRunned:     false,
		lastid:       -1,
		builtsLoaded: e.builtsLoaded,
	}
}

// Находим или создаем новый модуль в глобальном скоупе
func (e *Env) NewModule(n string) *Env {
	//ni := strings.ToLower(n)
	id := names.UniqueNames.Set(n)
	if v, err := e.Get(id); err == nil {
		if vv, ok := v.(*Env); ok {
			return vv
		}
	}

	m := e.NewEnv()
	m.name = n

	// на модуль можно ссылаться через переменную породившего глобального контекста
	e.DefineGlobal(id, m)
	return m
}

// func NewPackage(n string, w io.Writer) *Env {
// 	b := false

// 	return &Env{
// 		env:       make(map[string]reflect.Value),
// 		typ:       make(map[string]reflect.Type),
// 		parent:    nil,
// 		name:      strings.ToLower(n),
// 		interrupt: &b,
// 		stdout:    w,
// 	}
// }

func (e *Env) NewPackage(n string) *Env {
	return &Env{
		env:          make(map[int]VMValuer),
		typ:          make(map[int]reflect.Type),
		parent:       e,
		name:         strings.ToLower(n),
		interrupt:    e.interrupt,
		stdout:       e.stdout,
		goRunned:     false,
		lastid:       -1,
		builtsLoaded: e.builtsLoaded,
	}
}

// Destroy deletes current scope.
func (e *Env) Destroy() {
	if e.goRunned {
		e.Lock()
		defer e.Unlock()
	}

	if e.parent == nil {
		return
	}

	for k, v := range e.parent.env {
		if vv, ok := v.(*Env); ok {
			if vv == e {
				delete(e.parent.env, k)
			}
		}
	}
	e.parent = nil
	e.env = nil
}

func (e *Env) SetGoRunned(t bool) {
	for ee := e; ee != nil; ee = ee.parent {
		ee.Lock()
		ee.goRunned = t
		ee.Unlock()
	}
}

func (e *Env) SetBuiltsIsLoaded() {
	if e.goRunned {
		e.Lock()
	}
	e.builtsLoaded = true
	if e.goRunned {
		e.Unlock()
	}
}

func (e *Env) IsBuiltsLoaded() bool {
	for ee := e; ee != nil; ee = ee.parent {
		if ee.goRunned {
			ee.RLock()
			defer ee.RUnlock()
		}
		if ee.builtsLoaded {
			return true
		}
	}
	return false
}

// SetName sets a name of the scope. This means that the scope is module.
func (e *Env) SetName(n string) {
	if e.goRunned {
		e.Lock()
	}
	e.name = strings.ToLower(n)
	if e.goRunned {
		e.Unlock()
	}
}

// GetName returns module name.
func (e *Env) GetName() string {
	if e.goRunned {
		e.RLock()
		defer e.RUnlock()
	}

	return e.name
}

// TypeName определяет имя типа по типу значения
func (e *Env) TypeName(t reflect.Type) int {

	for ee := e; ee != nil; ee = ee.parent {
		if ee.goRunned {
			ee.RLock()
			defer ee.RUnlock()
		}
		for k, v := range ee.typ {
			if v == t {
				return k
			}
		}
	}
	return names.UniqueNames.Set(t.String())
}

// Type returns type which specified symbol. It goes to upper scope until
// found or returns error.
func (e *Env) Type(k int) (reflect.Type, error) {

	for ee := e; ee != nil; ee = ee.parent {
		if ee.goRunned {
			ee.RLock()
			defer ee.RUnlock()
		}
		if v, ok := ee.typ[k]; ok {
			return v, nil
		}
	}
	return nil, fmt.Errorf("Тип неопределен '%s'", names.UniqueNames.Get(k))
}

// Get returns value which specified symbol. It goes to upper scope until
// found or returns error.
func (e *Env) Get(k int) (VMValuer, error) {

	for ee := e; ee != nil; ee = ee.parent {
		if ee.goRunned {
			ee.RLock()
			defer ee.RUnlock()
		}
		if e.lastid == k {
			// это именно здесь, т.к. нужно учесть блокировку
			return e.lastval, nil
		}
		if v, ok := ee.env[k]; ok {
			return v, nil
		}
	}
	return nil, fmt.Errorf("Имя неопределено '%s'", names.UniqueNames.Get(k))
}

// Set modifies value which specified as symbol. It goes to upper scope until
// found or returns error.
func (e *Env) Set(k int, v VMValuer) error {

	for ee := e; ee != nil; ee = ee.parent {
		if ee.goRunned {
			ee.Lock()
			defer ee.Unlock()
		}
		if _, ok := ee.env[k]; ok {
			ee.env[k] = v
			e.lastid = k
			e.lastval = v
			return nil
		}
	}
	return fmt.Errorf("Имя неопределено '%s'", names.UniqueNames.Get(k))
}

// DefineGlobal defines symbol in global scope.
func (e *Env) DefineGlobal(k int, v VMValuer) error {
	for ee := e; ee != nil; ee = ee.parent {
		if ee.parent == nil {
			return ee.Define(k, v)
		}
	}
	return fmt.Errorf("Отсутствует глобальный контекст!")
}

// DefineType defines type which specifis symbol in global scope.
func (e *Env) DefineType(k int, t reflect.Type) error {
	for ee := e; ee != nil; ee = ee.parent {
		if ee.parent == nil {
			if ee.goRunned {
				ee.Lock()
				defer ee.Unlock()
			}
			ee.typ[k] = t
			// // пишем в кэш индексы полей и методов для структур
			// // для работы со структурами нам нужен конкретный тип
			// if typ.Kind() == reflect.Ptr {
			// 	typ = typ.Elem()
			// }
			// if typ.Kind() == reflect.Struct {
			// 	// методы берем в т.ч. у ссылки на структуру, они включают методы самой структуры
			// 	// это будут разные методы для разных reflect.Value
			// 	ptyp := reflect.TypeOf(reflect.New(typ).Interface())
			// 	basicpath := typ.PkgPath() + "." + typ.Name() + "."

			// 	//методы
			// 	nm := typ.NumMethod()
			// 	for i := 0; i < nm; i++ {
			// 		meth := typ.Method(i)
			// 		// только экспортируемые
			// 		if meth.PkgPath == "" {
			// 			namtyp := UniqueNames.Set(basicpath + meth.Name)
			// 			// fmt.Println("SET METHOD: "+basicpath+meth.Name, meth.Index)
			// 			// ast.StructMethodIndexes.Cache[namtyp] = meth.Index
			// 		}
			// 	}
			// 	nm = ptyp.NumMethod()
			// 	for i := 0; i < nm; i++ {
			// 		meth := ptyp.Method(i)
			// 		// только экспортируемые
			// 		if meth.PkgPath == "" {
			// 			namtyp := UniqueNames.Set(basicpath + "*" + meth.Name)
			// 			// fmt.Println("SET *METHOD: "+basicpath+"*"+meth.Name, meth.Index)
			// 			// ast.StructMethodIndexes.Cache[namtyp] = meth.Index
			// 		}
			// 	}

			// 	//поля
			// 	nm = typ.NumField()
			// 	for i := 0; i < nm; i++ {
			// 		field := typ.Field(i)
			// 		// только экспортируемые неанонимные поля
			// 		if field.PkgPath == "" && !field.Anonymous {
			// 			namtyp := UniqueNames.Set(basicpath + field.Name)
			// 			// fmt.Println("SET FIELD: "+basicpath+field.Name, field.Index)
			// 			// ast.StructFieldIndexes.Cache[namtyp] = field.Index
			// 		}
			// 	}
			// }
			return nil
		}
	}
	return fmt.Errorf("Отсутствует глобальный контекст!")
}

func (e *Env) DefineTypeS(k string, t reflect.Type) error {
	return e.DefineType(names.UniqueNames.Set(k), t)
}

// DefineTypeStruct регистрирует системную функциональную структуру, переданную в виде указателя!
func (e *Env) DefineTypeStruct(k string, t interface{}) error {
	gob.Register(t)
	return e.DefineType(names.UniqueNames.Set(k), reflect.Indirect(reflect.ValueOf(t)).Type())
}

// Define defines symbol in current scope.
func (e *Env) Define(k int, v VMValuer) error {
	if e.goRunned {
		e.Lock()
	}
	e.env[k] = v
	e.lastid = k
	e.lastval = v

	if e.goRunned {
		e.Unlock()
	}

	return nil
}

func (e *Env) DefineS(k string, v VMValuer) error {
	return e.Define(names.UniqueNames.Set(k), v)
}

// String return the name of current scope.
func (e *Env) String() string {
	if e.goRunned {
		e.RLock()
		defer e.RUnlock()
	}
	return e.name
}

// Dump show symbol values in the scope.
func (e *Env) Dump() {
	if e.goRunned {
		e.RLock()
	}
	var keys []int
	for k := range e.env {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		e.Printf("%d %s = %#v %T\n", k, names.UniqueNames.Get(k), e.env[k], e.env[k])
	}
	if e.goRunned {
		e.RUnlock()
	}
}

func (e *Env) Println(a ...interface{}) (n int, err error) {
	// e.RLock()
	// defer e.RUnlock()
	return fmt.Fprintln(e.stdout, a...)
}

func (e *Env) Printf(format string, a ...interface{}) (n int, err error) {
	// e.RLock()
	// defer e.RUnlock()
	return fmt.Fprintf(e.stdout, format, a...)
}

func (e *Env) Sprintf(format string, a ...interface{}) string {
	// e.RLock()
	// defer e.RUnlock()
	return fmt.Sprintf(format, a...)
}

func (e *Env) Print(a ...interface{}) (n int, err error) {
	// e.RLock()
	// defer e.RUnlock()
	return fmt.Fprint(e.stdout, a...)
}

func (e *Env) StdOut() reflect.Value {
	// e.RLock()
	// defer e.RUnlock()
	return reflect.ValueOf(e.stdout)
}

func (e *Env) SetStdOut(w io.Writer) {
	// e.Lock()
	//пренебрегаем возможными коллизиями при установке потока вывода, т.к. это совсем редкая операция
	e.stdout = w
	// e.Unlock()
}

func (e *Env) SetSid(s string) error {
	for ee := e; ee != nil; ee = ee.parent {
		if ee.parent == nil {
			ee.sid = s
			return ee.Define(names.UniqueNames.Set("ГлобальныйИдентификаторСессии"), VMString(s))
		}
	}
	return fmt.Errorf("Отсутствует глобальный контекст!")
}

func (e *Env) GetSid() string {
	for ee := e; ee != nil; ee = ee.parent {
		if ee.parent == nil {
			// пренебрегаем возможными коллизиями, т.к. изменение номера сессии - это совсем редкая операция
			return ee.sid
		}
	}
	return ""
}

func (e *Env) Interrupt() {
	e.Lock()
	*(e.interrupt) = true
	e.Unlock()
}

func (e *Env) CheckInterrupt() bool {
	if e.goRunned {
		e.Lock()
		defer e.Unlock()
	}
	if *(e.interrupt) {
		*(e.interrupt) = false
		return true
	}
	return false
}
