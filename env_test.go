package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type EnvTestSimple struct {
	StringData string
	IntData    int
	BoolDataT  bool
	BoolDataF  bool
}

type EnvTestNested struct {
	OuterString string
	NestedValue EnvTestSimple
}

type EnvTestName struct {
	StringData string        `config:"env:str"`
	IntData    int           `config:"env:number"`
	Nested     EnvTestSimple `config:"env:subval"`
}

type EnvTestSlice struct {
	EmptyList  []int
	List       []string
	NestedList []EnvTestSimple
}

type EnvTestArray struct {
	List [3]string
}

func TestEnvEmpty(t *testing.T) {
	withMockEnv(func(env map[string]string) {
		var conf EnvTestSimple
		assert.NoError(t, FromEnvironment("test", &conf))
	})
}

func TestEnvSimple(t *testing.T) {
	withMockEnv(func(env map[string]string) {
		env["TEST_STRINGDATA"] = "foobar"
		env["TEST_INTDATA"] = "42"
		env["TEST_BOOLDATAT"] = "true"
		env["TEST_BOOLDATAF"] = "false"

		var conf EnvTestSimple
		if assert.NoError(t, FromEnvironment("test", &conf)) {
			assert.Equal(t, "foobar", conf.StringData)
			assert.Equal(t, 42, conf.IntData)
			assert.True(t, conf.BoolDataT)
			assert.False(t, conf.BoolDataF)
		}
	})
}

func TestEnvInvalidInt(t *testing.T) {
	withMockEnv(func(env map[string]string) {
		env["TEST_INTDATA"] = "not numeric"

		var conf EnvTestSimple
		assert.Error(t, FromEnvironment("test", &conf))
	})
}

func TestEnvInvalidBool(t *testing.T) {
	withMockEnv(func(env map[string]string) {
		env["TEST_BOOLDATAT"] = "not a bool"

		var conf EnvTestSimple
		assert.Error(t, FromEnvironment("test", &conf))
	})
}

func TestEnvNested(t *testing.T) {
	withMockEnv(func(env map[string]string) {
		env["TEST_OUTERSTRING"] = "test"
		env["TEST_NESTEDVALUE_STRINGDATA"] = "foobar"
		env["TEST_NESTEDVALUE_INTDATA"] = "42"
		env["TEST_NESTEDVALUE_BOOLDATAT"] = "true"
		env["TEST_NESTEDVALUE_BOOLDATAF"] = "false"

		var conf EnvTestNested
		if assert.NoError(t, FromEnvironment("test", &conf)) {
			assert.Equal(t, "test", conf.OuterString)
			assert.Equal(t, "foobar", conf.NestedValue.StringData)
			assert.Equal(t, 42, conf.NestedValue.IntData)
			assert.True(t, conf.NestedValue.BoolDataT)
			assert.False(t, conf.NestedValue.BoolDataF)
		}
	})
}

func TestEnvName(t *testing.T) {
	withMockEnv(func(env map[string]string) {
		env["TEST_STR"] = "test str"
		env["TEST_NUMBER"] = "1337"
		env["TEST_SUBVAL_STRINGDATA"] = "foobar"
		env["TEST_SUBVAL_INTDATA"] = "42"
		env["TEST_SUBVAL_BOOLDATAT"] = "true"
		env["TEST_SUBVAL_BOOLDATAF"] = "false"

		var conf EnvTestName
		if assert.NoError(t, FromEnvironment("test", &conf)) {
			assert.Equal(t, "test str", conf.StringData)
			assert.Equal(t, 1337, conf.IntData)
			assert.Equal(t, "foobar", conf.Nested.StringData)
			assert.Equal(t, 42, conf.Nested.IntData)
			assert.True(t, conf.Nested.BoolDataT)
			assert.False(t, conf.Nested.BoolDataF)
		}
	})
}

func TestEnvSlice(t *testing.T) {
	withMockEnv(func(env map[string]string) {
		env["TEST_EMPTYLIST_NUM"] = "0"
		env["TEST_LIST_NUM"] = "3"
		env["TEST_LIST_0"] = "foo"
		env["TEST_LIST_1"] = "bar"
		env["TEST_LIST_2"] = "42"
		env["TEST_NESTEDLIST_NUM"] = "1"
		env["TEST_NESTEDLIST_0_STRINGDATA"] = "foobar"
		env["TEST_NESTEDLIST_0_INTDATA"] = "42"
		env["TEST_NESTEDLIST_0_BOOLDATAT"] = "true"
		env["TEST_NESTEDLIST_0_BOOLDATAF"] = "false"

		var conf EnvTestSlice
		if assert.NoError(t, FromEnvironment("test", &conf)) {
			assert.Equal(t, []int{}, conf.EmptyList)
			assert.Equal(t, []string{"foo", "bar", "42"}, conf.List)
			assert.Equal(t, []EnvTestSimple{EnvTestSimple{"foobar", 42, true, false}}, conf.NestedList)
		}
	})
}

func TestEnvArray(t *testing.T) {
	withMockEnv(func(env map[string]string) {
		env["TEST_LIST_0"] = "foo"
		env["TEST_LIST_1"] = "bar"
		env["TEST_LIST_2"] = "42"

		var conf EnvTestArray
		if assert.NoError(t, FromEnvironment("test", &conf)) {
			assert.Equal(t, [3]string{"foo", "bar", "42"}, conf.List)
		}
	})
}

func withMockEnv(f func(map[string]string)) {
	oldLookupEnv := lookupEnv
	defer func() { lookupEnv = oldLookupEnv }()

	env := make(map[string]string)
	lookupEnv = func(str string) (string, bool) {
		val, ok := env[str]
		return val, ok
	}

	f(env)
}
