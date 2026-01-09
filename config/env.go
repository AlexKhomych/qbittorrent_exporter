package config

import (
	"fmt"
	"os"
	"qbittorrent_exporter/lib/log"
	"reflect"
	"strconv"
	"strings"
)

func loadEnvs(v any) {
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	if typ.Kind() == reflect.Pointer {
		val = val.Elem()
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		log.Warn("loadEnvs: unsupported type " + typ.Kind().String())
		return
	}

	for i := range typ.NumField() {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		if envTag, ok := field.Tag.Lookup("env"); ok {
			env := strings.ToLower(os.Getenv(envTag))
			if len(env) != 0 && fieldValue.CanSet() {
				if err := setFieldValue(fieldValue, env); err != nil {
					log.Error(fmt.Sprintf("Failed to set %s from env %s: %v", field.Name, envTag, err))
				} else {
					log.Debug(fmt.Sprintf("Loaded %s from environment variable", field.Name))
				}
			}
		}

		if fieldValue.Kind() == reflect.Struct {
			loadEnvs(fieldValue.Addr().Interface())
		}
	}
}

func setFieldValue(fieldValue reflect.Value, envTag string) error {
	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(envTag)
	case reflect.Int:
		intVal, err := strconv.Atoi(envTag)
		if err != nil {
			return err
		}
		fieldValue.SetInt(int64(intVal))
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(envTag)
		if err != nil {
			return err
		}
		fieldValue.SetBool(boolVal)
	case reflect.Float64:
		floatVal, err := strconv.ParseFloat(envTag, 64)
		if err != nil {
			return err
		}
		fieldValue.SetFloat(floatVal)
	default:
		return fmt.Errorf("unsupported field type: %s", fieldValue.Kind())
	}
	return nil
}
