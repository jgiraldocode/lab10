package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jgiraldo/focustrack/internal/app"
	"github.com/jgiraldo/focustrack/internal/config"
	"github.com/jgiraldo/focustrack/internal/httpapi"
	"github.com/jgiraldo/focustrack/internal/httpapi/handlers"
	"github.com/jgiraldo/focustrack/internal/repository/sqlite"
	"github.com/jgiraldo/focustrack/internal/tracker"
)

func main() {
	cfg := config.Load()

	db, err := sqlite.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Repositories
	eventRepo := sqlite.NewSampleEventRepo(db)
	pomodoroRepo := sqlite.NewPomodoroRepo(db)

	// Services
	pomodoroSvc := app.NewPomodoroService(pomodoroRepo, time.Duration(cfg.CooldownSecs)*time.Second)
	segmentSvc := app.NewSegmentService(eventRepo)
	exportSvc := app.NewExportService(segmentSvc)

	macTracker := tracker.NewMacTracker()
	trackerSvc := app.NewTrackerService(macTracker, eventRepo, pomodoroSvc, cfg.PollInterval)

	// HTTP handlers
	trackingH := handlers.NewTrackingHandler(trackerSvc, segmentSvc, exportSvc)
	pomodoroH := handlers.NewPomodoroHandler(pomodoroSvc)
	allowlistH := handlers.NewAllowlistHandler(pomodoroSvc)

	router := httpapi.NewRouter(trackingH, pomodoroH, allowlistH)

	// Start tracker
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go trackerSvc.Run(ctx)

	// Start HTTP server
	addr := fmt.Sprintf(":%d", cfg.HTTPPort)
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("FocusTrack server listening on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	server.Shutdown(shutdownCtx)

	log.Println("Goodbye!")
}
