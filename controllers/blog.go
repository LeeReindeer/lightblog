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

// get blogId from "/blog/<id>"
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

// get blogId from "/blog?id=<id>"
func getBloIdFromParams(this *BlogController) (int64, bool) {
	blogIdStr := this.GetString("id")
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

// edit blog, GET: /blog/edit?id=<id>&redirect=<url>
func (this *BlogController) GetEditBlog() {
	this.Layout = "layout.html"
	this.TplName = "edit_blog.html"
	blogId, ok := getBloIdFromParams(this)
	if !ok {
		util.Redirect302(this.GetString("redirect"), this.Ctx)
		return
	}
	lightBlog := models.GetBlogById(blogId)
	this.Data["blog"] = lightBlog
}

// edit blog POST: /blog/edit?id=<id>&redirect=
func (this *BlogController) EditBlog() {
	// get from form
	content := this.GetString("content")
	if content == "" {
		util.Redirect302(this.GetString("redirect"), this.Ctx)
		return
	}
	id, ok := getBloIdFromParams(this)
	if ok {
		blog := models.Blog{BlogId: id, BlogContent: content}
		models.UpdateBlog(&blog)
	}
	util.Redirect302(this.GetString("redirect"), this.Ctx)
}

// GET: /blog/like?id=<id>&redirect=<url>
// GET: /blog/like?id=<id>
func (this *BlogController) LikeBlog() {
	uid, ok := util.GetUserIdFromCookie(this.Ctx)
	blogId, liked := int64(0), false
	if !ok {
		goto error
	}
	blogId, ok = getBloIdFromParams(this)
	if !ok {
		goto error
	}

	liked = models.IsUserLikeBlog(blogId, uid)
	log.Printf("user%d liked blog%d? %v", uid, blogId, liked)
	if liked {
		models.DecLikeBlog(blogId, uid)
	} else {
		models.IncLikeBlog(blogId, uid)
	}
error:
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
