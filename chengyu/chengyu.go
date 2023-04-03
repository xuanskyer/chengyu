package chengyu

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"
)

const (
	TableMaxLen = 9 //表格边长

	P9rNil      = 0 //placeholder 占位符状态：未使用
	P9rNilStr   = "  "
	P9rBlank    = 1 //placeholder 占位符状态：留空处
	P9rBlankStr = "▒▒"
	P9rUsed     = 2 //placeholder 占位符状态：有字
	P9rUsedStr  = "☑"

	CyLen = 4 //一个成语的字数

	CyTypeCol  = 0
	CyTypeLine = 1

	MaxResultCount = 500
)

var (
	tableXyStatus    = make(map[string]int, 0) //坐标状态
	formattedSetting []Blank                   //表格格式化后的遍历配置列表
	sortedCyPos      = make([]ChengYu, 0)      //成语位置出现顺序列表
	sortedCyPosLen   = 0                       //成语位置出现顺序列表长度
	chengYuMap       = make(map[string]bool)   //成语库map
	allCY            []ChengYu                 //所有成语位置列表
	allCYLen         = 0                       //所有成语位置列表长度
	allLineCY        []ChengYu                 //所有行中的成语位置列表
	allColCY         []ChengYu                 //所有列中的成语位置列表
	maxResultCount   = MaxResultCount
)

type CyCell struct { //位置坐标
	X int `json:"x"`
	Y int `json:"y"`
}

type ChengYu []CyCell   //成语
type QueueItem struct { //队列元素
	Cy   ChengYu `json:"cy"`
	Type int     `json:"type"`
}

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

// Init 表格坐标转化成遍历所需配置
func Init(table [][]int, chengYuList []string, max ...int) error {
	isValidTable := IsValidTemplate(table)
	if !isValidTable {
		return errors.New("表格模板非法！")
	}
	//变量初始化
	tableXyStatus = make(map[string]int, 0) //坐标状态
	formattedSetting = make([]Blank, 0)     //表格格式化后的遍历配置列表
	sortedCyPos = make([]ChengYu, 0)        //成语位置出现顺序列表
	sortedCyPosLen = 0                      //成语位置出现顺序列表长度
	chengYuMap = make(map[string]bool)      //成语库map
	allCY = make([]ChengYu, 0)              //所有成语位置列表
	allCYLen = 0                            //所有成语位置列表长度
	allLineCY = make([]ChengYu, 0)          //所有行中的成语位置列表
	allColCY = make([]ChengYu, 0)           //所有列中的成语位置列表
	maxResultCount = MaxResultCount

	if len(max) > 0 && max[0] > 0 {
		maxResultCount = max[0]
	}
	var setting []Blank
	sortIndex := 0
	cyMap := make(map[string]int, 0)
	cyQueue := make([]QueueItem, 0)
	crossPoint := make(map[string]bool, 0)
	// 使用 map 存储成语列表，方便去重
	for _, item := range chengYuList {
		chengYuMap[item] = true
	}
	for index, item := range table {
		for k, v := range item {
			tableXyStatus[fmt.Sprintf("%d,%d", index, k)] = v
		}
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

	allCY = append(append(allCY, allColCY...), allLineCY...)
	allCYLen = len(allCY)
	loopSettingPos(allColCY, allLineCY, &cyQueue)

	cyOutQueue(allColCY, allLineCY, &setting, &sortIndex, cyMap, &sortedCyPos, crossPoint, &cyQueue)

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
	for _, aloneCol := range allColCY {
		keyAloneCol := fmt.Sprintf("%s", fmt.Sprint(aloneCol))
		if _, ok := cyMap[keyAloneCol]; !ok {

			cyMap[keyAloneCol] = sortIndex
			sortIndex++
			sortedCyPos = append(sortedCyPos, aloneCol)
			setting = append(setting, Blank{
				Head:           -1,
				Foot:           -1,
				HeadUseCyIndex: cyMap[keyAloneCol],
				FootUseCyIndex: -1,
			})
		}
	}
	sortedCyPosLen = len(sortedCyPos)
	for index, info := range setting {
		if info.HeadUseCyIndex > info.FootUseCyIndex {
			setting[index].HeadUseCyIndex, setting[index].FootUseCyIndex = setting[index].FootUseCyIndex, setting[index].HeadUseCyIndex
			setting[index].Head, setting[index].Foot = setting[index].Foot, setting[index].Head
		}
	}

	sort.Sort(BlankSort(setting))

	//配置分组
	lastFootUseCyIndex := 0
	var groupSetting [][]Blank
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
				temp.FootUseCyIndex = val.FootUseCyIndex
				temp.HeadFoot = append(temp.HeadFoot, BlankItem{HeadUseCyIndex: val.HeadUseCyIndex, FootUseCyIndex: val.FootUseCyIndex, Head: val.Head, Foot: val.Foot})
			}
			formattedSetting = append(formattedSetting, temp)
		}
	}
	return nil
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

func Check(ones []string) bool {
	lengthOnes := len(ones)
	if lengthOnes != sortedCyPosLen {
		return false
	}
	onesMap := make(map[string]bool, 0)
	for _, one := range ones {
		onesMap[one] = true
	}
	if len(onesMap) != lengthOnes {
		return false
	}

	var c1, c2 string
	var e1, e2 error
	for index, info := range formattedSetting {
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
			if info.Head < 0 || info.Foot < 0 {
				continue
			}
			if index-1 >= 0 {
				//预判断下一个成语是否和当前成语相交
				if formattedSetting[index-1].FootUseCyIndex+1 < info.FootUseCyIndex {
					continue
				}
			}
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

func RecursionGenerate(depth int, selectedOnes []string, result map[string]bool, selectedMap map[string]bool) {

	if len(result) > maxResultCount-1 {
		return
	}
	onesMap := make(map[string]bool, 0)
	for _, one := range selectedOnes {
		onesMap[one] = true
	}
	if len(onesMap) != len(selectedOnes) {
		return
	}
	if len(selectedOnes) == allCYLen {
		// 已填好所有空白处配置，判断所选成语序列是否符合条件
		key := strings.Join(selectedOnes, ",")
		_, ok := selectedMap[key]
		if !ok {
			if Check(selectedOnes) {
				result[strings.Join(selectedOnes, ",")] = true
			}
			key = strings.Join(selectedOnes, ",")
			selectedMap[key] = true
		}
		return
	}

	var cell1, cell2 string
	var err1, err2 error
	// 处理当前递归层
ChengYuMapFor:
	for c := range chengYuMap {
		if depth > 0 {
			if isExisted(c, selectedOnes) {
				continue
			}
			blank := formattedSetting[depth-1]
			if len(blank.HeadFoot) > 0 {
				hitCount := 0
				for _, val := range blank.HeadFoot {
					cell1, err1 = GetChengYuPosStr(val.Head-1, val.Head, selectedOnes[val.HeadUseCyIndex])
					cell2, err2 = GetChengYuPosStr(val.Foot-1, val.Foot, c)
					if err2 != nil || cell2 == "" || cell1 == "" || cell1 != cell2 || err1 != nil {
						continue ChengYuMapFor
					} else {
						hitCount++
					}
				}
				if hitCount == len(blank.HeadFoot) {
					goto HitAndRecursion
				}
			} else {
				if blank.Head < 0 || blank.Foot < 0 {
					//无相交的成语直接跳过判断
					goto HitAndRecursion
				}
				if depth-2 >= 0 {
					//预判断下一个成语是否和当前成语相交
					if formattedSetting[depth-2].FootUseCyIndex+1 < blank.FootUseCyIndex && blank.HeadUseCyIndex >= len(selectedOnes) {
						RecursionGenerate(depth, append(selectedOnes, c), result, selectedMap)
						continue
					}
				}
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
		RecursionGenerate(depth+1, append(selectedOnes, c), result, selectedMap)
	}
}

// IsValidTemplate 判断模板是否合法：同一行/列 不能有多个成语相连或者重叠
func IsValidTemplate(table [][]int) bool {
	for _, item := range table {
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

// TableXyWordSet 打印输出一个结果
func TableXyWordSet(cyStr string) ([TableMaxLen][TableMaxLen]string, error) {

	tableXyWord := [TableMaxLen][TableMaxLen]string{} //坐标字
	if cyStr == "" {
		return tableXyWord, errors.New("empty cy str")
	}
	cyList := strings.Split(cyStr, ",")
	if len(cyList) <= 0 {
		return tableXyWord, errors.New("empty cy str")
	}
	if len(cyList) != sortedCyPosLen {
		return tableXyWord, errors.New("invalid cy str")
	}

	for k, v := range tableXyWord {
		for kk, _ := range v {
			if tableXyStatus[fmt.Sprintf("%d,%d", k, kk)] == P9rNil {
				tableXyWord[k][kk] = P9rNilStr
			} else if tableXyStatus[fmt.Sprintf("%d,%d", k, kk)] == P9rBlank {
				tableXyWord[k][kk] = P9rBlankStr
			} else {
				tableXyWord[k][kk] = P9rUsedStr
			}
		}
	}
	for index, cy := range cyList {
		point := sortedCyPos[index]
		if len(point) < CyLen || len(cy) < CyLen {
			continue
		}
		if tableXyWord[point[0].Y][point[0].X] == P9rUsedStr {
			word1, _ := GetChengYuPosStr(0, 1, cy)
			tableXyWord[point[0].Y][point[0].X] = word1
		}
		if tableXyWord[point[1].Y][point[1].X] == P9rUsedStr {
			word2, _ := GetChengYuPosStr(1, 2, cy)
			tableXyWord[point[1].Y][point[1].X] = word2
		}
		if tableXyWord[point[2].Y][point[2].X] == P9rUsedStr {
			word3, _ := GetChengYuPosStr(2, 3, cy)
			tableXyWord[point[2].Y][point[2].X] = word3
		}
		if tableXyWord[point[3].Y][point[3].X] == P9rUsedStr {
			word4, _ := GetChengYuPosStr(3, 4, cy)
			tableXyWord[point[3].Y][point[3].X] = word4
		}
	}
	return tableXyWord, nil
}

// GetSliceXN 竖向取二维切片的第N列
func GetSliceXN(table [][]int, col int) ([]int, error) {
	var column []int
	for i := 0; i < len(table); i++ {
		column = append(column, table[i][col])
	}
	return column, nil
}

// GetChengYu 从第n 行/列 取出成语（从0开始）
func GetChengYu(n int, slice []int, fixLine bool) []ChengYu {
	var cyList []ChengYu
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

// QueueIn 入队
func QueueIn(cy QueueItem, cyQueue []QueueItem) []QueueItem {
	cyQueue = append([]QueueItem{cy}, cyQueue...)
	return cyQueue
}

// QueueOut 出队
func QueueOut(cyQueue []QueueItem) ([]QueueItem, QueueItem, error) {
	l := len(cyQueue)
	if l <= 0 {
		return cyQueue, QueueItem{}, errors.New("empty queue")
	}
	out := cyQueue[l-1]
	return cyQueue[0 : l-1], out, nil
}

// 遍历成语：广度遍历 + 深度遍历
func loopSettingPos(allColCY, allLineCY []ChengYu, cyQueue *[]QueueItem) {

	for _, col := range allColCY {
		queueItem := QueueItem{
			Cy:   col,
			Type: CyTypeCol,
		}
		*cyQueue = QueueIn(queueItem, *cyQueue)
		cyInQueue(allColCY, allLineCY, queueItem, CyTypeCol, cyQueue)
		for _, line := range allLineCY {
			lineQueueItem := QueueItem{
				Cy:   line,
				Type: CyTypeLine,
			}
			*cyQueue = QueueIn(lineQueueItem, *cyQueue)
			cyInQueue(allColCY, allLineCY, lineQueueItem, CyTypeLine, cyQueue)
		}
	}
}

func cyInQueue(allColCY, allLineCY []ChengYu, queueItem QueueItem, yuType int, cyQueue *[]QueueItem) {
	if yuType == CyTypeCol {
		for _, item := range allLineCY {
			if _, _, _, err := getHitPoint(queueItem.Cy, item); err == nil {
				*cyQueue = QueueIn(QueueItem{Cy: item, Type: CyTypeLine}, *cyQueue)
			}
		}
	} else {
		for _, item := range allColCY {
			if _, _, _, err := getHitPoint(item, queueItem.Cy); err == nil {
				*cyQueue = QueueIn(QueueItem{Cy: item, Type: CyTypeCol}, *cyQueue)
			}
		}
	}
}

func cyOutQueue(allColCY, allLineCY []ChengYu, setting *[]Blank, sortIndex *int, cyMap map[string]int, sortedCyPos *[]ChengYu, crossPoint map[string]bool, cyQueue *[]QueueItem) {

	var err error
	var out QueueItem
	for {
		if len(*cyQueue) <= 0 {
			break
		}
		*cyQueue, out, err = QueueOut(*cyQueue)
		if err != nil {
			fmt.Println("cyOutQueue err:", err)
		}
		if out.Type == CyTypeCol {
			keyCol := fmt.Sprintf("%s", fmt.Sprint(out.Cy))
			for _, item := range allLineCY {
				var cell CyCell
				var colPos, linePos int
				if colPos, linePos, cell, err = getHitPoint(out.Cy, item); err == nil {
					keyLine := fmt.Sprintf("%s", fmt.Sprint(item))
					if _, ok := cyMap[keyCol]; !ok {
						*sortedCyPos = append(*sortedCyPos, out.Cy)
						cyMap[keyCol] = *sortIndex
						*sortIndex++
					}
					if _, ok := cyMap[keyLine]; !ok {
						*sortedCyPos = append(*sortedCyPos, item)
						cyMap[keyLine] = *sortIndex
						*sortIndex++
					}
					pointKey := fmt.Sprintf("%d,%d", cell.X, cell.Y)
					if !crossPoint[pointKey] {
						*setting = append(*setting, Blank{
							Head:           colPos,
							Foot:           linePos,
							HeadUseCyIndex: cyMap[keyCol],
							FootUseCyIndex: cyMap[keyLine],
						})
					}
					crossPoint[pointKey] = true
				}
			}
		} else {
			keyLine := fmt.Sprintf("%s", fmt.Sprint(out.Cy))
			for _, item := range allColCY {
				var cell CyCell
				var colPos, linePos int
				if colPos, linePos, cell, err = getHitPoint(item, out.Cy); err == nil {
					keyCol := fmt.Sprintf("%s", fmt.Sprint(item))
					if _, ok := cyMap[keyCol]; !ok {
						*sortedCyPos = append(*sortedCyPos, item)
						cyMap[keyCol] = *sortIndex
						*sortIndex++
					}
					if _, ok := cyMap[keyLine]; !ok {
						*sortedCyPos = append(*sortedCyPos, out.Cy)
						cyMap[keyLine] = *sortIndex
						*sortIndex++
					}
					pointKey := fmt.Sprintf("%d,%d", cell.X, cell.Y)
					if !crossPoint[pointKey] {
						*setting = append(*setting, Blank{
							Head:           linePos,
							Foot:           colPos,
							HeadUseCyIndex: cyMap[keyLine],
							FootUseCyIndex: cyMap[keyCol],
						})
					}
					crossPoint[pointKey] = true
				}
			}
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

func isExisted(cy string, cyList []string) bool {
	if len(cyList) <= 0 {
		return false
	}
	for _, item := range cyList {
		if cy == item {
			return true
		}
	}
	return false
}
