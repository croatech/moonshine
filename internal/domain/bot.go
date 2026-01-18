package domain

type Bot struct {
	Model
	Name    string `db:"name"`
	Slug    string `db:"slug"`
	Attack  uint8  `db:"attack"`
	Defense uint8  `db:"defense"`
	Hp      uint8  `db:"hp"`
	Level   uint8  `db:"level"`
	Avatar  string `db:"avatar"`
}
