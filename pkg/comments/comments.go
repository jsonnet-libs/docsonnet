package comments

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/formatter"
	"github.com/markbates/pkger"
)

type Blocks []string

func (b Blocks) String() string {
	return strings.Join(b, "---\n")
}

func Transform(filename, data string) (string, error) {
	return TransformStaged(filename, data, StageFormat)
}

type Stage int

const (
	StageScan Stage = iota
	StageTranslate
	StageEval
	StageJoin
	StageFormat
)

func TransformStaged(filename, data string, stage Stage) (string, error) {
	// Stage 0: Scan for comments
	blocks, err := Scan(data)
	if err != nil {
		return "", err
	}
	if stage == StageScan {
		return blocks.String(), nil
	}

	// Stage 1: Translate to DSL
	for i := range blocks {
		blocks[i] = Translate(blocks[i])
	}
	if stage == StageTranslate {
		return blocks.String(), nil
	}

	// Stage 2: Eval DSL
	for i := range blocks {
		blocks[i], err = Eval(blocks[i])
		if err != nil {
			return "", err
		}
	}
	if stage == StageEval {
		return blocks.String(), nil
	}

	// Stage 3: Join
	joined := data + Join(blocks)
	if stage == StageJoin {
		return joined, nil
	}

	// Stage 4: Format
	formatted, err := formatter.Format(filename, joined, formatter.DefaultOptions())
	if err != nil {
		return "", err
	}
	return formatted, nil
}

// Scan extracts comment blocks from the Jsonnet document
func Scan(doc string) (Blocks, error) {
	doc, err := formatter.Format("", doc, formatter.Options{
		CommentStyle: formatter.CommentStyleHash,
	})
	if err != nil {
		return nil, err
	}

	var blocks []string
	block := ""

	for _, l := range strings.Split(doc, "\n") {
		l := strings.TrimSpace(l)
		if !strings.HasPrefix(l, "#") {
			if block != "" {
				blocks = append(blocks, block)
				block = ""
			}
			continue
		}

		block += l + "\n"
	}

	return blocks, nil
}

// Translate converts the comment syntax into Jsonnet DSL invocations
func Translate(block string) string {
	block = strings.Replace(block, "# @", "+ ", -1)
	block = strings.Replace(block, "# ", "", -1)
	return block
}

// Eval converts the block into an actual object, by evaluating it in the
// context of our DSL
func Eval(block string) (string, error) {
	dsl := loadDSL()

	vm := jsonnet.MakeVM()
	out, err := vm.EvaluateSnippet("", dsl+"\n"+block)
	if err != nil {
		return "", err
	}

	return out, nil
}

// Join chains multiple comment blocks into a single patch
func Join(blocks Blocks) string {
	s := ""
	for _, b := range blocks {
		s += "+ " + b
	}
	return s
}

func loadDSL() string {
	p, err := pkger.Open("/pkg/comments/dsl.libsonnet")
	if err != nil {
		// This must work. panic if not
		panic(fmt.Errorf("Loading embedded file: %s. This build appears broken", err))
	}
	data, err := ioutil.ReadAll(p)
	if err != nil {
		panic(fmt.Errorf("Loading embedded file: %s. This build appears broken", err))
	}
	return string(data)
}
