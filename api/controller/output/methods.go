package output

import "github.com/emicklei/go-restful"

type DataContainer struct {
	Data [][]interface{} `json:"data"`
}

func getData(request *restful.Request, response *restful.Response) {}
