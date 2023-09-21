package main

import (
	"fmt"
	"strings"
)

type Message struct {
	Title   string
	Content string
}

type Location string

func (l Location) GetParts() (state, city, district string, err error) {
	parts := strings.Split(string(l), " ")
	if len(parts) != 3 {
		err = fmt.Errorf("location %v format not supported, it should be like '广东 深圳 南山区'", l)
		return
	}

	state = parts[0]
	city = parts[0]
	district = parts[0]
	return
}

type ConfigOption struct {
	Modals              []string `yaml:"modals"`
	NotifyUrl           []string `yaml:"notifyUrl"`
	Location            Location `yaml:"location"`
	SearchInterval      int      `yaml:"searchInterval"`
	StoreNumber			string 	`yaml:"storeNumber"`
	NotifyMergedByStore bool     `yaml:"notifyMergedByStore"`
}

type SearchResponse struct {
	// Head SearchRespHead `json:"head,omitempty"`
	Body SearchRespBody `json:"body"`
}

type SearchRespHead struct {
	Status string                 `json:"status"`
	Data   map[string]interface{} `json:"data"`
}

type SearchRespBody struct {
	Content Content `json:"content"`
}

type Content struct {
	PickupMessage PickupMessage1 `json:"pickupMessage"`
}

type PickupMessage1 struct {
	Stores []*Store `json:"stores"`
}

type Store struct {
	StoreName         string            `json:"storeName"`
	PartsAvailability PartsAvailability `json:"partsAvailability"`
}

type PartsAvailability map[string]PartsAvailabilityValue //型号 => info

type PartsAvailabilityValue struct {
	PickupSearchQuote string        `json:"pickupSearchQuote"` //可取货
	StorePickEligible bool 			`json:"StorePickEligible"`
	MessageTypes      *MessageTypes `json:"messageTypes"`
}

type MessageTypes struct {
	Expanded Expanded `json:"expanded,omitempty"`
	Regular  Expanded `json:"regular,omitempty"`
}

type Expanded struct {
	StorePickupProductTitle string `json:"storePickupProductTitle"` //iPhone 14 Pro 128GB 暗紫色
}