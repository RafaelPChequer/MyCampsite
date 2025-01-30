package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Request struct {
	Prompt string `json:"prompt"`
}

type Response struct {
	Response string `json:"response"`
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	// Obtendo o prompt da URL
	prompt := r.URL.Query().Get("prompt")
	if prompt == "" {
		http.Error(w, "Prompt é obrigatório", http.StatusBadRequest)
		return
	}

	// Formatar a requisição para o modelo
	data := []byte(`{"model":"mistral","prompt":"` + prompt + `"}`)

	// Enviar a requisição para o servidor local do Ollama
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(data))
	if err != nil {
		http.Error(w, "Erro ao chamar modelo", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Ler a resposta e agrupar as partes
	var responseString string
	decoder := json.NewDecoder(resp.Body)
	for {
		var part Response
		if err := decoder.Decode(&part); err != nil {
			if err.Error() == "EOF" {
				break
			}
			http.Error(w, "Erro ao ler resposta", http.StatusInternalServerError)
			return
		}
		responseString += part.Response
	}

	// Enviar a resposta completa para o cliente
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"response": "` + responseString + `"}`))
}

func main() {
	http.HandleFunc("/chat", chatHandler)
	fmt.Println("Servidor rodando em http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
