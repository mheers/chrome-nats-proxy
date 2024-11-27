package proxy

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/elazarl/goproxy"
	"github.com/mheers/chrome-nats-proxy/mqclient/models"
)

type ProxyRequest struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    interface{}       `json:"body"`
}

type ProxyResponse struct {
	Status     int               `json:"status"`
	StatusText string            `json:"statusText"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

type Proxy struct {
	mqClient *models.MQClient
}

func NewProxy(mqClient *models.MQClient) *Proxy {
	return &Proxy{
		mqClient: mqClient,
	}
}

func (p *Proxy) Start() {
	setCA(caCert, caKey)
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		response := goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusForbidden, "Blocked")

		p.handleRequest(response, req)

		return req, response
	})

	proxy.Verbose = true
	log.Fatal(http.ListenAndServe(":8080", proxy))
}

func (p *Proxy) handleRequest(w *http.Response, r *http.Request) {
	// Log incoming request
	log.Printf("Received request: %s %s\n", r.Method, r.URL.String())

	// Prepare the ProxyRequest object
	proxyReq := ProxyRequest{
		URL:     r.URL.String(),
		Method:  r.Method,
		Headers: extractHeaders(r.Header),
	}

	// Read and set body if the method supports it
	if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			httpError(w, "Failed to read request body", http.StatusInternalServerError)
			log.Printf("Error reading request body: %v\n", err)
			return
		}
		r.Body.Close()
		if len(bodyBytes) > 0 {
			var bodyJSON interface{}
			if err := json.Unmarshal(bodyBytes, &bodyJSON); err != nil {
				bodyJSON = string(bodyBytes) // Use raw string if JSON unmarshalling fails
			}
			proxyReq.Body = bodyJSON
		}
	}

	// Convert ProxyRequest to JSON
	reqJSON, err := json.Marshal(proxyReq)
	if err != nil {
		httpError(w, "Failed to serialize request", http.StatusInternalServerError)
		log.Printf("Error marshaling request: %v\n", err)
		return
	}

	// Publish the request to the message queue
	resp, err := p.mqClient.Connection.Request("instance.chrome-nats-proxy.http", reqJSON, 30*time.Second)
	if err != nil {
		httpError(w, fmt.Sprintf("Failed to send request to backend: %v", err), http.StatusInternalServerError)
		log.Printf("Error sending request to backend: %v\n", err)
		return
	}

	// Unmarshal the response from the message queue
	respBody := resp.Data

	var proxyResp ProxyResponse
	if err := json.Unmarshal(respBody, &proxyResp); err != nil {
		httpError(w, "Invalid response from backend", http.StatusInternalServerError)
		log.Printf("Error unmarshaling backend response: %v\n", err)
		return
	}

	httpResp := &http.Response{
		StatusCode: proxyResp.Status,
		Status:     proxyResp.StatusText,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(proxyResp.Body)),
	}

	// Copy headers
	for key, values := range resp.Header {
		for _, value := range values {
			httpResp.Header.Add(key, value)
		}
	}

	copyResponse(w, httpResp)
}

func httpError(w *http.Response, message string, status int) {
	w.Header.Set("Content-Type", "text/plain")
	w.Body = io.NopCloser(strings.NewReader(message))
	w.StatusCode = status
	w.Status = http.StatusText(status)
}

func copyResponse(w *http.Response, resp *http.Response) {
	// Copy status code
	w.StatusCode = resp.StatusCode

	// Copy headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header.Add(key, value)
		}
	}

	// Copy body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		httpError(w, "Failed to read response body", http.StatusInternalServerError)
		log.Printf("Error reading response body: %v\n", err)
		return
	}
	w.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
}

// Helper function to convert HTTP headers to a map
func extractHeaders(header http.Header) map[string]string {
	headers := make(map[string]string)
	for key, values := range header {
		headers[key] = strings.Join(values, ", ")
	}
	return headers
}
