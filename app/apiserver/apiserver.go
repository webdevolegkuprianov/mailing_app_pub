package apiserver

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	logger "github.com/webdevolegkuprianov/client_email_app/app/apiserver/logger"
	"github.com/webdevolegkuprianov/client_email_app/app/apiserver/store/httpstore"
	"github.com/webdevolegkuprianov/client_email_app/model"
)

func Start(conf *model.Config) error {

	//server binding
	engine := gin.New()

	store := httpstore.NewConf(engine, conf)

	server := newServer(store, engine, conf)

	var errChan = make(chan error, 1)

	go func() {

		errChan <- server.engine.Run(conf.Spec.Ports.Addr)

	}()

	var signalChan = make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-signalChan:
		logger.InfoLogger.Println("interrupt, exit")
	case err := <-errChan:
		if err != nil {
			logger.ErrorLogger.Println("error api, exit")
		}
	}

	return nil

}
