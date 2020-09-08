package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"gopkg.in/yaml.v2"
)

func Parse(data []byte, v interface{}) error {
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return fmt.Errorf("second argument should be pointer")
	}

	if data == nil {
		return parseENV(v, true)
	}

	switch true {
	case json.Unmarshal(data, v) == nil:
		return parseENV(v, false)
	case yaml.Unmarshal(data, v) == nil:
		return parseENV(v, false)
	default:
		return errors.New("unsupported config file type")
	}

	return nil
}
