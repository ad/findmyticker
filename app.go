package main

import (
	"fmt"
)

func run() {
	if menuError != nil {
		menuError.Hide()
	}

	if config.AllowItems {
		items, errParseItems := ParseItems()
		if errParseItems != nil {
			fmt.Printf("Error: %+v\n", errParseItems)

			if menuError != nil {
				menuError.SetTitle(fmt.Sprintf("error: %s", errParseItems.Error()))
				menuError.Show()
			}
		}

		if items != nil {
			sendItemsToHomeAssistant(&items)
		}
	}

	if config.AllowDevices {
		devices, errParseItems := ParseDevices()
		if errParseItems != nil {
			fmt.Printf("Error: %+v\n", errParseItems)

			if menuError != nil {
				menuError.SetTitle(fmt.Sprintf("error: %s", errParseItems.Error()))
				menuError.Show()
			}

			return
		}

		sendDevicesToHomeAssistant(&devices)
	}
}
