package main

import (
	"fmt"
	"github.com/caffix/amass/amass"
)

func main() {
	output := make(chan *amass.AmassOutput)
	go func() {
		for result := range output {
			fmt.Println(result.Name)
		}
	}()
	// Setup the most basic amass configuration
	config := amass.CustomConfig(&amass.AmassConfig{Output: output})
	config.AddDomains([]string{"example.com"})
	// Begin the enumeration process
	amass.StartEnumeration(config)
}
