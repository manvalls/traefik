// +build !windows

package server

import (
	"os/signal"
	"syscall"
	"time"

	"github.com/manvalls/traefik/log"
)

func (server *Server) configureSignals() {
	signal.Notify(server.signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
}

func (server *Server) listenSignals() {
	for {
		sig := <-server.signals
		switch sig {
		case syscall.SIGUSR1:
			log.Infof("Closing and re-opening log files for rotation: %+v", sig)

			if server.accessLoggerMiddleware != nil {
				if err := server.accessLoggerMiddleware.Rotate(); err != nil {
					log.Errorf("Error rotating access log: %s", err)
				}
			}

			if err := log.RotateFile(); err != nil {
				log.Errorf("Error rotating traefik log: %s", err)
			}
		default:
			log.Infof("I have to go... %+v", sig)
			reqAcceptGraceTimeOut := time.Duration(server.globalConfiguration.LifeCycle.RequestAcceptGraceTimeout)
			if reqAcceptGraceTimeOut > 0 {
				log.Infof("Waiting %s for incoming requests to cease", reqAcceptGraceTimeOut)
				time.Sleep(reqAcceptGraceTimeOut)
			}
			log.Info("Stopping server gracefully")
			server.Stop()
		}
	}
}
