package slug

import (
	"regexp"
	"strconv"
	"strings"
)

type Slugger struct {
	occurences map[string]int
}

var (
	expWhitespace = regexp.MustCompile(`\s`)
	expSpecials   = regexp.MustCompile("[\u2000-\u206F\u2E00-\u2E7F\\'!\"#$%&()*+,./:;<=>?@[\\]^`{|}~’]")
)

func New() *Slugger {
	return &Slugger{
		occurences: make(map[string]int),
	}
}

func (s *Slugger) Slug(str string) string {
	str = expWhitespace.ReplaceAllString(str, "-")
	str = expSpecials.ReplaceAllString(str, "")

	old := str
	if o := s.occurences[str]; o > 0 {
		str += "-" + strconv.Itoa(o)
	}
	s.occurences[old] = s.occurences[old] + 1

	return strings.ToLower(str)
}
