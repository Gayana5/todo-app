package llm

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OllamaClient struct {
	Model string
}

func NewOllamaClient(model string) *OllamaClient {
	return &OllamaClient{Model: model}
}

type LLM interface {
	GetAdvice(title, description string) (string, error)
}
