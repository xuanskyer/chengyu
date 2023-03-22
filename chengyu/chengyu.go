package chengyu

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

type Blank struct {
	Head           int `json:"head"`
	Foot           int `json:"foot"`
	HeadUseCyIndex int `json:"head_use_cy_index"` //Head 使用第几个成语匹配
	FootUseCyIndex int `json:"foot_use_cy_index"` //foot 使用第几个成语匹配
}

func GenerateResult(chengYuMap map[string]bool, blankSetting []Blank, validCount, depth int, selectedOnes []string, result *[][]string, selectedMap map[string]bool) {

	onesMap := make(map[string]bool, 0)
	for _, one := range selectedOnes {
		onesMap[one] = true
	}
	if len(onesMap) != len(selectedOnes) {
		return
	}
	if len(selectedOnes) == validCount || depth >= validCount {
		// 已填好所有空白处配置，判断所选成语序列是否符合条件
		_, ok := selectedMap[fmt.Sprint(selectedOnes)]
		if !ok {
			if Check(selectedOnes, blankSetting, validCount) {
				//fmt.Println("answer: ", selectedOnes)
				*result = append(*result, selectedOnes)
			}
			selectedMap[fmt.Sprint(selectedOnes)] = true
		}
		return
	}
	var cell1, cell2 string
	var err1, err2 error
	// 处理当前递归层
	if depth > 0 {
		blank := blankSetting[depth-1]
		cell1, err1 = GetChengyuPosStr(blank.Head-1, blank.Head, selectedOnes[depth-1])
	}
	for c := range chengYuMap {
		if depth > 0 {
			blank := blankSetting[depth-1]
			//cell1, err1 := GetChengyuPosStr(blank.Head-1, blank.Head, selectedOnes[depth-1])
			cell2, err2 = GetChengyuPosStr(blank.Foot-1, blank.Foot, c)
			if err1 != nil || err2 != nil {
				fmt.Println(err1, err2)
			}
			if cell1 == "" || cell2 == "" || cell1 != cell2 || err1 != nil || err2 != nil {
				continue
			}
		}
		// 递归处理下一个空白处
		GenerateResult(chengYuMap, blankSetting, validCount, depth+1, append(selectedOnes, c), result, selectedMap)
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

	lengOnes := len(ones)
	if lengOnes != count {
		return false
	}
	onesMap := make(map[string]bool, 0)
	for _, one := range ones {
		onesMap[one] = true
	}
	if len(onesMap) != lengOnes {
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
		if c1 == "" || c2 == "" || c1 != c2 || e1 != nil || e2 != nil {
			return false
		}
	}
	return true
}

func RecursionGenerate(chengYuMap map[string]bool, blankSetting []Blank, validCount, depth int, selectedOnes []string, result *[][]string, selectedMap map[string]bool) {

	onesMap := make(map[string]bool, 0)
	for _, one := range selectedOnes {
		onesMap[one] = true
	}
	if len(onesMap) != len(selectedOnes) {
		return
	}
	if len(selectedOnes) == validCount || depth >= len(blankSetting) {
		// 已填好所有空白处配置，判断所选成语序列是否符合条件
		key := strings.Join(selectedOnes, ",")
		_, ok := selectedMap[key]
		if !ok {
			if Check(selectedOnes, blankSetting, validCount) {
				*result = append(*result, selectedOnes)
			}
			selectedMap[key] = true
		}
		return
	}
	var cell1, cell2 string
	var err1, err2 error
	// 处理当前递归层
	if depth > 0 {
		//{HeadUseCyIndex: 0, Head: 2, FootUseCyIndex: 1, Foot: 1},
		//{HeadUseCyIndex: 0, Head: 3, FootUseCyIndex: 2, Foot: 1},
		//{HeadUseCyIndex: 0, Head: 4, FootUseCyIndex: 3, Foot: 1},
		//{HeadUseCyIndex: 4, Head: 1, FootUseCyIndex: 1, Foot: 2},
		//{HeadUseCyIndex: 4, Head: 2, FootUseCyIndex: 2, Foot: 2},
		//{HeadUseCyIndex: 4, Head: 3, FootUseCyIndex: 3, Foot: 2},
		//{HeadUseCyIndex: 5, Head: 2, FootUseCyIndex: 1, Foot: 3},
		//{HeadUseCyIndex: 5, Head: 3, FootUseCyIndex: 2, Foot: 3},
		//{HeadUseCyIndex: 5, Head: 4, FootUseCyIndex: 3, Foot: 3},
		//{HeadUseCyIndex: 6, Head: 1, FootUseCyIndex: 1, Foot: 4},
		//{HeadUseCyIndex: 6, Head: 2, FootUseCyIndex: 2, Foot: 4},
		//{HeadUseCyIndex: 6, Head: 3, FootUseCyIndex: 3, Foot: 4},
		blank := blankSetting[depth-1]
		fmt.Printf("%d, %v, %+v, %v \n", depth, selectedOnes, blank, selectedOnes[blank.HeadUseCyIndex])
		if blank.HeadUseCyIndex > len(selectedOnes)-1 {
			return
		}
		cell1, err1 = GetChengyuPosStr(blank.Head-1, blank.Head, selectedOnes[blank.HeadUseCyIndex])

		fmt.Printf("%d, %v, %+v, %v, %v \n", depth, selectedOnes, blank, selectedOnes[blank.HeadUseCyIndex], cell1)
	}
	for c := range chengYuMap {
		if depth > 0 {
			fmt.Println(c)
			blank := blankSetting[depth-1]
			//cell1, err1 := GetChengyuPosStr(blank.Head-1, blank.Head, selectedOnes[depth-1])
			cell2, err2 = GetChengyuPosStr(blank.Foot-1, blank.Foot, c)
			if err1 != nil || err2 != nil {
				fmt.Println(err1, err2)
			}
			if cell1 == "" || cell2 == "" || cell1 != cell2 || err1 != nil || err2 != nil {
				continue
			}
		}
		// 递归处理下一个空白处
		RecursionGenerate(chengYuMap, blankSetting, validCount, depth+1, append(selectedOnes, c), result, selectedMap)
	}
}
