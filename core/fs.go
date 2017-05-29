package core

import (
	"regexp"
	"os"
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	"time"
	"eye/integ"
)

type Fs struct {
	Name        string
	File        string
	PingRequest *QueryRequest
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
			o.pingCheck, err = o.newСheck(&QueryRequest{})
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

func (o *FsService) NewСheck(req *QueryRequest) (ret Check, err error) {
	return o.newСheck(req)
}

func (o *FsService) newСheck(req *QueryRequest) (ret *FsCheck, err error) {
	var pattern *regexp.Regexp
	if pattern, err = compilePattern(req.Expr); err == nil {
		if req.Query != "" {
			ret = &FsCheck{info: req.CheckKey("Fs"), file: filepath.Join(o.Fs.File, req.Query),
				pattern:     pattern, not: req.Not}
		} else {
			ret = &FsCheck{info: req.CheckKey("Fs"), file: o.Fs.File,
				pattern:     pattern, not: req.Not}
		}
		ret.files = integ.NewObjectCache(func() (interface{}, error) { return ret.Files() })
	}
	return
}

func (o *FsService) NewExporter(req *ExportRequest) (ret Exporter, err error) {
	return
}

//buildCheck
type FsCheck struct {
	info    string
	file    string
	not     bool
	pattern *regexp.Regexp
	files   integ.ObjectCache
}

func (o FsCheck) Info() string {
	return o.info
}

func (o FsCheck) Validate() (err error) {
	return validate(o, o.pattern, o.not)
}

func (o FsCheck) Query() (data QueryResult, err error) {
	var value interface{}
	if value, err = o.files.Get(); err == nil {
		data = value.(QueryResult)
	}
	return

}

func (o FsCheck) Files() (data QueryResult, err error) {
	var file os.FileInfo
	if file, err = os.Stat(o.file); err == nil {
		if file.IsDir() {
			var files []os.FileInfo
			files, err = ioutil.ReadDir(o.file)
			list := make([]FileInfo, len(files))
			for i, entry := range files {
				list[i] = toFileInfo(entry)
			}
			data, err = json.Marshal(list)
		} else {
			data, err = json.Marshal(toFileInfo(file))
		}
	}
	return
}

func toFileInfo(item os.FileInfo) FileInfo {
	f := FileInfo{
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
