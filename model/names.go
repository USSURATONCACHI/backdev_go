package model

import (
	"github.com/google/uuid"
	"crypto/sha512"
)


func (model *Model) CreateSyllablesWord(hash [6]byte) string {
	startSylId := int(hash[0] + hash[1]) % len(model.StartSyllables)
	middleSylId := int(hash[2] + hash[3]) % len(model.MiddleSyllables)
	finalSylId := int(hash[4] + hash[5]) % len(model.FinalSyllables)

	word := model.StartSyllables[startSylId] + model.MiddleSyllables[middleSylId] + model.FinalSyllables[finalSylId]
	word = CapitalizeFirst(word)

	return word
}

func (model *Model) GenerateNameFromUuid(u uuid.UUID) string {
	hashsum := sha512.Sum512(u[:])

	var slice1 [6]byte
	copy(slice1[:], hashsum[0:6])

	var slice2 [6]byte
	copy(slice2[:], hashsum[6:12])

	firstName := model.CreateSyllablesWord(slice1)
	lastName := model.CreateSyllablesWord(slice2)

	fullName := firstName + " " + lastName

	return fullName
}