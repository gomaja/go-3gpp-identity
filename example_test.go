package identity_test

import (
	"fmt"

	identity "github.com/gomaja/go-3gpp-identity"
)

func ExampleSplitIMSI() {
	mcc, mnc, msin, err := identity.SplitIMSI("001010123456789", 2)
	fmt.Println(mcc, mnc, msin, err)
	// Output: 001 01 0123456789 <nil>
}

func ExampleIMEICheckDigit() {
	cd, err := identity.IMEICheckDigit("26053179311383")
	fmt.Printf("%c %v\n", cd, err)
	// Output: 7 <nil>
}

func ExampleVerifyIMEICheckDigit() {
	fmt.Println(identity.VerifyIMEICheckDigit("260531793113837"))
	// Output: <nil>
}

func ExampleNewPLMN() {
	p, err := identity.NewPLMN("001", "01")
	fmt.Printf("%s % x %v\n", p, p.Bytes(), err)
	// Output: 001-01 00 f1 10 <nil>
}

func ExampleDecodePLMN() {
	p, err := identity.DecodePLMN([]byte{0x13, 0x00, 0x62})
	fmt.Println(p.String(), err)
	// Output: 310-260 <nil>
}
