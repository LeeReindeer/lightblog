package models

import (
	"github.com/LeeReindeer/lightblog/util"
	"github.com/astaxie/beego/orm"
	"log"
)

func GetUserByName(name string) *User {
	o := orm.NewOrm()

	var user User
	user.UserName = name
	err := o.Raw("SELECT * FROM user WHERE user_name=?", name).QueryRow(&user)
	if err != nil {
		panic(err)
		return nil
	}
	return &user
}

func GetUserById(id int64) *User {
	o := orm.NewOrm()

	user := User{UserId: id}
	err := o.Raw("SELECT * FROM user WHERE user_id=?", id).QueryRow(&user)
	util.CheckDBErr(err)
	return &user
}

func GetPassHashByName(name string) string {
	o := orm.NewOrm()

	user := User{UserName: name}
	err := o.Raw("SELECT * FROM user WHERE user_name=?", name).QueryRow(&user)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	return user.UserPassword
}

func SaveUser(user User) (uid int64, err error) {
	o := orm.NewOrm()
	uid, err = o.Insert(&user)
	util.CheckAndFatal(err)
	return
}
