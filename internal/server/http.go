package server

import (
	"net/http"
	"grc_be/internal/conf"
	"grc_be/internal/service"
	_ "grc_be/api/docs"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// NewHTTPServer create a HTTP server.
func NewHTTPServer(c *conf.Server, tenant *service.TenantService, property *service.PropertyService, reg *service.RegulationService, ass *service.AssessmentService, auth *service.AuthService, logger log.Logger) *khttp.Server {
	var opts = []khttp.ServerOption{
		khttp.Middleware(
			recovery.Recovery(),
		),
	}
	if c.HTTP.Addr != "" {
		opts = append(opts, khttp.Address(c.HTTP.Addr))
	}
	srv := khttp.NewServer(opts...)
	
	// Gunakan gorilla mux untuk routing REST yang fleksibel
	router := mux.NewRouter()

	// Auth API
	aAuth := router.PathPrefix("/api/v1/auth").Subrouter()
	aAuth.HandleFunc("/login", auth.Login).Methods("POST")
	aAuth.HandleFunc("/register", auth.Register).Methods("POST")
	
	// Tenant API
	t := router.PathPrefix("/api/v1/tenants").Subrouter()
	t.HandleFunc("", tenant.ListTenants).Methods("GET")
	t.HandleFunc("", tenant.CreateTenant).Methods("POST")
	t.HandleFunc("/{id}", tenant.GetTenant).Methods("GET")
	t.HandleFunc("/{id}", tenant.UpdateTenant).Methods("PUT")
	t.HandleFunc("/{id}", tenant.DeleteTenant).Methods("DELETE")
	t.HandleFunc("/{id}/properties", tenant.ListTenantProperties).Methods("GET")

	// Property API
	p := router.PathPrefix("/api/v1/properties").Subrouter()
	p.HandleFunc("", property.ListProperties).Methods("GET")
	p.HandleFunc("", property.CreateProperty).Methods("POST")
	p.HandleFunc("/{id}", property.GetProperty).Methods("GET")
	p.HandleFunc("/{id}", property.UpdateProperty).Methods("PUT")
	p.HandleFunc("/{id}", property.DeleteProperty).Methods("DELETE")

	// Regulation API
	r := router.PathPrefix("/api/v1/regulations").Subrouter()
	r.HandleFunc("", reg.ListRegulations).Methods("GET")
	r.HandleFunc("", reg.CreateRegulation).Methods("POST")
	r.HandleFunc("/{id}", reg.GetRegulation).Methods("GET")
	r.HandleFunc("/{id}", reg.UpdateRegulation).Methods("PUT")
	r.HandleFunc("/{id}", reg.DeleteRegulation).Methods("DELETE")
	r.HandleFunc("/{id}/items", reg.ListItems).Methods("GET")
	r.HandleFunc("/{id}/items", reg.CreateItem).Methods("POST")
	r.HandleFunc("/{id}/items/{item_id}", reg.GetItem).Methods("GET")
	r.HandleFunc("/{id}/items/{item_id}", reg.UpdateItem).Methods("PUT")
	r.HandleFunc("/{id}/items/{item_id}", reg.DeleteItem).Methods("DELETE")
	r.HandleFunc("/{id}/mappings", reg.ListMappings).Methods("GET")
	r.HandleFunc("/{id}/mappings", reg.AddMapping).Methods("POST")
	r.HandleFunc("/{id}/mappings/{mapping_id}", reg.DeleteMapping).Methods("DELETE")

	// Assessment API
	a := router.PathPrefix("/api/v1/assessments").Subrouter()
	a.HandleFunc("/sessions", ass.ListSessions).Methods("GET")
	a.HandleFunc("/sessions", ass.CreateSession).Methods("POST")
	a.HandleFunc("/sessions/{id}", ass.GetSession).Methods("GET")
	a.HandleFunc("/sessions/{id}", ass.UpdateSession).Methods("PUT")
	a.HandleFunc("/sessions/{id}", ass.DeleteSession).Methods("DELETE")
	a.HandleFunc("/sessions/{id}/results", ass.ListResults).Methods("GET")
	a.HandleFunc("/sessions/{id}/results", ass.SubmitResult).Methods("POST")
	a.HandleFunc("/sessions/{id}/results/{result_id}", ass.DeleteResult).Methods("DELETE")
	a.HandleFunc("/sessions/{id}/summaries", ass.GetSummaries).Methods("GET")

	// Swagger UI & OpenAPI Spec
	router.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./api/docs/swagger.json")
	})
	
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"), 
	))

	// CORS configuration
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Allow all Origins for Dev
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(router)

	srv.HandlePrefix("/", corsHandler)
	return srv
}
