package util

import (
	"github.com/astaxie/beego/context"
	"log"
	"strconv"
	"strings"
)

func CheckDBErr(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

func CheckErr(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

func CheckAndFatal(err error) {
	if err != nil {
		panic(err)
		log.Fatalln("exit(1)")
	}
}

func ClearCookies(ctx *context.Context) {
	ctx.SetCookie("uid", "")
	ctx.SetCookie("username", "")
	ctx.SetCookie("login", "false")
	ctx.SetCookie("p", "")
}

func StringToInt64(str string) (int64, bool) {
	i, err := strconv.Atoi(str)
	if err != nil {
		log.Println(err.Error())
		return 0, false
	}
	return int64(i), true
}

func GetUserIdFromCookie(ctx *context.Context) (int64, bool) {
	uidStr := ctx.GetCookie("uid")
	if len(uidStr) == 0 {
		return 0, false
	}
	uid, ok := StringToInt64(uidStr)
	if !ok {
		return 0, false
	}
	return uid, true
}

func Redirect302(url string, ctx *context.Context) {
	ctx.Redirect(302, url)
}

func GetContentTag(content string) (string, bool, int) {
	if !IsContentTagged(content) {
		return "", false, -1
	}
	index := strings.Index(content, " ")
	return content[1:index], true, index + 1
}

// tag with prefix '#', and less than 10 chars, no less than 1
func IsContentTagged(content string) bool {
	contentRune := []rune(content)
	//log.Println("len of rune: ", len(contentRune))
	index := -1
	for i, value := range contentRune {
		if string(value) == " " {
			index = i
			break
		}
	}
	//log.Println("index of space: ", index)
	//log.Println("tagged: ", strings.HasPrefix(content, "#"))
	return strings.HasPrefix(content, "#") && index <= 25 && index > 1
}
