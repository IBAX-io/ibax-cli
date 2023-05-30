package models

import (
	"fmt"
	"github.com/IBAX-io/ibax-cli/conf"
	"github.com/peterh/liner"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	historyFn       string
	wordCompletions []string

	exit = regexp.MustCompile(`^\s*exit\s*;*\s*$`)

	consoleMode bool

	ErrSignal = make(chan ErrInfo, 3)
)

type ErrInfo struct {
	err    error
	isExit bool
}

type linerConsole struct {
	*liner.State
	Nonce int
}

func SendErrSignal(err error, exit bool) {
	var info ErrInfo
	info.err = err
	info.isExit = exit
	select {
	case ErrSignal <- info:
	default:
	}
}

func IsConsoleMode() bool {
	return consoleMode
}

func AddWordsCompletions(worlds []string) {
	wordCompletions = append(wordCompletions, worlds...)
}

func NewConsole() *linerConsole {
	initLinerConfig()
	consoleMode = true
	return &linerConsole{
		liner.NewLiner(),
		1,
	}
}

func (p *linerConsole) Close() {
	p.State.Close()
}

func initLinerConfig() {
	err := MakeDirectory(conf.Config.LinerPath)
	if err != nil {
		log.Fatalf("make liner directory failed:%s", err.Error())
		return
	}
	historyFn = filepath.Join(conf.Config.LinerPath, ".liner_history")

}

// MakeDirectory makes directory if is not exists
func MakeDirectory(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return os.Mkdir(dir, 0775)
		}
		return err
	}
	return nil
}

func NewTerminalLiner(line *linerConsole) {
	line.SetCtrlCAborts(true)
	line.SetMultiLineMode(true)
	line.SetTabCompletionStyle(liner.TabPrints)

	line.SetWordCompleter(func(line string, pos int) (head string, completions []string, tail string) {
		keyWorld := strings.ToLower(line)
		for _, n := range wordCompletions {
			word := strings.ToLower(n)
			if strings.HasPrefix(word, keyWorld) {
				completions = append(completions, n)
			}
		}
		return
	})

	//line.SetCompleter(func(line string) (c []string) {
	//	for _, n := range wordCompletions {
	//		if strings.HasPrefix(n, line) {
	//			c = append(c, n)
	//		}
	//	}
	//	return
	//})

	if f, err := os.Open(historyFn); err == nil {
		line.ReadHistory(f)
		f.Close()
	}
	defer lineConsoleClose(line)

	go StartDaemon(line)
	for {
		if command, err := line.Prompt(">"); err == nil {
			if command == "" {
				continue
			}
			if exit.MatchString(command) {
				log.Print("Exit")
				return
			}
			line.AppendHistory(command)

			list := strings.Split(command, " ")
			var completionArgs string
			var args []string
			var findIndents bool
			for _, str := range list {
				indents := countIndents(completionArgs + str)
				//fmt.Printf("str:%s,indents:%d\n", str, indents)
				if indents <= 0 {
					if findIndents {
						completionArgs += " " + str
						findIndents = false
					} else {
						completionArgs += str
					}
					args = append(args, completionArgs)
					completionArgs = ""
				} else {
					if !findIndents {
						findIndents = true
					}
					if completionArgs != "" {
						completionArgs += " " + str
					} else {
						completionArgs += str
					}
				}
			}
			os.Args = append(os.Args[:1], args...)
			//for k, v := range os.Args {
			//	fmt.Printf("k[%d]:%s\n", k, v)
			//}

			subCmd, _, err := globalCmd.Find(args)
			if err != nil {
				fmt.Println(err.Error())
				continue
			} else {
				resetAllFlags(subCmd)
				err = subCmd.Execute()
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
			}
		} else if err == liner.ErrPromptAborted {
			continue
		} else {
			if err.Error() == "EOF" {
				fmt.Println("EOF")
				return
			}
			log.Print("Error reading line: ", err)
			return
		}
	}

}

// countIndents returns the number of indentations for the given input.
// In case of invalid input such as var a = } the result can be negative.
func countIndents(input string) int {
	var (
		indents     = 0
		inString    = false
		strOpenChar = ' '   // keep track of the string open char to allow var str = "I'm ....";
		charEscaped = false // keep track if the previous char was the '\' char, allow var str = "abc\"def";
	)

	for _, c := range input {
		switch c {
		case '\\':
			// indicate next char as escaped when in string and previous char isn't escaping this backslash
			if !charEscaped && inString {
				charEscaped = true
			}
		case '\'', '"':
			if inString && !charEscaped && strOpenChar == c { // end string
				inString = false
				indents--
			} else if !inString && !charEscaped { // begin string
				indents++
				inString = true
				strOpenChar = c
			}
			charEscaped = false
		case '{', '(':
			if !inString { // ignore brackets when in string, allow var str = "a{"; without indenting
				indents++
			}
			charEscaped = false
		case '}', ')':
			if !inString {
				indents--
			}
			charEscaped = false
		default:
			charEscaped = false
		}
	}

	return indents
}

func resetAllFlags(cmd *cobra.Command) {
	if cmd != nil {
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if f.Name != "path" && f.Name != "rcpConnect" && f.Name != "rpcPort" {
				switch f.Value.Type() {
				case "bool":
					f.Value.Set(f.DefValue)
				case "int", "int64":
					f.Value.Set(f.DefValue)
				case "string":
					f.Value.Set(f.DefValue)
				}
			}
		})
	}
}

func lineConsoleClose(line *linerConsole) {
	if f, err := os.Create(historyFn); err != nil {
		log.Print("Error writing history file: ", err)
		return
	} else {
		line.WriteHistory(f)
		f.Close()
	}
}

func StartDaemon(line *linerConsole) {
	for {
		select {
		case info := <-ErrSignal:
			if info.err != nil {
				fmt.Println(info.err)
				if info.isExit {
					lineConsoleClose(line)
					os.Exit(1)
				}
			}
		}
	}
}
