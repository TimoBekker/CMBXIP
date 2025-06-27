package config

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var config = &Configuration{}
var nonFlagArgs []string
var leaveWindowsConsole bool
var configpath, panicdumppath string
var programversion string
var debugLogger *log.Logger

func init() {
	// read flags
	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagSet.BoolVar(&leaveWindowsConsole, "winconsole", false, "leave windows console")
	debug := flagSet.Bool("debug", false, "enable debug")
	debugfile := flagSet.String("debugfile", "-", "debug log file path, \"-\" for stdout")
	flagSet.StringVar(&configpath, "config", "config.json", "config path")
	flagSet.StringVar(&panicdumppath, "panicdump", "panic.log", "panic dump path")
	flagSetUsage := flagSet.Usage
	flagSet.Usage = func() {
		fmt.Printf("CompanyMediaBitrixImportProgram v%s. ", programversion)
		flagSetUsage()
	}
	flagSet.Parse(os.Args[1:])
	nonFlagArgs = flagSet.Args()

	//create debug logger
	logWriter := io.Discard
	if *debug {
		logWriter = os.Stdout
		if *debugfile != "-" {
			file, err := os.OpenFile(*debugfile, os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				panic(err)
			}
			logWriter = file
		}
	}
	debugLogger = log.New(logWriter, "[DEBUG] ", log.LstdFlags)

	// create configuration if doesn't exists and exit
	if _, err := os.Stat(configpath); os.IsNotExist(err) {
		config = &Configuration{
			CompanyMedia: &CompanyMedia{},
			Bitrix: &Bitrix{
				TaskTitleFormat:       "{{.Type}} № {{.RegNumPrefix}}{{.RegNumber}}{{.RegNumSuffix}} \"{{.Title}}\" от {{.RegDate}}{{if .Correspondent.Organization.FullName}} - {{.Correspondent.Organization.FullName}}{{end}}",
				TaskDescriptionFormat: "{{.URL}}",
			},
		}
		if err = Save(); err != nil {
			panic(err)
		}
		log.Printf("example config file %s created", configpath)
		os.Exit(1)
	}

	data, err := os.ReadFile(configpath)
	if err != nil {
		panic(err)
	}
	jsonDecoder := json.NewDecoder(bytes.NewReader(data))
	jsonDecoder.DisallowUnknownFields()
	if err = jsonDecoder.Decode(config); err != nil {
		panic(err)
	}

	if config.CompanyMedia == nil {
		config.CompanyMedia = &CompanyMedia{}
	}
	if config.Bitrix == nil {
		config.Bitrix = &Bitrix{}
	}
}

func Save() error {
	toSaveConfig := *config
	if !toSaveConfig.CompanyMedia.SaveAuth {
		toSaveConfig.CompanyMedia.APIEntry = ""
		toSaveConfig.CompanyMedia.Auth = ""
	}
	if !toSaveConfig.Bitrix.SaveInWebHook {
		toSaveConfig.Bitrix.InWebHook = ""
	}
	data, err := json.MarshalIndent(toSaveConfig, "", "    ")
	if err != nil {
		return err
	}
	if err = os.WriteFile(configpath, data, 0600); err != nil {
		return err
	}
	return nil
}

func NonFlagArgs() []string { return nonFlagArgs }

func LeaveWindowsConsole() bool { return leaveWindowsConsole }

func DebugLogger() *log.Logger { return debugLogger }

func Version() string { return programversion }

func PanicDumpPath() string { return panicdumppath }
