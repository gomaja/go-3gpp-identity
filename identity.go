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
// §§2.2 and 2.3: a three-digit MCC, a two- or three-digit MNC and at least
// one MSIN digit. The six-digit minimum follows from that structure; §2.3
// states the fifteen-digit maximum.
func ValidateIMSI(digits string) error {
	return checkIdentityDigits(imsi, digits)
}

// ValidateIMEI checks for 15 decimal digits per 3GPP TS 23.003 V20.1.0 §6.2.1.
// It does not verify the Annex B check digit or require a zero spare digit
// (3GPP TS 23.003 V20.1.0 §6.2.1).
// Callers handling the MAP IMEI type must select ValidateIMEISV for its
// 16-digit SVN form (3GPP TS 29.002 V19.1.0 §17.7.8).
func ValidateIMEI(digits string) error {
	return checkIdentityDigits(imei, digits)
}

// ValidateIMEISV checks for 16 decimal digits per 3GPP TS 23.003 V20.1.0 §6.2.2.
// It accepts reserved SVN 99 for receivers; use ValidateSentIMEISV for the
// sender restriction (3GPP TS 23.003 V20.1.0 §6.3).
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

// SplitIMSI separates the MCC, MNC and MSIN per 3GPP TS 23.003 V20.1.0
// §§2.2, 2.3. The caller supplies the configured MNC length (2 or 3);
// it cannot be inferred from the IMSI (§2.2 NOTE 1). At least one MSIN
// digit must remain. All returned strings are empty on error.
func SplitIMSI(digits string, mncLength int) (mcc, mnc, msin string, err error) {
	if err := ValidateIMSI(digits); err != nil {
		return "", "", "", err
	}
	if mncLength != 2 && mncLength != 3 {
		return "", "", "", fmt.Errorf("%w (got MNC length %d)", ErrMNCLength, mncLength)
	}
	if len(digits) < 3+mncLength+1 {
		return "", "", "", fmt.Errorf("%w (got %d digits; a %d-digit MNC needs at least %d)", ErrIMSILength, len(digits), mncLength, 3+mncLength+1)
	}
	return digits[:3], digits[3 : 3+mncLength], digits[3+mncLength:], nil
}

// IMEICheckDigit computes the check digit for 14 decimal TAC/SNR digits per
// 3GPP TS 23.003 V20.1.0 §6.2.1 and Annex B.1–B.2. The result is an ASCII
// digit ('0'..'9'), suitable for appending to tacSNR. For an IMEISV, pass
// its first 14 digits, excluding the SVN (Annex B.1). On error it returns 0.
func IMEICheckDigit(tacSNR string) (byte, error) {
	if err := checkIdentityDigits(identity{min: 14, max: 14, err: ErrTACSNRLength}, tacSNR); err != nil {
		return 0, err
	}
	sum := 0
	for i := range len(tacSNR) {
		d := int(tacSNR[i] - '0')
		// Annex B.2 steps 1–2: D1 is the last SNR digit, so odd labels
		// occur at odd zero-based indexes. Sum the digits of each product.
		if i%2 == 1 {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
	}
	// Annex B.2 step 3 includes zero when the sum is a multiple of ten.
	return byte((10-sum%10)%10) + '0', nil
}

// VerifyIMEICheckDigit validates an IMEI including its check digit per
// 3GPP TS 23.003 V20.1.0 §6.2.1 and Annex B, for example an IMEI printed
// on or entered for a device. A transmitted IMEI carries the spare digit 0
// instead (§6.2.1; MAP: TS 29.002 V19.1.0 §17.7.8); use
// ValidateIMEISpareDigit for that form.
func VerifyIMEICheckDigit(digits string) error {
	if err := ValidateIMEI(digits); err != nil {
		return err
	}
	cd, _ := IMEICheckDigit(digits[:14]) // ValidateIMEI established 14 decimal TAC/SNR digits.
	if digits[14] != cd {
		return fmt.Errorf("%w (got %q, want %q)", ErrIMEICheckDigit, digits[14], cd)
	}
	return nil
}

// ValidateSentIMEISV applies the sender check of 3GPP TS 23.003 V20.1.0
// §§6.2.2, 6.3: 16 decimal digits with SVN other than reserved value 99.
// Receivers can use ValidateIMEISV to retain a reserved SVN.
func ValidateSentIMEISV(digits string) error {
	if err := ValidateIMEISV(digits); err != nil {
		return err
	}
	if digits[14:] == "99" {
		return fmt.Errorf("%w (got SVN %s)", ErrIMEISVReservedSVN, digits[14:])
	}
	return nil
}
