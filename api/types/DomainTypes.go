package types

type WordRequestModel struct {
	WordModel
}

type WordModel struct {
	Id           string
	Word         string
	Translations map[string]string
	Description  string
}

type WordsRequestModel struct {
	Words []WordModel
}
