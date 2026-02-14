package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(person Person) string {
	v := reflect.ValueOf(person)
	t := reflect.TypeOf(person)
	res := []string{}
	num := v.NumField()
	for i := range num {
		tag, ok := t.Field(i).Tag.Lookup("properties")
		if !ok {
			continue
		}
		name := tag
		oe := false
		if strings.Contains(tag, "omitempty") {
			name = strings.Replace(name, ",omitempty", "", 1)
			oe = true
		}
		fmt.Println(name, oe, v.Field(i).IsZero())
		if oe && v.Field(i).IsZero() {
			continue
		}

		if !oe && v.Field(i).IsZero() {
			res = append(res, name+"="+emptyValue(t.Field(i)))
			continue
		}

		res = append(res, name+"="+nonEmptyValue(t.Field(i), v.Field(i)))

		// fmt.Println(v.Field(i))
		// fmt.Println(tag)
		// fmt.Println(v.Field(i).Kind().String())
		// fmt.Println(reflect.Zero(t.Field(i).Type))
	}
	return strings.Join(res, "\n")
}

func emptyValue(t reflect.StructField) string {
	// рассматриваю не все типы, только те что есть в задании, чисто для ускорения
	switch t.Type.Kind() {
	case reflect.Int:
		return "0"
	case reflect.Bool:
		return "false"
	default:
		// по хорошему нужны ещё пустые массивы и объекты
		return ""
	}
}

func nonEmptyValue(t reflect.StructField, v reflect.Value) string {
	switch t.Type.Kind() {
	case reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(int64(v.Int()), 10)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return v.String()
	default:
		return ""
	}
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
