package coupon

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const URL = "https://api.iijmio.jp/mobile/d/v2/coupon/"

type Body struct {
	ReturnCode string        `json:"returnCode,omitempty"`
	CouponInfo []*CouponInfo `json:"couponInfo,omitempty"`
}

type CouponInfo struct {
	HddServiceCode string     `json:"hddServiceCode,omitempty"`
	Plan           Plan       `json:"plan,omitempty"`
	HdoInfo        []*HdoInfo `json:"hdoInfo,omitempty"`
	HduInfo        []*HduInfo `json:"hduInfo,omitempty"`
	HdxInfo        []*HdxInfo `json:"hdxInfo,omitempty"`
	Coupon         []*Coupon  `json:"coupon,omitempty"`
	History        []*History `json:"history,omitempty"`
	Remains        int        `json:"remains,omitempty"`
}

type Plan string

const (
	PlanFamilyShare  Plan = "Family Share"
	PlanMinimumStart Plan = "Light Start"
	PlanEcoMinimum   Plan = "Eco Minimum"
	PlanEcoStandard  Plan = "Eco Standard"
)

// HdoInfo ドコモのSIM
type HdoInfo struct {
	HdoServiceCode string    `json:"hdoServiceCode,omitempty"`
	Number         string    `json:"number,omitempty"`
	Iccid          string    `json:"iccid,omitempty"`
	Regulation     bool      `json:"regulation,omitempty"`
	Sms            bool      `json:"sms,omitempty"`
	Voice          bool      `json:"voice,omitempty"`
	CouponUse      bool      `json:"couponUse"`
	Coupon         []*Coupon `json:"coupon,omitempty"`
}

// HduInfo auのSIM
type HduInfo struct {
	HduServiceCode string    `json:"hduServiceCode,omitempty"`
	Number         string    `json:"number,omitempty"`
	Iccid          string    `json:"iccid,omitempty"`
	Regulation     bool      `json:"regulation,omitempty"`
	Sms            bool      `json:"sms,omitempty"`
	Voice          bool      `json:"voice,omitempty"`
	CouponUse      bool      `json:"couponUse,omitempty"`
	Coupon         []*Coupon `json:"coupon,omitempty"`
}

// HdxInfo eSIM
type HdxInfo struct {
	HdxServiceCode string    `json:"hdxServiceCode,omitempty"`
	Number         string    `json:"number,omitempty"`
	Iccid          string    `json:"iccid,omitempty"`
	Regulation     bool      `json:"regulation,omitempty"`
	Sms            bool      `json:"sms,omitempty"`
	Voice          bool      `json:"voice,omitempty"`
	CouponUse      bool      `json:"couponUse,omitempty"`
	Coupon         []*Coupon `json:"coupon,omitempty"`
}

type Coupon struct {
	Volume int    `json:"volume,omitempty"`
	Expire string `json:"expire,omitempty"`
	Type   Type   `json:"type,omitempty"`
}

type Type string

const (
	TypeBundle Type = "bundle"
	TypeTopUp  Type = "topup"
	TypeSim    Type = "sim"
)

type History struct {
	Date   string `json:"date,omitempty"`
	Event  string `json:"event,omitempty"`
	Volume int    `json:"volume,omitempty"`
	Type   Type   `json:"type,omitempty"`
}

func Get(developerID, accessToken string) (*Body, int, error) {
	req, _ := http.NewRequest(http.MethodGet, URL, nil)
	req.Header.Set("X-IIJmio-Developer", developerID)
	req.Header.Set("X-IIJmio-Authorization", accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	res := &Body{}
	if err := json.Unmarshal(b, res); err != nil {
		return nil, resp.StatusCode, err
	}
	return res, resp.StatusCode, nil
}

func Put(developerID, accessToken string, body *Body) (*Body, int, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, 0, err
	}
	r := bytes.NewReader(b)
	req, _ := http.NewRequest(http.MethodPut, URL, r)
	req.Header.Set("X-IIJmio-Developer", developerID)
	req.Header.Set("X-IIJmio-Authorization", accessToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	res := &Body{}
	if err := json.Unmarshal(b, res); err != nil {
		return nil, resp.StatusCode, err
	}
	return res, resp.StatusCode, nil
}
