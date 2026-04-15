package open

import (
	"fmt"
	"strings"

	"hybridsearch/internal/infra/browser"
	"hybridsearch/internal/search"
)

func OpenResult(r search.SearchResult) error {
	target := strings.TrimSpace(r.Target)
	if target == "" {
		return fmt.Errorf("empty target")
	}

	return browser.Open(target)
}
