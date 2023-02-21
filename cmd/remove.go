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

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "You can remove the wholesaler data you want by choosing his name",
	Long:  `Anytime you want to remove the wholesaler data that you saved before. Using this command will show you the entire list of names so you can choose anyone and delete his data in an easy way.`,
	Run: func(cmd *cobra.Command, args []string) {
		util.ConfigFolderCreate()
		names := util.GetNameFilesConfig()
		name := askForWholesale(names)
		if name == "" {
			log.Println("You need choose one wholesaler or save one before try to remove one.")
			return
		}

		if err := askYouSureRemove(name); err != nil {
			log.Fatalf("Can't remove wholesaler data: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
