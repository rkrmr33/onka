package compiler

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ghodss/yaml"
	"github.com/rkrmr33/onka/pkg/proto/v1alpha1"
	"github.com/rkrmr33/onka/pkg/util"
	log "github.com/sirupsen/logrus"
)

// Compilation Errors
var (
	ErrNoEntry = errors.New("pipeline has no entry point")
)

//go:embed builtins
var builtinsFS embed.FS

var builtins = map[string]string{}

// Compile receives yaml of a pipeline, compiles it and returns a tree
// consists tasks in the order they should be executed in.
func Compile(orgPipelineYAML []byte) (p *v1alpha1.Pipeline, entrypoints []v1alpha1.PipelineTask, err error) {
	p = &v1alpha1.Pipeline{}

	orgPipelineYAML, err = render(orgPipelineYAML)
	if err != nil {
		return
	}
	log.Debugf("after pre-processor: \n%s\n", string(orgPipelineYAML))

	pipelineJSON, err := yaml.YAMLToJSON(orgPipelineYAML)
	if err != nil {
		return
	}
	log.Debugf("json pipe: \n%s\n", string(pipelineJSON))

	if err = json.Unmarshal(pipelineJSON, p); err != nil {
		return
	}

	tasks, err := p.GetOrederdTasksYAML(orgPipelineYAML)
	if err != nil {
		return
	}

	entrypoints, err = calcEntrypoints(tasks)
	if err != nil {
		return
	}

	return p, entrypoints, nil
}

func render(pipeline []byte) ([]byte, error) {
	buf := bytes.Buffer{}
	tpl, err := template.New("").Funcs(funcMap()).Parse(string(pipeline))
	if err != nil {
		return nil, err
	}

	err = tpl.Execute(&buf, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func init() {
	files, err := builtinsFS.ReadDir("builtins")
	util.Must(err)

	for _, f := range files {
		baseName := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
		data, err := builtinsFS.ReadFile("builtins/" + f.Name())
		util.Must(err)

		builtins[baseName] = string(data)
	}
}
