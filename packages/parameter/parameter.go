package parameter

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type state struct {
	Args  []string
	Index int
	Must  bool
}

func New(args []string) *state {
	return &state{Args: args}
}

func (p *state) Set(index int, must bool) *state {
	if p == nil {
		log.Fatalf("args params undefined")
	}
	p.Index = index
	p.Must = must
	return p
}

func (p *state) NumberInt64() (value int64, err error) {
	if len(p.Args) > p.Index {
		value, err = strconv.ParseInt(p.Args[p.Index], 10, 64)
		if err != nil {
			err = fmt.Errorf("args[%d] invalid:%s err:%s", p.Index, p.Args[p.Index], err.Error())
			return
		}
		return
	}
	if p.Must {
		return 0, fmt.Errorf("args[%d] Can't Not Be Empty", p.Index)
	}

	return
}

func (p *state) NumberInt() (value int, err error) {
	if len(p.Args) > p.Index {
		value, err = strconv.Atoi(p.Args[p.Index])
		if err != nil {
			err = fmt.Errorf("args[%d] invalid:%s err:%s", p.Index, p.Args[p.Index], err.Error())
			return
		}
		return
	}
	if p.Must {
		return 0, fmt.Errorf("args[%d] Can't Not Be Empty", p.Index)
	}

	return
}

func (p *state) NumberUint64() (value uint64, err error) {
	if len(p.Args) > p.Index {
		value, err = strconv.ParseUint(p.Args[p.Index], 10, 64)
		if err != nil {
			err = fmt.Errorf("args[%d] invalid:%s err:%s", p.Index, p.Args[p.Index], err.Error())
			return
		}
		return
	}
	if p.Must {
		return 0, fmt.Errorf("args[%d] Can't Not Be Empty", p.Index)
	}

	return
}

func (p *state) String() (value string, err error) {
	if len(p.Args) > p.Index {
		value = p.Args[p.Index]
		return
	}
	if p.Must {
		return "", fmt.Errorf("args[%d] Can't Not Be Empty", p.Index)
	}

	return
}
