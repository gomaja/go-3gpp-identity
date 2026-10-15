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
			return fmt.Errorf("identity: MCC digit %d: %w", i+1, ErrPLMNIdentityDigits)
		}
	}
	for i, d := range [3]byte{octets[2] & 0xf, octets[2] >> 4, octets[1] >> 4} {
		if d > 9 && (i != 2 || d != 0xf) {
			return fmt.Errorf("identity: MNC digit %d: %w", i+1, ErrPLMNIdentityDigits)
		}
	}
	return nil
}
