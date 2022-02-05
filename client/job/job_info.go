package job

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	pb "github.com/plarun/scheduler/client/data"
	"github.com/plarun/scheduler/client/model"
)

// JobInfo wraps the SubmitJilClient and manages parsing and prevalidation on JIL
type JobInfoController struct {
	client pb.SubmitJilClient
}

func NewJobInfoController(client pb.SubmitJilClient) *JobInfoController {
	return &JobInfoController{client: client}
}

// SubmitJil submits the JIL to grpc server after parsing and building the Job info from JIL
func (controller JobInfoController) SubmitJil(inputFilename string) error {
	log.Println("jil submitted for parsing...")
	// parse the raw JIL
	parsedJils, err := controller.Parse(inputFilename)
	if err != nil {
		return err
	}

	jilList := []*pb.Jil{}
	for _, parsedJil := range parsedJils {
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
func (controller JobInfoController) Parse(inputFile string) ([]model.JilData, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, fmt.Errorf("file: %s doesn't exist", inputFile)
	}
	defer file.Close()

	// singleParseInProgress indicates whether parser is parsing a job
	singleParseInProgress := false

	var parsedJils []map[string]string
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
				parsedJils = append(parsedJils, chunk)
				chunk = make(map[string]string)
			}
			continue // ignore spaces between different jil data
		}

		if line[:2] == "/*" {
			commentLine := strings.TrimRight(line, " \t")
			commentClose := commentLine[len(commentLine)-2:]
			if commentClose == "*/" {
				if singleParseInProgress {
					return nil, controller.logErr(lineNum, line, "comment should only be mentioned in begining of attribures")
				}
				continue // ignore comment on top of each jil data
			}
		}

		if len(line) != 0 {
			parsedLine := strings.SplitN(line, ":", 2)
			if len(parsedLine) < 2 {
				return nil, controller.logErr(lineNum, line, "line unparsable")
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
		parsedJils = append(parsedJils, chunk)
	}

	var jilList []model.JilData
	for _, parsedJil := range parsedJils {
		builder := JobInfoBuilder{parsedJil: parsedJil}
		jil, err := builder.buildJil()
		if err != nil {
			return nil, err
		}
		jilList = append(jilList, jil)
	}

	return jilList, nil
}

// logErr returns error with line number, line content and error message
func (JobInfoController) logErr(lineNum int64, line string, errMsg string) error {
	return fmt.Errorf("line no: %d\nline: %s\nerror: %s", lineNum, line, errMsg)
}

// actionAttribute checks if attribute is one of valid JIL action
func (JobInfoController) actionAttribute(attribute string) (string, bool) {
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
func (controller JobInfoController) valueAttribute(attribute string) bool {
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
