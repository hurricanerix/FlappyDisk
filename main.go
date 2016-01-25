// Copyright 2015-2016 Richard Hawkins
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
	"log"
	"os"
	"runtime"

	"github.com/hurricanerix/FlappyDisk/game"
	"github.com/hurricanerix/FlappyDisk/gen"
	"github.com/hurricanerix/shade/display"
	sgen "github.com/hurricanerix/shade/gen"
	"github.com/hurricanerix/shade/splash"
)

//go:generate ./gen_build_info.sh

var gitURL = "https://github.com/hurricanerix/FlappyDisk"

var (
	resetConf bool
	cheat     bool
	nosplash  bool
	version   bool
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func init() {
	flag.BoolVar(&version, "version", false, "print version and build info.")
	flag.BoolVar(&cheat, "cheat", false, "cheat mode.")
	flag.BoolVar(&nosplash, "nosplash", false, "don't show splash screen.")
}

func main() {
	flag.Parse()

	if version {
		fmt.Printf("FlappyDisk Copyright 2015-2016 Richard Hawkins\n")
		fmt.Printf("Licensed under the Apache License, Version 2.0\n")
		fmt.Printf("Project code can be found at: %s\n", gitURL)
		fmt.Printf("Build Info:\n")
		fmt.Printf("TODO: add build info back")
		fmt.Printf("  %s\n", gen.Version)
		fmt.Printf("  built on %s\n", gen.BuildDate)
		fmt.Printf("  built from %s/commit/%s\n", gitURL, gen.BuildHash)
		fmt.Printf("  built using Shade SDK\n")
		fmt.Printf("	%s\n", sgen.Version)
		fmt.Printf("	built from %s/commit/%s\n", sgen.GitURL, sgen.Hash)
		os.Exit(0)
	}

	config := game.Config{Cheat: cheat}

	screen, err := display.SetMode("FlappyDisk", 640, 480)
	if err != nil {
		log.Fatalln("failed to set display mode:", err)
	}

	g, err := game.New(screen)
	if err != nil {
		log.Fatalln("failed to create game:", err)
	}

	if !nosplash {
		splash.Main(screen)
	}
	g.Main(screen, config)
}
