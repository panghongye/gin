package lib

import (
	"math/rand"
	"time"
)

func GetRandomString(size int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// func StrMd5(str string) string {
	// return "333"
	//has := md5.Sum([]byte(str))
	//md5str := fmt.Sprintf("%x", has)
	//return md5str
// }
