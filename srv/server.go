package srv

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/ctrlaltdel121/cert-server/cert"
	"github.com/ctrlaltdel121/cert-server/cert/store"
	"github.com/gorilla/mux"
)

// Server is the main struct for the server, holding business logic and a storage mechanism
type Server struct {
	Store store.CertStorer
}

// NewServer creates a new server with the given storage
func NewServer(storeType string) *Server {
	switch storeType {
	case "s3":
		return &Server{Store: store.NewS3Store()}
	default:
		return &Server{Store: store.NewFileStore()}
	}
}

// Serve starts the server
func (s *Server) Serve() error {
	// blocks until finished
	return http.ListenAndServe(":"+os.Getenv("PORT"), s.createRouter())
}

// this is a separate function so we can expose the router to tests
func (s *Server) createRouter() *mux.Router {
	router := mux.NewRouter()
	router.Methods("POST").Path("/certificates").HandlerFunc(errorHandler(s.createCert))
	router.Methods("GET").Path("/certificates/{serial}").HandlerFunc(errorHandler(s.getCert))
	router.Methods("DELETE").Path("/certificates/{serial}").HandlerFunc(errorHandler(s.deleteCert))
	return router
}

// type definition of a function that handles incoming requests
type myHandlerFunc func(w http.ResponseWriter, r *http.Request) error

// writeResponse takes a code and response struct and JSON marshals and writes it
func writeResponse(w http.ResponseWriter, code int, body interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
	return nil
}

func (s *Server) createCert(w http.ResponseWriter, r *http.Request) error {
	if r.Body == nil {
		return BadRequestError("Request must have a body")
	}
	var c cert.CertInput
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		return BadRequestError("Body is not valid JSON")
	}

	if len(c.Names) == 0 {
		return BadRequestError("Request must contain at least one name for certificate")
	}

	cert, err := c.Generate()
	if err != nil {
		return InternalError("Could not generate cert: " + err.Error())
	}

	err = s.Store.Write(cert)
	if err != nil {
		return InternalError("Could not write cert: " + err.Error())
	}

	return writeResponse(w, 201, cert)
}

func (s *Server) getCert(w http.ResponseWriter, r *http.Request) error {
	serial, err := strconv.Atoi(mux.Vars(r)["serial"])
	if err != nil {
		return BadRequestError("Serial must be an integer")
	}

	cert, err := s.Store.Read(int64(serial))
	if err != nil {
		if os.IsNotExist(err) {
			return NotFoundError("Certificate not found")
		}
		return InternalError("Could not read certificate")
	}

	return writeResponse(w, 200, cert)
}

func (s *Server) deleteCert(w http.ResponseWriter, r *http.Request) error {
	serial, err := strconv.Atoi(mux.Vars(r)["serial"])
	if err != nil {
		return BadRequestError("Serial must be an integer")
	}

	err = s.Store.Delete(int64(serial))
	if err != nil {
		if os.IsNotExist(err) {
			return NotFoundError("Certificate not found")
		}
		return InternalError("Could not read certificate")
	}
	w.WriteHeader(202)
	return nil
}
