# structmasker

Simple Golang package for Custom Tags.
Allows You to bind custom funcions to custom tags and then proceed them.

# Getting started

``` bash
$ go get -u github.com/afv1/custom_tags
```

# How to use

Example:
``` golang
package main

import (
	"encoding/json"
	"fmt"
	"strings"
	
	"github.com/afv1/custom_tags"
)

// EGStruct define maskable struct.
type EGStruct struct {
    Number      string `json:"card_number" mask:"cardnumber"`
    Holder      string `json:"card_holder"`
    ExpireMonth string `json:"expire_month"`
    ExpireYear  string `json:"expire_year"`
    CVV         string `json:"cvv" mask:"cvv"`
}

// CardMaskHandler must implement structmask.Handler.
func CardMaskHandler(input string) string {
    input = strings.ReplaceAll(input, " ", "")

    const (
        binLn = 6
        tailLn = 4
        symbol = "*"
    )
    
    maskLn := len(input) - binLn - tailLn
    
    exp := fmt.Sprintf(`(\d{%d}).*(\d{%d})`, binLn, tailLn)
    repl := "$1" + strings.Repeat(symbol, maskLn) + "$2"
    
    reg, _ := regexp.Compile(exp)
    
    return string(reg.ReplaceAll([]byte(input), []byte(repl)))
}

// CVVMaskHandler must implement structmask.Handler.
func CVVMaskHandler(_ string) string {
	return "***"
}

func main() {
    // create struct instance.
    egStruct := EGStruct{
        Number:      "0000 0000 0000 0000",
        Holder:      "TEST DEV",
        ExpireMonth: "05",
        ExpireYear:  "25",
        CVV:         "123",
    }
    
    // init Custom Tags with tag name.
    ct := customtags.NewCustomTags("mask")
	customtags.Bind("test", Test1)
	
	// bind handlers to custom tag labels.
	customtags.Bind("cardnumber", CardMaskHandler)
	customtags.Bind("cardholder", CardHolderHandler)

	// print initial marshaled struct.
	initialJSON, _ := json.Marshal(egStruct)
	fmt.Println(string(initialJSON))

	// print masked marshaled struct.
	modifiedJSON := ct.Proceed(egStruct)
	jsn, _ := json.Marshal(modifiedJSON)
	fmt.Println(string(jsn))
}
```

Output:

```
{"card_number":"0000 0000 0000 0000","card_holder":"TEST DEV","expire_month":"05","expire_year":"25","cvv":"123"}
{"card_number":"000000******0000","card_holder":"TEST DEV","expire_month":"05","expire_year":"25","cvv":"***"}
```
