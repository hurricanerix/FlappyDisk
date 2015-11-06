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
	"os"
	"os/user"

	"github.com/hurricanerix/FlappyDisk/app"
	"gopkg.in/gcfg.v1"
)

func getConfigFileName() string {
	usr, _ := user.Current()
	return usr.HomeDir + "/.config/flappy-disk/app.conf"
}

func main() {
	configFile := getConfigFileName()

	var a app.Config
	err := gcfg.ReadFileInto(&a, configFile)
	if err != nil {
		// TODO: If does not exist, create it from default.conf instead of exiting.
		println(err)
		os.Exit(1)
	}

	// TODO: Verify config settings are valid.

	a.Run()
}
