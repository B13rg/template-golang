# cmd

```go
import "github.com/b13rg/template-golang/cmd"
```

## Index

- [Variables](<#variables>)
- [func Colorize\(input interface\{\}, colorNum int, disabled bool\) string](<#Colorize>)
- [func ConfigureLogger\(debug bool\)](<#ConfigureLogger>)
- [func Execute\(ver string\)](<#Execute>)
- [func InitConfig\(\)](<#InitConfig>)
- [func ProfilingFinalizer\(\)](<#ProfilingFinalizer>)
- [func ProfilingInitializer\(\)](<#ProfilingInitializer>)
- [func SetupLogger\(enableColor bool\) zerolog.Logger](<#SetupLogger>)
- [type CmdRootOptions](<#CmdRootOptions>)
- [type Stamp](<#Stamp>)


## Variables

<a name="RootCmd"></a>RootCmd represents the base command when called without any subcommands.

```go
var RootCmd = &cobra.Command{
    Use:   "",
    Short: "Short desc",
    Long:  `Long description`,
}
```

<a name="VersionCmd"></a>Print out versions of packages in use. Chore\(\) \- Updated manually.

```go
var VersionCmd = &cobra.Command{
    Use:   "version",
    Short: "Get version",
    Long:  `Get the current version of tool`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println(RootCmd.Use + "+ Version: " + version)
        info, ok := debug.ReadBuildInfo()
        if !ok {
            log.Fatal().Msg("could not read build info")
        }
        stamp := retrieveStamp(info)
        fmt.Printf("  Built with %s on %s\n", stamp.InfoGoCompiler, stamp.InfoBuildTime)
        fmt.Printf("  VCS revision: %s\n", stamp.VCSRevision)
        fmt.Printf("  Go version %s, GOOS %s, GOARCH %s\n", info.GoVersion, stamp.InfoGOOS, stamp.InfoGOARCH)
        fmt.Print("  Dependencies:\n")
        for _, mod := range retrieveDepends(info) {
            fmt.Printf("    %s\n", mod)
        }

    },
}
```

<a name="Colorize"></a>
## func [Colorize](<https://github.com:b13rg/template-golang/blob/main/cmd/logging.go#L67>)

```go
func Colorize(input interface{}, colorNum int, disabled bool) string
```

Colorize function from zerolog console.go file to replicate their coloring functionality. Source: https://github.com/rs/zerolog/blob/a21d6107dcda23e36bc5cfd00ce8fdbe8f3ddc23/console.go#L389 Replicated here because it's a private function.

<a name="ConfigureLogger"></a>
## func [ConfigureLogger](<https://github.com:b13rg/template-golang/blob/main/cmd/logging.go#L12>)

```go
func ConfigureLogger(debug bool)
```



<a name="Execute"></a>
## func [Execute](<https://github.com:b13rg/template-golang/blob/main/cmd/root.go#L31>)

```go
func Execute(ver string)
```

Execute adds all child commands to the root command sets flags appropriately. This is called by main.main\(\). It only needs to happen once to the rootCmd.

<a name="InitConfig"></a>
## func [InitConfig](<https://github.com:b13rg/template-golang/blob/main/cmd/root.go#L75>)

```go
func InitConfig()
```

InitConfig reads in config file and ENV variables if set.

<a name="ProfilingFinalizer"></a>
## func [ProfilingFinalizer](<https://github.com:b13rg/template-golang/blob/main/cmd/root.go#L91>)

```go
func ProfilingFinalizer()
```

Stop profiling and write cpu and memory profiling files if configured.

<a name="ProfilingInitializer"></a>
## func [ProfilingInitializer](<https://github.com:b13rg/template-golang/blob/main/cmd/root.go#L116>)

```go
func ProfilingInitializer()
```

Sets up program profiling.

<a name="SetupLogger"></a>
## func [SetupLogger](<https://github.com:b13rg/template-golang/blob/main/cmd/logging.go#L38>)

```go
func SetupLogger(enableColor bool) zerolog.Logger
```

Configure zerolog with some defaults and cleanup error formatting.

<a name="CmdRootOptions"></a>
## type [CmdRootOptions](<https://github.com:b13rg/template-golang/blob/main/cmd/root.go#L40-L51>)

Default options that are available to all commands.

```go
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
```

<a name="RootConfig"></a>

```go
var RootConfig CmdRootOptions
```

<a name="Stamp"></a>
## type [Stamp](<https://github.com:b13rg/template-golang/blob/main/cmd/version.go#L14-L21>)



```go
type Stamp struct {
    InfoGoVersion  string
    InfoGoCompiler string
    InfoGOARCH     string
    InfoGOOS       string
    InfoBuildTime  string
    VCSRevision    string
}
```