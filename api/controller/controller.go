package controller

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
)

var Container = restful.NewContainer()

func AddResource(ws *restful.WebService) {
	Container.Add(ws)
}

type fs struct {
	d http.Dir
}

func (self fs) Open(name string) (http.File, error) {
	log.Println(name)
	return self.d.Open(name)
}

func Init() {
	fs := &fs{http.Dir("../views")}
	Container.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(fs)))

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
		WebServicesUrl:  "http://go-faster.devfest.com:8080",
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
