package util

import (
	"NULL/casbin/pkg/setting"
	"crypto/rand"
	r "math/rand"
	"time"
)

// Setup Initialize the util
func Setup() {
	jwtSecret = []byte(setting.AppSetting.JwtSecret)
}

// 随机字符串
func RandomString(n int, alphabets ...byte) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	var randby bool
	if num, err := rand.Read(bytes); num != n || err != nil {
		r.Seed(time.Now().UnixNano())
		randby = true
	}
	for i, b := range bytes {
		if len(alphabets) == 0 {
			if randby {
				bytes[i] = alphanum[r.Intn(len(alphanum))]
			} else {
				bytes[i] = alphanum[b%byte(len(alphanum))]
			}
		} else {
			if randby {
				bytes[i] = alphabets[r.Intn(len(alphabets))]
			} else {
				bytes[i] = alphabets[b%byte(len(alphabets))]
			}
		}
	}
	return string(bytes)
}

//计算秒数差
func TimeDifference(timeNow, timePast string) int {
	var timeLayoutStr = "2006-01-02 15:04:05"
	stNow, _ := time.Parse(timeLayoutStr, timeNow)   //string转time
	stPast, _ := time.Parse(timeLayoutStr, timePast) //string转time
	if stNow.Before(stPast) {
		stNow, stPast = stPast, stNow
	}
	dif := int(stNow.Sub(stPast).Seconds())
	return dif
}
