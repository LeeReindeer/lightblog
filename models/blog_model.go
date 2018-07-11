package models

import (
	"log"
	"time"

	"github.com/LeeReindeer/lightblog/util"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func getBlogPreview(content string) (preview string) {
	contentRune := []rune(content)
	if len(contentRune) > 140 {
		preview = string(contentRune[:140]) + "..."
	} else {
		preview = content
	}
	return
}

func GetTimeLineByUid(uid int64) (blogs []LightBlog) {
	return GetTimeLineByUidWithPaging(uid, 1)
}

// every page has 20 blog
// page start from 1
func GetTimeLineByUidWithPaging(uid int64, page int) []LightBlog {
	db, _ := orm.GetDB()
	rows, _ := db.Query("SELECT * FROM blogdetail WHERE blog_uid IN "+
		"(SELECT user_to FROM fan_follow WHERE user_from = ?) OR blog_uid= ? ORDER BY blog_time DESC LIMIT ?, ?", uid, uid, (page-1)*20, 20)

	blogs := make([]LightBlog, 20)
	count := 0
	for ; rows.Next(); count++ {
		err := rows.Scan(&blogs[count].BlogId, &blogs[count].BlogUid, &blogs[count].BlogTagId, &blogs[count].BlogContent,
			&blogs[count].BlogTime, &blogs[count].BlogLike, &blogs[count].BlogUnlike, &blogs[count].BlogComment,
			&blogs[count].BlogUsername, &blogs[count].BlogUserAvatar)
		if err != nil {
			return nil
		}
	}
	blogs = blogs[:count]

	local, _ := time.LoadLocation("Asia/Shanghai")
	for i, _ := range blogs {
		if blogs[i].BlogTagId != 0 {
			blogs[i].Tag = *GetTagById(blogs[i].BlogTagId)
		}
		blogs[i].BlogPreview = getBlogPreview(blogs[i].BlogContent)
		blogs[i].BlogTimeString = blogs[i].BlogTime.In(local).Format("2006-01-02 15:04:05")
	}
	return blogs
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

func GetBlogDetailFromView(id int64) (lightblog LightBlog) {
	db, _ := orm.GetDB()
	row := db.QueryRow("select * from blogdetail where blog_id=?", id)

	err := row.Scan(&lightblog.BlogId, &lightblog.BlogUid, &lightblog.BlogTagId, &lightblog.BlogContent,
		&lightblog.BlogTime, &lightblog.BlogLike, &lightblog.BlogUnlike, &lightblog.BlogComment,
		&lightblog.BlogUsername, &lightblog.BlogUserAvatar)

	if err != nil {
		log.Println(err.Error())
	}
	local, _ := time.LoadLocation("Asia/Shanghai")
	lightblog.BlogPreview = getBlogPreview(lightblog.BlogContent)
	lightblog.BlogTimeString = lightblog.BlogTime.In(local).Format("2006-01-02 15:04:05")
	if lightblog.BlogTagId != 0 {
		lightblog.Tag = *GetTagById(lightblog.BlogTagId)
		log.Println("tag name: ", lightblog.TagName)
	}
	return
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
	if lightblog.BlogTagId != 0 {
		lightblog.Tag = *GetTagById(lightblog.BlogTagId)
		log.Println("tag name: ", lightblog.TagName)
	}
	lightblog.BlogPreview = getBlogPreview(blog.BlogContent)
	lightblog.BlogTimeString = blog.BlogTime.Format("2006-01-02 15:04:05")
	lightblog.BlogUsername = user.UserName
	lightblog.BlogUserAvatar = user.UserAvatar
	return
}

func GetBlogsByUid(uid int64) []LightBlog {
	db, _ := orm.GetDB()
	rows, _ := db.Query("SELECT * FROM blogdetail WHERE blog_uid=? ORDER BY blog_time DESC", uid)

	blogs := make([]LightBlog, 20)
	count := 0
	for ; rows.Next(); count++ {
		err := rows.Scan(&blogs[count].BlogId, &blogs[count].BlogUid, &blogs[count].BlogTagId, &blogs[count].BlogContent,
			&blogs[count].BlogTime, &blogs[count].BlogLike, &blogs[count].BlogUnlike, &blogs[count].BlogComment,
			&blogs[count].BlogUsername, &blogs[count].BlogUserAvatar)
		if err != nil {
			return nil
		}
	}
	blogs = blogs[:count]

	local, _ := time.LoadLocation("Asia/Shanghai")
	for i, _ := range blogs {
		if blogs[i].BlogTagId != 0 {
			blogs[i].Tag = *GetTagById(blogs[i].BlogTagId)
		}
		blogs[i].BlogPreview = getBlogPreview(blogs[i].BlogContent)
		blogs[i].BlogTimeString = blogs[i].BlogTime.In(local).Format("2006-01-02 15:04:05")
	}
	return blogs
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
