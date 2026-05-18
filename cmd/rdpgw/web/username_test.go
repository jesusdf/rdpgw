package web

import "testing"

func TestSanitizeHostLabel(t *testing.T) {
	cases := map[string]string{
		"John Doe":           "John-Doe",
		"john.doe@corp.com":  "john.doe-corp.com",
		"  spaced  ":         "spaced",
		"":                   "user",
		"already-valid":      "already-valid",
	}
	for in, want := range cases {
		if got := sanitizeHostLabel(in); got != want {
			t.Errorf("sanitizeHostLabel(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestApplyPreferredUsernameToHost(t *testing.T) {
	host := applyPreferredUsernameToHost("vm-{{ preferred_username }}-internal:3389", "John Doe")
	want := "vm-John-Doe-internal:3389"
	if host != want {
		t.Fatalf("got %q want %q", host, want)
	}
	if err := validateHostAddress(host); err != nil {
		t.Fatalf("validateHostAddress: %v", err)
	}
}
