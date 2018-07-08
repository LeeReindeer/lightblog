package models

import (
	"github.com/astaxie/beego/orm"
	"log"
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
