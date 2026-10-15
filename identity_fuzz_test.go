package identity_test

import (
	"fmt"
	"strings"
	"testing"

	identity "github.com/gomaja/go-3gpp-identity"
)

// The reference scans runes using the standard library, returning byte offsets
// even for invalid UTF-8, independently of the production byte loop.
func referenceDigits(s string, min, max int, length error) (string, error) {
	if s == "" {
		return "", identity.ErrEmpty
	}
	if i := strings.IndexFunc(s, func(r rune) bool { return r < '0' || r > '9' }); i >= 0 {
		return fmt.Sprintf(": %q at index %d", s[i], i), identity.ErrNotDigits
	}
	if len(s) < min || len(s) > max {
		return fmt.Sprintf(" (got %d digits)", len(s)), length
	}
	return "", nil
}

// Try every candidate using the full-number ISO/IEC 7812 checksum, traversing
// right to left and summing decimal quotient/remainder instead of subtracting 9.
func referenceCD(s string) byte {
	for candidate := byte('0'); candidate <= '9'; candidate++ {
		all := s + string(candidate)
		sum := 0
		for i := len(all) - 1; i >= 0; i-- {
			d := int(all[i] - '0')
			if (len(all)-1-i)%2 == 1 {
				d *= 2
			}
			sum += d/10 + d%10
		}
		if sum%10 == 0 {
			return candidate
		}
	}
	panic("no Luhn digit")
}

var stringSeeds = []string{"", "0", "12345", "123456", "1234567", "001010123456789", "1234567890123456", "12345678901234567", "260531793113837", "490154203237518", "352099001761481", "0000000000000", "00000000000000", "260531793113830", "2605317931138398", "2605317931138399", "x", "12x34", "\x00\xff", "１２３４５６", "١٢٣٤٥٦"}

func FuzzIMSI(f *testing.F) {
	for _, s := range stringSeeds {
		for _, n := range []int{-1, 0, 1, 2, 3, 4} {
			f.Add(s, n)
		}
	}
	f.Fuzz(func(t *testing.T, s string, n int) {
		ctx, want := referenceDigits(s, 6, 15, identity.ErrIMSILength)
		assertError(t, identity.ValidateIMSI(s), want, ctx)
		if want == nil {
			switch {
			case n != 2 && n != 3:
				want = identity.ErrMNCLength
				ctx = fmt.Sprintf(" (got MNC length %d)", n)
			case len(s) < 3+n+1:
				want = identity.ErrIMSILength
				ctx = fmt.Sprintf(" (got %d digits; a %d-digit MNC needs at least %d)", len(s), n, 3+n+1)
			}
		}
		a, b, c, err := identity.SplitIMSI(s, n)
		assertError(t, err, want, ctx)
		if want != nil {
			if a != "" || b != "" || c != "" {
				t.Fatal("partial output on failure")
			}
			return
		}
		if len(a) != 3 || len(b) != n || len(c) < 1 || a+b+c != s {
			t.Fatalf("invalid split: %q %q %q", a, b, c)
		}
	})
}

func FuzzIMEI(f *testing.F) {
	for _, s := range stringSeeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		ctx, want := referenceDigits(s, 15, 15, identity.ErrIMEILength)
		assertError(t, identity.ValidateIMEI(s), want, ctx)
		spareWant, spareCtx := want, ctx
		if want == nil && s[14] != '0' {
			spareWant = identity.ErrIMEISpareDigitNotZero
			spareCtx = fmt.Sprintf(" (got %q)", s[14])
		}
		assertError(t, identity.ValidateIMEISpareDigit(s), spareWant, spareCtx)
		if want == nil {
			cd := referenceCD(s[:14])
			if cd != s[14] {
				want = identity.ErrIMEICheckDigit
				ctx = fmt.Sprintf(" (got %q, want %q)", s[14], cd)
			}
		}
		assertError(t, identity.VerifyIMEICheckDigit(s), want, ctx)
		ctx, want = referenceDigits(s, 14, 14, identity.ErrTACSNRLength)
		cd, err := identity.IMEICheckDigit(s)
		assertError(t, err, want, ctx)
		if want != nil {
			if cd != 0 {
				t.Fatal("nonzero CD on error")
			}
			return
		}
		if cd != referenceCD(s) {
			t.Fatalf("CD %q for %q", cd, s)
		}
	})
}

func FuzzIMEISV(f *testing.F) {
	for _, s := range stringSeeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		ctx, want := referenceDigits(s, 16, 16, identity.ErrIMEISVLength)
		assertError(t, identity.ValidateIMEISV(s), want, ctx)
		if want == nil && strings.HasSuffix(s, "99") {
			want = identity.ErrIMEISVReservedSVN
			ctx = " (got SVN 99)"
		}
		assertError(t, identity.ValidateSentIMEISV(s), want, ctx)
	})
}

func FuzzDecodePLMN(f *testing.F) {
	for _, b := range [][]byte{nil, {}, {0, 0}, {0, 0xf1, 0x10}, {0x13, 0, 0x62}, {0x62, 0xf2, 0x10}, {0xff, 0xff, 0xff}, {0, 0, 0, 0}, {0xa0, 0, 0}, {0, 0xa0, 0}} {
		f.Add(b)
	}
	f.Fuzz(func(t *testing.T, b []byte) {
		ctx, want := referencePLMNError(b)
		assertError(t, identity.ValidatePLMNIdentity(b), want, ctx)
		p, err := identity.DecodePLMN(b)
		assertError(t, err, want, ctx)
		if want != nil {
			if p != (identity.PLMN{}) {
				t.Fatal("partial PLMN")
			}
			return
		}
		expected, ok := referencePLMN(b)
		if !ok || p != expected {
			t.Fatalf("%x: %+v want %+v", b, p, expected)
		}
		raw := p.Bytes()
		if string(raw[:]) != string(b) {
			t.Fatalf("%x re-encoded as %x", b, raw)
		}
	})
}

func referenceField(s, field string, min, max int, length error) (string, error) {
	ctx, err := referenceDigits(s, min, max, length)
	if err == nil {
		return "", nil
	}
	switch err {
	case length:
		return fmt.Sprintf(" (got %d %s digits)", len(s), field), err
	case identity.ErrEmpty:
		return " (got an empty " + field + ")", err
	}
	return ctx + " of the " + field, err
}

func FuzzNewPLMN(f *testing.F) {
	for _, p := range [][2]string{{"", ""}, {"001", ""}, {"001", "01"}, {"310", "260"}, {"262", "01"}, {"000", "00"}, {"999", "999"}, {"a", "01"}, {"001", "x"}, {"00", "01"}, {"0000", "01"}, {"001", "1"}, {"001", "1234"}, {"１２３", "01"}} {
		f.Add(p[0], p[1])
	}
	f.Fuzz(func(t *testing.T, mcc, mnc string) {
		ctx, want := referenceField(mcc, "MCC", 3, 3, identity.ErrMCCLength)
		if want == nil {
			ctx, want = referenceField(mnc, "MNC", 2, 3, identity.ErrMNCLength)
		}
		p, err := identity.NewPLMN(mcc, mnc)
		assertError(t, err, want, ctx)
		raw := (identity.PLMN{MCC: mcc, MNC: mnc}).Bytes()
		decoded, decodeErr := identity.DecodePLMN(raw[:])
		if want != nil {
			if p != (identity.PLMN{}) || decodeErr == nil || raw != [3]byte{255, 255, 255} {
				t.Fatal("invalid PLMN accepted")
			}
			return
		}
		if decodeErr != nil || decoded != p || p != (identity.PLMN{MCC: mcc, MNC: mnc}) || raw != referenceEncode(mcc, mnc) {
			t.Fatalf("%q-%q: %+v %x %v", mcc, mnc, p, raw, decodeErr)
		}
	})
}
