package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type PrintTestSimple struct {
	Str     string
	Number  int
	Boolean bool
}

type PrintTestNested struct {
	Foo PrintTestSimple
	Bar PrintTestSimple
}

type PrintTestAnnoatedSimple struct {
	Default      string          `config:"print:default-val"`
	Secret       string          `config:"print:-"`
	NestedSecret PrintTestSimple `config:"print:-"`
}

func TestSprintSimple(t *testing.T) {
	conf := PrintTestSimple{"foobar", 42, true}
	str, err := Sprint("Stuff", conf)
	if assert.NoError(t, err) {
		assert.Equal(t, "Stuff.Str:     foobar\nStuff.Number:  42\nStuff.Boolean: true", str)
	}
}

func TestSprintNested(t *testing.T) {
	conf := PrintTestNested{PrintTestSimple{"foo", 42, true}, PrintTestSimple{"bar", 1337, false}}
	str, err := Sprint("Pre", conf)
	if assert.NoError(t, err) {
		assert.Equal(t, "Pre.Foo.Str:     foo\nPre.Foo.Number:  42\nPre.Foo.Boolean: true\nPre.Bar.Str:     bar\nPre.Bar.Number:  1337\nPre.Bar.Boolean: false", str)
	}
}

func TestSprintAnnoatedSimple(t *testing.T) {
	conf := PrintTestAnnoatedSimple{"foobar", "super-secret-password", PrintTestSimple{"bar", 1337, false}}
	str, err := Sprint("Stuff", conf)
	if assert.NoError(t, err) {
		assert.Equal(t, "Stuff.default-val: foobar", str)
	}
}