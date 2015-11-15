package input

import (
	"time"

	"github.com/emicklei/go-restful"
)

type Entry struct {
	Time time.Time `json:"time"`
}

func addData(request *restful.Request, response *restful.Response) {}
