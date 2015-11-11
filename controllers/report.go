package controllers

import (
	"LTE_Security/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"strconv"
)

type ReportData struct {
	Mac   string `json:"mac"`
	Mid   int64  `json:"mid"`
	Psn   int64  `json:"psn"`
	Error int64  `json:"error"`
	Code  int64  `json:"code"`
}

type ReportActData struct {
	Mac   string `json:"mac"`
	Mid   int64  `json:"mid"`
	Psn   int64  `json:"psn"`
	Error int64  `json:"error"`
	Code  int64  `json:"code"`
	Act   string `json:"act"`
}

type ReportController struct {
	beego.Controller
}

var (
	RE_devicesCount int64 = 0  //default device numbers in database
	RE_totalPages   int64 = 0  //default page numbers to show
	RE_curPage      int64 = 1  //default current page
	RE_listCount    int64 = 10 //show numbers of deviceinfo every page
	RE_curCount     int64 = 0  //total numbers of deviceinfo had been showed
)

func (this *ReportController) Get() {
	//check islogin
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}

	this.TplNames = "report.html"

	//重置参数
	devicesCount = 0
	totalPages = 0
	curPage = 1
	listCount = 10
	curCount = 0

	//recive request listcount info
	listcount := this.Input().Get("listcount")
	beego.Debug("listcount=", listcount)
	if !(listcount == "") {
		count, _ := strconv.Atoi(listcount)
		switch count {
		case 20:
			RE_listCount = 20
			RE_curCount = 0
			RE_curPage = 1
		case 30:
			RE_listCount = 30
			RE_curCount = 0
			RE_curPage = 1
		default:
			RE_listCount = 10
			RE_curCount = 0
			RE_curPage = 1
		}
	}

	//calc variables to show
	RE_devicesCount = models.GetReportDevicesCount()
	beego.Debug("devicesCount=", RE_devicesCount)
	if RE_devicesCount%RE_listCount > 0 {
		RE_totalPages = RE_devicesCount/RE_listCount + 1
	} else {
		RE_totalPages = RE_devicesCount / listCount
	}

	//recive request op info
	ope := this.Input().Get("op")
	beego.Debug("ope=", ope)
	switch ope {
	case "firstpage":
		RE_curCount = 0
		RE_curPage = 1
	case "prepage":
		if RE_curPage > 1 {
			RE_curCount -= listCount
			RE_curPage -= 1
		}
	case "nextpage":
		if RE_curPage < RE_totalPages {
			RE_curCount += listCount
			RE_curPage += 1
		}
	case "lastpage":
		RE_curCount = RE_listCount * (RE_totalPages - 1)
		RE_curPage = RE_totalPages
	}

	devices, nums, ok := models.GetReportDevices(RE_listCount, RE_curCount)
	if ok {
		beego.Info("GetDevices success!")
		this.Data["Devices"] = devices
		this.Data["DevicesNum"] = RE_devicesCount
		this.Data["CurPage"] = curPage
		this.Data["TotalPages"] = RE_totalPages
	}
	if nums <= 0 {
		this.Data["CurPage"] = 0
		this.Data["NoInfo"] = "没有记录!"
	}
}

/* ret
(0)report success
(-1)uname and pwd not match
(-2)input data error
(-3)report failed
*/

func (this *ReportController) Post() {
	data := ReportData{}

	//check auth
	uname, pwd, ok := this.Ctx.Request.BasicAuth()
	if !ok {
		beego.Info("get client  Request.BasicAuth failed!")
		data.Code = -1
		writeContent, _ := json.Marshal(data)
		this.Ctx.WriteString(string(writeContent))
		return
	}
	user := models.Userinfo{
		Username: uname,
		Password: pwd,
	}
	ok = models.CheckAccount(&user)
	if !ok {
		beego.Info("user/pwd not matched!")
		data.Code = -1
		writeContent, _ := json.Marshal(data)
		this.Ctx.WriteString(string(writeContent))
		return
	}
	//get client data
	var deviceinfo models.Deviceinfo
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &deviceinfo)
	if err != nil {
		beego.Error(err)
		data.Code = -2
		writeContent, _ := json.Marshal(data)
		this.Ctx.WriteString(string(writeContent))
		return
	}
	if deviceinfo.Mac == "" || deviceinfo.Mid == 0 || deviceinfo.Psn == 0 {
		beego.Error(" input content error!")
		data.Code = -2
		writeContent, _ := json.Marshal(data)
		this.Ctx.WriteString(string(writeContent))
		return
	}

	data.Mac = deviceinfo.Mac
	data.Mid = deviceinfo.Mid
	data.Psn = deviceinfo.Psn
	data.Error = deviceinfo.Error
	//report
	var act string
	act, data.Code = models.ReportDeivce(&deviceinfo)

	//action不为空，加入返回值
	if act != "" {
		ret := RegisterActData{
			Mac:  data.Mac,
			Mid:  data.Mid,
			Psn:  data.Psn,
			Code: data.Code,
			Act:  act,
		}
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	} else {
		writeContent, _ := json.Marshal(data)
		this.Ctx.WriteString(string(writeContent))
		return
	}
}
