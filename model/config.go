package model

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

//config yaml struct
type Config struct {
	APIVersion string `yaml:"apiVersion"`
	Spec       struct {
		Ports struct {
			Name string `yaml:"name"`
			Addr string `yaml:"bind_addr"`
		} `yaml:"ports"`
		SendpulseConfig struct {
			ClientId           string  `yaml:"client_id"`
			ClientSecret       string  `yaml:"client_secret"`
			TimetoRefreshToken int     `yaml:"time_to_refresh_token"`
			TokenTermLife      float64 `yaml:"token_term_life"`
		} `yaml:"sendpulse_config"`
		DadataConfig struct {
			DadataKey    string `yaml:"dadata_key"`
			DadataSecret string `yaml:"dadata_secret"`
			DadataToken  string `yaml:"dadata_token"`
		} `yaml:"dadata_config"`
		Client struct {
			Timeout int `yaml:"timeout"`
		} `yaml:"client"`
		App struct {
			UrlOwnerTemplatesList string `yaml:"url_owner_templates_list"`
			UrlOwnerIp            string `yaml:"url_owner_ip"`
			UrlTemplateInfo       string `yaml:"url_template_info"`
			UrlOauth              string `yaml:"url_oauth"`
			UrlSmtp               string `yaml:"url_smtp"`
			UrlDadataGeo          string `yaml:"url_dadata_geo"`
			TimeFormat            string `yaml:"time_format"`
		} `yaml:"app"`
	} `yaml:"spec"`
}

//New config
func NewConfig() (*Config, error) {

	var service *Config

	f, err := filepath.Abs("/root/config/app_email.yaml")
	if err != nil {
		return nil, err
	}

	y, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(y, &service); err != nil {
		return nil, err
	}

	return service, nil

}
