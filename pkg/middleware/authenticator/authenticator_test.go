package authenticator

import (
	"bytes"
	"context"
	"github.com/go-chi/chi"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticatorHTTPMux(t *testing.T) {
	mux := http.NewServeMux()
	authenticatorMd := Authenticator(func(ctx context.Context) (*string, error) {
		id := "192.0.2.1"
		return &id, nil
	}, func(ctx context.Context, id *string) (interface{}, error) {
		return "USERAUTH", nil
	})
	mux.Handle(
		"/get",
		authenticatorMd(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			profile, err := Authentication(request.Context())
			if err != nil {
				if err == ErrNoAuthentication {
					writer.WriteHeader(http.StatusUnauthorized)
					return
				}
				t.Fatal(err)
			}
			data := profile.(string)

			if data == "USERAUTH" && request.RemoteAddr != "192.0.2.1" {
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}

			_, err = writer.Write([]byte(data))
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
		{name: "GET", args: args{method: "GET", path: "/get", addr: "192.0.2.1"}, wantCode: 200, wantBody: []byte("USERAUTH")},
		{name: "POST", args: args{method: "POST", path: "/get", addr: "192.0.2.1"}, wantCode: 200, wantBody: []byte("USERAUTH")},
		{name: "Not found", args: args{method: "POST", path: "/post", addr: "192.0.2.1"}, wantCode: 404, wantBody: []byte("404 page not found\n")},
		{name: "No access", args: args{method: "POST", path: "/get", addr: "127.0.0.1"}, wantCode: 401, wantBody: []byte{}},
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
			t.Errorf("%s: got %s, want %s", tt.name, gotBytes, tt.wantBody)
		}
	}
}

func TestAuthenticatorChi(t *testing.T) {
	router := chi.NewRouter()
	authenticatorMd := Authenticator(func(ctx context.Context) (*string, error) {
		id := "192.0.2.1"
		return &id, nil
	}, func(ctx context.Context, id *string) (interface{}, error) {
		return "USERAUTH", nil
	})
	router.With(authenticatorMd).Get(
		"/get",
		func(writer http.ResponseWriter, request *http.Request) {
			profile, err := Authentication(request.Context())
			if err != nil {
				if err == ErrNoAuthentication {
					writer.WriteHeader(http.StatusUnauthorized)
					return
				}
				t.Fatal(err)
			}
			data := profile.(string)

			if data == "USERAUTH" && request.RemoteAddr != "192.0.2.1" {
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}

			_, err = writer.Write([]byte(data))
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
		{name: "GET", args: args{method: "GET", path: "/get", addr: "192.0.2.1"}, wantCode: 200, wantBody: []byte("USERAUTH")},
		{name: "POST", args: args{method: "POST", path: "/get", addr: "192.0.2.1"}, wantCode: 405, wantBody: []byte{}},
		{name: "Not found", args: args{method: "POST", path: "/post", addr: "192.0.2.1"}, wantCode: 404, wantBody: []byte("404 page not found\n")},
		{name: "No access", args: args{method: "GET", path: "/get", addr: "127.0.0.1"}, wantCode: 401, wantBody: []byte{}},
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
			t.Errorf("%s: got %s, want %s", tt.name, gotBytes, tt.wantBody)
		}
	}
}
