package youtubedownload

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	youtube "github.com/kkdai/youtube"
	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	// currentFile, _ := filepath.Abs(os.Args[0])
	if logerr := initlog(); logerr != nil {
		logrus.Error(logerr)
	}
	currentFile, _ := filepath.Abs("./a.mp4")
	logrus.Println("download to file=", currentFile)

	// NewYoutube(debug) if debug parameter will set true we can log of messages
	y := youtube.NewYoutube(true, true)
	y.DecodeURL("https://www.youtube.com/watch?v=rFejpH_tAHM")
	y.StartDownload(currentFile)
}

func findFunc(f *runtime.Frame) (string, string) {
	s := strings.Split(f.Function, ".")
	funcname := s[len(s)-1]

	return funcname, ""
}

var fieldSeq = map[string]int{
	"time":  0,
	"level": 1,
	"func":  2,
}

func sortCustom(fields []string) {
	sort.Slice(fields, func(i, j int) bool {
		if fields[i] == "msg" {
			return false
		}
		if iIdx, oki := fieldSeq[fields[i]]; oki {
			if jIdx, okj := fieldSeq[fields[j]]; okj {
				return iIdx < jIdx
			}
			return true
		}
		return false
	})
}
func initlog() error {
	//Create directory
	logrusFormatter := &logrus.TextFormatter{
		DisableColors:    true,
		CallerPrettyfier: findFunc,
		ForceQuote:       true,
		DisableSorting:   false,
		SortingFunc:      sortCustom,
	}
	logrus.SetReportCaller(true)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(logrusFormatter)

	currdir, direrr := os.Getwd()
	if direrr != nil {
		return direrr
	}

	currdir += "/log"
	err := os.MkdirAll(fmt.Sprintf("%s", currdir), os.ModePerm)
	if err != nil {
		logrus.Error(err)
		return err
	}

	filename := fmt.Sprintf("%s/logrus.log", currdir)

	logrus.SetOutput(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     3,
		Compress:   false,
	})

	// custom reposrtcaller on

	logrus.Infof("filename : %s", filename)
	return nil
}
