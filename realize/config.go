package realize

import (
	"os"
	"gopkg.in/yaml.v2"
	"errors"
	"gopkg.in/urfave/cli.v2"
	"io/ioutil"
	"github.com/fatih/color"
	"fmt"
)

type Config struct {
	file string `yaml:"app_file,omitempty"`
	Version string `yaml:"version,omitempty"`
	Projects []Project
}

type Project struct {
	Name string `yaml:"app_name,omitempty"`
	Run bool `yaml:"app_run,omitempty"`
	Build bool `yaml:"app_build,omitempty"`
	Main string `yaml:"app_main,omitempty"`
	Watcher Watcher `yaml:"app_watcher,omitempty"`
}

type Watcher struct{
	Before []string `yaml:"before,omitempty"`
	After []string `yaml:"after,omitempty"`
	Paths []string `yaml:"paths,omitempty"`
	Exts []string `yaml:"exts,omitempty"`
}

// Default value
func New(params *cli.Context) *Config{
	return &Config{
		file: "realize.config.yaml",
		Version: "1.0",
		Projects: []Project{
			{
				Name: params.String("name"),
				Main: params.String("main"),
				Run: params.Bool("run"),
				Build: params.Bool("build"),
				Watcher: Watcher{
				Paths: []string{"/"},
				Exts: []string{"go"},
				},
			},
		},
	}
}

// check for duplicates
func Duplicates(value Project, arr []Project) bool{
	for _, val := range arr{
		if value.Main == val.Main || value.Name == val.Name{
			return true
		}
	}
	return false
}

// Remove duplicate projects
func (h *Config) Clean() {
	arr := h.Projects
	for key, val := range arr {
		 if Duplicates(val, arr[key+1:]) {
			 h.Projects = append(arr[:key], arr[key+1:]...)
			 break
		 }
	}
}

// Check, Read and remove duplicates from the config file
func (h *Config) Read() error{
	if file, err :=  ioutil.ReadFile(h.file); err == nil{
		err = yaml.Unmarshal(file, h)
		if err == nil {
			h.Clean()
		}
		return err
	}else{
		return err
	}
}

// write and marshal
func (h *Config) Write() error{
	y, err := yaml.Marshal(h)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(h.file, y, 0755)
}

// Create config yaml file
func (h *Config) Create(params *cli.Context) error{
	if h.Read() != nil {
		if err := h.Write(); err != nil {
			os.Remove(h.file)
			return err
		}else{
			return err
		}
	}
	return errors.New("The config file already exists, check for realize.config.yaml")
}

// Add another project
func (h *Config) Add(params *cli.Context) error{
	if err := h.Read(); err == nil {
		new := Project{
			Name: params.String("name"),
			Main: params.String("main"),
			Run: params.Bool("run"),
			Build: params.Bool("build"),
			Watcher: Watcher{
				Paths: []string{"/"},
				Exts: []string{"go"},
			},
		}
		if Duplicates(new, h.Projects) {
			return errors.New("There is already one project with same path or name")
		}
		h.Projects = append(h.Projects, new)
		return h.Write()
	}else{
		return err
	}
}

// remove a project in list
func (h *Config) Remove(params *cli.Context) error{
	if err := h.Read(); err == nil {
		for key, val := range h.Projects {
			if params.String("name") == val.Name {
				h.Projects = append(h.Projects[:key], h.Projects[key+1:]...)
				return h.Write()
			}
		}
		return errors.New("No project found")
	}else{
		return err
	}
}

// List of projects
func (h *Config) List() error{
	if err := h.Read(); err == nil {
		green := color.New(color.FgGreen, color.Bold).SprintFunc()
		greenl := color.New(color.FgHiGreen).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()
		for _, val := range h.Projects {
			fmt.Println(green("|"), green(val.Name))
			fmt.Println(greenl("|"),"\t", green("Main Path:"), red(val.Main))
			fmt.Println(greenl("|"),"\t", green("Run:"), red(val.Run))
			fmt.Println(greenl("|"),"\t", green("Build:"), red(val.Build))
			fmt.Println(greenl("|"),"\t", green("Watcher:"))
			fmt.Println(greenl("|"),"\t\t", green("After:"), red(val.Watcher.After))
			fmt.Println(greenl("|"),"\t\t", green("Before:"), red(val.Watcher.Before))
			fmt.Println(greenl("|"),"\t\t", green("Extensions:"), red(val.Watcher.Exts))
			fmt.Println(greenl("|"),"\t\t", green("Paths:"), red(val.Watcher.Paths))
		}
		return nil
	}else{
		return err
	}
}

