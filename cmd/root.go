package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"grafana_to_datadog/dashboard"
	"grafana_to_datadog/grafana"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	outputPath string
	inputPath  string
	debug      bool
)

func init() {
	rootCmd.Flags().StringVarP(&inputPath, "input", "i", "", "name of the input file or directory")
	rootCmd.Flags().StringVarP(&outputPath, "output", "o", "", "name of the output directory")
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "ennable debug mode")
	rootCmd.Flag("debug").NoOptDefVal = "true"

	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("output")

	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
	// log.SetReportCaller(true)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getJSONFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".json" {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

var rootCmd = &cobra.Command{
	Use:   "grafana_to_datadog",
	Short: "Convert Grafana Dashboard to Datadog",
	Run: func(cmd *cobra.Command, args []string) {

		if debug {
			log.SetLevel(log.DebugLevel)
		}

		if _, err := os.Stat(outputPath); errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(outputPath, os.ModePerm)
			check(err)
		}

		fileInfo, err := os.Stat(inputPath)
		if err != nil {
			check(err)
		}

		var files []string

		if fileInfo.IsDir() {
			files, err = getJSONFiles(inputPath)
			check(err)
		} else {
			files = []string{inputPath}
		}
		for _, file := range files {
			directoryName, fileName := path.Split(file)

			outputDirectory := outputPath

			fStructure := strings.Split(directoryName, "/")
			if fStructure[0] == ".." || fStructure[0] == "." {
				fStructure = fStructure[1:]
			}

			if len(fStructure) > 1 {
				fStructure[0] = outputPath
				outputDirectory = path.Join(fStructure...)
			}

			if _, err := os.Stat(outputDirectory); errors.Is(err, os.ErrNotExist) {
				err := os.MkdirAll(outputDirectory, os.ModePerm)
				check(err)
			}

			contextLogger := log.WithFields(log.Fields{
				"dashboard": fileName,
			})

			dat, err := os.ReadFile(file)
			check(err)
			dash := &grafana.Dashboard{}
			json.Unmarshal(dat, dash)
			res, err := dashboard.ConvertDashboard(dash, contextLogger).MarshalJSON()
			check(err)
			output := path.Join(outputDirectory, fileName)
			err = os.WriteFile(output, res, 0644)
			check(err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
