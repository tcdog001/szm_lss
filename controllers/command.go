package controllers

import (
	"LTE_Security/models"
	"github.com/astaxie/beego"
)

type CommandController struct {
	beego.Controller
}

func (this *CommandController) Get() {
	//检查登录状态
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}
	this.TplNames = "command.html"

	//获取传递数据
	mac := this.Input().Get("mac")
	this.Data["Mac"] = mac
}

func (this *CommandController) Post() {
	//检查登录状态
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}

	//获取表单信息
	mac := this.Input().Get("mac")
	beego.Debug("mac=", mac)
	commandContent := this.Input().Get("commandContent")
	beego.Debug("commandContent=", commandContent)

	//下发action
	device := models.Deviceinfo{
		Mac: mac,
		Act: commandContent,
	}
	ok := models.SetAct(&device)
	if !ok {
		beego.Info("set action failed!")
	}

	//返回设备页面
	this.Redirect("/home", 302)
	return
}
