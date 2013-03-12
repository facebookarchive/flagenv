package flagenv

import (
	"flag"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	os.Setenv("foo", "bar")
	os.Setenv("foo_bar", "barfoo")

	var flagFoo = flag.String("foo", "", "")
	var flagFoobar = flag.String("foo_bar", "", "")
	var flagDotSeparator = flag.String("foo.bar", "", "")

	if err := parse(); err != nil {
		t.Log("blaa")
		t.Error(err)
	}

	if *flagFoo != "bar" {
		t.Fail()
	}
	if *flagFoobar != "barfoo" {
		t.Fail()
	}

	// Testing . separator
	if *flagDotSeparator != "barfoo" {
		t.Fail()
	}

	os.Setenv("foo_int", "i should not be a string")
	flag.Int("foo_int", 0, "")
	if err := parse(); err == nil {
		t.Fail()
	}
}
