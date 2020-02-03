package streamprocessor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Server struct holds and httprouter
type Server struct {
	Router *httprouter.Router
}

func (s *Server) ready(ctx context.Context) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var data struct {
			Status string
		}

		if ctx.Err() != nil {
			data.Status = "Server shutdown"
		} else {
			data.Status = "OK"
		}

		resp, err := json.Marshal(data)

		if err != nil {
			fmt.Println("JSON marshalling error: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(resp)
		if err != nil {
			fmt.Println("Error writing a response: ", err)
		}
	}
}
