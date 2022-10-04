package cmd

import (
	"fmt"
	"os"
)

/*
#!/usr/bin/env bash
version=0.1.0
go build -ldflags="-X 'github.com/kore3lab/kore/cmd.BuildTime=$(date -u +%FT%T%Z)' -X 'github.com/kore3lab/kore/cmd.BuildVersion=$version'" .
*/
var (
	BuildVersion string = ""
	BuildTime    string = ""
)

type Options struct {
	OutStream *os.File // output stream
	Values    []string
	//ConfigFile string // config file
	//Output     string // output format (json/yaml)
	Filename    string // file
	ProfileName string
	Namespace   string
}

func (o *Options) GetFilename() string {
	return o.Filename
}

func (o *Options) Println(format string, params ...interface{}) {
	msg := fmt.Sprintf(format+"\n", params...)
	if o.OutStream != nil {
		o.OutStream.WriteString(msg)
	} else {
		os.Stdout.WriteString(msg)
	}
}
func (o *Options) PrintlnError(err error) {
	o.Println("%+v\n", err)
}
