package env

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

const tagName = `env`

// Parses options and field name. Returning name and slice of opts
func parseTag(tag string) (string, []string) {
	if i := strings.Index(tag, ","); i != -1 {
		return tag[:i], strings.Split(tag[i+1:], ",")
	}
	return tag, []string{}
}

// Helper function that splits on first '='
func splitEnv(in string) (string, string, error) {
	splits := strings.Split(in, "=")
	if len(splits) == 1 {
		return "", "", fmt.Errorf("invalid input for splitting env variables")
	}
	key := splits[0]
	val := strings.Join(splits[1:], "=")
	return key, val, nil
}

// Takes a list of strings in format key=val and splits to a map
func parse(data []string) (map[string]string, error) {
	items := make(map[string]string)
	for _, item := range data {
		key, val, err := splitEnv(item)
		if err != nil {
			return items, err
		}
		items[key] = val
	}
	return items, nil
}

// Takes a list of strings in format key=value and unmarshals to struct
func Unmarshal(data []string, out interface{}) error {
	envVars, err := parse(data)
	if err != nil {
		return err
	}

	typeStruct := reflect.TypeOf(out).Elem()
	valueStruct := reflect.ValueOf(out).Elem()

	for i := 0; i < valueStruct.NumField(); i++ {
		tn, ok := typeStruct.Field(i).Tag.Lookup(tagName)
		if !ok {
			// ToDo, do type assertion on struct field and see if value exists with key=FieldName
			log.Printf("Error! The package is not inferring by struct type yet. Skipping: %s", tagName)
			continue
		}
		name, opt := parseTag(tn)

		val, ok := envVars[name]
		if len(opt) == 0 || opt[0] == "string" {
			if !ok {
				return fmt.Errorf("can not find required value: '%s'", name)
			}

			// We assume string for now, see comment above
			if valueStruct.Field(i).IsValid() && valueStruct.Field(i).Kind() == reflect.String {
				valueStruct.Field(i).SetString(val)
			}
			continue
		}
		switch opt[0] {
		case "int64":
			v, err := strconv.ParseInt(envVars[name], 10, 64)
			if err != nil {
				return err
			}
			valueStruct.Field(i).SetInt(v)
		case "uint64":
			v, err := strconv.ParseUint(envVars[name], 10, 64)
			if err != nil {
				return err
			}
			valueStruct.Field(i).SetUint(v)
		case "bool":
			v, err := strconv.ParseBool(envVars[name])
			if err != nil {
				return err
			}
			valueStruct.Field(i).SetBool(v)
		}
	}

	return nil
}
