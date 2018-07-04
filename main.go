package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/barryz/goci/config"
	"github.com/barryz/goci/job"
	"github.com/barryz/goci/util"
)

var (
	configFile = flag.String("c", "", "config file, YAML")
	showBuild  = flag.Bool("v", false, "show build version")

	// Build information
	Build = ""

	cfg *config.Config
)

func init() {
	flag.Parse()
}

func main() {
	if *showBuild {
		fmt.Printf("Build Info: %s\n", Build)
		os.Exit(0)
	}

	if *configFile != "" {
		err := config.InitConfig(*configFile)
		if err != nil {
			fmt.Println("Config Error: ", err)
			os.Exit(1)
		}
	} else {
		if err := config.InitDefaultConfig(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	cfg = config.DefaultConfig
	jobs := []job.Job{}

	cwd, _ := os.Getwd()
	fmt.Println("[goci] Current Work Dir: ", cwd)

	if len(cfg.RealPkgs()) != 0 {
		if !inArray(cfg.Skips, "fmt") {
			jobs = append(jobs, job.GoFmtJob("gofmt"))
		}
		if !inArray(cfg.Skips, "lint") {
			jobs = append(jobs, job.GoLintJob("lint"))
		}
		if !inArray(cfg.Skips, "vet") {
			jobs = append(jobs, job.GoVetJob("govet"))
		}
	}

	if cfg.Build != "" {
		jobs = append(jobs, job.GoBuildJob("build"))
	}
	if cfg.Test != "" {
		jobs = append(jobs, job.NewTestJob("test"))
	}
	if cfg.Race != nil && cfg.Race.Main != "" {
		jobs = append(jobs, job.GoRaceJob("race"))
	}

	for _, v := range jobs {
		jobTypeName := v.Name()
		fmt.Println("")
		util.Carve("Job " + jobTypeName + " Start")

		msg, err := v.Do()
		if err != nil {
			fmt.Printf("[failure] Job %s Error: %s\n", jobTypeName, err.Error())
			util.Carve("Job " + jobTypeName + " Finish")
			if v.IsFailTerminate() {
				break
			}
		} else {
			fmt.Printf("[success] Job %s Success %s\n", jobTypeName, msg)
			util.Carve("Job " + jobTypeName + " Finish")
		}
	}
}

func inArray(arr []string, v string) bool {
	for _, str := range arr {
		if v == str {
			return true
		}
	}
	return false
}
