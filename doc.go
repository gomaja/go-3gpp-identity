// Package identity validates and splits IMSIs, validates IMEIs and IMEISVs,
// computes and verifies IMEI check digits, and encodes and decodes PLMN-Ids.
// It follows 3GPP TS 23.003 V20.1.0 §§2.2, 2.3, 6.2.1, 6.2.2, 6.3, 12.1
// and Annex B, and TS 29.002 V19.1.0 §§17.7.8, 17.7.13.
//
// ValidateIMSI checks decimal digits and length. The six-digit minimum
// follows from the structure in TS 23.003 V20.1.0 §2.2: a three-digit MCC,
// a two- or three-digit MNC, and at least one MSIN digit. Section 2.3 states
// the fifteen-digit maximum. SplitIMSI requires the configured MNC length;
// it cannot infer that length from the IMSI (§2.2 NOTE 1).
//
// ValidateIMEI and ValidateIMEISV check decimal digits and length under
// TS 23.003 V20.1.0 §§6.2.1, 6.2.2. IMEICheckDigit and VerifyIMEICheckDigit
// handle the Annex B check digit, used for manual entry. Sender checks are
// separate: ValidateIMEISpareDigit requires the transmitted spare digit 0
// (§6.2.1); ValidateSentIMEISV rejects reserved SVN 99 (§6.3). Receivers can
// use ValidateIMEISV to retain a reserved SVN.
//
// NewPLMN validates decimal MCC and MNC strings (TS 23.003 V20.1.0 §2.2).
// DecodePLMN, PLMN.Bytes and ValidatePLMNIdentity handle the three-octet
// PLMN-Id coding (TS 29.002 V19.1.0 §17.7.8). PLMN.String preserves the
// component strings, including leading zeros.
//
// Errors match exported sentinels through errors.Is. String checks report
// empty input before non-decimal bytes, then length, then any specific rule.
// Wrapped errors add the byte index, field, digit or octet that failed.
//
// TBCD coding of whole identities, as MAP carries them (TS 29.002 V19.1.0
// §17.7.8), belongs to github.com/gomaja/go-tbcd.
package identity
