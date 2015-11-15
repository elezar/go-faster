package controller

import (
	"database/sql"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
)

var Container = restful.NewContainer()

func AddResource(ws *restful.WebService) {
	Container.Add(ws)
}

func Init() {
	cors := restful.CrossOriginResourceSharing{
		AllowedMethods: []string{"GET", "PUT", "POST", "DELETE", "OPTIONS", "HEAD"},
		ExposeHeaders:  []string{restful.HEADER_AccessControlAllowOrigin, restful.HEADER_AccessControlAllowMethods},
		AllowedHeaders: []string{"Accept", "Content-Type"},
		CookiesAllowed: false,
		Container:      Container,
	}
	Container.Filter(cors.Filter)
	Container.Filter(Container.OPTIONSFilter)

	config := swagger.Config{
		WebServices:     Container.RegisteredWebServices(),
		WebServicesUrl:  "http://go-faster.devfest.com:8080/",
		ApiPath:         "/apidocs.json",
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "./swagger",
	}

	swagger.RegisterSwaggerService(config, Container)
}

func CreateHandler(response *restful.Response, json interface{}, err error) {
	if err != nil {
		Error(response, err)
		return
	}

	response.WriteAsJson(json)
}

func DeleteHandler(response *restful.Response, err error) {
	if err != nil {
		Error(response, err)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}

func UpdateHandler(response *restful.Response, json interface{}, err error) {
	if err != nil {
		Error(response, err)
		return
	}

	response.WriteAsJson(json)
}

func GetHandler(response *restful.Response, json interface{}, err error) {
	if err != nil {
		Error(response, err)
		return
	}

	response.WriteAsJson(json)
}

func Error(response *restful.Response, err error) {
	response.AddHeader("Content-Type", "text/plain")
	if err == nil {
		response.WriteHeader(http.StatusInternalServerError)
	} else if err == sql.ErrNoRows {
		response.WriteHeader(http.StatusNotFound)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}
