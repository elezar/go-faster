package output

import (
	"time"

	"github.com/emicklei/go-restful"
	"github.com/pakohan/go-faster/api/controller"
	"github.com/pakohan/go-faster/api/model"
)

func getData(request *restful.Request, response *restful.Response) {
	d, err := model.GetData(time.Time{}, time.Time{})
	controller.GetHandler(response, d, err)
}
