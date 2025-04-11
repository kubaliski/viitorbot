package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Iniciando bot...")

	// Cargar variables de entorno desde archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error al cargar archivo .env: %v", err)
	}

	// Obtener token desde variables de entorno
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN no está configurado en el archivo .env")
	}

	// Crear una nueva sesión de Discord
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error al crear la sesión de Discord: %v", err)
		return
	}

	// Añadir manejador para el evento Ready
	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Bot está conectado como: %v#%v\n", s.State.User.Username, s.State.User.Discriminator)
	})

	// Añadir manejador para mensajes
	dg.AddHandler(messageCreate)

	// Abrir conexión a Discord
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error al abrir conexión: %v", err)
		return
	}

	fmt.Println("Bot en ejecución. Presiona CTRL+C para salir.")
	select {} // Mantener el bot ejecutándose indefinidamente
}

// Esta función será llamada cada vez que se reciba un mensaje
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignorar mensajes del propio bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Verificar si el mensaje contiene "viitorbot hazlotuyo"
	if strings.Contains(strings.ToLower(m.Content), "viitorbot hazlotuyo") {
		// Obtener contenido de Wikipedia
		title, extract, url, err := getRandomWikipediaArticle()
		if err != nil {
			log.Printf("Error al obtener artículo: %v", err)
			s.ChannelMessageSend(m.ChannelID, "Error al obtener información de Wikipedia: "+err.Error())
			return
		}

		// Enviar la respuesta al canal
		response := fmt.Sprintf("**%s**\n\n%s\n\nMás información: %s", title, extract, url)
		_, err = s.ChannelMessageSend(m.ChannelID, response)
		if err != nil {
			log.Printf("Error al enviar mensaje a Discord: %v", err)
		}
	}
}

// Estructuras para manejar la respuesta de la API de Wikipedia
type WikipediaResponse struct {
	Title       string      `json:"title"`
	Extract     string      `json:"extract"`
	PageID      int         `json:"pageid"`
	ContentURLs ContentURLs `json:"content_urls"`
}

type ContentURLs struct {
	Desktop struct {
		Page string `json:"page"`
	} `json:"desktop"`
}

// Función para obtener un artículo aleatorio de Wikipedia
func getRandomWikipediaArticle() (string, string, string, error) {
	// URL para obtener un artículo aleatorio de Wikipedia en español
	url := "https://es.wikipedia.org/api/rest_v1/page/random/summary"

	// Realizar la petición HTTP
	resp, err := http.Get(url)
	if err != nil {
		return "", "", "", fmt.Errorf("error en la petición HTTP: %w", err)
	}
	defer resp.Body.Close()

	// Verificar código de respuesta
	if resp.StatusCode != http.StatusOK {
		return "", "", "", fmt.Errorf("respuesta HTTP no válida: %d", resp.StatusCode)
	}

	// Leer el cuerpo de la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", fmt.Errorf("error al leer la respuesta: %w", err)
	}

	// Decodificar la respuesta JSON
	var wikiResp WikipediaResponse
	err = json.Unmarshal(body, &wikiResp)
	if err != nil {
		return "", "", "", fmt.Errorf("error al decodificar JSON: %w", err)
	}

	articleURL := wikiResp.ContentURLs.Desktop.Page

	return wikiResp.Title, wikiResp.Extract, articleURL, nil
}
