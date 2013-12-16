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
	set := make(map[*flag.Flag]bool, 0)
	flag.Visit(func(f *flag.Flag) {
		set[f] = true
	})

	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	flag.VisitAll(func(f *flag.Flag) {
		if !set[f] {
			r := strings.NewReplacer(".", "_", "-", "_")
			name := r.Replace(f.Name)
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

// For each declared flag, Parse() will get the value of the corresponding
// environment variable and will set it. If dots or dash are presents in the
// flag name, they will be converted to underscores. If you want flag names to
// be converted to uppercase, you can set `UseUpperCaseFlagNames` to `true`.
//
// If Parse fails, a fatal error is issued.
func Parse() {
	if err := parse(); err != nil {
		log.Fatalln(err)
	}
}
