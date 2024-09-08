package bazicore

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"math"
	"os"
	"path"
	"strings"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type BaziMeta struct {
	Longitude    float64 //经度
	Name         string
	Gender       bool
	BzAst        bool
	Bzdt         *BZDateTime
	HeadJieQi    *BZDateTime
	TailJieQi    *BZDateTime
	HeadJieQiAst *BZDateTime
	TailJieQiAst *BZDateTime
	Jydt         *BZDateTime
	Astdt        *BZDateTime
	Lifa         int
	Bazi         [8]int
	ShiShen      []int
	QiYunDelta   [3]int
}

func (bzm *BaziMeta) New(dt *BZDateTime, gender bool, ast bool, lon float64, lifa int, name string) *BaziMeta {
	if len(name) == 0 {
		name = "无名氏"
	}
	astDt := &BZDateTime{}
	if ast {
		astDt = calcAST(dt, lon)
	} else {
		astDt = dt
	}

	initDt := &BZDateTime{}
	initDt = initDt.New(BzDataTimeInitTime)

	bz := &BaziMeta{
		Longitude:    lon,
		Name:         name,
		Gender:       gender, //男：false, 女：true
		BzAst:        ast,
		Bzdt:         dt,
		HeadJieQi:    initDt,
		TailJieQi:    initDt,
		HeadJieQiAst: initDt,
		TailJieQiAst: initDt,
		Jydt:         initDt,
		Astdt:        astDt,
		Lifa:         lifa,
		Bazi:         [8]int{-1, -1, -1, -1, -1, -1, -1, -1},
		ShiShen:      []int{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		QiYunDelta:   [3]int{0, 0, 0},
	}

	return bz
}

func (bzm *BaziMeta) calcSolarTerms() (int, float64, float64) {
	mindx := 0
	headJieqi := &BZDateTime{}
	headJieqi = headJieqi.New(BzDataTimeInitTime)
	tailJieqi := &BZDateTime{}
	tailJieqi = tailJieqi.New(BzDataTimeInitTime)

	year := bzm.Bzdt.Year
	month := bzm.Bzdt.Month
	b, k, b1, k1 := bkCalc(year, month, 12)
	daxue := b - k
	lichun := b + k
	JD := bzm.Bzdt.ToJD()
	var D float64
	if JD < b {
		mindx = 0
		headJieqi.SetFromJD(b - k)
		tailJieqi.SetFromJD(b)
	} else if b <= JD && JD < b1 {
		n := math.Floor((JD - b) / k)
		D = k*n + b
		mindx = int(n) + 1
		headJieqi.SetFromJD(D)
		tailJieqi.SetFromJD(D + k)
	} else {
		n := math.Floor((JD - b1) / k1)
		D = k1*n + b1
		mindx = 13
		headJieqi.SetFromJD(b1)
		tailJieqi.SetFromJD(b1 + k1)
	}

	bzm.HeadJieQi = headJieqi
	bzm.TailJieQi = tailJieqi
	if bzm.BzAst {
		bzm.HeadJieQiAst = calcAST(headJieqi, bzm.Longitude)
		bzm.TailJieQiAst = calcAST(tailJieqi, bzm.Longitude)
	}
	return mindx, daxue, lichun
}

func (bzm *BaziMeta) calcJiaoYunDate(tFlag bool) {
	var bd, hJq, tJq float64

	if bzm.BzAst {
		bd = bzm.Astdt.ToJD()
		hJq = bzm.HeadJieQiAst.ToJD()
		tJq = bzm.TailJieQiAst.ToJD()
	} else {
		bd = bzm.Bzdt.ToJD()
		hJq = bzm.HeadJieQi.ToJD()
		tJq = bzm.TailJieQi.ToJD()
	}
	var dt1, dt2 float64
	if tFlag {
		dt1 = hJq
		dt2 = bd
	} else {
		dt1 = bd
		dt2 = tJq
	}

	offset := dt2 - dt1
	deltaY := math.Floor(offset / (365.2422 / 120.0))
	offset = (offset/(365.2422/120.0) - deltaY) * 365.2422
	bzm.Jydt.SetFromJD(bd + 365.2422*deltaY + offset)

	deltaM := math.Floor(offset / (365.2422 / 12.0))
	offset -= deltaM * (365.2422 / 12.0)
	deltaD := math.Floor(offset)
	bzm.QiYunDelta = [3]int{int(deltaY), int(deltaM), int(deltaD)}
}

func (bzm *BaziMeta) CalcBazi() {
	bd := &BZDateTime{}
	if bzm.BzAst {
		bd = bzm.Astdt
	} else {
		bd = bzm.Bzdt
	}
	mindx, daxue, lichun := bzm.calcSolarTerms()
	if bzm.Bzdt.ToJD() < lichun {
		lichun -= 365
	}

	//1984年立春作为起点
	//lichun -= SolarTerms.J2000 - 365.25*16 + 35
	//self.lichun = lichun
	lichun -= 2445735
	Nianzhu := math.Floor(lichun/365.2422+0.5) + 9000
	//#daxue -= SolarTerms.J2000 - 24 - 365
	//#1998 12 7日 大雪作为起点
	daxue -= 2451155
	Yuezhu := math.Floor(daxue/365.2422+0.4)*12 + float64(mindx) + 90000
	//Nianzhu = int(Nianzhu)
	//Yuezhu = int(Yuezhu)

	bzm.Bazi[0] = int(Nianzhu) % 10
	bzm.Bazi[1] = int(Nianzhu) % 12
	bzm.Bazi[2] = int(Yuezhu) % 10
	bzm.Bazi[3] = int(Yuezhu) % 12

	//#2000年1月7日作为起点
	Rizhu := int(math.Mod(bd.JdD-JulianCalendarJ2000-6+9000000, 60))
	h := bd.Hour
	t := int(h / 2)
	Dflag := 0
	if h%2 == 1 {
		t += 1
	}

	if t == 12 {
		Dflag = 1
		t = 0
	}

	if Dflag == 1 {
		Rizhu += 1
		Rizhu %= 60
	}

	Shizhu := int((Rizhu%5)*12 + t)
	bzm.Bazi[4] = Rizhu % 10
	bzm.Bazi[5] = Rizhu % 12
	bzm.Bazi[6] = Shizhu % 10
	bzm.Bazi[7] = Shizhu % 12

	j := bzm.Bazi[4]
	if j%2 == 0 {
		for i := 0; i < 10; i++ {
			bzm.ShiShen[j] = i
			j += 1
			j %= 10
		}
	} else {
		for k := 0; k < 9; k += 2 {
			bzm.ShiShen[j] = k
			bzm.ShiShen[j-1] = k + 1
			j += 2
			j %= 10
		}
	}
	bzm.ShiShen = append(bzm.ShiShen, 10)
	isFemale := 0
	if bzm.Gender {
		isFemale = 1 //女性
	} else {
		isFemale = 0
	}
	flag := false
	if (bzm.Bazi[0]%2)^(isFemale) == 1 {
		flag = true
	} else {
		flag = false
	}
	bzm.calcJiaoYunDate(flag)
}

func (bzm *BaziMeta) GetBazi() string {
	var output []string
	isFemale := 0
	gender := ""
	if bzm.Gender {
		isFemale = 1 //女性
		gender = "坤"
	} else {
		isFemale = 0
		gender = "乾"
	}
	opt := fmt.Sprintf("%s:%s%s %s%s %s%s %s%s", gender,
		gTianGan[bzm.Bazi[0]], gDiZhi[bzm.Bazi[1]],
		gTianGan[bzm.Bazi[2]], gDiZhi[bzm.Bazi[3]],
		gTianGan[bzm.Bazi[4]], gDiZhi[bzm.Bazi[5]],
		gTianGan[bzm.Bazi[6]], gDiZhi[bzm.Bazi[7]])

	output = append(output, opt)
	offsets := [2]int{0, 0}
	if ((bzm.Bazi[0] % 2) ^ (isFemale)) == 1 {
		offsets = [2]int{9, 11}
	} else {
		offsets = [2]int{1, 1}
	}

	opt = "大运:"
	j := bzm.Bazi[2]
	k := bzm.Bazi[3]
	for i := 0; i < 8; i++ {
		j += offsets[0]
		j %= 10
		k += offsets[1]
		k %= 12
		opt += fmt.Sprintf("%s%s ", gTianGan[j], gDiZhi[k])
	}
	output = append(output, opt)
	jydt := bzm.Jydt
	opt = fmt.Sprintf("  %d年%d月%d日交运", jydt.Year, jydt.Month, jydt.Day)
	output = append(output, opt)
	ret := strings.Join(output, "   ")
	return ret
}

func (bzm *BaziMeta) PrintBazi() string {
	isFemale := 0
	if bzm.Gender {
		isFemale = 1 //女性
	} else {
		isFemale = 0
	}
	var output []string
	opt := fmt.Sprintf("            %s        %s        %s        %s",
		gShiShen[bzm.ShiShen[bzm.Bazi[0]]],
		gShiShen[bzm.ShiShen[bzm.Bazi[2]]], "日元",
		gShiShen[bzm.ShiShen[bzm.Bazi[6]]])
	output = append(output, opt)

	opt = fmt.Sprintf("%s       %s           %s           %s           %s",
		gGender[isFemale], gTianGan[bzm.Bazi[0]], gTianGan[bzm.Bazi[2]],
		gTianGan[bzm.Bazi[4]], gTianGan[bzm.Bazi[6]])
	output = append(output, opt)

	opt = fmt.Sprintf("              %s           %s           %s           %s",
		gDiZhi[bzm.Bazi[1]], gDiZhi[bzm.Bazi[3]],
		gDiZhi[bzm.Bazi[5]], gDiZhi[bzm.Bazi[7]])
	output = append(output, opt)

	for i := 0; i < 3; i++ {
		opt = "         "
		for j := 1; j < 8; j += 2 {
			k := gCangGan[bzm.Bazi[j]][i]
			if k < 10 {
				opt += fmt.Sprintf("%s %s    ", gTianGan[k], gShiShen[bzm.ShiShen[k]])
			} else {
				opt += "                "
			}
		}
		output = append(output, opt)
	}
	output = append(output, "")
	flag := (bzm.Bazi[0] % 2) ^ (isFemale)
	offsets := [2]int{0, 0}
	if flag == 1 {
		offsets = [2]int{9, 11}
	} else {
		offsets = [2]int{1, 1}
	}

	j := bzm.Bazi[2]
	k := bzm.Bazi[3]
	opt = "大运"
	for i := 0; i < 8; i++ {
		j += offsets[0]
		j %= 10
		k += offsets[1]
		k %= 12

		opt += fmt.Sprintf("  %s%s", gTianGan[j], gDiZhi[k])
	}
	output = append(output, opt)
	ret := strings.Join(output, "\n")
	return ret
}

func (bzm *BaziMeta) GetAge(now *BZDateTime) string {
	bdJd := bzm.Bzdt.ToJD()
	if now.ToJD() < bdJd {
		age := "『年龄』 " + "未出生"
		return age
	} else {
		startLichun := calcLichun(bzm.Bzdt)
		endLichun := calcLichun(now)
		age := math.Floor((endLichun-startLichun)/365.2422+0.4) + 1
		if age > 120.0 {
			return "『年龄』 " + "历史人物"
		}
		return fmt.Sprintf("『年龄』 虚岁%d", int(age))
	}
}

func (bzm *BaziMeta) GetLifa() string {
	if bzm.Lifa >= 1 && bzm.Lifa <= 10 {
		s := fmt.Sprintf("依据 %s 数据拟合", gLifaList[bzm.Lifa-1])
		return "『定气方式』 " + s
	} else if bzm.Lifa > 10 && bzm.Lifa < 13 {
		return "『定气方式』 依据尤氏子平历计算节气交接日期"
	} else {
		return "『定气方式』 现代农历定气"
	}
}

func (bzm *BaziMeta) GetAst() string {
	var output []string
	opt := fmt.Sprintf("『出生地经度』%f", bzm.Longitude)
	output = append(output, opt)
	opt = fmt.Sprintf("『出生地真太阳时』 %s", bzm.Astdt.getFormatStr())
	output = append(output, opt)
	ret := strings.Join(output, "\n")
	return ret
}

func (bzm *BaziMeta) GetUserName() string {
	return "『姓名』  " + bzm.Name
}

func (bzm *BaziMeta) GetUserGender() string {
	if bzm.Gender {
		return "『性别』  " + "女"
	} else {
		return "『性别』  " + "男"
	}
}

func (bzm *BaziMeta) GetSolarTerms() string {
	var output []string
	var astStr string
	hJq := &BZDateTime{}
	tJq := &BZDateTime{}
	if bzm.BzAst {
		astStr = "出生地真太阳时"
		hJq = bzm.HeadJieQiAst
		tJq = bzm.TailJieQiAst
	} else {
		astStr = ""
		hJq = bzm.HeadJieQi
		tJq = bzm.TailJieQi
	}
	opt := fmt.Sprintf("『%s%s』 %s", gJieQi[bzm.Bazi[3]], astStr, hJq.getFormatStr())
	output = append(output, opt)
	opt = fmt.Sprintf("『%s%s』 %s", gJieQi[bzm.Bazi[3]+1], astStr, tJq.getFormatStr())
	output = append(output, opt)
	ret := strings.Join(output, "\n")
	return ret
}

func (bzm *BaziMeta) GetShenXiao() string {
	lySx := gShengXiao[bzm.Bazi[1]]
	return fmt.Sprintf("『生肖』 %s", lySx)
}

func (bzm *BaziMeta) GetSolarBirth() string {
	birthStr := ""
	if bzm.Bzdt.IsGcal {
		birthStr = "『公历生日』 "
	} else {
		birthStr = "『儒略历生日』 "
	}

	birthStr += bzm.Bzdt.getFormatStr()
	return birthStr
}

func (bzm *BaziMeta) GetLunarBirth() string {
	_, _, lunarInfo := bzm.Bzdt.GetLunarInfo()
	res := "『农历生日』 " + lunarInfo
	return res
}

func (bzm *BaziMeta) GetDateOfBirth() string {
	lyNh, _, _ := bzm.Bzdt.GetLunarInfo()
	return fmt.Sprintf("『出生年代』 %s", lyNh)
}

func (bzm *BaziMeta) GetJieLing() string {
	res := ""
	if bzm.BzAst {
		res += bzm.GetAst()
		res += "\n"
	}

	res += bzm.GetSolarTerms()
	res += "\n\n"
	return res
}

func (bzm *BaziMeta) GetSiZhu() [][]string {
	bazi := bzm.Bazi
	ssgx := bzm.ShiShen
	isFemale := 0
	gender := ""
	if bzm.Gender {
		isFemale = 1 //女性
		gender = gGender[isFemale]
	} else {
		isFemale = 0
		gender = gGender[isFemale]
	}

	var shiShenArray, tianGanArray, diZhiArray []string
	//shishen
	shiShenArray = append(shiShenArray, "")
	shiShenArray = append(shiShenArray, gShiShen[ssgx[bazi[0]]])
	shiShenArray = append(shiShenArray, gShiShen[ssgx[bazi[2]]])
	shiShenArray = append(shiShenArray, "日元")
	shiShenArray = append(shiShenArray, gShiShen[ssgx[bazi[6]]])
	//tiangan
	tianGanArray = append(tianGanArray, gender)
	tianGanArray = append(tianGanArray, gTianGan[bazi[0]])
	tianGanArray = append(tianGanArray, gTianGan[bazi[2]])
	tianGanArray = append(tianGanArray, gTianGan[bazi[4]])
	tianGanArray = append(tianGanArray, gTianGan[bazi[6]])
	//dizhi
	diZhiArray = append(diZhiArray, "")
	diZhiArray = append(diZhiArray, gDiZhi[bazi[1]])
	diZhiArray = append(diZhiArray, gDiZhi[bazi[3]])
	diZhiArray = append(diZhiArray, gDiZhi[bazi[5]])
	diZhiArray = append(diZhiArray, gDiZhi[bazi[7]])

	var output [][]string
	output = append(output, shiShenArray)
	output = append(output, tianGanArray)
	output = append(output, diZhiArray)

	for i := 0; i < 3; i++ {
		var cangGan []string
		opt := ""

		if i == 0 {
			opt = fmt.Sprintf("%s", "藏干")
			cangGan = append(cangGan, opt)
		} else {
			cangGan = append(cangGan, "")
		}

		for j := 1; j < 8; j += 2 {
			k := gCangGan[bazi[j]][i]
			if k < 10 {
				opt = fmt.Sprintf("%s %s", gTianGan[k], gShiShen[ssgx[k]])
				cangGan = append(cangGan, opt)
			} else {
				cangGan = append(cangGan, "")
			}
		}
		output = append(output, cangGan)
	}

	return output
}

func (bzm *BaziMeta) GetBaziBasic() string {
	bazi := bzm.Bazi
	ssgx := bzm.ShiShen
	isFemale := 0
	gender := ""
	if bzm.Gender {
		isFemale = 1 //女性
		gender = gGender[isFemale]
	} else {
		isFemale = 0
		gender = gGender[isFemale]
	}

	resultlist := ""
	resultlist += fmt.Sprintf("%-18s%-18s%-18s%-18s%-18s\n",
		"  ", gShiShen[ssgx[bazi[0]]], gShiShen[ssgx[bazi[2]]], "日元", gShiShen[ssgx[bazi[6]]])

	resultlist += fmt.Sprintf("%-16s%-20s%-20s%-20s%-20s\n", gender, gTianGan[bazi[0]], gTianGan[bazi[2]], gTianGan[bazi[4]], gTianGan[bazi[6]])

	resultlist += fmt.Sprintf("%-18s%-20s%-20s%-20s%-20s\n", "  ", gDiZhi[bazi[1]], gDiZhi[bazi[3]], gDiZhi[bazi[5]], gDiZhi[bazi[7]])

	//resultlist += "\n"

	for i := 0; i < 3; i++ {
		opt := ""
		for j := 1; j < 8; j += 2 {
			k := gCangGan[bazi[j]][i]
			if k < 10 {
				opt += fmt.Sprintf("%s %s%-16s", gTianGan[k], gShiShen[ssgx[k]], "  ")
			} else {
				opt += fmt.Sprintf("%-20s", "  ")
			}
		}
		if i == 0 {
			resultlist += fmt.Sprintf("%-16s%-14s\n", "藏干", opt)
		} else {
			resultlist += fmt.Sprintf("%-16s%-16s\n", "  ", opt)
		}

	}
	return resultlist
}

func (bzm *BaziMeta) GetQiYun() string {
	qy := bzm.QiYunDelta
	res := fmt.Sprintf("起运 命主于出生后%d年%d个月%d天后起运\n", qy[0], qy[1], qy[2])
	return res
}

func (bzm *BaziMeta) GetJiaoYun() string {
	jydt := bzm.Jydt
	res := fmt.Sprintf("交运 命主于%s交运\n", jydt.getFormatStr())
	return res
}

func (bzm *BaziMeta) GetDaYunList() []string {
	isFemale := 0
	if bzm.Gender {
		isFemale = 1 //女性
	} else {
		isFemale = 0
	}

	flag := (bzm.Bazi[0] % 2) ^ (isFemale)
	offsets := [2]int{0, 0}
	if flag == 1 {
		offsets = [2]int{9, 11}
	} else {
		offsets = [2]int{1, 1}
	}

	j := bzm.Bazi[2]
	k := bzm.Bazi[3]
	opt := ""
	var output []string
	for i := 0; i < 8; i++ {
		j += offsets[0]
		j %= 10
		k += offsets[1]
		k %= 12

		opt = fmt.Sprintf("%s%s", gTianGan[j], gDiZhi[k])
		output = append(output, opt)
	}

	return output
}

func (bzm *BaziMeta) GetStartYearList() []int32 {
	jydt := bzm.Jydt
	res, err := jydt.ToLunarDate()
	if err != nil {
		fmt.Printf("error msg: %s\n", err)
	}

	jyYear := res.Year
	var output []int32
	for i := 0; i < 8; i++ {
		output = append(output, int32(jyYear))
		jyYear += 10
	}
	return output
}

func (bzm *BaziMeta) GetFleetingYearList() []string {
	jydt := bzm.Jydt
	res1, err1 := jydt.ToLunarDate()
	if err1 != nil {
		fmt.Printf("error msg: %s\n", err1)
	}

	jyYear := res1.Year

	p := (jyYear + 6) % 10
	q := (jyYear + 8) % 12
	var output []string
	opt := ""
	for i := 0; i < 10; i++ {
		m := q
		for j := 0; j < 8; j++ {
			opt = fmt.Sprintf("%s%s", gTianGan[p], gDiZhi[m])
			m += 10
			m %= 12
			output = append(output, opt)
		}
		p += 1
		p %= 10
		q += 1
		q %= 12
	}
	return output
}

func (bzm *BaziMeta) GetEndYearList() []int32 {
	jydt := bzm.Jydt
	res, err := jydt.ToLunarDate()
	if err != nil {
		fmt.Printf("error msg: %s\n", err)
	}

	jyYear := res.Year
	endYear := jyYear + 9
	var output []int32
	for i := 0; i < 8; i++ {
		output = append(output, int32(endYear))
		endYear += 10
	}

	return output
}

func (bzm *BaziMeta) RenderImage(filePath string) (*image.RGBA, error) {
	srcWidth := 475
	srcHeight := 720
	img := image.NewRGBA(image.Rect(0, 0, srcWidth, srcHeight))
	nowStr := FormatDateTime(time.Now(), YYYYMMDD_HH_MM_SS_SSS_EN)
	now := &BZDateTime{}
	now = now.New(nowStr)
	nowJd := now.ToJD()

	zodiacIcons := []string{"rat", "cow", "tiger", "rabbit", "dragon", "snake", "horse", "goat", "monkey", "hen", "dog", "pig"}
	zodiacsImagePath := path.Join(filePath, "zodiacs", zodiacIcons[bzm.Bazi[1]]+".png")

	zodiac, err := os.Open(zodiacsImagePath)
	if err != nil {
		return img, err
	}

	defer func() {
		zodiac.Close()
	}()
	zodiacImg, _, err := image.Decode(zodiac)
	if err != nil {
		return img, err
	}

	backImagePath := path.Join(filePath, "zodiacs", "paper.jpeg")
	backGround, err := os.Open(backImagePath)
	defer backGround.Close()
	backImg, _, err := image.Decode(backGround) // 打开背景图片
	if err != nil {
		return img, err
	}

	//插入背景图和生肖图
	draw.Draw(img, img.Bounds(), backImg, image.Pt(0, 0), draw.Over)
	draw.Draw(img, image.Rect(28, 150, srcWidth, srcHeight), zodiacImg, image.Pt(0, 0), draw.Over)

	//字体设置
	fontPath := path.Join(filePath, "fonts", "WenYue_GuDianMingChaoTi_JRFC.otf")
	//fontPath := path.Join(filePath, "fonts", "WenYue.otf")
	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		return img, err
	}

	f, err := opentype.Parse(fontBytes)
	if err != nil {
		return img, err
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    10,
		DPI:     100,
		Hinting: font.HintingNone,
	})

	defer func(face font.Face) {
		_ = face.Close()
	}(face)

	if err != nil {
		return img, err
	}

	//姓名，性别，生肖，年龄
	resultlist := bzm.GetUserName() + bzm.GetUserGender() + bzm.GetShenXiao() + bzm.GetAge(now) + "\n"
	//公历生日
	resultlist += bzm.GetSolarBirth() + "\n"
	//出生年代，农历生日
	lyNh, _, lunarInfo := bzm.Bzdt.GetLunarInfo()
	resultlist += fmt.Sprintf("『出生年代』 %s\n", lyNh)
	resultlist += fmt.Sprintf("『农历生日』 %s\n", lunarInfo)
	//定气方式
	resultlist += bzm.GetLifa() + "\n"
	//节气
	if bzm.BzAst {
		resultlist += bzm.GetAst()
	}
	resultlist += bzm.GetJieLing()

	//四柱藏干盘面排列
	bazi := bzm.Bazi
	ssgx := bzm.ShiShen

	isFemale := 0
	gender := ""
	if bzm.Gender {
		isFemale = 1 //女性
		gender = gGender[isFemale]
	} else {
		isFemale = 0
		gender = gGender[isFemale]
	}

	resultlist += fmt.Sprintf("%-18s%-17s%-17s%-17s%-17s\n", "  ", gShiShen[ssgx[bazi[0]]], gShiShen[ssgx[bazi[2]]], "日元", gShiShen[ssgx[bazi[6]]])

	resultlist += fmt.Sprintf("%-14s%-20s%-20s%-20s%-20s\n", gender, gTianGan[bazi[0]], gTianGan[bazi[2]], gTianGan[bazi[4]], gTianGan[bazi[6]])

	resultlist += fmt.Sprintf("%-20s%-20s%-20s%-20s%-20s\n", "  ", gDiZhi[bazi[1]], gDiZhi[bazi[3]], gDiZhi[bazi[5]], gDiZhi[bazi[7]])

	resultlist += "\n"

	multilns := strings.Split(resultlist, "\n")
	d := font.Drawer{}
	for i, v := range multilns {
		d = font.Drawer{
			Dst:  img,
			Src:  image.Black,
			Face: face,
			Dot:  fixed.P(36, 19*i+19),
		}
		d.DrawString(v)
	}

	for i := 0; i < 3; i++ {
		if 0 == i {
			d = font.Drawer{
				Dst:  img,
				Src:  image.Black,
				Face: face,
				Dot:  fixed.P(36, 19*11+19),
			}
			d.DrawString("藏干")
		}
		for j := 1; j < 8; j += 2 {
			k := gCangGan[bazi[j]][i]
			if k < 10 {
				d := font.Drawer{
					Dst:  img,
					Src:  image.Black,
					Face: face,
					Dot:  fixed.P(95+(j-1)*40, 19+19*(11+i)),
				}
				opt := fmt.Sprintf("%s %s", gTianGan[k], gShiShen[ssgx[k]])
				d.DrawString(opt)
			}
		}
	}
	resultlist = bzm.GetQiYun()
	resultlist += bzm.GetJiaoYun()
	resultlist += "\n"
	multilns = strings.Split(resultlist, "\n")
	for i, v := range multilns {
		d = font.Drawer{
			Dst:  img,
			Src:  image.Black,
			Face: face,
			Dot:  fixed.P(36, 19*i+19+19*15),
		}
		d.DrawString(v)
	}

	//大运排列
	d = font.Drawer{
		Dst:  img,
		Src:  image.Black,
		Face: face,
		Dot:  fixed.P(36, 19*18+19),
	}
	d.DrawString("大运")

	flag := (bzm.Bazi[0] % 2) ^ (isFemale)
	offsets := [2]int{0, 0}
	if flag == 1 {
		offsets = [2]int{9, 11}
	} else {
		offsets = [2]int{1, 1}
	}

	j := bzm.Bazi[2]
	k := bzm.Bazi[3]

	tmp := bzm.Jydt.ToJD()

	for i := 0; i < 8; i++ {
		j += offsets[0]
		j %= 10
		k += offsets[1]
		k %= 12
		selectedColor := image.Black
		if tmp <= nowJd && nowJd < tmp+365.2422*10 {
			selectedColor = image.NewUniform(color.RGBA{255, 0, 0, 255})
		}

		d = font.Drawer{
			Dst:  img,
			Src:  selectedColor,
			Face: face,
			Dot:  fixed.P(81+i*47, 19*19),
		}

		opt := fmt.Sprintf("%s%s", gTianGan[j], gDiZhi[k])
		d.DrawString(opt)
		tmp += 365.2422 * 10
	}

	//起运时间
	d = font.Drawer{
		Dst:  img,
		Src:  image.Black,
		Face: face,
		Dot:  fixed.P(36, 19*20+19),
	}
	d.DrawString("起于     ")

	jydt := bzm.Jydt
	res, err := jydt.ToLunarDate()
	if err != nil {
		return img, err
	}

	nowLy, err := now.ToLunarDate()
	if err != nil {
		return img, err
	}

	jyYear := res.Year
	lnindx := []int{0, 8}

	for i := 0; i < 8; i++ {
		opt := fmt.Sprintf("%-04d", jyYear)
		d = font.Drawer{
			Dst:  img,
			Src:  image.Black,
			Face: face,
			Dot:  fixed.P(81+i*47, 19*21),
		}
		d.DrawString(opt)

		if jyYear <= nowLy.Year && nowLy.Year < jyYear+10 {
			lnindx[1] = i
		}
		jyYear += 10
	}

	//流年列表
	d = font.Drawer{
		Dst:  img,
		Src:  image.Black,
		Face: face,
		Dot:  fixed.P(36, 19*22+19),
	}
	d.DrawString("流年")

	lnindx[0] = ((nowLy.Year%10 - jyYear%10) + 10) % 10
	jyYear -= 80
	p := (jyYear + 6) % 10
	q := (jyYear + 8) % 12

	for i := 0; i < 10; i++ {
		m := q
		for j := 0; j < 8; j++ {
			opt := fmt.Sprintf("%s%s", gTianGan[p], gDiZhi[m])
			if i == lnindx[0] && j == lnindx[1] {
				if i < 5 {
					d = font.Drawer{
						Dst:  img,
						Src:  image.NewUniform(color.RGBA{255, 0, 0, 255}),
						Face: face,
						Dot:  fixed.P(81+j*47, 19*23+i*19),
					}
					d.DrawString(opt)
				} else {
					d = font.Drawer{
						Dst:  img,
						Src:  image.NewUniform(color.RGBA{255, 0, 0, 255}),
						Face: face,
						Dot:  fixed.P(81+j*47, 19*24+i*19),
					}
					d.DrawString(opt)
				}
			} else if i < 5 {
				d = font.Drawer{
					Dst:  img,
					Src:  image.Black,
					Face: face,
					Dot:  fixed.P(81+j*47, 19*23+i*19),
				}
				d.DrawString(opt)
			} else {
				d = font.Drawer{
					Dst:  img,
					Src:  image.Black,
					Face: face,
					Dot:  fixed.P(81+j*47, 19*24+i*19),
				}
				d.DrawString(opt)
			}

			m += 10
			m %= 12
		}
		p += 1
		p %= 10
		q += 1
		q %= 12
	}

	//交运结束列表
	jyYear += 9
	d = font.Drawer{
		Dst:  img,
		Src:  image.Black,
		Face: face,
		Dot:  fixed.P(36, 19*35),
	}
	d.DrawString("止于")
	for i := 0; i < 8; i++ {
		opt := fmt.Sprintf("%-04d", jyYear)
		d = font.Drawer{
			Dst:  img,
			Src:  image.Black,
			Face: face,
			Dot:  fixed.P(81+i*47, 19*35),
		}
		d.DrawString(opt)
		jyYear += 10
	}

	return img, err

}
