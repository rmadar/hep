// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"image/color"
	"math"

	"go-hep.org/x/hep/hbook"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
)

// BinnedErrBand implements the plot.Plotter interface,
// drawing a colored band for the error on any binned
// quantity.
type BinnedErrBand struct {

	// Data for every bins.
	Counts []hbook.Count

	// LineStyle is the style of the line
	// contouring the band.
	// Use zero width to disable.
	draw.LineStyle

	// FillColor is the color to fill the area
	// between the top and bottom data points.
	// Use nil to disable the filling.
	FillColor color.Color
}

// NewBinnedErrBand creates a binned error band
// from a slice of count.
func NewBinnedErrBand(cs []hbook.Count) *BinnedErrBand {
	return &BinnedErrBand{
		Counts: cs,
	}
}

// Plot implements the Plotter interface,
// drawing a colored box defined by width
// of bins (x-axis) and error (y-axis).
func (b *BinnedErrBand) Plot(c draw.Canvas, plt *plot.Plot) {

	for _, count := range b.Counts {

		// Don't do anything if both errors are zero
		if count.Err.Low == 0 && count.Err.High == 0 {
			continue
		}

		// Get four corner of the ith bin
		xmin, xmax := count.XRange.Min, count.XRange.Max
		y, ydo, yup := count.Val, count.Err.Low, count.Err.High
		xys := plotter.XYs{
			plotter.XY{X: xmin, Y: y - ydo},
			plotter.XY{X: xmin, Y: y + yup},
			plotter.XY{X: xmax, Y: y + yup},
			plotter.XY{X: xmax, Y: y - ydo},
		}

		// Polygon
		poly := plotter.Polygon{XYs: []plotter.XYs{xys}, Color: b.FillColor}
		poly.Plot(c, plt)

		// Bottom line
		xysBo := plotter.XYs{xys[0], xys[3]}
		lBo := plotter.Line{XYs: xysBo, LineStyle: b.LineStyle}
		lBo.Plot(c, plt)

		// Upper line
		xysUp := plotter.XYs{xys[1], xys[2]}
		lUp := plotter.Line{XYs: xysUp, LineStyle: b.LineStyle}
		lUp.Plot(c, plt)
	}
}

// DataRange returns the minimum and maximum x and
// y values, implementing the plot.DataRanger interface.
func (b *BinnedErrBand) DataRange() (xmin, xmax, ymin, ymax float64) {

	n := len(b.Counts) - 1
	xmin, xmax = b.Counts[0].XRange.Min, b.Counts[n].XRange.Max

	ymin, ymax = math.Inf(+1), math.Inf(-1)
	for _, c := range b.Counts {

		// Ignore bins for which there are no error
		if c.Err.Low == 0 && c.Err.High == 0 {
			continue
		}

		y, ydo, yup := c.Val, c.Err.Low, c.Err.High
		ymin = math.Min(ymin, y-ydo)
		ymax = math.Max(ymax, y+yup)
	}

	return xmin, xmax, ymin, ymax
}

var (
	_ plot.Plotter    = (*BinnedErrBand)(nil)
	_ plot.DataRanger = (*BinnedErrBand)(nil)
)
