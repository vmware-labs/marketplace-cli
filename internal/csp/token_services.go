// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package csp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/golang-jwt/jwt"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
)

type TokenServices struct {
	keyfunc jwt.Keyfunc
	keyPem  string
	CSPHost string
	Client  pkg.HTTPClient
}

type RedeemResponse struct {
	AccessToken string `json:"access_token"`
}

func (csp *TokenServices) Redeem(refreshToken string) (*Claims, error) {
	requestURL := pkg.MakeURL(csp.CSPHost, "/csp/gateway/am/api/auth/api-tokens/authorize", nil)
	formData := url.Values{
		"refresh_token": []string{refreshToken},
	}

	retried := false
	resp, err := csp.Client.PostForm(requestURL, formData)
	if err != nil {
		return nil, fmt.Errorf("failed to redeem token: %w", err)
	}

	if resp.StatusCode == http.StatusServiceUnavailable {
		retried = true
		resp, err = csp.Client.PostForm(requestURL, formData)
		if err != nil {
			return nil, fmt.Errorf("failed to redeem token on second attempt: %w", err)
		}
	}

	if resp.StatusCode != http.StatusOK {
		if !retried {
			return nil, fmt.Errorf("failed to exchange refresh token for access token: %s", resp.Status)
		}
		return nil, fmt.Errorf("failed twice to exchange refresh token for access token: %s", resp.Status)
	}

	var body RedeemResponse
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redeem response: %w", err)
	}

	claims := &Claims{}
	_, _ = jwt.ParseWithClaims(body.AccessToken, claims, func(t *jwt.Token) (interface{}, error) {
		// token was just retrieved, no need to validate
		return "not a valid key anyway", nil
	})
	// err != nil here are the token validation has failed

	claims.Token = body.AccessToken
	return claims, nil
}

func (csp *TokenServices) Validate(jwtAccessToken string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(jwtAccessToken, claims, csp.keyfunc)
	return claims, err
}

func (csp *TokenServices) VerificationKey() string {
	return csp.keyPem
}

func NewTokenServices(cspHost string, client pkg.HTTPClient) (*TokenServices, error) {
	keyData, err := fetchPublicKey(cspHost, client)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return nil, fmt.Errorf("failed to make public key structure: %w", err)
	}

	rsa := func(*jwt.Token) (interface{}, error) {
		return publicKey, nil
	}

	return &TokenServices{
		CSPHost: cspHost,
		Client:  client,
		keyfunc: rsa,
		keyPem:  string(keyData),
	}, nil
}

func fetchPublicKey(cspHost string, client pkg.HTTPClient) ([]byte, error) {
	resp, err := client.Get(pkg.MakeURL(cspHost, "/csp/gateway/am/api/auth/token-public-key", nil))
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

	return []byte(s), nil
}
