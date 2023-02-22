package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jackgris/goscrapy-cli/util"
)

func emailValidator() survey.Validator {
	return func(val interface{}) error {
		email := val.(string)
		if len(email) == 0 {
			return fmt.Errorf("Email is required")
		}
		if !util.IsValidEmail(email) {
			return fmt.Errorf("Invalid email address")
		}
		return nil
	}
}

func urlValidator() survey.Validator {
	return func(val interface{}) error {
		urlStr := val.(string)
		if len(urlStr) == 0 {
			return fmt.Errorf("URL is required")
		}
		if !util.IsValidURL(urlStr) {
			return fmt.Errorf("Invalid URL")
		}
		return nil
	}
}
