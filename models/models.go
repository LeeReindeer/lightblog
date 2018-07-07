package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

// user for orm register
type Blog struct {
	BlogId      int64 `orm:"pk"`
	BlogUid     int64
	BlogContent string
	BlogTime    time.Time

	BlogLike    int
	BlogUnlike  int
	BlogComment int
}

const (
	LIKE_COUNT   = iota // 0
	UNLIKE_COUNT = iota // 1
)

type LightBlog struct {
	Blog
	//extra for convenience
	BlogPreview    string //博客简述，取前140个字符
	BlogUsername   string
	BlogUserAvatar string
	BlogTimeString string
}

func (blog Blog) String() string {
	return fmt.Sprintf("[%d] creator:%d, created:%s, content:%s",
		blog.BlogId, blog.BlogUid, blog.BlogTime.Format("2006-01-02 15:04:05"), blog.BlogContent)
}

const DefaultAvatar string = "https://avatars2.githubusercontent.com/u/24387694?s=100&v=4"

type User struct {
	UserId        int64 `orm:"pk"`
	UserName      string
	UserAvatar    string
	UserPassword  string
	UserRegister  time.Time
	UserBio       string
	UserFollowers int
	UserFollowing int
}

//see time/format.go
func (user User) String() string {
	return fmt.Sprintf("[%d %s] bio:%s, register:%s, followers:%d, following:%d",
		user.UserId, user.UserName, user.UserBio, user.UserRegister.Format("2006-01-02 15:04:05"),
		user.UserFollowers, user.UserFollowing)
}

type Comment struct {
	CommId      int64 `orm:"pk"`
	CommBlogId  int64
	CommFromUid int64
	CommToUid   int64
	CommContent string
	CommTime    time.Time
	CommLike    int

	// 冗余信息
	CommFromName   string
	CommFromAvatar string
	CommToName     string
	CommToAvatar   string
}

type Tag struct {
	TagId    int64 `orm:"pk"`
	TagName  string
	TageTime time.Time
}

func init() {
	orm.RegisterModel(new(Blog), new(User), new(Comment), new(Tag))
}
