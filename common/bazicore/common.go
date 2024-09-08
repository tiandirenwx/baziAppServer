package bazicore

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"reflect"
	"strconv"
)

func rad2mrad(v float64) float64 {
	value := v - v/(2*math.Pi)*(2*math.Pi) //v % (2 * math.Pi), a % b = a - a / b * b
	if value < 0 {
		return value + 2*math.Pi
	}
	return value
}

func rad2rrad(v float64) float64 {
	value := v - v/(2*math.Pi)*(2*math.Pi) //v % (2 * math.Pi), a % b = a - a / b * b
	if (value - 0.0) <= -math.Pi {
		return value + 2*math.Pi
	}

	if (value - 0.0) > math.Pi {
		return value - 2*math.Pi
	}
	return value
}

// 取整数部分,向零取整
func int2(x float64) int64 {
	value := int64(math.Floor(x))
	return value
}

// 将弧度转为字串
// # tim=0输出格式示例: -23°59' 48.23"
// # tim=1输出格式示例:  18h 29m 44.52s
func rad2str(d float64, tim bool) string {
	var s, w1, w2, w3 = "+", "°", "’", "”"
	if d < 0 {
		d = -d
		s = "-"
	}

	if tim {
		d = d * 12 / math.Pi
		w1, w2, w3 = "h ", "m ", "s "
	} else {
		d = d * 180 / math.Pi
	}

	var a, b, c float64
	a = math.Floor(d)
	d = (d - a) * 60
	b = math.Floor(d)
	d = (d - b) * 60
	c = math.Floor(d)
	d = (d - c) * 100
	d = math.Floor(d + 0.5)

	if d >= 100 {
		d = d - 100
		c = c + 1
	}

	if c >= 60 {
		c = c - 60
		b = b + 1
	}

	if b >= 60 {
		b = b - 60
		a = a + 1
	}

	stra := "   " + strconv.FormatInt(int64(a), 10)
	strb := "0" + strconv.FormatInt(int64(b), 10)
	strc := "0" + strconv.FormatInt(int64(c), 10)
	strd := "0" + strconv.FormatInt(int64(d), 10)

	s += stra[len(stra)-3:] + w1
	s += strb[len(strb)-2:] + w2
	s += strc[len(strc)-2:] + "."
	s += strd[len(strd)-2:] + w3

	return s
}

func Index[T comparable](v T, array []T) (int, bool) {
	for i, item := range array {
		if v == item {
			return i, true
		}
	}

	return -1, false
}

func IndexPlus[T comparable](v T, array []T, start, end int) (int, bool) {
	arrayLen := len(array)
	if -1 == end {
		end = arrayLen
	}

	if end > arrayLen || end < start {
		return -1, false
	}

	for i := start; i < end; i++ {
		if v == array[i] {
			return i, true
		}
	}

	return -1, false
}

func deepCopyByGob(src, dst interface{}) error {
	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(&buffer).Decode(dst)
}

func deepCopyByReflect(src, dst interface{}) error {
	srcValue := reflect.ValueOf(src)
	dstValue := reflect.ValueOf(dst)

	if srcValue.Kind() != dstValue.Kind() {
		return fmt.Errorf("src and dst must be of the same type")
	}

	if srcValue.Type().AssignableTo(reflect.TypeOf(dst)) {
		dstValue.Set(srcValue)
		return nil
	}

	switch srcValue.Kind() {
	case reflect.Ptr:
		srcElem := srcValue.Elem()
		dstElem := reflect.New(srcElem.Type()).Elem()
		if err := deepCopyByReflect(srcElem.Interface(), dstElem.Addr().Interface()); err != nil {
			return err
		}
		dstValue.Set(dstElem.Addr())
	case reflect.Struct:
		for i := 0; i < srcValue.NumField(); i++ {
			if err := deepCopyByReflect(srcValue.Field(i).Interface(), dstValue.Field(i).Addr().Interface()); err != nil {
				return err
			}
		}
	case reflect.Slice:
		dstSlice := reflect.MakeSlice(srcValue.Type(), srcValue.Len(), srcValue.Cap())
		for i := 0; i < srcValue.Len(); i++ {
			if err := deepCopyByReflect(srcValue.Index(i).Interface(), dstSlice.Index(i).Addr().Interface()); err != nil {
				return err
			}
		}
		dstValue.Set(dstSlice)
	default:
		if srcValue.CanInterface() {
			dstValue.Set(reflect.ValueOf(srcValue.Interface()))
		} else {
			return fmt.Errorf("src cannot be converted to interface{}")
		}
	}

	return nil
}

func saveImageAsJpeg(fileName string, img *image.RGBA) error {
	outFile, err := os.Create(fileName)
	defer outFile.Close()
	if err != nil {
		fmt.Printf("error occur: %s", err.Error())
		return err
	}
	/*
		b := bufio.NewWriter(outFile)
		var opt jpeg.Options
		opt.Quality = 100
		err = jpeg.Encode(b, img, &opt)
		if err != nil {
			fmt.Printf("error occur: %s", err.Error())
			return err
		}
		err = b.Flush()
		if err != nil {
			fmt.Printf("error occur: %s", err.Error())
			return err
		}

	*/
	var opt jpeg.Options
	opt.Quality = 100
	// ok, write out the data into the new JPEG file
	err = jpeg.Encode(outFile, img, &opt) // put quality to 80%
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return nil
}

func saveImageAsPng(fileName string, img *image.RGBA) error {
	outFile, err := os.Create(fileName)
	defer outFile.Close()
	if err != nil {
		fmt.Printf("error occur: %s", err.Error())
		return err
	}

	err = png.Encode(outFile, img)
	if err != nil {
		fmt.Printf("error occur: %s", err.Error())
		return err
	}
	return nil
}

func saveImageAsGif(fileName string, img *image.RGBA) error {
	outFile, err := os.Create(fileName)
	defer outFile.Close()
	if err != nil {
		fmt.Printf("error occur: %s", err.Error())
		return err
	}

	var opt gif.Options
	opt.NumColors = 256 // you can add more parameters if you want
	// ok, write out the data into the new GIF file
	err = gif.Encode(outFile, img, &opt) // put num of colors to 256
	if err != nil {
		fmt.Printf("error occur: %s", err.Error())
		return err
	}

	return nil
}
