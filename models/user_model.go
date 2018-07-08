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

// is fromId followed toId
func IsUserFollow(fromId, toId int64) bool {
	db, _ := orm.GetDB()
	row := db.QueryRow("SELECT user_to FROM fan_follow WHERE user_from=? AND user_to=?", fromId, toId)
	var userToId int64
	err := row.Scan(&userToId)
	if err != nil {
		return false
	}
	return userToId == toId
}

func FollowUser(fromId, toId int64) bool {
	db, _ := orm.GetDB()
	_, err := db.Exec("INSERT INTO fan_follow VALUES(?, ?)", fromId, toId)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	// incr follower of user:toId and incr following of user:fromId
	_, err = db.Exec("UPDATE user SET user_followers=user_followers+1 WHERE user_id=?", toId)
	_, err = db.Exec("UPDATE user SET user_following=user_following+1 WHERE user_id=?", fromId)
	if err != nil {
		return false
	}
	return true
}

func UnFollowUser(fromId, toId int64) bool {
	db, _ := orm.GetDB()
	_, err := db.Exec("DELETE FROM fan_follow WHERE user_from=? AND user_to=?", fromId, toId)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	// incr follower of user:toId and incr following of user:fromId
	_, err = db.Exec("UPDATE user SET user_followers=user_followers-1 WHERE user_id=?", toId)
	_, err = db.Exec("UPDATE user SET user_following=user_following-1 WHERE user_id=?", fromId)
	if err != nil {
		return false
	}
	return true
}
