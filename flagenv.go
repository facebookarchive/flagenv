// Package flagenv provides the ability to populate flags from
// environment variables.
package flagenv

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

// Specify a prefix for environment variables.
var Prefix = ""

func contains(list []*flag.Flag, f *flag.Flag) bool {
	for _, i := range list {
		if i == f {
			return true
		}
	}
	return false
}

// ParseSet parses the given flagset. The specified prefix will be applied to
// the environment variable names.
func ParseSet(prefix string, set *flag.FlagSet) error {
	mu.Lock()
	defer mu.Unlock()
	var explicit []*flag.Flag
	var all []*flag.Flag
	set.Visit(func(f *flag.Flag) {
		explicit = append(explicit, f)
	})

	var err error
	set.VisitAll(func(f *flag.Flag) {
		if err != nil {
			return
		}
		all = append(all, f)
		if !contains(explicit, f) {
			defaultName := strings.Replace(f.Name, ".", "_", -1)
			defaultName = strings.Replace(defaultName, "-", "_", -1)
			if prefix != "" {
				defaultName = prefix + defaultName
			}
			defaultName = strings.ToUpper(defaultName)
			envNames := []string{defaultName}
			if setNames, ok := names[f.Name]; ok {
				envNames = []string{}
			set:
				for _, v := range setNames {
					if v == "" {
						envNames = append(envNames, defaultName)
						continue set
					}
					envNames = append(envNames, v)
				}
			}
			for _, v := range envNames {
				val := os.Getenv(v)
				if val != "" {
					if ferr := f.Value.Set(val); ferr != nil {
						err = fmt.Errorf("failed to set flag %q with value %q", f.Name, val)
					}
					return
				}
			}
		}
	})
	return err
}

// Parse will set each defined flag from its corresponding environment
// variable . If dots or dash are presents in the flag name, they will be
// converted to underscores.
//
// If Parse fails, a fatal error is issued.
func Parse() {
	if err := ParseSet(Prefix, flag.CommandLine); err != nil {
		log.Fatalln(err)
	}
}

var (
	names = make(map[string][]string, 0)
	mu    sync.Mutex
)

// SetNames replaces automatically generated environment variables for a
// specific flag and supports multiple names for a flag. Environment variable
// prefixes are ignored for manually specified names. If an emtpy string is
// given as an envName argument it will be replaced with the default generated
// environment variable name.
func SetNames(flagName string, envName ...string) {
	mu.Lock()
	defer mu.Unlock()
	if len(envName) == 0 {
		delete(names, flagName)
		return
	}
	names[flagName] = envName
}
