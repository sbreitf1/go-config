package config

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type JSONTestSimple struct {
	StringData string
	IntData    int
	BoolDataT  bool
	BoolDataF  bool
}

func TestFromFileSimple(t *testing.T) {
	tmp, err := ioutil.TempFile("", "go-config-test-")
	if assert.NoError(t, err) {
		tmpFile := tmp.Name()
		tmp.Close()
		defer os.Remove(tmpFile)

		if assert.NoError(t, ioutil.WriteFile(tmpFile, []byte(`{"StringData":"foobar","IntData":42,"BoolDataT":true,"BoolDataF":false}`), os.ModePerm)) {
			var conf EnvTestSimple
			if assert.NoError(t, FromFile(tmpFile, &conf)) {
				assert.Equal(t, "foobar", conf.StringData)
				assert.Equal(t, 42, conf.IntData)
				assert.True(t, conf.BoolDataT)
				assert.False(t, conf.BoolDataF)
			}
		}
	}
}

func TestFromJSONSimple(t *testing.T) {
	var conf EnvTestSimple
	if assert.NoError(t, FromJSON([]byte(`{"StringData":"foobar","IntData":42,"BoolDataT":true,"BoolDataF":false}`), &conf)) {
		assert.Equal(t, "foobar", conf.StringData)
		assert.Equal(t, 42, conf.IntData)
		assert.True(t, conf.BoolDataT)
		assert.False(t, conf.BoolDataF)
	}
}

func TestFromJSONDate(t *testing.T) {
	var conf time.Time
	if assert.NoError(t, FromJSON([]byte(`"2020-02-25T17:20:34Z"`), &conf)) {
		assert.Equal(t, time.Date(2020, time.February, 25, 17, 20, 34, 0, time.UTC), conf)
	}
}

func TestFromJSONDuration(t *testing.T) {
	var conf time.Duration
	if assert.NoError(t, FromJSON([]byte(`"1h34m17s"`), &conf)) {
		assert.Equal(t, 1*time.Hour+34*time.Minute+17*time.Second, conf)
	}
}

type JSONTestTag struct {
	StringData string `json:"str"`
	IntData    int    `json:"-"`
	BoolDataT  bool   `json:"yup"`
	BoolDataF  bool
}

func TestFromJSONTag(t *testing.T) {
	var conf JSONTestTag
	if assert.NoError(t, FromJSON([]byte(`{"str":"foobar","IntData":42,"yup":true,"BoolDataF":false}`), &conf)) {
		assert.Equal(t, "foobar", conf.StringData)
		assert.Equal(t, 0, conf.IntData)
		assert.True(t, conf.BoolDataT)
		assert.False(t, conf.BoolDataF)
	}
}
