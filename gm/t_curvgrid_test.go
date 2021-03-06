// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func TestCurvGrid01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("CurvGrid01. 2d ring")

	// coordinates
	R := utl.LinSpace(-1, 1, 5)
	S := utl.LinSpace(-1, 1, 5)

	// curvgrid
	a, b := 1.0, 2.0
	trf := FactoryTfinite.Surf2dQuarterRing(a, b)
	cg := new(CurvGrid)
	cg.SetTransfinite2d(trf, R, S)

	// check
	π := math.Pi
	A := (b - a) / 2.0 // dρ/dr
	B := π / 4.0       // dα/ds
	for n := 0; n < cg.N1; n++ {
		for m := 0; m < cg.N0; m++ {
			mtr := cg.M2d[n][m]
			ρ := a + (1.0+mtr.U[0])*A // cylindrical coordinates
			α := (1.0 + mtr.U[1]) * B // cylindrical coordinates
			cα, sα := math.Cos(α), math.Sin(α)
			chk.Array(tst, "x      ", 1e-15, mtr.X, []float64{ρ * cα, ρ * sα})
			chk.Array(tst, "CovG0  ", 1e-15, mtr.CovG0, []float64{cα * A, sα * A})
			chk.Array(tst, "CovG1  ", 1e-15, mtr.CovG1, []float64{-ρ * sα * B, ρ * cα * B})
			chk.Deep2(tst, "CovGmat", 1e-15, mtr.CovGmat.GetDeep2(), [][]float64{
				{A * A, 0.0},
				{0.0, ρ * ρ * B * B},
			})
			chk.Deep2(tst, "CntGmat", 1e-14, mtr.CntGmat.GetDeep2(), [][]float64{
				{1.0 / (A * A), 0.0},
				{0.0, 1.0 / (ρ * ρ * B * B)},
			})
			chk.Deep3(tst, "GammaS", 1e-15, mtr.GammaS, [][][]float64{
				{
					{0, 0},
					{0, -ρ * B * B / A},
				},
				{
					{0, A / ρ},
					{A / ρ, 0},
				},
			})
			chk.Array(tst, "L", 1e-14, mtr.L, []float64{-1.0 / (ρ * A), 0})
		}
	}

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		trf.Draw([]int{11, 21}, false, &plt.A{C: plt.C(2, 9)}, &plt.A{C: plt.C(3, 9), Lw: 2})
		cg.DrawBases(0.15, nil, nil, nil)
		plt.HideAllBorders()
		plt.Equal()
		plt.Save("/tmp/gosl/gm", "curvgrid01")
	}
}
