package domain

type Bot struct {
	Model
	Name    string `json:"name"`
	Attack  uint8  `json:"attack"`
	Defense uint8  `json:"defense"`
	Hp      uint8  `json:"hp"`
	Level   uint8  `json:"level"`
	Avatar  string `json:"avatar"`
}
