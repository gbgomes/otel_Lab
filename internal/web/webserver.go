package webserver

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gbgomes/GoExpert/otel_Lab/internal/entity"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/propagation/b3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
)

type TemplateData struct {
	OTELTracer trace.Tracer
	ZpTracer   zipkin.Tracer
}

type WebServer struct {
	TemplateData *TemplateData
}

type Resposta struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type PostParam struct {
	Cep string `json:"cep"`
}

func NewServer(TemplateData *TemplateData) *WebServer {
	return &WebServer{
		TemplateData: TemplateData,
	}
}

func (we *WebServer) CreateServer() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(zipkinhttp.NewServerMiddleware(
		&we.TemplateData.ZpTracer,
		zipkinhttp.SpanName("otel_Lab")),
	)
	router.Handle("/metrics", promhttp.Handler())
	router.Post("/", we.HandleInicio)
	router.Get("/clima", we.HandleClima)

	return router
}

func (we *WebServer) HandleInicio(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var param PostParam
	err := decoder.Decode(&param)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("invalid zipcode"))
		return
	}

	_, err = strconv.Atoi(param.Cep)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("invalid zipcode"))
		return
	}

	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	//criação dos spans
	ctx, span := we.TemplateData.OTELTracer.Start(ctx, "Serviço A - Inicio")
	defer span.End()
	zipspan := we.TemplateData.ZpTracer.StartSpan("Serviço A - Inicio")
	defer zipspan.Finish()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/clima?cep="+param.Cep, nil)
	if err != nil {
		log.Fatal("erro no início, criando req")
	}

	//Injeta as propagações de span
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	injector := b3.InjectHTTP(req)
	injector(zipspan.Context())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("erro no início, executando req")
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bodyBytes)
}

func (we *WebServer) HandleClima(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	//criação dos spans
	ctx, span := we.TemplateData.OTELTracer.Start(ctx, "Serviço B - busca CEP")
	//zipspan := we.TemplateData.ZpTracer.StartSpan("Serviço B - Coleta Localidade")
	zipspan, _ := we.TemplateData.ZpTracer.StartSpanFromContext(ctx, "Serviço B - busca CEP")

	cep := r.URL.Query().Get("cep")

	local := entity.NewLocalidadeViaCep()
	localidade, err := local.ColetaLocalidade(cep)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(err.StatusCode)
		w.Write([]byte(err.Err.Error()))
		return
	}
	span.End()
	zipspan.Finish()

	//criação dos spans
	_, span2 := we.TemplateData.OTELTracer.Start(ctx, "Serviço B - busca temperatura")
	zipspan2 := we.TemplateData.ZpTracer.StartSpan("Serviço B - busca temperatura", zipkin.Parent(zipspan.Context()))
	defer span2.End()
	defer zipspan2.Finish()

	weather := entity.NewWeather()
	tempo, err := weather.ColetaTempo(localidade.Localidade)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(err.StatusCode)
		w.Write([]byte(err.Err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&Resposta{localidade.Localidade, tempo.Current.TempC, tempo.Current.TempF, tempo.Current.TempC + 273})
}
