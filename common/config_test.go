package common

import (
	"fmt"
	"testing"
)

func TestConfig_WriteFile(t *testing.T) {
	CreateConfigFile()
	//CreateLocalConfigFile()
}

func TestParseConfigFile(t *testing.T) {
	conf := ParseConfig("config_127.0.0.1_5010_5")
	fmt.Println(conf)
}
