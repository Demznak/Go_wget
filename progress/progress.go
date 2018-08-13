package progress

type Progress struct {
	progress map[string]int64
}

var progress Progress

func init() {
	progress.progress = make(map[string]int64)
}

func GetProgress() *Progress {
	return &progress
}

func (pr *Progress) GetFileNames() []string {
	files := []string{}
	for key, _ := range progress.progress {
		files = append(files, key)
	}
	return files
}

func (pr *Progress) GetProgressFile(fileName string) (int64, bool) {
	if val, ok := progress.progress[fileName]; ok {
		return val, true
	}
	return 0, false
}

func (pr *Progress) SetProgressFile(fileName string, value int64) {
	progress.progress[fileName] = value
}
