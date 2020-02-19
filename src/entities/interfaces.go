package entities

type DAppEntity interface {
	GetKeys(*string) []string
	MapItemToModel(string, map[string]string) *DAppEntity
	UpdateAll(*map[string]string) []DAppEntity
}

