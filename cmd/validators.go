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
