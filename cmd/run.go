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

	"github.com/jackgris/goscrapy-cli/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start scraping the web site you choice with config setup",
	Long: `Once you configure all the necesary data for the wholesaler who you want scrape.
Can run this command for start saving data, the default opcion is a JSON file.`,
	Run: func(cmd *cobra.Command, args []string) {
		var conf Config
		var ok bool

		setupEnv()

		if conf, ok = setup(); !ok {
			log.Fatalln("You need to setup the wholesaler data again")
		}

		log.Println("You will start scraping: ", conf.Name)

		w := internal.Wholesalers{
			Login:        conf.Login,
			User:         conf.User,
			Pass:         conf.Pass,
			Searchpage:   conf.SearchPage,
			Name:         conf.Name,
			EndPhrase:    "No tenemos",
			EndPhraseDiv: "div.text-center",
		}

		file, err := os.OpenFile("products.json", os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalln("Error creating file:", err)
		}
		defer file.Close()

		log := logrus.New()
		log.SetLevel(logrus.InfoLevel)
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:          true,
			TimestampFormat:        "2006-01-02 15:04:05",
			ForceColors:            true,
			DisableLevelTruncation: true,
		})

		if n, err := file.WriteString("["); err != nil {
			log.Fatalln(n, err)
		}

		if err = internal.GetData(file, w, log); err != nil {
			log.Fatalln(err)
		}
		// Remove the last character from the file, so we remove the last , we added before know
		// we append the last element
		if _, err = file.Seek(-2, 2); err != nil { // Set the file pointer to the last character
			log.Fatalln(err)
		}

		if _, err = file.Write([]byte{}); err != nil {
			log.Fatalln(err)
		}

		if n, err := file.WriteString("]"); err != nil {
			log.Fatalln(n, err)
		}

	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
