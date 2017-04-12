package core

import (
	"errors"
	"fmt"
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
	if err = o.Init(); err != nil {
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
				pattern:     pattern}
		} else {
			ret = &FsCheck{info: req.CheckKey("Fs"), file: o.Fs.File,
				pattern:     pattern}
		}
		ret.files = integ.NewObjectCache(func() (interface{}, error) { return ret.Files() })
	}
	return
}

//buildCheck
type FsCheck struct {
	info    string
	file    string
	pattern *regexp.Regexp
	files   integ.ObjectCache
}

func (o FsCheck) Info() string {
	return o.info
}

func (o FsCheck) Validate() (err error) {
	data, err := o.Query()

	if err == nil && o.pattern != nil {
		if !o.pattern.Match(data) {
			err = errors.New(fmt.Sprintf("No match for %v", o.info))
		}
	}
	return
}

func (o FsCheck) Query() (data QueryResult, err error) {
	var value interface{}
	value, err = o.files.Get()
	data = value.(QueryResult)
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
	l.Info("%s", data)
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
