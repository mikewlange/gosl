// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/utl"
)

// Metrics holds data related to a position in a space represented by curvilinear coordinates
type Metrics struct {
	U          la.Vector     // reference coordinates {r,s,t}
	X          la.Vector     // physical coordinates {x,y,z}
	CovG0      la.Vector     // covariant basis g_0 = d{x}/dr
	CovG1      la.Vector     // covariant basis g_1 = d{x}/ds
	CovG2      la.Vector     // covariant basis g_2 = d{x}/dt
	CovGmat    *la.Matrix    // covariant metrics g_ij = g_i ⋅ g_j
	CntGmat    *la.Matrix    // contravariant metrics g^ij = g^i ⋅ g^j
	DetCovGmat float64       // determinant of covariant g matrix = det(CovGmat)
	GammaS     [][][]float64 // [k][i][j] Christoffel coefficients of second kind
	L          []float64     // [3] L-coefficients = sum(Γ_ij^k ⋅ g^ij)
}

// NewMetrics2d allocate new 2D metrics structure
func NewMetrics2d(u, x, dxdr, dxds, ddxdrr, ddxdss, ddxdrs la.Vector) (o *Metrics) {

	// input
	o = new(Metrics)
	o.U = u.GetCopy()
	o.X = x.GetCopy()
	o.CovG0 = dxdr.GetCopy()
	o.CovG1 = dxds.GetCopy()

	// covariant metrics
	o.CovGmat = la.NewMatrix(2, 2)
	o.CovGmat.Set(0, 0, la.VecDot(o.CovG0, o.CovG0))
	o.CovGmat.Set(1, 1, la.VecDot(o.CovG1, o.CovG1))
	o.CovGmat.Set(0, 1, la.VecDot(o.CovG0, o.CovG1))
	o.CovGmat.Set(1, 0, o.CovGmat.Get(0, 1))

	// contravariant metrics
	o.CntGmat = la.NewMatrix(2, 2)
	o.DetCovGmat = la.MatInvSmall(o.CntGmat, o.CovGmat, 1e-13)

	// contravariant basis vectors
	cntG0, cntG1 := la.NewVector(2), la.NewVector(2)
	for i := 0; i < 2; i++ {
		cntG0[i] += o.CntGmat.Get(0, 0)*o.CovG0[i] + o.CntGmat.Get(0, 1)*o.CovG1[i]
		cntG1[i] += o.CntGmat.Get(1, 0)*o.CovG0[i] + o.CntGmat.Get(1, 1)*o.CovG1[i]
	}

	// Christoffel vectors
	Γ00, Γ11, Γ01 := ddxdrr, ddxdss, ddxdrs

	// Christoffel symbols of second kind
	o.GammaS = utl.Deep3alloc(2, 2, 2)
	o.GammaS[0][0][0] = la.VecDot(Γ00, cntG0)
	o.GammaS[0][1][1] = la.VecDot(Γ11, cntG0)
	o.GammaS[0][0][1] = la.VecDot(Γ01, cntG0)
	o.GammaS[0][1][0] = o.GammaS[0][0][1]
	o.GammaS[1][0][0] = la.VecDot(Γ00, cntG1)
	o.GammaS[1][1][1] = la.VecDot(Γ11, cntG1)
	o.GammaS[1][0][1] = la.VecDot(Γ01, cntG1)
	o.GammaS[1][1][0] = o.GammaS[1][0][1]

	// L-coefficients
	o.L = make([]float64, 2)
	o.L[0] = o.GammaS[0][0][0]*o.CntGmat.Get(0, 0) + o.GammaS[0][1][1]*o.CntGmat.Get(1, 1) + 2.0*o.GammaS[0][0][1]*o.CntGmat.Get(0, 1)
	o.L[1] = o.GammaS[1][0][0]*o.CntGmat.Get(0, 0) + o.GammaS[1][1][1]*o.CntGmat.Get(1, 1) + 2.0*o.GammaS[1][0][1]*o.CntGmat.Get(0, 1)
	return
}
