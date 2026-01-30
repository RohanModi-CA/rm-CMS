package config

import (
	"cms/misc"
	"github.com/BurntSushi/toml"
	"os"
)

type configTOML struct {
	HeaderPath string
}

type Config struct {
	HeadContents string
}

func path_from_user(inp_path string) string {
	// If a relative path, prepends the relative path to user.

	if len(inp_path) > 0 && inp_path[0] != '\\' {
		inp_path = "../config/user/" + inp_path
	}
	return inp_path
}

func ProcessConfig() Config {
	// For now we assume we're running from src/
	config_toml_path := "../config/config.toml"

	data_bytes, err_read := os.ReadFile(config_toml_path)
	misc.ErrorHandlePanic(err_read)
	data := string(data_bytes)

	var confTOML configTOML

	_, err := toml.Decode(data, &confTOML)
	misc.ErrorHandlePanic(err)
	confTOML.HeaderPath = path_from_user(confTOML.HeaderPath)

	var conf Config

	header_bytes, err_header_read := os.ReadFile(confTOML.HeaderPath)
	misc.ErrorHandlePanic(err_header_read)

	conf.HeadContents = string(header_bytes)

	return conf
}
