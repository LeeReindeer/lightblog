package tests

import (
	"fmt"
	"github.com/LeeReindeer/lightblog/models"
	"github.com/LeeReindeer/lightblog/util"
	"github.com/astaxie/beego/orm"
	"testing"
)

func Test(t *testing.T) {
	orm.Debug = true

	err := orm.RegisterDriver("mysql", orm.DRMySQL)
	util.CheckAndFatal(err)
	err = orm.RegisterDataBase("default", "mysql",
		fmt.Sprintf("%s:%s@/lightblog?charset=utf8&parseTime=true",
			"leeR", "0915"))
	util.CheckAndFatal(err)

	o := orm.NewOrm()
	leer := models.User{UserId: 1}
	kwok := models.User{UserId: 2}

	err = o.Read(&leer)
	util.CheckDBErr(err)
	err = o.Read(&kwok)
	util.CheckDBErr(err)

	fmt.Println(leer)
	fmt.Println(kwok)

	// test blog
	var blog models.Blog
	//err = o.QueryTable("blog").Filter("blog_uid", leer.UserId).One(&blog)
	err = o.Raw("SELECT * FROM blog WHERE blog_id = ?", leer.UserId).QueryRow(&blog)
	util.CheckDBErr(err)
	fmt.Println(blog)

	// test comments
	var comment models.Comment
	err = o.Raw("SELECT * FROM comment WHERE comm_host_id = ?", blog.BlogId).QueryRow(&comment)
	util.CheckDBErr(err)
	fmt.Println(comment)

	// test TL
	var blogs []models.Blog
	num, err := o.Raw("SELECT * FROM blog WHERE blog_uid IN "+
		"(SELECT user_to FROM fan_follow WHERE user_from = ? ORDER BY blog_time DESC)", leer.UserId).QueryRows(&blogs)
	if err == nil {
		fmt.Println("timeline nums: ", num)
		fmt.Println(blogs)
	} else {
		panic(err)
	}
}
