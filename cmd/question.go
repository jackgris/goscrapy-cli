/*
	Copyright © 2023 Gabriel Pozo

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
		Name:     "login",
		Prompt:   &survey.Input{Message: "Enter login URL:"},
		Validate: urlValidator(),
	},
	{
		Name:     "user",
		Prompt:   &survey.Input{Message: "Enter user name or email:"},
		Validate: emailValidator(),
	},
	{
		Name:     "pass",
		Prompt:   &survey.Password{Message: "Enter the password:"},
		Validate: survey.Required,
	},
	{
		Name:     "searchpage",
		Prompt:   &survey.Input{Message: "Enter URL where are the products:"},
		Validate: urlValidator(),
	},
	{
		Name:     "endphrase",
		Prompt:   &survey.Input{Message: "Enter the phrase that you see when don't see any product:"},
		Validate: survey.Required,
	},
	{
		Name:     "endphrasediv",
		Prompt:   &survey.Input{Message: "Enter the specific div where the end phrase appear:"},
		Validate: survey.Required,
	},
}

// question for update fields of wholesalers
var updateDataQs = []*survey.Question{
	{
		Name: "field",
		Prompt: &survey.Select{
			Message: "Choose the field to update:",
			Options: []string{"login", "user", "pass", "searhpage", "endphrase", "endphrasediv"},
			Default: "login",
		},
	},
	{
		Name:      "value",
		Prompt:    &survey.Input{Message: "What is the value?"},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
}

// askingWholesalerData collect data from user about wholesaler
func askingWholesalerData(name string) error {
	var conf Config
	var err error
	viper.AddConfigPath("config/")
	viper.SetConfigName(name)
	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("Wholesaler file not found: %s", err)
		}
		// Config file was found but another error was produced
		log.Fatalf("Problems reading wholesaler data: %s", err)
	}

	conf.Name = viper.GetString("name")
	conf.Login = viper.GetString("login")
	conf.User = viper.GetString("user")
	conf.Pass = viper.GetString("pass")
	conf.SearchPage = viper.GetString("searchpage")
	conf.EndPhrase = viper.GetString("endphrase")
	conf.EndPhraseDiv = viper.GetString("endphrasediv")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("Wholesaler file not found: %s", err)
		}
		// Config file was found but another error was produced
		log.Fatalf("Problems reading wholesaler data: %s", err)
	}

	viper.Set("name", conf.Name)
	viper.Set("login", conf.Login)
	viper.Set("pass", conf.Pass)
	viper.Set("User", conf.User)
	viper.Set("searchpage", conf.SearchPage)
	viper.Set("endphrase", conf.EndPhrase)
	viper.Set("endphrasediv", conf.EndPhraseDiv)

	if err = viper.WriteConfig(); err != nil {
		log.Printf("Error while write config file: %s", err)
		return err
	}

	return err
}

type Update struct {
	Field string
	Value string
}

func askYouSureRemove(file string) error {
	name := util.FileNameWithoutExtSliceNotation(file)

	var askSure = &survey.Confirm{
		Message: "Are you sure you want to remove " + name + " data?",
	}
	var sure bool
	// perform the questions
	err := survey.AskOne(askSure, &sure)
	if err != nil {
		return err
	}

	if sure {
		err := util.RemoveFile("config/" + file)
		return err
	}

	return nil
}

func askForWholesale(files []string) string {
	var file string

	// question for update fields of wholesalers
	var chooseWSalerQs = []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Select{
				Message: "Choose the wholesaler:",
				Options: files,
				Default: files[0],
			},
		},
	}

	// perform the questions
	if err := survey.Ask(chooseWSalerQs, &file); err != nil {
		log.Printf("Error when get prompt input: %s", err.Error())
		return ""
	}

	return file
}

func askForWholesaleUpdate(files []string) string {
	var name string
	var names []string

	for _, file := range files {
		fileName := util.FileNameWithoutExtSliceNotation(file)
		names = append(names, fileName)
	}
	// question for update fields of wholesalers
	var chooseWSalerQs = []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Select{
				Message: "Choose the name of the wholesaler:",
				Options: names,
				Default: names[0],
			},
		},
	}

	// perform the questions
	err := survey.Ask(chooseWSalerQs, &name)
	if err != nil {
		log.Printf("Error when get prompt input: %s", err.Error())
		return ""
	}
	return name
}

func askFieldToUpdate(name string) {

	KEY = os.Getenv("GOSCRAPY_KEY")
	if KEY == "" {
		log.Fatalln("You need setup the env variable GOSCRAPY_KEY in your terminal")
	}

	viper.SetConfigType("yaml")
	viper.AddConfigPath("config/")
	viper.SetConfigName(name)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Can't read the config file %s", err)
	}

	var update Update
	// perform the questions
	err := survey.Ask(updateDataQs, &update)
	if err != nil {
		log.Printf("Error when get prompt input: %s", err.Error())
		return
	}

	if update.Field == "pass" {
		passHash, err := util.Encrypt(update.Value, KEY)
		if err != nil {
			log.Printf("Can't hash password: %s", err)
		}
		update.Value = passHash
	}

	viper.Set(update.Field, update.Value)
	if err = viper.WriteConfig(); err != nil {
		log.Printf("Error while write config file: %s", err)
	}
}
