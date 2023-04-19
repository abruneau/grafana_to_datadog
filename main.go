package main

import (
	"encoding/json"
	"fmt"
	"grafana_to_datadog/dashboard"
	"grafana_to_datadog/grafana"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func init() {
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

func main() {

	if len(os.Args) < 2 {
		fmt.Println("path to a dashboard definition is required")
		os.Exit(1)
	}

	dPath := os.Args[1]

	fileInfo, err := os.Stat(dPath)
	if err != nil {
		check(err)
	}

	if fileInfo.IsDir() {
		files, err := getJSONFiles(dPath)
		check(err)
		for _, file := range files {
			pathElements := strings.Split(file, "/")
			directory := pathElements[len(pathElements)-2]
			fileName := pathElements[len(pathElements)-1]

			contextLogger := log.WithFields(log.Fields{
				"dashboard": fileName,
			})

			dat, err := os.ReadFile(file)
			check(err)
			dash := &grafana.Dashboard{}
			json.Unmarshal(dat, dash)
			res, err := dashboard.ConvertDashboard(dash, contextLogger).MarshalJSON()
			check(err)
			os.Mkdir(fmt.Sprintf("./output/%s", directory), os.ModePerm)
			outputPath := fmt.Sprintf("./output/%s/%s", directory, fileName)
			err = os.WriteFile(outputPath, res, 0644)
			check(err)
		}
	} else {
		contextLogger := log.WithFields(log.Fields{
			"dashboard": dPath,
		})

		dat, err := os.ReadFile(dPath)
		check(err)

		dash := &grafana.Dashboard{}
		json.Unmarshal(dat, dash)

		res, err := dashboard.ConvertDashboard(dash, contextLogger).MarshalJSON()
		check(err)
		err = os.WriteFile("./dd_dashboard.json", res, 0644)
		check(err)
	}
}
