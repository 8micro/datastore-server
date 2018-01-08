package main

import "github.com/8micro/datastore-server/api"
import "github.com/8micro/datastore-server/etc"
import "github.com/8micro/datastore-server/server"
import "github.com/8micro/gounits/flocker"
import "github.com/8micro/gounits/logger"
import "github.com/8micro/gounits/rand"

import (
	"flag"
	"time"
)

//DataStoreServer is exported
type DataStoreServer struct {
	RetryStartup bool
	Locker       *flocker.FileLocker
	APIServer    *api.Server
	DataServer   *server.DataServer
}

//NewDataStoreServer is exported
func NewDataStoreServer() (*DataStoreServer, error) {

	var confFile string
	flag.StringVar(&confFile, "f", "./etc/config.yaml", "datastore server etc file.")
	flag.Parse()

	if err := etc.New(confFile); err != nil {
		return nil, err
	}

	logger.OPEN(etc.LoggerArgs())
	key, err := rand.UUIDFile("./storedata_server.key") //服务器唯一标识文件
	if err != nil {
		return nil, err
	}

	var fLocker *flocker.FileLocker
	if pidFile := etc.PidFile(); pidFile != "" {
		fLocker = flocker.NewFileLocker(pidFile, 0)
	}

	dataServer := server.NewDataServer(key, etc.PrestConfig())
	api.RegisterStore("Key", key)
	api.RegisterStore("SystemConfig", etc.SystemConfig)
	api.RegisterStore("DataServer", dataServer)
	listen := etc.SystemConfig.Listen
	apiServer := api.NewServer(listen.Hosts, listen.EnableCors, nil)
	return &DataStoreServer{
		RetryStartup: etc.RetryStartup(),
		Locker:       fLocker,
		APIServer:    apiServer,
		DataServer:   dataServer,
	}, nil
}

//Startup is exported
func (dataStoreServer *DataStoreServer) Startup() error {

	var err error
	for {
		if err != nil {
			if dataStoreServer.RetryStartup == false {
				return err
			}
			time.Sleep(time.Second * 10) //retry, after sleep 10 seconds.
		}

		dataStoreServer.Locker.Unlock()
		if err = dataStoreServer.Locker.Lock(); err != nil {
			logger.ERROR("[#main#] pidfile lock error, %s", err)
			continue
		}
		break
	}

	go func() {
		logger.INFO("[#main#] API listener: %s", dataStoreServer.APIServer.ListenHosts())
		if err := dataStoreServer.APIServer.Startup(); err != nil {
			logger.ERROR("[#main#] API startup error, %s", err.Error())
		}
	}()
	logger.INFO("[#main#] datastore server started.")
	logger.INFO("[#main#] key:%s", dataStoreServer.DataServer.Key)
	return nil
}

//Stop is exported
func (dataStoreServer *DataStoreServer) Stop() error {

	dataStoreServer.Locker.Unlock()
	logger.INFO("[#main#] datastore server stoped.")
	logger.CLOSE()
	return nil
}
