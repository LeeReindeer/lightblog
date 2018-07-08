package models

import (
	"github.com/LeeReindeer/lightblog/util"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

func GetTimeLineByUid(uid int64) (blogs []LightBlog) {
	return GetTimeLineByUidWithPaging(uid, 1)
}

func getBlogPreview(content string) (preview string) {
	contentRune := []rune(content)
	if len(contentRune) > 140 {
		preview = string(contentRune[:140]) + "..."
	} else {
		preview = content
	}
	return
}

// every page has 20 blog
// page start from 1
func GetTimeLineByUidWithPaging(uid int64, page int) (blogs []LightBlog) {
	o := orm.NewOrm()
	_, err := o.Raw("SELECT * FROM blog WHERE blog_uid IN "+
		"(SELECT user_to FROM fan_follow WHERE user_from = ?) OR blog_uid= ? ORDER BY blog_time DESC LIMIT ?, ?", uid, uid, (page-1)*20, 20).QueryRows(&blogs)
	util.CheckDBErr(err)

	userMap := make(map[int64]*User)
	for i, _ := range blogs {
		uid := blogs[i].BlogUid
		if userMap[uid] == nil {
			user := User{UserId: uid}
			// only read a user once
			err = o.Read(&user)
			util.CheckDBErr(err)
			userMap[user.UserId] = &user
		}
		blogs[i].BlogPreview = getBlogPreview(blogs[i].BlogContent)
		blogs[i].BlogTimeString = blogs[i].BlogTime.Format("2006-01-02 15:04:05")
		blogs[i].BlogUsername = userMap[uid].UserName
		blogs[i].BlogUserAvatar = userMap[uid].UserAvatar
	}
	return
}

func GetTimeLineCount(uid int64) (count int64) {
	o := orm.NewOrm()
	var blogs []Blog
	nums, err := o.Raw("SELECT * FROM blog WHERE blog_uid IN "+
		"(SELECT user_to FROM fan_follow WHERE user_from = ?) OR blog_uid= ?", uid, uid).QueryRows(&blogs)
	if err != nil {
		log.Println(err)
		return 0
	}
	return nums
}

func GetBlogById(id int64) (lightblog LightBlog) {
	o := orm.NewOrm()
	blog := Blog{BlogId: id}
	err := o.Read(&blog)
	util.CheckDBErr(err)
	lightblog.Blog = blog

	user := User{UserId: blog.BlogUid}
	err = o.Read(&user)
	util.CheckDBErr(err)

	// fill
	lightblog.BlogPreview = getBlogPreview(blog.BlogContent)
	lightblog.BlogTimeString = blog.BlogTime.Format("2006-01-02 15:04:05")
	lightblog.BlogUsername = user.UserName
	lightblog.BlogUserAvatar = user.UserAvatar
	return
}

func GetBlogsByUid(uid int64) (blogs []LightBlog) {
	o := orm.NewOrm()
	_, err := o.Raw("SELECT * FROM blog WHERE blog_uid=? ORDER BY blog_time DESC", uid).QueryRows(&blogs)
	if err != nil {
		return nil
	}
	user := User{UserId: uid}
	// only read a user once
	err = o.Read(&user)
	util.CheckDBErr(err)
	for i, _ := range blogs {
		blogs[i].BlogPreview = getBlogPreview(blogs[i].BlogContent)
		blogs[i].BlogTimeString = blogs[i].BlogTime.Format("2006-01-02 15:04:05")
		blogs[i].BlogUsername = user.UserName
		blogs[i].BlogUserAvatar = user.UserAvatar
	}
	return
}

func SaveBlog(blog *Blog) int64 {
	o := orm.NewOrm()
	id, err := o.Insert(blog)
	log.Println("new blog id: ", id)
	util.CheckDBErr(err)
	return id
}

// blog in param must have BlogId, BlogContent
// so, this update BlogContent and BlogTime
func UpdateBlog(blog *Blog) {
	o := orm.NewOrm()
	blog.BlogTime = time.Now()
	_, err := o.Update(blog, "BlogContent", "BlogTime")
	util.CheckDBErr(err)
}

func DeleteBlog(blogId int64) int64 {
	o := orm.NewOrm()
	blog := Blog{BlogId: blogId}
	id, err := o.Delete(&blog)
	log.Println("delete blog id: ", id)
	util.CheckDBErr(err)
	return id
}

/** counter func start**/
func IncLikeBlog(blogId, uid int64) {
	db, err := orm.GetDB()
	_, err = db.Exec("UPDATE blog SET blog_like=blog_like+1 WHERE blog_id=?", blogId)
	if err != nil {
		log.Println(err.Error())
		return
	}

	_, err = db.Exec("INSERT INTO blog_counter VALUES(?, ?, ?)", blogId, uid, LIKE_COUNT)
	util.CheckDBErr(err)
	//_, err := o.Raw("UPDATE blog SET blog_like=blog_like+1 WHERE blog_id=?", blogId).Exec()
	//util.CheckDBErr(err)
}

func DecLikeBlog(blogId int64, uid int64) {
	db, err := orm.GetDB()
	util.CheckDBErr(err)
	row := db.QueryRow("SELECT blog_like FROM blog WHERE blog_id=? ", blogId)
	util.CheckDBErr(err)

	like := 0
	err = row.Scan(&like)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("like count: ", like)
	if like > 0 {
		_, err = db.Exec("UPDATE blog SET blog_like=blog_like-1 WHERE blog_id=?", blogId)
		if err != nil {
			log.Println(err.Error())
			return
		}
		_, err = db.Exec("DELETE FROM blog_counter WHERE blog_id=? AND user_id=? AND count_type=?", blogId, uid, LIKE_COUNT)
		util.CheckDBErr(err)
	} else {
		log.Println("like counter < 0")
	}
}

// select count(*) from blog_counter where blog_id=<blogId> and user_id=<id> and count_type=0;
func IsUserLikeBlog(blogId, uid int64) bool {
	db, err := orm.GetDB()
	util.CheckDBErr(err)
	row := db.QueryRow("SELECT count(*) FROM blog_counter WHERE blog_id=? AND user_id=? AND count_type=0", blogId, uid)
	// count can only be zero or one
	count := 0
	row.Scan(&count)

	return count != 0
}

func IncUnlikeBlog(blogId, uid int64) {
	db, err := orm.GetDB()
	_, err = db.Exec("UPDATE blog SET blog_unlike=blog_unlike+1 WHERE blog_id=?", blogId)
	if err != nil {
		log.Println(err.Error())
		return
	}

	_, err = db.Exec("INSERT INTO blog_counter VALUES(?, ?, ?)", blogId, uid, UNLIKE_COUNT)
	util.CheckDBErr(err)
}

func DecUnlikeBlog(blogId, uid int64) {
	db, err := orm.GetDB()
	rows, err := db.Query("SELECT blog_unlike FROM blog WHERE blog_id=? ", blogId)
	util.CheckDBErr(err)
	if rows.Next() {
		unlike := 0
		rows.Scan(&unlike)
		if unlike > 0 {
			_, err = db.Exec("UPDATE blog SET blog_unlike=blog_unlike-1 WHERE blog_id=?", blogId)
			if err != nil {
				log.Println(err.Error())
				return
			}
			_, err = db.Exec("DELETE FROM blog_counter WHERE blog_id=? AND user_id=? AND count_type=?", blogId, uid, UNLIKE_COUNT)
			util.CheckDBErr(err)
		} else {
			log.Println("unlike counter < 0")
		}
	}
}

func IsUserDisLikeBlog(blogId, uid int64) bool {
	db, err := orm.GetDB()
	util.CheckDBErr(err)
	row := db.QueryRow("SELECT count(*) FROM blog_counter WHERE blog_id=? AND user_id=? AND count_type=1", blogId, uid)
	count := 0
	row.Scan(&count)

	return count != 0
}

func IncBlogComment(blogId int64) {
	db, err := orm.GetDB()
	_, err = db.Exec("UPDATE blog SET blog_comment=blog_comment+1 WHERE blog_id=?", blogId)
	util.CheckDBErr(err)
}

func DecBlogComment(blogId int64) {
	db, err := orm.GetDB()
	_, err = db.Exec("UPDATE blog SET blog_comment=blog_comment-1 WHERE blog_id=?", blogId)
	util.CheckDBErr(err)
}

/**counter func end**/
