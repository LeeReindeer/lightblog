package models

import (
	"log"

	"github.com/LeeReindeer/lightblog/util"

	"github.com/astaxie/beego/orm"
)

func GetAllTags() (tags []Tag) {
	return nil
}

// a blog only has a tag
func GetTagById(tagId int64) *Tag {
	o := orm.NewOrm()
	tag := Tag{TagId: tagId}
	o.Read(&tag)
	return &tag
}

func GetBlogsWithTag(tagId int64) (blogs []LightBlog) {
	if tagId == 0 {
		return nil
	}

	o := orm.NewOrm()
	_, err := o.Raw("SELECT * FROM blog WHERE blog_tag_id=? ORDER BY blog_time DESC", tagId).QueryRows(&blogs)
	if err != nil {
		return nil
	}

	userMap := make(map[int64]*User)
	for i, blog := range blogs {
		if userMap[blog.BlogUid] == nil {
			user := User{UserId: blog.BlogUid}
			// only read a user once
			err = o.Read(&user)
			util.CheckDBErr(err)
			userMap[blog.BlogUid] = &user
		}

		blogs[i].Tag = *GetTagById(tagId)
		blogs[i].BlogPreview = getBlogPreview(blog.BlogContent)
		blogs[i].BlogTimeString = blog.BlogTime.Format("2006-01-02 15:04:05")
		blogs[i].BlogUsername = userMap[blog.BlogUid].UserName
		blogs[i].BlogUserAvatar = userMap[blog.BlogUid].UserAvatar
	}
	return
}

// Insert a new tag, if no err, return TagId, true
// if tag exits return (TagId, false)
func SaveTag(tag *Tag) (int64, bool) {
	tagId, isExist := IsTagExist(tag.TagName)
	if !isExist {
		// insert new tag
		o := orm.NewOrm()
		tagId, err := o.Insert(tag)
		if err != nil {
			log.Println(err.Error())
			return 0, false
		}
		return tagId, true
	}
	return tagId, !isExist
}

func IsTagExist(tagName string) (int64, bool) {
	o := orm.NewOrm()
	tag := Tag{TagName: tagName}
	err := o.Read(&tag, "TagName")
	if err != nil {
		log.Println(err.Error())
		return 0, false
	}
	return tag.TagId, true
}
