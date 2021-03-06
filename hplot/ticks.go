// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"fmt"
	"math"
	"strconv"

	"go-hep.org/x/hep/hplot/internal/talbot"
	"gonum.org/v1/gonum/floats/scalar"
	"gonum.org/v1/plot"
)

const (
	// displayPrecision is a sane level of float precision for a plot.
	displayPrecision = 4
)

// FreqTicks implements a simple plot.Ticker scheme.
// FreqTicks will generate N ticks where 1 every Freq tick will be labeled.
type FreqTicks struct {
	N    int // number of ticks
	Freq int // frequency of labeled ticks
}

// Ticks returns Ticks in a specified range
func (ft FreqTicks) Ticks(min, max float64) []plot.Tick {
	prec := maxInt(precisionOf(min), precisionOf(max))
	ticks := make([]plot.Tick, ft.N)
	for i := range ticks {
		v := min + float64(i)*(max-min)/float64(len(ticks)-1)
		label := ""
		if i%ft.Freq == 0 {
			label = formatFloatTick(v, prec)
		}
		ticks[i] = plot.Tick{Value: v, Label: label}
	}
	return ticks
}

// formatFloatTick returns a g-formated string representation of v
// to the specified precision.
func formatFloatTick(v float64, prec int) string {
	return strconv.FormatFloat(scalar.Round(v, prec), 'g', displayPrecision, 64)
}

// precisionOf returns the precision needed to display x without e notation.
func precisionOf(x float64) int {
	return int(math.Max(math.Ceil(-math.Log10(math.Abs(x))), displayPrecision))
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// NoTicks implements plot.Ticker but does not display any tick.
type NoTicks struct{}

// Ticks returns Ticks in a specified range
func (NoTicks) Ticks(min, max float64) []plot.Tick {
	return nil
}

// Ticks implements plot.Ticker.
// Ticks allows to specify the maximum number of major ticks to display.
// The zero value of Ticks display a maximum number of 3 major ticks.
type Ticks struct {
	N int // N is the suggested number of major ticks to display.

	// Format is an optional major-tick formatter.
	// If empty, a format will be automatically chosen.
	Format string
}

func (tck Ticks) Ticks(min, max float64) []plot.Tick {
	if tck.N == 0 {
		tck.N = 3 // same default than plot.DefaultTicks
	}

	ticks := talbot.Ticks(min, max, tck.N)
	if xfmt := tck.Format; xfmt != "" {
		for i, tck := range ticks {
			if tck.IsMinor() {
				continue
			}
			ticks[i].Label = fmt.Sprintf(xfmt, tck.Value)
		}
	}
	return ticks
}
