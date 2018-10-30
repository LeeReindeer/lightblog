package controllers

import (
	"encoding/base64"
	"log"
	"strconv"
	"time"

	"github.com/LeeReindeer/lightblog/models"
	"github.com/LeeReindeer/lightblog/util"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

type IndexController struct {
	beego.Controller
}

// GET: /
func (this *IndexController) TimeLine() {
	uid, login := CheckLogin(this.Ctx)
	log.Println("uid: ", uid)
	if !login {
		log.Println("you're not login")
		util.ClearCookies(this.Ctx)
		this.Redirect("/login", 302)
		return
	}
	var blogs []models.LightBlog
	curPage, err := this.GetInt("page")
	if err != nil {
		curPage = 1
		blogs = models.GetTimeLineByUid(uid)
	} else {
		blogs = models.GetTimeLineByUidWithPaging(uid, curPage)
	}

	this.Data["blogs"] = blogs
	this.Data["thisUser"] = models.GetUserById(uid)
	this.Data["redirect"] = this.Ctx.Input.URL()
	paging(uid, curPage, this)

	this.Layout = "layout.html"
	this.TplName = "index.html"
}

func paging(uid int64, curPage int, that *IndexController) {
	cnt := models.GetTimeLineCount(uid)
	log.Println("sum of blogs:", cnt)
	pages := int(cnt / 20)
	if cnt%20 != 0 {
		pages += 1
	}
	log.Println("pages: ", pages)
	pages_ := make([]int, pages)
	for i := 0; i < pages; i++ {
		pages_[i] = i + 1
	}
	log.Println(pages_)
	that.Data["pages"] = pages_
	if curPage < pages {

	}
	if curPage < pages {
		that.Data["next"] = curPage + 1
	}
	if curPage > 1 {
		that.Data["prev"] = curPage - 1
	}
}

func CheckLogin(ctx *context.Context) (int64, bool) {
	nameEncoded := ctx.GetCookie("username")
	nameByte, err := base64.StdEncoding.DecodeString(nameEncoded)
	if err != nil {
		return 0, false
	}
	name := string(nameByte)
	log.Println("name from cookie: ", name)
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
	// content can not be empty
	if err != nil || content == "" {
		this.Redirect("/", 302)
		return
	}
	blog := models.Blog{BlogUid: int64(uid), BlogContent: content, BlogTime: time.Now()}
	tagName, hasTag, contentIndex := util.GetContentTag(content)
	if hasTag {
		log.Println("blog has tag")
		tagId, _ := models.SaveTag(&models.Tag{TagName: tagName, TagTime: time.Now()})
		blog.BlogTagId = tagId
		blog.BlogContent = content[contentIndex:]
	}
	models.SaveBlog(&blog)
	this.Redirect("/", 302)
}
