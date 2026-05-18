package web

import (
	"fmt"
	"strings"
)

const preferredUsernameHostPlaceholder = "{{ preferred_username }}"

// applyPreferredUsernameToHost substitutes {{ preferred_username }} using a DNS-safe label.
// The RDP login username may still contain spaces; hostnames and PAA remoteServer may not.
func applyPreferredUsernameToHost(host, username string) string {
	return strings.Replace(host, preferredUsernameHostPlaceholder, sanitizeHostLabel(username), 1)
}

// sanitizeHostLabel makes a string safe for use inside a hostname template (no spaces or @).
func sanitizeHostLabel(username string) string {
	s := strings.TrimSpace(username)
	if s == "" {
		return "user"
	}
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "@", "-")
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '.' {
			b.WriteRune(r)
		}
	}
	out := b.String()
	if out == "" {
		return "user"
	}
	return out
}

func validateHostAddress(host string) error {
	if strings.Contains(host, " ") {
		return fmt.Errorf("host address %q contains spaces (check host templates and URL encoding)", host)
	}
	return nil
}
