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

func contains(list []*flag.Flag, f *flag.Flag) bool {
	for _, i := range list {
		if i == f {
			return true
		}
	}
	return false
}

func parse() (err error) {
	explicit := make([]*flag.Flag, 0)
	all := make([]*flag.Flag, 0)
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
			val := os.Getenv(strings.Replace(f.Name, ".", "_", -1))
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

func Parse() {
	if err := parse(); err != nil {
		log.Fatalln(err)
	}
}
