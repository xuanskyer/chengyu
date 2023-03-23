package chengyu

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

type Blank struct {
	Head           int         `json:"head"`
	Foot           int         `json:"foot"`
	HeadFoot       []BlankItem `json:"head_foot"`
	HeadUseCyIndex int         `json:"head_use_cy_index"` //Head 使用第几个成语匹配
	FootUseCyIndex int         `json:"foot_use_cy_index"` //foot 使用第几个成语匹配
}

type BlankItem struct {
	Head           int `json:"head"`
	HeadUseCyIndex int `json:"head_use_cy_index"` //Head 使用第几个成语匹配
	FootUseCyIndex int `json:"foot_use_cy_index"` //foot 使用第几个成语匹配
	Foot           int `json:"foot"`
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
	fmt.Println("check: ", ones, setting, count)
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

	fmt.Println("check2: ", len(onesMap), lengOnes, onesMap, ones)
	var c1, c2 string
	var e1, e2 error
	for _, info := range setting {
		if len(info.HeadFoot) > 0 {
			for _, val := range info.HeadFoot {
				c1, e1 = GetChengyuPosStr(val.Head-1, val.Head, ones[val.HeadUseCyIndex])
				c2, e2 = GetChengyuPosStr(val.Foot-1, val.Foot, ones[val.FootUseCyIndex])
				if e1 != nil || e2 != nil {
					fmt.Println(e1, e2)
				}
				if c1 == "" || c2 == "" || c1 != c2 || e1 != nil || e2 != nil {
					return false
				}
			}
		} else {

			c1, e1 = GetChengyuPosStr(info.Head-1, info.Head, ones[info.HeadUseCyIndex])
			c2, e2 = GetChengyuPosStr(info.Foot-1, info.Foot, ones[info.FootUseCyIndex])
			if e1 != nil || e2 != nil {
				fmt.Println(e1, e2)
			}
			if c1 == "" || c2 == "" || c1 != c2 || e1 != nil || e2 != nil {
				return false
			}
		}
	}
	return true
}

func RecursionGenerate(chengYuMap map[string]bool, blankSetting []Blank, validCount, depth int, selectedOnes []string, result *[][]string, selectedMap map[string]bool) {

	//fmt.Println("begin: ", selectedOnes, depth)
	onesMap := make(map[string]bool, 0)
	for _, one := range selectedOnes {
		onesMap[one] = true
	}
	if len(onesMap) != len(selectedOnes) {
		return
	}
	//fmt.Println("begin2: ", selectedOnes, depth, len(onesMap), len(selectedOnes), validCount, len(blankSetting))
	if len(selectedOnes) == validCount || depth > len(blankSetting)+1 {
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
ChengyuMapFor:
	for c := range chengYuMap {
		if depth > 0 {
			blank := blankSetting[depth-1]
			if len(blank.HeadFoot) > 0 {
				hitCount := 0
				for _, val := range blank.HeadFoot {
					cell1, err1 = GetChengyuPosStr(val.Head-1, val.Head, selectedOnes[val.HeadUseCyIndex])
					cell2, err2 = GetChengyuPosStr(val.Foot-1, val.Foot, c)
					//fmt.Printf("aaaa: %d, %v, %+v, cell1: %v,cell2: %v, old: %v, new: %v \n", depth, selectedOnes, blank, cell1, cell2, selectedOnes[val.HeadUseCyIndex], c)
					if err2 != nil || cell2 == "" || cell1 == "" || cell1 != cell2 || err1 != nil {
						continue ChengyuMapFor
					} else {
						hitCount++
					}
				}
				if hitCount == len(blank.HeadFoot) {
					goto HitAndRecursion
				}
			} else {
				if blank.HeadUseCyIndex > len(selectedOnes)-1 {
					return
				}
				cell2, err2 = GetChengyuPosStr(blank.Foot-1, blank.Foot, c)
				if err2 != nil || cell2 == "" {
					continue
				}
				cell1, err1 = GetChengyuPosStr(blank.Head-1, blank.Head, selectedOnes[blank.HeadUseCyIndex])
				if cell1 == "" || cell2 == "" || cell1 != cell2 || err1 != nil || err2 != nil {
					continue
				}
			}
		}

		// 递归处理下一个空白处
	HitAndRecursion:
		RecursionGenerate(chengYuMap, blankSetting, validCount, depth+1, append(selectedOnes, c), result, selectedMap)
	}
}
