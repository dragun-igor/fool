package models

type CardStruct struct {
	Denomination string `json:"denomination"`
	Suit         string `json:"suit"`
}

type Hand []CardStruct

type PutCard struct {
	Card CardStruct `json:"card"`
}
