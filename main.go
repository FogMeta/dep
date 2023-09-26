package main

import (
	"bytes"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/FogMeta/dep/docker"
	"github.com/FogMeta/dep/lagrange"
	"github.com/docker/docker/api/types/registry"
	"github.com/urfave/cli/v2"
)

const (
	defaultConf = "dep.conf"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	app := &cli.App{
		Name:     "dep",
		Flags:    []cli.Flag{},
		Commands: []*cli.Command{initCmd, buildCmd},
		Usage:    "A tool to deploy the cross-platform applications",
	}

	if err := app.Run(os.Args); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

var initCmd = &cli.Command{
	Name:  "init",
	Usage: "init conf file",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "path",
			Usage: "conf save path",
		},
	},
	Action: func(ctx *cli.Context) (err error) {
		path := ctx.String("path")
		if path == "" {
			path = defaultConf
		} else {
			path = filepath.Join(path, defaultConf)
		}
		b, err := tomlMarshal(Config{
			WorkDir:  ".",
			Registry: &Registry{},
		})
		if err != nil {
			return
		}
		if err = os.WriteFile(path, b, 0666); err != nil {
			return
		}
		log.Println("config file saved to ", path)
		return
	},
}

var buildCmd = &cli.Command{
	Name:  "build",
	Usage: "build from lagrange url",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "conf",
			Usage: "conf file path",
		},
		&cli.StringFlag{
			Name:  "work-dir",
			Usage: "download directory",
		},
		&cli.StringFlag{
			Name:     "url",
			Usage:    "lagrange space url",
			Required: true,
		},
	},
	Action: func(ctx *cli.Context) (err error) {
		confPath := ctx.String("conf")
		if confPath == "" {
			confPath = defaultConf
		}
		_, err = os.Stat(confPath)
		if err != nil {
			return errors.New("not set conf path, you can run init to create it")
		}
		conf, err := Init(confPath)
		if err != nil {
			return
		}
		image, err := downloadAndBuild(ctx, conf)
		if err != nil {
			return
		}
		log.Printf("image: %s build successfully\n", image)
		return nil
	},
}

func downloadAndBuild(ctx *cli.Context, conf *Config) (image string, err error) {
	path, err := lagrange.DownloadSpace(ctx.String("url"), conf.WorkDir)
	if err != nil {
		return
	}
	image = filepath.Base(path)
	dockerService := docker.NewDockerService()
	if err = dockerService.BuildImage(path, image); err != nil {
		log.Println("Error building Docker image: ", err)
		return "", err
	}
	if conf != nil && conf.Registry != nil && conf.Registry.UserName != "" {
		reg := conf.Registry
		dockerService.PushImage(image, &registry.AuthConfig{
			Username:      reg.UserName,
			Password:      reg.Password,
			ServerAddress: reg.ServerAddress,
		})
	}
	return
}

type Config struct {
	WorkDir  string    `yaml:"work_dir"`
	Registry *Registry `yaml:"registry"`
}

type Registry struct {
	ServerAddress string `yaml:"server_address"`
	UserName      string `yaml:"user_name"`
	Password      string `yaml:"password"`
}

func tomlMarshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Init(path string) (*Config, error) {
	var conf Config
	_, err := toml.DecodeFile(path, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
