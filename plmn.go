package identity

import "fmt"

// ValidatePLMNIdentity checks a three-octet PLMN-Id per 3GPP TS 29.002
// V19.1.0 §17.7.8 (also used in §17.7.13). Low/high nibbles contain
// MCC1/MCC2, MCC3/MNC3 and MNC1/MNC2 respectively. All digits are decimal;
// only MNC3 may be 1111, indicating a two-digit MNC.
func ValidatePLMNIdentity(octets []byte) error {
	if len(octets) != 3 {
		return fmt.Errorf("%w (got %d octets)", ErrPLMNIdentityLength, len(octets))
	}
	for i, d := range [3]byte{octets[0] & 0xf, octets[0] >> 4, octets[1] & 0xf} {
		if d > 9 {
			return fmt.Errorf("%w: MCC digit %d in octet %d is %#x", ErrPLMNIdentityDigits, i+1, i/2+1, d)
		}
	}
	for i, d := range [3]byte{octets[2] & 0xf, octets[2] >> 4, octets[1] >> 4} {
		if d > 9 && (i != 2 || d != 0xf) {
			return fmt.Errorf("%w: MNC digit %d in octet %d is %#x", ErrPLMNIdentityDigits, i+1, 3-i/2, d)
		}
	}
	return nil
}

// PLMN holds decimal MCC and MNC strings, preserving leading zeros and the
// MNC length (3GPP TS 23.003 V20.1.0 §§2.2, 12.1). Use NewPLMN to validate
// a value; Bytes also checks values constructed directly.
type PLMN struct{ MCC, MNC string }

// NewPLMN checks a three-digit MCC followed by a two- or three-digit MNC
// per 3GPP TS 23.003 V20.1.0 §2.2. On failure it returns the zero PLMN.
func NewPLMN(mcc, mnc string) (PLMN, error) {
	if err := checkPLMNField(mcc, "MCC", identity{min: 3, max: 3, err: ErrMCCLength}); err != nil {
		return PLMN{}, err
	}
	if err := checkPLMNField(mnc, "MNC", identity{min: 2, max: 3, err: ErrMNCLength}); err != nil {
		return PLMN{}, err
	}
	return PLMN{MCC: mcc, MNC: mnc}, nil
}

func checkPLMNField(digits, field string, id identity) error {
	if digits == "" {
		return fmt.Errorf("%w (got an empty %s)", ErrEmpty, field)
	}
	for i := range len(digits) {
		if digits[i] < '0' || digits[i] > '9' {
			return fmt.Errorf("%w: %q at index %d of the %s", ErrNotDigits, digits[i], i, field)
		}
	}
	if len(digits) < id.min || len(digits) > id.max {
		return fmt.Errorf("%w (got %d %s digits)", id.err, len(digits), field)
	}
	return nil
}

// DecodePLMN reads exactly three octets per 3GPP TS 29.002 V19.1.0
// §17.7.8 (also used in §17.7.13). It applies ValidatePLMNIdentity first,
// preserving a two-digit MNC when MNC digit 3 is filler. On failure it
// returns the zero PLMN.
func DecodePLMN(b []byte) (PLMN, error) {
	if err := ValidatePLMNIdentity(b); err != nil {
		return PLMN{}, err
	}
	digits := [6]byte{
		'0' + (b[0] & 0xf), '0' + (b[0] >> 4), '0' + (b[1] & 0xf),
		'0' + (b[2] & 0xf), '0' + (b[2] >> 4), '0' + (b[1] >> 4),
	}
	n := len(digits)
	if b[1]>>4 == 0xf {
		n-- // a two-digit MNC
	}
	s := string(digits[:n])
	return PLMN{MCC: s[:3], MNC: s[3:]}, nil
}

// Bytes encodes p per 3GPP TS 29.002 V19.1.0 §17.7.8, using filler in
// MNC digit 3 for a two-digit MNC. If NewPLMN rejects p, Bytes returns
// [3]byte{0xff, 0xff, 0xff}, which DecodePLMN rejects. It never panics.
func (p PLMN) Bytes() [3]byte {
	if _, err := NewPLMN(p.MCC, p.MNC); err != nil {
		return [3]byte{0xff, 0xff, 0xff}
	}
	third := byte(0xf)
	if len(p.MNC) == 3 {
		third = p.MNC[2] - '0'
	}
	return [3]byte{(p.MCC[0] - '0') | (p.MCC[1]-'0')<<4, (p.MCC[2] - '0') | third<<4, (p.MNC[0] - '0') | (p.MNC[1]-'0')<<4}
}

// String returns p.MCC + "-" + p.MNC, including for an invalid PLMN.
func (p PLMN) String() string { return p.MCC + "-" + p.MNC }
