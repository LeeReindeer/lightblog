package models

import (
	"github.com/LeeReindeer/lightblog/util"
	"github.com/astaxie/beego/orm"
	"log"
)

func GetAllComments(blogId int64) (comments []Comment) {
	o := orm.NewOrm()
	_, err := o.QueryTable("comment").Filter("comm_blog_id", blogId).All(&comments)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return
}

// not support emoji, and blog also.
func SaveComment(comment *Comment) (commId int64, ok bool) {
	o := orm.NewOrm()
	commId, err := o.Insert(comment)
	if err != nil {
		log.Println(err.Error())
		return 0, false
	}
	log.Println("new comment id: ", commId)
	ok = true
	return
}

func DeleteComment(comment *Comment) (commId int64) {
	o := orm.NewOrm()
	commId, err := o.Delete(comment)
	util.CheckDBErr(err)
	log.Println("new comment id: ", commId)
	return
}

func LikeComment(commId int64) {
	if commId == 0 {
		log.Println("error: comm_id = 0")
		return
	}
	db, err := orm.GetDB()
	_, err = db.Exec("UPDATE comment SET comm_like=comm_like+1 WHERE comm_id=?", commId)
	util.CheckDBErr(err)
}
