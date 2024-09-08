package bazicore

import (
	"fmt"
	"math"
	"strconv"
)

type BZDateTime struct {
	Year    int
	Month   int
	Day     int
	Hour    int
	Minute  int
	Second  int
	BzJ2000 int64
	BzDts   []float64
	IsGcal  bool
	JdD     float64
	JdF     float64
}

// /* yyyyMMdd-HH:mm:SS.sss */
func (bzt *BZDateTime) New(timeStr string) *BZDateTime {
	dt, err := DateTimeSByStr(timeStr)
	if err != nil {
		return nil
	}
	goTime := DateTimeByStr(timeStr)
	var jdD0 float64 = 0
	var isGcal = false
	year := int64(dt.Year)
	month := int64(dt.Month)
	day := int64(dt.Day)
	if dt.Year < 1582 {
		jdRes1, jdRes2 := jcal2jd(year, month, day)
		jdD0 = jdRes1 + jdRes2 + 0.5
		isGcal = false
	} else if dt.Year > 1582 {
		jcRes1, jcRes2 := gcal2jd(year, month, day)
		jdD0 = jcRes1 + jcRes2 + 0.5
		isGcal = true
	} else {
		gcalSt := DateTimeByStr(JulianCalendar1582Spec)
		if goTime >= gcalSt {
			jcRes1, jcRes2 := gcal2jd(year, month, day)
			jdD0 = jcRes1 + jcRes2 + 0.5
			isGcal = true
		} else {
			jdRes1, jdRes2 := jcal2jd(year, month, day)
			jdD0 = jdRes1 + jdRes2 + 0.5
			isGcal = false
		}
	}
	jdD := math.Floor(jdD0)
	jdF := (float64(dt.Hour)+float64(dt.Minute)/60.0+float64(dt.Second)/3600.0)/24.0 - 0.5
	bzDt := &BZDateTime{
		Year:    dt.Year,
		Month:   dt.Month,
		Day:     dt.Day,
		Hour:    dt.Hour,
		Minute:  dt.Minute,
		Second:  dt.Second,
		BzJ2000: JulianCalendarJ2000, //2000年前儒略日数(2000-1-1 12:00:00格林威治平时
		BzDts:   gDts,
		IsGcal:  isGcal,
		JdD:     jdD,
		JdF:     jdF,
	}

	return bzDt
}

func (bzt *BZDateTime) ToJD() float64 {
	return float64(bzt.JdD) + bzt.JdF
}

func (bzt *BZDateTime) SetFromJD(jd float64) {
	D := math.Floor(jd + 0.5)
	F := jd + 0.5 - D
	bzt.JdD = D
	bzt.JdF = F - 0.5
	if D >= 2299161 {
		year, month, day, _ := jd2gcal(MJD_0, jd-MJD_0)
		bzt.Year = int(year)
		bzt.Month = int(month)
		bzt.Day = int(day)
		bzt.IsGcal = true
	} else {
		year, month, day, _ := jd2jcal(MJD_0, jd-MJD_0)
		bzt.Year = int(year)
		bzt.Month = int(month)
		bzt.Day = int(day)
		bzt.IsGcal = false
	}

	F *= 24
	bzt.Hour = int(math.Floor(F))
	F -= float64(bzt.Hour)
	F *= 60
	bzt.Minute = int(math.Floor(F))
	F -= float64(bzt.Minute)
	F *= 60
	bzt.Second = int(F)
}

func (bzt *BZDateTime) ToLunarDate() (SolarToLunarResS, error) {
	Y := bzt.Year
	M := bzt.Month
	D := bzt.Day
	h := bzt.Hour

	res, err := Solar2lunar(Y, M, D, h)
	if err != nil {
		fmt.Printf("error msg: %s\n", err)
		return res, err
	}

	return res, nil
}

func (bzt *BZDateTime) GetLunarInfo() (string, string, string) {
	res, err := Solar2lunar(bzt.Year, bzt.Month, bzt.Day, bzt.Hour)
	if err != nil {
		fmt.Printf("error msg: %s\n", err)
		//return res, err
	}
	lySx := gShengXiao[(res.Year-1984+90000)%12]
	lyNh := bzt.getNianHao(res.Year)
	lyXh := ((res.Year - gBaseYear) + 57) % 60
	lyGan := gTianGan[(lyXh % 10)]
	lyZhi := gDiZhi[(lyXh % 12)]
	lmMc := gYmc[res.MonthIdx]
	ldMc := gRmc[res.DayIdx]

	leapMsg := ""
	if res.IsLeapM {
		leapMsg = "闰"
	}
	specMsg := ""
	if res.NextM {
		s := fmt.Sprintf("%s", lmMc)
		specMsg = fmt.Sprintf(" 〖有两个%s月，此为后一个〗", s)
	}

	retMsg := fmt.Sprintf("农历%s%s年%s%s月%s %s", lyGan, lyZhi, leapMsg, lmMc, ldMc, specMsg)

	return lyNh, lySx, retMsg
}

func (bzt *BZDateTime) getNianHao(Y int) string {
	msg := ""
	for i := 0; i < len(gJNB); i += 7 {
		n, _ := strconv.Atoi(gJNB[i])
		m, _ := strconv.Atoi(gJNB[i+1])
		if Y < n || Y >= n+m {
			continue
		}
		p, _ := strconv.Atoi(gJNB[i+2])
		r := gJNB[i+3]
		s := gJNB[i+6]

		c := fmt.Sprintf("%s%d年", s, Y-n+1+p)
		if len(msg) == 0 {
			msg += "" + fmt.Sprintf("%s--%s", r, c)
		} else {
			msg += ";" + fmt.Sprintf("%s--%s", r, c)
		}
	}
	return msg
}

func (bzt *BZDateTime) GetSolarDaysFromBase() float64 {
	return bzt.JdD - float64(gBaseJD)
}

func (bzt *BZDateTime) dtExt(y, jsd float64) float64 {
	dy := (y - 1820) / (100 - 0.0)
	return jsd*dy*dy - 20
}

func (bzt *BZDateTime) dtCalc(y float64) float64 {
	y0 := bzt.BzDts[len(bzt.BzDts)-2]
	t0 := bzt.BzDts[len(bzt.BzDts)-1]

	var jsd float64
	if y >= y0 {
		jsd = 31
		if y > (y0 + 100) {
			return bzt.dtExt(y, jsd)
		}
		v := bzt.dtExt(y, jsd)
		dv := bzt.dtExt(y0, jsd) - t0
		return v - dv*(y0+100-y)/(100-0.0)
	}

	i := 0
	for {
		if (i + 5) >= len(bzt.BzDts) {
			break
		}

		if y >= bzt.BzDts[i+5] {
			i += 5
		} else {
			break
		}
	}

	t1 := (y - bzt.BzDts[i]) / (bzt.BzDts[i+5] - bzt.BzDts[i] - 0.0) * 10
	t2 := t1 * t1
	t3 := t2 * t1
	return bzt.BzDts[i+1] + bzt.BzDts[i+2]*t1 + bzt.BzDts[i+3]*t2 + bzt.BzDts[i+4]*t3
}

// 计算世界时与原子时之差,传入年
func (bzt *BZDateTime) deltatT(y float64) float64 {
	d := bzt.BzDts

	i := 0
	for i < 100 {
		if y < d[i+5] || i == 95 {
			break
		}
		i = i + 5
	}

	t1 := (y - d[i]) / (d[i+5] - d[i] - 0.0) * 10
	t2 := t1 * t1
	t3 := t2 * t1
	return d[i+1] + d[i+2]*t1 + d[i+3]*t2 + d[i+4]*t3
}

func (bzt *BZDateTime) dtT2(jd float64) float64 {
	jd -= float64(bzt.BzJ2000)
	return bzt.dtCalc(jd/365.2425+2000) / 86400.0
}

func (bzt *BZDateTime) deltatT2(jd float64) float64 {
	return bzt.deltatT(jd/365.2425+2000) / 86400.0
}

func (bzt *BZDateTime) getMeridiem() string {
	hm := bzt.Hour*100 + bzt.Minute
	if hm < 600 {
		return "凌晨"
	} else if hm < 900 {
		return "早上"
	} else if hm < 1130 {
		return "上午"
	} else if hm < 1230 {
		return "中午"
	} else if hm < 1800 {
		return "下午"
	} else {
		return "晚上"
	}
}

func (bzt *BZDateTime) getFormatStr() string {
	Y := bzt.Year
	M := bzt.Month
	D := bzt.Day
	h := bzt.Hour
	if bzt.Hour < 13 {
		h = bzt.Hour
	} else {
		h = bzt.Hour - 12
	}

	m := bzt.Minute
	md := bzt.getMeridiem()
	return fmt.Sprintf("%04d年%02d月%02d日 %s%02d点%02d", Y, M, D, md, h, m)
}

func (bzt *BZDateTime) toStr() string {
	return fmt.Sprintf("%d年%d月%d日%d时%d分",
		bzt.Year, bzt.Month, bzt.Day, bzt.Hour, bzt.Minute)
}

func (bzt *BZDateTime) toRepr() string {
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d}",
		bzt.Year, bzt.Month, bzt.Day, bzt.Hour, bzt.Minute, bzt.Second)
}
func (bzt *BZDateTime) toJDOld(UTC bool) float64 {
	y := bzt.Year
	m := bzt.Month
	var n int64 = 0 //取出年月
	if m <= 2 {
		m += 12
		y -= 1
	}

	//判断是否为格里高利历日1582*372+10*31+15
	if bzt.Year*372+bzt.Month*31+bzt.Day >= 588829 {
		n = int2((float64(y) / 100))
		n = 2 - n + int2((float64(n) / 4)) //加百年闰
	}

	n += int2(365.2500001 * (float64(y) + 4716))  //加上年引起的偏移日数
	n += int2(30.6*float64(m+1)) + int64(bzt.Day) //加上月引起的偏移日数及日偏移数

	nRet := ((float64(bzt.Second)/(60-0.0)+float64(bzt.Minute))/(60-0.0)+float64(bzt.Hour))/(24-0.0) - 1524.5
	nRet += float64(n)

	if UTC {
		return nRet + bzt.dtT2(nRet-float64(bzt.BzJ2000))
	}
	return nRet
}

/*
算出:jd转到当地UTC后,UTC日数的整数部分或小数部分
基于J2000力学时jd的起算点是12:00:00时,所以跳日时刻发生在12:00:00,这与日历计算发生矛盾
把jd改正为00:00:00起算,这样儒略日的跳日动作就与日期的跳日同步
改正方法为jd=jd+0.5-deltatT+shiqu/24
把儒略日的起点移动-0.5(即前移12小时)
式中shiqu是时区,北京的起算点是-8小时,shiqu取8
*/
func (bzt *BZDateTime) Dint_dec(jd, shiqu float64) (int64, float64) {
	u := jd + 0.5 - bzt.dtT2(jd) + shiqu/(24-0.0)
	u1 := int64(math.Floor(u))
	u2 := u - float64(u1)
	return u1, u2
}
