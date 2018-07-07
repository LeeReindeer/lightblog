package controllers

import (
	"github.com/LeeReindeer/lightblog/models"
	"github.com/LeeReindeer/lightblog/util"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"log"
	"strconv"
	"time"
)

type IndexController struct {
	beego.Controller
}

// GET: /
func (this *IndexController) TimeLine() {
	uid, login := CheckLogin(this.Ctx)
	log.Println("uid: ", uid)
	if !login {
		util.ClearCookies(this.Ctx)
		this.Redirect("/login", 302)
		return
	}

	blogs := models.GetTimeLineByUid(uid)
	log.Println("number of blogs: ", len(blogs))
	this.Data["blogs"] = models.GetTimeLineByUid(uid)
	this.Data["user"] = models.GetUserById(uid)
	this.Data["redirect"] = this.Ctx.Input.URL()

	this.Layout = "layout.html"
	this.TplName = "index.html"
}

func CheckLogin(ctx *context.Context) (int64, bool) {
	name := ctx.GetCookie("username")
	passHash := models.GetPassHashByName(name)
	passHashFromCookie, ok := ctx.GetSecureCookie(passHash, "p")

	// check username and hash
	if name == "" || !ok || passHashFromCookie != passHash {
		log.Println("not login: cookie error")
		return 0, false
	}
	uid := models.GetUserByName(name).UserId
	if uid == 0 {
		log.Println("not login: error id")
		return 0, false
	}
	// check uid
	uidFromCookie, ok := util.GetUserIdFromCookie(ctx)
	if !ok || uidFromCookie != uid {
		return 0, false
	}
	return uid, true
}

// POST: /
func (this *IndexController) NewLight() {
	log.Println("new blog")
	content := this.GetString("content")
	uid, err := strconv.Atoi(this.Ctx.GetCookie("uid"))
	if err != nil {
		this.Redirect("/", 302)
		return
	}
	blog := models.Blog{BlogUid: int64(uid), BlogContent: content, BlogTime: time.Now()}
	models.SaveBlog(&blog)
	this.Redirect("/", 302)
}
