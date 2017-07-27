package format

import (
	"strings"
	"strconv"
	"errors"
	"fmt"
)

var CnNums [10]rune = [10]rune{rune('零'), rune('一'), rune('二'), rune('三'), rune('四'),
							   rune('五'), rune('六'), rune('七'), rune('八'), rune('九')}

var CnUnit [5]rune = [5]rune{rune('十'), rune('百'), rune('千'), rune('万'), rune('亿'), }

func ChineseNumberToInt(s string) (res int, err error) {
	// TODO: NOT IMPLEMENT
	return 0, nil
}

// ChineseSuffixStringToInt transform "1.28亿" to corresponding integer
func ChineseSuffixStringToInt(s string) (res int64, err error) {
	r := []rune(s)
	n := len(r)

	var mutiplier float64;
	switch r[n-1] {
	case rune('万'):
		mutiplier = 10000
		r = r[0:n-1]
	case rune('亿'):
		mutiplier = 100000000
		r = r[0:n-1]
	default:
		mutiplier = 1
	}

	numStr := string(r)
	if dotInd := strings.Index(numStr, "."); dotInd == -1 {
		// 没有小数点
		if i, err := strconv.Atoi(numStr); err != nil {
			return 0, err
		} else {
			return int64(float64(i) * mutiplier), nil
		}
	} else {
		// 有小数点,判断小数位数并移除小数点
		for i := 0; i < len(numStr)-dotInd-1; i++ {
			mutiplier /= 10
		}

		numStr = strings.Replace(numStr, ".", "", 1)
		if i, err := strconv.Atoi(numStr); err != nil {
			return 0, err
		} else {
			return int64(float64(i) * mutiplier), nil
		}
	}
}

// PrefixedBytesToInt用于将形如"128k" 转换为相应的字节数
func PrefixedBytesToInt(s string) (res int64, err error) {
	var i, nFrac int
	var val int64
	var c byte
	var dot bool

	// parse numeric val (omit dot), and length of frac part
Loop:
	for i < len(s) {
		c = s[i]
		switch {
		case '0' <= c && c <= '9':
			val *= 10
			val += int64(c - '0')
			if dot {
				nFrac ++
			}
			i++
		case c == '.':
			dot = true
			i++
		default:
			break Loop
		}
	}

	fmt.Println(val)
	unit := strings.ToUpper(strings.TrimSpace(s[i:]))

	switch unit {
	case "", "B":
	case "KB", "K":
		val <<= 10
	case "MB", "M":
		val <<= 20
	case "GB", "G":
		val <<= 30
	case "TB", "T":
		val <<= 40
	case "PB", "P":
		val <<= 50
	case "EB", "E":
		val <<= 60
	default:
		return 0, errors.New("invalid prefix")
	}

	// handle frac
	for j := 0; j < nFrac; j++ {
		val /= 10
	}

	return val, nil
}
