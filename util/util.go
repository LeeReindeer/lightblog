package util

import (
	"github.com/astaxie/beego/context"
	"log"
	"strconv"
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
