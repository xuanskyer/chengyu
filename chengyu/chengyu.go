package chengyu

import (
	"errors"
)

type Blank struct {
	Head int `json:"head"`
	Foot int `json:"foot"`
}

func GenerateResult(chengYuMap map[string]bool, blankSetting []Blank, depth int, selectedOnes []string, result *[][]string) {
	if depth == len(blankSetting)+1 {
		// 已填好所有空白处配置，判断所选成语序列是否符合条件
		if ok := Check(selectedOnes, blankSetting); ok {
			//fmt.Println("answer: ", selectedOnes)
			*result = append(*result, selectedOnes)
		}
		return
	}
	// 处理当前递归层
	for c := range chengYuMap {
		if depth > 0 {
			blank := blankSetting[depth-1]
			cell1, _ := GetChengyuPosStr(blank.Head-1, blank.Head, selectedOnes[depth-1])
			cell2, _ := GetChengyuPosStr(blank.Foot-1, blank.Foot, c)
			if cell1 != cell2 {
				continue
			}
		}
		// 递归处理下一个空白处
		GenerateResult(chengYuMap, blankSetting, depth+1, append(selectedOnes, c), result)
	}
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
