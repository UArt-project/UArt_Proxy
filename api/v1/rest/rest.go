package rest

import (
	"net/http"

	"github.com/UArt-project/UArt-proxy/domain/authdomain"
	"github.com/UArt-project/UArt-proxy/internal/service"
	"github.com/UArt-project/UArt-proxy/pkg/jsonoperations"
	"github.com/UArt-project/UArt-proxy/pkg/logger"
	"github.com/gorilla/mux"
)

// API is responsible for handling REST API requests.
type API struct {
	// The application service.
	appService service.AppService
	// Logger.
	loggr *logger.Logger
	// Router.
	router *mux.Router
}

// NewAPI creates a new instance of the API.
func NewAPI(appService service.AppService, loggr *logger.Logger) *API {
	router := mux.NewRouter()

	api := &API{
		appService: appService,
		loggr:      loggr,
		router:     router,
	}

	api.HandleFunc()

	return api
}

// HandleFunc registers handlers for REST API requests.
func (r *API) HandleFunc() {
	r.router.HandleFunc("/v1/market/{page}", r.getMarketPage).Methods(http.MethodGet)
	r.router.HandleFunc("/v1/auth", r.getAuth).Methods(http.MethodGet)
	r.router.HandleFunc("/login/oauth2/code/google", r.getAuthCallback).Methods(http.MethodGet)
}

// ServeHTTP handles REST API requests.
func (r *API) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

// getMarketPage handles the request for getting a page of market items.
func (r *API) getMarketPage(responseWriter http.ResponseWriter, req *http.Request) {
	page, err := getPathNumber(req)
	if err != nil {
		r.loggr.Error("getting the page number from the path: %v", err)
		responseWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	items, err := r.appService.GetMarketPage(page)
	if err != nil {
		r.loggr.Error("getting the page of items: %v", err)
		responseWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	response := itemsToResponse(page, items)

	encData, err := jsonoperations.Encode(response)
	if err != nil {
		r.loggr.Error("encoding the response body: %v", err)
		responseWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)

	_, err = responseWriter.Write(encData)
	if err != nil {
		r.loggr.Error("writing the response body: %v", err)

		return
	}
}

// getAuth handles the request for getting the auth url.
func (r *API) getAuth(responseWriter http.ResponseWriter, req *http.Request) {
	url, err := r.appService.GetAuthPage()
	if err != nil {
		r.loggr.Error("getting the auth page: %v", err)
		responseWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	// redirect to the auth url
	responseWriter.Header().Set("Location", url)
	responseWriter.WriteHeader(http.StatusSeeOther)
}

// getAuthCallback handles the request for getting the auth token.
func (r *API) getAuthCallback(w http.ResponseWriter, req *http.Request) {
	// cookies := req.Cookies()

	// cookie := ""

	// if len(cookies) <= 1 {
	// r.loggr.Error("no cookies in the request")
	// w.WriteHeader(http.StatusBadRequest)

	// return
	// } else {
	// cookie = cookies[1].Value
	// }

	callbackData := authdomain.CallbackRequest{
		// Cookie:   cookie,
		State:    req.URL.Query().Get("state"),
		Code:     req.URL.Query().Get("code"),
		Scope:    req.URL.Query().Get("scope"),
		AuthUser: req.URL.Query().Get("authuser"),
		Prompt:   req.URL.Query().Get("prompt"),
	}

	token, err := r.appService.GetAuthToken(callbackData)
	if err != nil {
		r.loggr.Error("getting the auth token: %v", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	// send token.RespCode to the w
	// w.Header().Set("Location", token.RedirectURL)
	// w.Header().Add("Set-Cookie", token.RespCode.Header.Get("Set-Cookie"))

	// redirect to the redirect url
	// w.WriteHeader(token.RespCode.StatusCode)

	// redirect to the redirect url and set the token
	w.Header().Set("token", token.AuthToken)
	w.Header().Set("Location", token.RedirectURL)
	w.WriteHeader(http.StatusSeeOther)
}
