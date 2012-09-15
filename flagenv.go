// Package flagenv provides the ability to populate flags from
// environment variables.
package flagenv

import (
	"flag"
	"log"
	"os"
	"strings"
)

func contains(list []*flag.Flag, f *flag.Flag) bool {
	for _, i := range list {
		if i == f {
			return true
		}
	}
	return false
}

func Parse() {
	explicit := make([]*flag.Flag, 0)
	all := make([]*flag.Flag, 0)
	flag.Visit(func(f *flag.Flag) {
		explicit = append(explicit, f)
	})
	flag.VisitAll(func(f *flag.Flag) {
		all = append(all, f)
		if !contains(explicit, f) {
			val := os.Getenv(strings.Replace(f.Name, ".", "_", -1))
			if val != "" {
				err := f.Value.Set(val)
				if err != nil {
					log.Fatalf("Failed to set flag %s with value %s", f.Name, val)
				}
			}
		}
	})
}
