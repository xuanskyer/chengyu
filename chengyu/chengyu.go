package chengyu

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

type Blank struct {
	Head int `json:"head"`
	Foot int `json:"foot"`
}

func GenerateResult(chengYuMap map[string]bool, blankSetting []Blank, validCount, depth int, selectedOnes []string, result *[][]string) {
	if depth == validCount {
		// 已填好所有空白处配置，判断所选成语序列是否符合条件
		if Check(selectedOnes, blankSetting, validCount) {
			//fmt.Println("answer: ", selectedOnes)
			*result = append(*result, selectedOnes)
			selectedOnes = make([]string, 0)
		}
		return
	}
	// 处理当前递归层
	for c := range chengYuMap {
		if depth > 0 {
			blank := blankSetting[depth-1]
			cell1, err1 := GetChengyuPosStr(blank.Head-1, blank.Head, selectedOnes[depth-1])
			cell2, err2 := GetChengyuPosStr(blank.Foot-1, blank.Foot, c)
			if err1 != nil || err2 != nil {
				fmt.Println(err1, err2)
			}
			if cell1 != "" && cell2 != "" && cell1 != cell2 {
				continue
			}
		}
		// 递归处理下一个空白处
		GenerateResult(chengYuMap, blankSetting, validCount, depth+1, append(selectedOnes, c), result)
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
	han := string(rs[begin:end])
	_, size := utf8.DecodeRuneInString(han)
	if size != len(han) {
		return "", errors.New("err: Invalid Chinese Characters")
	}
	return han, nil
}

func Check(ones []string, setting []Blank, count int) bool {

	if len(ones) != count {
		return false
	}
	var c1, c2 string
	var e1, e2 error

	for index, info := range setting {
		c1, e1 = GetChengyuPosStr(info.Head-1, info.Head, ones[index])
		c2, e2 = GetChengyuPosStr(info.Foot-1, info.Foot, ones[index+1])
		if e1 != nil || e2 != nil {
			fmt.Println(e1, e2)
		}
		if c1 != "" && c2 != "" && c1 != c2 {
			return false
		}
	}
	return true
}
