package main

import (
	"fmt"
)

func run() {
	if menuError != nil {
		menuError.Hide()
	}

	items, errParseItems := ParseItems()
	if errParseItems != nil {
		fmt.Printf("Error: %+v\n", errParseItems)

		if menuError != nil {
			menuError.SetTitle(fmt.Sprintf("error: %s", errParseItems.Error()))
			menuError.Show()
		}

		return
	}

	sendToHomeAssistant(&items)
}
