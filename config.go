package main

//
//import (
//	"fmt"
//
//	csptoken "gitlab.eng.vmware.com/vdp/common-go/csp/token"
//)
////
////type Config struct {
////	CloudAPIToken   string `long:"api-token" env:"CLOUD_API_TOKEN" description:"Cloud Services API Token"`
////	CSPHost         string `long:"csp-hostname" env:"CSP_HOSTNAME" default:"console.cloud.vmware.com"`
////	MarketplaceHost string `long:"marketplace-hostname" env:"MARKETPLACE_HOSTNAME" default:"gtwstg.market.csp.vmware.com"`
////
////	CSPRefreshToken *csptoken.CspClaims
////}
//
//func (c *Config) GetRefreshToken() error {
//	fmt.Printf("Using CSP: %s\n", c.CSPHost)
//
//	fmt.Println("Initializing token services")
//	tokenServices, err := csptoken.NewCspTokenServices(fmt.Sprintf("https://%s/", c.CSPHost))
//	if err != nil {
//		return err
//	}
//
//	fmt.Println("Token Exchange")
//	claims, err := tokenServices.Redeem(c.CloudAPIToken)
//	if err != nil {
//		return err
//	}
//
//	fmt.Printf("claims: %+v\n", claims)
//	c.CSPRefreshToken = claims
//	return nil
//}
