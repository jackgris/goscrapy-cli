package cmd

import (
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jackgris/goscrapy-cli/util"
	"github.com/spf13/viper"
)

var KEY string

func setupEnv() {

	KEY = os.Getenv("GOSCRAPY_KEY")
	if KEY == "" {
		log.Fatalln("You need setup the env variable GOSCRAPY_KEY in your terminal")
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
}

// the questions to ask
var wholesalerQs = []*survey.Question{
	{
		Name:      "name",
		Prompt:    &survey.Input{Message: "Enter wholesaler name:"},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
	{
		Name:      "login",
		Prompt:    &survey.Input{Message: "Enter login URL:"},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
	{
		Name:     "user",
		Prompt:   &survey.Input{Message: "Enter user name or email:"},
		Validate: survey.Required,
	},
	{
		Name:     "pass",
		Prompt:   &survey.Password{Message: "Enter the password:"},
		Validate: survey.Required,
	},
	{
		Name:     "searchpage",
		Prompt:   &survey.Input{Message: "Enter URL where are the products:"},
		Validate: survey.Required,
	},
}

// askingWholesalerData collect data from user about wholesaler
func askingWholesalerData() {
	var conf Config
	// perform the questions
	err := survey.Ask(wholesalerQs, &conf)
	if err != nil {
		log.Printf("Error when get prompt input: %s", err.Error())
		return
	}

	pass := conf.Pass
	passHash, err := util.Encrypt(pass, KEY)
	if err != nil {
		log.Printf("Can't hash password: %s", err)
	}

	viper.Set("name", conf.Name)
	viper.Set("login", conf.Login)
	viper.Set("pass", passHash)
	viper.Set("User", conf.User)
	viper.Set("searchpage", conf.SearchPage)
	if err = viper.WriteConfig(); err != nil {
		log.Printf("Error while write config file: %s", err)
	}
}
