package main

import (
	"bytes"
	"text/template"
)

type BitrixTaskTitleAndDescriptionFormat struct {
	URL   string
	Title string
	Type  string

	RegDate string

	RegNumPrefix string
	RegNumber    int
	RegNumSuffix string

	Correspondent struct {
		Organization struct {
			FullName string
		}
	}
}

func parseTemplateFromStruct(templ string, data interface{}) (string, error) {
	t, err := template.New("").Parse(templ)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer([]byte{})
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
