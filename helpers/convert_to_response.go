package helpers

import (
	"strconv"
	"training/models"
	"training/repos"

	"gorm.io/gorm"
)

func Convert(db *gorm.DB, userCharacters *[]models.UserCharacter, userCharactersResponses *[]models.UserCharacterResponse) (err error) {
	for _, userCharacter := range *userCharacters {
		// Get character by character ID
		var character models.Character
		err := repos.GetCharacter(db, &character, uint(userCharacter.CharacterID))
		if err != nil {
			return err
		}
		// Convert fileds into string according to the reponse
		*userCharactersResponses = append(*userCharactersResponses, models.UserCharacterResponse{
			UserCharacterID: strconv.Itoa(int(userCharacter.ID)),
			CharacterID:     strconv.Itoa(int(userCharacter.CharacterID)),
			Name:            character.Name,
		})
	}
	return nil
}
