package cli

import (
	"github.com/baasapi/baasapi/api"

	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
)

type pairList []baasapi.Pair

// Set implementation for a list of baasapi.Pair
func (l *pairList) Set(value string) error {
	parts := strings.SplitN(value, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("expected NAME=VALUE got '%s'", value)
	}
	p := new(baasapi.Pair)
	p.Name = parts[0]
	p.Value = parts[1]
	*l = append(*l, *p)
	return nil
}

// String implementation for a list of pair
func (l *pairList) String() string {
	return ""
}

// IsCumulative implementation for a list of pair
func (l *pairList) IsCumulative() bool {
	return true
}

func pairs(s kingpin.Settings) (target *[]baasapi.Pair) {
	target = new([]baasapi.Pair)
	s.SetValue((*pairList)(target))
	return
}
