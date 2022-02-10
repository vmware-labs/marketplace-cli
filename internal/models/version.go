// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package models

import (
	"sort"
	"strings"

	"github.com/coreos/go-semver/semver"
)

type Version struct {
	Number           string `json:"versionnumber"`
	Details          string `json:"versiondetails"`
	Status           string `json:"status,omitempty"`
	Instructions     string `json:"versioninstruction"`
	CreatedOn        int32  `json:"createdon,omitempty"`
	HasLimitedAccess bool   `json:"haslimitedaccess,omitempty"`
	Tag              string `json:"tag,omitempty"`
}

func (product *Product) GetVersion(version string) *Version {
	if version == "" {
		return product.GetLatestVersion()
	}

	for _, v := range product.AllVersions {
		if v.Number == version {
			return v
		}
	}
	return nil
}

func (product *Product) GetLatestVersion() *Version {
	if len(product.AllVersions) == 0 {
		return nil
	}

	version, err := product.getLatestVersionSemver()
	if err != nil {
		version = product.getLatestVersionAlphanumeric()
	}

	return version
}

func (product *Product) getLatestVersionSemver() (*Version, error) {
	latestVersion := product.AllVersions[0]
	version, err := semver.NewVersion(latestVersion.Number)
	if err != nil {
		return nil, err
	}
	for _, v := range product.AllVersions {
		otherVersion, err := semver.NewVersion(v.Number)
		if err != nil {
			return nil, err
		}
		if version.LessThan(*otherVersion) {
			latestVersion = v
			version = otherVersion
		}
	}

	return latestVersion, nil
}

func (product *Product) getLatestVersionAlphanumeric() *Version {
	latestVersion := product.AllVersions[0]
	for _, v := range product.AllVersions {
		if strings.Compare(latestVersion.Number, v.Number) < 0 {
			latestVersion = v
		}
	}
	return latestVersion
}

type Versions []*Version

func (v Versions) Len() int {
	return len(v)
}

func (v Versions) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v Versions) Less(i, j int) bool {
	return v[i].LessThan(*v[j])
}

func (a Version) LessThan(b Version) bool {
	semverA, errA := semver.NewVersion(a.Number)
	semverB, errB := semver.NewVersion(b.Number)

	if errA != nil || errB != nil {
		return strings.Compare(a.Number, b.Number) < 0
	}

	return semverA.LessThan(*semverB)
}

func Sort(versions []*Version) {
	sort.Sort(sort.Reverse(Versions(versions)))
}

func (product *Product) HasVersion(version string) bool {
	return product.GetVersion(version) != nil
}
