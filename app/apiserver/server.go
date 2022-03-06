package apiserver

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	logger "github.com/webdevolegkuprianov/client_email_app/app/apiserver/logger"
	"github.com/webdevolegkuprianov/client_email_app/app/apiserver/store"
	"github.com/webdevolegkuprianov/client_email_app/model"
)

//server configure
type server struct {
	store  store.Store
	engine *gin.Engine
	config *model.Config
}

func newServer(store store.Store, engine *gin.Engine, config *model.Config) *server {
	s := &server{
		store:  store,
		engine: engine,
		config: config,
	}
	s.configureRouter()
	return s
}

//handler set
func (s *server) configureRouter() {

	s.engine.POST("/emailbooking", s.handleEmailBooking)

}

//booking
func (s *server) handleEmailBooking(context *gin.Context) {

	//verify token duration
	err := s.store.Serv().VerifyTokenDuration(s.config)
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		logger.ErrorLogger.Println(err)
		return
	}

	//get token
	token, err := s.store.Serv().GetToken()
	if err != nil {
		logger.ErrorLogger.Println(err)
		return
	}

	req := model.DataBooking{}

	if err := context.BindJSON(&req); err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		logger.ErrorLogger.Println(err)
		return
	}

	//request smtp booking
	resp, err := s.store.Data().SmtpBooking(req, s.config, token)
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		logger.ErrorLogger.Println(err)
	}
	respbodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorLogger.Println(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		logger.ErrorLogger.Println(string(respbodyBytes))
		context.AbortWithStatusJSON(resp.StatusCode, string(respbodyBytes))
	} else {
		logger.InfoLogger.Println("mail ok")
		context.AbortWithStatusJSON(resp.StatusCode, string(respbodyBytes))
	}

}
