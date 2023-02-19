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

	"github.com/jackgris/goscrapy-cli/util"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update data from wholesalers who already have been saved",
	Long: `With this command, you can update data from wholesalers who already have been saved in yaml files
in the config folder. You can choose one wholesaler and one field and update his value.`,
	Run: func(cmd *cobra.Command, args []string) {
		util.ConfigFolderCreate()
		names := util.GetNameFilesConfig()
		// FIXME we need to use the name to ask the user what file want to update
		log.Println(names)
		askFieldToUpdate("Ahora")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
