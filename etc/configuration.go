package etc

import "github.com/8micro/datastore-server/server"
import "github.com/8micro/gounits/logger"
import "gopkg.in/yaml.v2"

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	SystemConfig *Configuration = nil
)

var (
	ErrConfigFileNotFound      = errors.New("config file not found.")
	ErrConfigGenerateFailure   = errors.New("config file generated failure.")
	ErrConfigFormatInvalid     = errors.New("config file format invalid.")
	ErrConfigServerDataInvalid = errors.New("config server data invalid.")
)

// Configuration is exported
type Configuration struct {
	Version      string `yaml:"version" json:"version"`
	PidFile      string `yaml:"pidfile" json:"pidfile"`
	RetryStartup bool   `yaml:"retrystartup" json:"retrystartup"`

	Prest struct {
		DataBase     string `yaml:"database" json:"database"`
		Schema       string `yaml:"schema" json:"schema"`
		Hosts        string `yaml:"hosts" json:"hosts"`
		Jwt struct {
			Secret  string `yaml:"secret" json:"secret"`
			Expired string `yaml:"expired" json:"expired"`
		} `yaml:"jwt" json:"jwt"`
	} `yaml:"prest" json:"prest"`

	Listen struct {
		Hosts      []string `yaml:"hosts" json:"hosts"`
		EnableCors bool     `yaml:"enablecors" json:"enablecors"`
	} `yaml:"listen" json:"listen"`

	Logger struct {
		LogFile  string `yaml:"logfile" json:"logfile"`
		LogLevel string `yaml:"loglevel" json:"loglevel"`
		LogSize  int64  `yaml:"logsize" json:"logsize"`
	} `yaml:"logger" json:"logger"`
}

// New is exported
func New(file string) error {

	buf, err := readConfigurationFile(file)
	if err != nil {
		return fmt.Errorf("config error %s", err.Error())
	}

	conf := &Configuration{RetryStartup: true}
	if err := yaml.Unmarshal(buf, conf); err != nil {
		return ErrConfigFormatInvalid
	}

	if err = conf.parseEnv(); err != nil {
		return fmt.Errorf("config parse env %s", err.Error())
	}

	SystemConfig = parseDefaultParmeters(conf)
	log.Printf("[#etc#] version: %s\n", SystemConfig.Version)
	log.Printf("[#etc#] pidfile: %s\n", SystemConfig.PidFile)
	log.Printf("[#etc#] retrystartup: %s\n", strconv.FormatBool(SystemConfig.RetryStartup))
	log.Printf("[#etc#] prest: %+v\n", SystemConfig.Prest)
	log.Printf("[#etc#] listen: %+v\n", SystemConfig.Listen)
	log.Printf("[#etc#] logger: %+v\n", SystemConfig.Logger)
	return nil
}

//PidFile is exported
func PidFile() string {

	if SystemConfig != nil {
		return SystemConfig.PidFile
	}
	return ""
}

//RetryStartup is exported
func RetryStartup() bool {

	if SystemConfig != nil {
		return SystemConfig.RetryStartup
	}
	return false
}

//PrestConfig is exported
func PrestConfig() *server.PrestConfig {

	if SystemConfig != nil {
		hosts := []string{}
		prestHosts := strings.Split(SystemConfig.Prest.Hosts, ",")
		for _, hostPort := range prestHosts {
			if _, _, err := net.SplitHostPort(hostPort); err == nil {
				hosts = append(hosts, hostPort)
			}
		}

		var jwt *server.JwtConfig
		if SystemConfig.Prest.Jwt.Secret != "" {
			jwt = &server.JwtConfig {
				Secret: SystemConfig.Prest.Jwt.Secret,
				Expired: SystemConfig.Prest.Jwt.Expired,
			}
		}

		return &server.PrestConfig{
			DataBase:     SystemConfig.Prest.DataBase,
			Schema:       SystemConfig.Prest.Schema,
			Hosts:        hosts,
			Jwt: 		  jwt,
		}
	}
	return nil
}

//LoggerArgs is exported
func LoggerArgs() *logger.Args {

	if SystemConfig != nil {
		return &logger.Args{
			FileName: SystemConfig.Logger.LogFile,
			Level:    SystemConfig.Logger.LogLevel,
			MaxSize:  SystemConfig.Logger.LogSize,
		}
	}
	return nil
}

func readConfigurationFile(file string) ([]byte, error) {

	fd, err := os.OpenFile(file, os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}

	defer fd.Close()
	buf, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func parseDefaultParmeters(conf *Configuration) *Configuration {

	if conf.Logger.LogLevel == "" {
		conf.Logger.LogLevel = "info"
	}

	if conf.Logger.LogSize == 0 {
		conf.Logger.LogSize = 20971520
	}
	return conf
}
