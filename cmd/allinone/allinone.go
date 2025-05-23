//go:build linux

/*
 * Copyright (C) 2016 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy ofthe License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specificlanguage governing permissions and
 * limitations under the License.
 *
 */

package allinone

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/avast/retry-go"
	"github.com/kardianos/osext"
	"github.com/spf13/cobra"

	"github.com/skydive-project/skydive/analyzer"
	"github.com/skydive-project/skydive/cmd"
	cmdconfig "github.com/skydive-project/skydive/cmd/config"
	"github.com/skydive-project/skydive/config"
	"github.com/skydive-project/skydive/graffiti/http"
	"github.com/skydive-project/skydive/graffiti/logging"
	"github.com/skydive-project/skydive/graffiti/service"
)

var (
	analyzerUID uint32
)

// AllInOneCmd describes the allinone root command
var AllInOneCmd = &cobra.Command{
	Use:          "allinone",
	Short:        "Skydive All-In-One",
	Long:         "Skydive All-In-One (Analyzer + Agent)",
	SilenceUsage: true,
	Run: func(c *cobra.Command, args []string) {
		if uid := os.Geteuid(); uid != 0 {
			fmt.Fprintln(os.Stderr, "All-In-One mode has to be run as root")
			os.Exit(1)
		}

		skydivePath, _ := osext.Executable()
		logFile := config.GetString("logging.file.path")
		extension := filepath.Ext(logFile)
		logFile = strings.TrimSuffix(logFile, extension)
		os.Setenv("SKYDIVE_LOGGING_FILE_PATH", logFile+"-analyzer"+extension)

		analyzerAttr := &os.ProcAttr{
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
			Env:   os.Environ(),
		}

		if analyzerUID != 0 {
			analyzerAttr.Sys = &syscall.SysProcAttr{
				Credential: &syscall.Credential{
					Uid: analyzerUID,
				},
			}
		}

		args = []string{"skydive"}

		for _, cfgFile := range cmd.CfgFiles {
			args = append(args, "-c")
			args = append(args, cfgFile)
		}

		if cmd.CfgBackend != "" {
			args = append(args, "-b")
			args = append(args, cmd.CfgBackend)
		}

		analyzerArgs := make([]string, len(args))
		copy(analyzerArgs, args)

		if len(cmd.CfgFiles) != 0 {
			if err := cmdconfig.LoadConfiguration(cmd.CfgBackend, cmd.CfgFiles); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to initialize config: %s", err.Error())
				os.Exit(1)
			}
		}

		analyzerProcess, err := os.StartProcess(skydivePath, append(analyzerArgs, "analyzer"), analyzerAttr)
		if err != nil {
			logging.GetLogger().Errorf("Can't start Skydive analyzer: %v", err)
			os.Exit(1)
		}

		authOptions := analyzer.ClusterAuthenticationOpts()
		svcAddr, _ := service.AddressFromString(config.GetString("analyzer.listen"))
		tlsConfig, err := config.GetTLSClientConfig(true)
		if err != nil {
			logging.GetLogger().Errorf("Can't start Skydive analyzer: %v", err)
			os.Exit(1)
		}
		addr := svcAddr.Addr
		if addr == "0.0.0.0" {
			addr = "127.0.0.1"
		}

		url, err := config.GetURL("http", addr, svcAddr.Port, "")
		if err != nil {
			logging.GetLogger().Errorf("Invalid service address: %s", err)
			os.Exit(1)
		}
		restClient := http.NewRestClient(url, authOptions, tlsConfig)

		err = retry.Do(func() error {
			_, err := restClient.Request("GET", "/api", nil, nil)
			return err
		})

		if err != nil {
			logging.GetLogger().Errorf("Failed to start Skydive analyzer: %v", err)
			os.Exit(1)
		}

		os.Setenv("SKYDIVE_ANALYZERS", fmt.Sprintf("%s:%d", addr, svcAddr.Port))
		os.Setenv("SKYDIVE_LOGGING_FILE_PATH", logFile+"-agent"+extension)

		agentAttr := &os.ProcAttr{
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
			Env:   os.Environ(),
		}

		agentArgs := make([]string, len(args))
		copy(agentArgs, args)

		agentProcess, err := os.StartProcess(skydivePath, append(agentArgs, "agent"), agentAttr)
		if err != nil {
			logging.GetLogger().Errorf("Can't start Skydive agent: %s. Please check you are root", err.Error())
			os.Exit(1)
		}

		logging.GetLogger().Notice("Skydive All-in-One starting !")
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch

		analyzerProcess.Signal(os.Interrupt)
		agentProcess.Signal(os.Interrupt)

		analyzerProcess.Wait()
		agentProcess.Wait()

		logging.GetLogger().Notice("Skydive All-in-One stopped.")
	},
}

func init() {
	var uid uint32

	skydivePath, _ := osext.Executable()

	fi, err := os.Stat(skydivePath)
	if err == nil {
		uid = fi.Sys().(*syscall.Stat_t).Uid
	}

	AllInOneCmd.Flags().Uint32VarP(&analyzerUID, "analyzer-uid", "", uid, "uid used for analyzer process")
}
