package bazicore

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	GolangTimeFormat = "20060102-15:04:05.000" //go的诞生时间
)

type DateStyle string

const (
	MM_DD                    = "MM-dd"
	YYYYMM                   = "yyyyMM"
	YYYY_MM                  = "yyyy-MM"
	YYYY_MM_DD               = "yyyy-MM-dd"
	YYYYMMDD                 = "yyyyMMdd"
	YYYYMMDDHHMMSS           = "yyyyMMddHHmmss"
	YYYYMMDD_HH_MM_SS        = "yyyyMMdd HH:mm:ss"
	YYYYMMDD_HH_MM_SS_SSS    = "yyyyMMdd HH:mm:ss.SSS"
	YYYYMMDD_HH_MM_SS_SSS_EN = "yyyyMMdd-HH:mm:ss.SSS"
	YYYYMMDDHHMM             = "yyyyMMddHHmm"
	YYYYMMDDHH               = "yyyyMMddHH"
	YYMMDDHHMM               = "yyMMddHHmm"
	MM_DD_HH_MM              = "MM-dd HH:mm"
	MM_DD_HH_MM_SS           = "MM-dd HH:mm:ss"
	YYYY_MM_DD_HH_MM         = "yyyy-MM-dd HH:mm"
	YYYY_MM_DD_HH_MM_SS      = "yyyy-MM-dd HH:mm:ss"
	YYYY_MM_DD_HH_MM_SS_SSS  = "yyyy-MM-dd HH:mm:ss.SSS"

	MM_DD_EN                   = "MM/dd"
	YYYY_MM_EN                 = "yyyy/MM"
	YYYY_MM_DD_EN              = "yyyy/MM/dd"
	MM_DD_HH_MM_EN             = "MM/dd HH:mm"
	MM_DD_HH_MM_SS_EN          = "MM/dd HH:mm:ss"
	YYYY_MM_DD_HH_MM_EN        = "yyyy/MM/dd HH:mm"
	YYYY_MM_DD_HH_MM_SS_EN     = "yyyy/MM/dd HH:mm:ss"
	YYYY_MM_DD_HH_MM_SS_SSS_EN = "yyyy/MM/dd HH:mm:ss.SSS"

	MM_DD_CN               = "MM月dd日"
	YYYY_MM_CN             = "yyyy年MM月"
	YYYY_MM_DD_CN          = "yyyy年MM月dd日"
	MM_DD_HH_MM_CN         = "MM月dd日 HH:mm"
	MM_DD_HH_MM_SS_CN      = "MM月dd日 HH:mm:ss"
	YYYY_MM_DD_HH_MM_CN    = "yyyy年MM月dd日 HH:mm"
	YYYY_MM_DD_HH_MM_SS_CN = "yyyy年MM月dd日 HH:mm:ss"

	HH_MM       = "HH:mm"
	HH_MM_SS    = "HH:mm:ss"
	HH_MM_SS_MS = "HH:mm:ss.SSS"
)

type DateTimeS struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
	Second int
}
type DateTime uint64
type Date uint64
type Time uint64

/* yyyyMMdd-HH:mm:SS.sss */
func DateTimeStrByBazi(year, month, day, hour, min int) string {
	s := fmt.Sprintf("%04d%02d%02d-%02d:%02d", year, month, day, hour, min)
	return s + ":00.000"
}

// 日期转字符串
func FormatDateTime(date time.Time, dateStyle DateStyle) string {
	layout := string(dateStyle)
	layout = strings.Replace(layout, "yyyy", "2006", 1)
	layout = strings.Replace(layout, "yy", "06", 1)
	layout = strings.Replace(layout, "MM", "01", 1)
	layout = strings.Replace(layout, "dd", "02", 1)
	layout = strings.Replace(layout, "HH", "15", 1)
	layout = strings.Replace(layout, "mm", "04", 1)
	layout = strings.Replace(layout, "ss", "05", 1)
	layout = strings.Replace(layout, "SSS", "000", -1)

	return date.Format(layout)
}

/* yyyyMMdd  */
func DateByStr(date string) Date {
	date = strings.TrimSpace(date)
	t, err := strconv.ParseInt(date, 10, 64)
	if err != nil {
		return Date(0)
	}

	return Date(t)
}

/* HH:mm:SS */
func TimeBySecStr(timeStr string) Time {
	t, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		return Time(0)
	}

	return Time(uint64(t.Hour())*10000000 + uint64(t.Minute())*100000 + uint64(t.Second())*1000)
}

/* HH:mm:SS.sss */
func TimeByMSStr(timeStr string) Time {
	t, err := time.Parse("15:04:05.000", timeStr)
	if err != nil {
		return Time(0)
	}

	return Time(uint64(t.Hour())*10000000 + uint64(t.Minute())*100000 + uint64(t.Second())*1000 + uint64(t.Nanosecond())/1000000)
}

/* yyyyMMdd-HH:mm:SS.sss */
func DateTimeByStr(timeStr string) DateTime {
	return DateTimeByFormatStr(GolangTimeFormat, timeStr)
}

func DateTimeSByStr(timeStr string) (DateTimeS, error) {
	s0 := DateTimeS{
		Year:   1,
		Month:  1,
		Day:    1,
		Hour:   0,
		Minute: 0,
		Second: 0,
	}
	t, err := time.Parse(GolangTimeFormat, timeStr)
	if err != nil {
		return s0, err
	}
	s0 = DateTimeS{
		Year:   t.Year(),
		Month:  int(t.Month()),
		Day:    t.Day(),
		Hour:   t.Hour(),
		Minute: t.Minute(),
		Second: t.Second(),
	}

	return s0, nil
}
func DateTimeByFormatStr(format string, timeStr string) DateTime {
	t, err := time.Parse(format, timeStr)
	if err != nil {
		return DateTime(0)
	}

	return DateTime(uint64(t.Year())*10000000000000 + uint64(t.Month())*100000000000 + uint64(t.Day())*1000000000 + uint64(t.Hour())*10000000 + uint64(t.Minute())*100000 + uint64(t.Second())*1000 + uint64(t.Nanosecond())/1000000)
}

func DateTimeByTime(t time.Time) DateTime {
	return DateTime(uint64(t.Year())*10000000000000 + uint64(t.Month())*100000000000 + uint64(t.Day())*1000000000 + uint64(t.Hour())*10000000 + uint64(t.Minute())*100000 + uint64(t.Second())*1000 + uint64(t.Nanosecond())/1000000)
}

func DateTimeByDateTimeS(dt DateTimeS) DateTime {
	return DateTime(uint64(dt.Year)*10000000000000 + uint64(dt.Month)*100000000000 + uint64(dt.Day)*1000000000 + uint64(dt.Hour)*10000000 + uint64(dt.Minute)*100000 + uint64(dt.Second)*1000)
}

func DateTimeByUTCTimestamp(timestamp int64) DateTime {
	tt := time.Unix(timestamp/1000000000, timestamp%1000000000)
	return DateTimeByTime(tt)
}
func (d Date) ToDateTime() DateTime {
	return DateTime(uint64(d) * 1000000000)
}

func (d Date) Value() uint64 {
	return uint64(d)
}

func (t Time) Value() uint64 {
	return uint64(t)
}

func (dt DateTime) Value() uint64 {
	return uint64(dt)
}

func (dt DateTime) Date() Date {
	return Date(uint64(dt) / 1000000000)
}

func (dt DateTime) Time() Time {
	return Time(uint64(dt) % 1000000000)
}

func (dt DateTime) RemoveMSec() DateTime {
	return DateTime(uint64(dt) / 1000 * 1000)
}

func (dt *DateTime) SetDate(d Date) {
	*dt = DateTime(uint64(dt.Time()) + uint64(d.ToDateTime()))
}

func (dt *DateTime) SetTime(t Time) {
	*dt = DateTime(uint64(t) + uint64(dt.Date().ToDateTime()))
}

func (dt DateTime) GoTime() time.Time {
	t, _ := time.Parse("20060102150405.000", fmt.Sprintf("%d.%03d", dt/1000, dt%1000))
	return t
}

func (t1 DateTime) Sub(t2 DateTime) time.Duration {
	tt1 := t1.GoTime()
	tt2 := t2.GoTime()

	return tt1.Sub(tt2)
}
func (t1 DateTime) SubDateTime(t2 DateTime) DateTime {
	td := t2 / 1000000000
	tt := t2 % 1000000000
	y, m, d := td/10000, (td%10000)/100, td%100
	ot := t1.GoTime().AddDate(int(-y), int(-m), int(-d))
	h, M, s, ms := tt/10000000, (tt%10000000)/100000, (tt%100000)/1000, tt%1000
	duration := time.Duration(h*3600+M*60+s)*1000000000 + time.Duration(ms)*1000000
	ot = ot.Add(-duration)
	y1, m1, d1 := ot.Date()
	h1 := ot.Hour()
	M1 := ot.Minute()
	s1 := ot.Second()
	ms1 := ot.Nanosecond() / 1000000
	return DateTime(y1)*10000000000000 +
		DateTime(m1)*100000000000 +
		DateTime(d1)*1000000000 +
		DateTime(h1)*10000000 +
		DateTime(M1)*100000 +
		DateTime(s1)*1000 +
		DateTime(ms1)
}

// GetBetweenDates 根据开始日期和结束日期计算出时间段内所有日期
// 参数为日期格式，如：2020-01-01
func GetBetweenDates(sdate, edate string) []string {
	d := []string{}
	timeFormatTpl := "2006-01-02 15:04:05"
	if len(timeFormatTpl) != len(sdate) {
		timeFormatTpl = timeFormatTpl[0:len(sdate)]
	}
	date, err := time.Parse(timeFormatTpl, sdate)
	if err != nil {
		// 时间解析，异常
		return d
	}
	date2, err := time.Parse(timeFormatTpl, edate)
	if err != nil {
		// 时间解析，异常
		return d
	}
	if date2.Before(date) {
		// 如果结束时间小于开始时间，异常
		return d
	}
	// 输出日期格式固定
	timeFormatTpl = "2006-01-02"
	date2Str := date2.Format(timeFormatTpl)
	d = append(d, date.Format(timeFormatTpl))
	for {
		date = date.AddDate(0, 0, 1)
		dateStr := date.Format(timeFormatTpl)
		d = append(d, dateStr)
		if dateStr == date2Str {
			break
		}
	}
	return d
}
