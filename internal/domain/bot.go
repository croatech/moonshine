package domain

type Bot struct {
	Model
	Name    string `db:"name"`
	Slug    string `db:"slug"`
	Avatar  string `db:"avatar"`
	Attack  uint16 `db:"attack"`
	Defense uint16 `db:"defense"`
	Hp      uint16 `db:"hp"`
	Level   uint8  `db:"level"`
}
