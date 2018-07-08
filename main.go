package main

import (
	"fmt"
	"github.com/LeeReindeer/lightblog/controllers"
	"github.com/LeeReindeer/lightblog/util"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func init() {
	// todo remove when release
	orm.Debug = true

	err := orm.RegisterDriver("mysql", orm.DRMySQL)
	util.CheckAndFatal(err)
	err = orm.RegisterDataBase("default", "mysql",
		fmt.Sprintf("%s:%s@/lightblog?charset=utf8&parseTime=true",
			beego.AppConfig.String("mysqluser"),
			beego.AppConfig.String("mysqlpass")))
	util.CheckAndFatal(err)

	// GET: TimeLine, POST: new blog
	beego.Router("/", &controllers.IndexController{}, "get:TimeLine;post:NewLight")

	// user process
	beego.Router("/login", &controllers.UserController{}, "get:GetLogin;post:LoginUser")
	beego.Router("/register", &controllers.UserController{}, "get:GetRegister;post:RegisterUser")
	beego.Router("/logout", &controllers.UserController{}, "get:LogoutUser")

	// user detail, update user and delete user
	beego.Router("/user/:username([\\w]+)", &controllers.UserController{})

	// view blog detail with comments, edit blog and delete blog
	beego.Router("/blog/:id:int", &controllers.BlogController{}, "get:DetailBlog")
	beego.Router("blog/edit", &controllers.BlogController{}, "post:EditBlog")
	// like or dislike blog: GET: /blog/like?id=1
	beego.Router("/blog/like", &controllers.BlogController{}, "get:LikeBlog")
	beego.Router("/blog/dislike", &controllers.BlogController{}, "get:DisLikeBlog")

	// new comment(/comment) and delete comment: GET: /comment?comm_id=<id>
	beego.Router("/comment", &controllers.CommentController{},
		"post:NewComment;get:DeleteComment")
	//like a comment : GET: /comment/like?comm_id=
	beego.Router("comment/like", &controllers.CommentController{}, "get:LikeComment")

	beego.Router("/tag", &controllers.TagController{},
		"get:ShowAllTags")
	beego.Router("/tag/:tag_id:int", &controllers.TagController{},
		"get:ShowBlogWithTag")
}

func main() {
	beego.Run()
}
