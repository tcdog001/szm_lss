package controllers

import (
	"LTE_Security/models"
	"encoding/json"
	"github.com/astaxie/beego"
)

type VerifyData struct {
	Mac  string `json:"mac"`
	Mid  int64  `json:"mid"`
	Psn  int64  `json:"psn"`
	Guid string `json:"guid"`
	Code int64  `json:"code"`
}

type VerifyActData struct {
	Mac  string `json:"mac"`
	Mid  int64  `json:"mid"`
	Psn  int64  `json:"psn"`
	Guid string `json:"guid"`
	Code int64  `json:"code"`
	Act  string `json:"act"`
}

type VerifyController struct {
	beego.Controller
}

func (this *VerifyController) Get() {
	//check islogin
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}
	this.TplNames = "login.html"
}

/* return value
(0)verify success
(-1)uname and pwd not match
(-2)input data error
(-3)verify failed
*/
func (this *VerifyController) Post() {
	data := VerifyData{}
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

	if deviceinfo.Mac == "" || deviceinfo.Guid == "" || deviceinfo.Mid == 0 || deviceinfo.Psn == 0 {
		beego.Error(" input content error!")
		data.Code = -2
		writeContent, _ := json.Marshal(data)
		this.Ctx.WriteString(string(writeContent))
		return
	}

	data.Mac = deviceinfo.Mac
	data.Mid = deviceinfo.Mid
	data.Psn = deviceinfo.Psn
	data.Guid = deviceinfo.Guid

	//verify
	var act string
	act, data.Code = models.VerifyDevice(&deviceinfo)

	//action不为空，加入返回值
	if act != "" {
		ret := VerifyActData{
			Mac:  data.Mac,
			Mid:  data.Mid,
			Psn:  data.Psn,
			Guid: data.Guid,
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
