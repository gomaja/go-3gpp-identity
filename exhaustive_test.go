package identity_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	identity "github.com/gomaja/go-3gpp-identity"
)

// Read one integer, select base-16 places in decimal-field order, and form
// decimal numbers. This reference does not use production's nibble expressions.
func referencePLMN(b []byte) (identity.PLMN, bool) {
	if len(b) != 3 {
		return identity.PLMN{}, false
	}
	v := int(b[0]) + 256*int(b[1]) + 65536*int(b[2])
	places := [6]int{1, 16, 256, 65536, 1048576, 4096}
	var d [6]int
	for i, p := range places {
		d[i] = v / p % 16
		if d[i] > 9 && (i != 5 || d[i] != 15) {
			return identity.PLMN{}, false
		}
	}
	mnc := fmt.Sprintf("%02d", 10*d[3]+d[4])
	if d[5] != 15 {
		mnc += strconv.Itoa(d[5])
	}
	return identity.PLMN{MCC: fmt.Sprintf("%03d", 100*d[0]+10*d[1]+d[2]), MNC: mnc}, true
}

func referencePLMNError(b []byte) (string, error) {
	if len(b) != 3 {
		return fmt.Sprintf(" (got %d octets)", len(b)), identity.ErrPLMNIdentityLength
	}
	v := int(b[0]) + 256*int(b[1]) + 65536*int(b[2])
	for i, p := range []int{1, 16, 256, 65536, 1048576, 4096} {
		d := v / p % 16
		if d <= 9 || (i == 5 && d == 15) {
			continue
		}
		field := "MCC"
		digit := i + 1
		octet := []int{1, 1, 2, 3, 3, 2}[i]
		if i >= 3 {
			field = "MNC"
			digit = i - 2
		}
		return fmt.Sprintf(": %s digit %d in octet %d is %#x", field, digit, octet, d), identity.ErrPLMNIdentityDigits
	}
	return "", nil
}

// Encoding as a base-16 integer from independently parsed decimal components.
func referenceEncode(mcc, mnc string) [3]byte {
	c, _ := strconv.Atoi(mcc)
	n, _ := strconv.Atoi(mnc)
	third := 15
	if len(mnc) == 3 {
		third = n % 10
		n /= 10
	}
	v := c/100 + 16*(c/10%10) + 256*(c%10) + 4096*third + 65536*(n/10) + 1048576*(n%10)
	return [3]byte{byte(v % 256), byte(v / 256 % 256), byte(v / 65536)}
}

func TestPLMNExhaustiveOctets(t *testing.T) {
	// The sweep runs in one goroutine, so the race detector has nothing to
	// find in it and only makes it about ten times slower.
	if testing.Short() || raceEnabled {
		t.Skip("all 2^24 three-octet inputs; skipped with -short and -race")
	}
	start := time.Now()
	valid := 0
	for v := 0; v < 1<<24; v++ {
		b := [3]byte{byte(v), byte(v >> 8), byte(v >> 16)}
		want, ok := referencePLMN(b[:])
		p, err := identity.DecodePLMN(b[:])
		validateErr := identity.ValidatePLMNIdentity(b[:])
		if (err == nil) != ok || (validateErr == nil) != ok {
			t.Fatalf("%x: %+v %v / %v; valid=%v", b, p, err, validateErr, ok)
		}
		if !ok {
			if p != (identity.PLMN{}) || !errors.Is(err, identity.ErrPLMNIdentityDigits) || !errors.Is(validateErr, identity.ErrPLMNIdentityDigits) {
				t.Fatalf("%x: invalid failure result", b)
			}
			continue
		}
		valid++
		if p != want || p.Bytes() != b || p.String() != want.MCC+"-"+want.MNC {
			t.Fatalf("%x: %+v want %+v", b, p, want)
		}
	}
	if valid != 1100000 {
		t.Fatalf("valid count %d", valid)
	}
	t.Logf("16777216 inputs; %d valid; elapsed %s", valid, time.Since(start))
}

func TestPLMNExhaustiveValid(t *testing.T) {
	for c := range 1000 {
		mcc := fmt.Sprintf("%03d", c)
		for n := range 1100 {
			var mnc string
			if n < 100 {
				mnc = fmt.Sprintf("%02d", n)
			} else {
				mnc = fmt.Sprintf("%03d", n-100)
			}
			p, err := identity.NewPLMN(mcc, mnc)
			if err != nil || p != (identity.PLMN{MCC: mcc, MNC: mnc}) {
				t.Fatalf("%s-%s: %+v %v", mcc, mnc, p, err)
			}
			b := p.Bytes()
			if b != referenceEncode(mcc, mnc) {
				t.Fatalf("%+v: %x", p, b)
			}
			decoded, err := identity.DecodePLMN(b[:])
			if err != nil || decoded != p {
				t.Fatalf("%+v: %+v %v", p, decoded, err)
			}
		}
	}
}
