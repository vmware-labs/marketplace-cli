// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package internal

import (
	"fmt"
	"io"
	"time"

	"github.com/schollz/progressbar/v3"
)

//go:generate counterfeiter . ProgressBar
type ProgressBar interface {
	WrapWriter(source io.Writer) io.Writer
	WrapReader(source io.Reader) io.Reader
}

//go:generate counterfeiter . ProgressBarMaker
type ProgressBarMaker func(description string, length int64, output io.Writer) ProgressBar

var MakeProgressBar = NewProgressBar

func NewProgressBar(description string, length int64, output io.Writer) ProgressBar {
	return &ProgressBarImpl{
		description: description,
		length:      length,
		output:      output,
	}
}

type ProgressBarImpl struct {
	description string
	length      int64
	output      io.Writer
}

func (b *ProgressBarImpl) makeProgressBar() *progressbar.ProgressBar {
	bar := progressbar.NewOptions64(
		b.length,
		progressbar.OptionSetDescription(b.description),
		progressbar.OptionSetWriter(b.output),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			_, _ = fmt.Fprintln(b.output, "")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
	)
	_ = bar.RenderBlank()
	return bar
}

func (b *ProgressBarImpl) WrapWriter(source io.Writer) io.Writer {
	return io.MultiWriter(source, b.makeProgressBar())
}

func (b *ProgressBarImpl) WrapReader(source io.Reader) io.Reader {
	reader := progressbar.NewReader(source, b.makeProgressBar())
	return &reader
}
