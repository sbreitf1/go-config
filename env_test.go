package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type EnvTestSimple struct {
	StringData string
	IntData    int
	BoolDataT  bool
	BoolDataF  bool
}

func TestEnvNoPointer(t *testing.T) {
	withMockEnv(func(env map[string]string) {
		env["TEST_STRINGDATA"] = "foobar"
		env["TEST_INTDATA"] = "42"
		env["TEST_BOOLDATAT"] = "true"
		env["TEST_BOOLDATAF"] = "false"

		var conf EnvTestSimple
		if assert.Error(t, FromEnvironment("test", conf)) {
			assert.Equal(t, "", conf.StringData)
			assert.Equal(t, 0, conf.IntData)
			assert.False(t, conf.BoolDataT)
			assert.False(t, conf.BoolDataF)
		}
	})
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

type EnvTestPrivate struct {
	Public  string
	private string
}

func TestEnvPrivate(t *testing.T) {
	withMockEnv(func(env map[string]string) {
		env["TEST_PUBLIC"] = "foobar"
		env["TEST_PRIVATE"] = "should not be used"

		conf := EnvTestPrivate{"", "keep it"}
		assert.NoError(t, FromEnvironment("test", &conf))
		assert.Equal(t, "foobar", conf.Public)
		assert.Equal(t, "keep it", conf.private)
	})
}

type EnvTestDefault struct {
	StringData string `config:"default:foobar"`
	IntData    int    `config:"default:42"`
	BoolDataT  bool   `config:"default:true"`
	BoolDataF  bool   `config:"default:false"`
}

func TestEnvDefault(t *testing.T) {
	withMockEnv(func(env map[string]string) {
		var conf EnvTestDefault
		if assert.NoError(t, FromEnvironment("test", &conf)) {
			assert.Equal(t, "foobar", conf.StringData)
			assert.Equal(t, 42, conf.IntData)
			assert.True(t, conf.BoolDataT)
			assert.False(t, conf.BoolDataF)
		}
	})
}

func TestEnvDefaultOverride(t *testing.T) {
	withMockEnv(func(env map[string]string) {
		env["TEST_STRINGDATA"] = "NEW STR"
		env["TEST_INTDATA"] = "1337"
		env["TEST_BOOLDATAT"] = "false"
		env["TEST_BOOLDATAF"] = "true"

		var conf EnvTestDefault
		if assert.NoError(t, FromEnvironment("test", &conf)) {
			assert.Equal(t, "NEW STR", conf.StringData)
			assert.Equal(t, 1337, conf.IntData)
			assert.False(t, conf.BoolDataT)
			assert.True(t, conf.BoolDataF)
		}
	})
}

type EnvTestNested struct {
	OuterString string
	NestedValue EnvTestSimple
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

type EnvTestName struct {
	StringData string        `config:"env:str"`
	IntData    int           `config:"env:number"`
	Nested     EnvTestSimple `config:"env:subval"`
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

type EnvTestSlice struct {
	EmptyList  []int
	List       []string
	NestedList []EnvTestSimple
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

type EnvTestArray struct {
	List [3]string
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

type EnvTestDateTime struct {
	Val                  time.Time
	Default              time.Time `config:"default:2020-02-17 09:06:21"`
	DefaultDay           time.Time `config:"default:2020-02-17"`
	DefaultUTC           time.Time `config:"default:2020-02-17T09:06:21Z"`
	DefaultShortTimeZone time.Time `config:"default:2020-02-17 07:06:21+0100"`
	DefaultTimeZone      time.Time `config:"default:2020-02-17 09:06:21+01:00"`
}

func TestEnvDateTime(t *testing.T) {
	withMockEnv(func(env map[string]string) {
		env["TEST_VAL"] = "2020-02-17T09:08:42"

		var conf EnvTestDateTime
		if assert.NoError(t, FromEnvironment("test", &conf)) {
			assert.Equal(t, time.Date(2020, time.February, 17, 9, 8, 42, 0, time.Local), conf.Val, "conf.Val")
			assert.Equal(t, time.Date(2020, time.February, 17, 9, 6, 21, 0, time.Local), conf.Default, "conf.Default")
			assert.Equal(t, time.Date(2020, time.February, 17, 0, 0, 0, 0, time.Local), conf.DefaultDay, "conf.DefaultDay")
			assert.Equal(t, time.Date(2020, time.February, 17, 9, 6, 21, 0, time.UTC), conf.DefaultUTC, "conf.DefaultUTC")

			// check hour in given time zone utc+1 (7)
			assert.Equal(t, 7, conf.DefaultShortTimeZone.Hour())
			// in UTC it must be 6
			assert.Equal(t, time.Date(2020, time.February, 17, 6, 6, 21, 0, time.UTC), conf.DefaultShortTimeZone.UTC(), "conf.DefaultShortTimeZone")

			// check hour in given time zone utc+1 (9)
			assert.Equal(t, 9, conf.DefaultTimeZone.Hour())
			// in UTC it must be 8
			assert.Equal(t, time.Date(2020, time.February, 17, 8, 6, 21, 0, time.UTC), conf.DefaultTimeZone.UTC(), "conf.DefaultTimeZone")
		}
	})
}

type EnvTestDuration struct {
	Val           time.Duration
	Default       time.Duration `config:"default:P1Y2M3DT4H5M6S"`
	DefaultYear   time.Duration `config:"default:P1Y"`
	DefaultMonth  time.Duration `config:"default:P1M"`
	DefaultWeek   time.Duration `config:"default:P1W"`
	DefaultDay    time.Duration `config:"default:P1D"`
	DefaultHour   time.Duration `config:"default:T1H"`
	DefaultMinute time.Duration `config:"default:T1M"`
	DefaultSecond time.Duration `config:"default:T1S"`
}

func TestEnvDuration(t *testing.T) {
	second := time.Duration(time.Second)
	minute := time.Duration(time.Minute)
	hour := time.Duration(time.Hour)
	day := 24 * hour
	week := 7 * day
	month := 30 * day
	year := 365 * day

	withMockEnv(func(env map[string]string) {
		env["TEST_VAL"] = "1y 2d 3h 4m 45s"

		var conf EnvTestDuration
		if assert.NoError(t, FromEnvironment("test", &conf)) {
			assert.Equal(t, 1*year+2*day+3*hour+4*minute+45*second, conf.Val, "conf.Val")
			assert.Equal(t, 1*year+2*month+3*day+4*hour+5*minute+6*second, conf.Default, "conf.Default")
			assert.Equal(t, 1*year, conf.DefaultYear, "conf.DefaultYear")
			assert.Equal(t, 1*month, conf.DefaultMonth, "conf.DefaultMonth")
			assert.Equal(t, 1*week, conf.DefaultWeek, "conf.DefaultWeek")
			assert.Equal(t, 1*day, conf.DefaultDay, "conf.DefaultDay")
			assert.Equal(t, 1*hour, conf.DefaultHour, "conf.DefaultHour")
			assert.Equal(t, 1*minute, conf.DefaultMinute, "conf.DefaultMinute")
			assert.Equal(t, 1*second, conf.DefaultSecond, "conf.DefaultSecond")
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
