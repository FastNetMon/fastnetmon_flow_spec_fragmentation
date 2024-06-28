package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/fastnetmon/fastnetmon-go"
)

var fast_logger = log.New(os.Stderr, fmt.Sprintf(" %d ", os.Getpid()), log.LstdFlags)

type Configuration struct {
	Log_path string `json:"log_path" fastnetmon_type:"string"`
}

func main() {
	conf := Configuration{Log_path: "/var/log/fastnetmon/fastnetmon_flow_spec_fragmentation.log"}

	log_file, err := os.OpenFile(conf.Log_path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fast_logger.Fatalf("Cannot open log file: %v", err)
	}

	defer log_file.Close()

	multi_writer := io.MultiWriter(os.Stdout, log_file)

	fast_logger.SetOutput(multi_writer)

	fast_logger.Printf("Prepared to read data from stdin")
	stdin_data, err := ioutil.ReadAll(os.Stdin)

	if err != nil {
		fast_logger.Fatal("Cannot read data from stdin")
	}

	callback_data := fastnetmon.CallbackDetails{}

	fast_logger.Printf("Callback raw data: %s", stdin_data)

	err = json.Unmarshal([]byte(stdin_data), &callback_data)

	if err != nil {
		fast_logger.Printf("Raw data: %s", stdin_data)
		fast_logger.Fatalf("Cannot unmarshal data: %v", err)
	}

	// Scope can be per host or total hostgroups
	alert_scope := callback_data.AlertScope

	if alert_scope != "host" {
		fast_logger.Fatalf("Unknown scope: %s Only host scope is supported", alert_scope)
	}

	fast_logger.Printf("Flow Spec rules: %+v", callback_data.FlowSpecRules)
}
