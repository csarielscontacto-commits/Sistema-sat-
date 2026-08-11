package main

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "strconv"
    "syscall"
    "time"

    "github.com/gorilla/mux"
    "github.com/segmentio/kafka-go"
    "github.com/sistema-sat/internal/models"
    "github.com/sistema-sat/internal/sat"
    "github.com/sistema-sat/pkg/logger"
)

type Config struct {
    SATEndpoint    string
    PACEndpoint    string
    KafkaBrokers   []string
    KafkaTopic     string
    WorkerCount    int
    TimeoutSeconds int
}

func loadConfig() *Config {
    return &Config{
        SATEndpoint:    getEnv("SAT_ENDPOINT", "https://cfdi.sat.gob.mx/ws/soap"),
        PACEndpoint:    getEnv("PAC_ENDPOINT", "https://api.pac.com.mx/v1/timbrar"),
        KafkaBrokers:   []string{getEnv("KAFKA_BROKER", "localhost:9092")},
        KafkaTopic:     getEnv("KAFKA_TOPIC", "incoming_cfdis"),
        WorkerCount:    getEnvInt("WORKER_COUNT", 10),
        TimeoutSeconds: getEnvInt("TIMEOUT_SECONDS", 30),
    }
}

func main() {
    log := logger.NewLogger("ingest-sat")
    cfg := loadConfig()
    satClient := sat.NewClient(cfg.SATEndpoint, cfg.TimeoutSeconds)

    kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
        Brokers:  cfg.KafkaBrokers,
        Topic:    cfg.KafkaTopic,
        Balancer: &kafka.LeastBytes{},
    })
    defer kafkaWriter.Close()

    router := mux.NewRouter()
    router.HandleFunc("/api/v1/ingest", ingestHandler(satClient, kafkaWriter, log)).Methods("POST")
    router.HandleFunc("/health", healthCheck).Methods("GET")

    server := &http.Server{
        Addr:         ":8080",
        Handler:      router,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 30 * time.Second,
    }

    stopChan := make(chan os.Signal, 1)
    signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        log.Info("🚀 Servidor iniciado en :8080")
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Error al iniciar servidor: %v", err)
        }
    }()

    <-stopChan
    log.Info("🛑 Recibida señal de shutdown, cerrando servidor...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("Error en graceful shutdown: %v", err)
    }
    log.Info("✅ Servidor cerrado correctamente")
}

func ingestHandler(satClient *sat.Client, kafkaWriter *kafka.Writer, log *logger.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var request struct {
            RFC      string `json:"rfc"`
            Password string `json:"password"`
            Period   string `json:"period"`
        }

        if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
            http.Error(w, "JSON inválido", http.StatusBadRequest)
            return
        }

        cfdiData, err := satClient.DownloadCFDIS(r.Context(), request.RFC, request.Password, request.Period)
        if err != nil {
            log.Errorw("Error descargando CFDI", "rfc", request.RFC, "error", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        for _, cfdi := range cfdiData {
            if err := validateCFDI(cfdi); err != nil {
                log.Warnw("CFDI inválido", "rfc", request.RFC, "error", err)
                continue
            }

            cfdiJSON, _ := json.Marshal(cfdi)
            err := kafkaWriter.WriteMessages(r.Context(), kafka.Message{
                Key:   []byte(cfdi.UUID),
                Value: cfdiJSON,
                Time:  time.Now(),
            })
            if err != nil {
                log.Errorw("Error enviando a Kafka", "error", err)
                http.Error(w, "Error procesando datos", http.StatusInternalServerError)
                return
            }
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "status":  "success",
            "message": fmt.Sprintf("%d CFDI procesados", len(cfdiData)),
        })
    }
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func validateCFDI(cfdi *models.CFDI) error {
    if cfdi.RFC == "" || cfdi.UUID == "" {
        return fmt.Errorf("RFC o UUID vacío")
    }
    return nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intVal, err := strconv.Atoi(value); err == nil {
            return intVal
        }
    }
    return defaultValue
}