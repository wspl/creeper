package creeper

import (
	"regexp"
)

func stringsMap(l []string, cb func(string) string) []string {
	nl := []string{}
	for i := range l {
		nl = append(nl, cb(l[i]))
	}
	return nl
}

func stringsFilter(l []string, cb func(string) bool) []string {
	nl := []string{}
	for _, e := range l {
		if cb(e) {
			nl = append(nl, e)
		}
	}
	return nl
}

func stringsMatch(l []string, r *regexp.Regexp) []string {
	return stringsFilter(l, func(s string) bool {
		return r.MatchString(s)
	})
}
