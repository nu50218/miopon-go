package miopon

import (
	"net/url"
	"sync"
	"time"

	"github.com/nu50218/miopon-go/packet"

	"github.com/nu50218/miopon-go/coupon"

	"github.com/nu50218/miopon-go/authorization"
)

// Client IIJmioクーポンスイッチAPIのクライアント
type Client struct {
	DeveloperID string
	mutex       Mutex
	settings    *Settings
}

// Mutex API、アクセストークン種別でスリープさせるもの
type Mutex struct {
	GetCoupon map[string]*sync.Mutex
	PutCoupon map[string]*sync.Mutex
	GetPacket map[string]*sync.Mutex
}

// Settings 設定
type Settings struct {
	GetCouponInterval time.Duration
	PutCouponInterval time.Duration
	GetPacketInterval time.Duration
}

// New 新しいClientを作成して返す
func New(developerID string, settings *Settings) *Client {
	return &Client{
		DeveloperID: developerID,
		settings:    settings,
	}
}

// MakeAuthorizationURL アクセストークンを取得するためのURLを作成する
func (client *Client) MakeAuthorizationURL(redirectURI, state string) string {
	u, _ := url.Parse(authorization.URL)
	q := u.Query()
	q.Add(authorization.ParameterResponseType, "token")
	q.Add(authorization.ParameterClientID, client.DeveloperID)
	q.Add(authorization.ParameterRedirectURI, redirectURI)
	q.Add(authorization.ParameterState, state)
	u.RawQuery = q.Encode()
	return u.String()
}

// GetCoupon クーポン残量照会・クーポンのON/OFF状態照会
// ステータスコードも返す
func (client *Client) GetCoupon(accessToken string) (*coupon.Body, int, error) {
	client.mutex.GetCoupon[accessToken].Lock()
	defer func(token string) {
		time.Sleep(client.settings.GetCouponInterval)
		client.mutex.GetCoupon[token].Unlock()
	}(accessToken)

	return coupon.Get(client.DeveloperID, accessToken)
}

// PutCoupon クーポンのON/OFF
// ステータスコードも返す
func (client *Client) PutCoupon(accessToken string, hdoInfo []*coupon.HdoInfo, hduInfo []*coupon.HduInfo, hdxInfo []*coupon.HdxInfo) (*coupon.Body, int, error) {
	client.mutex.PutCoupon[accessToken].Lock()
	defer func(token string) {
		time.Sleep(client.settings.PutCouponInterval)
		client.mutex.PutCoupon[token].Unlock()
	}(accessToken)

	return coupon.Put(
		client.DeveloperID,
		accessToken,
		&coupon.Body{
			CouponInfo: []*coupon.CouponInfo{
				&coupon.CouponInfo{
					HdoInfo: hdoInfo,
					HduInfo: hduInfo,
					HdxInfo: hdxInfo,
				},
			},
		},
	)
}

// GetPacket データ利用量照会
// ステータスコードも返す
func (client *Client) GetPacket(accessToken string) (*packet.Body, int, error) {
	client.mutex.GetPacket[accessToken].Lock()
	defer func(token string) {
		time.Sleep(client.settings.GetPacketInterval)
		client.mutex.GetPacket[token].Unlock()
	}(accessToken)

	return packet.Get(client.DeveloperID, accessToken)
}
