package main

import (
	"github.com/jideji/servicelauncher/procs"
	"github.com/magiconair/properties"
	"log"
	"os"
)

type Config struct {
	Pid     int    `properties:"pid"`
	Command string `properties:"command"`
}

func main() {
	p := properties.MustLoadFile(os.ExpandEnv("$HOME/.slcfg"), properties.UTF8)
	var cfg Config
	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}

	pr, err := procs.FindByCommandLine(cfg.Command)
	if err != nil {
		panic(err)
	}
	println(pr.Pid)
	println(pr.CommandLine)
}
