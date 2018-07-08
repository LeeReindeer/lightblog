package controllers

import (
	"github.com/LeeReindeer/lightblog/models"
	"github.com/LeeReindeer/lightblog/util"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"log"
	"strconv"
)

type BlogController struct {
	beego.Controller
}

func getBlogIdFromUrl(ctx *context.Context) (int64, bool) {
	blogIdStr := ctx.Input.Param(":id")
	if len(blogIdStr) == 0 {
		return 0, false
	}
	blogId, err := strconv.Atoi(blogIdStr)
	if err != nil {
		log.Println(err.Error())
		return 0, false
	}
	return int64(blogId), true
}

// GET: /blog/<id>, delete a blog : GET: /blog/<id>?delete=true
func (this *BlogController) DetailBlog() {
	log.Println("URL: ", this.Ctx.Input.URL())
	blogId, ok := getBlogIdFromUrl(this.Ctx)
	if !ok || blogId == 0 {
		this.Redirect("/", 302)
		return
	}
	log.Println("blog id: ", blogId)
	lightBlog := models.GetBlogById(blogId)

	if this.GetString("delete") != "" {
		// check user login
		if idFromC, _ := util.GetUserIdFromCookie(this.Ctx); idFromC == lightBlog.BlogUid {
			models.DeleteBlog(blogId)
		}
		this.Redirect("/", 302)
		return
	}

	this.Data["blog"] = lightBlog

	comments := models.GetAllComments(blogId)
	this.Data["comments"] = comments

	uid, ok := util.GetUserIdFromCookie(this.Ctx)
	if !ok {
		util.ClearCookies(this.Ctx)
		this.Redirect("/login", 302)
	}
	this.Data["user"] = models.GetUserById(uid)
	this.Data["redirect"] = this.Ctx.Input.URL()

	//log.Println("test: ", this.GetString("id"))
	this.Layout = "layout.html"
	this.TplName = "detail.html"
}

// POST: /blog/<id>
func (this *BlogController) EditBlog() {

}

// DELETE: /blog/<id>
func (this *BlogController) DeleteBlog() {

}

// GET: /blog/like?id=<id>&redirect=<url>
// GET: /blog/like?id=<id>
func (this *BlogController) LikeBlog() {
	uid, _ := strconv.Atoi(this.Ctx.GetCookie("uid"))
	blogIdStr := this.GetString("id")
	blogId, err := strconv.Atoi(blogIdStr)
	if err != nil {
		log.Println(err.Error())
		return
	}

	liked := models.IsUserLikeBlog(int64(blogId), int64(uid))
	log.Printf("user%d liked blog%d? %v", uid, blogId, liked)
	if liked {
		models.DecLikeBlog(int64(blogId), int64(uid))
	} else {
		models.IncLikeBlog(int64(blogId), int64(uid))
	}

	this.Redirect(this.GetString("redirect"), 302)
}

// GET: /blog/dislike?id=<id>
func (this *BlogController) DisLikeBlog() {
	uid, _ := strconv.Atoi(this.Ctx.GetCookie("uid"))
	blogIdStr := this.GetString("id")
	blogId, err := strconv.Atoi(blogIdStr)
	if err != nil {
		log.Println(err.Error())
		return
	}

	disliked := models.IsUserDisLikeBlog(int64(blogId), int64(uid))
	log.Printf("user%d liked blog%d? %v", uid, blogId, disliked)
	if disliked {
		models.DecUnlikeBlog(int64(blogId), int64(uid))
	} else {
		models.IncUnlikeBlog(int64(blogId), int64(uid))
	}

	this.Redirect(this.GetString("redirect"), 302)
}
