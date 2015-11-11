package routers

import (
	"LTE_Security/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//管理平台
	beego.Router("/", &controllers.LoginController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/home", &controllers.HomeController{})
	beego.Router("/modify", &controllers.ModifyController{})
	beego.Router("/noregister", &controllers.NoregisterController{})
	beego.Router("/import", &controllers.ImportController{})
	beego.Router("/command", &controllers.CommandController{})
	beego.Router("/report", &controllers.ReportController{})

	//与设备交互
	beego.Router("/service1", &controllers.RegisterController{})
	beego.Router("/service2", &controllers.VerifyController{})
	beego.Router("/service3", &controllers.ReportController{})
}
