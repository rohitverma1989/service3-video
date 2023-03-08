package checkgrp

import (
	"encoding/json"
	"net/http"
	"os"

	"go.uber.org/zap"
)

type Handlers struct {
	Build string
	Log   *zap.SugaredLogger
}

func (h Handlers) Readiness(w http.ResponseWriter, r *http.Request) {
	statusCode := http.StatusOK
	data := struct {
		Build  string
		Status string `json:"status"`
	}{
		Status: "OK",
		Build:  h.Build,
	}

	if err := response(w, http.StatusOK, data); err != nil {
		h.Log.Errorw("readiness", "ERROR", err)
	}
	h.Log.Infow("readiness", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteAddr", r.RemoteAddr)
}

func (h Handlers) Liveness(w http.ResponseWriter, r *http.Request) {
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	data := struct {
		Status    string `json:"status,omitempty"`
		Build     string `json:"build,omitempty"`
		Host      string `json:"host,omitempty"`
		Name      string `json:"name,omitempty"`
		PodIP     string `json:"podIP,omitempty"`
		Node      string `json:"node,omitempty"`
		Namespace string `json:"namespace,omitempty"`
	}{
		Status:    "up",
		Build:     h.Build,
		Host:      host,
		Name:      os.Getenv("KUBERNETES_NAME"),
		PodIP:     os.Getenv("KUBERNETES_POD_IP"),
		Node:      os.Getenv("KUBERNETES_NODE_NAME"),
		Namespace: os.Getenv("KUBERNETES_NAMESPACE"),
	}

	statusCode := http.StatusOK
	if err := response(w, statusCode, data); err != nil {
		h.Log.Errorw("liveness", "ERROR", err)
	}

	// THIS IS A FREE TIMER. WE COULD UPDATE THE METRIC GOROUTINE COUNT HERE.
	h.Log.Infow("liveness", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)
}

func response(w http.ResponseWriter, statusCode int, data interface{}) error {
	// convert the response value to json
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	//set the content type and header once the marshelling is confirmed
	w.Header().Set("Content-Type", "application/json")

	// write the status code to response
	w.WriteHeader(statusCode)

	// Send the result back to client
	// json.NewEncoder(w).Encode(status)
	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}
