// Package flagenv populates flags from environment variables.
package flagenv

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// If set to true, the corresponding environment variable name for a flag will be uppercased.
// ie: For a flag named "foobar", the corresponding environment variable will be "FOOBAR"
var UseUpperCaseFlagNames = false

func parse() (err error) {
	// Record which flags were set by command line args so that we don't overwrite them.
	set := make(map[*flag.Flag]bool, 0)
	flag.Visit(func(f *flag.Flag) {
		set[f] = true
	})

	flag.VisitAll(func(f *flag.Flag) {
		if !set[f] && err == nil {
			r := strings.NewReplacer(".", "_", "-", "_")
			name := r.Replace(f.Name)
			if UseUpperCaseFlagNames {
				name = strings.ToUpper(name)
			}
			val := os.Getenv(name)
			if val != "" {
				if seterr := f.Value.Set(val); seterr != nil {
					err = fmt.Errorf("Failed to set flag %s with value %s", f.Name, val)
				}
			}
		}
	})

	return
}

// For each declared flag, Parse() will get the value of the corresponding
// environment variable and will set it. If dots or dashes are present in the
// flag name, they will be converted to underscores. If you want flag names to
// be converted to uppercase, you can set `UseUpperCaseFlagNames` to `true`.
//
// If Parse fails, a fatal error is issued.
func Parse() {
	if err := parse(); err != nil {
		log.Fatalln(err)
	}
}

// ParseEnv populates flag values from environment variables. See Parse for
// details. Unlike Parse, errors are returned rather than exiting the process.
func ParseEnv() error {
	return parse()
}
