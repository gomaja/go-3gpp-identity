package identity

import "fmt"

type identity struct {
	min, max int
	err      error
}

var (
	imsi   = identity{min: 6, max: 15, err: ErrIMSILength}
	imei   = identity{min: 15, max: 15, err: ErrIMEILength}
	imeisv = identity{min: 16, max: 16, err: ErrIMEISVLength}
)

// ValidateIMSI checks for 6..15 decimal digits per 3GPP TS 23.003 V20.1.0
// §§2.2 and 2.3: a three-digit MCC, a two- or three-digit MNC and an MSIN.
func ValidateIMSI(digits string) error {
	return checkIdentityDigits(imsi, digits)
}

// ValidateIMEI checks for 15 decimal digits per 3GPP TS 23.003 V20.1.0 §6.2.1.
// It does not verify the Annex B check digit or require a zero spare digit.
// Callers handling the MAP IMEI type must select ValidateIMEISV for its
// 16-digit SVN form (3GPP TS 29.002 V19.1.0 §17.7.8).
func ValidateIMEI(digits string) error {
	return checkIdentityDigits(imei, digits)
}

// ValidateIMEISV checks for 16 decimal digits per 3GPP TS 23.003 V20.1.0 §6.2.2.
// It does not check SVN allocation restrictions in §6.3.
func ValidateIMEISV(digits string) error {
	return checkIdentityDigits(imeisv, digits)
}

// ValidateIMEISpareDigit validates an IMEI transmitted with its spare digit:
// 15 decimal digits ending in zero, per 3GPP TS 23.003 V20.1.0 §6.2.1.
func ValidateIMEISpareDigit(digits string) error {
	if err := ValidateIMEI(digits); err != nil {
		return err
	}
	if digits[len(digits)-1] != '0' {
		return fmt.Errorf("%w (got %q)", ErrIMEISpareDigitNotZero, digits[len(digits)-1])
	}
	return nil
}

func checkIdentityDigits(id identity, digits string) error {
	if digits == "" {
		return ErrEmpty
	}
	for i := 0; i < len(digits); i++ {
		if digits[i] < '0' || digits[i] > '9' {
			return fmt.Errorf("%w: %q at index %d", ErrNotDigits, digits[i], i)
		}
	}
	if len(digits) < id.min || len(digits) > id.max {
		return fmt.Errorf("%w (got %d digits)", id.err, len(digits))
	}
	return nil
}
