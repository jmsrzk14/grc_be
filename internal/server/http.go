package server

import (
	_ "grc_be/api/docs"
	"grc_be/internal/conf"
	"grc_be/internal/service"
	"net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// NewHTTPServer create a HTTP server.
func NewHTTPServer(c *conf.Server, tenant *service.TenantService, property *service.PropertyService, reg *service.RegulationService, ass *service.AssessmentService, auth *service.AuthService, risk *service.RiskService, logger log.Logger) *khttp.Server {
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

	// Apply AuthMiddleware to all routes
	router.Use(AuthMiddleware)

	// Auth API
	aAuth := router.PathPrefix("/api/v1/auth").Subrouter()
	aAuth.HandleFunc("/login", auth.Login).Methods("POST")
	aAuth.HandleFunc("/register", auth.Register).Methods("POST")

	// Tenant API
	t := router.PathPrefix("/api/v1/tenants").Subrouter()
	t.HandleFunc("", tenant.ListTenants).Methods("GET")
	t.HandleFunc("", AdminOnly(tenant.CreateTenant)).Methods("POST")
	t.HandleFunc("/{id}", tenant.GetTenant).Methods("GET")
	t.HandleFunc("/{id}", AdminOnly(tenant.UpdateTenant)).Methods("PUT")
	t.HandleFunc("/{id}", AdminOnly(tenant.DeleteTenant)).Methods("DELETE")
	t.HandleFunc("/{id}/properties", tenant.ListTenantProperties).Methods("GET")

	// Property API
	p := router.PathPrefix("/api/v1/properties").Subrouter()
	p.HandleFunc("", property.ListProperties).Methods("GET")
	p.HandleFunc("", AdminOnly(property.CreateProperty)).Methods("POST")
	p.HandleFunc("/{id}", property.GetProperty).Methods("GET")
	p.HandleFunc("/{id}", AdminOnly(property.UpdateProperty)).Methods("PUT")
	p.HandleFunc("/{id}", AdminOnly(property.DeleteProperty)).Methods("DELETE")

	// Regulation API
	regR := router.PathPrefix("/api/v1/regulations").Subrouter()
	regR.HandleFunc("", reg.ListRegulations).Methods("GET")
	regR.HandleFunc("", reg.CreateRegulation).Methods("POST")
	regR.HandleFunc("/upsert", AdminOnly(reg.UpsertRegulation)).Methods("POST")
	regR.HandleFunc("/{id}", reg.GetRegulation).Methods("GET")
	regR.HandleFunc("/{id}", reg.UpdateRegulation).Methods("PUT")
	regR.HandleFunc("/{id}", reg.DeleteRegulation).Methods("DELETE")
	regR.HandleFunc("/{id}/items", reg.ListItems).Methods("GET")
	regR.HandleFunc("/{id}/items", AdminOnly(reg.CreateItem)).Methods("POST")
	regR.HandleFunc("/{id}/items/upsert", AdminOnly(reg.UpsertItem)).Methods("POST")
	regR.HandleFunc("/{id}/items/{item_id}", reg.GetItem).Methods("GET")
	regR.HandleFunc("/{id}/items/{item_id}", AdminOnly(reg.UpdateItem)).Methods("PUT")
	regR.HandleFunc("/{id}/items/{item_id}", AdminOnly(reg.DeleteItem)).Methods("DELETE")
	regR.HandleFunc("/{id}/mappings", reg.ListMappings).Methods("GET")
	regR.HandleFunc("/{id}/mappings", AdminOnly(reg.AddMapping)).Methods("POST")
	regR.HandleFunc("/{id}/mappings/{mapping_id}", AdminOnly(reg.DeleteMapping)).Methods("DELETE")

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
	a.HandleFunc("/sessions/{id}/sync", ass.SyncSession).Methods("POST")

	// Risk Management API
	rc := router.PathPrefix("/api/v1/risk-categories").Subrouter()
	rc.HandleFunc("", risk.ListCategories).Methods("GET")
	rc.HandleFunc("", AdminOnly(risk.CreateCategory)).Methods("POST")
	rc.HandleFunc("/{id}", risk.GetCategory).Methods("GET")
	rc.HandleFunc("/{id}", AdminOnly(risk.UpdateCategory)).Methods("PUT")
	rc.HandleFunc("/{id}", AdminOnly(risk.DeleteCategory)).Methods("DELETE")
	rc.HandleFunc("/{id}/settings", risk.GetCategoryTenant).Methods("GET")
	rc.HandleFunc("/{id}/settings", risk.SaveCategoryTenant).Methods("POST")

	ri := router.PathPrefix("/api/v1/risks").Subrouter()
	ri.HandleFunc("", risk.ListRisks).Methods("GET")
	ri.HandleFunc("", risk.CreateRisk).Methods("POST")
	ri.HandleFunc("/{id}", risk.GetRisk).Methods("GET")
	ri.HandleFunc("/{id}", risk.UpdateRisk).Methods("PUT")
	ri.HandleFunc("/{id}", risk.DeleteRisk).Methods("DELETE")

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
