package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type PrintTestSimple struct {
	Str     string
	Number  int
	Boolean bool
}

func TestToStringSimple(t *testing.T) {
	conf := PrintTestSimple{"foobar", 42, true}
	assert.Equal(t, "Stuff.Str:     foobar\nStuff.Number:  42\nStuff.Boolean: true", ToString("Stuff", conf))
}

func TestToStringSimplePtr(t *testing.T) {
	conf := PrintTestSimple{"foobar", 42, true}
	assert.Equal(t, "Stuff.Str:     foobar\nStuff.Number:  42\nStuff.Boolean: true", ToString("Stuff", &conf))
}

type PrintTestPrivate struct {
	Public  string
	private string
}

func TestToStringPrivate(t *testing.T) {
	conf := PrintTestPrivate{"foobar", "unexported"}
	assert.Equal(t, "Stuff.Public: foobar", ToString("Stuff", conf))
}

type PrintTestPointer struct {
	Int *int
}

func TestToStringPointer(t *testing.T) {
	val := 42
	conf := PrintTestPointer{&val}
	assert.Equal(t, "Stuff.Int: 42", ToString("Stuff", conf))
}

func TestToStringNilPointer(t *testing.T) {
	conf := PrintTestPointer{}
	assert.Equal(t, "", ToString("Stuff", conf))
}

type PrintTestNested struct {
	Foo PrintTestSimple
	Bar PrintTestSimple
}

func TestToStringNested(t *testing.T) {
	conf := PrintTestNested{PrintTestSimple{"foo", 42, true}, PrintTestSimple{"bar", 1337, false}}
	assert.Equal(t, "Pre.Foo.Str:     foo\nPre.Foo.Number:  42\nPre.Foo.Boolean: true\nPre.Bar.Str:     bar\nPre.Bar.Number:  1337\nPre.Bar.Boolean: false", ToString("Pre", conf))
}

type PrintTestAnnoatedSimple struct {
	Default      string          `config:"print:default-val"`
	Secret       string          `config:"print:-"`
	NestedSecret PrintTestSimple `config:"print:-"`
}

func TestToStringAnnoatedSimple(t *testing.T) {
	conf := PrintTestAnnoatedSimple{"foobar", "super-secret-password", PrintTestSimple{"bar", 1337, false}}
	assert.Equal(t, "Stuff.default-val: foobar", ToString("Stuff", conf))
}

type PrintTestSlice struct {
	List []string
}

func TestToStringSlice(t *testing.T) {
	conf := PrintTestSlice{[]string{"foo", "bar"}}
	assert.Equal(t, "Stuff.List[0]: foo\nStuff.List[1]: bar", ToString("Stuff", conf))
}

type PrintTestArray struct {
	List [2]string
}

func TestToStringArray(t *testing.T) {
	conf := PrintTestArray{[2]string{"foo", "bar"}}
	assert.Equal(t, "Stuff.List[0]: foo\nStuff.List[1]: bar", ToString("Stuff", conf))
}

type PrintTestNonZero struct {
	Str     string `config:"print:[nonzero]"`
	StrZero string `config:"print:[nonzero]"`
	Int     int    `config:"print:[nonzero]"`
	IntZero int    `config:"print:[nonzero]"`
}

func TestToStringNonZero(t *testing.T) {
	conf := PrintTestNonZero{"foobar", "", 42, 0}
	assert.Equal(t, "Stuff.Str: foobar\nStuff.Int: 42", ToString("Stuff", conf))
}

type PrintTestLen struct {
	Str   string `config:"print:[len]"`
	Slice []int  `config:"print:[len]"`
	Array [5]int `config:"print:[len]"`
}

func TestToStringLen(t *testing.T) {
	conf := PrintTestLen{"foobar", []int{42, 1337}, [5]int{1, 2, 3, 4, 5}}
	assert.Equal(t, "Stuff.Str:   6\nStuff.Slice: 2\nStuff.Array: 5", ToString("Stuff", conf))
}

type PrintTestMasked struct {
	ZeroStr string `config:"print:[mask]"`
	Str     string `config:"print:[mask]"`
	ZeroInt int    `config:"print:[mask]"`
	Int     int    `config:"print:[mask]"`
}

func TestToStringMasked(t *testing.T) {
	conf := PrintTestMasked{"", "some test string", 0, 42}
	assert.Equal(t, "Stuff.Str: ******\nStuff.Int: ******", ToString("Stuff", conf))
}

type PrintTestTime struct {
	Time     time.Time
	Duration time.Duration
}

func TestToStringTime(t *testing.T) {
	conf := PrintTestTime{
		time.Date(2020, time.May, 5, 13, 49, 44, 0, time.UTC),
		2*time.Hour + 15*time.Minute,
	}
	assert.Equal(t, "Stuff.Time:     2020-05-05 13:49:44 +0000 UTC\nStuff.Duration: 2h15m0s", ToString("Stuff", conf))
}
