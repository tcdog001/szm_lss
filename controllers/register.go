package controllers

import (
	"LTE_Security/models"
	"encoding/json"
	"github.com/astaxie/beego"
)

type RegisterData struct {
	Mac  string `json:"mac"`
	Mid  int64  `json:"mid"`
	Psn  int64  `json:"psn"`
	Guid string `json:"guid"`
	Code int64  `json:"code"`
}

type RegisterActData struct {
	Mac  string `json:"mac"`
	Mid  int64  `json:"mid"`
	Psn  int64  `json:"psn"`
	Guid string `json:"guid"`
	Code int64  `json:"code"`
	Act  string `json:"act"`
}

type RegisterController struct {
	beego.Controller
}

func (this *RegisterController) Get() {
	//check islogin
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}

	this.TplNames = "login.html"
}

/* ret
(0)deivce register success
(-1)uname and pwd not match
(-2)invaild mac address
(-3)repeated registration
(-4)device register over time
(-5)insert db failed
(-6)input data error
*/

func (this *RegisterController) Post() {
	data := RegisterData{}
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
		data.Code = -6
		writeContent, _ := json.Marshal(data)
		this.Ctx.WriteString(string(writeContent))
		return
	}
	if deviceinfo.Mac == "" || deviceinfo.Mid == 0 || deviceinfo.Psn == 0 {
		beego.Error(" input content error!")
		data.Code = -6
		writeContent, _ := json.Marshal(data)
		this.Ctx.WriteString(string(writeContent))
		return
	}
	data.Mac = deviceinfo.Mac
	data.Mid = deviceinfo.Mid
	data.Psn = deviceinfo.Psn
	//register
	var act string
	act, data.Guid, data.Code = models.RegisterDeivce(&deviceinfo)

	//action不为空，加入返回值
	if act != "" {
		ret := RegisterActData{
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
