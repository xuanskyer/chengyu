package chengyu

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"
)

const (
	TableMaxLen = 9

	P9rNil   = 0 //placeholder 占位符状态：未使用
	P9rBlank = 1 //placeholder 占位符状态：空白
	P9rUsed  = 2 //placeholder 占位符状态：有字

	CyLen = 4 //一个成语的字数
)

type CyCell struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type ChengYu []CyCell

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
		cell1, err1 = GetChengYuPosStr(blank.Head-1, blank.Head, selectedOnes[depth-1])
	}
	for c := range chengYuMap {
		if depth > 0 {
			blank := blankSetting[depth-1]
			//cell1, err1 := GetChengYuPosStr(blank.Head-1, blank.Head, selectedOnes[depth-1])
			cell2, err2 = GetChengYuPosStr(blank.Foot-1, blank.Foot, c)
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

func GetChengYuPosStr(begin, end int, item string) (string, error) {
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
	for _, info := range setting {
		if len(info.HeadFoot) > 0 {
			for _, val := range info.HeadFoot {
				c1, e1 = GetChengYuPosStr(val.Head-1, val.Head, ones[val.HeadUseCyIndex])
				c2, e2 = GetChengYuPosStr(val.Foot-1, val.Foot, ones[val.FootUseCyIndex])
				if e1 != nil || e2 != nil {
					fmt.Println(e1, e2)
				}
				if c1 == "" || c2 == "" || c1 != c2 || e1 != nil || e2 != nil {
					return false
				}
			}
		} else {

			c1, e1 = GetChengYuPosStr(info.Head-1, info.Head, ones[info.HeadUseCyIndex])
			c2, e2 = GetChengYuPosStr(info.Foot-1, info.Foot, ones[info.FootUseCyIndex])
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
			if depth == 6 {
				//fmt.Println("blank: ", blank)
			}
			if len(blank.HeadFoot) > 0 {
				hitCount := 0
				for _, val := range blank.HeadFoot {
					cell1, err1 = GetChengYuPosStr(val.Head-1, val.Head, selectedOnes[val.HeadUseCyIndex])
					cell2, err2 = GetChengYuPosStr(val.Foot-1, val.Foot, c)
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
				//if blank.HeadUseCyIndex > len(selectedOnes)-1 {
				//	return
				//}
				cell2, err2 = GetChengYuPosStr(blank.Foot-1, blank.Foot, c)
				if err2 != nil || cell2 == "" {
					continue
				}
				cell1, err1 = GetChengYuPosStr(blank.Head-1, blank.Head, selectedOnes[blank.HeadUseCyIndex])
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

// 判断模板是否合法
func IsValidTemplate(table [][]int) bool {
	for _, item := range table {
		fmt.Println(item)
		count := 0
		for _, val := range item {
			if val == P9rNil {
				count = 0
			} else {
				count++
			}
		}
		if count > CyLen {
			return false
		}
	}
	for i := 0; i < TableMaxLen; i++ {
		column, _ := GetSliceXN(table, i)
		count := 0
		for _, val := range column {
			if val == P9rNil {
				count = 0
			} else {
				count++
			}
		}
		if count > CyLen {
			return false
		}
	}
	return true
}

// 输出一个结果表格
func PrintResult2Table(one []string, sortedCyPos []ChengYu) {
	if len(one) <= 0 {
		return
	}
	if len(one) != len(sortedCyPos) {
		return
	}

	tableString := [TableMaxLen][TableMaxLen]string{}
	for k, v := range tableString {
		for kk, _ := range v {
			tableString[k][kk] = "  "
		}
	}
	for index, cy := range one {
		point := sortedCyPos[index]
		if len(point) < CyLen || len(cy) < CyLen {
			continue
		}
		word1, _ := GetChengYuPosStr(0, 1, cy)
		tableString[point[0].Y][point[0].X] = word1
		word2, _ := GetChengYuPosStr(1, 2, cy)
		tableString[point[1].Y][point[1].X] = word2
		word3, _ := GetChengYuPosStr(2, 3, cy)
		tableString[point[2].Y][point[2].X] = word3
		word4, _ := GetChengYuPosStr(3, 4, cy)
		tableString[point[3].Y][point[3].X] = word4
	}
	//fmt.Printf("%+v\n", one)
	for _, item := range tableString {
		fmt.Printf("%+v\n", item)
	}
}

// 竖向取二维切片的第N列
func GetSliceXN(table [][]int, col int) ([]int, error) {
	var column []int
	for i := 0; i < len(table); i++ {
		column = append(column, table[i][col])
	}
	return column, nil
}

// 从第n 行/列 取出成语（从0开始）
func GetChengYu(n int, slice []int, fixLine bool) []ChengYu {
	cyList := []ChengYu{}
	chengYu := ChengYu{}
	count := 0
	lenSlice := len(slice)
	for index, val := range slice {
		if count == CyLen {
			cyList = append(cyList, chengYu)
		}
		if val == P9rNil || count == CyLen {
			count = 0
			chengYu = ChengYu{}
			continue
		} else {
			var cell CyCell
			if fixLine {
				cell = CyCell{n, index}
			} else {
				cell = CyCell{index, n}
			}
			chengYu = append(chengYu, cell)
			count++
			if count == CyLen && index+1 == lenSlice {
				cyList = append(cyList, chengYu)
			}
		}
	}
	return cyList
}

// 从表格生成模板配置
func Table2Setting(table [][]int) ([]Blank, []ChengYu, error) {
	setting := []Blank{}
	allLineCY := []ChengYu{}
	allColCY := []ChengYu{}
	for index, item := range table {
		cy := GetChengYu(index, item, false)
		if len(cy) > 0 {
			allLineCY = append(allLineCY, cy...)
		}
	}
	for i := 0; i < TableMaxLen; i++ {
		column, _ := GetSliceXN(table, i)
		cyCol := GetChengYu(i, column, true)
		if len(cyCol) > 0 {
			allColCY = append(allColCY, cyCol...)
		}
	}
	sortIndex := 0
	cyMap := make(map[string]int, 0)
	sortedCyPos := make([]ChengYu, 0)
	crossPoint := make(map[string]bool, 0)
	setSettingPos(allColCY, allLineCY, &setting, &sortIndex, cyMap, &sortedCyPos, crossPoint)

	//剩下的无相交的独立行成语计入
	for _, aloneLine := range allLineCY {
		keyAloneLine := fmt.Sprintf("%s", fmt.Sprint(aloneLine))
		if _, ok := cyMap[keyAloneLine]; !ok {

			cyMap[keyAloneLine] = sortIndex
			sortIndex++
			sortedCyPos = append(sortedCyPos, aloneLine)
			setting = append(setting, Blank{
				Head:           -1,
				Foot:           -1,
				HeadUseCyIndex: cyMap[keyAloneLine],
				FootUseCyIndex: -1,
			})
		}
	}

	//配置排序
	for index, info := range setting {
		if info.HeadUseCyIndex > info.FootUseCyIndex {
			setting[index].HeadUseCyIndex, setting[index].FootUseCyIndex = setting[index].FootUseCyIndex, setting[index].HeadUseCyIndex
			setting[index].Head, setting[index].Foot = setting[index].Foot, setting[index].Head
		}
	}

	sort.Sort(BlankSort(setting))

	//配置分组
	lastFootUseCyIndex := 0
	groupSetting := [][]Blank{}
	formattedSetting := []Blank{}
	for _, val := range setting {
		if lastFootUseCyIndex == 0 || lastFootUseCyIndex != val.FootUseCyIndex {
			groupSetting = append(groupSetting, []Blank{val})
		} else if val.FootUseCyIndex == lastFootUseCyIndex {
			length := len(groupSetting)
			groupSetting[length-1] = append(groupSetting[length-1], val)
		}
		lastFootUseCyIndex = val.FootUseCyIndex
	}

	//分组配置格式化
	for _, item := range groupSetting {
		if len(item) <= 0 {
			continue
		} else if len(item) == 1 {
			formattedSetting = append(formattedSetting, item[0])
		} else {
			temp := Blank{
				HeadFoot: make([]BlankItem, 0),
			}
			for _, val := range item {
				temp.HeadFoot = append(temp.HeadFoot, BlankItem{HeadUseCyIndex: val.HeadUseCyIndex, FootUseCyIndex: val.FootUseCyIndex, Head: val.Head, Foot: val.Foot})
			}
			formattedSetting = append(formattedSetting, temp)
		}
	}
	return formattedSetting, sortedCyPos, nil
}

func setSettingPos(allColCY, allLineCY []ChengYu, setting *[]Blank, sortIndex *int, cyMap map[string]int, sortedCyPos *[]ChengYu, crossPoint map[string]bool) {

	for _, col := range allColCY {
		keyCol := fmt.Sprintf("%s", fmt.Sprint(col))
		if _, ok := cyMap[keyCol]; !ok || cyMap[keyCol] <= 0 {
			*sortedCyPos = append(*sortedCyPos, col)
			cyMap[keyCol] = *sortIndex
			*sortIndex++
		}
		isAloneCol := true //是否是无相交的独立成语
		for _, line := range allLineCY {
			keyLine := fmt.Sprintf("%s", fmt.Sprint(line))
			var colPos, linePos int
			var point CyCell
			var err error
			if colPos, linePos, point, err = getHitPoint(col, line); err != nil {
				//fmt.Println("getHitPoint err: ", err, line, col)
			} else {
				isAloneCol = false
				if _, ok := cyMap[keyLine]; !ok || cyMap[keyLine] <= 0 {
					*sortedCyPos = append(*sortedCyPos, line)
					cyMap[keyLine] = *sortIndex
					*sortIndex++
				}
				if !crossPoint[fmt.Sprintf("%d,%d", point.X, point.Y)] {

					*setting = append(*setting, Blank{
						Head:           colPos,
						Foot:           linePos,
						HeadUseCyIndex: cyMap[keyCol],
						FootUseCyIndex: cyMap[keyLine],
					})
				}
				crossPoint[fmt.Sprintf("%d,%d", point.X, point.Y)] = true
				for _, innerCol := range allColCY {
					//剩余的与当前行有交点的列也记录
					keyInnerCol := fmt.Sprintf("%s", fmt.Sprint(innerCol))

					if colPos, linePos, point, err = getHitPoint(innerCol, line); err == nil {

						if _, ok := cyMap[keyInnerCol]; !ok {
							*sortedCyPos = append(*sortedCyPos, innerCol)
							cyMap[keyInnerCol] = *sortIndex
							*sortIndex++
						}
						if !crossPoint[fmt.Sprintf("%d,%d", point.X, point.Y)] {
							*setting = append(*setting, Blank{
								Head:           colPos,
								Foot:           linePos,
								HeadUseCyIndex: cyMap[keyInnerCol],
								FootUseCyIndex: cyMap[keyLine],
							})

							crossPoint[fmt.Sprintf("%d,%d", point.X, point.Y)] = true
							setSettingPos([]ChengYu{innerCol}, allLineCY, setting, sortIndex, cyMap, sortedCyPos, crossPoint)
						}

					}
				}
			}
		}
		if isAloneCol {
			//无相交的成语
			*setting = append(*setting, Blank{
				Head:           -1,
				Foot:           -1,
				HeadUseCyIndex: cyMap[keyCol],
				FootUseCyIndex: -1,
			})

		}
	}
}

// 获取成语交叉点位置
func getHitPoint(col, line ChengYu) (colPos, linePos int, point CyCell, err error) {
	if len(line) < CyLen || len(col) < CyLen {
		return 0, 0, CyCell{}, errors.New("invalid len")
	}
	for indexLine, pointX := range line {
		for indexCol, pointY := range col {
			if pointX.X == pointY.X && pointX.Y == pointY.Y {
				return indexCol + 1, indexLine + 1, pointX, nil
			}
		}
	}
	return 0, 0, CyCell{}, errors.New("no hit point")
}

type BlankSort []Blank

func (a BlankSort) Len() int      { return len(a) }
func (a BlankSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a BlankSort) Less(i, j int) bool {
	if a[i].FootUseCyIndex == a[j].FootUseCyIndex {
		return a[i].HeadUseCyIndex < a[j].HeadUseCyIndex
	}
	return a[i].FootUseCyIndex < a[j].FootUseCyIndex
}
