package filter

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/webgpu/gputeachingkit-labbuilder/pkg"
	pf "github.com/webgpu/gputeachingkit-labbuilder/pkg/pandocfilter"
	log "github.com/sirupsen/logrus"
	"github.com/Unknwon/com"
	"github.com/mitchellh/go-homedir"
)

func toJSON(inputFilePath string) (string, error) {
	tmpDir := os.TempDir()
	outputFile := filepath.Join(tmpDir, "pandocJsonOutput.json")
	log.Debug("Generating pandoc json from markdown file")
	args := []string{
		"-o",
		outputFile,
		"-f",
		pandoc.MarkdownFormat,
		"-t",
		"json",
		inputFilePath,
	}
	args = append(args, pandoc.DefaultFilter...)
	cmd := exec.Command("pandoc", args...)
	cmd.Dir = tmpDir
	buf, err := cmd.CombinedOutput()
	log.WithError(err).WithField("command_out", string(buf)).Debug("Ran pandoc to json command")
	if err != nil {
		return "", errors.New(string(buf) + " .. Error: " + err.Error())
	}
	return outputFile, err
}

func fromJSON(outputFilePath, inputFilePath string) error {
	tmpDir := os.TempDir()
	args := []string{
		"-o",
		outputFilePath,
		"-f",
		"json",
		"-t",
		pandoc.MarkdownFormat,
		"-S",
		"-s",
		inputFilePath,
	}
	args = append(args, pandoc.DefaultFilter...)
	cmd := exec.Command("pandoc", args...)
	cmd.Dir = tmpDir
	buf, err := cmd.CombinedOutput()
	log.WithError(err).WithField("command_out", string(buf)).Debug("Ran pandoc to markdown command")
	if err != nil {
		return errors.New(string(buf) + " .. Error: " + err.Error())
	}
	return err
}

func fileName(pth string) string {
	base := filepath.Base(pth)
	return base[:len(base)-len(filepath.Ext(pth))]
}

func isMarkdownExt(pth string) bool {
	return filepath.Ext(pth) == ".markdown" ||
		filepath.Ext(pth) == ".md"
}

type pandocJSON struct {
	Blocks []struct {
		T string        `json:"t"`
		C []interface{} `json:"c"`
	} `json:"blocks"`
	PandocAPIVersion []int `json:"pandoc-api-version"`
	Meta             struct {
		Module struct {
			T string `json:"t"`
			C string `json:"c"`
		} `json:"module"`
		Author struct {
			T string `json:"t"`
			C []struct {
				T string `json:"t"`
				C string `json:"c,omitempty"`
			} `json:"c"`
		} `json:"author"`
		Title struct {
			T string `json:"t"`
			C []struct {
				T string `json:"t"`
				C string `json:"c,omitempty"`
			} `json:"c"`
		} `json:"title"`
	} `json:"meta"`
}

func Filter(outputFileDir, inputFilePath string, format string) (string, error) {

	var (
		doc               pandocJSON
		jsonInputFilePath string
	)

	inputFilePath, _ = homedir.Expand(inputFilePath)
	inputFilePath, _ = filepath.Abs(inputFilePath)
	jsonOutpuFilePath := filepath.Join(outputFileDir, fileName(inputFilePath)+"-filter.json")
	outputFilePath := filepath.Join(outputFileDir, fileName(inputFilePath)+"-filter.markdown")
	log.Debug("Input file is set to " + inputFilePath)
	log.Debug("Output file is set to " + outputFilePath)

	if !com.IsFile(inputFilePath) {
		return "", errors.New("input file does not exist")
	}

	if isMarkdownExt(inputFilePath) {
		var err error
		log.Debug("File has a markdown extension... converting to JSON format.")
		jsonInputFilePath, err = toJSON(inputFilePath)
		if err != nil {
			return "", err
		}
	} else {
		jsonInputFilePath = inputFilePath
	}

	pandoc.Clear()

	inputData, err := ioutil.ReadFile(jsonInputFilePath)
	if err != nil {
		return "", err
	}

	if err = json.Unmarshal(inputData, &doc); err != nil {
		return "", err
	}

	meta := doc.Meta
	var ndoc interface{}
	for _, filter := range pandoc.Filters {
		ndoc = pf.Walk(doc, filter, format, meta)
	}
	b, err := json.Marshal(ndoc)
	if err != nil {
		return "", err
	}
	if err := ioutil.WriteFile(jsonOutpuFilePath, b, 0644); err != nil {
		return "", err
	}

	if isMarkdownExt(inputFilePath) {
		if err := fromJSON(outputFilePath, jsonOutpuFilePath); err != nil {
			return "", err
		}
	}

	return outputFilePath, nil
}
