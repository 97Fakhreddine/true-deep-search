package bleve

import (
	blevev2 "github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
)

func newIndexMapping() mapping.IndexMapping {
	indexMapping := blevev2.NewIndexMapping()

	docMapping := blevev2.NewDocumentMapping()

	textFieldMapping := blevev2.NewTextFieldMapping()
	textFieldMapping.Store = true
	textFieldMapping.Index = true
	textFieldMapping.IncludeInAll = true

	keywordFieldMapping := blevev2.NewKeywordFieldMapping()
	keywordFieldMapping.Store = true
	keywordFieldMapping.Index = true
	keywordFieldMapping.IncludeInAll = false

	docMapping.AddFieldMappingsAt("title", textFieldMapping)
	docMapping.AddFieldMappingsAt("content", textFieldMapping)
	docMapping.AddFieldMappingsAt("path", keywordFieldMapping)

	indexMapping.DefaultMapping = docMapping
	indexMapping.DefaultType = "_default"

	return indexMapping
}
