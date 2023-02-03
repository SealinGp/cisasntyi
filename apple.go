package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	SearchUrl = "https://www.apple.com.cn/shop/fulfillment-messages"
)

type Apple struct {
	cli *http.Client

	configOption *ConfigOption
}

type AppleOption struct {
	ConfigOption *ConfigOption
}

func NewApple(opt *AppleOption) *Apple {

	apple := &Apple{
		cli: &http.Client{},

		configOption: opt.ConfigOption,
	}

	return apple
}

func (apple *Apple) Serve() {

	var timer *time.Timer

	for {
		if timer != nil {
			<-timer.C
		}

		apple.ReqSearch()
		timer = time.NewTimer(time.Duration(apple.configOption.SearchInterval) * time.Second)
	}
}
func (apple *Apple) ReqSearch() error {
	log.Printf("[I] 开始查询苹果接口. 查询位置:%v", apple.configOption.Location)
	appleUrl, err := apple.makeUrl()
	if err != nil {
		log.Printf("[E] make url failed. err:%v", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, appleUrl, nil)
	if err != nil {
		log.Printf("[E] new request failed. err:%v", err)
	}
	req.Header.Add("sec-fetch-site", "same-origin")

	resp, err := apple.cli.Do(req)
	if err != nil {
		log.Printf("[E] do request failed. url:%v, err:%v", appleUrl, err)
		return nil
	}

	defer resp.Body.Close()

	searchResponse, err := apple.unMarshalResp(resp)
	if err != nil {
		log.Printf("[E] unMarshalResp failed. err:%v", err)
		return nil
	}

	pickupQuote2IphoneModels := make(map[string][]string)
	stores := searchResponse.Body.Content.PickupMessage.Stores
	for _, store := range stores {
		for _, info := range store.PartsAvailability {
			msgTypes := info.MessageTypes
			pickupQuote := info.PickupSearchQuote
			pickStore := store.StoreName
			iphoneModal := msgTypes.Expanded.StorePickupProductTitle
			if iphoneModal == "" {
				iphoneModal = msgTypes.Regular.StorePickupProductTitle
			}

			log.Printf("[I] 型号:%+v 地点:%v %v", iphoneModal, pickStore, pickupQuote)
			if !apple.hasStockOffline(info.PickupSearchQuote) {
				continue
			}

			pickupContent := fmt.Sprintf("%v %v", pickStore, pickupQuote)

			iphoneModals, ok := pickupQuote2IphoneModels[pickupContent]
			if !ok {
				iphoneModals = make([]string, 0)
			}
			iphoneModals = append(iphoneModals, iphoneModal)
			pickupQuote2IphoneModels[pickupContent] = iphoneModals
		}
	}

	messages := make([]*Message, 0, 10)
	for pickupStore, iphoneModels := range pickupQuote2IphoneModels {
		if apple.configOption.NotifyMergedByStore {
			messages = append(messages, &Message{
				Title:   pickupStore,
				Content: strings.Join(iphoneModels, ","),
			})
			continue
		}

		for _, iphoneModel := range iphoneModels {
			messages = append(messages, &Message{
				Title:   iphoneModel,
				Content: pickupStore,
			})
		}
	}

	apple.sendNotificationToBarkApp(messages...)
	return nil
}

func (apple *Apple) sendNotificationToBarkApp(messages ...*Message) {

	for _, msg := range messages {

		for _, notifyUrl := range apple.configOption.NotifyUrl {
			url := fmt.Sprintf("%v/%v/%v", notifyUrl, msg.Title, msg.Content)
			_, err := apple.cli.Get(url)
			if err != nil {
				log.Printf("[E] send stock message failed. err:%v", err)
			}
		}
	}
}

func (apple *Apple) hasStockOffline(s string) bool {
	return strings.Contains(s, "可取货") && !strings.Contains(s, "不")
}

func (apple *Apple) unMarshalResp(resp *http.Response) (*SearchResponse, error) {
	searchResponse := &SearchResponse{}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[E] read body failed. err:%v", err)
		return nil, err
	}

	err = json.Unmarshal(data, searchResponse)
	if err != nil {
		log.Printf("[E] unmarshal body failed. data:%s, err:%v", data, err)
		return nil, err
	}

	return searchResponse, nil
}

func (apple *Apple) makeUrl() (string, error) {
	values := make(url.Values)

	state, city, district, err := apple.configOption.Location.GetParts()
	if err != nil {
		return "", err
	}
	values.Add("state", state)
	values.Add("city", city)
	values.Add("district", district)
	values.Add("geoLocated", "true")

	// values.Add("mt", "regular")
	// values.Add("little", "false")
	// values.Add("pl", "true")

	for i, modal := range apple.configOption.Modals {
		key := fmt.Sprintf("parts.%d", i)
		values.Add(key, modal)
	}

	query := values.Encode()
	return fmt.Sprintf("%v?%v", SearchUrl, query), nil
}
