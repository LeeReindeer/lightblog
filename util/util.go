package util

import "log"

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
