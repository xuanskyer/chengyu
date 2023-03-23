package main

import (
	"errors"
	"fmt"
	"github.com/chengyu/chengyu"
	"os"
	"time"
)

const (
	MaxX = 9
	MaxY = 9

	p9rNil   = 0 //占位符状态：未使用
	p9rBlank = 1 //占位符状态：空白
	p9rUsed  = 2 //占位符状态：有字

	ChengYuLen = 4
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
		"鞭长莫及", "不要回答", "长命百岁", "同舟共济", "争风吃醋", "争相吃瓜", "省吃俭用", "大吃一惊", "寅吃卯粮", "百家争鸣",
		"六神无主", "阿迪达斯", "好吃懒做", "千言万语", "只争朝夕", "草长莺飞", "挟长挟贵", "大富大贵", "骑马找马", "先入为主",
		"古今中外", "水秀山明", "你追我赶", "先来后到", "瓜田李下", "回天无力", "桃红柳绿", "前思后想", "落地生根", "大吃大喝",
		"千军万马", "取长补短", "一言为定", "一事无成", "万水千山", "明明白白", "张灯结彩", "一心一意", "一团和气", "美中不足",
		"深情厚谊", "马到成功", "张三李四", "千人一面", "再三再四", "成家立业", "一五一十", "狗急跳墙", "安身立命", "十指连心",
		"万里长城", "心直口快", "安居乐业", "长年累月", "窗明几净", "冷言冷语", "天罗地网", "天网恢恢", "人来人往", "目中无人",
		"三心二意", "莺歌燕舞", "热火朝天", "落花流水", "毛手毛脚", "旁若无人", "答非所问", "众多非一", "十字路口", "星星点点",
		"红男绿女", "百无是处", "快言快语", "别有洞天", "水深火热", "春风化雨", "一路平安", "自由自在", "棋逢对手", "叶公好龙",
		"后会无期", "守株待兔", "凤凰于飞", "一生一世", "花好月圆", "世外桃源", "韬光养晦", "画蛇添足", "青梅竹马", "风花雪月",
		"滥竽充数", "没完没了", "总而言之", "欣欣向荣", "时光荏苒", "差强人意", "好好先生", "无懈可击", "随波逐流", "袖手旁观",
		"行尸走肉", "金蝉脱壳", "百里挑一", "金玉满堂", "背水一战", "霸王别姬", "天上人间", "不吐不快", "海阔天空", "情非得已", "满腹经纶", "兵临城下", "春暖花开", "插翅难逃", "黄道吉日", "天下无双", "偷天换日", "两小无猜", "卧虎藏龙", "珠光宝气", "簪缨世族", "花花公子", "绘声绘影", "国色天香", "相亲相爱", "八仙过海", "金玉良缘", "掌上明珠", "皆大欢喜", "逍遥法外", "生财有道", "极乐世界", "情不自禁", "愚公移山", "魑魅魍魉", "龙生九子", "精卫填海", "海市蜃楼", "高山流水", "卧薪尝胆", "壮志凌云", "金枝玉叶", "四海一家", "穿针引线", "无忧无虑", "无地自容", "三位一体", "落叶归根", "相见恨晚", "惊天动地", "滔滔不绝", "相濡以沫", "长生不死", "原来如此", "女娲补天", "三皇五帝", "万箭穿心", "水木清华", "窈窕淑女", "破釜沉舟", "天涯海角", "牛郎织女", "倾国倾城", "飘飘欲仙", "福星高照", "妄自菲薄", "永无止境", "学富五车", "饮食男女", "英雄豪杰", "国士无双", "塞翁失马", "万家灯火", "石破天惊", "精忠报国", "养生之道", "覆雨翻云", "六道轮回", "鹰击长空", "日日夜夜", "厚德载物", "亡羊补牢", "万里长城", "黄金时代", "出生入死", "一路顺风", "随遇而安", "千军万马", "郑人买履", "棋逢对手", "叶公好龙", "后会无期", "守株待兔", "凤凰于飞", "一生一世", "花好月圆", "世外桃源", "韬光养晦", "画蛇添足", "青梅竹马", "风花雪月", "滥竽充数", "没完没了", "总而言之", "欣欣向荣", "时光荏苒", "差强人意", "好好先生", "无懈可击", "随波逐流", "袖手旁观", "群雄逐鹿", "血战到底", "唯我独尊", "买椟还珠", "龙马精神", "一见钟情", "喜闻乐见", "负荆请罪", "三人成虎", "河东狮吼", "程门立雪", "金戈铁马", "笑逐颜开", "千钧一发", "纸上谈兵", "风和日丽", "邯郸学步", "大器晚成", "庖丁解牛", "甜言蜜语", "雷霆万钧", "浮生若梦", "大开眼界", "汗牛充栋", "百鸟朝凤", "以德服人", "白驹过隙", "难兄难弟", "鬼哭神嚎", "声色犬马", "指鹿为马", "龙争虎斗", "雾里看花", "男大当婚", "未雨绸缪", "南辕北辙", "三从四德", "一丝不挂", "高屋建瓴", "阳春白雪", "杯弓蛇影", "闻鸡起舞", "四面楚歌", "登堂入室", "张灯结彩", "而立之年", "饮鸩止渴", "杏雨梨云", "龙凤呈祥", "勇往直前", "左道旁门", "莫衷一是", "马踏飞燕", "掩耳盗铃", "大江东去", "凿壁偷光", "色厉内荏", "花容月貌", "越俎代庖", "鳞次栉比", "美轮美奂", "缘木求鱼", "再接再厉", "马到成功", "红颜知己", "赤子之心", "迫在眉睫", "风流韵事", "相形见绌", "诸子百家", "鬼迷心窍", "星火燎原", "画地为牢", "岁寒三友", "花花世界", "纸醉金迷", "狐假虎威", "纵横捭阖", "沧海桑田", "不求甚解", "暴殄天物", "吃喝玩乐", "乐不思蜀", "身不由己", "小家碧玉", "文不加点", "天马行空", "人来人往", "千方百计", "天高地厚", "万人空巷", "争分夺秒", "如火如荼", "大智若愚", "斗转星移", "七情六欲", "大禹治水", "空穴来风", "孟母三迁", "绘声绘色", "九五之尊", "随心所欲", "干将莫邪", "相得益彰", "借刀杀人", "浪迹天涯", "刚愎自用", "镜花水月", "黔驴技穷", "肝胆相照", "多多益善", "叱咤风云", "杞人忧天", "作茧自缚", "一飞冲天", "殊途同归", "风卷残云", "因果报应", "无可厚非", "赶尽杀绝", "天长地久", "飞龙在天", "桃之夭夭", "南柯一梦", "口是心非", "江山如画", "风华正茂", "一帆风顺", "一叶知秋", "草船借箭", "铁石心肠", "望其项背", "头晕目眩", "大浪淘沙", "纵横天下", "有问必答", "无为而治", "釜底抽薪", "吹毛求疵", "好事多磨", "空谷幽兰", "悬梁刺股", "白手起家", "完璧归赵", "忍俊不禁", "沐猴而冠", "白云苍狗", "贼眉鼠眼", "围魏救赵", "烟雨蒙蒙", "炙手可热", "尸位素餐", "出水芙蓉", "礼仪之邦", "一丘之貉", "鹏程万里", "叹为观止", "韦编三绝", "今生今世", "草木皆兵", "宁缺毋滥", "回光返照", "露水夫妻", "讳莫如深", "贻笑大方", "紫气东来", "万马奔腾", "一诺千金", "老马识途", "五花大绑", "捉襟见肘", "瓜田李下", "水漫金山", "苦心孤诣", "可见一斑", "五湖四海", "虚怀若谷", "欲擒故纵", "风声鹤唳", "毛遂自荐", "蛛丝马迹", "中庸之道", "迷途知返", "自由自在", "龙飞凤舞", "树大根深", "雨过天晴", "乘风破浪", "筚路蓝缕", "朝三暮四", "患得患失", "君子好逑", "鞭长莫及", "竭泽而渔", "飞黄腾达", "囊萤映雪", "飞蛾扑火", "自怨自艾", "风驰电掣", "白马非马", "退避三舍", "三山五岳", "称心如意", "望梅止渴", "茕茕孑立", "振聋发聩", "运筹帷幄", "逃之夭夭", "杯水车薪", "有的放矢", "矫枉过正", "睚眦必报", "姗姗来迟", "一鸣惊人", "孜孜不倦", "一马平川", "入木三分", "沆瀣一气", "天伦之乐", "兄弟阋墙", "藕断丝连", "心猿意马", "想入非非", "盲人摸象", "眉飞色舞", "三教九流", "高楼大厦", "锲而不舍", "过犹不及", "狗尾续貂", "斗酒学士", "高山仰止", "形影不离", "小心翼翼", "返璞归真", "见贤思齐", "按图索骥", "枪林弹雨", "桀骜不驯", "遇人不淑", "道貌岸然", "名扬四海", "虚与委蛇", "门可罗雀", "水落石出", "不卑不亢", "无法无天", "拔苗助长", "大快朵颐", "因地制宜", "单刀直入", "时来运转", "天方夜谭", "一蹴而就", "踌躇满志", "战无不胜", "插翅难飞", "图穷匕见", "鬼话连篇", "亢龙有悔", "望洋兴叹", "爱屋及乌", "惊鸿一瞥", "风华绝代", "名胜古迹", "如履薄冰", "持之以恒", "潜移默化", "昙花一现", "巫山云雨", "狡兔三窟", "栉风沐雨", "骇人听闻", "断章取义", "曲突徙薪", "谢天谢地", "脱颖而出", "垂帘听政", "一马当先", "不耻下问", "不以为然", "春华秋实", "欲盖弥彰", "人琴俱亡", "投鼠忌器", "歧路亡羊", "金风玉露", "落花流水", "春风化雨", "心如刀割", "锱铢必较", "一叶障目", "来历不明", "名副其实", "中流砥柱", "绕梁三日", "安步当车", "放荡不羁", "天衣无缝", "自相矛盾", "神机妙算", "沧海一粟", "冲锋陷阵", "龙虎风云", "言简意赅", "九死一生", "铁树开花", "画龙点睛", "风雨无阻", "坐井观天", "奇货可居", "浮光掠影", "牝鸡司晨", "沽名钓誉", "天作之合", "甚嚣尘上", "铩羽而归", "劫后余生", "泾渭分明", "节哀顺变", "有恃无恐", "不绝如缕", "马革裹尸", "监守自盗", "耳濡目染", "金屋藏娇", "不约而同", "逐鹿中原", "龙潭虎穴", "江郎才尽", "明日黄花", "栩栩如生", "人山人海", "面面相觑", "唇亡齿寒", "相敬如宾", "知法犯法", "欢聚一堂", "曾几何时", "纷至沓来", "李代桃僵", "毛骨悚然", "衣冠禽兽", "有凤来仪", "见微知著", "旗鼓相当", "无与伦比", "摸金校尉", "牛头马面", "凤毛麟角", "难得糊涂", "衣香鬓影", "马到功成", "鸠占鹊巢", "狭路相逢", "春秋笔法", "厉兵秣马", "约法三章", "豁然开朗", "平步青云", "步步为营", "蝇营狗苟", "心如止水", "从善如流", "殚精竭虑", "十字路口", "矢志不渝", "九九归一", "井底之蛙", "居安思危", "不一而足", "周而复始", "望穿秋水", "秦晋之好", "不落窠臼", "司空见惯", "怙恶不悛", "百年好合", "出神入化", "身体力行", "敬谢不敏", "嗤之以鼻", "天之骄子", "贤妻良母", "能说会道", "进退维谷", "甘之如饴", "人心不古", "颐指气使", "墨守成规", "左右逢源", "回心转意", "插科打诨", "别来无恙", "翩翩公子", "穷兵黩武", "舌战群儒", "字字珠玑", "义无反顾", "举重若轻", "钟灵毓秀", "水滴石穿", "防微杜渐", "衣冠楚楚", "卧冰求鲤", "觥筹交错", "络绎不绝", "自强不息", "秀色可餐", "至理名言", "分庭抗礼", "萍水相逢", "水性杨花", "戛然而止", "气喘吁吁", "沉鱼落雁", "望尘莫及", "亦步亦趋", "川流不息", "千锤百炼", "谈笑风生", "高朋满座", "丧心病狂", "天下无敌", "惊弓之鸟", "耿耿于怀", "心照不宣", "荦荦大端", "噤若寒蝉", "上下其手", "弄假成真", "天网恢恢", "夜郎自大", "鞭辟入里", "义薄云天", "所向披靡", "点石成金", "回眸一笑", "巴山夜雨", "兢兢业业", "克己复礼", "风起云涌", "不惑之年", "义愤填膺", "门当户对", "声名狼藉", "时过境迁", "念念不忘", "鞠躬尽瘁", "不言而喻", "人生如梦", "琴棋书画", "酸甜苦辣", "走马观花", "全力以赴", "人面桃花", "王侯将相", "青山不老", "朝令夕改", "小时了了", "玩世不恭", "人情世故", "聊胜于无", "为虎作伥", "休戚相关", "三阳开泰", "五子登科", "熙熙攘攘", "开源节流", "绝处逢生", "一石二鸟", "鬼斧神工", "青天白日", "病入膏肓", "横行霸道", "对牛弹琴", "诚惶诚恐", "胡服骑射", "虎视眈眈", "十万火急", "断袖之癖", "得陇望蜀", "分道扬镳", "壮士断腕", "自惭形秽", "云淡风轻", "巾帼英雄", "眼花缭乱", "不可一世", "沁人心脾", "侃侃而谈", "闻过则喜", "班门弄斧", "舍我其谁", "潸然泪下", "肆无忌惮", "心旷神怡", "物竞天择", "东山再起", "丹凤朝阳", "和光同尘", "心力衰竭", "事半功倍", "阿鼻地狱", "九关虎豹", "劝百讽一", "琳琅满目", "一丝不苟", "逝者如斯", "同仇敌忾", "朝秦暮楚", "不亦乐乎", "哭笑不得", "重见天日", "集腋成裘", "风月无边", "乐此不疲", "咫尺天涯", "宠辱不惊", "安然无恙", "一事无成", "若即若离", "本末倒置", "秋风落叶", "无价之宝", "金刚怒目", "以儆效尤", "波涛汹涌", "花团锦簇", "海枯石烂", "目无全牛", "颠倒乾坤", "当仁不让", "车水马龙", "天下为公", "火中取栗", "众矢之的", "尽善尽美", "欢天喜地", "今非昔比", "天府之国", "不可名状", "异想天开", "粉墨登场", "根深蒂固", "钟鸣鼎食", "历历在目", "不法之徒", "出人头地", "以德报怨", "梨花带雨", "抛砖引玉", "优柔寡断", "开门见山", "参差不齐", "温文尔雅", "暗度陈仓", "甘心情愿", "挑肥拣瘦", "阿猫阿狗", "心有余悸", "数典忘祖", "喜出望外", "文过饰非", "连锁反应", "将心比心", "无动于衷", "鹤唳华亭", "妙手空空", "登峰造极", "惊涛骇浪", "自欺欺人", "绿树成荫", "岂有此理", "万马齐喑", "世态炎凉", "冠冕堂皇", "天罗地网", "踽踽独行", "兔死狐悲", "众志成城", "耳提面命", "待字闺中", "女扮男装", "东张西望", "马首是瞻", "物极必反", "蔚然成风", "迫不及待", "淋漓尽致", "风尘仆仆", "外强中干", "求全责备", "人浮于事", "安居乐业", "珠联璧合", "一网打尽", "任重道远", "循循善诱", "移花接木", "不知所措", "柳暗花明", "白虹贯日", "首鼠两端", "前仆后继", "醉生梦死", "惺惺相惜", "焚膏继晷", "金童玉女", "横扫千军", "闭门造车", "峰回路转", "涸辙之鲋", "锦上添花", "亭亭玉立", "干柴烈火", "香草美人", "新亭对泣", "鹤立鸡群", "一往无前", "吴下阿蒙", "草长莺飞", "兔死狗烹", "姹紫嫣红", "因材施教", "长生不老", "爱莫能助", "洗耳恭听", "信手拈来", "时不我待", "举一反三", "蠢蠢欲动", "苟延残喘", "正襟危坐", "助人为乐", "火树银花", "齐大非偶", "无影无踪", "不胫而走", "笨鸟先飞", "精打细算", "尾大不掉", "词不达意", "门庭若市", "落英缤纷", "戎马倥偬", "上行下效", "提纲挈领", "蹉跎岁月",
	}
	fmt.Println("表格模板：")
	//demo1
	//table := [][]int{
	//	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	{0, 1, 0, 1, 0, 0, 0, 0, 0},
	//	{0, 1, 1, 1, 1, 0, 0, 0, 0},
	//	{0, 1, 1, 1, 1, 0, 0, 0, 0},
	//	{0, 1, 1, 1, 1, 0, 0, 0, 0},
	//	{0, 0, 1, 0, 1, 0, 0, 0, 0},
	//	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	//}
	//
	//v2Setting := []chengyu.Blank{
	//	{HeadUseCyIndex: 0, Head: 2, FootUseCyIndex: 1, Foot: 1},
	//	{HeadUseCyIndex: 0, Head: 3, FootUseCyIndex: 2, Foot: 1},
	//	{HeadUseCyIndex: 0, Head: 4, FootUseCyIndex: 3, Foot: 1},
	//	{HeadFoot: []chengyu.BlankItem{
	//		chengyu.BlankItem{HeadUseCyIndex: 1, FootUseCyIndex: 4, Head: 2, Foot: 1},
	//		chengyu.BlankItem{HeadUseCyIndex: 2, FootUseCyIndex: 4, Head: 2, Foot: 2},
	//		chengyu.BlankItem{HeadUseCyIndex: 3, FootUseCyIndex: 4, Head: 2, Foot: 3},
	//	}},
	//	{HeadFoot: []chengyu.BlankItem{
	//		chengyu.BlankItem{HeadUseCyIndex: 1, FootUseCyIndex: 5, Head: 3, Foot: 2},
	//		chengyu.BlankItem{HeadUseCyIndex: 2, FootUseCyIndex: 5, Head: 3, Foot: 3},
	//		chengyu.BlankItem{HeadUseCyIndex: 3, FootUseCyIndex: 5, Head: 3, Foot: 4},
	//	}},
	//	{HeadFoot: []chengyu.BlankItem{
	//		chengyu.BlankItem{HeadUseCyIndex: 1, FootUseCyIndex: 6, Head: 4, Foot: 1},
	//		chengyu.BlankItem{HeadUseCyIndex: 2, FootUseCyIndex: 6, Head: 4, Foot: 2},
	//		chengyu.BlankItem{HeadUseCyIndex: 3, FootUseCyIndex: 6, Head: 4, Foot: 3},
	//	}},
	//}

	//demo2
	table := [][]int{
		{0, 0, 1, 0, 0, 0, 0, 0, 0},
		{0, 0, 1, 1, 1, 1, 0, 0, 0},
		{0, 0, 1, 0, 1, 0, 1, 0, 0},
		{0, 0, 1, 0, 1, 1, 1, 1, 0},
		{0, 0, 0, 0, 1, 0, 1, 0, 0},
		{0, 0, 0, 0, 0, 0, 1, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	//v2Setting := []chengyu.Blank{
	//	{HeadUseCyIndex: 0, FootUseCyIndex: 1, Head: 2, Foot: 1},
	//	{HeadUseCyIndex: 1, FootUseCyIndex: 2, Head: 3, Foot: 1},
	//	{HeadUseCyIndex: 2, FootUseCyIndex: 3, Head: 3, Foot: 1},
	//	{HeadUseCyIndex: 3, FootUseCyIndex: 4, Head: 3, Foot: 2},
	//}
	v2Setting, sortedCyPos, _ := table2Setting(table)
	//for _, item := range v2Setting {
	//	fmt.Printf("%+v\n", item)
	//}

	//for _, val := range sortedCyPos {
	//	fmt.Printf("%+v\n", val)
	//}
	isValidTable := IsValidTemplate(table)
	if !isValidTable {
		fmt.Println("表格模板非法！")
		return
	}
	allCY := []ChengYu{}
	allLineCY := []ChengYu{}
	allColCY := []ChengYu{}
	for index, item := range table {
		cy := getChengYu(index, item, false)
		if len(cy) > 0 {
			allLineCY = append(allLineCY, cy...)
		}
	}
	for i := 0; i < MaxY; i++ {
		column, _ := getSliceXN(table, i)
		cyCol := getChengYu(i, column, true)
		if len(cyCol) > 0 {
			allColCY = append(allColCY, cyCol...)
		}
	}
	//fmt.Println("所有行中的成语： ")
	//for _, item := range allLineCY {
	//	fmt.Println(item)
	//}
	//fmt.Println("所有列中的成语： ")
	//for _, item := range allColCY {
	//	fmt.Println(item)
	//}
	allCY = append(append(allCY, allColCY...), allLineCY...)

	//fmt.Println("所有成语位： ")
	//for index, item := range allCY {
	//	fmt.Println(index, item)
	//}
	//fmt.Printf("成语列表(总数：%d)\n", len(ChengYuList))

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
	fmt.Println("结果示例：")
	for index, one := range result {
		if index < 2 {
			printResult2Table(one, sortedCyPos)
			fmt.Println("\n")
		}
	}
	fmt.Printf("该函数执行完成耗时：%v，答案数：%d\n", elapsed, len(result))
}

func printResult2Table(one []string, sortedCyPos []ChengYu) {
	if len(one) <= 0 {
		return
	}
	if len(one) != len(sortedCyPos) {
		return
	}

	tableString := [9][9]string{
		{"  ", "  ", "  ", "  ", "  ", "  ", "  ", "  ", "  "},
		{"  ", "  ", "  ", "  ", "  ", "  ", "  ", "  ", "  "},
		{"  ", "  ", "  ", "  ", "  ", "  ", "  ", "  ", "  "},
		{"  ", "  ", "  ", "  ", "  ", "  ", "  ", "  ", "  "},
		{"  ", "  ", "  ", "  ", "  ", "  ", "  ", "  ", "  "},
		{"  ", "  ", "  ", "  ", "  ", "  ", "  ", "  ", "  "},
		{"  ", "  ", "  ", "  ", "  ", "  ", "  ", "  ", "  "},
		{"  ", "  ", "  ", "  ", "  ", "  ", "  ", "  ", "  "},
		{"  ", "  ", "  ", "  ", "  ", "  ", "  ", "  ", "  "},
	}
	for index, cy := range one {
		point := sortedCyPos[index]
		if len(point) < ChengYuLen || len(cy) < ChengYuLen {
			continue
		}
		word1, _ := chengyu.GetChengyuPosStr(0, 1, cy)
		tableString[point[0].Y][point[0].X] = word1
		word2, _ := chengyu.GetChengyuPosStr(1, 2, cy)
		tableString[point[1].Y][point[1].X] = word2
		word3, _ := chengyu.GetChengyuPosStr(2, 3, cy)
		tableString[point[2].Y][point[2].X] = word3
		word4, _ := chengyu.GetChengyuPosStr(3, 4, cy)
		tableString[point[3].Y][point[3].X] = word4
	}
	//fmt.Printf("%+v\n", one)
	for _, item := range tableString {
		fmt.Printf("%+v\n", item)
	}
}

// 判断模板是否合法
func IsValidTemplate(table [][]int) bool {
	for _, item := range table {
		fmt.Println(item)
		count := 0
		for _, val := range item {
			if val == p9rNil {
				count = 0
			} else {
				count++
			}
		}
		if count > ChengYuLen {
			return false
		}
	}
	for i := 0; i < MaxY; i++ {
		column, _ := getSliceXN(table, i)
		count := 0
		for _, val := range column {
			if val == p9rNil {
				count = 0
			} else {
				count++
			}
		}
		if count > ChengYuLen {
			return false
		}
	}
	return true
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

func table2Setting(table [][]int) ([]chengyu.Blank, []ChengYu, error) {
	setting := []chengyu.Blank{}
	allLineCY := []ChengYu{}
	allColCY := []ChengYu{}
	for index, item := range table {
		cy := getChengYu(index, item, false)
		if len(cy) > 0 {
			allLineCY = append(allLineCY, cy...)
		}
	}
	for i := 0; i < MaxY; i++ {
		column, _ := getSliceXN(table, i)
		cyCol := getChengYu(i, column, true)
		if len(cyCol) > 0 {
			allColCY = append(allColCY, cyCol...)
		}
	}
	//fmt.Println("所有行中的成语位： ")
	//for _, item := range allLineCY {
	//	fmt.Println(item)
	//	//key := fmt.Sprint(item)
	//	//fmt.Println("key: ", key)
	//}
	//fmt.Println("所有列中的成语位： ")
	//for _, item := range allColCY {
	//	fmt.Println(item)
	//}

	//TODO 判断一个成语到底是行的还是列的
	cyMap := make(map[string]int, 0)
	sortedCyPos := make([]ChengYu, 0)
	for _, col := range allColCY {
		keyCol := fmt.Sprintf("%s", fmt.Sprint(col))
		if _, ok := cyMap[keyCol]; !ok || cyMap[keyCol] <= 0 {
			cyMap[keyCol] = len(cyMap)
			sortedCyPos = append(sortedCyPos, col)
		}
		for _, line := range allLineCY {
			keyLine := fmt.Sprintf("%s", fmt.Sprint(line))
			var colPos, linePos int
			var err error
			if colPos, linePos, err = getHitPoint(col, line); err != nil {
				//fmt.Println("getHitPoint err: ", err, line, col)
			} else {
				if _, ok := cyMap[keyLine]; !ok || cyMap[keyLine] <= 0 {
					cyMap[keyLine] = len(cyMap)
					sortedCyPos = append(sortedCyPos, line)
				}
				setting = append(setting, chengyu.Blank{
					Head:           colPos,
					Foot:           linePos,
					HeadUseCyIndex: cyMap[keyCol],
					FootUseCyIndex: cyMap[keyLine],
				})
			}
		}
	}

	//配置排序
	for index, info := range setting {
		if info.HeadUseCyIndex > info.FootUseCyIndex {
			setting[index].HeadUseCyIndex, setting[index].FootUseCyIndex = setting[index].FootUseCyIndex, setting[index].HeadUseCyIndex
			setting[index].Head, setting[index].Foot = setting[index].Foot, setting[index].Head
		}
	}

	//配置分组
	lastFootUseCyIndex := 0
	groupSetting := [][]chengyu.Blank{}
	formattedSetting := []chengyu.Blank{}
	for _, val := range setting {
		if lastFootUseCyIndex == 0 || lastFootUseCyIndex != val.FootUseCyIndex {
			groupSetting = append(groupSetting, []chengyu.Blank{val})
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
			temp := chengyu.Blank{
				HeadFoot: make([]chengyu.BlankItem, 0),
			}
			for _, val := range item {
				temp.HeadFoot = append(temp.HeadFoot, chengyu.BlankItem{HeadUseCyIndex: val.HeadUseCyIndex, FootUseCyIndex: val.FootUseCyIndex, Head: val.Head, Foot: val.Foot})
			}
			formattedSetting = append(formattedSetting, temp)
		}
	}
	return formattedSetting, sortedCyPos, nil
}

// 获取成语交叉点位置
func getHitPoint(col, line ChengYu) (colPos, linePos int, err error) {
	if len(line) < ChengYuLen || len(col) < ChengYuLen {
		return 0, 0, errors.New("invalid len")
	}
	for indexLine, pointX := range line {
		for indexCol, pointY := range col {
			if pointX.X == pointY.X && pointX.Y == pointY.Y {
				return indexCol + 1, indexLine + 1, nil
			}
		}
	}
	return 0, 0, errors.New("no hit point")
}
