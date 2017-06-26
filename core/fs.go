package core

import (
	"os"
	"io/ioutil"
	"path/filepath"
	"time"
	"eye/integ"
	"gopkg.in/Knetic/govaluate.v2"
	"io"
)

type Fs struct {
	Name        string
	File        string
	PingRequest *ValidationRequest
}

type FsService struct {
	Fs        *Fs
	pingCheck *FsCheck
}

func (o *FsService) Name() string {
	return o.Fs.Name
}

func (o *FsService) Init() (err error) {
	if o.pingCheck == nil {
		if o.Fs.PingRequest != nil {
			o.pingCheck, err = o.newСheck(o.Fs.PingRequest)
		} else {
			o.pingCheck, err = o.newСheck(&ValidationRequest{})
		}
		if err != nil {
			o.Close()
		}
	}
	return
}

func (o *FsService) Close() {
	o.pingCheck = nil
}

func (o *FsService) Ping() (err error) {
	if err = o.Init(); err == nil {
		err = o.pingCheck.Validate()
	}
	return
}

func (o FsService) Files(file string) (ret []*FileInfo, err error) {
	var fileInfo os.FileInfo

	if fileInfo, err = os.Stat(file); err == nil {
		if fileInfo.IsDir() {
			var files []os.FileInfo
			files, err = ioutil.ReadDir(file)
			ret = make([]*FileInfo, len(files))
			for i, entry := range files {
				ret[i] = toFileInfo(entry)
			}
		} else {
			ret = make([]*FileInfo, 1)
			ret[0] = toFileInfo(fileInfo)
		}
	}
	return
}

func (o *FsService) buildPath(pathElement string) string {
	var filePath string
	if len(pathElement) > 0 {
		filePath = filepath.Join(o.Fs.File, pathElement)
	} else {
		filePath = o.Fs.File
	}
	return filePath
}

func (o *FsService) queryToWriter(file string, writer MapWriter) (err error) {
	var items []*FileInfo
	if items, err = o.Files(file); err == nil {
		for _, fileInfo := range items {
			writer.WriteMap(fileInfo.ToMap())
		}
	}
	return
}

func (o *FsService) NewСheck(req *ValidationRequest) (ret Check, err error) {
	return o.newСheck(req)
}

func (o *FsService) newСheck(req *ValidationRequest) (ret *FsCheck, err error) {
	var eval *govaluate.EvaluableExpression
	if eval, err = compileEval(req.EvalExpr); err != nil {
		return
	}

	ret = &FsCheck{
		info: req.CheckKey("Fs"),
		file: o.buildPath(req.Query),
		eval: eval, all: req.All}
	ret.files = integ.NewObjectCache(func() (interface{}, error) { return ret.Files() })
	return
}

func (o *FsService) NewExporter(req *ExportRequest) (ret Exporter, err error) {
	ret = &fsExporter{info: req.ExportKey(o.Name()), req: req, service: o}
	return
}

//buildCheck
type FsCheck struct {
	info    string
	file    string
	all     bool
	service *FsService
	eval    *govaluate.EvaluableExpression
	files   integ.ObjectCache
}

func (o *FsCheck) Info() string {
	return o.info
}

func (o *FsCheck) Validate() (err error) {
	return validate(o, o.eval, o.all)
}

func (o *FsCheck) Query() (ret QueryResults, err error) {
	if err = o.service.Init(); err == nil {
		writer := NewQueryResultMapWriter()
		if err = o.service.queryToWriter(o.file, writer); err == nil {
			ret = writer.Data
		}
	}
	return
}

func (o *FsCheck) Files() (ret []*FileInfo, err error) {
	return o.service.Files(o.file)
}

func toFileInfo(item os.FileInfo) *FileInfo {
	f := &FileInfo{
		Name:    item.Name(),
		Size:    item.Size(),
		Mode:    item.Mode(),
		ModTime: item.ModTime(),
		IsDir:   item.IsDir(),
	}
	return f
}

type FileInfo struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	IsDir   bool
}

func (o *FileInfo) ToMap() (ret map[string]interface{}) {
	return map[string]interface{}{
		"Name": o.Name, "Size": o.Size, "Mode": o.Mode, "ModTime": o.ModTime, "IsDir": o.IsDir}
}

type fsExporter struct {
	info    string
	req     *ExportRequest
	service *FsService
}

func (o *fsExporter) Info() string {
	return o.info
}

func (o *fsExporter) Export(params map[string]string) (err error) {
	if err = o.service.Init(); err != nil {
		return
	}

	var out io.WriteCloser
	if out, err = o.req.CreateOut(params); err != nil {
		return
	}
	defer out.Close()

	err = o.service.queryToWriter(o.service.buildPath(o.req.Query), &WriteCloserMapWriter{Convert: o.req.Convert, Out: out})
	return
}
