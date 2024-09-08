package bazicore

import (
	"math"
)

// translate from https://github.com/phn/jdcal/blob/master/jdcal.py
var MJD_0 = 2400000.5
var MJD_JD2000 = 51544.5

func IsLeap(year int64) bool {
	//"""Leap year or not in the Gregorian calendar."""
	//# Divisible by 4 and,
	//# either not divisible by 100 or divisible by 400.
	leap := year%400 == 0 || (year%4 == 0 && year%100 != 0)
	return leap
}

func gcal2jd(year, month, day int64) (float64, float64) {
	/*
	   """Gregorian calendar date to Julian date.
	   The input and output are for the proleptic Gregorian calendar,
	   i.e., no consideration of historical usage of the calendar is
	   made.
	   Parameters
	   ----------
	   year : int
	   Year as an integer.
	   month : int
	   Month as an integer.
	   day : int
	   Day as an integer.
	   Returns
	   -------
	   jd1, jd2: 2-element tuple of floats
	   When added together, the numbers give the Julian date for the
	   given Gregorian calendar date. The first number is always
	   MJD_0 i.e., 2400000.5. So the second is the MJD.
	   Examples
	   --------
	   >>> gcal2jd(2000,1,1)
	   (2400000.5, 51544.0)
	   >>> 2400000.5 + 51544.0 + 0.5
	   2451545.0
	   >>> year = [-4699, -2114, -1050, -123, -1, 0, 1, 123, 1678.0, 2000,
	   ....: 2012, 2245]
	   >>> month = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12]
	   >>> day = [1, 12, 23, 14, 25, 16, 27, 8, 9, 10, 11, 31]
	   >>> x = [gcal2jd(y, m, d) for y, m, d in zip(year, month, day)]
	   >>> for i in x: print i
	   (2400000.5, -2395215.0)
	   (2400000.5, -1451021.0)
	   (2400000.5, -1062364.0)
	   (2400000.5, -723762.0)
	   (2400000.5, -679162.0)
	   (2400000.5, -678774.0)
	   (2400000.5, -678368.0)
	   (2400000.5, -633797.0)
	   (2400000.5, -65812.0)
	   (2400000.5, 51827.0)
	   (2400000.5, 56242.0)
	   (2400000.5, 141393.0)
	   Negative months and days are valid. For example, 2000/-2/-4 =>
	   1999/+12-2/-4 => 1999/10/-4 => 1999/9/30-4 => 1999/9/26.
	   >>> gcal2jd(2000, -2, -4)
	   (2400000.5, 51447.0)
	   >>> gcal2jd(1999, 9, 26)
	   (2400000.5, 51447.0)
	   >>> gcal2jd(2000, 2, -1)
	   (2400000.5, 51573.0)
	   >>> gcal2jd(2000, 1, 30)
	   (2400000.5, 51573.0)
	   >>> gcal2jd(2000, 3, -1)
	   (2400000.5, 51602.0)
	   >>> gcal2jd(2000, 2, 28)
	   (2400000.5, 51602.0)
	   Month 0 becomes previous month.
	   >>> gcal2jd(2000, 0, 1)
	   (2400000.5, 51513.0)
	   >>> gcal2jd(1999, 12, 1)
	   (2400000.5, 51513.0)
	   Day number 0 becomes last day of previous month.
	   >>> gcal2jd(2000, 3, 0)
	   (2400000.5, 51603.0)
	   >>> gcal2jd(2000, 2, 29)
	   (2400000.5, 51603.0)
	   If `day` is greater than the number of days in `month`, then it
	   gets carried over to the next month.
	   >>> gcal2jd(2000,2,30)
	   (2400000.5, 51604.0)
	   >>> gcal2jd(2000,3,1)
	   (2400000.5, 51604.0)
	   >>> gcal2jd(2001,2,30)
	   (2400000.5, 51970.0)
	   >>> gcal2jd(2001,3,2)
	   (2400000.5, 51970.0)
	   Notes
	   -----
	   The returned Julian date is for mid-night of the given date. To
	   find the Julian date for any time of the day, simply add time as a
	   fraction of a day. For example Julian date for mid-day can be
	   obtained by adding 0.5 to either the first part or the second
	   part. The latter is preferable, since it will give the MJD for the
	   date and time.
	   BC dates should be given as -(BC - 1) where BC is the year. For
	   example 1 BC == 0, 2 BC == -1, and so on.
	   Negative numbers can be used for `month` and `day`. For example
	   2000, -1, 1 is the same as 1999, 11, 1.
	   The Julian dates are proleptic Julian dates, i.e., values are
	   returned without considering if Gregorian dates are valid for the
	   given date.
	   The input values are truncated to integers.
	   """
	*/

	a := int64(float64(month-14) / 12.0)
	jd := int64(float64(1461*(year+4800+a)) / 4.0)
	jd += int64(float64(367*float64(month-2-12*a)) / 12.0)
	x := int64(float64(year+4900+a) / 100.0)
	jd -= int64(float64(3*x) / 4.0)
	jdRes := float64(jd)
	jdRes += float64(day) - 2432075.5 //# was 32075; add 2400000.5

	jdRes -= 0.5 //# 0 hours; above JD is for midday, switch to midnight.

	return MJD_0, jdRes
}

func jd2gcal(jd1, jd2 float64) (int64, int64, int64, float64) {
	/*
	    """Julian date to Gregorian calendar date and time of day.
	   The input and output are for the proleptic Gregorian calendar,
	   i.e., no consideration of historical usage of the calendar is
	   made.
	   Parameters
	   ----------
	   jd1, jd2: float
	   Sum of the two numbers is taken as the given Julian date. For
	   example `jd1` can be the zero point of MJD (MJD_0) and `jd2`
	   can be the MJD of the date and time. But any combination will
	   work.
	   Returns
	   -------
	   y, m, d, f : int, int, int, float
	   Four element tuple containing year, month, day and the
	   fractional part of the day in the Gregorian calendar. The first
	   three are integers, and the last part is a float.
	   Examples
	   --------
	   >>> jd2gcal(*gcal2jd(2000,1,1))
	   (2000, 1, 1, 0.0)
	   >>> jd2gcal(*gcal2jd(1950,1,1))
	   (1950, 1, 1, 0.0)
	   Out of range months and days are carried over to the next/previous
	   year or next/previous month. See gcal2jd for more examples.
	   >>> jd2gcal(*gcal2jd(1999,10,12))
	   (1999, 10, 12, 0.0)
	   >>> jd2gcal(*gcal2jd(2000,2,30))
	   (2000, 3, 1, 0.0)
	   >>> jd2gcal(*gcal2jd(-1999,10,12))
	   (-1999, 10, 12, 0.0)
	   >>> jd2gcal(*gcal2jd(2000, -2, -4))
	   (1999, 9, 26, 0.0)
	   >>> gcal2jd(2000,1,1)
	   (2400000.5, 51544.0)
	   >>> jd2gcal(2400000.5, 51544.0)
	   (2000, 1, 1, 0.0)
	   >>> jd2gcal(2400000.5, 51544.5)
	   (2000, 1, 1, 0.5)
	   >>> jd2gcal(2400000.5, 51544.245)
	   (2000, 1, 1, 0.24500000000261934)
	   >>> jd2gcal(2400000.5, 51544.1)
	   (2000, 1, 1, 0.099999999998544808)
	   >>> jd2gcal(2400000.5, 51544.75)
	   (2000, 1, 1, 0.75)
	   Notes
	   -----
	   The last element of the tuple is the same as
	   (hh + mm / 60.0 + ss / 3600.0) / 24.0
	   where hh, mm, and ss are the hour, minute and second of the day.
	   See Also
	   --------
	   gcal2jd
	   """
	*/
	jd1I, jd1F := math.Modf(jd1)
	jd2I, jd2F := math.Modf(jd2)

	jdI := jd1I + jd2I

	f := jd1F + jd2F

	//# Set JD to noon of the current date. Fractional part is the
	//# fraction from midnight of the current date.
	if -0.5 < f && f < 0.5 {
		f += 0.5
	} else if f >= 0.5 {
		jdI += 1
		f -= 0.5
	} else if f <= -0.5 {
		jdI -= 1
		f += 1.5
	}

	ell := jdI + 68569
	n := int64((4 * ell) / 146097.0)
	ell -= float64(int64(float64((146097*n)+3) / 4.0))
	i := int64((4000 * (ell + 1)) / 1461001)
	ell -= float64(int64(float64(1461*i)/4.0) - 31)
	j := int64((80 * ell) / 2447.0)
	day := ell - float64(int64(float64(2447*j)/80.0))
	ell = float64(int64(float64(j) / 11.0))
	month := float64(j) + 2 - (12 * ell)
	year := float64(100*(n-49)+i) + ell
	return int64(year), int64(month), int64(day), f
}

func jcal2jd(year, month, day int64) (float64, float64) {
	/*
		"""Julian calendar date to Julian date.
		The input and output are for the proleptic Julian calendar,
		i.e., no consideration of historical usage of the calendar is
		made.
		Parameters
		----------
		year : int
		Year as an integer.
		month : int
		Month as an integer.
		day : int
		Day as an integer.
		Returns
		-------
		jd1, jd2: 2-element tuple of floats
		When added together, the numbers give the Julian date for the
		given Julian calendar date. The first number is always
		MJD_0 i.e., 2451545.5. So the second is the MJD.
		Examples
		--------
		>>> jcal2jd(2000, 1, 1)
		(2400000.5, 51557.0)
		>>> year = [-4699, -2114, -1050, -123, -1, 0, 1, 123, 1678, 2000,
		...:  2012, 2245]
		>>> month = [1, 2, 3, 4, 5, 6, 7, 8, 10, 11, 12]
		>>> day = [1, 12, 23, 14, 25, 16, 27, 8, 9, 10, 11, 31]
		>>> x = [jcal2jd(y, m, d) for y, m, d in zip(year, month, day)]
		>>> for i in x: print i
		(2400000.5, -2395252.0)
		(2400000.5, -1451039.0)
		(2400000.5, -1062374.0)
		(2400000.5, -723765.0)
		(2400000.5, -679164.0)
		(2400000.5, -678776.0)
		(2400000.5, -678370.0)
		(2400000.5, -633798.0)
		(2400000.5, -65772.0)
		(2400000.5, 51871.0)
		(2400000.5, 56285.0)
		Notes
		-----
		Unlike `gcal2jd`, negative months and days can result in incorrect
		Julian dates.
		"""
	*/
	jd := float64(367 * year) //先保证jd是一个float
	x := int64(float64(month-9) / 7.0)
	jd -= float64(int64((7 * float64(year+5001+x)) / 4.0))
	jd += float64(int64(float64(275*month) / 9.0))
	jd += float64(day)
	jd += 1729777 - 2400000.5 //# Return 240000.5 as first part of JD.

	jd -= 0.5 //# Convert midday to midnight.

	return MJD_0, jd
}

func jd2jcal(jd1, jd2 float64) (int64, int64, int64, float64) {
	/*
		"""Julian calendar date for the given Julian date.
		The input and output are for the proleptic Julian calendar,
		i.e., no consideration of historical usage of the calendar is
		made.
		Parameters
		----------
		jd1, jd2: float
		Sum of the two numbers is taken as the given Julian date. For
		example `jd1` can be the zero point of MJD (MJD_0) and `jd2`
		can be the MJD of the date and time. But any combination will
		work.
		Returns
		-------
		y, m, d, f : int, int, int, float
		Four element tuple containing year, month, day and the
		fractional part of the day in the Julian calendar. The first
		three are integers, and the last part is a float.
		Examples
		--------
		>>> jd2jcal(*jcal2jd(2000, 1, 1))
		(2000, 1, 1, 0.0)
		>>> jd2jcal(*jcal2jd(-4000, 10, 11))
		(-4000, 10, 11, 0.0)
		>>> jcal2jd(2000, 1, 1)
		(2400000.5, 51557.0)
		>>> jd2jcal(2400000.5, 51557.0)
		(2000, 1, 1, 0.0)
		>>> jd2jcal(2400000.5, 51557.5)
		(2000, 1, 1, 0.5)
		>>> jd2jcal(2400000.5, 51557.245)
		(2000, 1, 1, 0.24500000000261934)
		>>> jd2jcal(2400000.5, 51557.1)
		(2000, 1, 1, 0.099999999998544808)
		>>> jd2jcal(2400000.5, 51557.75)
		(2000, 1, 1, 0.75)
		"""
	*/
	jd1I, jd1F := math.Modf(jd1)
	jd2I, jd2F := math.Modf(jd2)

	jdI := jd1I + jd2I

	f := jd1F + jd2F

	//# Set JD to noon of the current date. Fractional part is the
	//# fraction from midnight of the current date.
	if -0.5 < f && f < 0.5 {
		f += 0.5
	} else if f >= 0.5 {
		jdI += 1
		f -= 0.5
	} else if f <= -0.5 {
		jdI -= 1
		f += 1.5
	}

	j := jdI + 1402.0
	k := int64((j - 1) / 1461.0)
	ell := j - (1461.0 * float64(k))
	n := int64((ell-1)/365.0) - int64(ell/1461.0)
	i := ell - (365.0 * float64(n)) + 30.0
	j = float64(int64((80.0 * i) / 2447.0))
	day := i - float64(int64((2447.0*j)/80.0))
	i = float64(int64(j / 11.0))
	month := j + 2 - (12.0 * i)
	year := float64((4*k)+n) + i - 4716.0

	return int64(year), int64(month), int64(day), f
}
