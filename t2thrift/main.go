package parser

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

// m 管理命名空间
var m = map[string]string{}
// l 记录要通过 G() 翻译为 thrift 文本的类型
var l = []reflect.Type{}

func N(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Ptr:
		return "optional " + n(t)
	default:
		return "required " + n(t)
	}
}

// n 根据类型名生成 unique name
// 同时记录要跟踪的类型到 l
func n(t reflect.Type) string {
	var name string

	switch t.Kind() {
	case reflect.Array, reflect.Slice:
		return "list<" + n(t.Elem()) + "> "
	case reflect.Map:
		return fmt.Sprintf("map<%s,%s> ", n(t.Key()), n(t.Elem()))
	case reflect.Ptr:
		return n(t.Elem())
	case reflect.Struct:
		a := []byte(t.PkgPath())
		for k, b := range a {
			if !((b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')) {
				a[k] = '_'
			}
		}
		aa := strings.Split(string(a), "_")
		name = camelString(aa[len(aa)-1]) + t.Name()
		for i := len(aa) - 2; i >= 0; i-- {
			// 命名空间碰撞检测
			// 根据 v 是否等于 t.PkgPath() 来避免一个 t 生成 2 个name
			if v, ok := m[name]; ok {
				if v == t.PkgPath(){
					break
				}
			}else{
				m[name] = t.PkgPath()
				l = append(l, t)
				break
			}
			name += camelString(aa[i])
		}
	case reflect.String:
		name = "string"
	case reflect.Int8:
		name = "byte"
	case reflect.Int, reflect.Int64:
		name = "i64"
	case reflect.Int32:
		name = "i32"
	case reflect.Int16:
		name = "i16"
	case reflect.Float32, reflect.Float64:
		name = "double"
	case reflect.Bool:
		name = "bool"
	default:
		panic("bad case")
	}

	return name
}

// G 根据 type 生成 thrift 文本
func G(t reflect.Type) string {
	if len(t.Name()) > 4 && t.Name()[0:4] == "Enum" {
		o := "enum " + n(t) + "{\n"
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			o += fmt.Sprintf("\t%s\n", f.Name)
		}
		o += "}\n"
		return o
	} else {
		o := "struct " + n(t) + "{\n"
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			o += fmt.Sprintf("\t%d: %s %s\n", i, N(f.Type), f.Name)
		}
		o += "}\n"
		return o
	}
}

func Parse(service interface{}) {
	var output = []string{}

	{
		o := "service "
		t := reflect.TypeOf(service)

		o += t.Name() + "{\n"
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			o += fmt.Sprintf("\t%s %s(", n(f.Type.Out(0)), f.Name)
			// func params
			for i := 0; i < f.Type.NumIn(); i++ {
				p := f.Type.In(i)
				o += fmt.Sprintf("%d: %s arg%d,", i, n(p), i)
			}
			// thrift anno
			o += `)(` + strings.ReplaceAll(strings.ReplaceAll(string(f.Tag)+" ", `:"`, `='`), `" `, "',") + ")\n"
		}
		o += `}`
		output = append(output, o)
	}

	{
		var l2 []reflect.Type
		for len(l) > 0 {
			l, l2 = []reflect.Type{}, l
			for _, t := range l2 {
				output = append(output, G(t))
			}
		}
	}

	f, err := os.OpenFile("gen.thrift", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for i := len(output) - 1; i >= 0; i-- {
		if _, err = f.WriteString(output[i]); err != nil {
			panic(err)
		}
	}
}
