package packet

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const URL = "https://api.iijmio.jp/mobile/d/v2/log/packet/"

type Body struct {
	ReturnCode    string           `json:"returnCode,omitempty"`
	PacketLogInfo []*PacketLogInfo `json:"packetLogInfo,omitempty"`
}

type PacketLogInfo struct {
	HddServiceCode string     `json:"hddServiceCode,omitempty"`
	Plan           string     `json:"plan,omitempty"`
	HdoInfo        []*HdoInfo `json:"hdoInfo,omitempty"`
	HduInfo        []*HduInfo `json:"hduInfo,omitempty"`
	HdxInfo        []*HdxInfo `json:"hdxInfo,omitempty"`
}

type HdoInfo struct {
	HdoServiceCode string       `json:"hdoServiceCode,omitempty"`
	PacketLog      []*PacketLog `json:"packetLog,omitempty"`
}

type HduInfo struct {
	HduServiceCode string       `json:"hduServiceCode,omitempty"`
	PacketLog      []*PacketLog `json:"packetLog,omitempty"`
}

type HdxInfo struct {
	HdxServiceCode string       `json:"hdxServiceCode,omitempty"`
	PacketLog      []*PacketLog `json:"packetLog,omitempty"`
}

type PacketLog struct {
	Date          string `json:"date,omitempty"`
	WithCoupon    int    `json:"withCoupon,omitempty"`
	WithoutCoupon int    `json:"withoutCoupon,omitempty"`
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
