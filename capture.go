package middleware

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"net/http/httptest"
)

/*
Captures http.Response data from an httprouter.Handle.
Suitable for copying the response out to multiple sources.

Example:
    func copyResp(resp *httptest.ResponseRecorder, req *http.Request) {
        log.Println(resp.StatusCode)
    }
    func handler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
        io.WriteString(w, "OK")
    }
    h := middleware.CaptureResponse(handler, copyResp)
    router := httprouter.New()
    router.Get("/", h)
    http.ListenAndServe(":80", router)
*/
func CaptureResponse(handler httprouter.Handle, fn func(*httptest.ResponseRecorder, *http.Request)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w2 := httptest.NewRecorder()
		handler(w2, r, p)

		// write out to original writer
		w.WriteHeader(w2.Code)
		for k, values := range w2.HeaderMap {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
		_, err := w.Write(w2.Body.Bytes())
		if err != nil {
			log.Panic(err)
		}
		fn(w2, r)
	}
}
