package structmask_test

import (
	"encoding/json"
	"fmt"
	"github.com/afv1/custom_tags/customtags"
	"regexp"
	"strings"
)

type Card struct {
	Number string `json:"card_number" mask:"cardnumber"`
	Holder string `json:"card_holder"`
}

func Example() {
	// create struct instance.
	eg := Card{
		Number: "0000 0000 0000 0000",
		Holder: "TEST HOLDER",
	}

	// define handler for card number mask.
	// Could be separated from init function.
	cardMaskHandler := func(input string) string {
		input = strings.ReplaceAll(input, " ", "")

		const (
			binLn  = 6
			tailLn = 4
			symbol = "*"
		)

		maskLn := len(input) - binLn - tailLn

		exp := fmt.Sprintf(`(\d{%d}).*(\d{%d})`, binLn, tailLn)
		repl := "$1" + strings.Repeat(symbol, maskLn) + "$2"

		reg, _ := regexp.Compile(exp)

		return string(reg.ReplaceAll([]byte(input), []byte(repl)))
	}

	// init Custom Tags with tag name.
	ct := customtags.NewCustomTags("mask")
	// bind handler to custom tag label.
	customtags.Bind("cardnumber", cardMaskHandler)

	// print initial marshaled struct.
	initialJSON, _ := json.Marshal(eg)
	fmt.Println(string(initialJSON))

	// print masked marshaled struct.
	maskedJSON := ct.Proceed(eg)
	jsn, _ := json.Marshal(maskedJSON)
	fmt.Println(string(jsn))

	// Output:
	// {"card_number":"0000 0000 0000 0000","card_holder":"TEST HOLDER"}
	// {"card_number":"000000******0000","card_holder":"TEST HOLDER"}
}
