package identity_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	identity "github.com/gomaja/go-3gpp-identity"
)

func TestIdentityLengths(t *testing.T) {
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
		t.Run(tc.name, func(t *testing.T) {
			for n := 0; n <= 20; n++ {
				digits := strings.Repeat("0", n)
				var want error
				if n == 0 {
					want = identity.ErrEmpty
				} else if n < tc.min || n > tc.max {
					want = tc.length
				}
				if err := tc.validate(digits); !errors.Is(err, want) {
					t.Errorf("length %d: got %v, want %v", n, err, want)
				}
			}
		})
	}
}

func TestIdentityDecimalDigits(t *testing.T) {
	for _, tc := range []struct {
		name     string
		validate func(string) error
		digits   string
	}{
		{"IMSI", identity.ValidateIMSI, "204080012345678"},
		{"IMEI", identity.ValidateIMEI, "490154203237510"},
		{"IMEISV", identity.ValidateIMEISV, "4901542032375188"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for pos := range len(tc.digits) {
				for value := 0; value <= 255; value++ {
					input := []byte(tc.digits)
					input[pos] = byte(value)
					var want error
					if value < '0' || value > '9' {
						want = identity.ErrNotDigits
					}
					err := tc.validate(string(input))
					if !errors.Is(err, want) {
						t.Fatalf("byte %02x at %d: got %v, want %v", value, pos, err, want)
					}
					if err != nil && !strings.HasSuffix(err.Error(), fmt.Sprintf(": %q at index %d", byte(value), pos)) {
						t.Fatalf("missing byte context: %v", err)
					}
				}
			}
			for _, input := range []string{"*", "#", "abc", "１２３４５６", "١٢٣٤٥٦", "12345678901234567890x", "1 2", "\x00", "/:", "a#"} {
				if err := tc.validate(input); !errors.Is(err, identity.ErrNotDigits) {
					t.Errorf("%q: got %v, want ErrNotDigits before length", input, err)
				}
			}
		})
	}
}

func TestIMSIExamples(t *testing.T) {
	for _, tc := range []struct {
		digits string
		want   error
	}{
		{"123456", nil},
		{"204080", nil},
		{"204080012345678", nil},
		{"001", identity.ErrIMSILength},
		{"12345", identity.ErrIMSILength},
		{"20408", identity.ErrIMSILength},
		{"1234567890123456", identity.ErrIMSILength},
		{"2040800123456789", identity.ErrIMSILength},
		{"20408001234567a", identity.ErrNotDigits},
		{"20408001234567*", identity.ErrNotDigits},
		{"20408a", identity.ErrNotDigits},
		{"12345a", identity.ErrNotDigits},
		{"*#*#*#", identity.ErrNotDigits},
		{"", identity.ErrEmpty},
	} {
		t.Run(tc.digits, func(t *testing.T) {
			if err := identity.ValidateIMSI(tc.digits); !errors.Is(err, tc.want) {
				t.Fatalf("got %v, want %v", err, tc.want)
			}
		})
	}
}

func TestIMEISpareDigit(t *testing.T) {
	for digit := byte('0'); digit <= '9'; digit++ {
		digits := "49015420323751" + string(digit)
		if err := identity.ValidateIMEI(digits); err != nil {
			t.Fatalf("ValidateIMEI(%q): %v", digits, err)
		}
		var want error
		if digit != '0' {
			want = identity.ErrIMEISpareDigitNotZero
		}
		if err := identity.ValidateIMEISpareDigit(digits); !errors.Is(err, want) {
			t.Errorf("ValidateIMEISpareDigit(%q): got %v, want %v", digits, err, want)
		}
	}
	for _, tc := range []struct {
		digits string
		want   error
	}{
		{"", identity.ErrEmpty},
		{"1", identity.ErrIMEILength},
		{"49015420323751", identity.ErrIMEILength},
		{"4901542032375100", identity.ErrIMEILength},
		{"4901542032375199", identity.ErrIMEILength},
		{"49015420323751a", identity.ErrNotDigits},
		{"4901542032375a1", identity.ErrNotDigits},
		{"a", identity.ErrNotDigits},
	} {
		if err := identity.ValidateIMEISpareDigit(tc.digits); !errors.Is(err, tc.want) {
			t.Errorf("%q: got %v, want %v", tc.digits, err, tc.want)
		}
	}
}

func TestIMEISVSoftwareVersion(t *testing.T) {
	for _, digits := range []string{"4901542032375100", "4901542032375188", "4901542032375199"} {
		if err := identity.ValidateIMEISV(digits); err != nil {
			t.Errorf("%q: %v", digits, err)
		}
	}
}

func TestErrorDescriptions(t *testing.T) {
	for _, err := range []error{
		identity.ErrEmpty, identity.ErrNotDigits, identity.ErrIMSILength,
		identity.ErrIMEILength, identity.ErrIMEISVLength, identity.ErrIMEISpareDigitNotZero,
		identity.ErrPLMNIdentityDigits, identity.ErrPLMNIdentityLength,
	} {
		if !strings.HasPrefix(err.Error(), "identity: ") || !strings.Contains(err.Error(), "§") {
			t.Errorf("missing package prefix or clause: %v", err)
		}
	}
}

func TestIdentityErrorContext(t *testing.T) {
	for _, tc := range []struct {
		err  error
		want string
	}{
		{identity.ValidateIMSI("12a#"), `: 'a' at index 2`},
		{identity.ValidateIMSI("12345"), "(got 5 digits)"},
		{identity.ValidateIMEI("1234567890123456"), "(got 16 digits)"},
		{identity.ValidateIMEISV("123456789012345"), "(got 15 digits)"},
		{identity.ValidateIMEISpareDigit("490154203237518"), "(got '8')"},
	} {
		if tc.err == nil || !strings.HasSuffix(tc.err.Error(), tc.want) {
			t.Errorf("got %v, want suffix %q", tc.err, tc.want)
		}
	}
}
