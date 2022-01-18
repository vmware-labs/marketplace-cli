// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package models

type SKUPublisherInfo struct {
	SKUNumber            string   `json:"skunumber"`
	Price                string   `json:"price"`
	Currency             string   `json:"currency"`
	BillFrequency        string   `json:"billfreq"`
	TermLength           int32    `json:"termlength"`
	UnitOfMeasurement    string   `json:"unitofmeasurement"`
	PriceMeasurementUnit []string `json:"pricemeasurementunitlist"`
	IsMonthlySKUEnabled  bool     `json:"ismonthlyskuenabled"`
	MonthlySKUNumber     string   `json:"monthlyskunumber"`
}

type SKUPublisherView struct {
	SKUPublisherInfo *SKUPublisherInfo `json:"skupublisherinfo"`
	SKUID            string            `json:"skuid"`
	Description      string            `json:"description"`
	Status           string            `json:"status"`
}
