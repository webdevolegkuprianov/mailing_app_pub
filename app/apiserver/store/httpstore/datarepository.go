package httpstore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	"github.com/webdevolegkuprianov/client_email_app/model"

	logger "github.com/webdevolegkuprianov/client_email_app/app/apiserver/logger"
)

//Data repository
type DataRepository struct {
	store *Store
}

//smtp request
func (r *DataRepository) SmtpBooking(body model.DataBooking, conf *model.Config, token string) (*http.Response, error) {

	newSmtpParams := model.NewSmtpParams()
	var urlStaticMap string

	//handling parametres
	if body.Var7Tuning == "no_data" {
		body.Var7Tuning = " забронирован"
	}

	if body.Var2Name == "no_data" {
		body.Var1Surname = body.Var1_1Surname
		body.Var2Name = body.Var2_1Name
	}

	//add static yandex map
	//get geo
	resg, err := r.store.Data().RequestDadataGeo(conf, body.DeliveryFiasCode)
	if err != nil {
		urlStaticMap = ""
	} else {
		urlStaticMap = fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s",
			newSmtpParams.UrlStaticMaps,
			"ll=", resg.DgGeo,
			",", resg.ShGeo,
			"&z=15&size=450,450&l=map&pt=",
			resg.DgGeo,
			",",
			resg.ShGeo,
			",pm2lbl")
	}

	a := &model.Email1{
		Email: model.Email{
			Subject: newSmtpParams.Subject,
			Template: model.Template{
				IdTemplate: newSmtpParams.TemplateIdBooking,
				TemplateVariables: model.TemplateVariables{
					Var1Surname:          body.Var1Surname,
					Var2Name:             body.Var2Name,
					Var3Semeistvo:        body.Var3Semeistvo,
					Var4TipKuzov:         body.Var4TipKuzov,
					Var5Engine:           body.Var5Engine,
					Var6Base:             body.Var6Base,
					Var7Tuning:           body.Var7Tuning,
					Var8Modification:     body.Var8Modification,
					Var9Price:            body.Var9Price,
					Var10Vin:             body.Var10Vin,
					Var11Url:             body.Var11Url,
					Var12DeliveryAddress: body.Var12DeliveryAddress,
					Var13UrlStaticMap:    urlStaticMap,
				},
			},
			From: model.From{
				Name:  newSmtpParams.NameFrom,
				Email: newSmtpParams.EmailFrom,
			},
			To: []struct {
				Email string `json:"email"`
				Name  string `json:"name"`
			}{
				{Email: body.ToEmail, Name: body.Var2Name},
			},
			Attach: model.Attach{
				File: body.File,
			},
		},
	}

	//d_spaces, err := json.MarshalIndent(a, "", "    ")
	//if err != nil {
	//return nil, err
	//}

	d, err := json.Marshal(a)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	client := &http.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, conf.Spec.App.UrlSmtp, bytes.NewBuffer(d))
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	t := fmt.Sprintf("%s%s", "Bearer ", token)

	req.Header.Add("Authorization", t)

	resp, err := client.Do(req)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	return resp, nil
}

//dadata request
func (r *DataRepository) RequestDadataGeo(conf *model.Config, codefias string) (*model.ResultGeo, error) {

	res := model.Results{}

	var (
		dadata model.DadataRequestGeo
	)

	dadata.CodeFias = codefias

	client := http.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bodyBytes, err := json.Marshal(dadata)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, conf.Spec.App.UrlDadataGeo, bytes.NewBuffer(bodyBytes))
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	token := fmt.Sprintf("%s%s", "Token ", conf.Spec.DadataConfig.DadataToken)

	req.Header.Add("Authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	defer resp.Body.Close()

	resp_bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	defer resp.Body.Close()

	if err := json.Unmarshal(resp_bodyBytes, &res); err != nil {
		return nil, err
	}

	valSlice := reflect.ValueOf(res).FieldByName("Suggestions").Interface().([]model.Suggestions)

	if len(valSlice) == 0 {
		return &model.ResultGeo{
			ShGeo: "",
			DgGeo: "",
		}, nil
	}

	shGeo := valSlice[0].Data.Geo_lat
	dgGeo := valSlice[0].Data.Geo_lon

	resultGeo := &model.ResultGeo{
		ShGeo: shGeo,
		DgGeo: dgGeo,
	}

	return resultGeo, nil

}
