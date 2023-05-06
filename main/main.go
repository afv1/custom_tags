package main

import (
	"encoding/json"
	"fmt"
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

	//structMaskConfig := confidential.Config{
	//	Tags:     nil,
	//	Mappers:  nil,
	//	Handlers: nil,
	//}
	//
	//
	//masked := confidential.Proceed(tst)
	//
	//jsn, _ := json.Marshal(masked)
	//fmt.Println(string(jsn))
}
