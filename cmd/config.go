/*
Copyright © 2023 Gabriel Pozo <jackgris2@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure necessary data for login web page and collect prices",
	Long: `With this command you can load the data needed for login in the web page where
we want to collect product prices.`,
	Run: func(cmd *cobra.Command, args []string) {
		var conf Config
		var ok bool

		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")

		if conf, ok = setup(); ok {

			if ok := checkData(conf); !ok {
				asking()
			} else {
				if currentData(conf) {
					log.Println("We need start scraping, conf data: ", conf)
				} else {
					asking()
				}
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

// the questions to ask
var qs = []*survey.Question{
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

// asking collect data from user
func asking() {
	var conf Config
	// perform the questions
	err := survey.Ask(qs, &conf)
	if err != nil {
		log.Printf("Error when get prompt input: %s", err.Error())
		return
	}

	viper.Set("name", conf.Name)
	viper.Set("login", conf.Login)
	viper.Set("pass", conf.Pass)
	viper.Set("User", conf.User)
	viper.Set("searchpage", conf.SearchPage)
	if err = viper.WriteConfig(); err != nil {
		log.Printf("Error while write config file: %s", err)
	}
}

type Config struct {
	Name       string
	Login      string
	User       string
	Pass       string
	SearchPage string
}

// setup will check if you have the necessary data
func setup() (Config, bool) {
	var conf Config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Printf("File not found: %s", err)
			err = viper.SafeWriteConfig()
			if err != nil {
				log.Printf("Can't create config file: %s", err)
			}
			return conf, false
		} else {
			// Config file was found but another error was produced
			log.Printf("Occur an error when trying read config file: %s", err)
			return conf, false
		}
	}

	conf.Name = viper.GetString("name")
	conf.Login = viper.GetString("login")
	conf.User = viper.GetString("user")
	conf.Pass = viper.GetString("pass")
	conf.SearchPage = viper.GetString("searchpage")

	return conf, true
}

// checkData verify is any necessary data is empty
func checkData(conf Config) bool {
	if conf.Login == "" || conf.User == "" || conf.Pass == "" || conf.SearchPage == "" || conf.Name == "" {
		return false
	}
	return true
}

var currentAsk = &survey.Confirm{
	Message: "Do you want use the current data:",
}

// currentData ask to the user if want use the data saved on the config file
func currentData(conf Config) bool {

	var actual bool
	// perform the questions
	err := survey.AskOne(currentAsk, &actual)
	if err != nil {
		log.Printf("Error when get prompt input: %s", err.Error())
		return false
	}

	return actual
}
