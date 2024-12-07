package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

const ServerPort = ":3000"
var ServicesURL string


type Registry struct {
	services map[string]ServiceConfig
	mutex    *sync.RWMutex
}

func (r *Registry) add(reg ServiceConfig) error {
	r.mutex.Lock()
	r.services[string(reg.Name)] = reg
	r.mutex.Unlock()
	err := r.sendRequiredServices(reg)
	r.notify(patch{
		Added: []patchEntry{
			{Name: reg.Name, URL: reg.URL},
		},
	})
	return err
}

func (r *Registry) get(name string) (ServiceConfig, bool) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()
    service, exists := r.services[name]
    return service, exists
}


func (r *Registry) remove(url string) error {
	for k, v := range r.services {
		if v.URL == url {
			r.notify(patch{
				Removed: []patchEntry{
					{Name: ServiceName(k), URL: v.URL},
				},
			})
			r.mutex.Lock()
			delete(r.services, k)
			r.mutex.Unlock()
			return nil
		}
	}
	return fmt.Errorf("service at URL %v not found", url)
}

func (r Registry) notify(p patch) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for _, reg := range r.services {
		go func(reg ServiceConfig) {
			for _, reqService := range reg.RequiredServices {
				p := patch{Added: []patchEntry{}, Removed: []patchEntry{}}
				sendUpdate := false
				for _, added := range p.Added {
					if added.Name == reqService {
						p.Added = append(p.Added, added)
						sendUpdate = true
					}
				}
				for _, removed := range p.Removed {
					if removed.Name == reqService {
						p.Removed = append(p.Removed, removed)
						sendUpdate = true
					}
				}
				if sendUpdate {
					err := r.sendPatch(p, reg.UpdateURL)
					if err != nil {
						log.Println(err)
						return
					}
				}
			}
		}(reg)
	}
}

func (r Registry) sendPatch(p patch, url string) error {
	d, err := json.Marshal(p)
	if err != nil {
		log.Println("Marshall error", err)
		return err
	}
	_, err = http.Post(url, "application/json", bytes.NewBuffer(d))
	if err != nil {
		log.Println("Post error", err)
		return err
	}
	return nil
}

func (r Registry) sendRequiredServices(reg ServiceConfig) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var p patch
	for _, serviceReg := range r.services {
		for _, reqService := range reg.RequiredServices {
			if serviceReg.Name == reqService {
				p.Added = append(p.Added, patchEntry{
					Name: serviceReg.Name,
					URL:  serviceReg.URL,
				})
			}
		}
	}
	err := r.sendPatch(p, reg.UpdateURL)
	if err != nil {
		log.Println("sendPatch error", err)
		return err
	}
	return nil
}

func (r *Registry) heartbeat(freq time.Duration) {
	for {
		var wg sync.WaitGroup
		for _, reg := range r.services {
			wg.Add(1)
			go func(reg ServiceConfig) {
				defer wg.Done()
				success := true
				for attempts := 0; attempts < 3; attempts++ {
					res, err := http.Get(reg.HeartbeatURL)
					if err != nil {
						log.Println(err)
					} else if res.StatusCode == http.StatusOK {
						log.Printf("Heartbeat check passed for %v", reg.Name)
						if !success {
							r.add(reg)
						}
						break
					}
					log.Printf("Heartbeat check failed for %v", reg.Name)
					if success {
						success = false
						r.remove(reg.URL)
					}
					time.Sleep(3 * time.Second) // wait to try again
				}
			}(reg)
		}
		wg.Wait()
		time.Sleep(freq)
	}
}

var regi Registry = *NewRegistry()

var once sync.Once

func SetupRegistryService() {
	once.Do(func() {
		go regi.heartbeat(3 * time.Second)
	})
}

type RegistryService struct{}

func (s RegistryService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("server ServeHTTP Request received req %v\n", req)
	switch req.Method {
	case http.MethodPost:
		dec := json.NewDecoder(req.Body)
		var r ServiceConfig
		err := dec.Decode(&r)
		if err != nil {
			log.Println("server Post ServeHTTP Decode error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("server ServeHTTP Adding service: %v with URL: %v\n", r.Name, r.URL)
		err = regi.add(r)
		if err != nil {
			log.Println("server ServeHTTP add error", err)
			// w.WriteHeader(http.StatusBadRequest)  this could fail in kubenetes since the other server is not fully up yet
			return

		}
	case http.MethodGet:
		serviceName := req.URL.Query().Get("name")
		log.Printf("server ServeHTTP Getting service: %v \n", serviceName)
		 // Get the service configuration
		 service, exists := regi.get(string(serviceName))
		 if !exists {
			log.Println("Service not found", serviceName)
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		 }
 
		 // Encode the service configuration as JSON and write it to the response
		 w.Header().Set("Content-Type", "application/json")
		 if err := json.NewEncoder(w).Encode(service); err != nil {
			log.Println("Failed to encode service", err)
			http.Error(w, "Failed to encode service", http.StatusInternalServerError)
		 }

	case http.MethodDelete:
		payload, err := io.ReadAll(req.Body)
		if err != nil {
			log.Println("server ServeHTTP ReadAll error", err)
			
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		url := string(payload)
		log.Printf("Removing service at URL: %v", url)
		err = regi.remove(url)
		if err != nil {
			fmt.Println("server ServeHTTP remove error", err)
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
