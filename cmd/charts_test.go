// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package cmd_test

import (
	"fmt"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/vmware-labs/marketplace-cli/v2/cmd"
	"github.com/vmware-labs/marketplace-cli/v2/cmd/output/outputfakes"
	"github.com/vmware-labs/marketplace-cli/v2/internal/models"
	"github.com/vmware-labs/marketplace-cli/v2/pkg/pkgfakes"
	"github.com/vmware-labs/marketplace-cli/v2/test"
)

var _ = Describe("Charts", func() {
	var (
		marketplace *pkgfakes.FakeMarketplaceInterface
		output      *outputfakes.FakeFormat
		product     *models.Product
	)

	BeforeEach(func() {
		marketplace = &pkgfakes.FakeMarketplaceInterface{}
		cmd.Marketplace = marketplace

		marketplace.GetProductWithVersionStub = func(slug string, version string) (*models.Product, *models.Version, error) {
			Expect(slug).To(Equal("my-super-product"))
			Expect(version).To(Equal("1.2.3"))
			return product, &models.Version{Number: "1.2.3"}, nil
		}

		output = &outputfakes.FakeFormat{}
		cmd.Output = output
	})

	Describe("ListChartsCmd", func() {
		BeforeEach(func() {
			product = test.CreateFakeProduct("", "My Super Product", "my-super-product", "PENDING")
			test.AddVerions(product, "1.2.3")
			product.ChartVersions = []*models.ChartVersion{
				{
					Id:         uuid.New().String(),
					Version:    "5.0.0",
					AppVersion: "1.2.3",
					Repo: &models.Repo{
						Name: "Bitnami charts repo @ Github",
						Url:  "https://github.com/bitnami/charts/tree/master/bitnami/kube-prometheus",
					},
					HelmTarUrl: "https://charts.bitnami.com/bitnami/kube-prometheus-5.0.0.tgz",
					TarUrl:     "https://charts.bitnami.com/bitnami/kube-prometheus-5.0.0.tgz",
				},
			}
			cmd.ListChartsCmd.SetErr(gbytes.NewBuffer())
		})

		It("outputs the charts", func() {
			cmd.ChartProductSlug = "my-super-product"
			cmd.ChartProductVersion = "1.2.3"
			err := cmd.ListChartsCmd.RunE(cmd.ListChartsCmd, []string{""})
			Expect(err).ToNot(HaveOccurred())

			By("getting the product details", func() {
				Expect(marketplace.GetProductWithVersionCallCount()).To(Equal(1))
			})

			By("outputting the response", func() {
				Expect(output.RenderChartsCallCount()).To(Equal(1))
				charts := output.RenderChartsArgsForCall(0)
				Expect(charts).To(HaveLen(1))
				Expect(charts[0].Version).To(Equal("5.0.0"))
				Expect(charts[0].AppVersion).To(Equal("1.2.3"))
				Expect(charts[0].Repo.Name).To(Equal("Bitnami charts repo @ Github"))
				Expect(charts[0].Repo.Url).To(Equal("https://github.com/bitnami/charts/tree/master/bitnami/kube-prometheus"))
				Expect(charts[0].HelmTarUrl).To(Equal("https://charts.bitnami.com/bitnami/kube-prometheus-5.0.0.tgz"))
				Expect(charts[0].TarUrl).To(Equal("https://charts.bitnami.com/bitnami/kube-prometheus-5.0.0.tgz"))
			})
		})

		Context("Error fetching product", func() {
			BeforeEach(func() {
				marketplace.GetProductWithVersionReturns(nil, nil, fmt.Errorf("get product failed"))
			})

			It("prints the error", func() {
				cmd.ChartProductSlug = "my-super-product"
				cmd.ChartProductVersion = "1.2.3"
				err := cmd.ListChartsCmd.RunE(cmd.ListChartsCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("get product failed"))
			})
		})
	})

	Describe("GetChartCmd", func() {
		var chartId string
		BeforeEach(func() {
			product = test.CreateFakeProduct("", "My Super Product", "my-super-product", "PENDING")
			test.AddVerions(product, "1.2.3")
			chartId = uuid.New().String()
			product.ChartVersions = []*models.ChartVersion{
				{
					Id:         chartId,
					Version:    "5.0.0",
					AppVersion: "1.2.3",
					Repo: &models.Repo{
						Name: "Bitnami charts repo @ Github",
						Url:  "https://github.com/bitnami/charts/tree/master/bitnami/kube-prometheus",
					},
					HelmTarUrl: "https://charts.bitnami.com/bitnami/kube-prometheus-5.0.0.tgz",
					TarUrl:     "https://charts.bitnami.com/bitnami/kube-prometheus-5.0.0.tgz",
				},
			}
		})

		It("outputs the chart", func() {
			cmd.ChartProductSlug = "my-super-product"
			cmd.ChartProductVersion = "1.2.3"
			cmd.ChartID = chartId
			err := cmd.GetChartCmd.RunE(cmd.GetChartCmd, []string{""})
			Expect(err).ToNot(HaveOccurred())

			By("getting the product details", func() {
				Expect(marketplace.GetProductWithVersionCallCount()).To(Equal(1))
			})

			By("outputting the response", func() {
				Expect(output.RenderChartCallCount()).To(Equal(1))
				chart := output.RenderChartArgsForCall(0)
				Expect(chart.Id).To(Equal(chartId))
				Expect(chart.AppVersion).To(Equal("1.2.3"))
				Expect(chart.Version).To(Equal("5.0.0"))
				Expect(chart.Repo.Name).To(Equal("Bitnami charts repo @ Github"))
				Expect(chart.Repo.Url).To(Equal("https://github.com/bitnami/charts/tree/master/bitnami/kube-prometheus"))
				Expect(chart.HelmTarUrl).To(Equal("https://charts.bitnami.com/bitnami/kube-prometheus-5.0.0.tgz"))
				Expect(chart.TarUrl).To(Equal("https://charts.bitnami.com/bitnami/kube-prometheus-5.0.0.tgz"))
			})
		})

		Context("Error fetching product", func() {
			BeforeEach(func() {
				marketplace.GetProductWithVersionReturns(nil, nil, fmt.Errorf("get product failed"))
			})

			It("prints the error", func() {
				cmd.ChartProductSlug = "my-super-product"
				cmd.ChartProductVersion = "1.2.3"
				cmd.ChartID = chartId
				err := cmd.GetChartCmd.RunE(cmd.GetChartCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("get product failed"))
			})
		})

		Context("No matching chart found", func() {
			It("says that the version does not exist", func() {
				cmd.ChartProductSlug = "my-super-product"
				cmd.ChartProductVersion = "1.2.3"
				cmd.ChartID = "this-chart-id-does-not-exist"
				err := cmd.GetChartCmd.RunE(cmd.GetChartCmd, []string{""})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("my-super-product 1.2.3 does not have the chart \"this-chart-id-does-not-exist\""))
			})
		})
	})
})
