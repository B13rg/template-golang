// Defines the cli-interface commands available to the user.
//
//nolint:gochecknoinits,gochecknoglobals

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// exported version variable.
var version string

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "",
	Short: "Short desc",
	Long:  `Long description`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ver string) {
	version = ver
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// Default options that are available to all commands.
type CmdRootOptions struct {
	// log more information about what the tool is doing. Overrides --loglevel
	Debug bool
	// set log level
	LogLevel string
	// enable colorized output (default true). Set to false to disable")
	Color bool
	// Profiling output directory.  Only captured if set.
	ProfilingDir string
	// CPU profiling output file handle.
	ProfilingCPUFile *os.File
}

var RootConfig CmdRootOptions

func init() {
	// Ran before each command is ran
	cobra.OnInitialize(InitConfig, ProfilingInitializer)
	cobra.OnFinalize(ProfilingFinalizer)

	RootCmd.PersistentFlags().BoolVar(&RootConfig.Debug,
		"debug", false,
		"log additional information about what the tool is doing. Overrides --loglevel")
	RootCmd.PersistentFlags().StringVarP(&RootConfig.LogLevel,
		"loglevel", "L", "info",
		"set zerolog log level")
	RootCmd.PersistentFlags().BoolVar(&RootConfig.Color,
		"color", true,
		"enable colorized output")
	RootCmd.PersistentFlags().StringVarP(&RootConfig.ProfilingDir,
		"profiledir", "", "",
		"directory to write pprof profile data to")
}

// InitConfig reads in config file and ENV variables if set.
func InitConfig() {
	ConfigureLogger(RootConfig.Debug)

	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME") // adding home directory as first search path
	viper.SetEnvPrefix("GCFG")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debug().Msg("Using config file:" + viper.ConfigFileUsed())
	}
}

// Stop profiling and write cpu and memory profiling files if configured.
func ProfilingFinalizer() {
	if RootConfig.ProfilingDir != "" {
		pprof.StopCPUProfile()
		if RootConfig.ProfilingCPUFile != nil {
			_ = RootConfig.ProfilingCPUFile.Close()
		}

		runtime.GC() // get up-to-date statistics

		// Various types of profiles that can be collected:
		// https://cs.opensource.google/go/go/+/go1.24.2:src/runtime/pprof/pprof.go;l=178
		var err error
		heapFile, err := os.Create(filepath.Join(RootConfig.ProfilingDir, "profile_heap.pb.gz"))
		if err != nil {
			log.Fatal().Err(err).Msg("could not write memory profile: ")
		}
		if err = pprof.WriteHeapProfile(heapFile); err != nil {
			_ = heapFile.Close()
			log.Fatal().Err(err).Msg("could not write memory profile: ")
		}
		_ = heapFile.Close()
	}
}

// Sets up program profiling.
func ProfilingInitializer() {
	var err error
	if RootConfig.ProfilingDir != "" {
		RootConfig.ProfilingCPUFile, err = os.Create(filepath.Join(RootConfig.ProfilingDir, "profile_cpu.pb.gz"))
		if err != nil {
			log.Fatal().Err(err).Msg("could not create CPU profile: ")
		}
		if err := pprof.StartCPUProfile(RootConfig.ProfilingCPUFile); err != nil {
			log.Fatal().Err(err).Msg("could not create CPU profile: ")
		}
	}
}
