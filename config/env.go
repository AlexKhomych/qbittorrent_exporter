package config

import (
	"fmt"
	"os"
	"qbittorrent_exporter/lib/log"
	"reflect"
	"strconv"
	"strings"
)

// (v any) - must be a struct reference
func loadEnvs(v any) {
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	if typ.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		log.Warn("Unsupported type: " + typ.Kind().String())
		return
	}

	for i := range typ.NumField() {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		if envTag, ok := field.Tag.Lookup("env"); ok {
			log.Debug(fmt.Sprintf("Field: %s, Kind: %s, Addrssable: %v, CanSet: %v",
				field.Name, fieldValue.Kind(), fieldValue.CanAddr(), fieldValue.CanSet()))

			env := strings.ToLower(os.Getenv(envTag))
			if len(env) != 0 && fieldValue.CanSet() {
				if err := setFieldValue(fieldValue, env); err != nil {
					log.Error(fmt.Sprintf("Error setting field value: %v\n", err))
				} else {
					log.Debug(fmt.Sprintf("Updated [%s] with %s value", field.Name, env))
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
