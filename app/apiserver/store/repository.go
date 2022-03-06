package store

import (
	"net/http"

	"github.com/webdevolegkuprianov/client_email_app/model"
)

//serv repository
type ServRepository interface {
	VerifyTokenDuration(*model.Config) error
	DurationToken(*model.Config, string) (float64, error)
	GetToken() (string, error)
}

//data repository
type DataRepository interface {
	SmtpBooking(model.DataBooking, *model.Config, string) (*http.Response, error)
	RequestDadataGeo(*model.Config, string) (*model.ResultGeo, error)
}
