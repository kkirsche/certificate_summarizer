// Copyright © 2018 Kevin Kirsche <kev.kirsche[at]gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
  "fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// Default build-time variable.
// These values are overridden via ldflags
var (
	BuildBinary  = "certificate_summarizer"
	BuildVersion = "development-not_available"
	BuildHash    = "development-not_available"
	BuildTime    = "development-not_available"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version and build information about the tool",
	Long: `Display version and build information about the tool`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Binary:\t%s\n", BuildBinary)
		fmt.Printf("Version:\t%s\n", BuildVersion)
		fmt.Printf("Go Version:\t%s\n", runtime.Version())
		fmt.Printf("OS:\t\t%s\n", runtime.GOOS)
		fmt.Printf("Arch:\t\t%s\n", runtime.GOARCH)
		fmt.Printf("Git Hash:\t%s\n", BuildHash)
		fmt.Printf("Build Time:\t%s\n", BuildTime)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
