package model

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
)

//smtp emailing struct sendpulse
type Email1 struct {
	Email Email `json:"email"`
}

type Email struct {
	Subject  string `json:"subject"`
	Template `json:"template"`
	From     `json:"from"`
	To       `json:"to"`
	Attach   `json:"attachments_binary"`
}

type Template struct {
	IdTemplate        int `json:"id"`
	TemplateVariables `json:"variables"`
}

type TemplateVariables struct {
	Var1Surname          string `json:"surname"`   //фамилия
	Var2Name             string `json:"name"`      //имя
	Var3Semeistvo        string `json:"semeistvo"` //семейство
	Var4TipKuzov         string `json:"tipkuzov"`  //тип кузова
	Var5Engine           string `json:"engine"`    //двигатель
	Var6Base             string `json:"base"`      //база
	Var7Tuning           string `json:"tuning"`    //надстройка
	Var8Modification     string `json:"modification"`
	Var9Price            int    `json:"price"`
	Var10Vin             string `json:"vin"`
	Var11Url             string `json:"url"`
	Var12DeliveryAddress string `json:"deliveryaddress"`
	Var13UrlStaticMap    string `json:"staticmapurl"`
}

type From struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type To []struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Attach struct {
	File string `json:"file.pdf"`
}

type RespSmtpBooking struct {
	Result bool   `json:"result"`
	Id     string `json:"id"`
}

//geo
type Results struct {
	Suggestions []Suggestions
}

type Suggestions struct {
	Value              string `json:"value"`
	Unrestricted_value string `json:"unrestricted_value"`
	Data               struct {
		Geo_lat string `json:"geo_lat"`
		Geo_lon string `json:"geo_lon"`
	}
}

type ResultGeo struct {
	ShGeo string
	DgGeo string
}

type DataUserTemplate struct {
	TemplateId   string `json:"id"`
	TemplateName string `json:"name"`
}

type SmtpParams struct {
	Subject           string //тема письма
	TemplateIdBooking int
	NameFrom          string
	EmailFrom         string
	UrlStaticMaps     string
}

//New config
func NewSmtpParams() *SmtpParams {
	return &SmtpParams{
		Subject:           "Бронирование автомобиля",
		TemplateIdBooking: 26965,
		NameFrom:          "АО Современные транспортные технологии",
		EmailFrom:         "kupriyanovoe@st.tech",
		UrlStaticMaps:     "https://static-maps.yandex.ru/1.x/?",
	}
}

//get base64 from pdf
func PdfToBase64() (string, error) {

	//encode to base64
	PdfBytes, err := ioutil.ReadFile("счет_на_предоплату.pdf")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	str := base64.StdEncoding.EncodeToString(PdfBytes)

	return str, nil

}

//json tag to serialize json body
type DataBooking struct {
	Var1Surname          string `json:"surname"`       //фамилия
	Var2Name             string `json:"client_name"`   //имя
	Var3Semeistvo        string `json:"mod_family"`    //семейство
	Var4TipKuzov         string `json:"mod_body_type"` //тип кузова
	Var5Engine           string `json:"mod_engine"`    //двигатель
	Var6Base             string `json:"mod_base"`      //база
	Var7Tuning           string `json:"mod_tuning"`    //надстройка template changes если == no_data
	Var8Modification     string `json:"modification"`
	Var9Price            int    `json:"price"`
	Var10Vin             string `json:"vin"` //need mask
	Var11Url             string `json:"url_mod"`
	Var12DeliveryAddress string `json:"delivery_address"`
	Var1_1Surname        string `json:"representative_surname"`
	Var2_1Name           string `json:"representative_name"`
	ToEmail              string `json:"client_email"`
	File                 string `json:"file"`
	DeliveryFiasCode     string `json:"delivery_address_code"`
}

//dadata request geo struct
type DadataRequestGeo struct {
	CodeFias string `json:"query"`
}
