package config_reader

import (
	"os"
	"reflect"
	"strconv"
	"strings"
)

func ReadEnv(config interface{}) {
	readEnv(reflect.ValueOf(config).Elem())
}

func readEnv(val reflect.Value) {
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := field.Kind()
		tagValue, exist := val.Type().Field(i).Tag.Lookup("env")
		if fieldType == reflect.Struct {
			readEnv(field)
		}
		if !exist {
			continue
		}
		configValue, ok := os.LookupEnv(tagValue)
		if !ok {
			continue
		}
		if fieldType == reflect.String {
			field.SetString(configValue)
		} else if fieldType == reflect.Int {
			if val, ok := strconv.ParseInt(configValue, 10, 64); ok == nil {
				field.SetInt(val)
			}
		} else if fieldType == reflect.Bool {
			if strings.Contains(configValue, "true") || strings.Contains(configValue, "1") {
				field.SetBool(true)
			} else {
				field.SetBool(false)
			}
		}
	}
}
