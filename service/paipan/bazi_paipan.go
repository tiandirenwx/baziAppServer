package paipan

import (
	"bytes"
	"context"
	"encoding/base64"
	"html/template"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
	"trpc.app/app/baziAppServer/errorcode"
	"trpc.group/trpc-go/trpc-go"

	"trpc.app/app/baziAppServer/common/bazicore"
	"trpc.app/app/baziAppServer/common/utils"

	pb "trpc.bazi.paipan"

	"trpc.group/trpc-go/trpc-go/log"
)

type BaziPaipanImpl struct{}

func checkPaiPanParams(ctx context.Context, req *pb.PaiPanRequest) error {
	if req.GetYear() <= 0 || req.GetYear() > 9999 {
		return errorcode.WithMessage(errorcode.ParamsErr, "输入年份错误,目前支持年份为1~9999")
	}

	if req.GetMonth() <= 0 || req.GetMonth() > 13 {
		return errorcode.WithMessage(errorcode.ParamsErr, "输入月份错误,请保证输入的月数为1~12")
	}

	if req.GetDay() < 0 || req.GetDay() > 31 {
		return errorcode.WithMessage(errorcode.ParamsErr, "输入日期错误,请保证输入的日数为1~31")
	}

	if req.GetHour() < 0 || req.GetHour() > 24 {
		return errorcode.WithMessage(errorcode.ParamsErr, "输入时间小时数错误")
	}

	if req.GetMinute() < 0 || req.GetMinute() > 60 {
		return errorcode.WithMessage(errorcode.ParamsErr, "输入时间分钟数错误")
	}

	return nil
}

func calcBaziInfo(ctx context.Context, req *pb.PaiPanRequest) (bazi *bazicore.BaziMeta, err error) {
	var year, month, day int
	if req.GetCalendarType() == 1 {
		year = int(req.GetYear())
		month = int(req.GetMonth())
		day = int(req.GetDay())
	} else if req.GetCalendarType() == 2 {
		y, m, d, err := bazicore.Lunar2solar(int(req.GetYear()), int(req.GetMonth()), int(req.GetDay()),
			req.GetIsLeapMonth(), req.GetIsLeapMonth())
		if err != nil {
			log.ErrorContextf(ctx, "lunar date transfer to solar fail, error: %s", err)
			return nil, err
		}
		year = y
		month = m
		day = d
	} else {
		return nil, errorcode.WithMessage(errorcode.ParamsErr, "输入排盘历法错误,目前系统支持公历和农历排盘,其它暂不支持")
	}

	hour := int(req.GetHour())
	min := int(req.GetMinute())

	timeStr := bazicore.DateTimeStrByBazi(year, month, day, hour, min)
	log.DebugContextf(ctx, "time str: %s", timeStr)
	dt := &bazicore.BZDateTime{}
	dt = dt.New(timeStr)
	log.DebugContextf(ctx, "dt: %+v", dt)
	bz := &bazicore.BaziMeta{}
	bz = bz.New(dt, req.GetGender(), req.GetAstFlag(), float64(req.GetLongitude()), int(req.GetCalendar()), req.GetName())
	bz.CalcBazi()

	return bz, nil
}

func (s *BaziPaipanImpl) CreateBaziPaipan(ctx context.Context, req *pb.PaiPanRequest) (*pb.CreatePaiPanRsp, error) {
	// implement business logic here ...
	// ...
	rsp := &pb.CreatePaiPanRsp{}
	log.DebugContextf(ctx, "input param: %+v", req)

	if err := checkPaiPanParams(ctx, req); err != nil {
		return rsp, err
	}

	bz, err := calcBaziInfo(ctx, req)
	if bz == nil || err != nil {
		return rsp, err
	}

	// calc age
	nowStr := bazicore.FormatDateTime(time.Now(), bazicore.YYYYMMDD_HH_MM_SS_SSS_EN)
	now := &bazicore.BZDateTime{}
	now = now.New(nowStr)
	age := bz.GetAge(now)

	//sizhu info
	siZhuArray := bz.GetSiZhu()
	pBazi := &pb.BaziSiZhu{
		Shishen:  siZhuArray[0],
		Tiangan:  siZhuArray[1],
		Dizhi:    siZhuArray[2],
		Canggan1: siZhuArray[3],
		Canggan2: siZhuArray[4],
		Canggan3: siZhuArray[5],
	}

	//response
	rsp.Nonce = utils.Get32RandomString()
	rsp.Timestamp = nowStr
	rsp.UserName = bz.GetUserName()
	rsp.Gender = bz.GetUserGender()
	rsp.Shengxiao = bz.GetShenXiao()
	rsp.Age = age
	rsp.SolarBirth = bz.GetSolarBirth()
	rsp.LunarBirth = bz.GetLunarBirth()
	rsp.DateOfBirth = bz.GetDateOfBirth()
	rsp.DingQiType = bz.GetLifa()
	rsp.JieQi = bz.GetJieLing()
	rsp.Bazi = pBazi
	rsp.QiYun = bz.GetQiYun()
	rsp.JiaoYun = bz.GetJiaoYun()
	rsp.DaYun = bz.GetDaYunList()
	rsp.StartYear = bz.GetStartYearList()
	rsp.FleetingYear = bz.GetFleetingYearList()
	rsp.EndYear = bz.GetEndYearList()

	return rsp, nil
}

var ImageTemplate string = `<!DOCTYPE html>
<html lang="en"><head></head>
<body><img src="data:image/jpg;base64,{{.Image}}"></body>`

// Writeimagewithtemplate encodes an image "img" in jpeg format and writes it into ResponseWriter using a template.
func writeImageWithTemplate(ctx context.Context, w http.ResponseWriter, img *image.Image) {

	buffer := new(bytes.Buffer)
	var opt jpeg.Options
	opt.Quality = 100
	if err := jpeg.Encode(buffer, *img, &opt); err != nil {
		log.ErrorContext(ctx, "unable to encode image.")
	}

	str := base64.StdEncoding.EncodeToString(buffer.Bytes())
	if tmpl, err := template.New("image").Parse(ImageTemplate); err != nil {
		log.ErrorContext(ctx, "unable to parse image template.")
	} else {
		data := map[string]interface{}{"Image": str}
		if err = tmpl.Execute(w, data); err != nil {
			log.ErrorContext(ctx, "unable to execute template.")
		}
	}
}

// writeImage encodes an image "img" in jpeg format and writes it into ResponseWriter.
func writeImage(ctx context.Context, w http.ResponseWriter, img *image.Image) {

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.ErrorContext(ctx, "unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.ErrorContext(ctx, "unable to write image.")
	}
}

func (s *BaziPaipanImpl) RenderPaiPanImage(ctx context.Context, req *pb.PaiPanRequest) (*pb.RenderPaiPanImageRsp, error) {
	// implement business logic here ...
	// ...
	rsp := &pb.RenderPaiPanImageRsp{}
	if err := checkPaiPanParams(ctx, req); err != nil {
		return rsp, err
	}

	bz, err := calcBaziInfo(ctx, req)
	if bz == nil || err != nil {
		return rsp, err
	}

	configPath := path.Dir(trpc.ServerConfigPath)
	img, err := bz.RenderImage(configPath)
	if err != nil {
		log.ErrorContextf(ctx, "render image error: %s", err.Error())
		return rsp, err
	}

	buffer := new(bytes.Buffer)
	var opt jpeg.Options
	opt.Quality = 100
	if err := jpeg.Encode(buffer, img, &opt); err != nil {
		log.ErrorContext(ctx, "unable to encode image.")
		return rsp, err
	}

	str := base64.StdEncoding.EncodeToString(buffer.Bytes())
	rsp.Image = str

	fileName := path.Join(configPath, "result", bz.Name+".jpeg")
	outFile, err := os.Create(fileName)
	defer outFile.Close()
	if err != nil {
		log.ErrorContextf(ctx, "create file %s fail, error %s.", fileName, err.Error())
		return rsp, err
	}

	// ok, write out the data into the new JPEG file
	err = jpeg.Encode(outFile, img, &opt)
	if err != nil {
		log.ErrorContextf(ctx, "create file %s fail, error %s.", fileName, err.Error())
		return rsp, err
	}

	return rsp, nil
}
