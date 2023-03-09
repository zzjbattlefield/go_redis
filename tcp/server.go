package tcp

import (
	"context"
	"go_redis/interface/tcp"
	"go_redis/lib/logger"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Config struct {
	Address string
}

func ListenAndServerWithSignal(cfg *Config, handler tcp.Handler) error {
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	logger.Info("start listen")
	closeChan := make(chan struct{})
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigChan
		closeChan <- struct{}{}
	}()
	ListenAndServer(listener, handler, closeChan)
	return nil
}

func ListenAndServer(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}) {
	ctx := context.Background()
	go func() {
		<-closeChan
		logger.Info("shutting down")
		listener.Close()
		handler.Close()
	}()
	defer func() {
		listener.Close()
		handler.Close()
	}()
	var waitDone sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		logger.Info("accepted link")
		waitDone.Add(1)
		go func() {
			defer func() {
				waitDone.Done()
			}()
			handler.Handle(ctx, conn)
		}()
	}
	waitDone.Wait()
}
