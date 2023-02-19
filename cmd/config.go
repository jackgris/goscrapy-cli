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

		setupEnv()

		if conf, ok = setup(); ok {

			if ok := checkData(conf); !ok {
				askingWholesalerData()
			} else {
				if currentData(conf) {
					log.Println("All data is right, you can start web scraping")
				} else {
					askingWholesalerData()
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
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
		}
		// Config file was found but another error was produced
		log.Printf("Occur an error when trying read config file: %s", err)
		return conf, false
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
