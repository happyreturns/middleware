package middleware

import (
	"bytes"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
)

func copyBody(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return body, nil
}

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
		body, err := copyBody(r)
		if err != nil {
			log.Println(err)
		}
		handler(w2, r, p)

		// write out to original writer
		w.WriteHeader(w2.Code)
		for k, values := range w2.HeaderMap {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
		_, err = w.Write(w2.Body.Bytes())
		if err != nil {
			log.Panic(err)
		}

		// Reset the body because it's empty after being read
		if body != nil {
			r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}

		fn(w2, r)
	}
}
