package llm

type LLMClient struct {
}

func NewLLMClient() *LLMClient {
	return &LLMClient{}
}

type LLM interface {
	SendToHyperbolic(title, description string) (string, error)
}
