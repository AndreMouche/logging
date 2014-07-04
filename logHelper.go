package logging

import (
	"os"
	"sort"
	"strings"
	"time"
)

const (
	LogSuffixFormat      = "2006-01-02"  // 15:04:05.999999999"
	LogFileCheckDuration = 1 * time.Hour //1 * time.Second
	LogFileMaxSize       = 1024 * 1024
	MaxLogFiles          = 30
)

func InitLogToFile(logDir string, namePrefix string) error {
	const LogFileOpt = os.O_RDWR | os.O_CREATE | os.O_APPEND
	fp, err := os.OpenFile(logDir+"/"+namePrefix+"."+time.Now().Format(LogSuffixFormat), LogFileOpt, 0666)
	if err != nil {
		return err
	}
	SetHighlighting(false)
	SetOutput(fp)
	go updateByTime(logDir, namePrefix, fp)
	println("initlog")
	return nil
}

/**
  每隔一小时检测一次
*/
func updateByTime(logDir string, namePrefix string, lastFp *os.File) {
	const LogFileOpt = os.O_RDWR | os.O_CREATE | os.O_APPEND
	tickChan := time.Tick(LogFileCheckDuration)

	for {
		<-tickChan
		path := logDir + "/" + namePrefix + "." + time.Now().Format(LogSuffixFormat)
		finfo, err := os.Stat(path)
		if err != nil || finfo.Size() > LogFileMaxSize {
			fp, err := os.OpenFile(path, LogFileOpt, 0666)
			if err != nil {
				println("[ERROR]", err)
				continue
			}
			SetOutput(fp)
			lastFp.Close()
			lastFp = fp
			cleanOldLogs(logDir, namePrefix)
		}

	}
}

func cleanOldLogs(logDir string, namePrefix string) {
	dir, err := os.Open(logDir)
	if err != nil {
		println("[ERROR]", err)
		return
	}
	files, err := dir.Readdirnames(0)
	logFiles := make([]string, 0)
	for _, file := range files {
		if strings.HasPrefix(file, namePrefix) {
			logFiles = append(logFiles, file)
		}
	}

	sort.Sort(sort.StringSlice(logFiles))
	deNum := len(logFiles) - MaxLogFiles
	for id := 0; id < deNum; id++ {
		curPath := logDir + "/" + logFiles[id]
		err := os.Remove(curPath)
		if err != nil {
			println("[ERROR]", err)
		}
	}

}
