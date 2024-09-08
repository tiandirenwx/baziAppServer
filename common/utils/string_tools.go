package utils

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"trpc.group/trpc-go/trpc-go/log"
)

var formatStr = "2006-01-02 15:04:05"

// DefaultTimeStr 默认空值时间
var DefaultTimeStr = "0000-00-00 00:00:00"

// Get32RandomString
func Get32RandomString() string {
	return GetRandomString(32)
}

// Get64RandomString
func Get64RandomString() string {
	return GetRandomString(64)
}

// GetRandomString
func GetRandomString(n int) string {
	str, _ := RandomStr(n)
	return str
}

// GetRandomNumString
func GetRandomNumString(n int) string {
	str, _ := RandomNumStr(n)
	return str
}

// HashPassword 哈希加密密码
func HashSha256(password string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
}

// HashAndSalt 通过bcrypt加salt加密
func HashAndSalt(ctx context.Context, pwdStr string) string {
	pwd := []byte(pwdStr)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.ErrorContext(ctx, err)
		return ""
	}
	return string(hash)
}

// stringToUint32
func stringToUint32(s string) uint32 {
	num, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0
	}
	return uint32(num)
}

// stringToUint64
func stringToUint64(s string) uint64 {
	num, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return num
}

// stringToInt32
func stringToInt32(s string) int32 {
	num, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0
	}
	return int32(num)
}

// stringToInt32
func StringToInt32(s string) int32 {
	num, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0
	}
	return int32(num)
}

// Int32ToString
func Int32ToString(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

// Int64ToString ...
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// stringToInt
func stringToInt(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return num
}

// inetAtoN
func inetAtoN(ip string) uint32 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return uint32(ret.Int64())
}

/*
*
返回日期字符 2020-10-23 12:00:00
*/
func TimeDateString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// TimeToDateString 格式化时间为 yyyy-mm-dd HH:MM:ss 字符串
func TimeToDateString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// TimeFromString 把字符串格式时间对象转成 time对象
func TimeFromString(t string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", t)
}

// TimeMustFromString 确定 t 的格式一定正确的时候使用
func TimeMustFromString(t string) time.Time {
	ret, _ := TimeFromString(t)
	return ret
}

// TimeNowAddDaysAndTruncate 在当前时间基础上偏移days 天，并取0点（按0时区）
func TimeNowAddDaysAndTruncate(days int32) time.Time {
	return time.Now().Local().Add(time.Hour * 24 * time.Duration(days)).Truncate(time.Hour * 24)
}

/*
*
uint32转字符
*/
func Uint32ToString(n uint32) string {
	return strconv.FormatInt(int64(n), 10)
}

// Uint64ToString ...
func Uint64ToString(n uint64) string {
	return strconv.FormatInt(int64(n), 10)
}

// StringToUint32 ...
func StringToUint32(n string) uint32 {
	intNum, _ := strconv.Atoi(n)
	return uint32(intNum)
}

/*
*
字符数组转interface
*/
func StringListToInterface(stringList []string) []interface{} {
	var field []interface{}
	for _, v := range stringList {
		field = append(field, v)
	}
	return field
}

// RandomStr 生成指定长度的随机字符串 真随机
func RandomStr(n int) (string, error) {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	l := len(letters)

	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = byte(letters[int(b[i])%l])
	}
	return string(b), nil
}

// RandomNumStr 生成指定长度的随机数字字符串
func RandomNumStr(n int) (string, error) {
	letters := []rune("0123456789")
	l := len(letters)

	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = byte(letters[int(b[i])%l])
	}
	return string(b), nil
}

// GenerateRandomLocal 本地生成真随机数
func GenerateRandomLocal(numberOfBytes uint64) (string, error) {
	if numberOfBytes%8 != 0 {
		return "", fmt.Errorf("unsurported len, numberOfBytes mod 8 must be 0")
	}

	randomBytes := []byte{}
	round := numberOfBytes / 8
	for i := uint64(0); i < round; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(9223372036854775807))
		randomBytes = bytes.Join([][]byte{randomBytes, n.Bytes()}, []byte(""))
	}

	return base64.StdEncoding.EncodeToString(randomBytes), nil
}

// DateFormatToString 2022-11-03T00:00:00+08:00返回日期字符 2022-11-03 00:00:00
func DateFormatToString(date string) string {
	createTime, _ := time.Parse("2006-01-02T15:04:05+08:00", date)
	if createTime.IsZero() == true {
		log.Infof("timeIsZero;date:%s", date)
		return ""
	}
	return createTime.Format(formatStr)
}

// 密码hash,对应php
func PasswordHash(passwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(passwd), 12)
	return string(bytes), err
}

// 密码验证,对应php
func PasswordVerify(passwd, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd))
	return err == nil
}

// FuncDateTimeStringFormat 格式化含有T的时间串
func FuncDateTimeStringFormat(timeValue string) string {
	//不含有T的不格式化
	if timeValue == "" || !strings.Contains(timeValue, "T") {
		return timeValue
	}
	format := "2006-01-02 15:04:05"
	timeLayout := "2006-01-02T15:04:05+08:00"                      //时间字符串
	loc, _ := time.LoadLocation("Local")                           //***获取时区***
	theTime, _ := time.ParseInLocation(timeLayout, timeValue, loc) //使用模板在对应时区转化为time.time类型
	// 0001-01-01T00:00:00Z这里也表示时间为null
	if theTime.IsZero() {
		return ""
	} else {
		//时间戳转日期
		dataTimeStr := theTime.Format(format) //使用模板格式化为日期字符串
		return dataTimeStr
	}
}

// StringSliceContains contains
func StringSliceContains(ss []string, k string) bool {
	for _, s := range ss {
		if k == s {
			return true
		}
	}
	return false
}

// IsEmailLowercase 邮箱是否是全小写
// 如果不是邮箱则返回 false & ""
// 如果不是则返回 false & 小写邮箱地址
// 如果是则返回 true & 小写邮箱地址
func IsEmailLowercase(email string) (bool, string) {
	// 判断邮箱是否全小写
	lowerCaseEmail := strings.ToLower(email)
	if lowerCaseEmail == email {
		return false, lowerCaseEmail
	}
	return true, lowerCaseEmail
}

// BatchTransEmail2Lowercase 邮箱批量转换成小写
func BatchTransEmail2Lowercase(ctx context.Context, emailList []string) []string {

	var lowerCaseEmailList []string

	for _, email := range emailList {
		lowerCaseEmailList = append(lowerCaseEmailList, strings.ToLower(email))
	}

	return lowerCaseEmailList
}

// DistinctString 去重
func DistinctString(ss []string) []string {
	size := len(ss)
	if size <= 1 {
		return ss
	}
	if size == 2 {
		if ss[0] == ss[1] {
			return ss[0:1]
		}
		return ss
	}

	var ret []string
	dict := make(map[string]bool, size)
	for _, el := range ss {
		if _, ok := dict[el]; !ok {
			ret = append(ret, el)
			dict[el] = true
		}
	}
	return ret
}
