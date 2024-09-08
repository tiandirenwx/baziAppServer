package bazicore

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(test *testing.T) {
	dateStr := "20160524"
	timeStr := "15:04:01.123"
	dtStr := "19950101-03:04:59.999"

	d := DateByStr(dateStr)
	t := TimeByMSStr(timeStr)
	dt := DateTimeByStr(dtStr)

	if d != 20160524 {
		test.Log("d is ", d)
		test.FailNow()
	}

	if t != 150401123 {
		test.Log("t is ", t)
		test.FailNow()
	}

	if dt != 19950101030459999 {
		test.Log("dt is ", dt)
		test.FailNow()
	}

	dt2 := d.ToDateTime()
	d2 := dt.Date()
	t2 := dt.Time()

	if dt2 != 20160524000000000 {
		test.Log("dt2 is ", dt2)
		test.FailNow()
	}

	if d2 != 19950101 {
		test.Log("d2 is ", d2)
		test.FailNow()
	}

	if t2 != 30459999 {
		test.Log("t2 is ", t2)
		test.FailNow()
	}

	dt.SetDate(d)
	if dt != 20160524030459999 {
		test.Log("dt3 is ", dt)
		test.FailNow()
	}

	dt.SetTime(t)

	if dt != 20160524150401123 {
		test.Log("dt4 is ", dt)
		test.FailNow()
	}

	dt = dt.SubDateTime(1120101998)
	if dt != 20160523030259125 {
		test.Log("dt5 is ", dt)
		test.FailNow()
	}
}

func Test_FormatDateTime(test *testing.T) {
	fmt.Println(FormatDateTime(time.Now(), HH_MM))
	fmt.Println(FormatDateTime(time.Now(), HH_MM_SS))
	fmt.Println(FormatDateTime(time.Now(), HH_MM_SS_MS))
	fmt.Println(FormatDateTime(time.Now(), MM_DD))
	fmt.Println(FormatDateTime(time.Now(), MM_DD_CN))
	fmt.Println(FormatDateTime(time.Now(), MM_DD_EN))
	fmt.Println(FormatDateTime(time.Now(), MM_DD_HH_MM))
	fmt.Println(FormatDateTime(time.Now(), MM_DD_HH_MM_CN))
	fmt.Println(FormatDateTime(time.Now(), MM_DD_HH_MM_EN))
	fmt.Println(FormatDateTime(time.Now(), MM_DD_HH_MM_SS))
	fmt.Println(FormatDateTime(time.Now(), MM_DD_HH_MM_SS_CN))
	fmt.Println(FormatDateTime(time.Now(), MM_DD_HH_MM_SS_EN))
	fmt.Println(FormatDateTime(time.Now(), YYMMDDHHMM))
	fmt.Println(FormatDateTime(time.Now(), YYYYMM))
	fmt.Println(FormatDateTime(time.Now(), YYYYMMDD))
	fmt.Println(FormatDateTime(time.Now(), YYYYMMDDHH))
	fmt.Println(FormatDateTime(time.Now(), YYYYMMDDHHMM))
	fmt.Println(FormatDateTime(time.Now(), YYYYMMDDHHMMSS))
	fmt.Println(FormatDateTime(time.Now(), YYYYMMDD_HH_MM_SS))
	fmt.Println(FormatDateTime(time.Now(), YYYYMMDD_HH_MM_SS_SSS))
	fmt.Println(FormatDateTime(time.Now(), YYYY_MM))
	fmt.Println(FormatDateTime(time.Now(), YYYY_MM_CN))
	fmt.Println(FormatDateTime(time.Now(), YYYY_MM_DD))
	fmt.Println(FormatDateTime(time.Now(), YYYY_MM_DD_CN))
	fmt.Println(FormatDateTime(time.Now(), YYYY_MM_DD_EN))
	fmt.Println(FormatDateTime(time.Now(), YYYY_MM_DD_HH_MM))
	fmt.Println(FormatDateTime(time.Now(), YYYY_MM_DD_HH_MM_CN))
	fmt.Println(FormatDateTime(time.Now(), YYYY_MM_DD_HH_MM_EN))
	fmt.Println(FormatDateTime(time.Now(), YYYY_MM_DD_HH_MM_SS))
	fmt.Println(FormatDateTime(time.Now(), YYYY_MM_DD_HH_MM_SS_CN))
	fmt.Println(FormatDateTime(time.Now(), YYYY_MM_DD_HH_MM_SS_EN))
	fmt.Println(FormatDateTime(time.Now(), YYYY_MM_DD_HH_MM_SS_SSS))
	fmt.Println(FormatDateTime(time.Now(), YYYY_MM_DD_HH_MM_SS_SSS_EN))
	fmt.Println(FormatDateTime(time.Now(), YYYY_MM_EN))

	fmt.Println(FormatDateTime(time.Now(), "yyyy/MM/dd HH:mm:ss.SSSSSSSSS"))
	fmt.Println(FormatDateTime(time.Now(), "2006 01 02 15:04:05.00000000"))
}

func Test_DateTimeByDateTimes(test *testing.T) {
	nowStr := FormatDateTime(time.Now(), YYYYMMDD_HH_MM_SS_SSS_EN)
	goTime := DateTimeByStr(nowStr)
	dts, _ := DateTimeSByStr(nowStr)
	dtsTime := DateTimeByDateTimeS(dts)
	fmt.Printf(" now time form FormatDateTime nowStr: %s \n", nowStr)
	fmt.Printf(" now time from DateTimeByStr, goTime: %d \n", goTime)
	fmt.Printf(" now time from DateTimeByDateTimeS, dtsTime: %d \n", dtsTime)
}
