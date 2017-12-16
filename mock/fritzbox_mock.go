package mock

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Fritz represents the mock of the FB.
// codebeat:disable[TOO_MANY_IVARS]
type Fritz struct {
	LoginChallengeResponse string
	LoginResponse          string
	DeviceList             string
	Logs                   string
	LanDevices             string
	InetStats              string
	PhoneCalls             string
	SystemStatus           string
	Server                 *httptest.Server
}

// codebeat:enable[TOO_MANY_IVARS]

// New creates a new *mock.Fritz with default configuration.
func New() *Fritz {
	return &Fritz{
		LoginChallengeResponse: "../mock/login_challenge.xml",
		LoginResponse:          "../mock/login_response_success.xml",
		DeviceList:             "../mock/devicelist.xml",
		Logs:                   "../mock/logs.json",
		LanDevices:             "../mock/landevices.json",
		InetStats:              "../mock/traffic.json",
		PhoneCalls:             "../mock/calls.csv",
		SystemStatus:           "../mock/system_status.html",
	}
}

// Start creates and starts the server.
func (f *Fritz) Start() *Fritz {
	server := f.UnstartedServer()
	server.Start()
	f.Server = server
	return f
}

// Close closes the server.
func (f *Fritz) Close() {
	f.Server.Close()
}

// UnstartedServer sets up the routes and creates a server.
func (f *Fritz) UnstartedServer() *httptest.Server {
	router := f.fritzRoutes()
	return httptest.NewUnstartedServer(router)
}

func (f *Fritz) fritzRoutes() *httprouter.Router {
	router := httprouter.New()
	router.GET("/login_sid.lua", f.loginHandler)
	router.GET("/webservices/homeautoswitch.lua", f.homeAutoHandler)
	router.GET("/query.lua", f.queryHandler)
	router.GET("/internet/inetstat_monitor.lua", f.inetStatHandler)
	router.GET("/fon_num/foncalls_list.lua", f.phoneCallsHandler)
	router.GET("/cgi-bin/system_status", f.systemStatusHandler)
	return router
}

func (f *Fritz) loginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.URL.Query().Get("response") == "" {
		f.writeFromFs(w, f.LoginChallengeResponse)
	} else {
		f.writeFromFs(w, f.LoginResponse)
	}
}

func (f *Fritz) homeAutoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !f.preProcess(w, r) {
		return
	}
	switch r.URL.Query().Get("switchcmd") {
	case "getdevicelistinfos":
		f.writeFromFs(w, f.DeviceList)
	case "setswitchon":
		w.Write([]byte("1"))
	case "setswitchoff":
		w.Write([]byte("0"))
	case "setswitchtoggle":
		w.Write([]byte("1"))
	case "sethkrtsoll":
		w.Write([]byte("OK"))
	}
}

func (f *Fritz) queryHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.URL.Query().Get("network") != "" {
		f.writeFromFs(w, f.LanDevices)
	}
	if r.URL.Query().Get("mq_log") != "" {
		f.writeFromFs(w, f.Logs)
	}
}

func (f *Fritz) inetStatHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	f.writeFromFs(w, f.InetStats)
}

func (f *Fritz) writeFromFs(w http.ResponseWriter, path string) {
	file, err := os.Open(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer file.Close()
	io.Copy(w, file)
}

func (f *Fritz) preProcess(w http.ResponseWriter, r *http.Request) bool {
	ain := r.URL.Query().Get("ain")
	if strings.Contains(strings.ToLower(ain), "fail") {
		http.Error(w, "Operation on device '"+ain+"' failed.", 500)
		return false
	}
	return true
}

func (f *Fritz) phoneCallsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	f.writeFromFs(w, f.PhoneCalls)
}

func (f *Fritz) systemStatusHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	f.writeFromFs(w, f.SystemStatus)
}
