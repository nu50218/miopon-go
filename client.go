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
	mutex       *Mutex
	settings    *Settings
}

// Mutex API、アクセストークン種別でスリープさせるもの
type Mutex struct {
	GetCoupon      map[string]*sync.Mutex
	GetCouponMutex *sync.Mutex
	PutCoupon      map[string]*sync.Mutex
	PutCouponMutex *sync.Mutex
	GetPacket      map[string]*sync.Mutex
	GetPacketMutex *sync.Mutex
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
		mutex: &Mutex{
			GetCoupon:      map[string]*sync.Mutex{},
			GetCouponMutex: &sync.Mutex{},
			PutCoupon:      map[string]*sync.Mutex{},
			PutCouponMutex: &sync.Mutex{},
			GetPacket:      map[string]*sync.Mutex{},
			GetPacketMutex: &sync.Mutex{},
		},
		settings: settings,
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
	client.mutex.GetCouponMutex.Lock()
	if client.mutex.GetCoupon[accessToken] == nil {
		client.mutex.GetCoupon[accessToken] = &sync.Mutex{}
	}
	client.mutex.GetCouponMutex.Unlock()

	client.mutex.GetCoupon[accessToken].Lock()
	defer func(token string) {
		go func(token string) {
			time.Sleep(client.settings.GetCouponInterval)
			client.mutex.GetCoupon[token].Unlock()
		}(token)
	}(accessToken)

	return coupon.Get(client.DeveloperID, accessToken)
}

// PutCoupon クーポンのON/OFF
// ステータスコードも返す
func (client *Client) PutCoupon(accessToken string, couponInfo []*coupon.CouponInfo) (*coupon.Body, int, error) {
	client.mutex.PutCouponMutex.Lock()
	if client.mutex.PutCoupon[accessToken] == nil {
		client.mutex.PutCoupon[accessToken] = &sync.Mutex{}
	}
	client.mutex.PutCouponMutex.Unlock()

	client.mutex.PutCoupon[accessToken].Lock()
	defer func(token string) {
		go func(token string) {
			time.Sleep(client.settings.PutCouponInterval)
			client.mutex.PutCoupon[token].Unlock()
		}(token)
	}(accessToken)

	return coupon.Put(
		client.DeveloperID,
		accessToken,
		&coupon.Body{
			CouponInfo: couponInfo,
		},
	)
}

// GetPacket データ利用量照会
// ステータスコードも返す
func (client *Client) GetPacket(accessToken string) (*packet.Body, int, error) {
	client.mutex.GetPacketMutex.Lock()
	if client.mutex.GetPacket[accessToken] == nil {
		client.mutex.GetPacket[accessToken] = &sync.Mutex{}
	}
	client.mutex.GetPacketMutex.Unlock()

	client.mutex.GetPacket[accessToken].Lock()
	defer func(token string) {
		go func(token string) {
			time.Sleep(client.settings.GetPacketInterval)
			client.mutex.GetPacket[token].Unlock()
		}(token)
	}(accessToken)

	return packet.Get(client.DeveloperID, accessToken)
}
