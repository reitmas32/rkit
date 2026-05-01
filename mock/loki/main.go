package main

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// LokiPushRequest representa el formato esperado por la API de Loki
type LokiPushRequest struct {
	Streams []LokiStream `json:"streams"`
}

// LokiStream representa un stream de logs con labels
type LokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"` // [[timestamp_nanoseconds, log_line], ...]
}

var totalLogsReceived int64

func main() {
	router := gin.New()

	router.POST("/loki/api/v1/push", func(c *gin.Context) {
		// Leer el body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(400, gin.H{"error": "Failed to read body"})
			return
		}

		// Decodificar el payload de Loki
		var payload LokiPushRequest
		if err := json.Unmarshal(body, &payload); err != nil {
			fmt.Printf("❌ Error decodificando payload: %v\n", err)
			c.JSON(400, gin.H{"error": "Invalid JSON"})
			return
		}

		// Contar logs totales
		totalEntries := 0
		for _, stream := range payload.Streams {
			totalEntries += len(stream.Values)
		}

		totalLogsReceived += int64(totalEntries)

		// Mostrar información
		fmt.Printf("\n📊 === LOKI MOCK - Logs recibidos ===\n")
		fmt.Printf("⏰ Hora: %s\n", time.Now().Format("2006-01-02 15:04:05.000"))
		fmt.Printf("📦 Streams: %d\n", len(payload.Streams))
		fmt.Printf("📝 Logs en este batch: %d\n", totalEntries)
		fmt.Printf("🔢 Total acumulado: %d\n", totalLogsReceived)
		fmt.Printf("─────────────────────────────────────────\n")

		// Mostrar detalles de cada stream
		for i, stream := range payload.Streams {
			fmt.Printf("\n📋 Stream #%d:\n", i+1)
			fmt.Printf("   Labels: %v\n", stream.Stream)
			fmt.Printf("   Logs: %d\n", len(stream.Values))

			// Mostrar los primeros 3 logs de cada stream como ejemplo
			maxLogsToShow := 3
			if len(stream.Values) < maxLogsToShow {
				maxLogsToShow = len(stream.Values)
			}
			for j := 0; j < maxLogsToShow; j++ {
				timestamp := stream.Values[j][0]
				message := stream.Values[j][1]
				// Convertir timestamp de nanosegundos a tiempo legible
				var timestampInt int64
				if _, err := fmt.Sscanf(timestamp, "%d", &timestampInt); err == nil {
					ts := time.Unix(0, timestampInt)
					fmt.Printf("   [%s] %s\n", ts.Format("15:04:05.000"), message)
				} else {
					fmt.Printf("   [%s] %s\n", timestamp, message)
				}
			}
			if len(stream.Values) > maxLogsToShow {
				fmt.Printf("   ... y %d logs más\n", len(stream.Values)-maxLogsToShow)
			}
		}

		fmt.Printf("─────────────────────────────────────────\n\n")

		// Retornar éxito
		c.JSON(200, gin.H{
			"message":           "Logs received",
			"streams":           len(payload.Streams),
			"logs_in_batch":     totalEntries,
			"total_accumulated": totalLogsReceived,
		})
	})

	fmt.Println("🚀 Mock de Loki iniciado en http://localhost:3100")
	fmt.Println("📡 Esperando logs en /loki/api/v1/push\n")
	router.Run(":3100")
}
