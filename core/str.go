package core

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func Hex2bin(s string) []byte {
	ret, _ := hex.DecodeString(s)
	return ret
}

func Bin2Hex(b []byte) string {
	return hex.EncodeToString(b)
}

func EncodeIdToBase64(input string) string {
	s := base64.StdEncoding.EncodeToString(Hex2bin(input))
	base64Str := strings.Map(func(r rune) rune {
		switch r {
		case '+':
			return '-'
		case '/':
			return '_'
		}

		return r
	}, s)

	base64Str = strings.ReplaceAll(base64Str, "=", "")
	return base64Str
}

func DecodeIdToBase64(input string) string {
	base64Str := strings.Map(func(r rune) rune {
		switch r {
		case '-':
			return '+'
		case '_':
			return '/'
		}

		return r
	}, input)

	if pad := len(base64Str) % 4; pad > 0 {
		base64Str += strings.Repeat("=", 4-pad)
	}

	b, _ := base64.StdEncoding.DecodeString(base64Str)
	return Bin2Hex(b)
}

func UrlsafeBase64Encode(input string) string {
	s := base64.StdEncoding.EncodeToString([]byte(input))
	base64Str := strings.Map(func(r rune) rune {
		switch r {
		case '+':
			return '-'
		case '/':
			return '_'
		}

		return r
	}, s)

	base64Str = strings.ReplaceAll(base64Str, "=", "")
	return base64Str
}

func Uniqid(prefix string, has_entropy bool) string {
	now := time.Now()
	sec := now.Unix()
	usec := now.UnixNano() % 0x100000

	rand.Seed(time.Now().UnixNano())

	result := ""
	if has_entropy {
		min := 0.0
		max := 10.0
		entropy := min + rand.Float64()*(max-min)
		result = fmt.Sprintf("%s%08x%05x%.8f", prefix, sec, usec, entropy)
	} else {
		result = fmt.Sprintf("%s%08x%05x", prefix, sec, usec)
	}

	return result
}

func Microtime() string {
	now := time.Now()
	zr := fmt.Sprintf("%.8f", float64(now.Nanosecond())/float64(time.Second))
	ts := fmt.Sprintf("%s %d", zr, now.UnixNano()/int64(time.Second))
	return ts
}
