package convert

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/simplifiedchinese"
)

//float64 - float64
func Float64DivFloat64(a float64, b float64) float64 {
	da := decimal.NewFromFloat(a)
	db := decimal.NewFromFloat(b)
	res, _ := da.Div(db).Float64()
	return res
}

func IntDivInt(a int, b int) int {
	da := decimal.New(int64(a), 0)
	db := decimal.New(int64(b), 0)
	return int(da.Sub(db).IntPart())
}

func Float64AddFloat64(a float64, b float64) float64 {
	da := decimal.NewFromFloat(a)
	db := decimal.NewFromFloat(b)
	aa, _ := da.Add(db).Float64()
	return aa
}

func Float64SubFloat64(a float64, b float64) string {
	da := decimal.NewFromFloat(a)
	db := decimal.NewFromFloat(b)
	res, _ := da.Sub(db).Float64()

	return fmt.Sprintf("%.2f", res)
}

func ObjToJson(obj interface{}) string {

	if obj == nil {
		return ""
	}

	b, err := json.Marshal(obj)
	if err != nil {
		logrus.Info("ObjToJson, error, ", err)
		return ""
	}

	return string(b)
}

func JsonToObj(jsonString string, obj interface{}) error {

	if len(jsonString) == 0 {
		return errors.New("JSON字符串为空")
	}

	err := json.Unmarshal([]byte(jsonString), &obj)
	if err != nil {
		logrus.Info("JsonToObj, error, ", err)
		return err
	}

	return nil
}

func ConvertStr2GBK(str string) string {
	data, err := simplifiedchinese.GBK.NewEncoder().String(str)
	if err != nil {
		logrus.Error("ConvertStr2GBK err, ", err)
		return ""
	}

	return data
}

func ConvertGBK2Str(gbkStr string) string {
	data, err := simplifiedchinese.GBK.NewDecoder().String(gbkStr)
	if err != nil {
		logrus.Error("ConvertGBK2Str err, ", err)
		return ""
	}

	return data
}

//string转int
func StringToInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return n
}

//string转int32
func StringToInt32(s string) int32 {
	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0
	}

	return int32(n)
}

//string转int64
func StringToInt64(s string) int64 {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}

	return n
}

//convert number to big.Int type
func StringToBigInt(s string) *big.Int {
	// convert number to big.Int type
	ip := new(big.Int)
	ip.SetString(s, 10) //base 10

	return ip
}

//string转float32
func StringToFloat32(s string) float32 {
	if f, err := strconv.ParseFloat(s, 32); err == nil {
		return float32(f)
	}

	return 0
}

//string转float64
func StringToFloat64(s string) float64 {
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}

	return 0
}

// RandomString 返回指定位数随机字符串
func RandomString(size int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// StrToMd5 字符串转md5
func StrToMd5(str string) string {
	has := md5.Sum([]byte(str))
	md5str := fmt.Sprintf("%x", has)
	return md5str
}
