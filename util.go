package creeper

import (
	"regexp"
	"crypto/md5"
	"encoding/hex"
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

func MD5(s string) string {
	m := md5.New()
	m.Write([]byte(s))
	return hex.EncodeToString(m.Sum(nil))
}