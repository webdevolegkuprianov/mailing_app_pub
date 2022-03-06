package httpstore

import (
	"github.com/gin-gonic/gin"
	"github.com/webdevolegkuprianov/client_email_app/app/apiserver/store"
	"github.com/webdevolegkuprianov/client_email_app/model"
)

type Store struct {
	engine         *gin.Engine
	conf           *model.Config
	servRepository *ServRepository
	dataRepository *DataRepository
}

func NewConf(engine *gin.Engine, conf *model.Config) *Store {

	return &Store{
		engine: engine,
		conf:   conf,
	}

}

//Serv
func (s *Store) Serv() store.ServRepository {
	if s.servRepository != nil {
		return s.servRepository
	}

	s.servRepository = &ServRepository{
		store: s,
	}

	return s.servRepository
}

//Data
func (s *Store) Data() store.DataRepository {
	if s.dataRepository != nil {
		return s.dataRepository
	}

	s.dataRepository = &DataRepository{
		store: s,
	}

	return s.dataRepository
}
