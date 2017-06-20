package core

import (
	"os"
	"io/ioutil"
	"path/filepath"
	"time"
	"eye/integ"
	"gopkg.in/Knetic/govaluate.v2"
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

func (o *FsService) NewСheck(req *ValidationRequest) (ret Check, err error) {
	return o.newСheck(req)
}

func (o *FsService) newСheck(req *ValidationRequest) (ret *FsCheck, err error) {
	var eval *govaluate.EvaluableExpression
	if eval, err = compileEval(req.EvalExpr); err != nil {
		return
	}

	if req.Query != "" {
		ret = &FsCheck{
			info: req.CheckKey("Fs"),
			file: filepath.Join(o.Fs.File, req.Query),
			eval: eval, all: req.All}
	} else {
		ret = &FsCheck{
			info: req.CheckKey("Fs"),
			file: o.Fs.File,
			eval: eval, all: req.All}
	}
	ret.files = integ.NewObjectCache(func() (interface{}, error) { return ret.Files() })

	return
}

func (o *FsService) NewExporter(req *ExportRequest) (ret Exporter, err error) {
	return
}

//buildCheck
type FsCheck struct {
	info  string
	file  string
	all   bool
	eval  *govaluate.EvaluableExpression
	files integ.ObjectCache
}

func (o FsCheck) Info() string {
	return o.info
}

func (o FsCheck) Validate() (err error) {
	return validate(o, o.eval, o.all)
}

func (o FsCheck) Query() (ret QueryResults, err error) {
	var items []*FileInfo
	if items, err = o.Files(); err == nil {
		ret = make([]QueryResult, len(items))
		for i, fileInfo := range items {
			ret[i] = &MapQueryResult{fileInfo.ToMap()}
		}
	}
	return
}

func (o FsCheck) Files() (ret []*FileInfo, err error) {
	var file os.FileInfo
	if file, err = os.Stat(o.file); err == nil {
		if file.IsDir() {
			var files []os.FileInfo
			files, err = ioutil.ReadDir(o.file)
			ret = make([]*FileInfo, len(files))
			for i, entry := range files {
				ret[i] = toFileInfo(entry)
			}
		} else {
			ret = make([]*FileInfo, 1)
			ret[0] = toFileInfo(file)
		}
	}
	return
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
