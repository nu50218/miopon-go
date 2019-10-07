package miopon

import (
	"net/url"

	"github.com/nu50218/miopon-go/packet"

	"github.com/nu50218/miopon-go/coupon"

	"github.com/nu50218/miopon-go/authorization"
)

// Client IIJmioクーポンスイッチAPIのクライアント
type Client struct {
	DeveloperID string
}

// New 新しいClientを作成して返す
func New(developerID string) *Client {
	return &Client{
		DeveloperID: developerID,
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
	return coupon.Get(client.DeveloperID, accessToken)
}

// PutCoupon クーポンのON/OFF
// ステータスコードも返す
func (client *Client) PutCoupon(accessToken string, hdoInfo []*coupon.HdoInfo, hduInfo []*coupon.HduInfo, hdxInfo []*coupon.HdxInfo) (*coupon.Body, int, error) {
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
	return packet.Get(client.DeveloperID, accessToken)
}
