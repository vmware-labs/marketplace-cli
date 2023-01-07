// Copyright 2023 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package csp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
)

//go:generate counterfeiter . TokenParserFn
type TokenParserFn func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error)

type TokenServices struct {
	CSPHost     string
	Client      pkg.HTTPClient
	TokenParser TokenParserFn
}

type RedeemResponse struct {
	AccessToken  string      `json:"access_token"`
	StatusCode   int         `json:"statusCode,omitempty"`
	ModuleCode   int         `json:"moduleCode,omitempty"`
	Metadata     interface{} `json:"metadata,omitempty"` // I don't know what the appropriate type for this field is
	TraceID      string      `json:"traceId,omitempty"`
	CSPErrorCode string      `json:"cspErrorCode,omitempty"`
	Message      string      `json:"message,omitempty"`
	RequestID    string      `json:"requestId,omitempty"`
}

func (csp *TokenServices) Redeem(refreshToken string) (*Claims, error) {
	requestURL := pkg.MakeURL(csp.CSPHost, "/csp/gateway/am/api/auth/api-tokens/authorize", nil)
	formData := url.Values{
		"refresh_token": []string{refreshToken},
	}

	resp, err := csp.Client.PostForm(requestURL, formData)
	if err != nil {
		return nil, fmt.Errorf("failed to redeem token: %w", err)
	}

	var body RedeemResponse
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redeem response: %w", err)
	}

	if resp.StatusCode == http.StatusBadRequest && strings.Contains(body.Message, "invalid_grant: Invalid refresh token") {
		return nil, errors.New("the CSP API token is invalid or expired")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to exchange refresh token for access token: %s: %s", resp.Status, body.Message)
	}

	claims := &Claims{}
	token, err := csp.TokenParser(body.AccessToken, claims, csp.GetPublicKey)
	if err != nil {
		return nil, fmt.Errorf("invalid token returned from CSP: %w", err)
	}

	claims.Token = token.Raw
	return claims, nil
}

func (csp *TokenServices) GetPublicKey(*jwt.Token) (interface{}, error) {
	resp, err := csp.Client.Get(pkg.MakeURL(csp.CSPHost, "/csp/gateway/am/api/auth/token-public-key", nil))
	if err != nil {
		return nil, fmt.Errorf("failed to get CSP Public key: %w", err)
	}

	m := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSP Public key: %w", err)
	}

	pemData, ok := m["value"]
	if !ok {
		return nil, fmt.Errorf("public key does not contain value")
	}

	s, ok := pemData.(string)
	if !ok {
		return nil, fmt.Errorf("public key value was not in the expected format")
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(s))
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSP public key: %w", err)
	}

	return publicKey, nil
}
