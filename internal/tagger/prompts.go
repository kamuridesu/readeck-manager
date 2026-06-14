package tagger

const (
	TAGGER_PROMPT = `Categorize this bookmark using ONLY the available labels. 
Title: %s
Site: %s
Description: %s
Content: %s

Available Labels: %s
`
)
