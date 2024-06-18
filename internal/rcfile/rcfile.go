package rcfile

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

const (
	defaultMode = 0644
	urlTemplate = "https://%s/registry.terraform.io/hashicorp/random/index.json"
	rcFileName  = "/.terraformrc"
)

var (
	//go:embed terraformrc.tpl
	terrafromrcTemplateStr string
	terrafromrcTemplate    = template.Must(template.New("terrafromrc").Parse(terrafromrcTemplateStr))
)

func TemplFile() *template.Template {
	return terrafromrcTemplate
}

func writeFile(mirrorList []string) error {
	rcfilepath := os.Getenv("HOME") + rcFileName
	for _, elem := range mirrorList {
		if checkMirror(elem) {
			if err := os.WriteFile(rcfilepath, templateFile(elem), defaultMode); err != nil {
				return err
			}
			break
		}
	}
	return nil
}

func CreateFile(mirrorList []string) error {
	err := writeFile(mirrorList)
	return err
}

func checkMirror(mirr string) bool {
	link := fmt.Sprintf(urlTemplate, mirr)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, link, nil)
	if err != nil {
		return false
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func templateFile(site string) []byte {
	terraformrcBuffer := new(bytes.Buffer)
	template := TemplFile()
	err := template.Execute(terraformrcBuffer, site)
	if err != nil {
		log.Panic(err.Error())
	}
	return terraformrcBuffer.Bytes()
}
