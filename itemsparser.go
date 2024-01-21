package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"time"
)

var (
	lastModTime time.Time
)

type Items []struct {
	Identifier string `json:"identifier,omitempty"`
	Name       string `json:"name,omitempty"`
	Location   struct {
		LocationFinished   bool    `json:"locationFinished,omitempty"`
		HorizontalAccuracy float64 `json:"horizontalAccuracy,omitempty"`
		Latitude           float64 `json:"latitude,omitempty"`
		Longitude          float64 `json:"longitude,omitempty"`
		// IsOld              bool    `json:"isOld,omitempty"`
		// Altitude           int     `json:"altitude,omitempty"`
		// FloorLevel         int     `json:"floorLevel,omitempty"`
		// PositionType       string  `json:"positionType,omitempty"`
		// IsInaccurate       bool    `json:"isInaccurate,omitempty"`
		// TimeStamp          int64   `json:"timeStamp,omitempty"`
		// VerticalAccuracy   int     `json:"verticalAccuracy,omitempty"`
	} `json:"location,omitempty"`
	BatteryStatus int `json:"batteryStatus,omitempty"`
	// Address    struct {
	// SubAdministrativeArea interface{}   `json:"subAdministrativeArea,omitempty"`
	// FullThroroughfare     string        `json:"fullThroroughfare,omitempty"`
	// StateCode             interface{}   `json:"stateCode,omitempty"`
	// StreetAddress         string        `json:"streetAddress,omitempty"`
	// Country               string        `json:"country,omitempty"`
	// Label                 string        `json:"label,omitempty"`
	// AreaOfInterest        []interface{} `json:"areaOfInterest,omitempty"`
	// FormattedAddressLines []string      `json:"formattedAddressLines,omitempty"`
	// MapItemFullAddress string `json:"mapItemFullAddress,omitempty"`
	// CountryCode           string        `json:"countryCode,omitempty"`
	// AdministrativeArea    string        `json:"administrativeArea,omitempty"`
	// Locality              string        `json:"locality,omitempty"`
	// StreetName            string        `json:"streetName,omitempty"`
	// } `json:"address,omitempty"`
	// IsAppleAudioAccessory     bool `json:"isAppleAudioAccessory,omitempty"`
	// IsFirmwareUpdateMandatory bool `json:"isFirmwareUpdateMandatory,omitempty"`
	// ProductType               struct {
	// 	ProductInformation struct {
	// 		RequiresAudioSafetyAlert         bool   `json:"requiresAudioSafetyAlert,omitempty"`
	// 		ModelName                        string `json:"modelName,omitempty"`
	// 		DefaultHeroIcon2X                string `json:"defaultHeroIcon2x,omitempty"`
	// 		DefaultListIcon3X                string `json:"defaultListIcon3x,omitempty"`
	// 		ManufacturerName                 string `json:"manufacturerName,omitempty"`
	// 		AppBundleIdentifier              string `json:"appBundleIdentifier,omitempty"`
	// 		DefaultHeroIcon                  string `json:"defaultHeroIcon,omitempty"`
	// 		DefaultListIcon                  string `json:"defaultListIcon,omitempty"`
	// 		RequiresAdditionalConnectionTime bool   `json:"requiresAdditionalConnectionTime,omitempty"`
	// 		DefaultListIcon2X                string `json:"defaultListIcon2x,omitempty"`
	// 		ProductIdentifier                int    `json:"productIdentifier,omitempty"`
	// 		VendorIdentifier                 int    `json:"vendorIdentifier,omitempty"`
	// 		DefaultHeroIcon3X                string `json:"defaultHeroIcon3x,omitempty"`
	// 		AntennaPower                     int    `json:"antennaPower,omitempty"`
	// 	} `json:"productInformation,omitempty"`
	// 	Type string `json:"type,omitempty"`
	// } `json:"productType,omitempty"`
	// LostModeMetadata interface{} `json:"lostModeMetadata,omitempty"`
	// SerialNumber     string      `json:"serialNumber,omitempty"`
	// SystemVersion    string      `json:"systemVersion,omitempty"`
	// Capabilities     int         `json:"capabilities,omitempty"`
	// GroupIdentifier  interface{} `json:"groupIdentifier,omitempty"`
	// Role             struct {
	// 	Name       string `json:"name,omitempty"`
	// 	Identifier int    `json:"identifier,omitempty"`
	// 	Emoji      string `json:"emoji,omitempty"`
	// } `json:"role,omitempty"`
	// CrowdSourcedLocation struct {
	// 	IsInaccurate       bool    `json:"isInaccurate,omitempty"`
	// 	Altitude           int     `json:"altitude,omitempty"`
	// 	LocationFinished   bool    `json:"locationFinished,omitempty"`
	// 	PositionType       string  `json:"positionType,omitempty"`
	// 	IsOld              bool    `json:"isOld,omitempty"`
	// 	FloorLevel         int     `json:"floorLevel,omitempty"`
	// 	HorizontalAccuracy float64 `json:"horizontalAccuracy,omitempty"`
	// 	Longitude          float64 `json:"longitude,omitempty"`
	// 	TimeStamp          int64   `json:"timeStamp,omitempty"`
	// 	VerticalAccuracy   int     `json:"verticalAccuracy,omitempty"`
	// 	Latitude           float64 `json:"latitude,omitempty"`
	// } `json:"crowdSourcedLocation,omitempty"`
	// ProductIdentifier string `json:"productIdentifier,omitempty"`
	// SafeLocations     []struct {
	// 	Location struct {
	// 		TimeStamp          int64   `json:"timeStamp,omitempty"`
	// 		LocationFinished   bool    `json:"locationFinished,omitempty"`
	// 		IsInaccurate       bool    `json:"isInaccurate,omitempty"`
	// 		PositionType       string  `json:"positionType,omitempty"`
	// 		FloorLevel         int     `json:"floorLevel,omitempty"`
	// 		HorizontalAccuracy float64 `json:"horizontalAccuracy,omitempty"`
	// 		Altitude           int     `json:"altitude,omitempty"`
	// 		IsOld              bool    `json:"isOld,omitempty"`
	// 		VerticalAccuracy   float64 `json:"verticalAccuracy,omitempty"`
	// 		Latitude           float64 `json:"latitude,omitempty"`
	// 		Longitude          float64 `json:"longitude,omitempty"`
	// 	} `json:"location,omitempty"`
	// 	Identifier string `json:"identifier,omitempty"`
	// 	Name       string `json:"name,omitempty"`
	// 	Address    struct {
	// 		MapItemFullAddress    string        `json:"mapItemFullAddress,omitempty"`
	// 		CountryCode           string        `json:"countryCode,omitempty"`
	// 		FullThroroughfare     string        `json:"fullThroroughfare,omitempty"`
	// 		Locality              string        `json:"locality,omitempty"`
	// 		StreetAddress         interface{}   `json:"streetAddress,omitempty"`
	// 		StreetName            string        `json:"streetName,omitempty"`
	// 		AreaOfInterest        []interface{} `json:"areaOfInterest,omitempty"`
	// 		AdministrativeArea    string        `json:"administrativeArea,omitempty"`
	// 		SubAdministrativeArea interface{}   `json:"subAdministrativeArea,omitempty"`
	// 		StateCode             interface{}   `json:"stateCode,omitempty"`
	// 		Country               string        `json:"country,omitempty"`
	// 		Label                 string        `json:"label,omitempty"`
	// 		FormattedAddressLines []string      `json:"formattedAddressLines,omitempty"`
	// 	} `json:"address,omitempty"`
	// 	Type          int `json:"type,omitempty"`
	// 	ApprovalState int `json:"approvalState,omitempty"`
	// } `json:"safeLocations,omitempty"`
	// PartInfo interface{} `json:"partInfo,omitempty"`
	// Owner    string      `json:"owner,omitempty"`
}

func ParseItems() (Items, error) {
	u, err := user.Lookup(os.ExpandEnv("$USER"))
	if err != nil {
		return nil, fmt.Errorf("user lookup: %+v", err)
	}

	itemsPath := fmt.Sprintf("%s/Library/Caches/com.apple.findmy.fmipcore/Items.data", u.HomeDir)

	stat, errStat := os.Stat(itemsPath)
	if errStat != nil {
		return nil, fmt.Errorf("stat: %+v", errStat)
	}

	if lastModTime == stat.ModTime() {
		return nil, nil
	}

	lastModTime = stat.ModTime()

	// fmt.Println(itemsPath, "updated at", lastModTime)

	// Open items file
	f, errReadFile := os.ReadFile(itemsPath)
	if errReadFile != nil {
		return nil, fmt.Errorf("read file: %+v", errReadFile)
	}

	items := &Items{}

	// Parse items
	errUnmarshal := json.Unmarshal(f, &items)
	if errUnmarshal != nil {
		return nil, fmt.Errorf("unmarshal: %+v, %+v", errUnmarshal, f)
	}

	return *items, nil
}
