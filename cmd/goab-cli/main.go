package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	bundle "github.com/GizmoOAO/go-asset-bundle"
	"github.com/urfave/cli/v2"
)

type fileInfo struct {
	Path    string `json:"path"`
	ModTime int64  `json:"time"`
	Size    int64  `json:"size"`
	At      int64  `json:"-"`
}

func main() {
	app := &cli.App{
		Name:     "goab-cli",
		Usage:    "go-asset-bundle cli application",
		Commands: []*cli.Command{pack(), list(), extract(), extractFile()},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func pack() *cli.Command {
	return &cli.Command{
		Name:    "pack",
		Usage:   "Create an AssetBundle",
		Aliases: []string{"p"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    "Output bundle file",
				Required: true,
			},
			&cli.UintFlag{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Bundle version",
				Value:   1000,
			},
		},
		ArgsUsage: "BundleFolder",
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() != 1 {
				return errors.New("\"BundleFolder\" is required")
			}
			output := ctx.String("output")
			folder := ctx.Args().Get(0)
			version := ctx.Uint("version")
			_, err := os.Stat(output)
			if err == nil {
				_ = os.Remove(output)
			}
			ab, err := bundle.NewAssetBundle(output)
			if err != nil {
				return err
			}
			defer ab.Close()
			return ab.Bundle(folder, uint16(version))
		},
	}
}

func list() *cli.Command {
	return &cli.Command{
		Name:      "list",
		Aliases:   []string{"l"},
		Usage:     "List files of AssetBundle",
		ArgsUsage: "BundleFile",
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() != 1 {
				return errors.New("\"BundleFile\" is required")
			}
			file := ctx.Args().Get(0)
			ab, err := bundle.OpenAssetBundle(file)
			if err != nil {
				return err
			}
			defer ab.Close()
			var files []fileInfo
			for i := 0; i < len(ab.Files); i++ {
				info := fileInfo(ab.Files[i])
				files = append(files, info)
			}
			data, err := json.MarshalIndent(&files, "", "\t")
			if err != nil {
				return err
			}
			fmt.Println(string(data))
			return nil
		},
	}
}

func extract() *cli.Command {
	return &cli.Command{
		Name:    "extract",
		Usage:   "Extract AssetBundle",
		Aliases: []string{"e"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    "Output folder",
				Required: true,
			},
		},
		ArgsUsage: "BundleFile",
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() != 1 {
				return errors.New("\"BundleFile\" is required")
			}
			file := ctx.Args().Get(0)
			ab, err := bundle.OpenAssetBundle(file)
			if err != nil {
				return err
			}
			defer ab.Close()
			for _, info := range ab.Files {
				filename, err := filepath.Abs(filepath.Join(ctx.String("output"), info.Path))
				if err != nil {
					return err
				}
				if err = os.MkdirAll(filepath.Dir(filename), 0666); err != nil {
					return err
				}
				f2, err := ab.Open(info.Path)
				if err != nil {
					return err
				}
				data, err := ioutil.ReadAll(f2)
				if err != nil {
					return err
				}
				if err = ioutil.WriteFile(filename, data, 0666); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func extractFile() *cli.Command {
	return &cli.Command{
		Name:    "extract-file",
		Usage:   "Extract one file from AssetBundle",
		Aliases: []string{"ef"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "path",
				Aliases:  []string{"p"},
				Usage:    "Asset path",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    "Output filename",
				Required: true,
			},
		},
		ArgsUsage: "BundleFile",
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() != 1 {
				return errors.New("\"BundleFile\" is required")
			}
			file := ctx.Args().Get(0)
			filename, err := filepath.Abs(ctx.String("output"))
			if err != nil {
				return err
			}
			if err = os.MkdirAll(filepath.Dir(filename), 0666); err != nil {
				return err
			}
			ab, err := bundle.OpenAssetBundle(file)
			if err != nil {
				return err
			}
			defer ab.Close()
			f, err := ab.Open(ctx.String("path"))
			if err != nil {
				return err
			}
			data, err := ioutil.ReadAll(f)
			if err != nil {
				return err
			}
			if err = ioutil.WriteFile(filename, data, 0666); err != nil {
				return err
			}
			return nil
		},
	}
}
