package config

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

type JSONTestName struct {
	StringData string `config:"name:str"`
	IntData    int    `config:"name:number"`
	BoolDataT  bool   `config:"name:yup"`
	BoolDataF  bool
}

func TestFromJSONName(t *testing.T) {
	var conf JSONTestName
	if assert.NoError(t, FromJSON([]byte(`{"str":"foobar","IntData":42,"yup":true,"BoolDataF":false}`), &conf)) {
		assert.Equal(t, "foobar", conf.StringData)
		assert.Equal(t, 0, conf.IntData)
		assert.True(t, conf.BoolDataT)
		assert.False(t, conf.BoolDataF)
	}
}

type JSONTestSlice struct {
	SliceData1 []string
	SliceData2 []string
	SliceData3 []string
	SliceData4 []string
	SliceData5 []string
}

func TestFromJSONSlice(t *testing.T) {
	var conf JSONTestSlice
	require.NoError(t, FromJSON([]byte(`{"SliceData2":null,"SliceData3":[],"SliceData4":["foo"],"SliceData5":["bar","test","42"]}`), &conf))
	require.Nil(t, conf.SliceData1)
	require.Nil(t, conf.SliceData2)
	require.Equal(t, []string{}, conf.SliceData3)
	require.Equal(t, []string{"foo"}, conf.SliceData4)
	require.Equal(t, []string{"bar", "test", "42"}, conf.SliceData5)
}

type JSONTestArray struct {
	ArrayData1 [0]string
	ArrayData2 [1]string
	ArrayData3 [2]string
}

func TestFromJSONArray(t *testing.T) {
	var conf JSONTestArray
	require.NoError(t, FromJSON([]byte(`{"ArrayData1":[],"ArrayData2":["foo"],"ArrayData3":["bar","42"]}`), &conf))
	require.Equal(t, [0]string{}, conf.ArrayData1)
	require.Equal(t, [1]string{"foo"}, conf.ArrayData2)
	require.Equal(t, [2]string{"bar", "42"}, conf.ArrayData3)
}
