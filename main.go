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
	Log_path    string `json:"log_path" fastnetmon_type:"string"`
	ApiUser     string `json:"api_user"`
	ApiPassword string `json:"api_password"`
	ApiHost     string `json:"api_host"`
	ApiPort     uint32 `json:"api_port"`
}

var conf Configuration

func main() {
	conf.ApiUser = "admin"
	conf.ApiPassword = "test_password"
	conf.ApiHost = "127.0.0.1"
	conf.ApiPort = 10007

	file_as_array, err := ioutil.ReadFile("/etc/fastnetmon/fastnetmon_flow_spec_fragmentation.conf")

	if err != nil {
		log.Fatalf("Could not read configuration file with error: %v", err)
	}

	// This command will override our default configuration
	err = json.Unmarshal(file_as_array, &conf)

	if err != nil {
		log.Fatalf("Could not read json configuration: %v", err)
	}

	if conf.Log_path == "" {
		conf.Log_path = "/var/log/fastnetmon/fastnetmon_flow_spec_fragmentation.log"
	}

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

	if len(callback_data.FlowSpecRules) == 0 {
		fast_logger.Fatalf("Empty list of Flow Spec rules")
	}

	// We have at least one rule. Our core does not return more then one
	flow_spec_rule := callback_data.FlowSpecRules[0]

	if len(flow_spec_rule.Protocols) == 1 && flow_spec_rule.Protocols[0] == "udp" &&
		len(flow_spec_rule.SourcePorts) == 1 && flow_spec_rule.SourcePorts[0] == 0 &&
		len(flow_spec_rule.DestinationPorts) == 1 && flow_spec_rule.DestinationPorts[0] == 0 {
		fast_logger.Printf("Matched rule action: %s", callback_data.Action)
	} else {
		fast_logger.Fatalf("We do not handle such rules, ignoring")
	}

	// Time to pass it to FastNetMon

	fast_logger.Printf("Connect to FastNetMon API")

	fastnetmon_client, err := fastnetmon.NewClient(conf.ApiHost, conf.ApiPort, conf.ApiUser, conf.ApiPassword)

	if err != nil {
		fast_logger.Fatalf("Cannot connect to client: %v", err)
	}

	_ = fastnetmon_client

	// Adjsut flow spec rule

	// Remove ports from it
	flow_spec_rule.SourcePorts = []uint{}
	flow_spec_rule.DestinationPorts = []uint{}

	// Remove packet lengths
	flow_spec_rule.PacketLengths = []uint{}

	// Add fragmentation flags
	flow_spec_rule.FragmentationFlags = []string{"is-fragment"}

	fast_logger.Printf("Complementary Flow Spec rule: %+v", flow_spec_rule)

	// Encode to JSON for checking
	encodedJSON, err := json.Marshal(flow_spec_rule)

	if err != nil {
		fast_logger.Fatalf("Cannot encode JSON: %v", err)
	}

	fast_logger.Printf("Encoded Flow Spec in JSON: %s", string(encodedJSON))
}
