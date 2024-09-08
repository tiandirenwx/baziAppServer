package bazicore

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

func jieYa(s0 string) string {
	o := "0000000000"
	o2 := o + o
	s := s0
	s = strings.Replace(s, "J", "00", -1)
	s = strings.Replace(s, "I", "000", -1)
	s = strings.Replace(s, "H", "0000", -1)
	s = strings.Replace(s, "G", "00000", -1)
	s = strings.Replace(s, "t", "02", -1)
	s = strings.Replace(s, "s", "002", -1)
	s = strings.Replace(s, "r", "0002", -1)
	s = strings.Replace(s, "q", "00002", -1)
	s = strings.Replace(s, "p", "000002", -1)
	s = strings.Replace(s, "o", "0000002", -1)
	s = strings.Replace(s, "n", "00000002", -1)
	s = strings.Replace(s, "m", "000000002", -1)
	s = strings.Replace(s, "l", "0000000002", -1)
	s = strings.Replace(s, "k", "01", -1)
	s = strings.Replace(s, "j", "0101", -1)
	s = strings.Replace(s, "i", "001", -1)
	s = strings.Replace(s, "h", "001001", -1)
	s = strings.Replace(s, "g", "0001", -1)
	s = strings.Replace(s, "f", "00001", -1)
	s = strings.Replace(s, "e", "000001", -1)
	s = strings.Replace(s, "d", "0000001", -1)
	s = strings.Replace(s, "c", "00000001", -1)
	s = strings.Replace(s, "b", "000000001", -1)
	s = strings.Replace(s, "a", "0000000001", -1)
	s = strings.Replace(s, "A", o2+o2+o2, -1)
	s = strings.Replace(s, "B", o2+o2+o, -1)
	s = strings.Replace(s, "C", o2+o2, -1)
	s = strings.Replace(s, "D", o2+o, -1)
	s = strings.Replace(s, "E", o2, -1)
	s = strings.Replace(s, "F", o, -1)
	return s
}

// # 返回黄赤交角(常规精度),短期精度很高
func hcjj1(t float64) float64 {
	t1 := t / (36525 - 0.0)
	t2 := t1 * t1
	t3 := t2 * t1
	return (gHcjjB[0] + gHcjjB[1]*t1 + gHcjjB[2]*t2 + gHcjjB[3]*t3) / SecondsPerRadian
}

// # 黄赤转换(黄赤坐标旋转)
func HCconv(JW []float64, E float64) {
	HJ := rad2mrad(JW[0])
	HW := JW[1]
	sinE := math.Sin(E)
	cosE := math.Cos(E)
	sinW := cosE*math.Sin(HW) + sinE*math.Cos(HW)*math.Sin(HJ)
	J := math.Atan2(math.Sin(HJ)*cosE-math.Tan(HW)*sinE, math.Cos(HJ))

	JW[0] = rad2mrad(J)
	JW[1] = math.Asin(sinW)
}

// # 补岁差
func addPrece(jd float64, zb []float64) {
	t := 1.0
	v := 0.0
	t1 := jd / (365250 - 0.0)
	for i := 1; i < 8; i++ {
		t *= t1
		v += gPreceB[i] * t
	}

	zb[0] = rad2mrad(zb[0] + (v+2.9965*t1)/SecondsPerRadian)
}

// # 恒星周年光行差计算(黄道坐标中)
func addGxc(t float64, zb []float64) {
	t1 := t / (36525 - 0.0)
	t2 := t1 * t1
	t3 := t2 * t1
	t4 := t3 * t1
	L := ggxcL[0] + ggxcL[1]*t1 + ggxcL[2]*t2 + ggxcL[3]*t3 + ggxcL[4]*t4
	p := ggxcP[0] + ggxcP[1]*t1 + ggxcP[2]*t2
	e := ggxcE[0] + ggxcE[1]*t1 + ggxcE[2]*t2
	dL := L - zb[0]
	dP := p - zb[0]
	zb[0] -= ggxcK * (math.Cos(dL) - e*math.Cos(dP)) / math.Cos(zb[1])
	zb[1] -= ggxcK * math.Sin(zb[1]) * (math.Sin(dL) - e*math.Sin(dP))
	zb[0] = rad2mrad(zb[0])
}

type SLonObl struct {
	Lon float64
	Obl float64
}

// # 计算黄经章动及交角章动
func nutation(t float64, d *SLonObl) {
	d = &SLonObl{
		Lon: 0.0,
		Obl: 0.0,
	}
	t /= (36525 - 0.0)
	t1 := t
	t2 := t1 * t1
	t3 := t2 * t1
	t4 := t3 * t1

	//#     t5 = t4 * t1
	for i := 0; i < len(gNutB); i++ {
		c := gNutB[i] + gNutB[i+1]*t1 + gNutB[i+2]*t2 + gNutB[i+3]*t3 + gNutB[i+4]*t4
		d.Lon += (gNutB[i+5] + gNutB[i+6]*t/10) * math.Sin(c) //# 黄经章动
		d.Obl += (gNutB[i+7] + gNutB[i+8]*t/10) * math.Cos(c) //# 交角章动
	}
	d.Lon /= SecondsPerRadian * 10000
	d.Obl /= SecondsPerRadian * 10000
}

// # 本函数计算赤经章动及赤纬章动
func nutationRaDec(t float64, zb []float64) {
	Ra := zb[0]
	Dec := zb[1]
	E := hcjj1(t)
	sinE := math.Sin(E)
	cosE := math.Cos(E)
	d := &SLonObl{}
	nutation(t, d) //# 计算黄经章动及交角章动
	cosRa := math.Cos(Ra)
	sinRa := math.Sin(Ra)
	tanDec := math.Tan(Dec)
	zb[0] += (cosE+sinE*sinRa*tanDec)*d.Lon - cosRa*tanDec*d.Obl //# 赤经章动
	zb[1] += sinE*cosRa*d.Lon + sinRa*d.Obl                      //# 赤纬章动
	zb[0] = rad2mrad(zb[0])
}

// EnnT = 0  # 调用Enn前先设置EnnT时间变量
// # 计算E10,E11,E20等,即:某一组周期项或泊松项算出,计算前先设置EnnT时间

func Enn(F []float64) float64 {
	v := 0.0
	for i := 0; i < len(F); i += 3 {
		v += F[i] * math.Cos(F[i+1]+gEnnT*F[i+2])
	}

	return v
}

// # 返回地球位置,日心Date黄道分点坐标
func earCal(jd float64, llr []float64) {
	gEnnT = jd / (365250 - 0.0)

	t1 := gEnnT
	t2 := t1 * t1
	t3 := t2 * t1
	t4 := t3 * t1
	t5 := t4 * t1
	//# 测试
	//# print Enn(E10)  #1.75444665847之后的数据不一致
	llr[0] = (Enn(gE10) + Enn(gE11)*t1 + Enn(gE12)*t2 + Enn(gE13)*t3 + Enn(gE14)*t4 + Enn(gE15)*t5)
	llr[1] = (Enn(gE20) + Enn(gE21)*t1)
	llr[2] = (Enn(gE30) + Enn(gE31)*t1 + Enn(gE32)*t2 + Enn(gE33)*t3)

	llr[0] = rad2mrad(llr[0])
}

// //# 返回太阳视位置
func sunCal(jd float64, sun []float64) {
	earCal(jd, sun)
	sun[0] += math.Pi
	sun[1] = -sun[1] //# 计算太阳真位置
	addGxc(jd, sun)  //# 补周年黄经光行差
}

// # 返回太阳视位置
func sunCal2(jd float64, sun []float64) {
	earCal(jd, sun)
	sun[0] += math.Pi
	sun[1] = -sun[1] //# 计算太阳真位置
	d := &SLonObl{}
	nutation(jd, d)
	sun[0] = rad2mrad(sun[0] + d.Lon) //# 补章动
	addGxc(jd, sun)                   //# 补周年黄经光行差
}

// # 计算M10,M11,M20等,计算前先设置MnnT时间
func Mnn(F []float64) float64 {
	v := 0.0
	t1 := gMnnT
	t2 := t1 * t1
	t3 := t2 * t1
	t4 := t3 * t1
	for i := 0; i < len(F); i += 6 {
		v += F[i] * math.Sin(F[i+1]+t1*F[i+2]+t2*F[i+3]+t3*F[i+4]+t4*F[i+5])
	}

	return v
}

// # 返回月球位置,返回地心Date黄道坐标
func moonCal(jd float64, llr []float64) {
	gMnnT = jd / (36525 - 0.0)
	t1 := gMnnT
	t2 := t1 * t1
	t3 := t2 * t1
	t4 := t3 * t1

	llr[0] = ((Mnn(gM10) + Mnn(gM11)*t1 + Mnn(gM12)*t2) / SecondsPerRadian)
	llr[1] = ((Mnn(gM20) + Mnn(gM21)*t1) / SecondsPerRadian)
	llr[2] = ((Mnn(gM30) + Mnn(gM31)*t1) * 0.999999949827)
	llr[0] = llr[0] + gM1n[0] + gM1n[1]*t1 + gM1n[2]*t2 + gM1n[3]*t3 + gM1n[4]*t4
	llr[0] = rad2mrad(llr[0]) //# 地心Date黄道原点坐标(不含岁差)
	addPrece(jd, llr)         //# 补岁差
}

// 传回月球的地心视黄经及视黄纬
func moonCal2(jd float64, moon []float64) {
	moonCal(jd, moon)
	d := &SLonObl{}
	nutation(jd, d)
	moon[0] = rad2mrad(moon[0] + d.Lon) //# 补章动
}

// # 传回月球的地心视赤经及视赤纬
func moonCal3(jd float64, moon []float64) {
	moonCal(jd, moon)
	HCconv(moon, hcjj1(jd))
	nutationRaDec(jd, moon) //# 补赤经及赤纬章动
	//# 如果黄赤转换前补了黄经章动及交章动,就不能再补赤经赤纬章动
}

// #==================地心坐标中的日月位置计算===================
func jiaoCai(lx, t, jiao float64) float64 {
	//# lx=1时计算t时刻日月角距与jiao的差, lx=0计算t时刻太阳黄经与jiao的差
	var sun []float64
	earCal(t, sun) //# 计算太阳真位置(先算出日心坐标中地球的位置)
	sun[0] += math.Pi
	sun[1] = -sun[1] //# 转为地心坐标
	addGxc(t, sun)   //# 补周年光行差
	if lx == 0 {
		d := &SLonObl{}
		nutation(t, d)
		sun[0] += d.Lon //# 补黄经章动
		return rad2mrad(jiao - sun[0])
	}
	var moon []float64
	moonCal(t, moon) //# 日月角差与章动无关
	return rad2mrad(jiao - (moon[0] - sun[0]))
}

// #==================已知位置反求时间===================
func jiaoCal(t1, jiao, lx float64) float64 {
	/*
	   #     t1是J2000起算儒略日数
	   #     已知角度(jiao)求时间(t)
	   #     lx=0是太阳黄经达某角度的时刻计算(用于节气计算)
	   #     lx=1是日月角距达某角度的时刻计算(用于定朔望等)
	   #     传入的t1是指定角度对应真时刻t的前一些天
	   #     对于节气计算,应满足t在t1到t1+360天之间,对于Y年第n个节气(n=0是春分),t1可取值Y*365.2422+n*15.2
	   #     对于朔望计算,应满足t在t1到t1+25天之间,在此范围之外,求右边的根
	   #     print (jiao,t1,lx)  #当jiao为0的时候t1出现
	*/
	t2 := t1
	t := 0.0

	if lx == 0 {
		t2 += 360 //# 在t1到t2范围内求解(范气360天范围),结果置于t
	} else {
		t2 += 25
	}

	jiao *= math.Pi / (180 - 0.0) //# 待搜索目标角

	//# 利用截弦法计算
	v1 := jiaoCai(lx, t1, jiao) //# v1,v2为t1,t2时对应的黄经
	v2 := jiaoCai(lx, t2, jiao)

	if v1 < v2 {
		v2 -= 2 * math.Pi //# 减2pi作用是将周期性角度转为连续角度
	}

	k := 1.0 //# k是截弦的斜率

	//# 快速截弦求根,通常截弦三四次就已达所需精度
	for i := 0; i < 10; i++ {
		k2 := (v2 - v1) / (t2 - t1 - 0.0) //# 算出斜率

		if math.Abs(k2) > 1e-15 {
			k = k2 //# 差商可能为零,应排除
		}

		t = t1 - v1/(k-0.0)
		v := jiaoCai(lx, t, jiao) //# 直线逼近法求根(直线方程的根)
		if v > 1 {
			v -= 2 * math.Pi //# 一次逼近后,v1就已接近0,如果很大,则应减1周
		}

		if math.Abs(v) < 1e-8 {
			break //# 已达精度
		}

		t1 = t2
		v1 = v2
		t2 = t
		v2 = v //# 下一次截弦
	}
	return t
}

func ELon(t float64, n int) float64 {
	return XL0Calc(0, 0, t, n)
}

func MLon(t float64, n int) float64 {
	return XL1Calc(0, t, n)
}

func nutationLon2(t float64) float64 {
	t2 := t * t
	dL := 0.0
	var a float64
	for i := 0; i < len(gNutBL); i += 5 {
		if i == 0 {
			a = -1.742 * t
		} else {
			a = 0
		}

		dL += (gNutBL[i+3] + a) * math.Sin(gNutBL[i]+gNutBL[i+1]*t+gNutBL[i+2]*t2)
	}

	return dL / 100 / SecondsPerRadian
}

func gxcSunlon(t float64) float64 {
	v := -0.043126 + 628.301955*t - 0.000002732*t*t
	e := 0.016708634 - 0.000042037*t - 0.0000001267*t*t
	return (-20.49552 * (1 + e*math.Cos(v))) / SecondsPerRadian
}

func gxcMoonlon(t float64) float64 {
	return -3.4e-6
}

func MsALonT(W float64) float64 {
	v := 7771.37714500204
	t := (W + 1.08472) / v
	t += (W - MsALon(t, 3, 3)) / v
	v = MonthV(t) - EarV(t)
	t += (W - MsALon(t, 20, 10)) / v
	t += (W - MsALon(t, -1, 60)) / v
	return t
}

func MsALonT2(W float64) float64 {
	v := 7771.37714500204
	t := (W + 1.08472) / v
	t2 := t * t
	t -= (-0.00003309*t2 + 0.10976*math.Cos(0.784758+8328.6914246*t+0.000152292*t2) +
		0.02224*math.Cos(0.18740+7214.0628654*t-0.00021848*t2) - 0.03342*math.Cos(4.669257+628.307585*t)) / v
	L := MLon(t, 20) - (4.8950632 + 628.3319653318*t + 0.000005297*t*t + 0.0334166*math.Cos(4.669257+628.307585*t) +
		0.0002061*math.Cos(2.67823+628.307585*t)*t + 0.000349*math.Cos(4.6261+1256.61517*t) - 20.5/SecondsPerRadian)
	v = 7771.38 - 914*math.Sin(0.7848+8328.691425*t+0.0001523*t*t) - 179*math.Sin(2.543+15542.7543*t) - 160*math.Sin(0.1874+7214.0629*t)
	t += (W - L) / v
	return t
}

func MsALon(t float64, Mn, Sn int) float64 {
	return MLon(t, Mn) + gxcMoonlon(t) - (ELon(t, Sn) + gxcSunlon(t) + math.Pi)
}

func SALon(t float64, n int) float64 {
	return ELon(t, n) + nutationLon2(t) + gxcSunlon(t) + math.Pi
}

func EarV(t float64) float64 {
	f := 628.307585 * t
	return 628.332 + 21*math.Sin(1.527+f) + 0.44*math.Sin(1.48+f*2) + 0.129*math.Sin(5.82+f)*t + 0.00055*math.Sin(4.21+f)*t*t
}

func MonthV(t float64) float64 {
	v := 8399.71 - 914*math.Sin(0.7848+8328.691425*t+0.0001523*t*t)
	v -= 179*math.Sin(2.543+15542.7543*t) + 160*math.Sin(0.1874+7214.0629*t) +
		62*math.Sin(3.14+16657.3828*t) + 34*math.Sin(4.827+16866.9323*t) +
		22*math.Sin(4.9+23871.4457*t) + 12*math.Sin(2.59+14914.4523*t) +
		7*math.Sin(0.23+6585.7609*t) + 5*math.Sin(0.9+25195.624*t) +
		5*math.Sin(2.32-7700.3895*t) + 5*math.Sin(3.88+8956.9934*t) + 5*math.Sin(0.49+7771.3771*t)
	return v
}

func SaLonT(W float64) float64 {
	v := 628.3319653318
	t := (W - 1.75347 - math.Pi) / v
	v = EarV(t)
	t += (W - SALon(t, 10)) / v
	v = EarV(t)
	t += (W - SALon(t, -1)) / v

	return t
}

func SAlonT2(W float64) float64 {
	v := 628.3319653318
	t := (W - 1.75347 - math.Pi) / v
	t -= (0.000005297*t*t + 0.0334166*math.Cos(4.669257+628.307585*t) +
		0.0002061*math.Cos(2.67823+628.307585*t)*t) / v
	t += (W - ELon(t, 8) - math.Pi + (20.5+17.2*math.Sin(2.1824-33.75705*t))/SecondsPerRadian) / v
	return t
}

func dtExt(y, jsd float64) float64 {
	dy := (y - 1820) / (100 - 0.0)
	return jsd*dy*dy - 20
}

func dtCalc(y float64) float64 {
	y0 := gDts[len(gDts)-2]
	t0 := gDts[len(gDts)-1]
	if y >= y0 {
		jsd := 31.0
		if y > (y0 + 100) {
			return dtExt(y, jsd)
		}

		v := dtExt(y, jsd)
		dv := dtExt(y0, jsd) - t0
		return v - dv*(y0+100-y)/(100-0.0)
	}

	i := 0
	for y >= gDts[i+5] {
		i += 5
	}

	t1 := (y - gDts[i]) / (gDts[i+5] - gDts[i] - 0.0) * 10
	t2 := t1 * t1
	t3 := t2 * t1
	return gDts[i+1] + gDts[i+2]*t1 + gDts[i+3]*t2 + gDts[i+4]*t3
}

func dtT2(jd float64) float64 {
	return dtCalc(jd/365.2425+2000) / 86400.0
}

func soLow(W float64) float64 {
	v := 7771.37714500204
	t := (W + 1.08472) / v
	t -= (-0.0000331*t*t+0.10976*math.Cos(0.785+8328.6914*t)+0.02224*math.Cos(0.187+7214.0629*t)-
		0.03342*math.Cos(4.669+628.3076*t))/v + (32*(t+1.8)*(t+1.8)-20)/86400/36525
	return t*36525 + 8.0/24.0
}

func qiLow(W float64) float64 {
	v := 628.3319653318
	t := (W - 4.895062166) / v
	t -= (53*t*t + 334116*math.Cos(4.67+628.307585*t) + 2061*math.Cos(2.678+628.3076*t)*t) / v / 10000000
	L := 48950621.66 + 6283319653.318*t + 53*t*t + 334166*math.Cos(4.669257+628.307585*t) +
		3489*math.Cos(4.6261+1256.61517*t) + 2060.6*math.Cos(2.67823+628.307585*t)*t -
		994 - 834*math.Sin(2.1824-33.75705*t)

	t -= (L/10000000-W)/628.332 + (32*(t+1.8)*(t+1.8)-20)/86400/36525
	return t*36525 + 8.0/24.0
}

func qiHigh(W float64) float64 {
	t := SAlonT2(W) * 36525
	t = t - dtT2(t) + 8.0/24.0
	t0 := math.Mod(t+0.5, 1)
	//fmt.Println("t = %d ,t0 = %d ", t, t0)

	v := t0 * 86400
	if v < 1200 || v > (86400-1200) {
		t = SaLonT(W)*36525 - dtT2(t) + 8.0/24.0
	}

	return t
}

func soHigh(W float64) float64 {
	t := MsALonT2(W) * 36525
	t = t - dtT2(t) + 8.0/24
	t0 := math.Mod(t+0.5, 1)
	//fmt.Println("t =%d ,t0 =%d ", t, t0)
	v := t0 * 86400
	if v < 1800 || v > (86400-1800) {
		t = MsALonT(W)*36525 - dtT2(t) + 8.0/24
	}
	return t
}

func qiAccurate(W float64, astflg bool, L float64) float64 {
	d := SaLonT(W)
	t := d * 36525
	if astflg {
		return t - dtT2(t) + mstAst(d) + L/360.0
	} else {
		return t - dtT2(t) + 8/(24-0.0)
	}
}

func qiAccurate2(jd float64, astflg bool, L float64) float64 {
	d := math.Pi / 12
	w := math.Floor((jd+293)/365.2422*24) * d

	a := qiAccurate(w, astflg, L)
	if (a - jd) > 5 {
		return qiAccurate(w-d, astflg, L)
	} else if (a - jd) < -5 {
		return qiAccurate(w+d, astflg, L)
	} else {
		return a
	}
}
func dqCalc(jd float64) (float64, int64, int64, int64, float64) {
	year, month, day, gd := jd2gcal(MJD_0, jd-MJD_0)
	tmp := jd - float64(day) + 4 - JulianCalendarJ2000
	D := qiAccurate2(tmp, false, 120) + JulianCalendarJ2000
	return D, year, month, day, gd
}

func bkCalc(year, month, lifa int) (float64, float64, float64, float64) {
	M := month
	var Y int
	if lifa == 12 {
		Y = year - 1
		M = 12
	} else if lifa == 11 && M < 7 {
		Y = year - 1
		M = 6
	} else if lifa == 11 && M >= 7 {
		Y = year
		M = 6
	}

	res1, res2 := gcal2jd(int64(Y), int64(M), 21)
	stjd := res1 + res2 + 0.5 - JulianCalendarJ2000

	res1, res2 = gcal2jd(int64(Y+1), int64(M), 21)
	mdjd := res1 + res2 + 0.5 - JulianCalendarJ2000

	res1, res2 = gcal2jd(int64(Y+2), int64(M), 21)
	edjd := res1 + res2 + 0.5 - JulianCalendarJ2000
	pzq := qiAccurate2(stjd, false, 120) + JulianCalendarJ2000
	zq := qiAccurate2(mdjd, false, 120) + JulianCalendarJ2000
	nzq := qiAccurate2(edjd, false, 120) + JulianCalendarJ2000

	k1 := (zq - pzq) / 12
	b1 := pzq + k1/2
	k2 := (nzq - zq) / 12
	b2 := zq + k2/2

	return b1, k1, b2, k2
}

func calcLichun(bzTime *BZDateTime) float64 {
	year := bzTime.Year
	res1, res2 := gcal2jd(int64(year), 2, 2)
	jd := res1 + res2 + 0.5 - JulianCalendarJ2000
	lichun := qiAccurate2(jd, false, 120) + JulianCalendarJ2000
	if bzTime.ToJD() < lichun {
		lichun -= 365.2422
	}
	return lichun
}

func calcByShuoOrQi(jd, pc float64, array []float64) int {
	var i int
	for i = 0; i < len(array); i += 2 {
		if (jd + pc) < array[i+2] {
			break
		}
	}
	D0 := math.Floor((jd + pc - array[i]) / array[i+1])

	D1 := array[i] + array[i+1]*float64(int2(D0))
	D := int(D1 + 0.5)

	if D == 1683460 {
		D += 1
	}

	return D
}

func calc(jd float64, qs string) float64 {
	var pc, s1, s2, s3 float64
	if qs == "s" {
		pc = 14
		s1 = gShuoKB[0]
		s2 = gShuoKB[len(gShuoKB)-1]
	} else {
		pc = 7
		s1 = gQiKB[0]
		s2 = gQiKB[len(gQiKB)-1]
	}
	s1 -= pc
	s2 -= pc
	s3 = 2436935
	if jd < s1 || jd >= s3 {
		if qs == "s" {
			return JulianCalendarJ2000 + math.Floor(soHigh(math.Floor((jd+pc-2451551)/29.5306)*math.Pi*2)+0.5)
		} else {
			return JulianCalendarJ2000 + math.Floor(qiHigh(math.Floor((jd+pc-2451259)/365.2422*24)*math.Pi/12)+0.5)
		}
	} else if s1 <= jd && jd < s2 {

		if qs == "s" {
			D := calcByShuoOrQi(jd, pc, gShuoKB)
			return float64(D)
		} else {
			D := calcByShuoOrQi(jd, pc, gQiKB)
			return float64(D)
		}

	} else if s2 <= jd && jd < s3 {
		var D float64
		var n string
		if qs == "s" {
			D0 := int2(soLow(float64(int2((jd+pc-2451551)/29.5306))*math.Pi*2) + 0.5)
			n0 := gSTbl[int2((jd-s2)/29.5306)]
			D = float64(D0)
			n = string(n0)
		} else {
			D0 := int2(qiLow(float64(int2((jd+pc-2451259)/365.2422*24))*math.Pi/12) + 0.5)
			n0 := gQTbl[int2((jd-s2)/365.2422*24)]
			D = float64(D0)
			n = string(n0)
		}

		if n == "1" {
			D += 1
		}

		if n == "2" {
			D -= 1
		}

		return D + JulianCalendarJ2000

	}
	return -1
}

type SolarToLunarResS struct {
	Year       int
	MonthIdx   int
	DayIdx     int
	IsLeapM    bool
	NextM      bool
	YearGanZhi string
	CCMonth    string
	CCDay      string
}

func Solar2lunar(year, month, day, hour int) (res SolarToLunarResS, err error) {
	ret := SolarToLunarResS{
		Year:       0,
		MonthIdx:   0,
		DayIdx:     0,
		IsLeapM:    false,
		NextM:      false,
		YearGanZhi: "",
		CCMonth:    "",
		CCDay:      "",
	}
	Y := int64(year)
	M := int64(month)
	D := int64(day)
	dts1 := DateTimeS{
		Year:   year,
		Month:  month,
		Day:    day,
		Hour:   0,
		Minute: 0,
		Second: 0,
	}
	dtsTime := DateTimeByDateTimeS(dts1)
	var jd float64
	if Y < 1582 {
		res1, res2 := jcal2jd(Y, M, D)
		jd = res1 + res2 + 0.5
	} else if Y > 1582 {
		res1, res2 := gcal2jd(Y, M, D)
		jd = res1 + res2 + 0.5
	} else {
		gcalSt := DateTimeByStr(JulianCalendar1582Spec)
		if dtsTime >= gcalSt {
			jcRes1, jcRes2 := gcal2jd(Y, M, D)
			jd = jcRes1 + jcRes2 + 0.5
		} else {
			jdRes1, jdRes2 := jcal2jd(Y, M, D)
			jd = jdRes1 + jdRes2 + 0.5
		}
	}

	//#子时换日
	if hour >= 23 {
		jd += 1
	}

	W := float64(int2((jd-JulianCalendarJ2000-355+183)/365.2422))*365.2422 + 355 + JulianCalendarJ2000
	if calc(W, "q") > jd {
		W -= 365.2422
	}

	var A, B []float64
	for i := 0; i < 13; i++ {
		A = append(A, calc(W+30.4368*float64(i), "q"))
	}

	w := calc(A[0], "s")
	if w > A[0] {
		w -= 29.53
	}
	for i := 0; i < 15; i++ {
		B = append(B, calc(w+29.5306*float64(i), "s"))
	}

	leap := 0
	var dx []float64
	var ym []int
	var sm []bool
	var mcc []string
	for i := 0; i < 14; i++ {
		dx = append(dx, B[i+1]-B[i])
		ym = append(ym, i)
		sm = append(sm, false)
		mcc = append(mcc, "")
	}

	if B[13] <= A[12] {
		i := 1
		for B[i+1] > A[i] && i < 13 {
			i += 1
		}
		leap = i
		for i < 14 {
			ym[i] -= 1
			i += 1
		}
	}

	//("冬", "腊", "正", "二", "三", "四", "五", "六", "七", "八", "九", "十")
	//mc := []int{11, 12, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1}
	for i := 0; i < 14; i++ {
		idx := ym[i] % 12

		if 1724360 <= B[i] && B[i] <= 1729794 {
			idx = (ym[i] + 1) % 12
		} else if 1807724 <= B[i] && B[i] <= 1808699 {
			idx = (ym[i] + 1) % 12
		} else if 1999349 <= B[i] && B[i] <= 1999467 {
			idx = (ym[i] + 2) % 12
		} else if 1973067 <= B[i] && B[i] <= 1977052 {
			if ym[i]%12 == 0 {
				idx = 2
			}
			if ym[i] == 2 {
				idx = 12
				sm[i] = true
			}
		}
		if B[i] == 1729823 || B[i] == 1808729 {
			idx = 1
			sm[i] = true
		}

		if B[i] == 1999497 || B[i] == 1999526 || B[i] == 1977112 {
			sm[i] = true
		}
		mcc[i] = gYmc[idx]
	}

	i := 0
	for i <= 14 && B[i] <= jd {
		i += 1
	}

	if i >= 15 {
		return ret, errors.New("没有这个月")
	}

	i -= 1
	lm := mcc[i]
	j, isExist := IndexPlus("正", mcc, 0, -1)
	if !isExist {
		return ret, errors.New("月序有误，没有正月")
	}

	lyBase := int2((B[j]-JulianCalendarJ2000+5810)/365.2422+0.5) + 1984
	var ly int64
	var lygan, lyzhi string
	if jd < B[j] {
		ly = ((lyBase - 1 - int64(gBaseYear)) + 57) % 60
		lygan = gTianGan[int2(float64(ly%10))]
		lyzhi = gDiZhi[int2(float64(ly%12))]
		ly = lyBase - 1
	} else {
		//闰正月[leap January of Chinese calender]:
		// 公元689年的月序: ['冬', '腊', '正', '二', '三', '四', '五', '六', '七', '八', '九', '九', '十', '正']
		//公元690年的月序： ['正', '腊', '一', '二', '三', '四', '五', '六', '七', '八', '九', '十', '正', '腊']
		//闰正月的问题 所有从正月后面第二个月开始
		z, isExist := IndexPlus("正", mcc, j+2, -1)
		if isExist && jd >= B[z] {
			//这里查不到是正常的，只有很特殊的像689，690这样的历法改制
			lyBase += 1
		}

		ly = ((lyBase - int64(gBaseYear)) + 57) % 60
		lygan = gTianGan[int2(float64(ly%10))]
		lyzhi = gDiZhi[int2(float64(ly%12))]
		ly = lyBase
	}

	ld := gRmc[int2(jd-B[i])]
	if leap > 0 && leap == i {
		ret.IsLeapM = true
		fmt.Printf("\"----{%d} {%s}{%s}年 闰{%s}月 {%s}\n", ly, lygan, lyzhi, lm, ld)
	} else {
		ret.IsLeapM = false
		fmt.Printf("\"----{%d} {%s}{%s}年 {%s}月 {%s}\n", ly, lygan, lyzhi, lm, ld)
	}

	if sm[i] {
		fmt.Printf("特殊年份有两个{%s}月，此为后一个{%s}月)\n", lm, lm)
	}

	ret.Year = int(ly)
	ret.MonthIdx, _ = Index(lm, gYmc)
	ret.DayIdx = int(int2(jd - B[i]))
	ret.NextM = sm[i]
	ret.YearGanZhi = lygan + lyzhi
	ret.CCMonth = lm
	ret.CCDay = ld

	return ret, nil
}

func Lunar2solar(year, month, day int, isLeapMonth, next bool) (int, int, int, error) {
	Y := year
	M := month
	D := day
	stdM := 0
	if (9 <= Y && Y <= 23) || (237 <= Y && Y <= 239) {
		stdM = 11
	} else if 690 <= Y && Y <= 699 {
		stdM = 12
	} else {
		stdM = 10
	}

	var jd float64
	var isHead bool
	if (0 < M && M <= stdM) || (Y == 700 && M == 12 && next == false) {
		res1, res2 := gcal2jd(int64(Y), 2, 4)
		jd = res1 + res2 + 0.5
		isHead = true
	} else if stdM < M && M < 13 {
		res1, res2 := gcal2jd(int64(Y+1), 2, 4)
		jd = res1 + res2 + 0.5
		isHead = false
	} else {
		return 0, 0, 0, errors.New("无效的月份")
	}

	W := float64(int2((jd-JulianCalendarJ2000-355+183)/365.2422))*365.2422 + 355 + JulianCalendarJ2000
	if calc(W, "q") > jd {
		W -= 365.2422
	}

	var A, B []float64
	for i := 0; i < 13; i++ {
		A = append(A, calc(W+30.4368*float64(i), "q"))
	}

	w := calc(A[0], "s")
	if w > A[0] {
		w -= 29.53
	}

	for i := 0; i < 15; i++ {
		B = append(B, calc(w+29.5306*float64(i), "s"))
	}

	leap := 0
	var dx []float64
	var ym []int
	for i := 0; i < 14; i++ {
		dx = append(dx, B[i+1]-B[i])
		ym = append(ym, i)
	}

	if B[13] <= A[12] {
		i := 1
		for B[i+1] > A[i] && i < 13 {
			i += 1
		}
		leap = i
		for i < 14 {
			ym[i] -= 1
			i += 1
		}
	}

	mc := []int{11, 12, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := 0; i < 14; i++ {
		idx := ym[i] % 12

		if 1724360 <= B[i] && B[i] <= 1729794 {
			idx = (ym[i] + 1) % 12
		} else if 1807724 <= B[i] && B[i] <= 1808699 {
			idx = (ym[i] + 1) % 12
		} else if 1999349 <= B[i] && B[i] <= 1999467 {
			idx = (ym[i] + 2) % 12
		} else if 1973067 <= B[i] && B[i] <= 1977052 {
			if ym[i]%12 == 0 {
				idx = 2
			}
			if ym[i] == 2 {
				idx = 2
			}
		}
		if B[i] == 1729794 || B[i] == 1808699 {
			idx = 1
		}

		ym[i] = mc[idx]
	}

	var jdD float64
	var mHeadIdx int
	if isLeapMonth {
		if leap > 0 && M == ym[leap] {
			if 0 < D && float64(D) <= dx[leap] {
				jdD = B[leap] + float64(D) - 1
			} else {
				return 0, 0, 0, errors.New(fmt.Sprintf("该月没有%d日", D))
			}

		} else {
			return 0, 0, 0, errors.New(fmt.Sprintf("%d月不是闰月", M))
		}
	} else {
		if isHead {
			start, isExist := Index(1, ym)
			if !isExist {
				return 0, 0, 0, errors.New("月序有误，没有1月")
			}

			mHeadIdx, isExist = IndexPlus(M, ym, start, -1)
			if !isExist {
				if next {
					msg := fmt.Sprintf("该年没有第二个%d月\n", M)
					return 0, 0, 0, errors.New(msg)
				} else {
					msg := fmt.Sprintf("该年没有%d月", M)
					return 0, 0, 0, errors.New(msg)
				}
			}

			if next {
				mHeadIdx, isExist = IndexPlus(M, ym, mHeadIdx+1, -1)
				if !isExist {
					msg := fmt.Sprintf("该年没有第二个%d月\n", M)
					return 0, 0, 0, errors.New(msg)
				}
			}

			if 0 < D && float64(D) <= dx[mHeadIdx] {
				jdD = B[mHeadIdx] + float64(D) - 1
			} else {
				return 0, 0, 0, errors.New(fmt.Sprintf("该月没有%d日", D))
			}

		} else {
			end, isExist := Index(1, ym)
			if !isExist {
				return 0, 0, 0, errors.New("月序有误，没有1月")
			}
			//#农历700年特殊，头尾各有一个腊月
			mHeadIdx, isExist = IndexPlus(M, ym, 0, end)
			if !isExist {
				if next {
					msg := fmt.Sprintf("该年没有第二个%d月\n", M)
					return 0, 0, 0, errors.New(msg)
				} else {
					msg := fmt.Sprintf("该年没有%d月", M)
					return 0, 0, 0, errors.New(msg)
				}
			}
			if next && Y != 700 {
				mHeadIdx, isExist = IndexPlus(M, ym, mHeadIdx+1, -1)
				if !isExist {
					msg := fmt.Sprintf("该年没有第二个%d月\n", M)
					return 0, 0, 0, errors.New(msg)
				}
			}

			if 0 < D && float64(D) <= dx[mHeadIdx] {
				jdD = B[mHeadIdx] + float64(D) - 1
			} else {
				return 0, 0, 0, errors.New(fmt.Sprintf("该月没有%d日", D))
			}

		}
	}

	var res1, res2, res3 int64
	if jdD >= 2299160.5 {
		res1, res2, res3, _ = jd2gcal(MJD_0, jdD-MJD_0)
	} else {
		res1, res2, res3, _ = jd2jcal(MJD_0, jdD-MJD_0)
	}

	return int(res1), int(res2), int(res3), nil
}
