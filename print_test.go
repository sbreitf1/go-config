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

type PrintTestPrivate struct {
	Public  string
	private string
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

func TestToStringSimple(t *testing.T) {
	conf := PrintTestSimple{"foobar", 42, true}
	assert.Equal(t, "Stuff.Str:     foobar\nStuff.Number:  42\nStuff.Boolean: true", ToString("Stuff", conf))
}

func TestToStringPrivate(t *testing.T) {
	conf := PrintTestPrivate{"foobar", "unexported"}
	assert.Equal(t, "Stuff.Public: foobar", ToString("Stuff", conf))
}

func TestToStringNested(t *testing.T) {
	conf := PrintTestNested{PrintTestSimple{"foo", 42, true}, PrintTestSimple{"bar", 1337, false}}
	assert.Equal(t, "Pre.Foo.Str:     foo\nPre.Foo.Number:  42\nPre.Foo.Boolean: true\nPre.Bar.Str:     bar\nPre.Bar.Number:  1337\nPre.Bar.Boolean: false", ToString("Pre", conf))
}

func TestToStringAnnoatedSimple(t *testing.T) {
	conf := PrintTestAnnoatedSimple{"foobar", "super-secret-password", PrintTestSimple{"bar", 1337, false}}
	assert.Equal(t, "Stuff.default-val: foobar", ToString("Stuff", conf))
}
