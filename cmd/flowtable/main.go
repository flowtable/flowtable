package main

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	restful "github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
	"log"
	"net/http"
)

type IndexResponse struct {
	Message string
}

func NewWebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/hello").To(index).Doc("hello page").Returns(200, "OK", IndexResponse{}))
	return ws
}
func main() {
	restful.DefaultContainer.Add(NewWebService())

	config := restfulspec.Config{
		WebServices:                   restful.RegisteredWebServices(), // you control what services are visible
		APIPath:                       "/api.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject}
	restful.DefaultContainer.Add(restfulspec.NewOpenAPIService(config))

	// Optionally, you can install the Swagger Service which provides a nice Web UI on your REST API
	// You need to download the Swagger HTML5 assets and change the FilePath location in the config below.
	// Open http://localhost:8080/apidocs/?url=http://localhost:8080/apidocs.json
	//http.Handle("/api/", http.StripPrefix("/api/", http.FileServer(http.Dir("/Users/emicklei/Projects/swagger-ui/dist"))))

	log.Printf("start listening on http://localhost:8088")
	log.Fatal(http.ListenAndServe(":8088", nil))
}
func index(request *restful.Request, response *restful.Response) {
	response.WriteAsJson(map[string]string{"hello": "world"})
}
func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "FlowTable",
			Description: "hyper table with workflow,access control,expansion",
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{
					Name:  "Leon",
					Email: "leondevlifelog@gmail.com",
					URL:   "https://flowtable.cn",
				},
			},
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "Apache License 2.0",
					URL:  "https://www.apache.org/licenses/LICENSE-2.0.txt",
				},
			},
			Version: "0.0.1",
		},
	}
	swo.Tags = []spec.Tag{spec.Tag{TagProps: spec.TagProps{
		Name:        "example",
		Description: "example description"}}}
}
