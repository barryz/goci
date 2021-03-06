package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	defaultBuildCommand = "go build ."
	defaultTestCommand  = "go test ."
)

var (
	defaultExcludes = []string{"vendor"}
)

// Config a ci linter configuration object abstraction.
type Config struct {
	Build string `yaml:"build"`
	Test  string `yaml:"test"`
	// same as Pkgs
	Dirs     []string `yaml:"dirs"`
	Pkgs     []string `yaml:"pkgs"`
	Race     *Race    `yaml:"race"`
	Lint     *Lint    `yaml:"lint"`
	Excludes []string `yaml:"excludes"`
	realPkgs []string
	Skips    []string `yaml:"skips"`
}

// Lint config item for go lint command.
type Lint struct {
	IgnoreNoCommentError bool `yaml:"ignore_no_comment_error"`
}

// Race config item for go race command.
type Race struct {
	Main    string `yaml:"main"`
	MainCMD string `yaml:"main_cmd"`
	Script  string `yaml:"script"`
	Timeout int    `yaml:"timeout"`
}

// RealPkgs get the realistic packages.
func (c *Config) RealPkgs() []string {
	if c.realPkgs == nil {
		excludes := append(c.Excludes, defaultExcludes...)
		fmt.Printf("\n!!!!! Excludes:\n    %v\n\n", excludes)
		err := c.walk(".", excludes)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return c.realPkgs
}

// InExcludes check the path whether in excluded directive.
func (c *Config) InExcludes(path string) bool {
	excludes := append(c.Excludes, defaultExcludes...)
	for _, e := range excludes {
		if strings.HasPrefix(path, e) {
			return true
		}
	}
	return false
}

func (c *Config) walk(cur string, excludes []string) error {
	files, err := ioutil.ReadDir(cur)
	if err != nil {
		return err
	}
	for _, file := range files {
		isExclude := false
		fileName := cleaPath(cur + "/" + file.Name())
		for _, e := range excludes {
			if strings.HasPrefix(fileName, e) {
				isExclude = true
				break
			}
		}
		if !isExclude {
			if strings.HasSuffix(fileName, ".go") {
				c.realPkgs = append(c.realPkgs, fileName)
			} else if file.IsDir() && !isHiddenDir(fileName) {
				err := c.walk("./"+fileName, excludes)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func cleaPath(path string) string {
	return strings.TrimPrefix(path, "./")
}

func isHiddenDir(name string) bool {
	if len(name) <= 1 {
		return false
	}
	return name[0] == '.' && name[1] != '/'
}

var configFields = []string{"build", "test", "dirs", "pkgs", "race", "raceok"}

// IsConfigField check a field whether valid.
func IsConfigField(f string) bool {
	for i := 0; i < len(configFields); i++ {
		if configFields[i] == f {
			return true
		}
	}
	return false
}

var (
	// DefaultConfig default config object.
	DefaultConfig *Config
)

// InitConfig initialize config by config file
func InitConfig(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open Config File Error:%s", err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("read Config File Error:%s", err)
	}
	cfg := Config{}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return fmt.Errorf("unmarshal Config File Contents Error:%s ", err)
	}

	if cfg.Dirs != nil || cfg.Pkgs != nil {
		return errors.New("dirs and pkgs field in goci.yml has been cancelled! Use excludes instead")
	}

	for i, v := range cfg.Excludes {
		v = strings.TrimPrefix(v, "./")
		cfg.Excludes[i] = v
	}

	if cfg.Lint == nil {
		cfg.Lint = &Lint{}
	}

	DefaultConfig = &cfg
	return nil
}

// InitDefaultConfig initialize default configuration
func InitDefaultConfig() error {
	cfg := new(Config)

	cfg.Build = defaultBuildCommand
	cfg.Test = defaultTestCommand
	cfg.Excludes = defaultExcludes
	cfg.Lint = &Lint{IgnoreNoCommentError: true}

	DefaultConfig = cfg

	return nil
}
