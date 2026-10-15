package identity

import "errors"

var (
	// ErrNotDigits indicates a non-decimal byte (3GPP TS 23.003 V20.1.0
	// §§2.2, 2.3, 6.2.1, 6.2.2).
	ErrNotDigits = errors.New("identity: IMSI, IMEI, IMEISV, MCC and MNC must contain digits 0-9 only per 3GPP TS 23.003 V20.1.0 §§2.2, 2.3, 6.2.1, 6.2.2")
	// ErrEmpty indicates an empty identity or component (3GPP TS 23.003
	// V20.1.0 §§2.2, 2.3, 6.2.1, 6.2.2).
	ErrEmpty = errors.New("identity: IMSI, IMEI, IMEISV, MCC or MNC holds no digits per 3GPP TS 23.003 V20.1.0 §§2.2, 2.3, 6.2.1, 6.2.2")
	// ErrIMSILength indicates an IMSI outside the structure of 3GPP TS 23.003
	// V20.1.0 §§2.2, 2.3 (6..15 digits).
	ErrIMSILength = errors.New("identity: IMSI must have 6 to 15 digits: a 3-digit MCC, a 2- or 3-digit MNC and at least one MSIN digit per 3GPP TS 23.003 V20.1.0 §§2.2, 2.3")
	// ErrIMEILength indicates an IMEI without 15 digits (3GPP TS 23.003 V20.1.0 §6.2.1).
	ErrIMEILength = errors.New("identity: IMEI must have 15 digits per 3GPP TS 23.003 V20.1.0 §6.2.1")
	// ErrIMEISpareDigitNotZero indicates a nonzero transmitted spare digit
	// (3GPP TS 23.003 V20.1.0 §6.2.1).
	ErrIMEISpareDigitNotZero = errors.New("identity: an IMEI transmitted with its spare digit must end in 0 per 3GPP TS 23.003 V20.1.0 §6.2.1")
	// ErrIMEISVLength indicates an IMEISV without 16 digits (3GPP TS 23.003 V20.1.0 §6.2.2).
	ErrIMEISVLength = errors.New("identity: IMEISV must have 16 digits per 3GPP TS 23.003 V20.1.0 §6.2.2")
	// ErrPLMNIdentityDigits indicates an invalid MCC or MNC nibble
	// (3GPP TS 29.002 V19.1.0 §§17.7.8, 17.7.13).
	ErrPLMNIdentityDigits = errors.New("identity: MCC and MNC must contain decimal digits; only the third MNC digit may be filler per 3GPP TS 29.002 V19.1.0 §§17.7.8, 17.7.13")
	// ErrPLMNIdentityLength indicates a PLMN identity without three octets
	// (3GPP TS 29.002 V19.1.0 §17.7.8).
	ErrPLMNIdentityLength = errors.New("identity: PLMN-Id must have 3 octets per 3GPP TS 29.002 V19.1.0 §17.7.8")
	// ErrMCCLength indicates an MCC without three digits (3GPP TS 23.003 V20.1.0 §2.2).
	ErrMCCLength = errors.New("identity: MCC must have 3 digits per 3GPP TS 23.003 V20.1.0 §2.2")
	// ErrMNCLength indicates an MNC without two or three digits (3GPP TS 23.003 V20.1.0 §2.2).
	ErrMNCLength = errors.New("identity: MNC must have 2 or 3 digits per 3GPP TS 23.003 V20.1.0 §2.2")
	// ErrTACSNRLength indicates a TAC/SNR without 14 digits (3GPP TS 23.003 V20.1.0 §6.2.1, Annex B.1).
	ErrTACSNRLength = errors.New("identity: TAC and SNR must have 14 digits per 3GPP TS 23.003 V20.1.0 §6.2.1, Annex B.1")
	// ErrIMEICheckDigit indicates an incorrect IMEI check digit (3GPP TS 23.003 V20.1.0 §6.2.1, Annex B).
	ErrIMEICheckDigit = errors.New("identity: IMEI check digit must match the Luhn formula per 3GPP TS 23.003 V20.1.0 §6.2.1, Annex B")
	// ErrIMEISVReservedSVN indicates the reserved sender SVN (3GPP TS 23.003 V20.1.0 §6.3).
	ErrIMEISVReservedSVN = errors.New("identity: sent IMEISV must not use reserved SVN 99 per 3GPP TS 23.003 V20.1.0 §6.3")
)
