package gonfigs

import (
	"flag"
	"fmt"

	"github.com/Someone0nEarth/qogi"

	"io"
	"os"
	"reflect"
)

// Parse parses a *struct for "gonfigs" tags like `argName`, `envName`, `default` and using them to set the value
// of the corresponding fields.
//
// If more than one gonfigs tag are matching for a field, the priority for setting the value is:
// `argName` > `envName` > `default`
//
// Tags:
//   - `argName` looks for a flag.Flag (command line argument) with the name corresponding to the tag value.
//     Will be listed, when using '--help'
//   - `envName` looks for an environment variable with the name corresponding to the tag value
//   - `default` will use the tags value for the field
//   - `description` the description of the field. For documentation and/or if argName is although set, it will be
//     shown in the `usage` field when using `--help` in the command line.
//
// Struct fields will be set only, if their current values are zero (Value.IsZero()==true).
//
// At the moment, the following struct fields types are supported: unit, *uint, string and *string
//
// The method will panic, if `config` is not a struct, `config` is not a pointer and the tags are used on unsupported
// field types.
func Parse(config any) {
	if reflect.TypeOf(config).Kind() != reflect.Pointer {
		panic("Config struct is not a pointer.")
	}

	if reflect.TypeOf(config).Elem().Kind() != reflect.Struct {
		panic("Config is not a struct.")
	}

	for index := 0; index < reflect.TypeOf(config).Elem().NumField(); index++ {
		field := reflect.Indirect(reflect.ValueOf(config)).Field(index)
		fieldTags := reflect.TypeOf(config).Elem().Field(index).Tag

		existingValue := reflect.ValueOf(config).Elem().Field(index)

		if existingValue.IsZero() {
			value := lookupValueUsingConfigsTags(fieldTags)

			if value != nil {
				setValue(field, value)
			}
		}

		if argumentName, found := fieldTags.Lookup(`argName`); found {
			description, _ := fieldTags.Lookup("description")

			if envName, envNameFound := fieldTags.Lookup("envName"); envNameFound {
				description = fmt.Sprintf(description+" (Overrides ENV variable '%s')", envName)
			}

			if defaultValue, foundDefault := fieldTags.Lookup("default"); foundDefault {
				addArgumentToGlobalFlag(argumentName, &defaultValue, field, description)
			} else {
				addArgumentToGlobalFlag(argumentName, nil, field, description)
			}
		}
	}
}

func setValue(field reflect.Value, value *string) {
	if field.Kind() == reflect.Pointer {
		fieldValueType := field.Type().Elem()

		valueToSet := createValueForField(fieldValueType, value)

		field.Set(reflect.New(fieldValueType))
		field.Elem().Set(valueToSet)
	} else {
		valueToSet := createValueForField(field.Type(), value)

		field.Set(valueToSet)
	}
}

func createValueForField(fieldType reflect.Type, value *string) reflect.Value {
	if fieldType.Kind() == reflect.String {
		if value == nil {
			return reflect.ValueOf("")
		}
		return reflect.ValueOf(value).Elem()
	}

	if fieldType.Kind() == reflect.Uint {
		if value == nil {
			return reflect.ValueOf(uint(0))
		}

		uintValue, err := qogi.AtoUi(*value)
		if err != nil {
			panic(err)
		}
		return reflect.ValueOf(&uintValue).Elem()
	}

	panic("Unsupported kind and/or type: " + fieldType.Kind().String() + " / " + fieldType.String())
}

func lookupValueUsingConfigsTags(fieldTags reflect.StructTag) *string {
	if argumentName, found := fieldTags.Lookup(`argName`); found {
		if argumentValue, valueFound := lookupArgument(argumentName); valueFound {
			return &argumentValue
		}
	}

	if envName, found := fieldTags.Lookup(`envName`); found {
		if envValue, valueFound := os.LookupEnv(envName); valueFound {
			return &envValue
		}
	}

	if defaultValue, valueFound := fieldTags.Lookup("default"); valueFound {
		return &defaultValue
	}

	return nil
}

func addArgumentToGlobalFlag(argumentName string, defaultValue *string, field reflect.Value, description string) {
	if flag.Lookup(argumentName) == nil { //Prevent "flag redefined" panic
		var defaultFlagValue reflect.Value

		if field.Kind() == reflect.Pointer {
			defaultFlagValue = createValueForField(field.Type().Elem(), defaultValue)
		} else {
			defaultFlagValue = createValueForField(field.Type(), defaultValue)
		}

		if defaultFlagValue.Kind() == reflect.String {
			flag.String(argumentName, defaultFlagValue.String(), description)
		} else if defaultFlagValue.Kind() == reflect.Uint {
			flag.Uint(argumentName, uint(defaultFlagValue.Uint()), description)
		} else {
			panic("Unsupported kind and/or type: " + defaultFlagValue.Kind().String() + " / " + defaultFlagValue.Type().String())
		}
	}
}

func lookupArgument(argumentName string) (string, bool) {
	var argumentValue string

	flagSet := flag.NewFlagSet("", flag.ContinueOnError)
	flagSet.StringVar(&argumentValue, argumentName, "", "")

	flagSet.SetOutput(io.Discard)
	_ = flagSet.Parse(os.Args[1:])

	if argumentValue != "" {
		return argumentValue, true
	}

	return "", false
}
