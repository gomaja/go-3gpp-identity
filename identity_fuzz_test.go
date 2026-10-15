package identity_test

import (
	"errors"
	"strings"
	"testing"

	identity "github.com/gomaja/go-3gpp-identity"
)

func FuzzValidateIdentities(f *testing.F) {
	for _, digits := range []string{
		"", "001", "12345", "123456", "204080012345678", "1234567890123456",
		"490154203237510", "490154203237518", "4901542032375100", "4901542032375199",
		"4901542032375188", "12345678901234567", "20408001234567a", "20408001234567*",
		"*#abc", "１２３４５６", "١٢٣٤٥٦", "\x00\xff", "0\n", "+123456",
	} {
		f.Add(digits)
	}
	f.Fuzz(func(t *testing.T, digits string) {
		for _, tc := range []struct {
			name     string
			validate func(string) error
			min, max int
			length   error
		}{
			{"IMSI", identity.ValidateIMSI, 6, 15, identity.ErrIMSILength},
			{"IMEI", identity.ValidateIMEI, 15, 15, identity.ErrIMEILength},
			{"IMEISV", identity.ValidateIMEISV, 16, 16, identity.ErrIMEISVLength},
		} {
			var want error
			switch {
			case digits == "":
				want = identity.ErrEmpty
			case strings.ContainsFunc(digits, func(r rune) bool { return !strings.ContainsRune("0123456789", r) }):
				want = identity.ErrNotDigits
			case len(digits) < tc.min || len(digits) > tc.max:
				want = tc.length
			}
			if err := tc.validate(digits); !errors.Is(err, want) {
				t.Fatalf("%s(%q): got %v, want %v", tc.name, digits, err, want)
			}
		}
	})
}
