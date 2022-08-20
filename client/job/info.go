package job

import (
	"bufio"
	"context"
	"fmt"
	builder2 "github.com/plarun/scheduler/client/job/builder"
	"log"
	"os"
	"strings"

	pb "github.com/plarun/scheduler/client/data"
	"github.com/plarun/scheduler/client/model"
)

// InfoController wraps the SubmitJilClient and manages parsing and pre validation on JIL
type InfoController struct {
	client pb.SubmitJilClient
}

func NewJobInfoController(client pb.SubmitJilClient) *InfoController {
	return &InfoController{client: client}
}

// SubmitJil submits the JIL to grpc server after parsing and building the Job info from JIL
func (controller *InfoController) SubmitJil(inputFilename string) error {
	log.Println("jil submitted for parsing...")
	// parse the raw JIL
	parsedJil, err := controller.Parse(inputFilename)
	if err != nil {
		return err
	}

	var jilList []*pb.Jil
	for _, parsedJil := range parsedJil {
		jil := &pb.Jil{}
		jil.Data = &pb.JilData{}

		if parsedJil.Action == model.DELETE {
			jil.Action = pb.JilAction_DELETE
			jil.Data.JobName = parsedJil.JobName
		} else {
			if parsedJil.Action == model.INSERT {
				jil.Action = pb.JilAction_INSERT
			} else if parsedJil.Action == model.UPDATE {
				jil.Action = pb.JilAction_UPDATE
			}

			jil.Data.JobName = parsedJil.JobName
			jil.Data.Command = parsedJil.Command
			jil.Data.Conditions = parsedJil.Conditions
			jil.Data.StdOut = parsedJil.StdOutLog
			jil.Data.StdErr = parsedJil.StdErrLog
			jil.Data.Machine = parsedJil.Machine
			jil.Data.RunDays = parsedJil.RunDays
			jil.Data.StartTimes = parsedJil.StartTimes
			jil.AttributeFlag = parsedJil.AttributeFlag
		}

		jilList = append(jilList, jil)
	}

	submitReq := &pb.SubmitJilReq{
		Jil: jilList,
	}

	for _, jil := range jilList {
		log.Println(jil.Action, jil.Data)
	}

	res, err := controller.client.Submit(context.Background(), submitReq)
	if err != nil {
		return err
	}

	log.Println(res)
	return nil
}

// Parse parses the content in file at path inputFile
func (controller *InfoController) Parse(inputFile string) ([]model.JilData, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, fmt.Errorf("file: %s doesn't exist", inputFile)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err.Error())
		}
	}(file)

	// singleParseInProgress indicates whether parser is parsing a job
	singleParseInProgress := false

	var parsedJil []map[string]string
	chunk := make(map[string]string)

	// lineNum tracks current line number of file
	var lineNum int64 = 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		if line == "" {
			if singleParseInProgress {
				singleParseInProgress = false
				parsedJil = append(parsedJil, chunk)
				chunk = make(map[string]string)
			}
			continue // ignore spaces between different jil data
		}

		if line[:2] == "/*" {
			commentLine := strings.TrimRight(line, " \t")
			commentClose := commentLine[len(commentLine)-2:]
			if commentClose == "*/" {
				if singleParseInProgress {
					return nil, controller.logErr(lineNum, line, "comment should only be mentioned in beginning of attributes")
				}
				continue // ignore comment on top of each jil data
			}
		}

		if len(line) != 0 {
			parsedLine := strings.SplitN(line, ":", 2)
			if len(parsedLine) < 2 {
				return nil, controller.logErr(lineNum, line, "line is not parsable")
			}

			attribute := parsedLine[0]
			value := parsedLine[1]
			// remove leading spaces or leading tabs
			value = strings.TrimLeft(value, " \t")

			action, isActionAttribute := controller.actionAttribute(attribute)
			if isActionAttribute {
				if singleParseInProgress {
					return nil, controller.logErr(lineNum, line, "non standard format")
				}
				singleParseInProgress = true
				if value == "" {
					return nil, controller.logErr(lineNum, line, "jil action shouldn't be empty")
				}
				chunk["action"] = action
				chunk["job_name"] = value
			} else if !singleParseInProgress {
				return nil, controller.logErr(lineNum, line, "jil attributes of same job shouldn't be separated by empty line")
			}

			if !controller.valueAttribute(attribute) {
				return nil, controller.logErr(lineNum, line, "invalid attribute")
			} else {
				if !isActionAttribute {
					chunk[attribute] = value
				}
			}
		}
	}

	if singleParseInProgress {
		parsedJil = append(parsedJil, chunk)
	}

	var jilList []model.JilData
	for _, parsedJil := range parsedJil {
		builder := builder2.InfoBuilder{ParsedJil: parsedJil}
		jil, err := builder.BuildJil()
		if err != nil {
			return nil, err
		}
		jilList = append(jilList, jil)
	}

	return jilList, nil
}

// logErr returns error with line number, line content and error message
func (*InfoController) logErr(lineNum int64, line string, errMsg string) error {
	return fmt.Errorf("line no: %d\nline: %s\nerror: %s", lineNum, line, errMsg)
}

// actionAttribute checks if attribute is one of valid JIL action
func (*InfoController) actionAttribute(attribute string) (string, bool) {
	if attribute == "insert" {
		return "insert", true
	} else if attribute == "update" {
		return "update", true
	} else if attribute == "delete" {
		return "delete", true
	} else {
		return "", false
	}
}

// valueAttribute checks if attribute is not actionAttribute but one of the valid JIL attributes
func (controller *InfoController) valueAttribute(attribute string) bool {
	_, ok := controller.actionAttribute(attribute)
	if ok {
		return true
	}
	for _, staticAttr := range model.StaticAttributes {
		if attribute == staticAttr {
			return true
		}
	}
	return false
}
