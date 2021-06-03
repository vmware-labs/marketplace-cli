package lib_test

import (
	"io/ioutil"
	"net/url"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"gitlab.eng.vmware.com/marketplace-partner-eng/marketplace-cli/v2/lib"
)

var _ = Describe("Pagination", func() {
	It("returns a valid pagination URL value", func() {
		pagination := lib.Pagination(1, 25)
		Expect(pagination).To(HaveLen(1))
		Expect(pagination[0]).To(Equal(`{"page":1,"pagesize":25}`))
	})
})

var _ = Describe("MakeRequest", func() {
	BeforeEach(func() {
		viper.Set("csp.refresh-token", "secrets")
		viper.Set("marketplace.host", "marketplace.vmware.example")
	})

	It("Makes a valid request object", func() {
		content := strings.NewReader("everything totally passed")
		request, err := lib.MakeRequest(
			"POST",
			"/api/v1/unit-tests",
			url.Values{
				"color": []string{"blue", "green"},
			},
			map[string]string{
				"Content-Type": "text/plain",
			},
			content,
		)
		Expect(err).ToNot(HaveOccurred())

		Expect(request.Method).To(Equal("POST"))

		By("building the right url", func() {
			Expect(request.URL.Scheme).To(Equal("https"))
			Expect(request.URL.Host).To(Equal("marketplace.vmware.example"))
			Expect(request.URL.Path).To(Equal("/api/v1/unit-tests"))
			Expect(request.URL.Query().Encode()).To(Equal("color=blue&color=green"))
		})

		By("setting the right headers", func() {
			Expect(request.Header.Get("Accept")).To(Equal("application/json"))
			Expect(request.Header.Get("csp-auth-token")).To(Equal("secrets"))
			Expect(request.Header.Get("Content-Type")).To(Equal("text/plain"))
		})

		By("including the right content", func() {
			Expect(ioutil.ReadAll(request.Body)).To(Equal([]byte("everything totally passed")))
		})
	})
})

var _ = Describe("MakeGetRequest", func() {
	BeforeEach(func() {
		viper.Set("csp.refresh-token", "secrets")
		viper.Set("marketplace.host", "marketplace.vmware.example")
	})

	It("Makes a valid request object", func() {
		request, err := lib.MakeGetRequest(
			"/api/v1/unit-tests",
			url.Values{
				"color": []string{"blue", "green"},
			},
		)
		Expect(err).ToNot(HaveOccurred())

		Expect(request.Method).To(Equal("GET"))

		By("building the right url", func() {
			Expect(request.URL.Scheme).To(Equal("https"))
			Expect(request.URL.Host).To(Equal("marketplace.vmware.example"))
			Expect(request.URL.Path).To(Equal("/api/v1/unit-tests"))
			Expect(request.URL.Query().Encode()).To(Equal("color=blue&color=green"))
		})

		By("setting the right headers", func() {
			Expect(request.Header.Get("Accept")).To(Equal("application/json"))
			Expect(request.Header.Get("csp-auth-token")).To(Equal("secrets"))
		})
	})
})
