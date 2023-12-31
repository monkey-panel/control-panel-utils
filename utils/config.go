package utils

import "reflect"

const globalConfigPath = "data/local_data"

type BaseConfig interface {
	Default() any
}

func ConfigRead(path string, config BaseConfig) error {
	data := map[string]any{}
	ReadJsonFile(globalConfigPath, &data)
	if err := WriteJsonFile(
		globalConfigPath,
		mergeObj(data, toMap(config.Default())),
	); err != nil {
		return err
	}
	ReadJsonFile(globalConfigPath, &config)
	return nil
}

func toMap(data any) map[string]any {
	result := map[string]any{}
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get("json"); tagValue != "" {
			result[tagValue] = v.Field(i).Interface()
		}
	}
	return result
}

// merge two object
func mergeObj(old, new map[string]any) map[string]any {
	for k, value := range new {
		if oldValue, ok := old[k]; ok {
			oldValueType := reflect.TypeOf(oldValue).Kind()
			valueType := reflect.TypeOf(value).Kind()
			if valueType == reflect.Map && oldValueType == reflect.Map {
				// if is dict
				old[k] = mergeObj(oldValue.(map[string]any), value.(map[string]any))
				continue
			} else if valueType == reflect.Slice && oldValueType == reflect.Slice {
				// if is list
				continue
			} else if oldValueType == valueType {
				continue
				// if is number
			} else if isNumber(oldValueType) && isNumber(valueType) {
				continue
			}
		}
		old[k] = value
	}

	return old
}

func isNumber(kind reflect.Kind) bool {
	switch kind {
	case
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64:
		return true
	}
	return false
}
