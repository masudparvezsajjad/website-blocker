package util

import (
	"regexp"
	"strings"
)

var domainRegex = regexp.MustCompile(`^(?i)[a-z0-9][-a-z0-9.]*[a-z0-9]\.[a-z]{2,}$`)

func NormalizeDomain(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	s = strings.TrimPrefix(s, "http://")
	s = strings.TrimPrefix(s, "https://")

	// remove path if present
	if idx := strings.Index(s, "/"); idx >= 0 {
		s = s[:idx]
	}

	// remove port if present
	if idx := strings.Index(s, ":"); idx >= 0 {
		s = s[:idx]
	}

	return s
}

func IsValidDomain(s string) bool {
	if strings.ContainsAny(s, " \t\r\n\"'`;|&$()<>") {
		return false
	}
	return domainRegex.MatchString(s)
}
