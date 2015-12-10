// Copyright 2015 Richard Hawkins
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
package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/hurricanerix/FlappyDisk/app"
	"github.com/hurricanerix/FlappyDisk/gen"
	"gopkg.in/gcfg.v1"
)

//go:generate ./gen_build_info.sh

var GitURL = "https://github.com/hurricanerix/FlappyDisk"
var BuildURL = fmt.Sprintf("%s/commit/%s", GitURL, gen.BuildHash)
var BuildDate = gen.BuildDate

var (
	resetConf bool
	version   bool
)

func init() {
	flag.BoolVar(&resetConf, "reset-conf", false, "reset config to default.")
	flag.BoolVar(&version, "version", false, "print version and build info.")
}

func main() {
	flag.Parse()

	if version {
		fmt.Printf("FlappyDisk Copyright 2015 Richard Hawkins\n")
		fmt.Printf("Licensed under the Apache License, Version 2.0\n")
		fmt.Printf("Project code can be found at: %s\n", GitURL)
		fmt.Printf("Build Info:\n")
		fmt.Printf("  Built on %s\n", BuildDate)
		fmt.Printf("  Built from %s\n", BuildURL)
		os.Exit(0)
	}

	configPath, configName := getConfigPathName()

	if resetConf {
		fmt.Println("resetting config to defaults")
		err := createConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	var a app.Config
	err := gcfg.ReadFileInto(&a, configPath+configName)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			createConfig()
		} else {
			fmt.Println(err)
			os.Exit(2)
		}
	}

	// TODO: Verify config settings are valid.

	a.Run()
}

func getConfigPathName() (string, string) {
	usr, _ := user.Current()
	return usr.HomeDir + "/.config/flappy-disk/", "app.conf"
}

func createConfig() error {
	path, name := getConfigPathName()
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return err
	}

	f, err := os.Create(path + name)
	defer f.Close()
	if err != nil {
		return err
	}

	configData, err := gen.Asset("assets/default.conf")
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = f.Write(configData)
	if err != nil {
		return err
	}

	f.Sync()

	return nil
}
