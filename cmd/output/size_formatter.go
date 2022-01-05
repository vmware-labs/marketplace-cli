// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package output

import "fmt"

var units = []string{"B", "KB", "MB", "GB", "TB", "PB"}

func FormatSize(size int64) string {
	var unit string
	newSize := float64(size)
	for _, unit = range units {
		if newSize < 1000 {
			break
		}
		newSize /= 1000
	}
	return fmt.Sprintf("%.3g %s", newSize, unit)
}
