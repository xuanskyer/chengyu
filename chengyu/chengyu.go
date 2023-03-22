package chengyu

import (
	"errors"
	"fmt"
	"os"
)

type Blank struct {
	Head int `json:"head"`
	Foot int `json:"foot"`
}

func GenerateChengYu(chengYuMap map[string]bool, blankSetting []Blank, depth int, selectedOnes []string) {
	if depth == len(blankSetting)+1 {
		// 已填好所有空白处配置，判断所选成语序列是否符合条件
		if ok := Check(selectedOnes, blankSetting); ok {
			//fmt.Println("answer: ", selectedOnes)
			// 打开文件以进行附加写入，如果文件不存在，则创建它
			file, err := os.OpenFile("result.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println("Error opening file:", err)
				return
			}
			defer file.Close()

			// 写入字符串到文件中
			_, err = fmt.Fprintln(file, selectedOnes)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
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
		GenerateChengYu(chengYuMap, blankSetting, depth+1, append(selectedOnes, c))
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