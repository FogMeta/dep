package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/FogMeta/dep/akash"
	"github.com/FogMeta/dep/docker"
	"github.com/FogMeta/dep/lagrange"
	"github.com/docker/docker/api/types/registry"
	"github.com/urfave/cli/v2"
)

const (
	defaultConf = "dep.conf"
)

func main() {
	log.SetFlags(0)
	app := &cli.App{
		Name:     "dep",
		Flags:    []cli.Flag{},
		Commands: []*cli.Command{initCmd, buildCmd, createAccountCmd, deployCmd, statusCmd, closeCmd},
		Usage:    "A tool to deploy the cross-platform applications",
	}

	if err := app.Run(os.Args); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

var closeCmd = &cli.Command{
	Name:  "close",
	Usage: "close the deployment",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "name",
			Usage: "account name",
		},
		&cli.StringFlag{
			Name:  "dseq",
			Usage: "deployment dseq",
		},
	},
	Action: func(ctx *cli.Context) (err error) {
		name := ctx.String("name")
		dseq := ctx.String("dseq")
		if name == "" && dseq == "" && ctx.Args().Len() >= 2 {
			name = ctx.Args().First()
			dseq = ctx.Args().Get(1)
		}
		if name == "" {
			return errors.New("account name is required")
		}
		if dseq == "" {
			return errors.New("deployment dseq is required")
		}
		return akash.CloseDeployment(name, dseq)
	},
}

var statusCmd = &cli.Command{
	Name:  "status",
	Usage: "confirm status of the lease",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "name",
			Usage: "account name",
		},
		&cli.StringFlag{
			Name:  "dseq",
			Usage: "deployment dseq",
		},
		&cli.StringFlag{
			Name:  "provider",
			Usage: "provider address",
		},
	},
	Action: func(ctx *cli.Context) (err error) {
		name := ctx.String("name")
		dseq := ctx.String("dseq")
		provider := ctx.String("provider")
		if name == "" && dseq == "" && provider == "" && ctx.Args().Len() >= 3 {
			name = ctx.Args().First()
			dseq = ctx.Args().Get(1)
			provider = ctx.Args().Get(2)
		}
		if name == "" {
			return errors.New("account name is required")
		}
		if dseq == "" {
			return errors.New("deployment dseq is required")
		}
		if provider == "" {
			return errors.New("provider is required")
		}
		return akash.StatusLease(name, dseq, provider)
	},
}

var createAccountCmd = &cli.Command{
	Name:  "create-account",
	Usage: "create a new account",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "name",
			Usage: "account name",
		},
	},
	Action: func(ctx *cli.Context) (err error) {
		name := ctx.String("name")
		if name == "" && ctx.Args().Len() > 0 {
			name = ctx.Args().First()
		}
		if name == "" {
			return errors.New("account name is required")
		}
		return akash.CreateAccount(name)
	},
}

var deployCmd = &cli.Command{
	Name:  "deploy",
	Usage: "deploy SDL file",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "name",
			Usage: "account name created after 'create-account' cmd ",
		},
		&cli.StringFlag{
			Name:  "file",
			Usage: "deployment file path",
		},
	},
	Action: func(ctx *cli.Context) (err error) {
		name := ctx.String("name")
		path := ctx.String("file")
		if name == "" && path == "" && ctx.Args().Len() >= 2 {
			name = ctx.Args().First()
			path = ctx.Args().Get(1)
		}
		if name == "" {
			return errors.New("account name is required")
		}
		if path == "" {
			return errors.New("deployment file path is required")
		}
		return akash.Deploy(name, path)
	},
}

var initCmd = &cli.Command{
	Name:  "init",
	Usage: "init conf file",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "path",
			Usage: "conf save directory",
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
			Name:  "url",
			Usage: "lagrange space url",
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
		url := ctx.String("url")
		if url == "" && ctx.Args().Len() > 0 {
			url = ctx.Args().First()
		}
		if url == "" {
			return errors.New("url is required")
		}
		image, err := downloadAndBuild(ctx, conf, url)
		if err != nil {
			return
		}
		log.Printf("image: %s build successfully\n", image)
		return nil
	},
}

func downloadAndBuild(ctx *cli.Context, conf *Config, url string) (image string, err error) {
	path, err := lagrange.DownloadSpace(url, conf.WorkDir)
	if err != nil {
		return
	}
	image = filepath.Base(path)
	if conf != nil && conf.Registry != nil && conf.Registry.UserName != "" {
		image = fmt.Sprintf("%s/%s", conf.Registry.UserName, image)
	}
	dockerService := docker.NewDockerService()
	if err = dockerService.BuildImage(path, image); err != nil {
		log.Println("Error building Docker image: ", err)
		return "", err
	}
	if conf != nil && conf.Registry != nil && conf.Registry.UserName != "" {
		log.Println("push image failed: ", err)
		reg := conf.Registry
		err = dockerService.PushImage(image, &registry.AuthConfig{
			Username:      reg.UserName,
			Password:      reg.Password,
			ServerAddress: reg.ServerAddress,
		})
		if err != nil {
			log.Println("push image failed: ", err)
			return
		}
	}
	return
}

type Config struct {
	WorkDir  string    `toml:"work_dir"`
	Registry *Registry `toml:"registry"`
}

type Registry struct {
	ServerAddress string `toml:"server_address"`
	UserName      string `toml:"user_name"`
	Password      string `toml:"password"`
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
