package main

import (
	"fmt"
)

func run() {
	items, errParseItems := ParseItems()
	if errParseItems != nil {
		fmt.Printf("Error: %+v\n", errParseItems)

		if menuInfo != nil {
			menuInfo.SetTitle(fmt.Sprintf("error: %s", errParseItems.Error()))
		}

		return
	}

	sendToHomeAssistant(&items)
}
