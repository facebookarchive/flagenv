package flagenv

import (
	"flag"
	"os"
	"testing"
)

func TestParseEnv(t *testing.T) {
	os.Setenv("foo", "bar")
	os.Setenv("foo_bar", "barfoo")

	var flagFoo = flag.String("foo", "", "")
	var flagFoobar = flag.String("foo_bar", "", "")
	var flagDotSeparator = flag.String("foo.bar", "", "")
	var flagDashSeparator = flag.String("foo-bar", "", "")

	if err := ParseEnv(); err != nil {
		t.Error(err)
	}

	if *flagFoo != "bar" {
		t.Fatalf("want 'bar', have %q", *flagFoo)
	}
	if *flagFoobar != "barfoo" {
		t.Fatalf("want 'barfoo', have %q", *flagFoobar)
	}

	// Testing . separator
	if *flagDotSeparator != "barfoo" {
		t.Fatalf("want 'barfoo', have %q", *flagDotSeparator)
	}

	// Testing - separator
	if *flagDashSeparator != "barfoo" {
		t.Fatalf("want 'barfoo', have %q", *flagDashSeparator)
	}

	os.Setenv("FOOBAR", "bar")
	UseUpperCaseFlagNames = true
	var flagUppercase = flag.String("foobar", "", "")
	if err := ParseEnv(); err != nil {
		t.Error(err)
	}
	if *flagUppercase != "bar" {
		t.Fatalf("want 'bar', have %q", *flagUppercase)
	}
	UseUpperCaseFlagNames = false

	os.Setenv("foo_int", "i should not be a string")
	flag.Int("foo_int", 0, "")
	if err := ParseEnv(); err == nil {
		t.Fatal("expected error parsing non-integer flag, got none")
	}
}
