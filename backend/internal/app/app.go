package app

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/noedaka/clothing-visual-search/backend/internal/config"
	"github.com/noedaka/clothing-visual-search/backend/internal/handler"
	milvusclient "github.com/noedaka/clothing-visual-search/backend/internal/milvus-client"
	minioclient "github.com/noedaka/clothing-visual-search/backend/internal/minio-client"
	mlclient "github.com/noedaka/clothing-visual-search/backend/internal/ml-client"
	"github.com/noedaka/clothing-visual-search/backend/internal/repository"
	"github.com/noedaka/clothing-visual-search/backend/internal/service"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Run() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
		return
	}

	db, err := initDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
		return
	}
	defer db.Close()

	minioClient, err := minioclient.NewMinIOClient(cfg)
	if err != nil {
		log.Fatalf("failed to initialize MinIO client: %v", err)
		return
	}

	if err = minioClient.EnsureMinIOBucket(); err != nil {
		log.Fatalf("failed to ensure MinIO bucket: %v", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	milvusClient, err := milvusclient.NewClient(ctx, cfg.MilvusAddr)
	if err != nil {
		log.Fatalf("failed to initialize milvus client: %v", err)
		return
	}

	mlClient, err := mlclient.NewClient(cfg.MLServiceAddr)
	if err != nil {
		log.Fatalf("failed to initialize ml client: %v", err)
		return
	}

	categoryRepo := repository.NewCategoryRepo(db)
	productRepo := repository.NewProductRepo(db)
	imageRepo := repository.NewImageRepo(db, minioClient.Client, milvusClient, cfg)

	categoryService := service.NewCategoryServ(categoryRepo)
	productService := service.NewProductServ(productRepo, imageRepo)
	searchService := service.NewSearchServ(milvusClient)
	embeddingService := service.NewEmbeddingServ(mlClient)

	handler := handler.NewHandler(
		productService,
		categoryService,
		embeddingService,
		searchService,
	)

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Route("/api", func(r chi.Router) {
			r.Route("/category", func(r chi.Router) {
				r.Post("/", handler.AddCategoryHandler)
				r.Get("/", handler.ListCategoryHandler)
			})

			r.Route("/product", func(r chi.Router) {
				r.Post("/", handler.AddProductHandler)
				r.Post("/seacrh", handler.SearchByImageHandler)
			})
		})
	})

	server := &http.Server{
		Addr:    cfg.ServerPort,
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	cancel()

	time.Sleep(1 * time.Second)
}

func initDatabase(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.DataBaseURL)
	if err != nil {
		return nil, err
	}

	return db, nil
}
