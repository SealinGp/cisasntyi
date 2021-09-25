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
func (apple *Apple) ReqSearch() {
	appleUrl := apple.makeUrl()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, appleUrl, nil)
	if err != nil {
		log.Printf("[E] new request failed. err:%v", err)
	}

	resp, err := apple.cli.Do(req)
	if err != nil {
		log.Printf("[E] do request failed. url:%v, err:%v", appleUrl, err)
		return
	}

	defer resp.Body.Close()

	searchResponse, err := apple.unMarshalResp(resp)
	if err != nil {
		log.Printf("[E] unMarshalResp failed. err:%v", err)
		return
	}

	messages := make([]*Message, 0, 10)
	stores := searchResponse.Body.Content.PickupMessage.Stores
	for _, store := range stores {

		for _, info := range store.PartsAvailability {
			if !apple.hasStockOffline(info.PickupSearchQuote) {
				continue
			}

			iphoneModal := info.StorePickupProductTitle
			pickTime := info.PickupSearchQuote
			pickStore := store.StoreName

			msg := &Message{
				Title:   iphoneModal,
				Content: fmt.Sprintf("取货时间:%v 地点:%v", pickTime, pickStore),
			}

			messages = append(messages, msg)
		}
	}

	for _, msg := range messages {
		for _, notifyUrl := range apple.configOption.NotifyUrl {
			url := fmt.Sprintf("%v/%v/%v", notifyUrl, msg.Title, msg.Content)
			_, err = apple.cli.Get(url)
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

func (apple *Apple) makeUrl() string {
	values := make(url.Values)
	values.Add("mt", "regular")
	values.Add("little", "false")
	values.Add("pl", "true")
	values.Add("location", apple.configOption.Location)

	for i, modal := range apple.configOption.Modals {
		key := fmt.Sprintf("parts.%d", i)
		values.Add(key, modal)
	}

	query := values.Encode()
	return fmt.Sprintf("%v?%v", SearchUrl, query)
}
