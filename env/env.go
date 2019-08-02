package env

import (
	"encoding/base64"
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
// Infers from the tagname, if not available chooses field name in all caps
func Unmarshal(data []string, out interface{}) error {
	envVars, err := parse(data)
	if err != nil {
		return err
	}

	typeStruct := reflect.TypeOf(out).Elem()
	valueStruct := reflect.ValueOf(out).Elem()

	for i := 0; i < valueStruct.NumField(); i++ {
		tn, ok := typeStruct.Field(i).Tag.Lookup(tagName)
		name, opt := parseTag(tn)
		if !ok {
			name = strings.ToUpper(typeStruct.Field(i).Name)
			opt = []string{}
		}
		val, ok := envVars[name]
		if !ok && name != "-" {
			return fmt.Errorf("could not find variable: '%s'", name)
		}
		switch typeStruct.Field(i).Type.String() {

		case "string":
			valueStruct.Field(i).SetString(val)

		case "uint64":
			u, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				return fmt.Errorf("could not parse uint: %s for field: %s", val, name)
			}
			valueStruct.Field(i).SetUint(u)

		case "int64":
			in, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return fmt.Errorf("could not parse int64: %s for field: %s", val, name)
			}
			valueStruct.Field(i).SetInt(in)

		case "int":
			in, err := strconv.Atoi(val)
			if err != nil {
				return fmt.Errorf("could not parse int: %s for field: %s", val, name)
			}
			valueStruct.Field(i).Set(reflect.ValueOf(in))

		case "bool":
			b, err := strconv.ParseBool(val)
			if err != nil {
				return fmt.Errorf("could not parse bool: %s for field: %s", val, name)
			}
			valueStruct.Field(i).SetBool(b)

		case "map[string]bool":
			m := make(map[string]bool)
			strSlice := strings.Split(envVars[name], ";")

			for _, subStr := range strSlice {
				kv := strings.Split(subStr, ":")
				if len(kv) != 2 {
					return fmt.Errorf("'%s' is formatted incorrectly", typeStruct.Field(i).Name)
				}
				v, err := strconv.ParseBool(kv[1])
				if err != nil {
					return fmt.Errorf("can not find required value: '%s'", typeStruct.Field(i).Name)
				}

				m[kv[0]] = v
			}
			valueStruct.Field(i).Set(reflect.ValueOf(m))

		case "[]uint8":
			if len(opt) > 0 {
				if opt[0] == "b64" { // Base64 decode the value

					bytes, err := base64.StdEncoding.DecodeString(envVars[name])
					if err != nil {
						return fmt.Errorf("could not decode base64 from env var: %s", name)
					}
					valueStruct.Field(i).SetBytes(bytes)
				}
			}

		}

	}
	log.Printf("Data: %+v", out)
	return nil
}
