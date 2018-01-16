package graphql_reflect

import (
	"reflect"
	"strings"
)

const (
	TAG = "gql"
)

func getTag(f reflect.StructField) string {
	return f.Tag.Get(TAG)
}

func extractName(f reflect.StructField) string {
	return strings.Split(getTag(f), ",")[0]
}

func extractDescription(f reflect.StructField) string {
	var (
		tmp         []string
		description string
	)

	tmp = strings.Split(getTag(f), ",")
	for i := 1; i < len(tmp); i++ {
		description += tmp[i]
	}

	return description
}
