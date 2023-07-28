package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mzz2017/juicity/config"
	"github.com/mzz2017/juicity/pkg/log"
	"github.com/mzz2017/juicity/server"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func init() {
	runCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Config file of juicity-server.")
	runCmd.PersistentFlags().BoolVarP(&disableTimestamp, "disable-timestamp", "", false, "Disable timestamp.")
}

var (
	cfgFile          string
	disableTimestamp bool

	runCmd = &cobra.Command{
		Use:   "run",
		Short: "To run juicity-server in the foreground.",
		Run: func(cmd *cobra.Command, args []string) {
			if cfgFile == "" {
				log.Logger().
					Fatal().
					Msg("Argument \"--config\" or \"-c\" is required but not provided.")
			}

			// Read config from --config cfgFile.
			conf, err := config.ReadConfig(cfgFile)
			if err != nil {
				log.Logger().
					Fatal().
					Err(err).
					Msg("Failed to read config")
			}
			lvl, err := zerolog.ParseLevel(conf.LogLevel)
			if err != nil {
				log.Logger().
					Fatal().
					Err(err).
					Send()
			}
			log.Logger().Level(lvl)

			go func() {
				if err := Serve(conf); err != nil {
					log.Logger().Fatal().
						Err(err).
						Send()
				}
			}()
			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGILL)
			for sig := range sigs {
				log.Logger().Warn().
					Str("signal", sig.String()).
					Msg("Exiting")
				return
			}
		},
	}
)

func Serve(conf *config.Config) error {
	s, err := server.New(&server.Options{
		Users:             conf.Users,
		Certificate:       conf.Certificate,
		PrivateKey:        conf.PrivateKey,
		CongestionControl: conf.CongestionControl,
	})
	if err != nil {
		return err
	}
	if conf.Listen == "" {
		return fmt.Errorf(`"Listen" is required`)
	}
	if err = s.Serve(conf.Listen); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(runCmd)
}
