package bazicore

import "math"

const (

	//--#========角度变换===============
	SecondsPerRadian float64 = 180 * 3600 / math.Pi //seconds per radian,每弧度的角秒数
	AnglesPerRadian  float64 = 180 / math.Pi        //angles per radian,每弧度的角度数

	JulianCalendarJ2000    = 2451545 //2000年前儒略日数(2000-1-1 12:00:00格林威治平时)
	JulianCalendar1582Spec = "15821015-00:00:00.000"
	BzDataTimeInitTime     = "00010101-00:00:00.000" //初始化BZDateTime的特殊时间字符串
)
