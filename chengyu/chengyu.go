package chengyu

import (
	"errors"
)

type Blank struct {
	Head int `json:"head"`
	Foot int `json:"foot"`
}


func GetChengyuPosStr(begin, end int, item string) (string, error) {
	if begin < 0 || end > 4 || begin > end || begin+1 != end {
		return "", errors.New("位置非法")
	}
	rs := []rune(item)
	lth := len(rs)
	if lth < end {
		return "", errors.New("err")
	}
	return string(rs[begin:end]), nil
}

func getChengyu(list []string, begin, end int, char string) ([]string, string, error) {
	if begin < 0 || end > 4 || begin > end || begin+1 != end {
		return list, "", errors.New("位置非法")
	}
	originLen := len(list)
	for index, val := range list {
		rs := []rune(val)
		lth := len(rs)
		if lth < end {
			continue
		}
		if string(rs[begin:end]) == char {
			if index+1 == originLen {
				list = list[:originLen-1]
			} else {
				list = append(list[0:index], list[index+1:]...)
			}
			return list, val, nil
		}
	}
	return list, "", errors.New("not found")
}


func Check(ones []string, setting []Blank) bool {

	if len(ones) != len(setting)+1 {
		return false
	}
	mapTrue := make(map[string]bool, 0)
	for _, item := range ones {
		mapTrue[item] = true
	}

	if len(mapTrue) != len(ones) {
		return false
	}
	var cell1, cell2 string

	for index, info := range setting {
		cell1, _ = GetChengyuPosStr(info.Head-1, info.Head, ones[index])
		cell2, _ = GetChengyuPosStr(info.Foot-1, info.Foot, ones[index+1])
		if cell1 != cell2 {
			return false
		}
	}
	return true
}