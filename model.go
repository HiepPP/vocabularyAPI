package main

type Vocabulary struct {
	WordTitle       string
	DefinitionShort string
	DefinitionLong  string
	Definition      []Definition
}

type Definition struct {
	Title    string
	Type     string
	Example  string
	Synonyms DeepDefinition
	Antonyms DeepDefinition
	Types    DeepDefinition
}

type DeepDefinition struct {
	ListWord    []string
	Description string
}

type FetchWord struct{
	Word string
	ShortDescription string
}