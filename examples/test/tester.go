package main

import (
	"fmt"
	gonadlan "github.com/devbrain/go-nadlan"
)

func main() {
	body, err := gonadlan.GetYad2Data(1, 9000, true)
	if err != nil {
		fmt.Printf("Error %v\n", err)
	}
	fmt.Printf("Code %s\n", body.Message)
}
