// Copyright 2024
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/penny-vault/edgar/financials"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "edgar",
	Short: "Parse 10-Q filings from SEC EDGAR database",
	Long:  `Download and extract fundamental data from XBRL`,
	Run: func(cmd *cobra.Command, args []string) {

		// get a list of filings by CIK
		// https://data.sec.gov/submissions/CIK0000096223.json
		statement, err := financials.ParseXBRL(args[0])
		if err != nil {
			log.Fatal().Err(err).Str("fn", args[0]).Msg("could not parse XBRL from file")
		}

		fmt.Println(statement)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.edgar.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func initLogging() {
	// Logging configuration
	if err := viper.BindEnv("log.level", "EDGAR_LOG_LEVEL"); err != nil {
		log.Panic().Err(err).Msg("could not bind EDGAR_LOG_LEVEL")
	}
	rootCmd.PersistentFlags().String("log-level", "warning", "Logging level")
	if err := viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("log-level")); err != nil {
		log.Panic().Err(err).Msg("could not bind log-level")
	}

	if err := viper.BindEnv("log.report_caller", "EDGAR_LOG_REPORT_CALLER"); err != nil {
		log.Panic().Err(err).Msg("could not bind EDGAR_LOG_REPORT_CALLER")
	}
	rootCmd.PersistentFlags().Bool("log-report-caller", false, "Log function name that called log statement")
	if err := viper.BindPFlag("log.report_caller", rootCmd.PersistentFlags().Lookup("log-report-caller")); err != nil {
		log.Panic().Err(err).Msg("could not bind log-report-caller")
	}

	if err := viper.BindEnv("log.output", "EDGAR_LOG_OUTPUT"); err != nil {
		log.Panic().Err(err).Msg("could not bind EDGAR_LOG_OUTPUT")
	}
	rootCmd.PersistentFlags().String("log-output", "stdout", "Write logs to specified output one of: file path, `stdout`, or `stderr`")
	if err := viper.BindPFlag("log.output", rootCmd.PersistentFlags().Lookup("log-output")); err != nil {
		log.Panic().Err(err).Msg("could not bind log-output")
	}

	if err := viper.BindEnv("log.otlp_url", "OTLP_URL"); err != nil {
		log.Panic().Err(err).Msg("could not bind OTLP_URL")
	}
	rootCmd.PersistentFlags().String("log-otlp-url", "", "OTLP server to send traces to, if blank don't send traces")
	if err := viper.BindPFlag("log.otlp_url", rootCmd.PersistentFlags().Lookup("log-otlp-url")); err != nil {
		log.Panic().Err(err).Msg("could not bind log-otlp-url")
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".edgar" (without extension).
		viper.AddConfigPath("/etc/") // path to look for the config file in
		viper.AddConfigPath(fmt.Sprintf("%s/.config", home))
		viper.AddConfigPath(".")

		viper.SetConfigType("toml")
		viper.SetConfigName("edgar")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info().Str("ConfigFile", viper.ConfigFileUsed()).Msg("Loaded config file")
	} else {
		log.Error().Stack().Err(err).Msg("error reading config file")
		os.Exit(1)
	}
}
