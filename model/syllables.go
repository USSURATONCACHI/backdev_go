package model

import (
	"github.com/google/uuid"
	"crypto/sha512"
)

type Syllables struct {
	Start  []string
	Middle []string
	Final  []string
}

func (syl *Syllables) Word(hash [6]byte) string {
	startSylId := int(hash[0] + hash[1]) % len(syl.Start)
	middleSylId := int(hash[2] + hash[3]) % len(syl.Middle)
	finalSylId := int(hash[4] + hash[5]) % len(syl.Final)

	word := syl.Start[startSylId] + syl.Middle[middleSylId] + syl.Final[finalSylId]
	word = CapitalizeFirst(word)

	return word
}

func (syl *Syllables) HumanNameFromUuid(u uuid.UUID) string {
	hashsum := sha512.Sum512(u[:])

	var slice1 [6]byte
	copy(slice1[:], hashsum[0:6])

	var slice2 [6]byte
	copy(slice2[:], hashsum[6:12])

	firstName := syl.Word(slice1)
	lastName := syl.Word(slice2)

	fullName := firstName + " " + lastName

	return fullName
}