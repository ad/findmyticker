package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"
)

type HAItem struct {
	DevID       string     `json:"dev_id"`
	Gps         [2]float64 `json:"gps"`
	GpsAccuracy float64    `json:"gps_accuracy"`
	HostName    string     `json:"host_name"`
	Battery     float64    `json:"battery"`
	// LocationName string     `json:"location_name"`
}

func sendItemsToHomeAssistant(items *Items) {
	for _, item := range *items {
		if config.Update.Ignore != nil && slices.Contains(config.Update.Ignore, item.Identifier) {
			continue
		}

		if !item.Location.LocationFinished {
			continue
		}

		oldValue, ok := lruCache.Get(item.Identifier)
		if ok {
			if oldValue[0] == item.Location.Latitude && oldValue[1] == item.Location.Longitude {
				continue
			}
		}

		lruCache.Set(item.Identifier, [2]float64{item.Location.Latitude, item.Location.Longitude})

		if menuInfo != nil {
			menuInfo.SetTitle(fmt.Sprintf("last update: %s", time.Now().Format("15:04:05")))
		}

		batteryLevel := 0.00

		if item.BatteryStatus > 0 {
			batteryLevel = 100 / float64(item.BatteryStatus)
		}

		haItem := HAItem{
			DevID: fmt.Sprintf("findmy_%s", strings.Replace(item.Identifier, "-", "", -1)),
			Gps: [2]float64{
				item.Location.Latitude,
				item.Location.Longitude,
			},
			GpsAccuracy: item.Location.HorizontalAccuracy,
			// LocationName: item.Address.MapItemFullAddress,
			HostName: item.Name,
			Battery:  batteryLevel,
		}

		_ = processHomeassistant(haItem)
	}
}

func sendDevicesToHomeAssistant(devices *Devices) {
	for _, item := range *devices {
		if config.Update.Ignore != nil && slices.Contains(config.Update.Ignore, item.Identifier) {
			continue
		}

		if !item.Location.LocationFinished {
			continue
		}

		if item.Location.HorizontalAccuracy > float64(config.Update.MinimalAccuracy) {
			continue
		}

		oldValue, ok := lruCache.Get(item.Identifier)
		if ok {
			if oldValue[0] == item.Location.Latitude && oldValue[1] == item.Location.Longitude {
				continue
			}
		}

		lruCache.Set(item.Identifier, [2]float64{item.Location.Latitude, item.Location.Longitude})

		if menuInfo != nil {
			menuInfo.SetTitle(fmt.Sprintf("last update: %s", time.Now().Format("15:04:05")))
		}

		batteryLevel := 0.00

		if item.BatteryStatus > 0 {
			batteryLevel = float64(item.BatteryStatus) * 100
		}

		haItem := HAItem{
			DevID: fmt.Sprintf("findmy_%s", strings.Replace(item.Identifier, "-", "", -1)),
			Gps: [2]float64{
				item.Location.Latitude,
				item.Location.Longitude,
			},
			GpsAccuracy: item.Location.HorizontalAccuracy,
			// LocationName: item.Address.MapItemFullAddress,
			HostName: item.Name,
			Battery:  batteryLevel,
		}

		_ = processHomeassistant(haItem)
	}
}

func processHomeassistant(haItem HAItem) error {
	jsonStr, errMarshal := json.Marshal(haItem)
	if errMarshal != nil {
		return errMarshal
	}

	url := config.Homeassistant.URL

	req, errNewRequest := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if errNewRequest != nil {
		return errNewRequest
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.Homeassistant.Token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, errDo := client.Do(req)
	if errDo != nil {
		return errDo
	}
	defer resp.Body.Close()

	fmt.Println("request body:", string(jsonStr))
	if resp.StatusCode != 200 {
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	}

	return nil
}
