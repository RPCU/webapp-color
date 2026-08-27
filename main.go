package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

var availableColors = map[string]string{
	"red":    "#FF0000",
	"yellow": "#FFD700",
	"orange": "#FFA500",
	"lime":   "#00FF00",
	"green":  "#008000",
	"blue":   "#0000FF",
	"navy":   "#000080",
	"purple": "#800080",
	"pink":   "#FF00FF",
	"brown":  "#A52A2A",
	"grey":   "#808080",
	"black":  "#000000",
}

func checkEnvColor(c string) (string, error) {
	// Check if var is empty
	if _, exists := os.LookupEnv("APP_COLOR"); !exists {
		return "", fmt.Errorf("APP_COLOR env var should not be empty")
	}
	// Check if color is available
	if value, ok := availableColors[c]; ok {
		return value, nil
	}
	colorKeys := make([]string, 0, len(availableColors))
	for k := range availableColors {
		colorKeys = append(colorKeys, k)
	}
	return "", fmt.Errorf("color not supported: get=%s available=%v", c, colorKeys)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	// Get color code if value is set and exist
	cValue, err := checkEnvColor(os.Getenv("APP_COLOR"))
	if err != nil {
		slog.Error("invalid color configuration", "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Get hostname
	h, _ := os.Hostname()

	// Populate struct
	data := struct {
		Hostname string
		Color    string
	}{
		Hostname: h,
		Color:    cValue,
	}

	// Html template rendering
	t, err := template.ParseFiles("hello.html")
	if err != nil {
		slog.Error("failed to parse template", "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, data); err != nil {
		slog.Error("failed to render template", "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Structured JSON logging (slog emits uppercase levels: INFO, WARN, ERROR)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	// Configure application port if env var is set
	appPort, exists := os.LookupEnv("APP_PORT")
	if exists {
		appPort = ":" + appPort
	} else {
		appPort = ":8080"
	}

	// Check APP_COLOR env var and start webserver
	if _, err := checkEnvColor(os.Getenv("APP_COLOR")); err != nil {
		slog.Error("startup configuration error", "err", err)
		os.Exit(1)
	}

	slog.Info("starting webapp-color", "port", appPort)
	http.HandleFunc("/", viewHandler)
	if err := http.ListenAndServe(appPort, nil); err != nil {
		slog.Error("server stopped", "err", err)
		os.Exit(1)
	}
}
