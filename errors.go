package identity

import "errors"

var (
	// ErrNotDigits indicates a byte outside the decimal alphabet.
	ErrNotDigits = errors.New("identity: IMSI, IMEI and IMEISV must contain digits 0-9 only per 3GPP TS 23.003 V20.1.0 §§2.3, 6.2.1, 6.2.2")
	// ErrEmpty indicates an identity with no digits.
	ErrEmpty = errors.New("identity: IMSI, IMEI or IMEISV holds no digits per 3GPP TS 23.003 V20.1.0 §§2.2, 6.2.1, 6.2.2")
	// ErrIMSILength indicates an IMSI outside 6..15 digits.
	ErrIMSILength = errors.New("identity: IMSI must have 6 to 15 digits per 3GPP TS 23.003 V20.1.0 §§2.2, 2.3")
	// ErrIMEILength indicates an IMEI without exactly 15 digits.
	ErrIMEILength = errors.New("identity: IMEI must have 15 digits per 3GPP TS 23.003 V20.1.0 §6.2.1")
	// ErrIMEISpareDigitNotZero indicates a nonzero transmitted spare digit.
	ErrIMEISpareDigitNotZero = errors.New("identity: an IMEI transmitted with its spare digit must end in 0 per 3GPP TS 23.003 V20.1.0 §6.2.1")
	// ErrIMEISVLength indicates an IMEISV without exactly 16 digits.
	ErrIMEISVLength = errors.New("identity: IMEISV must have 16 digits per 3GPP TS 23.003 V20.1.0 §6.2.2")
	// ErrPLMNIdentityDigits indicates an invalid MCC or MNC digit.
	ErrPLMNIdentityDigits = errors.New("identity: MCC and MNC must contain decimal digits; only the third MNC digit may be filler per 3GPP TS 29.002 V19.1.0 §§17.7.8, 17.7.13")
	// ErrPLMNIdentityLength indicates a PLMN identity without three octets.
	ErrPLMNIdentityLength = errors.New("identity: PLMN-Id must have 3 octets per 3GPP TS 29.002 V19.1.0 §17.7.8")
)
