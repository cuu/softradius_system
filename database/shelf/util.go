package shelf

import "github.com/pkg4go/camelcase"
import "github.com/pkg4go/convert"
import "reflect"

func Type(v interface{}) string {
	return reflect.ValueOf(v).Type().Name()
}

func getTableName(args ...interface{}) string {
	name := camelcase.Reverse(Type(args[0]))

	if len(args) == 1 {
		return name
	}

	if len(args) == 2 {
		if n := convert.String(args[1]); n != "" {
			return n
		}
	}

	return name
}
