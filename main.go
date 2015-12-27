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
// Package main reads CLI flags and creates and runs the app.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hurricanerix/FlappyDisk/app"
	"github.com/hurricanerix/FlappyDisk/gen"
)

//go:generate ./gen_build_info.sh

var gitURL = "https://github.com/hurricanerix/FlappyDisk"

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
		fmt.Printf("Project code can be found at: %s\n", gitURL)
		fmt.Printf("Build Info:\n")
		fmt.Printf("  %s\n", gen.Version)
		fmt.Printf("  built on %s\n", gen.BuildDate)
		fmt.Printf("  built from %s/commit/%s\n", gitURL, gen.BuildHash)
		os.Exit(0)
	}

	a, err := app.New(resetConf)
	if err != nil {
		panic(err)
	}
	defer a.Terminate()

	a.Run()
}
