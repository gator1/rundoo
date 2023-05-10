package log

import (
	"fmt"
	"io/ioutil"
	stlog "log"
	"net/http"
	"os"
)

var log *stlog.Logger

type fileLog string

func (fl fileLog) Write(data []byte) (int, error) {
	f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		fmt.Println("cannot open app.log")
		return 0, err
	}
	defer f.Close()
	return f.Write(data)
}

func Run(destination string) {
	log = stlog.New(fileLog(destination), "", stlog.LstdFlags)
}

type LogHandler struct {}

func (lh *LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    msg, err := ioutil.ReadAll(r.Body)
    if err != nil || len(msg) == 0 {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    write(string(msg))
}

func HttpHandler() {
    http.Handle("/log", &LogHandler{})
}

func write(message string) {
	log.Printf("%v\n", message)
}
