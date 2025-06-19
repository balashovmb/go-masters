package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/pprof"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"go-masters/final-project/internal/config"
	"go-masters/final-project/internal/db"
	"go-masters/final-project/internal/db/postgres"
	"go-masters/final-project/internal/llm"
	"go-masters/final-project/internal/metrics"
	"go-masters/final-project/internal/models"
	"go-masters/final-project/internal/rating"
	"go-masters/final-project/internal/telemetry"
)

type Server struct {
	cfg    *config.Cfg
	router *chi.Mux
	server *http.Server
	db     db.DB
	llm    *llm.LLM
}

type averageRatingResponse struct {
	Rating string `json:"rating"`
	Id     int    `json:"id"`
}

type listRequest struct {
	Filter string `json:"filter"`
	Id     int    `json:"id"`
}

func New(cfg *config.Cfg) (*Server, error) {
	r := chi.NewRouter()

	db, err := postgres.New(cfg.DBConnStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации БД: %w", err)
	}

	LLM := llm.New(cfg.LlmModel, cfg.LlmURL, cfg.LlmPort)

	s := Server{
		cfg:    cfg,
		router: r,
		server: &http.Server{
			Addr:         fmt.Sprintf(":%v", cfg.Port),
			Handler:      r,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
		db:  db,
		llm: LLM,
	}

	s.endpoints()

	return &s, nil
}

func (s *Server) endpoints() {
	// Настройка middleware
	s.router.Use(
		middleware.RequestID,                 // Добавляет X-Request-Id в заголовки
		telemetry.TracingMiddleware,          // OpenTelemetry трейсинг
		metrics.PrometheusMiddleware,         // Метрики Prometheus
		RequestLoggerMiddleware(&log.Logger), // Логирование запросов
		middleware.Recoverer,                 // Восстановление после паник
	)

	// Эндпоинты pprof
	// http://localhost:8080/debug/pprof/
	s.router.Get("/debug/pprof/", pprof.Index)
	s.router.Get("/debug/pprof/cmdline", pprof.Cmdline)
	s.router.Get("/debug/pprof/profile", pprof.Profile)
	s.router.Get("/debug/pprof/symbol", pprof.Symbol)
	s.router.Get("/debug/pprof/trace", pprof.Trace)
	s.router.Get("/debug/pprof/allocs", pprof.Handler("allocs").ServeHTTP)
	s.router.Get("/debug/pprof/block", pprof.Handler("block").ServeHTTP)
	s.router.Get("/debug/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
	s.router.Get("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)
	s.router.Get("/debug/pprof/mutex", pprof.Handler("mutex").ServeHTTP)
	s.router.Get("/debug/pprof/threadcreate", pprof.Handler("threadcreate").ServeHTTP)

	// HealthCheck - статус системы
	s.router.Get("/health", healthHandler)

	// Эндпоинт для Prometheus
	s.router.Get("/metrics", promhttp.Handler().ServeHTTP)

	// Инициализация маршрутов
	s.router.Post("/reviews", s.addReviewHandler)
	s.router.Get("/reviews", s.listReviewsHandler)
	s.router.Get("/reviews/{id}/average", s.averageRatingHandler)

}

func (s *Server) Start(ctx context.Context) error {
	log.Info().Msg("Инициализация телеметрии")

	shutdown, err := telemetry.SetupOTelSDK(ctx, "http://localhost:4318")
	if err != nil {
		return err
	}
	defer func() {
		err = shutdown(ctx)
		if err != nil {
			log.Err(err).Msg("cannot shutdown OTel")
		}
	}()

	log.Info().Str("addr", s.server.Addr).Msg("Запуск HTTP сервера")

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		log.Info().Msg("Остановка HTTP сервера")
		if err := s.server.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("Ошибка при остановке сервера")
		}
	}()

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Обработчики запросов

func healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)
	defer span.End()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *Server) addReviewHandler(w http.ResponseWriter, r *http.Request) {
	span := trace.SpanFromContext(r.Context())
	defer span.End()

	log.Info().Msg("Обработка запроса addReview")
	span.AddEvent("Обработка запроса addReview")

	var req models.Review

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		span.SetStatus(codes.Error, "не удалось декодировать запрос")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := s.db.AddReview(r.Context(), req)

	if err != nil {
		span.SetStatus(codes.Error, "не удалось добавить отзыв в БД")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go s.updateReviewRating(context.Background(), span, id, req.Text)
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) listReviewsHandler(w http.ResponseWriter, r *http.Request) {
	span := trace.SpanFromContext(r.Context())
	defer span.End()

	var req listRequest
	var filter string
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error().Err(err).Msg("не удалось декодировать запрос")
		span.SetStatus(codes.Error, "не удалось декодировать запрос")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	switch req.Filter {
	case "user":
		filter = "user_id"
	case "object":
		filter = "object_id"
	default:
		filter = ""
	}
	if filter == "" {
		log.Error().Msg("неверный фильтр")
		http.Error(w, "неверный фильтр", http.StatusBadRequest)
		span.SetStatus(codes.Error, "неверный фильтр")
		return
	}

	if req.Id == 0 {
		log.Error().Msg("неверный id")
		http.Error(w, "неверный id", http.StatusBadRequest)
		span.SetStatus(codes.Error, "неверный id")
		return
	}

	log.Info().Msg("Обработка запроса listReviews")
	span.AddEvent("Обработка запроса listReviews")

	reviews, err := s.db.ListReviews(ctx, filter, req.Id)
	if err != nil {
		log.Error().Err(err).Msg("не удалось получить отзывы")
		span.SetStatus(codes.Error, "не удалось получить отзывы")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	reviewRepresentations := make([]models.ReviewRepresentation, len(reviews))

	for i, review := range reviews {
		reviewRepresentations[i] = models.ReviewRepresentation{
			ID:       review.ID,
			UserID:   review.UserID,
			ObjectID: review.ObjectID,
			Text:     review.Text,
			Rating:   rating.RatingToString(review.Rating),
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviewRepresentations)
}

func (s *Server) averageRatingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)
	defer span.End()

	log.Info().Msg("Обработка запроса averageRating")
	span.AddEvent("Обработка запроса averageRating")

	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		span.SetStatus(codes.Error, "не удалось преобразовать id")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	avgRat, err := s.db.AverageRating(ctx, idInt)
	if err != nil {
		span.SetStatus(codes.Error, "не удалось получить среднюю оценку")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	res := averageRatingResponse{Rating: rating.AverageRatingToString(avgRat), Id: idInt}
	json.NewEncoder(w).Encode(res)
}

// RequestLoggerMiddleware - middleware для логирования запросов
func RequestLoggerMiddleware(logger *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			logger.Info().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_addr", r.RemoteAddr).
				Str("user_agent", r.UserAgent()).
				Int("status", ww.Status()).
				Int("bytes", ww.BytesWritten()).
				Dur("duration", time.Since(start)).
				Str("request_id", middleware.GetReqID(r.Context())).
				Msg("Обработан HTTP запрос")
		})
	}
}

func (s *Server) updateReviewRating(ctx context.Context, span trace.Span, id int, review string) {
	span.AddEvent("Обновление рейтинга")
	llmResponse, err := s.llm.GetRating(review)
	if err != nil {
		span.SetStatus(codes.Error, "не удалось обновить рейтинг")
		log.Error().Err(err).Msg("не удалось обновить рейтинг")
		return
	}
	s.db.UpdateReviewRating(ctx, id, llmResponse)
}
