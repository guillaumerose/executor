package main

import (
	"net"
	"os"
	"strings"

	"google.golang.org/grpc"
	"yunion.io/x/log"

	"yunion.io/x/executor/apis"
	"yunion.io/x/executor/server"
)

type SExecuteService struct {
}

func NewExecuteService() *SExecuteService {
	return &SExecuteService{}
}

func (s *SExecuteService) fixPathEnv() error {
	var paths = []string{
		"/usr/local/sbin",
		"/usr/local/bin",
		"/sbin",
		"/bin",
		"/usr/sbin",
		"/usr/bin",
	}
	return os.Setenv("PATH", strings.Join(paths, ":"))
}

func (s *SExecuteService) prepareEnv() error {
	if err := s.fixPathEnv(); err != nil {
		return err
	}
	return nil
}

func (s *SExecuteService) runService() {
	grpcServer := grpc.NewServer()
	apis.RegisterExecutorServer(grpcServer, &server.Executor{})
	if _, err := os.Stat(socketPath); !os.IsNotExist(err) {
		// socket file already exist, remove first
		if err := os.Remove(socketPath); err != nil {
			log.Fatalln(err)
		}
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()
	log.Infof("Init net listener on %s succ", socketPath)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalln(err)
	}
}

func (s *SExecuteService) initService() {
	if len(socketPath) == 0 {
		log.Fatalf("missing socket path")
	}
	if err := s.prepareEnv(); err != nil {
		log.Fatalln(err)
	}
}

func (s *SExecuteService) Run() {
	s.initService()
	s.runService()
}

func Server() {
	NewExecuteService().Run()
}
