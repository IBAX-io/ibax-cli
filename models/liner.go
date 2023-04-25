package models

import (
	"fmt"
	"github.com/IBAX-io/ibax-cli/conf"
	"github.com/peterh/liner"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	historyFn       string
	wordCompletions []string

	exit = regexp.MustCompile(`^\s*exit\s*;*\s*$`)
)

type linerConsole struct {
	*liner.State
	Nonce int
}

func AddWordsCompletions(worlds []string) {
	wordCompletions = append(wordCompletions, worlds...)
}

func NewConsole() *linerConsole {
	initLinerConfig()
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
		for _, n := range wordCompletions {
			if strings.HasPrefix(n, line) {
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
	defer func() {
		if f, err := os.Create(historyFn); err != nil {
			log.Print("Error writing history file: ", err)
			return
		} else {
			line.WriteHistory(f)
			f.Close()
		}
	}()

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

			args := strings.Split(command, " ")
			os.Args = append(os.Args[:1], args...)

			subCmd, _, err := RootCmd.Find(args)
			if err != nil {
				fmt.Println(err.Error())
				continue
			} else {
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
