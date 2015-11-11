package controllers

import (
	"LTE_Security/models"
	"github.com/astaxie/beego"
	"strconv"
)

var (
	devicesCount int64 = 0  //default device numbers in database
	totalPages   int64 = 0  //default page numbers to show
	curPage      int64 = 1  //default current page
	listCount    int64 = 10 //show numbers of deviceinfo every page
	curCount     int64 = 0  //total numbers of deviceinfo had been showed
)

type HomeController struct {
	beego.Controller
}

func (this *HomeController) Get() {
	//session认证
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}
	this.TplNames = "home.html"

	//重置参数
	devicesCount = 0
	totalPages = 0
	curPage = 1
	listCount = 10
	curCount = 0

	//获取每页显示条数
	listcount := this.Input().Get("ListCount")
	if listcount != "" {
		listCount, _ = strconv.ParseInt(listcount, 10, 0)
		beego.Debug("listCount=", listCount)
	}

	//获取设备总数
	devicescount := this.Input().Get("DevicesCount")
	if devicescount != "" {
		devicesCount, _ = strconv.ParseInt(devicescount, 10, 0)
		beego.Debug("devicesCount=", devicesCount)
	} else {
		devicesCount = models.GetDevicesCount()
	}

	//计算出总页数
	if devicesCount%listCount > 0 {
		totalPages = devicesCount/listCount + 1
	} else {
		totalPages = devicesCount / listCount
	}

	//获取当前页数
	curpage := this.Input().Get("CurPage")
	if curpage != "" {
		curPage, _ = strconv.ParseInt(curpage, 10, 0)
		beego.Debug("curPage=", curPage)
	} else {
		curPage = 1
	}
	if curPage > totalPages {
		curPage = totalPages
	}

	//计算出当前总条数
	curCount = listCount * (curPage - 1)

	//获取操作
	ope := this.Input().Get("op")
	switch ope {
	case "firstpage":
		curCount = 0
		curPage = 1
	case "prepage":
		if curPage > 1 {
			curCount -= listCount
			curPage -= 1
		}
	case "nextpage":
		if curPage < totalPages {
			curCount += listCount
			curPage += 1
		}
	case "lastpage":
		curCount = listCount * (totalPages - 1)
		curPage = totalPages
	}

	devices, nums, ok := models.GetDevices(listCount, curCount)
	if ok {
		this.Data["Devices"] = devices
		this.Data["DevicesNum"] = devicesCount
		this.Data["CurPage"] = curPage
		this.Data["ListCount"] = listCount
		this.Data["TotalPages"] = totalPages
	}
	if nums <= 0 {
		this.Data["CurPage"] = 0
		this.Data["NoInfo"] = "没有注册设备!"
	}
}

func (this *HomeController) Post() {
	//check islogin
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}

	this.TplNames = "login.html"
}
