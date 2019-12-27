package exporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/smileinnovation/imannotate/api/annotation"
)

type CSV struct {
}

func (c *CSV) Export(annotations []*annotation.Annotation) io.Reader {

	lines := []string{}
	for _, ann := range annotations {
		im := ann.Image
		for _, b := range ann.Boxes {
			lines = append(lines, fmt.Sprintf("%s,%s,%s,%s", im, b.Label, b.Username, b.Timestamp))
		}
	}

	return strings.NewReader(strings.Join(lines, "\n"))
}
