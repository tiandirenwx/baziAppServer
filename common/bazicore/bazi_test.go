package bazicore

import (
	"fmt"
	"path"
	"testing"
	"time"
)

func Test_calcBazi(t *testing.T) {

	//FamilyBirthday := [][]int{{1984, 1, 9}, {2017, 9, 2}, {2019, 10, 22}, {1987, 11, 17}, {1954, 4, 14},
	//{1953, 12, 19}, {1978, 12, 29}, {2003, 10, 15}, {2008, 4, 9}}

	inputArray1 := []string{"46990101-12:13:45.000", "21140212-12:13:45.000", "10500323-12:13:45.000",
		"01230414-12:13:45.000", "00010525-12:13:45.000", "00000616-12:13:45.000", "00010727-12:13:45.000",
		"01230808-12:13:45.000", "16780909-12:13:45.000", "20001010-12:13:45.000", "20121111-12:13:45.000",
		"22451231-12:13:45.000", JulianCalendar1582Spec, "06890218-12:13:45.000", "06901228-12:13:45.000",
		"00231231-12:13:45.000", "02400112-12:13:45.000", "07620429-12:13:45.000", "07620528-12:13:45.000",
		"07010114-12:13:45.000"}

	for _, v := range inputArray1 {
		dt := &BZDateTime{}
		dt = dt.New(v)
		bz := &BaziMeta{}
		bz = bz.New(dt, false, false, 120, 13, "")
		bz.CalcBazi()
		ret1 := bz.GetBazi()
		ret2 := bz.GetBaziBasic()
		nowStr := FormatDateTime(time.Now(), YYYYMMDD_HH_MM_SS_SSS_EN)
		now := &BZDateTime{}
		now = now.New(nowStr)
		ret3 := bz.GetAge(now)
		ret4 := bz.GetLifa()
		ret5 := bz.GetAst()
		ret6 := bz.GetSolarTerms()
		fmt.Printf("input: %s,age: %s, lifa: %s, ast: %s, solar: %s,  bazi info: %s \n,print bazi: %s \n, ",
			v, ret3, ret4, ret5, ret6, ret1, ret2)

	}
}

func Test_PrintBazi(t *testing.T) {
	FamilyBirthday := [][]int{{1984, 1, 9}, {2017, 9, 2}, {2019, 10, 22}, {1987, 11, 17}, {1954, 4, 14},
		{1953, 12, 19}, {1978, 12, 29}, {2003, 10, 15}, {2008, 4, 9}}
	shiChen := []string{"07:37", "11:40", "21:17", "3:30", "3:30", "21:30", "11:50", "14:20", "8:10"}
	genderList := []bool{false, true, false, true, true, false, false, true, false}
	nameList := []string{"xxx1", "xxx2", "xxx3", "xxx3", "xxx4", "xxx5", "xxx5", "xxx5", "xxx6"}
	for i, v := range FamilyBirthday {
		year, month, day, err := Lunar2solar(v[0], v[1], v[2], false, false)
		if err != nil {
			fmt.Printf("error occur: %s", err)
			continue
		}
		sc := shiChen[i] + ":00.000"
		birthTimeStr := fmt.Sprintf("%04d%02d%02d-%s", year, month, day, sc)
		dt := &BZDateTime{}
		dt = dt.New(birthTimeStr)
		bz := &BaziMeta{}
		bz = bz.New(dt, genderList[i], false, 120, 13, nameList[i])
		bz.CalcBazi()
		ret1 := bz.GetBazi()
		ret2 := bz.GetBaziBasic()
		nowStr := FormatDateTime(time.Now(), YYYYMMDD_HH_MM_SS_SSS_EN)
		now := &BZDateTime{}
		now = now.New(nowStr)
		ret3 := bz.GetAge(now)
		ret4 := bz.GetLifa()
		ret5 := bz.GetAst()
		ret6 := bz.GetSolarTerms()
		ret7 := bz.GetQiYun()
		ret8 := bz.GetJiaoYun()
		ret9 := bz.GetStartYearList()
		ret10 := bz.GetFleetingYearList()
		ret11 := bz.GetEndYearList()
		ret12 := bz.GetSiZhu()
		ret13 := bz.GetLunarBirth()
		ret14 := bz.GetSolarBirth()
		fmt.Printf("input: %s,age: %s, lifa: %s, "+
			"ast: %s, solar: %s,  bazi info: %s \n,%s \n,%s\n,%s\n,%+v \n, %+v\n, %+v\n ,%+v \n, %s\n,%s \n",
			birthTimeStr, ret3, ret4, ret5, ret6, ret1, ret2, ret7, ret8, ret9, ret10, ret11, ret12, ret13, ret14)
	}

}

func Test_RenderImage(t *testing.T) {
	FamilyBirthday := [][]int{{1984, 1, 9}, {2017, 9, 2}, {2019, 10, 22}, {1987, 11, 17}, {1954, 4, 14},
		{1953, 12, 19}, {1978, 12, 29}, {2003, 10, 15}, {2008, 4, 9}}
	shiChen := []string{"07:37", "11:40", "21:17", "3:30", "3:30", "21:30", "11:50", "14:20", "8:10"}
	genderList := []bool{false, true, false, true, true, false, false, true, false}
	nameList := []string{"xxx1", "xxx2", "xxx3", "xxx4", "xxx4", "xxx5", "xxx6", "xxx7", "xxx8"}
	for i, v := range FamilyBirthday {
		year, month, day, err := Lunar2solar(v[0], v[1], v[2], false, false)
		if err != nil {
			fmt.Printf("error occur: %s", err)
			continue
		}
		sc := shiChen[i] + ":00.000"
		birthTimeStr := fmt.Sprintf("%04d%02d%02d-%s", year, month, day, sc)
		dt := &BZDateTime{}
		dt = dt.New(birthTimeStr)
		bz := &BaziMeta{}
		bz = bz.New(dt, genderList[i], false, 120, 13, nameList[i])
		bz.CalcBazi()
		nowStr := FormatDateTime(time.Now(), YYYYMMDD_HH_MM_SS_SSS_EN)
		now := &BZDateTime{}
		now = now.New(nowStr)
		configPath := "/Users/fumingzhang/Documents/Code/goCode/src/fuming.com/baziAppServer/config"
		img, err := bz.RenderImage(configPath)
		if err != nil {
			fmt.Printf("error occur: %s", err.Error())
			continue
		}

		fileName := path.Join(configPath, "result", bz.Name+".jpeg")
		err = saveImageAsJpeg(fileName, img)
		if err != nil {
			fmt.Printf("error occur: %s", err.Error())
			continue
		}

		fileName = path.Join(configPath, "result", bz.Name+".png")
		err = saveImageAsPng(fileName, img)
		if err != nil {
			fmt.Printf("error occur: %s", err.Error())
			continue
		}

		fileName = path.Join(configPath, "result", bz.Name+".gif")
		err = saveImageAsGif(fileName, img)
		if err != nil {
			fmt.Printf("error occur: %s", err.Error())
			continue
		}
	}
}
