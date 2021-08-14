package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"

	"github.com/go-co-op/gocron"
	"github.com/kelseyhightower/envconfig"
	"github.com/oleiade/lane"
	"gopkg.in/mcuadros/go-syslog.v2"
	"gopkg.in/yaml.v2"
)

var lastLog *lane.Deque

type Config struct {
	Server struct {
		Host      string `yaml:"host",envconfig:"SERVER_HOST"`
		UdpPort   string `yaml:"udpPort",envconfig:"SERVER_PORT"`
		TcpPort   string `yaml:"tcpPort",envconfig:"SERVER_PORT_TCP"`
		LogFormat string `yaml:"logFormat",envconfig:"SERVER_LOGFORMAT"`
	} `yaml:"server"`
	Log struct {
		Path     string `yaml:"path",envconfig:"LOG_PATH"`
		Filename string `yaml:"filename",envconfig:"LOG_FILENAME"`
		Rotate   string `yaml:"rotate",envconfig:"LOG_ROTATE"`
		Keep     string `yaml:"keep",envconfig:"LOG_KEEP"`
	} `yaml:"log"`
	Database struct {
		Username string `yaml:"user",envconfig:"DB_USERNAME"`
		Password string `yaml:"pass",envconfig:"DB_PASSWORD"`
	} `yaml:"database"`
	App struct {
		Debug string `yam:"debug",envconfig:"APP_DEBUG"`
	} `yaml:"app"`
}

func getIPPort(addrPort string) string {
	results := strings.Split(addrPort, ":")
	return results[1]
}

func getIPAddr(addrPort string) string {
	results := strings.Split(addrPort, ":")
	return results[0]
}

var cfg Config
var totalBytesWritten int

const time_in_seconds = 10

func main() {

	// Config
	readFile(&cfg)
	readEnv(&cfg)
	//fmt.Printf("%+v", cfg)

	//Init lastLog Queue
	lastLog = lane.NewCappedDeque(10)

	//Log files init
	touchLogFile()

	//Scheduler init
	schdlr := gocron.NewScheduler(time.UTC)
	schdlr.StartAsync()

	schdlr.Every(60).Seconds().Do(func() {
		printStats()
	})

	mlc := make(chan string, 10) // main loop channel

	//Print some info
	printIntro()

	// Start the syslogd server
	go syslogServer()

	// Start the keypress listener
	go keyPressListener(mlc)

	// Main loop
	for msg := range mlc {
		//fmt.Println("This is the main loop")
		//fmt.Println("main: " + msg)
		switch msg {
		case "shutdown":
			os.Exit(0)
		case "showstats":
			printStats()

		case "showinfo":
			fmt.Println("TODO: show info")
		case "dump":
			printLastLog()
			// case "peek":
			// 	fmt.Println(lastLog.Full())
			//
		}
	}

}

func keyPressListener(mlc chan string) {
	//kplChan := make(chan string, 2) // keypress listener channel

	// Keyboard setup
	keysEvents, err := keyboard.GetKeys(2)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Println("Press [CTRL]+[X] to quit, [i] for info, [s] for stats, [d] dump last 10 log lines")

	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}
		//fmt.Printf("You pressed: rune %q, key %X\r\n", event.Rune, event.Key)
		if event.Key == keyboard.KeyCtrlX {
			mlc <- "shutdown"
			break
		}

		if event.Rune == 'i' {
			mlc <- "showinfo"
		}
		if event.Rune == 's' {
			mlc <- "showstats"
		}
		if event.Rune == 'd' {
			mlc <- "dump"
		}
		if event.Rune == 'p' {
			mlc <- "peek"
		}
	}

	// for kplChan := range kplChan {
	// 	fmt.Println("This is the keypresslistener loop")
	// 	fmt.Println("kpl: " + kplChan)

	// 	event := <-keysEvents
	// 	if event.Err != nil {
	// 		panic(event.Err)
	// 	}
	// 	fmt.Printf("You pressed: rune %q, key %X\r\n", event.Rune, event.Key)
	// 	if event.Key == keyboard.KeyEsc {
	// 		break
	// 	}
	// 	time.Sleep(1)
	// }
}

func printIntro() {
	// Print with default helper functions
	d := color.New(color.FgCyan, color.Bold)
	d.Printf("[octalspan] ")
	fmt.Printf("syslog server starting...\n")
	d.Printf("[octalspan] ")
	fmt.Printf("Listening on ")
	m := color.New(color.FgMagenta)
	y := color.New(color.FgYellow)
	m.Printf("%s:", cfg.Server.Host)
	y.Printf("%s", cfg.Server.UdpPort)
	fmt.Println("")
}

func printLastLog() {
	var logLines []string = make([]string, lastLog.Size())

	for i := 0; i < len(logLines); i++ {
		value := lastLog.Shift()
		logLines[i] = value.(string)
	}

	fmt.Println(strings.Join(logLines, "\n"))
}

func syslogServer() {
	// Begin syslog lib setup
	channel := make(syslog.LogPartsChannel)

	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.RFC3164)
	server.SetHandler(handler)

	server.ListenUDP(fmt.Sprintf("%s:%v", cfg.Server.Host, cfg.Server.UdpPort))
	server.SetTimeout(10)
	server.ListenTCP(fmt.Sprintf("%s:%v", cfg.Server.Host, cfg.Server.TcpPort))

	server.Boot()

	go func(channel syslog.LogPartsChannel) {

		f, err := os.OpenFile(cfg.Log.Path+cfg.Log.Filename, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
		l := 0

		for logParts := range channel {

			hostIP := getIPAddr(logParts["client"].(string))

			logLine := fmt.Sprintf("%v - %v - %v", logParts["timestamp"], hostIP, logParts["content"])
			// Write log line to file
			l, err = fmt.Fprintln(f, logLine)

			// Add log line to lastLog queue
			lastLog.Append(logLine)

			if err != nil {
				fmt.Println(err)
				f.Close()
				return
			}
			totalBytesWritten += l
			if cfg.App.Debug == "TRUE" {
				fmt.Println(l, "bytes written successfully")
			}
		}

		err = f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}

	}(channel)

	go func(channel syslog.LogPartsChannel) {
		if cfg.App.Debug == "TRUE" {
			for logParts := range channel {
				fmt.Printf("%T\n\n", logParts)
				fmt.Printf("  %+9s = %v\n", "client", logParts["client"])
				fmt.Printf("  %+9s = %v\n", "hostname", logParts["hostname"])
				fmt.Printf("  %+9s = %v\n", "tls_peer", logParts["tls_peer"])
				fmt.Printf("  %+9s = %v\n", "facility", logParts["facility"])
				fmt.Printf("  %+9s = %v\n", "priority", logParts["prority"])
				fmt.Printf("  %+9s = %v\n", "severity", logParts["severity"])
				fmt.Printf("  %+9s = %v\n", "tag", logParts["tag"])
				fmt.Printf("  %+9s = %v\n", "content", logParts["content"])
				fmt.Printf("  %+9s = %v\n", "timestamp", logParts["timestamp"])
				fmt.Println()
			}
		}
	}(channel)

	server.Wait()
}

func printStats() {
	d := color.New(color.FgCyan, color.Bold)
	d.Printf("[octalspan] ")
	fmt.Printf("%d bytes written to disk\n", totalBytesWritten)
}

// func for config and env loading
func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func readFile(cfg *Config) {
	f, err := os.Open("octalspan.yml")
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		processError(err)
	}
}

func touchLogFile() {

	exeDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		processError(err)
	}

	fmt.Println(exeDir)

	logDir := exeDir + "/" + cfg.Log.Path
	fullFilePath := cfg.Log.Path + cfg.Log.Filename

	if _, err = os.Stat(logDir); os.IsNotExist(err) {
		err := os.Mkdir(logDir, 0644)
		processError(err)
	}

	_, err = os.Stat(fullFilePath)
	if os.IsNotExist(err) {
		file, err := os.Create(fullFilePath)
		if err != nil {
			processError(err)
		}
		defer file.Close()
	} else {
		currentTime := time.Now().Local()
		err = os.Chtimes(fullFilePath, currentTime, currentTime)
		if err != nil {
			processError(err)
		}
	}
}
