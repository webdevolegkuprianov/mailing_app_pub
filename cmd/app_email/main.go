package main

import (
	"github.com/webdevolegkuprianov/client_email_app/app/apiserver"
	"github.com/webdevolegkuprianov/client_email_app/model"

	logger "github.com/webdevolegkuprianov/client_email_app/app/apiserver/logger"
)

func main() {

	conf, err := model.NewConfig()
	if err != nil {
		logger.ErrorLogger.Println(err)
	}

	if err := apiserver.Start(conf); err != nil {
		logger.ErrorLogger.Println(err)
	}

}
