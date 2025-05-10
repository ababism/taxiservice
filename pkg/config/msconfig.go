package msconfig

import (
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// ReplaceWithEnv заменяет все значения структуры на соответствующие значения из переменных окружения
func ReplaceWithEnv(config interface{}, prefix string) {
	v := reflect.ValueOf(config).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := field.Type()
		fieldName := v.Type().Field(i).Name

		if prefix != "" {
			fieldName = prefix + "_" + fieldName
		}

		if fieldType.Kind() == reflect.Ptr && fieldType.Elem().Kind() == reflect.Struct {
			// Рекурсивный вызов для вложенных структур, если не nil
			if !field.IsNil() {
				ReplaceWithEnv(field.Interface(), fieldName)
			}
		} else {
			envName := strings.ToUpper(fieldName)
			envValue := os.Getenv(envName)
			//fmt.Printf("%s = %s\n", envName, envValue)

			// Замена значения, если переменная окружения задана
			if envValue != "" && field.String() != envValue {
				log.Printf(
					"The configuration value of %s has been overwritten from '%s' to '%s'\n",
					envName,
					field.String(),
					envValue,
				)
				setField(&field, envValue)
			}
		}
	}
}

// setField устанавливает значение поля с учетом его типа
func setField(field *reflect.Value, value string) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			field.SetInt(intValue)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue, err := strconv.ParseUint(value, 10, 64)
		if err == nil {
			field.SetUint(uintValue)
		}
	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err == nil {
			field.SetFloat(floatValue)
		}
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err == nil {
			field.SetBool(boolValue)
		}
	default:
		return
	}
}
