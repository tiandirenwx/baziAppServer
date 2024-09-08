package bazicore

import (
	"fmt"
	"testing"
)

func getAllDaysFromYear(y int) [][]int {
	days := make([][]int, 0)
	year := int64(y)
	for month := 1; month <= 12; month++ {
		for day := 1; day <= 31; day++ {
			//如果是2月
			if month == 2 {
				if IsLeap(year) && day == 30 {
					//闰年2月29天
					break
				} else if !IsLeap(year) && day == 29 {
					//平年2月28天
					break
				} else {
					dayTemp := []int{int(year), month, day}
					days = append(days, dayTemp)
				}
			} else if month == 4 || month == 6 || month == 9 || month == 11 {
				//小月踢出来
				if day == 31 {
					break
				}
				dayTemp := []int{int(year), month, day}
				days = append(days, dayTemp)
			} else {
				dayTemp := []int{int(year), month, day}
				days = append(days, dayTemp)
			}
		}
	}
	return days
}

func Test_solar2lunar(t *testing.T) {

	ret, err := Solar2lunar(2017, 10, 21, 12)
	if err != nil {
		return
	}
	fmt.Printf("公元2017年10月21日对应农历: %+v", ret)
	fmt.Printf("#------------------------------------------------------------------------------------------#\n")

	ret, err = Solar2lunar(689, 2, 14, 12)
	if err != nil {
		return
	}
	fmt.Printf("公元689年2月14日对应农历: %+v", ret)
	fmt.Printf("#-------------------------------------------------------------------------------------------#\n")

	ret, err = Solar2lunar(690, 2, 14, 12)
	if err != nil {
		return
	}
	fmt.Printf("公元690年2月14日对应农历: %+v", ret)
	fmt.Printf("#-------------------------------------------------------------------------------------------#\n")

	inputArray1 := [][]int{{4699, 1, 1, 12}, {2114, 2, 12, 12}, {1050, 3, 23, 12}, {123, 4, 14, 12}, {1, 5, 25, 12},
		{0, 6, 16, 12}, {1, 7, 27, 12}, {123, 8, 8, 12}, {1678, 9, 9, 12}, {2000, 10, 10, 12}, {2012, 11, 11, 12}, {2245, 12, 31, 12}}

	for i, v := range inputArray1 {
		fmt.Printf("the input array1 is %+v,index is %d \n", v, i)
		ret, err = Solar2lunar(v[0], v[1], v[2], v[3])
		if err != nil {
			fmt.Printf("error occur: %s", err)
			continue
		}
		fmt.Printf("input %d,%d,%d,%d, output, res = %+v \n", v[0], v[1], v[2], v[3], ret)
	}

	ret, err = Solar2lunar(1, 1, 1, 12)
	if err != nil {
		return
	}
	fmt.Printf("公元-689年2月14日对应农历: %+v", ret)

	for i := 1; i < 9999; i++ {
		days := getAllDaysFromYear(i)
		for _, v := range days {
			ret, err := Solar2lunar(v[0], v[1], v[2], 12)
			if err != nil {
				fmt.Printf("error occur: %s", err)
				continue
			}
			fmt.Printf("input %d,%d,%d, output, res = %+v \n", v[0], v[1], v[2], ret)
		}
	}

}

func Test_lunar2solar(t *testing.T) {
	FamilyBirthday := [][]int{{1984, 1, 9}, {2017, 9, 2}, {2019, 10, 22}, {1987, 11, 17}, {1954, 4, 14},
		{1953, 12, 19}, {1978, 12, 29}, {2003, 10, 15}, {2008, 4, 9}}

	for _, v := range FamilyBirthday {
		year, month, day, err := Lunar2solar(v[0], v[1], v[2], false, false)
		if err != nil {
			fmt.Printf("error occur: %s", err)
			continue
		}
		fmt.Printf("input %d,%d,%d, output, lunar year = %d, month = %d, day = %d \n", v[0], v[1], v[2], year, month, day)
	}
	// 公元700的特殊情况，['冬', '腊', '正', '二', '三', '四', '五', '六', '七', '八', '九', '十', '冬', '腊']
	//year, month, day, err := lunar2solar(v[0], v[1], v[2], false, false)
}
