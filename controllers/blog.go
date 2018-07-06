package controllers

import (
	"fmt"
	"github.com/LeeReindeer/lightblog/models"
	"github.com/astaxie/beego"
	"log"
	"strconv"
)

type BlogController struct {
	beego.Controller
}

// GET: /blog/<id>
func (this *BlogController) DetailBlog() {
	log.Println("blog id: ", this.Ctx.Input.Param(":id"))
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

// GET: /blog/like?id=<id>&page=1
// GET: /blog/like?id=<id>
func (this *BlogController) LikeBlog() {
	uid, err := strconv.Atoi(this.Ctx.GetCookie("uid"))
	page, err := this.GetInt("page")

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
	if page != 0 {
		// page start from 1
		this.Redirect(fmt.Sprintf("/?page=%d", page), 302)
	} else {
		this.Redirect("/", 302)
	}
}

// GET: /blog/dislike?id=<id>
func (this *BlogController) DisLikeBlog() {
	uid, err := strconv.Atoi(this.Ctx.GetCookie("uid"))
	page, err := this.GetInt("page")

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
	if page != 0 {
		// page start from 1
		this.Redirect(fmt.Sprintf("/?page=%d", page), 302)
	} else {
		this.Redirect("/", 302)
	}
}
