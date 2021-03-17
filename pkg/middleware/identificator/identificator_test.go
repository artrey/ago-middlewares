package identificator

import (
	"bytes"
	"github.com/go-chi/chi"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIdentificatorHTTPMux(t *testing.T) {
	mux := http.NewServeMux()
	identificatorMd := Identificator
	mux.Handle("/get",
		identificatorMd(http.HandlerFunc(
			func(writer http.ResponseWriter, request *http.Request) {
				identifier, err := Identifier(request.Context())
				if err != nil {
					if err == ErrNoIdentifier {
						writer.WriteHeader(http.StatusUnauthorized)
						return
					}
					t.Fatal(err)
				}
				_, err = writer.Write([]byte(*identifier))
				if err != nil {
					t.Fatal(err)
				}
			})),
	)

	type args struct {
		method string
		path   string
		addr   string
	}

	tests := []struct {
		name     string
		args     args
		wantCode int
		wantBody []byte
	}{
		{name: "GET", args: args{method: "GET", path: "/get", addr: "192.0.2.1:12345"}, wantCode: 200, wantBody: []byte("192.0.2.1")},
		{name: "GET2", args: args{method: "GET", path: "/get", addr: "127.0.0.1:666"}, wantCode: 200, wantBody: []byte("127.0.0.1")},
		{name: "GET3", args: args{method: "GET", path: "/get", addr: "127.0.0.1"}, wantCode: 401, wantBody: []byte{}},
		{name: "POST", args: args{method: "POST", path: "/get", addr: "192.0.2.1:12345"}, wantCode: 200, wantBody: []byte("192.0.2.1")},
		{name: "POST2", args: args{method: "POST", path: "/post", addr: "192.0.2.1:12345"}, wantCode: 404, wantBody: []byte("404 page not found\n")},
	}

	for _, tt := range tests {
		request := httptest.NewRequest(tt.args.method, tt.args.path, nil)
		request.RemoteAddr = tt.args.addr

		response := httptest.NewRecorder()
		mux.ServeHTTP(response, request)
		gotCode := response.Code
		if tt.wantCode != gotCode {
			t.Errorf("%s: got %d, wantCode %d", tt.name, gotCode, tt.wantCode)
		}
		gotBytes := response.Body.Bytes()
		if !bytes.Equal(tt.wantBody, gotBytes) {
			t.Errorf("%s: got %s, wantBody %s", tt.name, gotBytes, tt.wantBody)
		}
	}
}

func TestIdentificatorChi(t *testing.T) {
	router := chi.NewRouter()
	identificatorMd := Identificator
	router.With(identificatorMd).Get(
		"/get",
		func(writer http.ResponseWriter, request *http.Request) {
			identifier, err := Identifier(request.Context())
			if err != nil {
				if err == ErrNoIdentifier {
					writer.WriteHeader(http.StatusUnauthorized)
					return
				}
				t.Fatal(err)
			}
			_, err = writer.Write([]byte(*identifier))
			if err != nil {
				t.Fatal(err)
			}
		},
	)

	type args struct {
		method string
		path   string
		addr   string
	}

	tests := []struct {
		name     string
		args     args
		wantCode int
		wantBody []byte
	}{
		{name: "GET", args: args{method: "GET", path: "/get", addr: "192.0.2.1:12345"}, wantCode: 200, wantBody: []byte("192.0.2.1")},
		{name: "GET2", args: args{method: "GET", path: "/get", addr: "127.0.0.1:666"}, wantCode: 200, wantBody: []byte("127.0.0.1")},
		{name: "GET3", args: args{method: "GET", path: "/get", addr: "127.0.0.1"}, wantCode: 401, wantBody: []byte{}},
		{name: "POST", args: args{method: "POST", path: "/get", addr: "192.0.2.1:12345"}, wantCode: 405, wantBody: []byte{}},
		{name: "POST2", args: args{method: "POST", path: "/post", addr: "192.0.2.1:12345"}, wantCode: 404, wantBody: []byte("404 page not found\n")},
	}

	for _, tt := range tests {
		request := httptest.NewRequest(tt.args.method, tt.args.path, nil)
		request.RemoteAddr = tt.args.addr

		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		gotCode := response.Code
		if tt.wantCode != gotCode {
			t.Errorf("%s: got %d, wantCode %d", tt.name, gotCode, tt.wantCode)
		}
		gotBytes := response.Body.Bytes()
		if !bytes.Equal(tt.wantBody, gotBytes) {
			t.Errorf("%s: got %s, wantBody %s", tt.name, gotBytes, tt.wantBody)
		}
	}
}
