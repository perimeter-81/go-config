package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	fieldTag        = "env"
	fieldDefaultTag = "default"
	separator       = ","
)

var (
	// ErrNotAStructPtr is returned if you pass something that is not a pointer to a Struct to Parse
	ErrNotAStructPtr = errors.New("expected a pointer to a Struct")

	// ErrUnsupportedType if the struct field type is not supported by env
	ErrUnsupportedType = errors.New("type is not supported")

	// ErrUnsupportedSliceType if the slice element type is not supported by env
	ErrUnsupportedSliceType = errors.New("unsupported slice type")

	// Friendly names for reflect types
	sliceOfInts     = reflect.TypeOf([]int(nil))
	sliceOfInt64s   = reflect.TypeOf([]int64(nil))
	sliceOfStrings  = reflect.TypeOf([]string(nil))
	sliceOfBools    = reflect.TypeOf([]bool(nil))
	sliceOfFloat32s = reflect.TypeOf([]float32(nil))
	sliceOfFloat64s = reflect.TypeOf([]float64(nil))
)

// parse parses a struct containing `env` tags and loads its values from
// environment variables.
func parseENV(v interface{}) error {
	ptrRef := reflect.ValueOf(v)
	if ptrRef.Kind() != reflect.Ptr {
		return ErrNotAStructPtr
	}

	ref := ptrRef.Elem()
	if ref.Kind() != reflect.Struct {
		return ErrNotAStructPtr
	}

	return doParse(ref)
}

func doParse(ref reflect.Value) error {
	refType := ref.Type()
	for i := 0; i < ref.NumField(); i++ {
		switch ref.Field(i).Kind() {
		case reflect.Struct:
			if err := doParse(ref.Field(i)); err != nil {
				return err
			}
		default:
			value, err := get(refType.Field(i))
			if err != nil {
				return err
			}

			if value == "" {
				continue
			}

			if err = set(ref.Field(i), refType.Field(i), value); err != nil {
				return err
			}
		}
	}

	return nil
}

func get(field reflect.StructField) (val string, err error) {
	key, opts := parseKeyForOption(field.Tag.Get(fieldTag))

	defaultValue := field.Tag.Get(fieldDefaultTag)
	val = getOr(key, defaultValue)

	if len(opts) > 0 {
		for _, opt := range opts {
			// The only option supported is "required".
			switch opt {
			case "":
				break
			case "required":
				val, err = getRequired(key)
			default:
				err = fmt.Errorf("env tag option %s not supported", opt)
			}
		}
	}

	return val, err
}

// split the env tag's key into the expected key and desired option, if any.
func parseKeyForOption(key string) (string, []string) {
	opts := strings.Split(key, separator)

	return opts[0], opts[1:]
}

func getRequired(key string) (string, error) {
	if value := os.Getenv(key); value != "" {
		return value, nil
	}

	return "", fmt.Errorf("required environment variable %s is not set", key)
}

func getOr(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}

	return defaultValue
}

func set(field reflect.Value, refType reflect.StructField, value string) error {
	switch field.Kind() {
	case reflect.Slice:
		sep := refType.Tag.Get("envSeparator")

		return handleSlice(field, value, sep)
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		bVal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}

		field.SetBool(bVal)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		intValue, err := strconv.ParseInt(value, 10, field.Type().Bits())
		if err != nil {
			return err
		}

		field.SetInt(intValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		intValue, err := strconv.ParseUint(value, 10, field.Type().Bits())
		if err != nil {
			return err
		}

		field.SetUint(intValue)
	case reflect.Float32:
		v, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}

		field.SetFloat(v)
	case reflect.Float64:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}

		field.Set(reflect.ValueOf(v))
	case reflect.Int64:
		if refType.Type.String() == "time.Duration" {
			dValue, err := time.ParseDuration(value)
			if err != nil {
				return err
			}

			field.Set(reflect.ValueOf(dValue))
		} else {
			intValue, err := strconv.ParseInt(value, 10, field.Type().Bits())
			if err != nil {
				return err
			}

			field.SetInt(intValue)
		}
	default:
		return ErrUnsupportedType
	}

	return nil
}

func handleSlice(field reflect.Value, value, sep string) error {
	if sep == "" {
		sep = separator
	}

	splitData := strings.Split(value, sep)

	switch field.Type() {
	case sliceOfStrings:
		field.Set(reflect.ValueOf(splitData))
	case sliceOfInts:
		intData, err := parseInts(splitData)
		if err != nil {
			return err
		}

		field.Set(reflect.ValueOf(intData))
	case sliceOfInt64s:
		int64Data, err := parseInt64s(splitData)
		if err != nil {
			return err
		}

		field.Set(reflect.ValueOf(int64Data))

	case sliceOfFloat32s:
		data, err := parseFloat32s(splitData)
		if err != nil {
			return err
		}

		field.Set(reflect.ValueOf(data))
	case sliceOfFloat64s:
		data, err := parseFloat64s(splitData)
		if err != nil {
			return err
		}

		field.Set(reflect.ValueOf(data))
	case sliceOfBools:
		boolData, err := parseBools(splitData)
		if err != nil {
			return err
		}

		field.Set(reflect.ValueOf(boolData))
	default:
		return ErrUnsupportedSliceType
	}

	return nil
}

func parseInts(data []string) ([]int, error) {
	var intSlice []int

	for _, v := range data {
		intValue, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return nil, err
		}

		intSlice = append(intSlice, int(intValue))
	}

	return intSlice, nil
}

func parseInt64s(data []string) ([]int64, error) {
	var intSlice []int64

	for _, v := range data {
		intValue, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}

		intSlice = append(intSlice, int64(intValue))
	}

	return intSlice, nil
}

func parseFloat32s(data []string) ([]float32, error) {
	var float32Slice []float32

	for _, v := range data {
		obj, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return nil, err
		}

		float32Slice = append(float32Slice, float32(obj))
	}

	return float32Slice, nil
}

func parseFloat64s(data []string) ([]float64, error) {
	var float64Slice []float64

	for _, v := range data {
		obj, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}

		float64Slice = append(float64Slice, float64(obj))
	}

	return float64Slice, nil
}

func parseBools(data []string) ([]bool, error) {
	var boolSlice []bool

	for _, v := range data {
		bVal, err := strconv.ParseBool(v)
		if err != nil {
			return nil, err
		}

		boolSlice = append(boolSlice, bVal)
	}

	return boolSlice, nil
}
