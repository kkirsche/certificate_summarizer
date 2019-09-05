// Copyright © 2019 Kevin Kirsche <kevin.kirsche@verizon.com>
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
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/grantae/certinfo"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "certificate_summarizer [file(s) with URLs of websites]",
	Short: "A tool which checks a list of websites and summarizes the certificate information",
	Long:  `A tool which checks a list of websites and summarizes the certificate information`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, urlFiles []string) {
		timeoutDur, _ := time.ParseDuration("5s")
		dialer := &net.Dialer{
			Timeout: timeoutDur,
		}

		summary := map[string]int{}
		failures := map[string]string{}
		for _, urlFile := range urlFiles {
			urls, err := getURLsFromFile(urlFile)
			if err != nil {
				fmt.Println(err)
				continue
			}

			for _, rawURL := range urls {
				parsedURL, err := url.Parse(rawURL)
				if err != nil {
					log.Println("Failed to parse URL: " + err.Error())
					failures[rawURL] = err.Error()
					summary["failed_url_parse"]++
					continue
				}

				host := parsedURL.Host

				if !strings.Contains(host, ":") || strings.HasSuffix(host, "]") {
					host = host + ":443"
				}

				cfg := tls.Config{
					InsecureSkipVerify: true,
				}

				conn, err := tls.DialWithDialer(dialer, "tcp", host, &cfg)
				if err != nil {
					log.Println("TLS connection failed: " + err.Error())
					failures[host] = err.Error()
					summary["failed_tls_conn"]++
					continue
				}

				// Grab the last certificate in the chain
				certChain := conn.ConnectionState().PeerCertificates
				cert := certChain[len(certChain)-1]

				k := "O: " + strings.Join(cert.Subject.Organization, " ") + " OU: " + strings.Join(cert.Subject.OrganizationalUnit, " ")
				summary[k]++

				// Print the certificate
				result, err := certinfo.CertificateText(cert)
				if err != nil {
					log.Fatal(err)
				}
				err = writeCertificateToFile(host, result)
				if err != nil {
					log.Println(err.Error())
					summary["write_error"]++
					failures[host] = err.Error()
				}
			}
		}

		fmt.Println("Successful Connections\nO == Organization\nOU == Organizational Unit")
		for k, v := range summary {
			fmt.Printf("%s: %d\n", k, v)
		}

		fmt.Println("\nFailures")
		for k, v := range failures {
			fmt.Printf("%s: %s\n", k, v)
		}
	},
}

func getURLsFromFile(filename string) ([]string, error) {
	urls := []string{}
	f, err := os.Open(filename)
	if err != nil {
		return urls, fmt.Errorf("Failed to open URL file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	return urls, scanner.Err()
}

func writeCertificateToFile(host, certresult string) error {
	err := os.MkdirAll("results", 0755)
	if err != nil {
		return fmt.Errorf("Failed to create results directory: %w", err)
	}

	f, err := os.Create(fmt.Sprintf("results/%s.txt", strings.Split(host, ":")[0]))
	if err != nil {
		return fmt.Errorf("Failed to create host file for certificate: %w", err)
	}
	defer f.Close()

	_, err = f.WriteString(certresult)
	if err != nil {
		return fmt.Errorf("Failed to write certificate string to file: %w", err)
	}

	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
