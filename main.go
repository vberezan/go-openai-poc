package main

import (
	"github.com/gin-gonic/gin"
	se "go-openai-poc/extractor"
	oai "go-openai-poc/openai"
	"log"
	"os"
)

type chatRequest struct {
	Prompt string `json:"prompt"`
	URL    string `json:"url"`
}

func main() {
	var apiKey = os.Getenv("OPENAI_API_KEY")
	var apiUrl = os.Getenv("OPENAI_CHAT_URL")
	var apiGPTModel = os.Getenv("OPENAI_GPT_MODEL")

	log.Println("OpenAI API Key: ", apiKey)
	log.Println("OpenAI API URL: ", apiUrl)

	var router = gin.Default()

	router.POST("/chat", func(ctx *gin.Context) {
		var req chatRequest

		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(400, gin.H{"Invalid request payload": err.Error()})
			return
		}

		var source, sourceErr = se.Extract(req.URL)
		if sourceErr != nil {
			ctx.JSON(500, gin.H{"Failed to extract source": sourceErr.Error()})
			return
		}

		var openAIRequest = oai.CreateCompletionsRequest{}
		openAIRequest.Messages = []oai.Message{
			{Role: "user", Content: req.Prompt + " " + source},
		}
		openAIRequest.N = 2
		openAIRequest.Model = apiGPTModel
		openAIRequest.MaxTokens = 100

		log.Printf("OpenAI Request: %+v\n", openAIRequest)

		var responseText, err = oai.CallOpenAIChat(openAIRequest, apiKey, apiUrl)
		if err != nil {
			ctx.JSON(500, gin.H{"Failed to call OpenAI API": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"response": responseText})
	})

	err := router.Run(":" + os.Getenv("APPLICATION_PORT"))
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
