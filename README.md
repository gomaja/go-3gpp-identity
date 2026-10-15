# go-3gpp-identity

3GPP TS 23.003 identities for Go: validation of IMSI, IMEI and IMEISV digit
strings and of the three-octet PLMN identity.

## Status

The package holds the identity rules that `github.com/gomaja/go-asn1-gsmmap`
applied internally, moved here.
[Issue #1](https://github.com/gomaja/go-3gpp-identity/issues/1) describes the
complete scope this module is to implement, including IMSI splitting, the IMEI
check digit and PLMN identity encoding and decoding.

## Install

The module is used from its main branch; there are no tagged releases.

```
go get github.com/gomaja/go-3gpp-identity@main
```

## Usage

```go
err := identity.ValidateIMSI("001010123456789")
err = identity.ValidateIMEI("490154203237518")
err = identity.ValidateIMEISV("4901542032375181")
err = identity.ValidatePLMNIdentity([]byte{0x00, 0xf1, 0x10}) // MCC 001, MNC 01
```

TBCD coding of whole identities, as MAP carries them, is provided by
`github.com/gomaja/go-tbcd`.

## Specifications

[STANDARDS.md](STANDARDS.md) lists the documents and versions the code follows.

## License

MIT, see [LICENSE](LICENSE).
