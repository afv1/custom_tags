package main

import (
	"encoding/json"
	"fmt"
	"sexreflection/pkg/confidential"
	"sexreflection/pkg/mask"
)

type Test struct {
	Text     string `json:"text"`
	TestMask string `json:"test_mask" confidential:"default"`
	Card     Card
}

type Card struct {
	Number      string `json:"card_number" confidential:"cardnumber"`
	Holder      string `json:"card_holder" confidential:"cardholder"`
	ExpireMonth string `json:"expire_month"`
	ExpireYear  string `json:"expire_year"`
	CVV         string `json:"cvv" confidential:"cvv"`
}

func main() {
	tst := Test{
		Text:     "test field1",
		TestMask: "test string",
		Card: Card{
			Number:      "0000 0000 0000 0000",
			Holder:      "TEST DEV",
			ExpireMonth: "05",
			ExpireYear:  "25",
			CVV:         "123",
		},
	}

	j, _ := json.Marshal(tst)

	fmt.Println(string(j))

	structMaskMappers := confidential.Mapper{
		"cvv":        mask.CVV,
		"cardnumber": mask.CardNumber,
		"cardholder": mask.CardHolder,
		"def":        mask.Default,
		"ip":         mask.IP,
	}

	confidential.InitStructMask(&confidential.Config{
		TagName: "confidential",
		Mappers: structMaskMappers,
	})

	masked := confidential.StructMasker.Proceed(tst)

	jsn, _ := json.Marshal(masked)
	fmt.Println(string(jsn))
}
