package rest

import (
	"net/http"

	"github.com/UArt-project/UArt-proxy/internal/service"
	"github.com/UArt-project/UArt-proxy/pkg/jsonoperations"
	"github.com/UArt-project/UArt-proxy/pkg/logger"
	"github.com/gorilla/mux"
)

// RESTApi is responsible for handling REST API requests.
type RESTApi struct {
	// The application service.
	appService service.AppService
	// Logger.
	loggr *logger.Logger
	// Router.
	router *mux.Router
}

// NewRESTApi creates a new instance of the RESTApi.
func NewRESTApi(appService service.AppService, loggr *logger.Logger) *RESTApi {
	router := mux.NewRouter()

	api := &RESTApi{
		appService: appService,
		loggr:      loggr,
		router:     router,
	}

	return api
}

// HandleFuncs registers handlers for REST API requests.
func (r *RESTApi) HandleFuncs() {
	r.router.HandleFunc("/v1/market/{page}", r.getMarketPage).Methods(http.MethodGet)
}

// ServeHTTP handles REST API requests.
func (r *RESTApi) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

// getMarketPage handles the request for getting a page of market items.
func (r *RESTApi) getMarketPage(w http.ResponseWriter, req *http.Request) {
	page, err := getPathNumber(req)
	if err != nil {
		r.loggr.Error("getting the page number from the path: %v", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	items, err := r.appService.GetMarketPage(page)
	if err != nil {
		r.loggr.Error("getting the page of items: %v", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	response := itemsToResponse(items)

	encData, err := jsonoperations.Encode(response)
	if err != nil {
		r.loggr.Error("encoding the response body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(encData)
	if err != nil {
		r.loggr.Error("writing the response body: %v", err)

		return
	}
}
