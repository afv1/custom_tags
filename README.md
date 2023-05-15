# structmasker

Simple Golang Struct fields masking pkg with customizable mask handlers for every tag label

# Getting started

``` bash
$ go get -u github.com/afv1/structmasker
```

# How to use

Example:
``` golang
package main

import (
	"cardmasker/structmask"
	"encoding/json"
	"fmt"
	"strings"
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
    
    // define mappers.
    structMaskMappers := structmask.Mapper{
        "cvv":        CVVMaskHandler,
        "cardnumber": CardMaskHandler,
    }
    
    // init StructMask with config.
    structmask.InitStructMask(&structmask.Config{
        TagName: "mask",
        Mappers: structMaskMappers,
    })
    
    // print initial marshaled struct.
    initialJSON, _ := json.Marshal(egStruct)
    fmt.Println(string(initialJSON))
    
    // print masked marshaled struct.
    maskedJSON := structmask.StructMasker.Proceed(egStruct)
    jsn, _ := json.Marshal(maskedJSON)
    fmt.Println(string(jsn))
}
```

Output:

```
{"card_number":"0000 0000 0000 0000","card_holder":"TEST DEV","expire_month":"05","expire_year":"25","cvv":"123"}
{"card_number":"000000******0000","card_holder":"TEST DEV","expire_month":"05","expire_year":"25","cvv":"***"}
```
