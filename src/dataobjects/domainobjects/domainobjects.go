package domainobjects

type Video struct {
	ID                    string `json:"id"`
	Title                 string `json:"title"`
	StandardDefinitionURL string `json:"standard_definition_url"`
	PremiumDefinitionURL  string `json:"premium_definition_url"`
}
