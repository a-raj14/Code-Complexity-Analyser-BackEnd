package main

import (
	"code-complexity-api/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Println("⚠️ Could not load .env file")
	// }
	// apiKey := os.Getenv("OPENAI_API_KEY")
	// log.Println("Using API Key:", apiKey[:5]+"...") // Just for checking

	// response, err := handlers.GetChatGPTResponse("What's the capital of France?")
	// if err != nil {
	// 	fmt.Println("❌ Error:", err)
	// 	return
	// }
	// fmt.Println("🤖 ChatGPT says:", response)

    http.HandleFunc("/analyze", handlers.WithCORS(handlers.AnalyzeHandler))
	// http.HandleFunc("/chatgpt", withCORS(chatGPTHandler))

    fmt.Println("🚀 Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
