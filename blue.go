// Package blue is a small package for my personal use. It simply converts PUT request payload to bson.M.
package blue

import (
	"reflect"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

// TagName is the tag name that used for defining options on structs
const TagName = "blue"

// Spray gets a struct and converts it to bson.M
func Spray(s interface{}) bson.M {
	result := bson.M{}

	fields := extractFields(s)
	values := extractValues(s)

	for _, field := range fields {
		name, omitempty := scanField(field)
		val := values.FieldByName(name)

		if omitempty {
			zero := reflect.Zero(val.Type()).Interface()
			current := val.Interface()

			if reflect.DeepEqual(current, zero) {
				continue
			}
		}

		result[name] = val.Interface()
		// TODO : current approach does not work with nested object, It should get fixed.
	}

	return result
}

// extractValues extracts the values from given struct using reflect package
func extractValues(s interface{}) reflect.Value {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return v
}

// extractFields extracts fields from given struct using reflect package
func extractFields(s interface{}) []reflect.StructField {
	var fields []reflect.StructField

	vals := extractValues(s)
	t := vals.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if tag := field.Tag.Get(TagName); tag == "-" {
			continue
		}

		fields = append(fields, field)
	}

	return fields
}

// scanField read the tags of the field, returns name and omitempty flag. Will use property name if
// the name does ot exists in tag.
func scanField(field reflect.StructField) (string, bool) {
	var (
		name      string
		omitempty bool
	)

	name = field.Name
	tag := field.Tag.Get(TagName)

	parts := strings.Split(tag, ",")

	// search for user defined name, otherwise use property name
	if len(parts) > 0 {
		if parts[0] != "" {
			name = parts[0]
		}
	}

	// check how should behave with empty fields
	if len(parts) > 1 {
		if parts[1] == "omitempty" {
			omitempty = true
		}
	}

	return name, omitempty
}
