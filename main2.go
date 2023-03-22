package main

import (
	"fmt"
	"github.com/chengyu/chengyu"
	"os"
	"time"
)

const (
	MaxX = 10
	MaxY = 6

	p9rNil   = 0 //占位符状态：未使用
	p9rBlank = 1 //占位符状态：空白
	p9rUsed  = 2 //占位符状态：有字
)

type Setting struct {
	Sort int `json:"sort"`
}

type ChengYuCell struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type ChengYu []ChengYuCell

func main() {
	start := time.Now()
	var ChengYuList = []string{
		"门当户对", "声名狼藉", "时过境迁", "念念不忘",
		"当声过念", "户名境念", "对狼迁不",
	}
	table := [][]int{
		{0, 0, 0, 0, 0, 0},
		{0, 1, 0, 1, 0, 0},
		{0, 1, 1, 1, 1, 0},
		{0, 1, 1, 1, 1, 0},
		{0, 1, 1, 1, 1, 0},
		{0, 0, 1, 0, 1, 0},
	}
	fmt.Println("地图表格：")
	allCY := []ChengYu{}
	allLineCY := []ChengYu{}
	allColCY := []ChengYu{}
	for index, item := range table {
		fmt.Println(item)
		cy := getChengYu(index, item, false)
		if len(cy) > 0 {
			allLineCY = append(allLineCY, cy...)
		}
	}
	for i := 0; i < MaxY; i++ {
		column, _ := getSliceXN(table, i)
		cyCol := getChengYu(i, column, true)
		fmt.Println("列： ", i, "column: ", column, "成语： ", cyCol)
		if len(cyCol) > 0 {
			allColCY = append(allColCY, cyCol...)
		}
	}
	fmt.Println("所有行中的成语： ")
	for _, item := range allLineCY {
		fmt.Println(item)
	}
	fmt.Println("所有列中的成语： ")
	for _, item := range allColCY {
		fmt.Println(item)
	}
	allCY = append(append(allCY, allColCY...), allLineCY...)

	fmt.Println("所有成语： ")
	for index, item := range allCY {
		fmt.Println(index, item)
	}
	v2Setting := []chengyu.Blank{
		{HeadUseCyIndex: 0, Head: 2, FootUseCyIndex: 0, Foot: 1},
		{HeadUseCyIndex: 0, Head: 3, FootUseCyIndex: 0, Foot: 1},
		{HeadUseCyIndex: 0, Head: 4, FootUseCyIndex: 0, Foot: 1},
		{HeadUseCyIndex: 1, Head: 2, FootUseCyIndex: 0, Foot: 1},
		{HeadUseCyIndex: 1, Head: 3, FootUseCyIndex: 0, Foot: 2},
		{HeadUseCyIndex: 1, Head: 3, FootUseCyIndex: 0, Foot: 2},
		{HeadUseCyIndex: 2, Head: 2, FootUseCyIndex: 0, Foot: 3},
		{HeadUseCyIndex: 2, Head: 3, FootUseCyIndex: 0, Foot: 3},
		{HeadUseCyIndex: 2, Head: 4, FootUseCyIndex: 0, Foot: 3},
		{HeadUseCyIndex: 3, Head: 1, FootUseCyIndex: 0, Foot: 4},
		{HeadUseCyIndex: 3, Head: 2, FootUseCyIndex: 0, Foot: 4},
		{HeadUseCyIndex: 3, Head: 3, FootUseCyIndex: 0, Foot: 4},
	}
	fmt.Printf("成语列表(总数：%d)： %v\n", len(ChengYuList), ChengYuList)
	fmt.Println("空白处配置： ", v2Setting)

	// 使用 map 存储成语列表，方便去重
	chengYuMap := make(map[string]bool)
	for _, item := range ChengYuList {
		chengYuMap[item] = true
	}
	f, _ := os.OpenFile("result2.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	_, _ = fmt.Fprintln(f, "")
	// 递归处理，开始生成成语序列并判断
	result := [][]string{}
	selectedMap := make(map[string]bool, 0)
	chengyu.RecursionGenerate(chengYuMap, v2Setting, len(allCY), 0, []string{}, &result, selectedMap)

	filter := [][]string{}
	filterMap := make(map[string]bool, 0)
	for _, val := range result {
		key := fmt.Sprint(val)
		if _, ok := filterMap[key]; ok {
			continue
		}
		if chengyu.Check(val, v2Setting, len(v2Setting)+1) {
			filter = append(filter, val)
		}
		filterMap[key] = true
	}
	//每次执行前，先清空文件内容
	f, _ = os.OpenFile("result2.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	defer f.Close()
	_, _ = fmt.Fprintln(f, filter)

	elapsed := time.Since(start)
	fmt.Println("该函数执行完成耗时：", elapsed)
	fmt.Println("resultCount: ", len(result), "filterCount: ", len(filter))
}

// 竖向取二维切片的第N列
func getSliceXN(table [][]int, col int) ([]int, error) {
	var column []int
	for i := 0; i < len(table); i++ {
		column = append(column, table[i][col])
	}
	return column, nil
}

// 从第n 行/列 取出成语（从0开始）
func getChengYu(n int, slice []int, fixLine bool) []ChengYu {
	cyList := []ChengYu{}
	chengYu := ChengYu{}
	count := 0
	lenSlice := len(slice)
	for index, val := range slice {
		if count == 4 {
			cyList = append(cyList, chengYu)
		}
		if val == p9rNil || count == 4 {
			count = 0
			chengYu = ChengYu{}
			continue
		} else {
			var cell ChengYuCell
			if fixLine {
				cell = ChengYuCell{n, index}
			} else {
				cell = ChengYuCell{index, n}
			}
			chengYu = append(chengYu, cell)
			count++
			if count == 4 && index+1 == lenSlice {
				cyList = append(cyList, chengYu)
			}
		}
	}
	return cyList
}
