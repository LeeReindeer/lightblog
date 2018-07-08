package controllers

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"github.com/LeeReindeer/lightblog/models"
	"github.com/LeeReindeer/lightblog/util"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"log"
	"strconv"
	"time"
)

type UserController struct {
	beego.Controller
}

// GET: /user/<username>, user detail
// follow user: /user/<username>?follow=true
// unfollow user: /user/<username>?unfollow=true
func (this *UserController) Get() {
	username := this.Ctx.Input.Params()[":username"]
	redirectURL := this.Ctx.Input.URL()

	loginUserId, ok := util.GetUserIdFromCookie(this.Ctx)
	if !ok {
		Logout(this.Ctx)
		return
	}
	thatUser := models.GetUserByName(username)
	thisUser := models.GetUserById(loginUserId)
	if this.GetString("follow") != "" {
		log.Println("handle follow")
		followUser(thisUser.UserId, thatUser.UserId)
		util.Redirect302(redirectURL, this.Ctx)
		return
	}

	if this.GetString("unfollow") != "" {
		log.Println("handle unfollow")
		unfollowUser(thisUser.UserId, thatUser.UserId)
		util.Redirect302(redirectURL, this.Ctx)
		return
	}
	this.Data["thisUser"] = thisUser
	this.Data["user"] = thatUser
	this.Data["redirect"] = redirectURL
	this.Data["blogs"] = models.GetBlogsByUid(thatUser.UserId)
	this.Data["followed"] = models.IsUserFollow(thisUser.UserId, thatUser.UserId)
	this.Layout = "layout.html"
	this.TplName = "profile.html"
}

func followUser(fromId, toId int64) {
	if fromId == toId {
		log.Println("you can't follow youself")
		return
	}
	models.FollowUser(fromId, toId)
}

func unfollowUser(fromId, toId int64) {
	if fromId == toId {
		log.Println("you can't unfollow youself")
		return
	}
	models.UnFollowUser(fromId, toId)
}

//abandon Post: /user/<username>, update user
func (this *UserController) Post() {
	//todo <input type="hidden" name="_method" value="PUT">
	// handel delete here
	//this.Delete()
}

//abandon DELETE: /user/<username>, delete user
func (this *UserController) Delete() {

}

// GET: /register
func (this *UserController) GetRegister() {
	beego.ReadFromRequest(&this.Controller) // parse flash
	this.Data["register"] = true
	this.TplName = "login.html"
}

// handel register: POST: /register
func (this *UserController) RegisterUser() {
	flash := beego.NewFlash()

	username := this.GetString("username")
	pass := this.GetString("password")
	passAgain := this.GetString("passwordAgain")
	if pass != passAgain {
		flash.Error("两次输入的密码不同")
		flash.Store(&this.Controller)
		log.Println("password not match")
		this.Ctx.Redirect(302, "/register")
		return
	}
	if len(pass) > 16 || len(pass) < 4 {
		flash.Error("密码太短")
		flash.Store(&this.Controller)
		log.Println("password too short")
		this.Ctx.Redirect(302, "/register")
		return
	}
	h := PasswordHash(pass)
	user := models.User{UserName: username, UserAvatar: models.DefaultAvatar, UserPassword: h,
		UserRegister: time.Now(), UserBio: "", UserFollowers: 0, UserFollowing: 0}

	id, err := models.SaveUser(user)
	if err != nil {
		//register failed
		log.Println("register failed")
		this.Ctx.Redirect(302, "/register")
	}
	log.Println("register succeed. user id: ", id)
	flash.Warning("请登录吧")
	flash.Store(&this.Controller)
	this.Ctx.Redirect(302, "/login")
}

func PasswordHash(password string) (sig string) {
	passHash := hmac.New(sha1.New, []byte(password))
	sig = fmt.Sprintf("%02x", passHash.Sum(nil))
	return
}

// GET : /login
func (this *UserController) GetLogin() {
	_, ok := CheckLogin(this.Ctx)
	if ok {
		this.Redirect("/", 302)
		return
	}
	util.ClearCookies(this.Ctx)
	flash := beego.ReadFromRequest(&this.Controller)
	if _, ok := flash.Data["error"]; ok {
		log.Println("login: show error flash")
	}
	this.Data["login"] = true
	this.TplName = "login.html"
}

// handel login: POST: /login
func (this *UserController) LoginUser() {
	flash := beego.NewFlash()

	username := this.GetString("username")
	pass := this.GetString("password")
	h := PasswordHash(pass)
	log.Println("username: ", username)
	log.Println("password: ", pass)
	if h == models.GetPassHashByName(username) {
		this.Ctx.SetCookie("username", username)
		uid := models.GetUserByName(username).UserId
		this.Ctx.SetCookie("uid", strconv.FormatInt(uid, 10))
		this.Ctx.SetSecureCookie(h, "p", h)
		this.Ctx.SetCookie("login", "true")
		log.Println("login success")
		this.Ctx.Redirect(302, "/")
	} else {
		flash.Error("密码错误")
		flash.Store(&this.Controller)
		log.Println("login failed: password error")
		this.Ctx.Redirect(302, "/login")
	}
}

func Logout(ctx *context.Context) {
	util.ClearCookies(ctx)
	ctx.Redirect(302, "/login")
	log.Println("logout succeed")
}

func (this *UserController) LogoutUser() {
	Logout(this.Ctx)
}
