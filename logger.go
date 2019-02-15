package goSeasLog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
)

type SeasLog struct {
	Out, DebugOut, WarnOut, InfoOut, ErrorOut, FatalOut io.Writer
	//
	lock sync.Mutex
	//
	// %L - Level 日志级别。
	// %M - Message 日志信息。
	// %T - DateTime 如2017-08-16 19:15:02，受seaslog.default_datetime_format影响。
	// %H - HostName 主机名。
	// %P - ProcessId 进程ID。
	// %F - FileName:LineNo 文件名:行号
	// LogFormat string = "%L %P %T %Q %H %F %M"
	// LogFormat string = "%T %L %P %Q %M %H %F"
	// LogFormat  string = "%T %L %P %Q %M"
	LogFormat string
	//
	Level int
	//
	Gap string
	//
	TerminalOut bool
	//
	LogPath       string
	LogfileFormat string
	LogSplit      bool
	//
	DateFormat string
	//
	HostName string
}

const (
	Debug = 1 << iota
	Info
	Warn
	Error
	Fatal
)

const (
	F_Flag     int         = os.O_RDWR | os.O_CREATE | os.O_APPEND
	F_ModePerm os.FileMode = 0644 //os.ModePerm
)

func levelString(v int) string {
	switch v {
	case Debug:
		return "Debug"
	case Info:
		return "Info"
	case Warn:
		return "Warn"
	case Error:
		return "Error"
	case Fatal:
		return "Fatal"
	default:
		return "None"
	}
}

func New() *SeasLog {
	return &SeasLog{
		Out:         os.Stderr,
		LogFormat:   "%L %P %T %Q %H %F %M",
		Level:       Debug,
		Gap:         " | ",
		TerminalOut: true,
		DateFormat:  "2006/1/2 15:04:05",
		HostName:    "Local",
	}
}

func (this *SeasLog) Output(level int, message ...interface{}) {

	var _code_file string
	var _code_line int
	var _ok bool
	_, _code_file, _code_line, _ok = runtime.Caller(1)
	if _ok == false {
		_code_file = "???"
		_code_line = 0
	}
	_, _short_file := path.Split(_code_file)

	// for i := len(_code_file) - 1; i > 0; i-- {
	// 	if _code_file[i] == '/' {
	// 		_short_file = _code_file[i+1:]
	// 		break
	// 	}
	// }

	F := map[string][]byte{
		"%L": []byte(levelString(level)),
		// "%Q": GetGID(),
		"%P": GetGID(),
		"%M": *this.formatMsg(message...),
		"%T": []byte(time.Now().Format(this.DateFormat)),
		"%H": []byte(this.HostName),
		"%F": []byte(fmt.Sprintf("%s:%d", _short_file, _code_line)),
	}
	_format := strings.Split(this.LogFormat, " ")

	_log_line := new([]byte)

	for _, _val := range _format {
		if __val := F[_val]; __val != nil {
			if len(*_log_line) == 0 {
				*_log_line = append([]byte(fmt.Sprintf("%s", F[_val])))
			} else {
				*_log_line = append(*_log_line, []byte(this.Gap)...)
				*_log_line = append(*_log_line, F[_val]...)
			}
		}
	}
	if len(*_log_line) > 0 {
		*_log_line = append(*_log_line, '\n')
	}

	//
	if this.TerminalOut {
		this.Out.Write(*_log_line)
	}
	this.writerLog(level, _log_line)
}

func (this *SeasLog) SetLevel(level int) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.Level = level
}

func (this *SeasLog) SetLogFormat(formatter string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.LogFormat = formatter
}


func (this *SeasLog) SetDateFormat(dateFormat string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.DateFormat = dateFormat
}

func (this *SeasLog) SetGap(char string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.Gap = char
}

func (this *SeasLog) SetHostName(host string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.HostName = host 
}

func (this *SeasLog) SetTerminalOut(b bool) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.TerminalOut = b
}

//demo
//.SetLogFile("/var/logs", "2006010215", false)
func (this *SeasLog) SetLogFile(logPath, logfileFormat string, logSplit bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.LogPath = logPath
	this.LogfileFormat = logfileFormat
	this.LogSplit = logSplit
}

func (this *SeasLog) formatLogFileName(prefix string) string {
	var _log_file_name string
	if this.LogPath != "" {
		if this.LogSplit == true {
			_log_file_name = path.Join(this.LogPath, "/", prefix+time.Now().Format(this.LogfileFormat)+".log")
		} else {
			_log_file_name = path.Join(this.LogPath, "/", "goSeasLog"+time.Now().Format(this.LogfileFormat)+".log")
		}
	}
	return _log_file_name
}

func (this *SeasLog) writerLog(level int, log_message *[]byte) {

	_fileName := this.formatLogFileName(levelString(level))

	if PathExists(_fileName) == false {
		switch {
		case level == Debug:
			this.DebugOut = nil
		case level == Info:
			this.InfoOut = nil
		case level == Warn:
			this.WarnOut = nil
		case level == Error:
			this.ErrorOut = nil
		case level == Fatal:
			this.FatalOut = nil
		}
	}

	switch {
	case level == Debug:
		if this.DebugOut == nil {
			this.lock.Lock()
			if _f, err := os.OpenFile(_fileName, F_Flag, F_ModePerm); err == nil {
				this.DebugOut = _f
			}
			this.lock.Unlock()
		}
		if this.DebugOut != nil {
			this.DebugOut.Write(*log_message)
		}
	case level == Info:
		if this.InfoOut == nil {
			this.lock.Lock()
			if _f, err := os.OpenFile(_fileName, F_Flag, F_ModePerm); err == nil {
				this.InfoOut = _f
			}
			this.lock.Unlock()
		}
		if this.InfoOut != nil {
			this.InfoOut.Write(*log_message)
		}
	case level == Warn:
		if this.WarnOut == nil {
			this.lock.Lock()
			if _f, err := os.OpenFile(_fileName, F_Flag, F_ModePerm); err == nil {
				this.WarnOut = _f
			}
			this.lock.Unlock()
		}
		if this.WarnOut != nil {
			this.WarnOut.Write(*log_message)
		}
	case level == Error:
		if this.ErrorOut == nil {
			this.lock.Lock()
			if _f, err := os.OpenFile(_fileName, F_Flag, F_ModePerm); err == nil {
				this.ErrorOut = _f
			}
			this.lock.Unlock()
		}
		if this.ErrorOut != nil {
			this.ErrorOut.Write(*log_message)
		}
	case level == Fatal:
		if this.FatalOut == nil {
			this.lock.Lock()
			if _f, err := os.OpenFile(_fileName, F_Flag, F_ModePerm); err == nil {
				this.DebugOut = _f
			}
			this.lock.Unlock()
		}
		if this.FatalOut != nil {
			this.FatalOut.Write(*log_message)
		}
	}

}

func (this *SeasLog) getLogfileWriter(filename string) io.Writer {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil
	} else {
		return file
	}
}

func (this *SeasLog) Debug(v ...interface{}) {
	if this.Level <= Debug {
		this.Output(Debug, v...)
	}
}

func (this *SeasLog) Info(v ...interface{}) {
	if this.Level <= Info {
		// this.Output("\033[34m"+levelString(Info)+"\033[0m", v...)
		this.Output(Info, v...)
	}
}

func (this *SeasLog) Warn(v ...interface{}) {
	if this.Level <= Warn {
		// this.Output("\033[33m"+levelString(Warn)+"\033[0m", v...)
		this.Output(Warn, v...)
	}
}

func (this *SeasLog) Error(v ...interface{}) {
	if this.Level <= Error {
		// this.Output("\033[31m"+levelString(Error)+"\033[0m", v...)
		this.Output(Error, v...)
	}
}

func (this *SeasLog) Fatal(v ...interface{}) {
	if this.Level <= Fatal {
		this.Output(Fatal, v...)
		os.Exit(1)
	}
}

func (this *SeasLog) formatMsg(v ...interface{}) *[]byte {

	_msg := new([]byte)

	for _, _val := range v {
		vvv := fmt.Sprintf("%+v", _val)
		vvv = strings.Replace(vvv, "\n", " ", -1)

		if len(*_msg) == 0 {
			*_msg = append(*_msg, []byte(vvv)...)
		} else {
			*_msg = append(*_msg, []byte(this.Gap)...)
			*_msg = append(*_msg, []byte(vvv)...)
		}
	}

	return _msg
}

func GetGID() []byte {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]

	return b
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
