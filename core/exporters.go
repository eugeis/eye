package core

import (
	"bytes"
	"fmt"
	"github.com/ngaut/log"
	"github.com/pkg/errors"
	"gopkg.in/Knetic/govaluate.v2"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var fileNamePattern, _ = regexp.Compile("[^a-zA-Z0-9.-]")

func (o *Eye) registerExporters() {
	for _, item := range o.config.FieldsExporter {
		if len(item.Services) > 1 {
			for _, serviceName := range item.Services {
				o.registerFieldExporter(fmt.Sprintf("%v-%v", serviceName, item.Name),
					serviceName, item)
			}
		} else if len(item.Services) > 0 {
			o.registerFieldExporter(item.Name, item.Services[0], item)
		} else {
			Log.Info("No service defined for the exporter %v", item.Name)
		}
	}
	for _, item := range o.config.FileExporter {
		if len(item.Services) > 1 {
			for _, serviceName := range item.Services {
				o.registerFileExporter(fmt.Sprintf("%v-%v", serviceName, item.Name),
					serviceName, item)
			}
		} else if len(item.Services) > 0 {
			o.registerFileExporter(item.Name, item.Services[0], item)
		} else {
			Log.Info("No service defined for the exporter %v", item.Name)
		}
	}

}

func (o *Eye) registerFileExporter(exporterFullName string, serviceName string, exporter *FileExporter) {
	var err error
	var service Service
	if service, err = o.serviceFactory.Find(serviceName); err == nil {
		var item Exporter
		var eval *govaluate.EvaluableExpression
		eval, err = compileEval(exporter.SourceFileExpr)
		request := &ExportRequest{Query: exporter.Query, EvalExpr: exporter.EvalExpr, Convert: func(row map[string]interface{}) (ret io.Reader, err error) {
			var evalResult interface{}
			if evalResult, err = eval.Eval(&MapQueryResult{row}); err == nil {
				var fileName string
				if evalResult != nil {
					if stringResult, ok := evalResult.(string); ok {
						fileName = stringResult
					}
				}
				if len(fileName) > 0 {
					ret, err = os.Create(fileName)
				} else {
					err = errors.New(fmt.Sprintf(
						"No file name computed based on source file expression '%v' and query result %v",
						exporter.SourceFileExpr, row))
				}
			}
			if err != nil {
				log.Error("Can't convert because of '%v'", err)
			}
			return
		},
			CreateOut: func(params map[string]string) (ret io.WriteCloser, err error) {
				fileName := filepath.Join(o.config.ExportFolder, fmt.Sprintf("%v_%v.zip", exporterFullName, time.Now().Unix()))
				return os.Create(fileName)
			},
		}
		if item, err = service.NewExporter(request); err == nil {
			o.exporters[exporterFullName] = item
		}
	}

}

func (o *Eye) registerFieldExporter(exporterFullName string, serviceName string, exporter *FieldsExporter) {
	var err error
	var service Service
	if service, err = o.serviceFactory.Find(serviceName); err == nil {
		var item Exporter
		request := &ExportRequest{Query: exporter.Query, Convert: func(row map[string]interface{}) (ret io.Reader, err error) {
			var line bytes.Buffer
			for _, field := range exporter.Fields {
				if val, ok := row[field]; ok {
					if str, ok := val.(string); ok {
						line.WriteString(strings.Trim(str, "\r\n"))
					} else {
						line.WriteString(fmt.Sprintf("%v", val))
					}
				} else {
					line.WriteString(" ")
				}
				line.WriteString(exporter.Separator)
			}
			line.WriteString("\n")
			return strings.NewReader(line.String()), nil
		},
			CreateOut: func(params map[string]string) (ret io.WriteCloser, err error) {
				var fileName string
				if params != nil && len(params) > 0 {
					var nameBuffer bytes.Buffer
					nameBuffer.WriteString(exporterFullName)
					for _, v := range params {
						nameBuffer.WriteString("_")
						nameBuffer.WriteString(v)
					}
					nameBuffer.WriteString(".txt")
					fileName = nameBuffer.String()
				} else {
					fileName = fmt.Sprintf("%v.txt", exporterFullName)
				}
				fileName = fileNamePattern.ReplaceAllString(fileName, "_")
				fileName = strings.Replace(fileName, "__", "_", -1)
				fileName = filepath.Join(o.config.ExportFolder, fileName)

				return os.Create(fileName)
			},
		}

		if item, err = service.NewExporter(request); err == nil {
			o.exporters[exporterFullName] = item
		}
	}
	if err != nil {
		Log.Info("Can't build exporter '%v' because: %v", exporterFullName, err)
	}
	return
}
