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
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jackgris/goscrapy-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "With this you have the option of save multiple wholesaler data",
	Long:  `This will help you having differents setup for different wholesalers, so you will have the option of choose the wholesaler who you want to scrape with the config command.`,
	Run: func(cmd *cobra.Command, args []string) {
		util.ConfigFolderCreate()
		saveWholesalerData()
	},
}

func init() {
	rootCmd.AddCommand(saveCmd)
}

func saveWholesalerData() {

	var conf Config
	// perform the questions
	err := survey.Ask(wholesalerQs, &conf)
	if err != nil {
		log.Printf("Error when get prompt input: %s", err.Error())
		return
	}

	KEY = os.Getenv("GOSCRAPY_KEY")
	if KEY == "" {
		log.Fatalln("You need setup the env variable GOSCRAPY_KEY in your terminal")
	}

	viper.SetConfigType("yaml")
	viper.AddConfigPath("config/")

	pass := conf.Pass
	passHash, err := util.Encrypt(pass, KEY)
	if err != nil {
		log.Printf("Can't hash password: %s", err)
	}

	viper.SetConfigName(conf.Name)
	viper.Set("name", conf.Name)
	viper.Set("login", conf.Login)
	viper.Set("pass", passHash)
	viper.Set("User", conf.User)
	viper.Set("searchpage", conf.SearchPage)
	viper.Set("endphrase", conf.EndPhrase)
	viper.Set("endphrasediv", conf.EndPhraseDiv)
	if err = viper.SafeWriteConfig(); err != nil {
		log.Printf("Error while write config file: %s", err)
	}
}
