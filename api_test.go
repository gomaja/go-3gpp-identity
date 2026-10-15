package identity_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	identity "github.com/gomaja/go-3gpp-identity"
)

func assertError(t *testing.T, got, want error, context string) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if want != nil && !strings.HasSuffix(got.Error(), context) {
		t.Fatalf("got %v, want suffix %q", got, context)
	}
}

func TestSplitIMSI(t *testing.T) {
	for _, tc := range []struct {
		digits         string
		n              int
		mcc, mnc, msin string
	}{
		{"001010", 2, "001", "01", "0"}, {"0010100", 3, "001", "010", "0"},
		{"310260123456789", 3, "310", "260", "123456789"}, {"001010123456789", 2, "001", "01", "0123456789"},
	} {
		a, b, c, err := identity.SplitIMSI(tc.digits, tc.n)
		if err != nil || a != tc.mcc || b != tc.mnc || c != tc.msin {
			t.Fatalf("%q/%d: %q %q %q %v", tc.digits, tc.n, a, b, c, err)
		}
	}
	for _, tc := range []struct {
		digits  string
		n       int
		want    error
		context string
	}{
		{"", 4, identity.ErrEmpty, ""}, {"12a", 4, identity.ErrNotDigits, ": 'a' at index 2"},
		{"12345", 4, identity.ErrIMSILength, " (got 5 digits)"}, {"1234567890123456", 2, identity.ErrIMSILength, " (got 16 digits)"},
		{"001010", 3, identity.ErrIMSILength, " (got 6 digits; a 3-digit MNC needs at least 7)"},
	} {
		a, b, c, err := identity.SplitIMSI(tc.digits, tc.n)
		assertError(t, err, tc.want, tc.context)
		if a != "" || b != "" || c != "" {
			t.Fatalf("partial output on error: %q %q %q", a, b, c)
		}
	}
	for _, n := range []int{-1, 0, 1, 4} {
		a, b, c, err := identity.SplitIMSI("001010123456789", n)
		assertError(t, err, identity.ErrMNCLength, fmt.Sprintf(" (got MNC length %d)", n))
		if a != "" || b != "" || c != "" {
			t.Fatal("partial output on error")
		}
	}
}

func TestIMEICheckDigit(t *testing.T) {
	// First vector: 3GPP TS 23.003 V20.1.0 Annex B.3. Other vectors were
	// checked by an independent ISO/IEC 7812 Luhn implementation.
	for _, digits := range []string{"260531793113837", "490154203237518", "352099001761481", "000000000000000"} {
		cd, err := identity.IMEICheckDigit(digits[:14])
		if err != nil || cd != digits[14] {
			t.Fatalf("%s: %q %v", digits, cd, err)
		}
		assertError(t, identity.VerifyIMEICheckDigit(digits), nil, "")
	}
	for pos := range 15 {
		for d := byte('0'); d <= '9'; d++ {
			changed := []byte("260531793113837")
			if changed[pos] == d {
				continue
			}
			changed[pos] = d
			cd := referenceCD(string(changed[:14]))
			assertError(t, identity.VerifyIMEICheckDigit(string(changed)), identity.ErrIMEICheckDigit, fmt.Sprintf(" (got %q, want %q)", changed[14], cd))
		}
	}
	for _, tc := range []struct {
		input   string
		want    error
		context string
	}{
		{"", identity.ErrEmpty, ""}, {"a", identity.ErrNotDigits, ": 'a' at index 0"},
		{"0000000000000", identity.ErrTACSNRLength, " (got 13 digits)"}, {"000000000000000", identity.ErrTACSNRLength, " (got 15 digits)"},
	} {
		cd, err := identity.IMEICheckDigit(tc.input)
		assertError(t, err, tc.want, tc.context)
		if cd != 0 {
			t.Fatal("nonzero check digit on error")
		}
	}
	for _, s := range []string{"", "x", "00000000000000", "0000000000000000"} {
		ctx, want := referenceDigits(s, 15, 15, identity.ErrIMEILength)
		assertError(t, identity.VerifyIMEICheckDigit(s), want, ctx)
	}
}

func TestValidateSentIMEISV(t *testing.T) {
	for n := range 100 {
		s := fmt.Sprintf("26053179311383%02d", n)
		assertError(t, identity.ValidateIMEISV(s), nil, "")
		var want error
		ctx := ""
		if n == 99 {
			want = identity.ErrIMEISVReservedSVN
			ctx = " (got SVN 99)"
		}
		assertError(t, identity.ValidateSentIMEISV(s), want, ctx)
	}
	for _, s := range []string{"", "x", "000000000000000", "00000000000000000"} {
		ctx, want := referenceDigits(s, 16, 16, identity.ErrIMEISVLength)
		assertError(t, identity.ValidateSentIMEISV(s), want, ctx)
	}
}

func TestNewPLMN(t *testing.T) {
	for _, tc := range []struct {
		mcc, mnc string
		want     error
		context  string
	}{
		{"", "", identity.ErrEmpty, " (got an empty MCC)"}, {"a", "", identity.ErrNotDigits, ": 'a' at index 0 of the MCC"},
		{"001a", "", identity.ErrNotDigits, ": 'a' at index 3 of the MCC"}, {"00", "a", identity.ErrMCCLength, " (got 2 MCC digits)"},
		{"0000", "01", identity.ErrMCCLength, " (got 4 MCC digits)"}, {"001", "", identity.ErrEmpty, " (got an empty MNC)"},
		{"001", "1", identity.ErrMNCLength, " (got 1 MNC digits)"}, {"001", "1234", identity.ErrMNCLength, " (got 4 MNC digits)"},
		{"001", "0a", identity.ErrNotDigits, ": 'a' at index 1 of the MNC"}, {"001", "123a", identity.ErrNotDigits, ": 'a' at index 3 of the MNC"},
	} {
		p, err := identity.NewPLMN(tc.mcc, tc.mnc)
		assertError(t, err, tc.want, tc.context)
		if p != (identity.PLMN{}) {
			t.Fatal("partial PLMN on error")
		}
	}
	for _, field := range []string{"MCC", "MNC"} {
		for pos := range 3 {
			for d := 0; d < 256; d++ {
				if d >= '0' && d <= '9' {
					continue
				}
				s := []byte("001")
				s[pos] = byte(d)
				mcc, mnc := "001", "010"
				if field == "MCC" {
					mcc = string(s)
				} else {
					mnc = string(s)
				}
				_, err := identity.NewPLMN(mcc, mnc)
				assertError(t, err, identity.ErrNotDigits, fmt.Sprintf(": %q at index %d of the %s", byte(d), pos, field))
			}
		}
	}
}

func TestPLMNVectors(t *testing.T) {
	// Literal octets derived from 3GPP TS 29.002 V19.1.0 §17.7.8:
	// MCC1/MCC2, MCC3/MNC3-or-filler, MNC1/MNC2 in low/high nibbles.
	for _, tc := range []struct {
		mcc, mnc string
		raw      [3]byte
	}{
		{"001", "01", [3]byte{0x00, 0xf1, 0x10}}, {"310", "260", [3]byte{0x13, 0x00, 0x62}}, {"262", "01", [3]byte{0x62, 0xf2, 0x10}},
	} {
		p, err := identity.NewPLMN(tc.mcc, tc.mnc)
		if err != nil || p.MCC != tc.mcc || p.MNC != tc.mnc || p.Bytes() != tc.raw {
			t.Fatalf("%+v: %+v %v", tc, p, err)
		}
		decoded, err := identity.DecodePLMN(tc.raw[:])
		if err != nil || decoded != p || p.String() != tc.mcc+"-"+tc.mnc {
			t.Fatalf("%+v: %+v %v", tc, decoded, err)
		}
	}
	for _, p := range []identity.PLMN{{}, {MCC: "001"}, {MCC: "01", MNC: "01"}, {MCC: "001", MNC: "1"}, {MCC: "001", MNC: "1234"}, {MCC: "abc", MNC: "01"}, {MCC: "001", MNC: "a1"}} {
		raw := p.Bytes()
		if raw != [3]byte{255, 255, 255} {
			t.Fatalf("invalid %+v encoded as %x", p, raw)
		}
		_, err := identity.DecodePLMN(raw[:])
		assertError(t, err, identity.ErrPLMNIdentityDigits, ": MCC digit 1 in octet 1 is 0xf")
		if p.String() != p.MCC+"-"+p.MNC {
			t.Fatal("String changed invalid value")
		}
	}
}
