package web

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/rate-limiter/config"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/rate-limiter/internal/middleware"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/rate-limiter/internal/ratelimiter"
)

func StartServer(rl *ratelimiter.RateLimiter) {
	mux := http.NewServeMux()
	mux.HandleFunc("/docs/", docsHandler)

	mux.Handle("/", middleware.RateLimiterMiddleware(rl)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`
			<html>
			<head>
				<title>Redirecionando...</title>
				<meta http-equiv="refresh" content="2;url=/docs">
				<style>
					 * {
						color: #374151;
						font-size: 1.5rem;
					}
					 h2 {
       						color: #2563eb;
        					font-size: 2rem;
        					margin: 2rem 0 1rem;
      					}
				</style>
			</head>
			<body>
				<h2>Welcome to ratelimiter with Redis</h2>
				<p>You will be redirected to <a href="/docs">/docs</a> in 2 seconds.</p>
			</body>
			</html>
		`))
	})))

	mux.HandleFunc("/list-keys", func(w http.ResponseWriter, r *http.Request) {
		keys, err := rl.Store.ListKeys("*")
		if err != nil {
			http.Error(w, fmt.Sprintf("Erro ao listar chaves: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("["))
		for i, key := range keys {
			if i > 0 {
				w.Write([]byte(","))
			}
			w.Write([]byte(fmt.Sprintf(`"%s"`, key)))
		}
		w.Write([]byte("]"))
	})

	port := config.LoadConfig().AppPort
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func docsHandler(w http.ResponseWriter, r *http.Request) {
	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return
	}

	templatePath := filepath.Join(rootDir, "docs", "docs.html")

	t, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatal(err)
		return
	}

	if err := t.Execute(w, nil); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
