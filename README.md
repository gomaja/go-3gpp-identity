# go-3gpp-identity

3GPP TS 23.003 identities for Go: IMSI splitting, IMEI check digits,
IMSI/IMEI/IMEISV validation, and three-octet PLMN identity coding.

## Status

The module implements the scope of
[issue #1](https://github.com/gomaja/go-3gpp-identity/issues/1), using only the
standard library. Location identities, SUPI, SUCI and temporary identities
are outside this scope.

## Install

The module is used from its main branch; there are no tagged releases.

```
go get github.com/gomaja/go-3gpp-identity@main
```

## Usage

```go
import identity "github.com/gomaja/go-3gpp-identity"

mcc, mnc, msin, err := identity.SplitIMSI("001010123456789", 2)
// mcc = "001", mnc = "01", msin = "0123456789", err = nil

cd, err := identity.IMEICheckDigit("26053179311383")
// cd = '7' (ASCII), err = nil
err = identity.VerifyIMEICheckDigit("260531793113837") // nil
err = identity.ValidateIMEISpareDigit("260531793113830") // nil

err = identity.ValidateIMEISV("2605317931138399") // nil: receiver retains SVN
err = identity.ValidateSentIMEISV("2605317931138398") // nil
err = identity.ValidateSentIMEISV("2605317931138399") // ErrIMEISVReservedSVN

p, err := identity.NewPLMN("001", "01") // err = nil
raw := p.Bytes() // [3]byte{0x00, 0xf1, 0x10}
p, err = identity.DecodePLMN(raw[:]) // err = nil
text := p.String() // "001-01"
```

Choose the MNC length from configuration; it cannot be inferred from an IMSI
(TS 23.003 V20.1.0 §2.2 NOTE 1). The six-digit IMSI minimum follows from its
MCC/MNC/MSIN structure (§2.2); §2.3 specifies the fifteen-digit maximum.

The check digit is for manual entry, such as a printed IMEI. Transmitted
IMEIs carry a zero spare digit (TS 23.003 V20.1.0 §6.2.1; MAP: TS 29.002
V19.1.0 §17.7.8). Use `ValidateIMEISpareDigit` for a sender's transmitted
IMEI, and `ValidateSentIMEISV` to reject reserved SVN 99 (§6.3).
`ValidateIMEI` and `ValidateIMEISV` check the decimal structure independently
of those sender checks (§§6.2.1, 6.2.2).

Errors match the exported sentinels with `errors.Is`; wrapped errors give
position or length context. A directly constructed invalid `PLMN` encodes as
`ff ff ff`, which `DecodePLMN` rejects.

TBCD coding of whole identities, as MAP carries them (TS 29.002 V19.1.0
§17.7.8), is provided by `github.com/gomaja/go-tbcd`.

## Specifications

[STANDARDS.md](STANDARDS.md) lists the documents, versions and clauses used.

## License

MIT, see [LICENSE](LICENSE).
