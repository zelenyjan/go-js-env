package main

import (
	"os"
	"strings"
	"text/template"
)

const PREFIX = "REACT_APP_"
const EnvFilename = ".env"
const ResultFilename = "build/envGo.js"
const TemplateText = `window.envGo = {
{{- range $key, $value := . }}
  {{$key}}: "{{$value}}",
{{- end }}
}`

func getEnvFromENV() map[string]string {
	// store items to map
	envs := make(map[string]string)

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if strings.HasPrefix(pair[0], PREFIX) {
			envs[pair[0]] = pair[1]
		}
	}
	return envs
}

func getEnvFromFile() map[string]string {
	envs := make(map[string]string)
	// open file
	data, err := os.ReadFile(EnvFilename)
	if err != nil {
		return envs
	}

	// split data by new line
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		// split line by "="
		pair := strings.SplitN(line, "=", 2)
		if strings.HasPrefix(pair[0], PREFIX) {
			envs[pair[0]] = pair[1]
		}
	}

	return envs
}

func writeToFile(envs map[string]string) {
	err := os.MkdirAll("build", os.ModePerm)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	// delete file if exists
	err = os.Remove(ResultFilename)
	if err != nil && !os.IsNotExist(err) {
		println(err.Error())
		os.Exit(1)
	}

	// open file
	file, err := os.OpenFile(ResultFilename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	tmpl, err := template.New("t").Parse(TemplateText)
	if err != nil {
		// close file
		println(err.Error())
		os.Exit(1)
	}

	err = tmpl.Execute(file, envs)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func main() {
	dataFromEnv := getEnvFromENV()
	dataFromFile := getEnvFromFile()

	// merge data
	// overwrite data from file with data from env -> env from env has higher priority
	for k, v := range dataFromEnv {
		dataFromFile[k] = v
	}
	writeToFile(dataFromFile)
}
