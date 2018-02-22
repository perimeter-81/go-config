package config

import (
	"encoding/json"
	"errors"

	"gopkg.in/yaml.v2"
)

func Parse(data []byte, v interface{}) (err error) {
	if data == nil {
		return parseENV(v)
	}

	switch true {
	case json.Unmarshal(data, v) == nil:
		return parseENV(v)
	case yaml.Unmarshal(data, v) == nil:
		return parseENV(v)
	default:
		return errors.New("unsupported config file type")
	}

	return nil
}
