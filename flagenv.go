// Package flagenv provides the ability to populate flags from
// environment variables.
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

func contains(list []*flag.Flag, f *flag.Flag) bool {
	for _, i := range list {
		if i == f {
			return true
		}
	}
	return false
}

func parse() (err error) {
	var explicit []*flag.Flag
	var all []*flag.Flag
	flag.Visit(func(f *flag.Flag) {
		explicit = append(explicit, f)
	})

	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	flag.VisitAll(func(f *flag.Flag) {
		all = append(all, f)
		if !contains(explicit, f) {
			name := strings.Replace(f.Name, ".", "_", -1)
			name = strings.Replace(name, "-", "_", -1)
			if UseUpperCaseFlagNames {
				name = strings.ToUpper(name)
			}
			val := os.Getenv(name)
			if val != "" {
				err := f.Value.Set(val)
				if err != nil {
					panic(fmt.Errorf("Failed to set flag %s with value %s", f.Name, val))
				}
			}
		}
	})

	return nil
}

// Parse will set each defined flag from its corresponding environment
// variable . If dots or dash are presents in the flag name, they will be
// converted to underscores. If you want flag names to be converted to
// uppercase, you can set `UseUpperCaseFlagNames` to `true`.
//
// If Parse fails, a fatal error is issued.
func Parse() {
	if err := parse(); err != nil {
		log.Fatalln(err)
	}
}
