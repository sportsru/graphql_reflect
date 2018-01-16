package graphql_reflect

/*
Код взят отсюда https://github.com/graphql-go/graphql/blob/master/util.go
*/

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"reflect"
	"strings"
)

func getGraphType(tipe reflect.Type) graphql.Output {
	kind := tipe.Kind()
	switch kind {
	case reflect.String:
		return graphql.String
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		return graphql.Int
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		return graphql.Float
	case reflect.Bool:
		return graphql.Boolean
	case reflect.Slice:
		return getGraphList(tipe)
	}
	return graphql.String
}

func getGraphList(tipe reflect.Type) *graphql.List {
	if tipe.Kind() == reflect.Slice {
		switch tipe.Elem().Kind() {
		case reflect.Int:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			return graphql.NewList(graphql.Int)
		case reflect.Bool:
			return graphql.NewList(graphql.Boolean)
		case reflect.Float32:
			fallthrough
		case reflect.Float64:
			return graphql.NewList(graphql.Float)
		case reflect.String:
			return graphql.NewList(graphql.String)
		}
	}
	// finaly bind object
	//t := reflect.New(tipe.Elem())
	name := strings.Replace(fmt.Sprint(tipe.Elem()), ".", "_", -1)
	obj := graphql.NewObject(graphql.ObjectConfig{
		Name:   name,
		Fields: nil, //BindFields(t.Elem().Interface()),
	})
	return graphql.NewList(obj)
}
