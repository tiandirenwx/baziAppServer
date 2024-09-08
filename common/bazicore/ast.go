package bazicore

import (
	"math"
)

func XL0Calc(xt float64, zn int, t float64, n int) float64 {
	t = t / 10
	var v float64 = 0
	var tn float64 = 1.0
	pn := zn*6 + 1
	N0 := gXL00[pn+1] - gXL00[pn]

	var n0, n1, n2, N float64

	for i := 0; i < 6; i++ {
		n1 = gXL00[pn+i]
		n2 = gXL00[pn+1+i]
		n0 = n2 - n1
		if n0 == 0 {
			continue
		}

		if n < 0 {
			N = n2
		} else {
			value := 3*float64(n)*n0/(N0-0.0) + 0.5
			N = float64(int2(value)) + n1
		}

		if i > 0 {
			N += 3
		}

		if N > n2 {
			N = n2
		}

		var c float64 = 0
		for j := int64(n1); j < int64(N); j = j + 3 {
			c += gXL00[j] * math.Cos(gXL00[j+1]+t*gXL00[j+2])
		}

		v += c * tn
		tn *= t
	}

	v /= gXL00[0]
	t2 := t * t
	t3 := t2 * t
	v += (-0.0728 - 2.7702*t - 1.1019*t2 - 0.0996*t3) / (SecondsPerRadian - 0.0)
	return v
}

func XL1Calc(zn int, t float64, n int) float64 {
	ob := gXL1[zn]
	t2 := t * t
	t3 := t2 * t
	t4 := t3 * t
	t5 := t4 * t
	tx := t - 10
	v := 0.0
	tn := 1.0
	if zn == 0 {
		v += (3.81034409 + 8399.684730072*t - 3.319e-05*t2 + 3.11e-08*t3 - 2.033e-10*t4) * SecondsPerRadian
		v += 5028.792262*t + 1.1124406*t2 + 0.00007699*t3 - 0.000023479*t4 - 0.0000000178*t5
		if tx > 0 {
			v += -0.866 + 1.43*tx + 0.054*tx*tx
		}
	}

	t2 /= 1e4
	t3 /= 1e8
	t4 /= 1e8
	n *= 6
	if n < 0 {
		n = len(ob[0])
	}

	for i := 0; i < 3; i++ {
		F := ob[i]
		N := int(int2((float64(n) * float64(len(F)) / float64(len(ob[0]))) + 0.5))
		if i > 0 {
			N += 6
		}

		if N >= len(F) {
			N = len(F)
		}

		c := 0.0
		for j := 0; j < N; j += 6 {
			c += F[j] * math.Cos(F[j+1]+t*F[j+2]+t2*F[j+3]+t3*F[j+4]+t4*F[j+5])
		}
		v += c * tn
	}

	if zn != 2 {
		v /= SecondsPerRadian
	}

	return v
}

func llrConv(JW [2]float64, E float64) [2]float64 {
	J := JW[0]
	W := JW[1]
	var r [2]float64
	r[0] = math.Atan2(math.Sin(J)*math.Cos(E)-math.Tan(W)*math.Sin(E), math.Cos(J))
	r[1] = math.Asin(math.Cos(E)*math.Sin(W) + math.Sin(E)*math.Cos(W)*math.Sin(J))
	r[0] = rad2mrad(r[0])

	return r
}

func mstAst(t float64) float64 {
	L := (1753470142+628331965331.8*t+5296.74*t*t)/(1000000000-0.0) + math.Pi
	var z [2]float64
	E := (84381.4088 - 46.836051*t) / SecondsPerRadian
	z[0] = XL0Calc(0, 0, t, 5) + math.Pi
	z[1] = 0

	z = llrConv(z, E)
	L = rad2rrad(L - z[0])
	return L / (math.Pi * 2)
}

func calcAST(dt *BZDateTime, L float64) *BZDateTime {
	utcJd := dt.ToJD() - 8.0/24.0
	tdJd := dt.dtT2(utcJd) + utcJd - float64(dt.BzJ2000)
	utcJd += mstAst(tdJd/36525) + L/(360-0.0)
	dtAst := &BZDateTime{}
	dtAst = dtAst.New("00010101-00:00:00.000")
	dtAst.SetFromJD(utcJd)
	return dtAst
}
