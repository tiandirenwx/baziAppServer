package bazicore

import (
	"fmt"
	"testing"
)

func Test_gcal2jd(t *testing.T) {
	var year int64 = 2000
	var month int64 = 1
	var day int64 = 1
	res1, res2 := gcal2jd(year, month, day)
	fmt.Printf("gcal2jd input %d,%d,%d, output, res1 = %f, res2 = %f \n", year, month, day, res1, res2)
	fmt.Printf("#-------------------------------------------------------------------------------------------#\n")
	//year := []int64{-4699, -2114, -1050, -123, -1, 0, 1, 123, 1678.0, 2000, 2012, 2245}
	//month := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	//day := []int64{1, 12, 23, 14, 25, 16, 27, 8, 9, 10, 11, 31}
	inputArray1 := [][]int64{{-4699, 1, 1}, {-2114, 2, 12}, {-1050, 3, 23}, {-123, 4, 14}, {-1, 5, 25},
		{0, 6, 16}, {1, 7, 27}, {123, 8, 8}, {1678.0, 9, 9}, {2000, 10, 10}, {2012, 11, 11}, {2245, 12, 31}}
	for i, v := range inputArray1 {
		fmt.Printf("the input array1 is %+v,index is %d \n", v, i)
		res1, res2 := gcal2jd(v[0], v[1], v[2])
		fmt.Printf("gcal2jd input %d,%d,%d, output, res1 = %f, res2 = %f \n", v[0], v[1], v[2], res1, res2)
	}

	fmt.Printf("#-------------------------------------------------------------------------------------------#\n")
	inputArray2 := [][]int64{{2000, -2, -4}, {1999, 9, 26}, {2000, 2, -1}, {2000, 1, 30}, {2000, 3, -1}, {2000, 2, 28}, {2000, 0, 1},
		{1999, 12, 1}, {2000, 3, 0}, {2000, 2, 29}, {2000, 2, 30}, {2000, 3, 1}, {2001, 2, 30}, {2001, 3, 2}}

	for i, v := range inputArray2 {
		fmt.Printf("the input array2 is %+v,index is %d \n", v, i)
		res1, res2 := gcal2jd(v[0], v[1], v[2])
		fmt.Printf("gcal2jd input %d,%d,%d, output, res1 = %f, res2 = %f \n", v[0], v[1], v[2], res1, res2)
	}

}

func Test_jd2gcal(t *testing.T) {
	var year int64 = 2000
	var month int64 = 1
	var day int64 = 1
	r1, r2 := gcal2jd(year, month, day)
	res1, res2, res3, res4 := jd2gcal(r1, r2)
	fmt.Printf("gcal2jd input: (%d,%d,%d), output: (%f, %f),jd2gcal output: (%d,  %d, %d, %f) \n", year, month, day, r1, r2, res1, res2, res3, res4)
	fmt.Printf("#-------------------------------------------------------------------------------------------#\n")
	//year := []int64{-4699, -2114, -1050, -123, -1, 0, 1, 123, 1678.0, 2000, 2012, 2245}
	//month := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	//day := []int64{1, 12, 23, 14, 25, 16, 27, 8, 9, 10, 11, 31}
	inputArray1 := [][]int64{{-4699, 1, 1}, {-2114, 2, 12}, {-1050, 3, 23}, {-123, 4, 14}, {-1, 5, 25},
		{0, 6, 16}, {1, 7, 27}, {123, 8, 8}, {1678.0, 9, 9}, {2000, 10, 10}, {2012, 11, 11}, {2245, 12, 31}}
	for _, v := range inputArray1 {
		r1, r2 = gcal2jd(v[0], v[1], v[2])
		res1, res2, res3, res4 = jd2gcal(r1, r2)
		fmt.Printf("array1 for gcal2jd input: (%d,%d,%d), output: (%f, %f),jd2gcal output: (%d,  %d, %d, %f) \n", v[0], v[1], v[2], r1, r2, res1, res2, res3, res4)
	}

	fmt.Printf("#-------------------------------------------------------------------------------------------#\n")
	inputArray2 := [][]int64{{2000, -2, -4}, {1999, 9, 26}, {2000, 2, -1}, {2000, 1, 30}, {2000, 3, -1}, {2000, 2, 28}, {2000, 0, 1},
		{1999, 12, 1}, {2000, 3, 0}, {2000, 2, 29}, {2000, 2, 30}, {2000, 3, 1}, {2001, 2, 30}, {2001, 3, 2}}

	for _, v := range inputArray2 {
		r1, r2 = gcal2jd(v[0], v[1], v[2])
		res1, res2, res3, res4 = jd2gcal(r1, r2)
		fmt.Printf("array2 for gcal2jd input: (%d,%d,%d), output: (%f, %f),jd2gcal output: (%d,  %d, %d, %f) \n", v[0], v[1], v[2], r1, r2, res1, res2, res3, res4)
	}

	inputArray3 := [][]int64{{1950, 1, 1}, {1999, 10, 12}, {2000, 2, 30}, {-1999, 10, 12}, {2000, -2, -4}}

	for _, v := range inputArray3 {
		r1, r2 = gcal2jd(v[0], v[1], v[2])
		res1, res2, res3, res4 = jd2gcal(r1, r2)
		fmt.Printf("array3 for gcal2jd input: (%d,%d,%d), output: (%f, %f),jd2gcal output: (%d,  %d, %d, %f) \n", v[0], v[1], v[2], r1, r2, res1, res2, res3, res4)
	}

}
func Test_jcal2jd(t *testing.T) {
	var year int64 = 2000
	var month int64 = 1
	var day int64 = 1
	res1, res2 := jcal2jd(year, month, day)
	fmt.Printf("jcal2jd input %d,%d,%d, output: (%f, %f) \n", year, month, day, res1, res2)
	fmt.Printf("#-------------------------------------------------------------------------------------------#\n")
	//year := []int64{-4699, -2114, -1050, -123, -1, 0, 1, 123, 1678.0, 2000, 2012, 2245}
	//month := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	//day := []int64{1, 12, 23, 14, 25, 16, 27, 8, 9, 10, 11, 31}
	inputArray1 := [][]int64{{-4699, 1, 1}, {-2114, 2, 12}, {-1050, 3, 23}, {-123, 4, 14}, {-1, 5, 25},
		{0, 6, 16}, {1, 7, 27}, {123, 8, 8}, {1678.0, 9, 9}, {2000, 10, 10}, {2012, 11, 11}, {2245, 12, 31}}
	/*
		(2400000.5, -2395252.0)
		(2400000.5, -1451039.0)
		(2400000.5, -1062374.0)
		(2400000.5, -723765.0)
		(2400000.5, -679164.0)
		(2400000.5, -678776.0)
		(2400000.5, -678370.0)
		(2400000.5, -633798.0)
		(2400000.5, -65802.0)
		(2400000.5, 51840.0)
		(2400000.5, 56255.0)
		(2400000.5, 141408.0)
	*/
	for i, v := range inputArray1 {
		fmt.Printf("the input array1 is %+v,index is %d \n", v, i)
		res1, res2 := jcal2jd(v[0], v[1], v[2])
		fmt.Printf("array 1 jcal2jd input (%d,%d,%d), output: (%f, %f) \n", v[0], v[1], v[2], res1, res2)
	}
}

func Test_jd2jcal(t *testing.T) {
	var year int64 = 2000
	var month int64 = 1
	var day int64 = 1
	r1, r2 := jcal2jd(year, month, day)
	res1, res2, res3, res4 := jd2jcal(r1, r2)
	fmt.Printf("jcal2jd input: (%d,%d,%d), output: (%f, %f),jd2jcal output: (%d, %d, %d, %f) \n", year, month, day, r1, r2, res1, res2, res3, res4)
	fmt.Printf("#-------------------------------------------------------------------------------------------#\n")
	//year := []int64{-4699, -2114, -1050, -123, -1, 0, 1, 123, 1678.0, 2000, 2012, 2245}
	//month := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	//day := []int64{1, 12, 23, 14, 25, 16, 27, 8, 9, 10, 11, 31}
	inputArray1 := [][]int64{{-4699, 1, 1}, {-2114, 2, 12}, {-1050, 3, 23}, {-123, 4, 14}, {-1, 5, 25},
		{0, 6, 16}, {1, 7, 27}, {123, 8, 8}, {1678.0, 9, 9}, {2000, 10, 10}, {2012, 11, 11}, {2245, 12, 31}}
	for _, v := range inputArray1 {
		r1, r2 = jcal2jd(v[0], v[1], v[2])
		res1, res2, res3, res4 = jd2jcal(r1, r2)
		fmt.Printf("array1 for jcal2jd input: (%d,%d,%d), output: (%f, %f),jd2jcal output: (%d,  %d, %d, %f) \n", v[0], v[1], v[2], r1, r2, res1, res2, res3, res4)
	}

	fmt.Printf("#-------------------------------------------------------------------------------------------#\n")
	inputArray2 := [][]int64{{2000, -2, -4}, {1999, 9, 26}, {2000, 2, -1}, {2000, 1, 30}, {2000, 3, -1}, {2000, 2, 28}, {2000, 0, 1},
		{1999, 12, 1}, {2000, 3, 0}, {2000, 2, 29}, {2000, 2, 30}, {2000, 3, 1}, {2001, 2, 30}, {2001, 3, 2}}

	for _, v := range inputArray2 {
		r1, r2 = jcal2jd(v[0], v[1], v[2])
		res1, res2, res3, res4 = jd2jcal(r1, r2)
		fmt.Printf("array2 for jcal2jd input: (%d,%d,%d), output: (%f, %f),jd2jcal output: (%d,  %d, %d, %f) \n", v[0], v[1], v[2], r1, r2, res1, res2, res3, res4)
	}

	inputArray3 := [][]int64{{1950, 1, 1}, {1999, 10, 12}, {2000, 2, 30}, {-1999, 10, 12}, {2000, -2, -4}}
	for _, v := range inputArray3 {
		r1, r2 = jcal2jd(v[0], v[1], v[2])
		res1, res2, res3, res4 = jd2jcal(r1, r2)
		fmt.Printf("array3 for jcal2jd input: (%d,%d,%d), output: (%f, %f),jd2jcal output: (%d,  %d, %d, %f) \n", v[0], v[1], v[2], r1, r2, res1, res2, res3, res4)
	}

}
