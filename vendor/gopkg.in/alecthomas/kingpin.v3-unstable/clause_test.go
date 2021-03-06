package kingpin

import (
	"io/ioutil"
	"net/url"
	"os"

	"github.com/stretchr/testify/assert"

	"testing"
)

func TestParseStrings(t *testing.T) {
	p := Clause{}
	v := p.Strings()
	p.value.Set("a")
	p.value.Set("b")
	assert.Equal(t, []string{"a", "b"}, *v)
}

func TestStringsStringer(t *testing.T) {
	target := []string{}
	v := newAccumulator(&target, nil, func(v interface{}) Value { return newStringValue(v.(*string)) })
	v.Set("hello")
	v.Set("world")
	assert.Equal(t, "hello,world", v.String())
}

func TestParseStringMap(t *testing.T) {
	p := Clause{}
	v := p.StringMap()
	p.value.Set("a:b")
	p.value.Set("b:c")
	assert.Equal(t, map[string]string{"a": "b", "b": "c"}, *v)
}

func TestParseStringMapWithSeparator(t *testing.T) {
	p := Clause{}
	v := p.StringMap(Separator(";"))
	p.value.Set("a:b")
	p.value.Set("b:c;c:d")
	assert.Equal(t, map[string]string{"a": "b", "b": "c", "c": "d"}, *v)
}

func TestParseURL(t *testing.T) {
	p := Clause{}
	v := p.URL()
	p.value.Set("http://w3.org")
	u, err := url.Parse("http://w3.org")
	assert.NoError(t, err)
	assert.Equal(t, *u, **v)
}

func TestParseExistingFile(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	p := Clause{}
	v := p.ExistingFile()
	err = p.value.Set(f.Name())
	assert.NoError(t, err)
	assert.Equal(t, f.Name(), *v)
	err = p.value.Set("/etc/hostsDEFINITELYMISSING")
	assert.Error(t, err)
}

func TestFloat32(t *testing.T) {
	p := Clause{}
	v := p.Float32()
	err := p.value.Set("123.45")
	assert.NoError(t, err)
	assert.InEpsilon(t, 123.45, *v, 0.001)
}

func TestUnicodeShortFlag(t *testing.T) {
	app := newTestApp()
	f := app.Flag("long", "").Short('ä').Bool()
	_, err := app.Parse([]string{"-ä"})
	assert.NoError(t, err)
	assert.True(t, *f)
}
