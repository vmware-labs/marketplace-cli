package main

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"github.com/vmware-labs/marketplace-cli/v2/pkg"
)

var (
	Client      pkg.HTTPClient
	Marketplace pkg.MarketplaceInterface
)

func main() {
	logFile, err := os.OpenFile("agentlogfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	defer logFile.Close()

	// Set the log output to the log file
	log.SetOutput(logFile)

	initMarketplace()
	StartMonitoring()
}

func initMarketplace() {
	Client = pkg.NewClient(
		os.Stderr,
		viper.GetBool("debugging.enabled"),
		viper.GetBool("debugging.print-request-payloads"),
		viper.GetBool("debugging.print-response-payloads"),
	)

	Marketplace = &pkg.Marketplace{
		Host:          viper.GetString("marketplace.host"),
		APIHost:       viper.GetString("marketplace.api-host"),
		UIHost:        viper.GetString("marketplace.ui-host"),
		StorageBucket: viper.GetString("marketplace.storage.bucket"),
		StorageRegion: viper.GetString("marketplace.storage.region"),
		Client:        Client,
		Output:        os.Stderr,
	}
}

func StartMonitoring() {
	log.Println("Agent Monitoring process has started....")
	placeHolderApi()
	log.Println("Agent process has completed")
}

func placeHolderApi() {
	// TODO: Place holder api
	ProductSlug := "testlisting-solution-2"

	for {
		log.Println("Agent  process has executing. Fetching the product details")
		product, err := Marketplace.GetProduct(ProductSlug)
		if err != nil {
			log.Println("failed to fetch product from Marketplace...")
		}

		log.Println("Fetched product details :", product.DisplayName)

		time.Sleep(20 * time.Second)
	}
}
