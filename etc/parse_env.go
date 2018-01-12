package etc

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func (conf *Configuration) parseEnv() error {

	if pidFile := os.Getenv("PID_FILE"); pidFile != "" {
		conf.PidFile = pidFile
	}

	retryStartup := os.Getenv("RETRY_STARTUP")
	if retryStartup != "" {
		value, err := strconv.ParseBool(retryStartup)
		if err != nil {
			return fmt.Errorf("RETRY_STARTUP invalid, %s", err.Error())
		}
		conf.RetryStartup = value
	}

	//parse prest env
	if err := parsePrestEnv(conf); err != nil {
		return err
	}

	//parse API listen env
	if err := parseListenEnv(conf); err != nil {
		return err
	}

	//parse logger env
	return parseLogger(conf)
}

func parsePrestEnv(conf *Configuration) error {

	if prestDataBase := os.Getenv("PREST_DATABASE"); prestDataBase != "" {
		conf.Prest.DataBase = prestDataBase
	}

	if prestSchema := os.Getenv("PREST_SCHEMA"); prestSchema != "" {
		conf.Prest.Schema = prestSchema
	}

	if prestHosts := os.Getenv("PREST_HOSTS"); prestHosts != "" {
		conf.Prest.Hosts = prestHosts
	}

	if prestJwtSecret := os.Getenv("PREST_JWT_SECRET"); prestJwtSecret != "" {
		conf.Prest.Jwt.Secret = prestJwtSecret
	}

	if prestJwtExpired := os.Getenv("PREST_JWT_EXPIRED"); prestJwtExpired != "" {
		if _, err := time.ParseDuration(prestJwtExpired); err != nil {
			return fmt.Errorf("prest token expired invalid.")
		}
		conf.Prest.Jwt.Expired = prestJwtExpired
	}
	return nil
}

func parseListenEnv(conf *Configuration) error {

	if apiHost := os.Getenv("API_LISTEN"); apiHost != "" {
		hostIP, hostPort, err := net.SplitHostPort(apiHost)
		if err != nil {
			return fmt.Errorf("API_LISTEN invalid, %s", err.Error())
		}
		if hostIP != "" {
			if _, err := net.LookupHost(hostIP); err != nil {
				return fmt.Errorf("API_LISTEN invalid, %s", err.Error())
			}
		}
		conf.Listen.Hosts = []string{net.JoinHostPort(hostIP, hostPort)}
	}

	if enableCors := os.Getenv("API_ENABLECORS"); enableCors != "" {
		value, err := strconv.ParseBool(enableCors)
		if err != nil {
			return fmt.Errorf("API_ENABLECORS invalid, %s", err.Error())
		}
		conf.Listen.EnableCors = value
	}
	return nil
}

func parseLogger(conf *Configuration) error {

	if logFile := os.Getenv("LOG_FILE"); logFile != "" {
		conf.Logger.LogFile = logFile
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		conf.Logger.LogLevel = logLevel
	}

	if logSize := os.Getenv("LOG_SIZE"); logSize != "" {
		value, err := strconv.ParseInt(logSize, 10, 64)
		if err != nil {
			return fmt.Errorf("LOG_SIZE invalid, %s", err.Error())
		}
		conf.Logger.LogSize = value
	}
	return nil
}
