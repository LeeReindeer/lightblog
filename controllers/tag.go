package controllers

import (
	"github.com/LeeReindeer/lightblog/models"
	"github.com/LeeReindeer/lightblog/util"

	"github.com/astaxie/beego"
)

const TagPrefix string = "#"

type TagController struct {
	beego.Controller
}

//todo GET:/tag/
func (this *TagController) ShowAllTags() {

}

// GET:/tag/<tag_id>/
func (this *TagController) ShowBlogWithTag() {
	uid, ok := util.GetUserIdFromCookie(this.Ctx)
	if !ok {
		Logout(this.Ctx)
		return
	}
	tagId, ok := util.GetIdFromUrl(this.Ctx, "tag_id")
	if !ok {
		this.Redirect("/", 302)
		return
	}
	tag := models.GetTagById(tagId)
	blogs := models.GetBlogsWithTag(tagId)

	//modify title
	this.Data["title"] = TagPrefix + tag.TagName

	this.Data["tag"] = tag
	this.Data["thisUser"] = models.GetUserById(uid)
	this.Data["blogs"] = blogs
	this.Layout = "layout.html"
	this.TplName = "tag_blog.html"
}
