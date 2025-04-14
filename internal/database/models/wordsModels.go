package models

type Word struct {
	_id          string
	Word         string
	Translations map[string]string
	Description  string
}

type Dictionary struct {
	_id        string
	Dictionary []Word
}

type Words struct {
	Words []Word
}
