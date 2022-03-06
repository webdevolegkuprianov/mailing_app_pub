package httpstore

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/webdevolegkuprianov/client_email_app/model"

	logger "github.com/webdevolegkuprianov/client_email_app/app/apiserver/logger"
)

//Data repository
type ServRepository struct {
	store *Store
}

//Oauth body
type Data_oauth struct {
	Grant_type    string `json:"grant_type"`
	Client_id     string `json:"client_id"`
	Client_secret string `json:"client_secret"`
}

//Token body response
type Data_token struct {
	Token string `json:"access_token"`
}

//verify token
func (r *ServRepository) VerifyTokenDuration(conf *model.Config) error {

	currentTime := time.Now().Format(conf.Spec.App.TimeFormat)

	duration, err := r.store.Serv().DurationToken(conf, currentTime)
	if err != nil {
		return err
	}
	delta := conf.Spec.SendpulseConfig.TokenTermLife - duration

	if delta < float64(conf.Spec.SendpulseConfig.TimetoRefreshToken) {

		requestBody, err := json.Marshal(
			Data_oauth{
				Client_id:     conf.Spec.SendpulseConfig.ClientId,
				Client_secret: conf.Spec.SendpulseConfig.ClientSecret,
				Grant_type:    "client_credentials"})

		if err != nil {
			return err
		}

		//request refresh token oauth
		resp, err := http.Post(conf.Spec.App.UrlOauth, "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			logger.ErrorLogger.Println(err)
			return err
		}

		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.ErrorLogger.Println(err)
			return err
		}

		var data_token Data_token

		err = json.Unmarshal(bodyBytes, &data_token)
		if err != nil {
			logger.ErrorLogger.Println(err)
			return err
		}

		//write new time create and token
		file_token, err := os.OpenFile("/root/config/token_email/tokenBearer.txt", os.O_TRUNC|os.O_RDWR, 0600)
		if err != nil {
			logger.ErrorLogger.Println(err)
			return err
		}

		defer file_token.Close()

		if err := file_token.Truncate(0); err != nil {
			logger.ErrorLogger.Println(err)
			return err
		}

		_, err = file_token.WriteString(data_token.Token)
		if err != nil {
			logger.ErrorLogger.Println(err)
			return err
		}

		file_exp, err := os.OpenFile("/root/config/token_email/tokenExp.txt", os.O_TRUNC|os.O_RDWR, 0600)
		if err != nil {
			logger.ErrorLogger.Println(err)
			return err
		}

		defer file_exp.Close()

		if err := file_exp.Truncate(0); err != nil {
			logger.ErrorLogger.Println(err)
			return err
		}

		_, err = file_exp.WriteString(currentTime)
		if err != nil {
			logger.ErrorLogger.Println(err)
			return err
		}

		return nil

	} else {

		return nil

	}

}

//token duration
func (r *ServRepository) DurationToken(conf *model.Config, currentTime string) (float64, error) {

	file, err := os.OpenFile("/root/config/token_email/tokenExp.txt", os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return 0, err
	}

	defer file.Close()

	body, err := ioutil.ReadAll(file)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return 0, err
	}
	bodyTimeString := string(body)

	if len(bodyTimeString) == 0 {

		_, err := file.WriteString(currentTime)
		if err != nil {
			logger.ErrorLogger.Println(err)
			return 0, err
		}

		return conf.Spec.SendpulseConfig.TokenTermLife, nil

	} else {

		a, err := time.Parse(conf.Spec.App.TimeFormat, bodyTimeString)
		if err != nil {
			logger.ErrorLogger.Println(err)
			return 0, err
		}
		duration := time.Since(a).Minutes()

		return duration, nil

	}

}

//get token
func (r *ServRepository) GetToken() (string, error) {

	file, err := os.OpenFile("/root/config/token_email/tokenBearer.txt", os.O_RDONLY, 0600)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return "", err
	}

	defer file.Close()

	bodyBytes, err := ioutil.ReadAll(file)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return "", err
	}
	bodyString := string(bodyBytes)

	return bodyString, nil

}
