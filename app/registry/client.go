package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

// Providers contains the URLs for providers that the service
// requires.
type providers struct {
	services map[ServiceName][]string
	mutex    *sync.RWMutex
}

var srvprividers = providers{
	services: make(map[ServiceName][]string),
	mutex:    new(sync.RWMutex),
}

func (p *providers) Update(pat patch) {
	log.Printf("srvprividers Update %v\n", pat)
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for _, patchEntry := range pat.Added {
		if _, ok := p.services[patchEntry.Name]; !ok {
			p.services[patchEntry.Name] = make([]string, 0)
		}
		p.services[patchEntry.Name] = append(p.services[patchEntry.Name], patchEntry.URL)
	}
	for _, patchEntry := range pat.Removed {
		if providerURLs, ok := p.services[patchEntry.Name]; ok {
			for i := range providerURLs {
				if providerURLs[i] == patchEntry.URL {
					p.services[patchEntry.Name] = append(providerURLs[:i], providerURLs[i+1:]...)
				}
			}
		}
	}
}

func (p providers) get(name ServiceName) (string, error) {
	providers, ok := p.services[name]
	if !ok {
		log.Printf("registry client get  no providers available for service %v %v\n", name, srvprividers)
		return "", fmt.Errorf("no providers available for service %v", name)
	}
	idx := int(rand.Float32() * float32(len(providers)))
	log.Printf("registry client get providers %v idx %d provider %s\n", providers, idx, providers[idx])
		
	return providers[idx], nil
}

func GetProvider(name ServiceName) (string, error) {
	surl, err := srvprividers.get(name)
	fmt.Printf("registry client GetProvider %s url %s %v\n", name, surl, srvprividers)
	
	if err != nil {
		
    	serviceURL := fmt.Sprintf("%s?name=%s", ServicesURL, url.QueryEscape(string(name)))

		fmt.Printf("Client GetProvider, serviceURL: %s %s\n", ServicesURL, serviceURL)

		resp, err := http.Get(serviceURL)
		if err != nil {
			fmt.Println("Failed to send GET request: ", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Failed to get service configuration: %v\n", resp.Status)
		}

		var service ServiceConfig
		if err := json.NewDecoder(resp.Body).Decode(&service); err != nil {
			fmt.Printf("Failed to decode response: %v\n", err)
		}

		fmt.Printf("registry client GetProvider get Service Configuration: %s %+v \n", service.URL, service)
		surl = service.URL
	}
		
	return surl, nil
}

func RegisterService(r ServiceConfig) error {

	heartbeatURL, err := url.Parse(r.HeartbeatURL)
	if err != nil {
		log.Printf("RegisterService Parse HeartbeatURL for service %v \n", r.HeartbeatURL)
		return fmt.Errorf("failed to parse HeartbeatURL URL: %s %w", heartbeatURL, err)
	}
	http.HandleFunc(heartbeatURL.Path, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	serviceUpdateURL, err := url.Parse(r.UpdateURL)
	if err != nil {
		log.Printf("RegisterService Parse UpdateURL for service %v \n", r.HeartbeatURL)
		return fmt.Errorf("failed to parse UpdateURL URL: %s %w", serviceUpdateURL, err)
	}
	http.Handle(serviceUpdateURL.Path, &serviceUpdateHandler{})

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err = enc.Encode(r)
	if err != nil {
		log.Printf("RegisterService failed to Encode %v \n", err)
		return fmt.Errorf("failed to Encode  %w", err)
	}
	res, err := http.Post(ServicesURL, "application/json", buf)
	if err != nil {
		log.Printf("RegisterService failed to Post ServicesURL %s %v \n", ServicesURL, err)
		return fmt.Errorf("failed to Post ServicesURL %s  %w", ServicesURL, err)
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("RegisterService failed to register service ServicesURL %s %v \n", ServicesURL, err)
		return fmt.Errorf("failed to register service, ServicesURL %s. Registry service responded with code %v", ServicesURL, res.StatusCode)
	}
	log.Printf("RegisterService succeeds service ServicesURL %s %s %v \n", ServicesURL, r.Name, r)
		
	return nil
}

func RegisterServiceMux(r ServiceConfig) error {

	heartbeatURL, err := url.Parse(r.HeartbeatURL)
	if err != nil {
		log.Printf("RegisterServiceMux Parse HeartbeatURL for service %v \n", r.HeartbeatURL)
		return fmt.Errorf("failed to parse HeartbeatURL URL: %s %w", heartbeatURL, err)
	}
	r.Mux.HandleFunc(heartbeatURL.Path, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	serviceUpdateURL, err := url.Parse(r.UpdateURL)
	if err != nil {
		log.Printf("RegisterServiceMux Parse UpdateURL for service %v \n", r.HeartbeatURL)
		
		return fmt.Errorf("failed to parse UpdateURL URL: %s %w", serviceUpdateURL, err)
	}
	r.Mux.Handle(serviceUpdateURL.Path, &serviceUpdateHandler{})

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err = enc.Encode(r)
	if err != nil {
		log.Printf("RegisterServiceMux failed to Encode %v \n", err)
		return fmt.Errorf("failed to Encode  %w", err)
	}
	res, err := http.Post(ServicesURL, "application/json", buf)
	if err != nil {
		log.Printf("RegisterServiceMux failed to Post ServicesURL %s %v \n", ServicesURL, err)
		return fmt.Errorf("failed to Post ServicesURL %s  %w", ServicesURL, err)
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("RegisterServiceMux failed to register service ServicesURL %s %v \n", ServicesURL, err)
		return fmt.Errorf("failed to register service, ServicesURL %s. Registry service responded with code %v", ServicesURL, res.StatusCode)
	}
	log.Printf("RegisterServiceMux succeeds service ServicesURL %s %v \n", ServicesURL, r)
	
	return nil
}

func ShutdownService(serviceURL string) error {
	req, err := http.NewRequest(http.MethodDelete,
		ServicesURL,
		bytes.NewBuffer([]byte(serviceURL)))
	req.Header.Add("content-type", "text/plain")
	if err != nil {
		log.Printf("ShutdownService Add service %v \n", serviceURL)
		
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if res.StatusCode != http.StatusOK {
		log.Printf("ShutdownService DefaultClient.Do service %v \n", serviceURL)
		return fmt.Errorf("failed to deregister service. Registry service responded with code %v", res.StatusCode)
	}
	return err
}

type serviceUpdateHandler struct{}

func (suh serviceUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("ServeHTTP method is not MethodPost %v \n", r)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	dec := json.NewDecoder(r.Body)
	var p patch
	err := dec.Decode(&p)
	if err != nil {
		log.Println(err)
		fmt.Printf("ServeHTTP Decode err %v \n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("Client ServeHTTP srvprividers Update r %v p %v \n", r, p)
	srvprividers.Update(p)
}
