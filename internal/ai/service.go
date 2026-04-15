package ai

import "hybridsearch/internal/ai/gemini"

type Service struct {
	client *gemini.Client
}

func NewService() *Service {
	return &Service{
		client: gemini.NewClient(),
	}
}

func (s *Service) Ask(query string) (string, error) {
	return s.client.Generate(query)
}
