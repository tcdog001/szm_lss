package models

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"time"
)

type Userinfo struct {
	Uid           int64     `orm:"auto"`
	Username      string    `json:"uname"`
	Password      string    `json:"pwd"`
	CreatedTime   time.Time `json:"-"`
	LastLoginTime time.Time `orm:"null;auto_now;type(datetime)";json:"lastlogin"`
}

type Admininfo struct {
	Uid           int64     `orm:"auto"`
	Username      string    `json:"uname"`
	Password      string    `json:"pwd"`
	CreatedTime   time.Time `json:"-"`
	LastLoginTime time.Time `orm:"null;auto_now;type(datetime)";json:"lastlogin"`
}

type Deviceinfo struct {
	Mac              string    `orm:"pk";json:"mac"`
	Guid             string    `orm:"null";json:"guid"`
	Mid              int64     `orm:"null";json:"mid"`
	Psn              int64     `orm:"null";json:"psn"`
	Error            int64     `orm:"null";json:"error"`
	Act              string    `orm:"null";json:"act"`
	ImportTime       time.Time `json:"-"`
	InvalidTime      time.Time `json:"-"`
	ErrorTime        time.Time `orm:"null";json:"-"`
	RegistrationTime time.Time `orm:"null";json:"-"`
	UpdateTime       time.Time `orm:"null";json:"-"`
}

func init() {
	orm.RegisterModel(new(Userinfo), new(Admininfo), new(Deviceinfo))

	err := orm.RegisterDriver("mysql", orm.DR_MySQL)
	if err != nil {
		beego.Critical(err)
	}

	dbIp := beego.AppConfig.String("DbIp")
	dbPort := beego.AppConfig.String("DbPort")
	dbName := beego.AppConfig.String("DbName")
	dbUser := beego.AppConfig.String("DbUser")
	dbPassword := beego.AppConfig.String("DbPassword")

	dbUrl := dbUser + ":" + dbPassword + "@tcp(" + dbIp + ":" + dbPort + ")/" + dbName + "?charset=utf8&loc=Asia%2FShanghai"
	beego.Debug("dbUrl=", dbUrl)

	err = orm.RegisterDataBase("default", "mysql", dbUrl)
	if err != nil {
		beego.Critical(err)
	}
}

func CheckAdmin(admin *Admininfo) bool {
	o := orm.NewOrm()
	exist := o.QueryTable("admininfo").Filter("UserName", admin.Username).Filter("Password", admin.Password).Exist()
	return exist
}

func CheckAccount(user *Userinfo) bool {
	o := orm.NewOrm()
	exist := o.QueryTable("userinfo").Filter("UserName", user.Username).Filter("Password", user.Password).Exist()
	return exist
}

func UpdateAdminStatus(admin *Admininfo) bool {
	o := orm.NewOrm()
	var u Admininfo
	err := o.QueryTable("admininfo").Filter("UserName", admin.Username).One(&u)
	if err != nil {
		beego.Error(err)
		return false
	}
	admin.Uid = u.Uid
	admin.CreatedTime = u.CreatedTime
	admin.LastLoginTime = time.Now()
	_, err = o.Update(admin)
	if err != nil {
		beego.Error(err)
		return false
	}
	return true
}

func ImportDeviceCheck(device *Deviceinfo) bool {
	o := orm.NewOrm()
	exist := o.QueryTable("deviceinfo").Filter("mac", device.Mac).Exist()
	return exist
}

func ImportDevices(devices *[]Deviceinfo) bool {
	o := orm.NewOrm()
	//successNum, err := o.InsertMulti(100, devices)
	_, err := o.InsertMulti(100, devices)
	if err != nil {
		beego.Error(err)
		return false
	}
	//fmt.Println("successNum=", successNum)
	return true
}

func RegisterDeivce(deviceinfo *Deviceinfo) (string, string, int64) {
	act := GetAct(deviceinfo)
	if act != "" {
		return act, "", -3
	}

	o := orm.NewOrm()
	guid := GetGuid()

	var device Deviceinfo
	err := o.QueryTable("deviceinfo").Filter("mac", deviceinfo.Mac).One(&device)
	if err != nil {
		//device doesnot  exsited, return error info
		beego.Error(err)
		return act, "", -2
	} else {
		//check if guid exist
		exist := o.QueryTable("deviceinfo").Filter("mac", device.Mac).Filter("registration_time__isnull", false).Exist()
		if exist {
			return act, "", -3
		}
		//check device register time is over time
		if device.InvalidTime.Before(time.Now()) {
			return act, "", -4
		}
		//device exsited,  insert info
		device.Guid = guid
		device.Mid = deviceinfo.Mid
		device.Psn = deviceinfo.Psn
		device.RegistrationTime = time.Now()
		_, err := o.Update(&device)
		if err != nil {
			beego.Error(err)
			return act, "", -5
		}
	}
	return act, guid, 0
}

func VerifyDevice(deviceinfo *Deviceinfo) (string, int64) {
	act := GetAct(deviceinfo)
	if act != "" {
		return act, -3
	}

	o := orm.NewOrm()
	exist := o.QueryTable("deviceinfo").Filter("mac", deviceinfo.Mac).Filter("guid", deviceinfo.Guid).Filter("mid", deviceinfo.Mid).Filter("psn", deviceinfo.Psn).Exist()
	if !exist {
		return act, -3 //mac/guid/mid/psn not match
	}
	return act, 0 //verify success
}

func ModifyDevice(deviceinfo *Deviceinfo, newmac string) bool {
	o := orm.NewOrm()
	var device Deviceinfo
	err := o.QueryTable("deviceinfo").Filter("mac", deviceinfo.Mac).One(&device)
	if err != nil {
		beego.Error(err)
		return false
	}
	//delete old record
	_, err = o.Delete(&device)
	if err != nil {
		beego.Error(err)
		return false
	}
	//insert new mac record
	device.Mac = newmac
	device.Guid = ""
	device.UpdateTime = time.Now()
	_, err = o.Insert(&device)
	if err != nil {
		beego.Error(err)
		return false
	}
	return true
}

func ReportDeivce(deviceinfo *Deviceinfo) (string, int64) {
	act := GetAct(deviceinfo)
	if act != "" {
		return act, -3
	}

	o := orm.NewOrm()
	var device Deviceinfo
	err := o.QueryTable("deviceinfo").Filter("mac", deviceinfo.Mac).Filter("mid", deviceinfo.Mid).Filter("psn", deviceinfo.Psn).One(&device)
	if err != nil {
		beego.Error(err)
		return act, -3
	}
	device.ErrorTime = time.Now()
	device.Error = deviceinfo.Error
	_, err = o.Update(&device)
	if err != nil {
		beego.Error(err)
		return act, -3
	}
	return act, 0
}

func SetAct(deviceinfo *Deviceinfo) bool {
	o := orm.NewOrm()
	var device Deviceinfo
	err := o.QueryTable("deviceinfo").Filter("mac", deviceinfo.Mac).One(&device)
	if err != nil {
		beego.Error(err)
		return false
	}

	device.Act = deviceinfo.Act
	_, err = o.Update(&device)
	if err != nil {
		beego.Error(err)
		return false
	}
	return true
}

func GetAct(deviceinfo *Deviceinfo) string {
	o := orm.NewOrm()
	var device Deviceinfo
	err := o.QueryTable("deviceinfo").Filter("mac", deviceinfo.Mac).One(&device)
	if err != nil {
		beego.Error(err)
		return ""
	}
	return device.Act
}

func GetDevices(start, offset int64) ([]*Deviceinfo, int64, bool) {
	o := orm.NewOrm()
	//get all devices
	devices := make([]*Deviceinfo, 0)
	num, err := o.QueryTable("deviceinfo").Limit(start, offset).Filter("registration_time__isnull", false).All(&devices)
	if err != nil {
		beego.Error(err)
		return nil, 0, false
	}
	return devices, num, true
}

func GetNoregisterDevices(start, offset int64) ([]*Deviceinfo, int64, bool) {
	o := orm.NewOrm()
	//get all devices of noregister
	devices := make([]*Deviceinfo, 0)
	num, err := o.QueryTable("deviceinfo").Limit(start, offset).Filter("registration_time__isnull", true).All(&devices)
	if err != nil {
		beego.Error(err)
		return nil, 0, false
	}
	return devices, num, true
}

func GetReportDevices(start, offset int64) ([]*Deviceinfo, int64, bool) {
	o := orm.NewOrm()
	//get all devices of report
	devices := make([]*Deviceinfo, 0)
	num, err := o.QueryTable("deviceinfo").Limit(start, offset).Filter("error_time__isnull", false).All(&devices)
	if err != nil {
		beego.Error(err)
		return nil, 0, false
	}
	return devices, num, true
}

func GetDevicesCount() int64 {
	o := orm.NewOrm()
	cnt, _ := o.QueryTable("deviceinfo").Filter("registration_time__isnull", false).Count()
	return cnt
}

func GetNoregisterDevicesCount() int64 {
	o := orm.NewOrm()
	cnt, _ := o.QueryTable("deviceinfo").Filter("registration_time__isnull", true).Count()
	return cnt
}

func GetReportDevicesCount() int64 {
	o := orm.NewOrm()
	cnt, _ := o.QueryTable("deviceinfo").Filter("error_time__isnull", false).Count()
	return cnt
}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func GetGuid() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		beego.Error(err)
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}
