// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package csp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	jwt "github.com/golang-jwt/jwt"
)

type TokenServices struct {
	keyfunc jwt.Keyfunc
	keyPem  string
	CSPHost string
}

func (csp *TokenServices) Redeem(refreshToken string) (*Claims, error) {
	formData := url.Values{
		"refresh_token": []string{refreshToken},
	}

	resp, err := http.PostForm(csp.CSPHost+"/csp/gateway/am/api/auth/api-tokens/authorize", formData)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("cannot exchange refresh token for access token: %d", resp.StatusCode)
	}

	m := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return nil, fmt.Errorf("bad response from CSP: %+v", err)
	}

	accessToken := m["access_token"].(string)
	if accessToken == "" {
		return nil, fmt.Errorf("bad response from server, access_token expected")
	}

	claims := &Claims{}
	_, _ = jwt.ParseWithClaims(accessToken, claims, func(t *jwt.Token) (interface{}, error) {
		// token was just retrieved, no need to validate
		return "not a valid key anyway", nil
	})

	// err != nil here are the token validation has failed

	claims.Token = accessToken
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

func NewTokenServices(cspHost string) (*TokenServices, error) {
	keyData, err := fetchPublicKey(cspHost)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return nil, err
	}

	rsa := func(*jwt.Token) (interface{}, error) {
		return publicKey, nil
	}

	return &TokenServices{
		CSPHost: cspHost,
		keyfunc: rsa,
		keyPem:  string(keyData),
	}, nil
}

func fetchPublicKey(cspLink string) ([]byte, error) {
	u := cspLink + "/csp/gateway/am/api/auth/token-public-key"
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	m := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return nil, err
	}

	pemData, ok := m["value"]
	if !ok {
		return nil, fmt.Errorf("cannot find validation key, value expected")
	}

	s, ok := pemData.(string)
	if !ok {
		return nil, fmt.Errorf("cannot find validation key, string expected for value")
	}

	return []byte(s), nil
}
