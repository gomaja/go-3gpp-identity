package identity_test

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	identity "github.com/gomaja/go-3gpp-identity"
)

func TestPLMNIdentityDigits(t *testing.T) {
	for pos := 0; pos < 6; pos++ {
		for digit := byte(0); digit < 16; digit++ {
			raw := []byte{0x21, 0x63, 0x54}
			shift := uint(pos%2) * 4
			raw[pos/2] = raw[pos/2]&^(0xf<<shift) | digit<<shift
			original := bytes.Clone(raw)
			var want error
			if digit > 9 && (pos != 3 || digit != 15) {
				want = identity.ErrPLMNIdentityDigits
			}
			err := identity.ValidatePLMNIdentity(raw)
			if !errors.Is(err, want) {
				t.Fatalf("nibble %d=%x: got %v, want %v", pos, digit, err, want)
			}
			if !bytes.Equal(raw, original) {
				t.Fatalf("input changed: %x to %x", original, raw)
			}
			if err != nil {
				label := []string{"MCC digit 1", "MCC digit 2", "MCC digit 3", "MNC digit 3", "MNC digit 1", "MNC digit 2"}[pos]
				if !strings.Contains(err.Error(), label) || !strings.HasPrefix(err.Error(), "identity: ") {
					t.Fatalf("missing digit context for %s: %v", label, err)
				}
			}
		}
	}
	for _, raw := range [][]byte{{0x21, 0xf3, 0x54}, {0x21, 0x63, 0x54}, {0, 0xf0, 0}, {0x99, 0x99, 0x99}} {
		if err := identity.ValidatePLMNIdentity(raw); err != nil {
			t.Errorf("%x: %v", raw, err)
		}
	}
}

func TestPLMNIdentityLength(t *testing.T) {
	for n := 0; n <= 7; n++ {
		raw := make([]byte, n)
		var want error
		if n != 3 {
			want = identity.ErrPLMNIdentityLength
		}
		if err := identity.ValidatePLMNIdentity(raw); !errors.Is(err, want) {
			t.Errorf("length %d: got %v, want %v", n, err, want)
		}
	}
	if err := identity.ValidatePLMNIdentity(nil); !errors.Is(err, identity.ErrPLMNIdentityLength) {
		t.Errorf("nil: %v", err)
	}
}

func TestPLMNIdentityErrorOrder(t *testing.T) {
	for _, tc := range []struct {
		raw     []byte
		want    error
		context string
	}{
		{[]byte{0xff, 0xff}, identity.ErrPLMNIdentityLength, "got 2 octets"},
		{[]byte{0xff, 0xff, 0xff, 0xff}, identity.ErrPLMNIdentityLength, "got 4 octets"},
		{[]byte{0xaa, 0xff, 0xaa}, identity.ErrPLMNIdentityDigits, "MCC digit 1"},
		{[]byte{0xa1, 0xfa, 0xaa}, identity.ErrPLMNIdentityDigits, "MCC digit 2"},
		{[]byte{0x21, 0xaa, 0xaa}, identity.ErrPLMNIdentityDigits, "MCC digit 3"},
		{[]byte{0x21, 0xa3, 0xaa}, identity.ErrPLMNIdentityDigits, "MNC digit 1"},
		{[]byte{0x21, 0xa3, 0xa4}, identity.ErrPLMNIdentityDigits, "MNC digit 2"},
		{[]byte{0x21, 0xa3, 0x54}, identity.ErrPLMNIdentityDigits, "MNC digit 3"},
	} {
		t.Run(fmt.Sprintf("%x", tc.raw), func(t *testing.T) {
			err := identity.ValidatePLMNIdentity(tc.raw)
			if !errors.Is(err, tc.want) || !strings.Contains(err.Error(), tc.context) {
				t.Fatalf("got %v, want %v with %q", err, tc.want, tc.context)
			}
		})
	}
}
