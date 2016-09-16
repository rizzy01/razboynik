package app

import (
	"fmt"
	"fuzzer"
	"fuzzer/src/bash"
	"io"
	"log"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/urfave/cli"
)

type MainInterface struct {
	cmd         *cli.App
	completer   *readline.PrefixCompleter
	commands    []cli.Command
	isConnected bool
	running     bool
}

func CreateMainApp() *MainInterface {
	c := MainInterface{}
	c._buildCommand()
	c._buildCompleter()
	c.cmd = _createCli(&c.commands)

	return &c
}

func (main *MainInterface) _buildCompleter() {
	main.completer = readline.NewPrefixCompleter()
	lgt := len(main.commands)

	for i := 0; i < lgt; i++ {
		child := readline.PcItem(main.commands[i].Name)
		main._buildChild(&main.commands[i].Subcommands, child)

		main.completer.SetChildren(append(main.completer.GetChildren(), child))
	}
}

func (main *MainInterface) _buildChild(sub *cli.Commands, parent *readline.PrefixCompleter) {
	childLgt := len(*sub)

	for x := 0; x < childLgt; x++ {
		child := readline.PcItem((*sub)[x].Name)
		parent.SetChildren(append(parent.GetChildren(), child))

		subChild := (*sub)[x].Subcommands
		if len(subChild) > 0 {
			main._buildChild(&subChild, child)
		}
	}
}

func _addInformation(app *cli.App) {
	app.Name = "Fuzzer"
	app.Version = "2.0.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Kamikaze",
			Email: "leakleass@protomail.com",
		},
	}
	app.Copyright = "(c) 2016 Eat Bytes"
	app.Usage = "File upload meterpreter"
}

func _createCli(commands *[]cli.Command) *cli.App {
	client := cli.NewApp()
	client.Commands = *commands
	client.CommandNotFound = func(c *cli.Context, command string) {
		cli.ShowAppHelp(c)
	}

	_addInformation(client)

	return client
}

func (main *MainInterface) GetCommand(line string) []string {
	line = strings.TrimSpace(line)
	appName := []string{"Fuzzer"}
	args := strings.Fields(line)
	command := append(appName, args...)

	return command
}

func (main *MainInterface) GetPrompt() *readline.Instance {
	l, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[31m»\033[0m ",
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    main.completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})

	if err != nil {
		panic(err)
	}

	return l
}

func (main *MainInterface) Loop() {
	prompt := main.GetPrompt()
	defer prompt.Close()
	log.SetOutput(prompt.Stderr())

	for main.IsRunning() {
		line, err := prompt.Readline()
		if err == readline.ErrInterrupt || err == io.EOF {
			return
		}

		if len(line) == 0 {
			continue
		}

		command := main.GetCommand(line)
		main.Run(command)
	}
}

func (main *MainInterface) Run(command []string) {
	main.cmd.Run(command)
}

func (main *MainInterface) Help(c *cli.Context) {
	cli.ShowAppHelp(c)
}

func (main *MainInterface) Generate(c *cli.Context) {
	fmt.Println("generate")
}

func (main *MainInterface) Exit(c *cli.Context) {
	main.running = false
}

func (main *MainInterface) Start() {
	main.running = true
	main.Loop()
}

func (main *MainInterface) Stop() {
	main.running = false
}

func (main *MainInterface) IsRunning() bool {
	return main.running
}

func (main *MainInterface) Encode(c *cli.Context) {
	str := c.Args().Get(0)
	sEnc := fuzzer.Encode(str)
	fmt.Println(sEnc)
}

func (main *MainInterface) Decode(c *cli.Context) {
	str := c.Args().Get(0)
	sDec := fuzzer.Decode(str)
	fmt.Println(sDec)
}

func (main *MainInterface) StartBash(c *cli.Context) {
	bsh := bash.CreateBashApp()
	bsh.Start()
}
