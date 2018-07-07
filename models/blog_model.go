package models

import (
	"github.com/LeeReindeer/lightblog/util"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

//var DB *sql.DB

//func init() {
//	DB = OpenDB()
//	if DB == nil {
//		log.Fatalln("can't open database")
//	}
//}

//func OpenDB() *sql.DB {
//	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/lightblog?charset=utf8&parseTime=true",
//		beego.AppConfig.String("mysqluser"),
//		beego.AppConfig.String("mysqlpass")))
//	util.CheckDBErr(err)
//	if db.Ping() == nil {
//		log.Println("Connected to mysql")
//	} else {
//		CloseDB(db)
//	}
//	return db
//}
//
//func CloseDB(db *sql.DB) {
//	db.Close()
//}

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

func SaveBlog(blog *Blog) int64 {
	o := orm.NewOrm()
	id, err := o.Insert(&blog)
	log.Println("new blog id: ", id)
	util.CheckDBErr(err)
	return id
}

/** counter func start**/
func IncLikeBlog(blogId, uid int64) {
	db, err := orm.GetDB()
	_, err = db.Exec("UPDATE blog SET blog_like=blog_like+1 WHERE blog_id=?", blogId)
	util.CheckDBErr(err)

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
		util.CheckDBErr(err)
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
	util.CheckDBErr(err)

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
			util.CheckDBErr(err)
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

//
//func UpdateBlog(blog Blog) {
//	stmt, err := DB.Prepare("UPDATE blog SET blog_content=? WHERE blog_id=?")
//	util.CheckDBErr(err)
//
//	_, err = stmt.Exec(blog.BlogContent, blog.BlogId)
//	util.CheckDBErr(err)
//}
//
//func DeleteBlogById(id int) {
//	stmt, err := DB.Prepare("DELETE FROM blog WHERE blog_id=?")
//	util.CheckDBErr(err)
//
//	_, err = stmt.Exec(id)
//	util.CheckDBErr(err)
//}
