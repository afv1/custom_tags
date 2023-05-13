package structmask_test

import (
	"cardmasker/structmask"
	"encoding/json"
	"fmt"
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
		const (
			binLength  = 6
			tailLength = 4
			maskSymbol = "*"
			cardNumber = "CARD NUMBER"
		)

		input = strings.ReplaceAll(input, " ", "")

		cardNumberLength := len(input)
		if cardNumberLength != 16 {
			return cardNumber
		}

		bin := input[:binLength]
		tail := input[cardNumberLength-tailLength:]
		mask := ""
		if maskLength := cardNumberLength - binLength - tailLength; maskLength > 0 {
			mask = strings.Repeat(maskSymbol, maskLength)
		}

		return bin + mask + tail
	}

	// init mappers. Mapper is map[key]handler.
	// Key is tag label from example struct.
	mappers := structmask.Mapper{
		"cardnumber": cardMaskHandler,
	}

	// init StructMask with config.
	structmask.InitStructMask(&structmask.Config{
		TagName: "mask",
		Mappers: mappers,
	})

	// marshal initial and masked structs.
	initialJSON, _ := json.Marshal(eg)
	maskedJSON, _ := json.Marshal(structmask.StructMasker.Proceed(eg))

	// print.
	fmt.Println(string(initialJSON))
	fmt.Println(string(maskedJSON))

	// Output:
	// {"card_number":"0000 0000 0000 0000","card_holder":"TEST HOLDER"}
	// {"card_number":"000000******0000","card_holder":"TEST HOLDER"}
}
