package security

import "testing"

func TestPickUsername(t *testing.T) {
	const guid = "f47ac10b-58cc-4372-a567-0e02b2c3d479"

	cases := []struct {
		name                       string
		preferred, unique, upn, id string
		want                       string
	}{
		{"preferred wins", "jdieguez", "ignored", "ignored", guid, "jdieguez"},
		{"preferred with two names", "Jesús Diéguez", "", "", guid, "Jesús Diéguez"},
		{"falls back to unique_name", "", "DOMINIO\\jdieguez", "", guid, "DOMINIO\\jdieguez"},
		{"falls back to upn", "", "", "user@example.com", guid, "user@example.com"},
		{"falls back to subject when none", "", "", "", guid, guid},
		{"multi word upn preserved", "", "", "José María del Río", guid, "José María del Río"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := pickUsername(c.preferred, c.unique, c.upn, c.id)
			if got != c.want {
				t.Errorf("pickUsername(%q,%q,%q,%q) = %q, want %q",
					c.preferred, c.unique, c.upn, c.id, got, c.want)
			}
		})
	}
}
