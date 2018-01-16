package graphql_reflect

import (
	"github.com/graphql-go/graphql"
	"reflect"
)

const (
	TYPE_FIELD = "gqlType"
)

//Каждый тип должен участвовать в схеме ровно один раз
//В том числе, если в каком-то другом типе этот тип участвует как поле — это должен быть ровно тот же экземлпяр
//Для этого нужен кеш: при первом объявлении типа кладем его в кеш, а затем отдаем _тот_же_самый_ экземпляр
var typeCash map[string]*graphql.Object = make(map[string]*graphql.Object)

func ReflectType(obj interface{}) *graphql.Object {
	return doReflect(reflect.TypeOf(obj))
}

func doReflect(t reflect.Type) *graphql.Object {
	var (
		name, description, nameField string
		fields                       map[string]*graphql.Field
		tmp                          *graphql.Object
	)

	//По умолчанию: Name — по названию типа (структуры) в коде, а Description не указан
	name = t.Name()
	description = "Без описания :("

	//Name и Description можно уточнить через тег на специальном поле GQLType
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Name == TYPE_FIELD && getTag(t.Field(i)) != "-" {
			name = extractName(t.Field(i))
			description = extractDescription(t.Field(i))
		}

	}

	//Здесь мы определили название типа — и попробуем найти его в кеше
	if cashed, exists := typeCash[name]; exists {
		return cashed
	}

	//Теперь собираем поля
	fields = make(map[string]*graphql.Field)
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Name != TYPE_FIELD && getTag(t.Field(i)) != "-" {
			var f graphql.Field

			if t.Field(i).Type.Kind() == reflect.Struct {
				f.Type = doReflect(t.Field(i).Type)
			} else if t.Field(i).Type.Kind() == reflect.Slice && t.Field(i).Type.Elem().Kind() == reflect.Struct {
				f.Type = graphql.NewList(doReflect(t.Field(i).Type.Elem()))
			} else {
				f.Type = getGraphType(t.Field(i).Type)
			}

			//По умолчанию: Name — имя поля в структуре, а Description не указан
			nameField = t.Field(i).Name
			f.Description = ""

			//Уточнить можно через тег на поле
			if getTag(t.Field(i)) != "" {
				nameField = extractName(t.Field(i))
				f.Description = extractDescription(t.Field(i))
			}

			//...и добавляем описание поля в тип
			fields[nameField] = &f
		}

	}

	tmp = graphql.NewObject(graphql.ObjectConfig{Name: name, Description: description, Fields: graphql.Fields(fields)})
	typeCash[name] = tmp
	return tmp
}
