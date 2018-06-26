package server

import (
	"github.com/emicklei/go-restful"
	"github.com/gorilla/mux"
	handler "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/serviceprovider"
)

func (a *App) AppRoute() *mux.Router {
	router := mux.NewRouter()

	container := restful.NewContainer()

	container.Filter(globalLogging)

	container.Add(newVersionService(a.ServiceProvider))
	container.Add(newNetworkService(a.ServiceProvider))
	container.Add(newStorageProviderService(a.ServiceProvider))
	container.Add(newVolumeService(a.ServiceProvider))

	router.PathPrefix("/v1/").Handler(container)
	return router
}

func newVersionService(sp *serviceprovider.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/versions").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	webService.Route(webService.GET("/").To(versionHandler(sp)))
	return webService
}

func newNetworkService(sp *serviceprovider.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/networks").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	webService.Route(webService.GET("/").To(handler.RESTfulServiceHandler(sp, listNetworkHandler)))
	webService.Route(webService.GET("/{id}").To(handler.RESTfulServiceHandler(sp, getNetworkHandler)))
	webService.Route(webService.POST("/").To(handler.RESTfulServiceHandler(sp, createNetworkHandler)))
	webService.Route(webService.DELETE("/{id}").To(handler.RESTfulServiceHandler(sp, deleteNetworkHandler)))
	return webService
}

func newStorageProviderService(sp *serviceprovider.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/storageprovider").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	webService.Route(webService.POST("/").To(handler.RESTfulServiceHandler(sp, createStorageProvider)))
	webService.Route(webService.GET("/").To(handler.RESTfulServiceHandler(sp, listStorageProvider)))
	webService.Route(webService.DELETE("/{id}").To(handler.RESTfulServiceHandler(sp, deleteStorageProvider)))
	return webService
}

func newVolumeService(sp *serviceprovider.Container) *restful.WebService {
	webService := new(restful.WebService)
	webService.Path("/v1/volume").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	webService.Route(webService.POST("/").To(handler.RESTfulServiceHandler(sp, createVolume)))
	return webService
}
