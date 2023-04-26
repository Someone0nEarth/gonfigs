package gonfigs

import (
	"flag"
	"github.com/Someone0nEarth/qogi"
	"os"
	"reflect"
	"testing"

	. "github.com/onsi/gomega"
)

func setup(t *testing.T) (*WithT, func()) {
	g := NewWithT(t)

	oldArgs := os.Args
	oldCommandLine := flag.CommandLine
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	cleanupFunc := func() {
		os.Args = oldArgs
		flag.CommandLine = oldCommandLine
	}

	return g, cleanupFunc
}

func Test_Parse(t *testing.T) {
	g, cleanup := setup(t)
	defer cleanup()

	t.Setenv("TEST_ENV_VALUE_1", "env_value_1")
	t.Setenv("TEST_ENV_VALUE_2", "env_value_2")
	t.Setenv("TEST_ENV_VALUE_3", "env_value_3")
	t.Setenv("TEST_ENV_VALUE_4", "env_value_4")
	t.Setenv("TEST_ENV_VALUE_5", "env_value_5")

	type ExampleConfig struct {
		WithDefault                    *string `env_name:"NO_ENV_VALUE" default:"default_1"`
		WithoutEnvAndDefault           *string `env_name:"NO_ENV_VALUE2"`
		WithEnvAndDefault              *string `env_name:"TEST_ENV_VALUE_1" default:"default_2"`
		WithEnv                        *string `env_name:"TEST_ENV_VALUE_2"`
		WithValueSetAndEnvAndDefault   *string `env_name:"TEST_ENV_VALUE_3" default:"default_2"`
		WithValueSetAndOnlyEnv         *string `env_name:"TEST_ENV_VALUE_4"`
		WithValueSetAndOnlyDefault     *string `env_name:"NO_ENV_VALUE3" default:"default_3"`
		WithValueSetOnly               *string `env_name:"NO_ENV_VALUE4"`
		NoPointer                      string
		NoPointerWithDefault           string `default:"default_4"`
		NoPointerWithEnvNameAndDefault string `env_name:"TEST_ENV_VALUE_5" default:"default_5"`
		NoPointerWithValue             string
		WithNoTags                     *string
		UnsupportedTypeWithNoTags      complex64
	}

	config := ExampleConfig{
		WithValueSetAndEnvAndDefault: qogi.ToPointer("valueSet_1"),
		WithValueSetAndOnlyEnv:       qogi.ToPointer("valueSet_2"),
		WithValueSetAndOnlyDefault:   qogi.ToPointer("valueSet_3"),
		WithValueSetOnly:             qogi.ToPointer("valueSet_4"),
		NoPointerWithValue:           "valueSet_5"}

	Parse(&config)

	expectedConfig := ExampleConfig{
		WithDefault:                    qogi.ToPointer("default_1"),
		WithoutEnvAndDefault:           nil,
		WithEnvAndDefault:              qogi.ToPointer("env_value_1"),
		WithEnv:                        qogi.ToPointer("env_value_2"),
		WithValueSetAndEnvAndDefault:   qogi.ToPointer("valueSet_1"),
		WithValueSetAndOnlyEnv:         qogi.ToPointer("valueSet_2"),
		WithValueSetAndOnlyDefault:     qogi.ToPointer("valueSet_3"),
		WithValueSetOnly:               qogi.ToPointer("valueSet_4"),
		NoPointer:                      "",
		NoPointerWithDefault:           "default_4",
		NoPointerWithEnvNameAndDefault: "env_value_5",
		NoPointerWithValue:             "valueSet_5",
		WithNoTags:                     nil,
		UnsupportedTypeWithNoTags:      0,
	}

	g.Expect(config).Should(BeEquivalentTo(expectedConfig))
}

func Test_SetConfigsTagsValuesForUint(t *testing.T) {
	g, cleanup := setup(t)
	defer cleanup()

	t.Setenv("TEST_ENV_VALUE_1", "1")
	t.Setenv("TEST_ENV_VALUE_2", "2")

	type ExampleConfig struct {
		WithDefault                  *uint   `env_name:"NO_ENV_VALUE" default:"10"`
		StringWithoutEnvAndDefault   *string `env_name:"NO_ENV_VALUE2"`
		WithEnvAndDefault            *uint   `env_name:"TEST_ENV_VALUE_1" default:"20"`
		WithEnv                      *uint   `env_name:"TEST_ENV_VALUE_2"`
		WithValueSetAndEnvAndDefault *uint   `env_name:"TEST_ENV_VALUE_3" default:"30"`
		WithValueSetAndOnlyEnv       *uint   `env_name:"TEST_ENV_VALUE_4"`
		NoPointerWithDefault         uint    `default:"40"`
	}

	config := ExampleConfig{
		WithValueSetAndEnvAndDefault: qogi.ToPointer(uint(100)),
		WithValueSetAndOnlyEnv:       qogi.ToPointer(uint(200)),
	}

	Parse(&config)

	expectedConfig := ExampleConfig{
		WithDefault:                  qogi.ToPointer(uint(10)),
		StringWithoutEnvAndDefault:   nil,
		WithEnvAndDefault:            qogi.ToPointer(uint(1)),
		WithEnv:                      qogi.ToPointer(uint(2)),
		WithValueSetAndEnvAndDefault: qogi.ToPointer(uint(100)),
		WithValueSetAndOnlyEnv:       qogi.ToPointer(uint(200)),
		NoPointerWithDefault:         40,
	}

	g.Expect(config).Should(BeEquivalentTo(expectedConfig))
}

func Test_Parse_Not_a_Struct(t *testing.T) {
	g, cleanup := setup(t)
	defer cleanup()

	notStruct := "I am not a golang struct"
	g.Expect(func() { Parse(notStruct) }).To(Panic())
}

func Test_Parse_UnsupportedTypes(t *testing.T) {
	g, cleanup := setup(t)
	defer cleanup()

	withUnsupportedTypes := struct {
		UnsupportedPointerType *complex64 `default:"12345"`
		UnsupportedType        complex64  `default:"12345"`
	}{}

	g.Expect(func() { Parse(withUnsupportedTypes) }).To(Panic())
}

func Test_Parse_WithoutConfigsTags(t *testing.T) {
	g, cleanup := setup(t)
	defer cleanup()

	type ExampleConfig struct {
		AnyField1 *string
		AnyField2 *string
	}
	structWithoutTags := ExampleConfig{
		AnyField1: qogi.ToPointer("value"),
		AnyField2: nil}

	Parse(&structWithoutTags)

	expectedStruct := ExampleConfig{
		AnyField1: qogi.ToPointer("value"),
		AnyField2: nil}

	g.Expect(structWithoutTags).Should(BeEquivalentTo(expectedStruct))
}

func Test_Parse_GlobalFlags(t *testing.T) {
	g, cleanup := setup(t)
	defer cleanup()

	type ExampleConfig struct {
		AnyField1 *string `arg_name:"arg_1"`
		AnyField2 *string `arg_name:"arg_2" default:"arg_2_default"`
		AnyField3 uint    `arg_name:"arg_3" default:"3"`
		AnyField4 uint    `arg_name:"arg_4"`
	}
	g.Expect(flag.Lookup("arg_1")).Should(BeNil())
	g.Expect(flag.Lookup("arg_2")).Should(BeNil())

	Parse(&ExampleConfig{})

	flagArg1 := flag.Lookup("arg_1")
	g.Expect(flagArg1.Name).Should(BeIdenticalTo("arg_1"))
	g.Expect(reflect.TypeOf(flagArg1.Value).Elem().Kind()).Should(BeIdenticalTo(reflect.String))
	g.Expect(flagArg1.DefValue).Should(BeIdenticalTo(""))

	flagArg2 := flag.Lookup("arg_2")
	g.Expect(flagArg2.Name).Should(BeIdenticalTo("arg_2"))
	g.Expect(reflect.TypeOf(flagArg2.Value).Elem().Kind()).Should(BeIdenticalTo(reflect.String))
	g.Expect(flagArg2.DefValue).Should(BeIdenticalTo("arg_2_default"))

	flagArg3 := flag.Lookup("arg_3")
	g.Expect(flagArg3.Name).Should(BeIdenticalTo("arg_3"))
	g.Expect(reflect.TypeOf(flagArg3.Value).Elem().Kind().String()).Should(BeIdenticalTo(reflect.Uint.String()))
	g.Expect(flagArg3.DefValue).Should(BeIdenticalTo("3"))

	flagArg4 := flag.Lookup("arg_4")
	g.Expect(flagArg4.Name).Should(BeIdenticalTo("arg_4"))
	g.Expect(reflect.TypeOf(flagArg4.Value).Elem().Kind().String()).Should(BeIdenticalTo(reflect.Uint.String()))
	g.Expect(flagArg4.DefValue).Should(BeIdenticalTo("0"))
}

func Test_DescriptionTags(t *testing.T) {
	g, cleanup := setup(t)
	defer cleanup()

	type ExampleConfig struct {
		WithDescription    *uint `arg_name:"arg_name_1" description:"Usage of arg_name_1"`
		WithOutDescription *uint `arg_name:"arg_name_2"`
	}

	config := ExampleConfig{}

	Parse(&config)

	flagArg1 := flag.Lookup("arg_name_1")
	g.Expect(flagArg1.Usage).Should(BeIdenticalTo("Usage of arg_name_1"))

	flagArg2 := flag.Lookup("arg_name_2")
	g.Expect(flagArg2.Usage).Should(BeEmpty())
}

func Test_ArgDescriptionAndEnvTags(t *testing.T) {
	g, cleanup := setup(t)
	defer cleanup()

	type ExampleConfig struct {
		WithDescriptionOnly      *uint `arg_name:"arg_name_1" description:"Usage of arg_name_1."`
		WithDescriptionAndEnv    *uint `env_name:"env_name_2" arg_name:"arg_name_2" description:"Usage of arg_name_2."`
		WithOutDescriptionAndEnv *uint `env_name:"env_name_3" arg_name:"arg_name_3"`
	}

	config := ExampleConfig{}

	Parse(&config)

	flagArg1 := flag.Lookup("arg_name_1")
	g.Expect(flagArg1.Usage).Should(BeIdenticalTo("Usage of arg_name_1."))

	flagArg2 := flag.Lookup("arg_name_2")
	g.Expect(flagArg2.Usage).Should(BeIdenticalTo("Usage of arg_name_2. (Overrides ENV variable 'env_name_2')"))

	flagArg3 := flag.Lookup("arg_name_3")
	g.Expect(flagArg3.Usage).Should(BeIdenticalTo(" (Overrides ENV variable 'env_name_3')"))
}

func Test_lookupValueFromConfigsTags(t *testing.T) {
	g, cleanup := setup(t)
	defer cleanup()

	t.Setenv("TEST_ENV_VALUE_1", "env_value_1")

	fieldTags := reflect.StructTag(``)
	value := lookupValueUsingConfigsTags(fieldTags)
	g.Expect(value).Should(BeNil())

	fieldTags = `any_tag="any_value"`
	value = lookupValueUsingConfigsTags(fieldTags)
	g.Expect(value).Should(BeNil())

	fieldTags = `env_name:"NO_ENV_VALUE"`
	value = lookupValueUsingConfigsTags(fieldTags)
	g.Expect(value).Should(BeNil())

	fieldTags = `env_name:"TEST_ENV_VALUE_1"`
	value = lookupValueUsingConfigsTags(fieldTags)
	g.Expect(*value).Should(BeIdenticalTo("env_value_1"))

	fieldTags = `default:"default_1"`
	value = lookupValueUsingConfigsTags(fieldTags)
	g.Expect(*value).Should(BeIdenticalTo("default_1"))

	fieldTags = `env_name:"NO_ENV_VALUE" default:"default_1"`
	value = lookupValueUsingConfigsTags(fieldTags)
	g.Expect(*value).Should(BeIdenticalTo("default_1"))

	fieldTags = `env_name:"TEST_ENV_VALUE_1" default:"default_1"`
	value = lookupValueUsingConfigsTags(fieldTags)
	g.Expect(*value).Should(BeIdenticalTo("env_value_1"))
}

func Test_lookupValueFromConfigsTags2(t *testing.T) {
	g, cleanup := setup(t)
	defer cleanup()

	os.Args = []string{"command_name", "-arg_name_1", "arg_value_1"}

	fieldTags := reflect.StructTag(`arg_name:"arg_name_1"`)
	value := lookupValueUsingConfigsTags(fieldTags)
	g.Expect(*value).Should(BeEquivalentTo("arg_value_1"))

	os.Args = []string{"command_name", "-arg_name_1", "arg_value_1"}

	//TODO archive? https://stackoverflow.com/questions/68284402/get-rid-of-flag-provided-but-not-defined-when-using-flag-package

	fieldTags = `arg_name:"arg_name_2"`
	value = lookupValueUsingConfigsTags(fieldTags)
	g.Expect(value).Should(BeNil())

	os.Args = []string{"command_name", "--arg_name_1=arg_value_1"}

	fieldTags = `arg_name:"arg_name_1"`
	value = lookupValueUsingConfigsTags(fieldTags)
	g.Expect(*value).Should(BeEquivalentTo("arg_value_1"))

	os.Args = []string{"command_name", "--arg_name_1=arg_value_1"}

	fieldTags = `arg_name:"arg_name_2" default:"default_1"`
	value = lookupValueUsingConfigsTags(fieldTags)
	g.Expect(*value).Should(BeEquivalentTo("default_1"))

	os.Args = []string{"command_name", "--arg_name_1=arg_value_1", "-h"}

	fieldTags = `arg_name:"arg_name_2" default:"default_1"`
	value = lookupValueUsingConfigsTags(fieldTags)
	g.Expect(*value).Should(BeEquivalentTo("default_1"))
}

func Test_WithoutDefault(t *testing.T) {
	g, cleanup := setup(t)
	defer cleanup()

	type ExampleConfig struct {
		Field1 *uint   `arg_name:"arg_name_1"`
		Field2 uint    `arg_name:"arg_name_2"`
		Field3 *string `arg_name:"arg_name_3"`
		Field4 string  `arg_name:"arg_name_4"`
	}

	config := ExampleConfig{}

	Parse(&config)

	expectedConfig := ExampleConfig{}

	g.Expect(config).Should(BeComparableTo(expectedConfig))
}
