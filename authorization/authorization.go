package authorization

import (
	"net/url"
	"strconv"
)

const URL = "https://api.iijmio.jp/mobile/d/v1/authorization/"

// リクエスト
const (
	ParameterResponseType = "response_type"
	ParameterClientID     = "client_id"
	ParameterRedirectURI  = "redirect_uri"
	ParameterState        = "state"
)

// レスポンス
const (
	ParameterAccessToken = "access_token"
	ParameterTokenType   = "token_type"
	ParameterExpiresIn   = "expires_in"
)

// エラー
const (
	ParameterError            = "error"
	ParameterErrorDescription = "error_description"
)

func IsSuccessResponse(uri string) (bool, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return false, err
	}
	q, err := url.ParseQuery(u.Fragment)
	if err != nil {
		return false, err
	}
	return q.Get(ParameterAccessToken) != "", nil
}

type SuccessResponse struct {
	AccessToken string
	TokenType   string
	ExpiresIn   int64
	State       string
}

func ParseSuccessResponse(uri string) (*SuccessResponse, error) {
	var a int
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	q, err := url.ParseQuery(u.Fragment)
	if err != nil {
		return nil, err
	}
	res := &SuccessResponse{
		AccessToken: q.Get(ParameterAccessToken),
		TokenType:   q.Get(ParameterTokenType),
		State:       q.Get(ParameterState),
	}
	res.ExpiresIn, err = strconv.ParseInt(q.Get(ParameterExpiresIn), 10, 64)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type ErrorResponse struct {
	Error            string
	ErrorDescription string
	State            string
}

func ParseErrorResponse(uri string) (*ErrorResponse, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	q, err := url.ParseQuery(u.Fragment)
	if err != nil {
		return nil, err
	}
	res := &ErrorResponse{
		Error:            q.Get(ParameterError),
		ErrorDescription: q.Get(ParameterErrorDescription),
		State:            q.Get(ParameterState),
	}
	return res, nil
}
