package controllers

import (
	"log"
	"strconv"
	"time"

	"github.com/LeeReindeer/lightblog/models"
	"github.com/LeeReindeer/lightblog/util"

	"github.com/astaxie/beego"
)

type CommentController struct {
	beego.Controller
}

// GET: comment?comm_id=<id>&delete=true&redirect=<pre url>
func (this *CommentController) DeleteComment() {
	deleteComm := this.GetString("delete")
	commId, ok := util.StringToInt64(this.GetString("comm_id"))
	comment := models.GetCommentById(commId)
	if deleteComm != "true" || !ok {
		goto error
	}
	if _, ok = models.DeleteComment(&comment); ok {
		models.DecBlogComment(comment.CommBlogId)
	}
error:
	util.Redirect302(this.GetString("redirect"), this.Ctx)
}

// POST: /comment
func (this *CommentController) NewComment() {
	blogId, err := this.GetInt64("blogId")
	if err != nil {
		log.Println(err.Error())
		this.Redirect("/", 302)
		return
	}
	fromUserId, err := this.GetInt64("fromUserId")
	log.Println("comment from user:", fromUserId)
	if err != nil {
		log.Println(err.Error())
		this.Redirect("/", 302)
		return
	}

	var fromUser models.User
	var toUser models.User
	var toUserId = 0
	var comment *models.Comment
	content := this.GetString("commentContent")
	toUserIdStr := this.GetString("toUserId")

	fromUser = *models.GetUserById(fromUserId)
	if len(toUserIdStr) != 0 {
		//reply to comment
		toUserId, _ = strconv.Atoi(toUserIdStr)
		log.Println("reply to user:", toUserId)
		toUser = *models.GetUserById(int64(toUserId))
	}

	comment = &models.Comment{CommBlogId: blogId, CommFromUid: fromUserId, CommToUid: int64(toUserId), CommContent: content, CommTime: time.Now(),
		CommFromName: fromUser.UserName, CommFromAvatar: fromUser.UserAvatar, CommToName: toUser.UserName, CommToAvatar: toUser.UserAvatar}
	if _, ok := models.SaveComment(comment); ok {
		models.IncBlogComment(blogId)
	}
	this.Redirect(this.GetString("redirect"), 302)
}

// GET" /comment/like?comm_id=<id>
func (this *CommentController) LikeComment() {
	redirect := this.GetString("redirect")
	commId, ok := util.StringToInt64(this.GetString("comm_id"))
	if !ok {
		this.Redirect(redirect, 302)
		return
	}
	models.LikeComment(commId)
	this.Redirect(redirect, 302)
}
