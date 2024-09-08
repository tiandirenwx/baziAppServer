package bazicore

import (
	"fmt"
	"testing"
)

func Test_setFromJD(t *testing.T) {
	inputArray := []float64{1972764.0, 1729823, 1808729, 1999497, 1999526, 1977112}
	for _, jd := range inputArray {
		dt := &BZDateTime{}
		dt = dt.New("00010101-00:00:00.000")
		dt.SetFromJD(jd)
		fmt.Printf("input %f, output: %+v\n", jd, dt)
	}

}

func Test_getLunarInfo(t *testing.T) {
	inputArray1 := []string{"46990101-12:13:45.000", "21140212-12:13:45.000", "10500323-12:13:45.000",
		"01230414-12:13:45.000", "00010525-12:13:45.000", "00000616-12:13:45.000", "00010727-12:13:45.000",
		"01230808-12:13:45.000", "16780909-12:13:45.000", "20001010-12:13:45.000", "20121111-12:13:45.000",
		"22451231-12:13:45.000", JulianCalendar1582Spec, "06890218-12:13:45.000", "06901228-12:13:45.000",
		"00231231-12:13:45.000", "02400112-12:13:45.000", "07620429-12:13:45.000", "07620528-12:13:45.000",
		"07010114-12:13:45.000"}
	for _, v := range inputArray1 {
		dt := &BZDateTime{}
		dt = dt.New(v)
		s1, s2, s3 := dt.GetLunarInfo()
		fmt.Printf("input: %s, output %s,%s,%s\n", v, s1, s2, s3)
	}

}
