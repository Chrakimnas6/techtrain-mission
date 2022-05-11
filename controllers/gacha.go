package controllers

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"net/http"
	"strconv"
	"training/accounts"
	token "training/contracts"
	"training/models"
	"training/repos"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Gacha struct {
	ID    uint `json:"id"`
	Times uint `gorm:"not null" json:"times"`
}

// Draw with times get from request's body
func (controller *Controller) DrawGacha(c *gin.Context) {
	// Get user with the token from Header
	var user models.User
	userToken := c.GetHeader("x-token")
	err := repos.GetUser(controller.Db, &user, userToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// Get gacha times from JSON body
	var gacha Gacha
	// Default value
	gacha.ID = 1
	err = c.BindJSON(&gacha)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	// Get User account
	userAccount, err := accounts.ImportAccount(controller.Keystore, user.Keystore, "password")
	_ = err
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
	// 	return
	// }

	// Get user balance
	// If the instance is nil, get the instance first
	if controller.Instance == nil {
		tkn := models.Token{}
		err := repos.GetToken(controller.Db, &tkn, "MTK")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		// address := common.HexToAddress("0xD1e685605C02f812D4200A16D6844E354ddCDD3C")
		address := common.HexToAddress(tkn.Address)
		instance, err := token.NewToken(address, controller.Client)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		controller.Instance = instance
	}

	userBalance, err := token.GetTokenBalance(controller.Instance, userAccount.Address)
	userBalance = new(big.Int).Div(userBalance, big.NewInt(1000000000000000000))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	// Doens't have enough balance
	if userBalance.Cmp(big.NewInt(int64(gacha.Times))) == -1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":        "Not enough balance",
			"user_balance": userBalance.String(),
		})
		return
	}
	// Consume token
	tx, err := token.BurnToken(controller.Client, controller.Keystore, controller.Instance, user.Keystore, int(gacha.Times))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// Wait transaction to be mined
	err = token.CheckTransaction(tx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	userBalance, err = token.GetTokenBalance(controller.Instance, userAccount.Address)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	fmt.Printf("User balance after the gacha: %v\n", userBalance)

	// Get combined information of users
	type GachaResult struct {
		CharacterID string `json:"characterID"`
		Name        string `json:"name"`
	}

	var gachaCharacterOdds []models.GachaCharacterOdds
	var characters []models.Character

	err = repos.GetCharactersOddsComb(controller.Db, &gachaCharacterOdds, &characters, gacha.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	gachaResults := make([]GachaResult, 0)
	userCharacters := make([]models.UserCharacter, 0)

	// Suppose all cards' possibilities sum up to 1 in a certain type of gacha pool
	for i := 0; i < int(gacha.Times); i++ {
		num := rand.Float64()
		oddsSum := 0.0
		for i := range gachaCharacterOdds {
			oddsSum += gachaCharacterOdds[i].Odds
			if num <= oddsSum {
				userCharacters = append(userCharacters, models.UserCharacter{
					UserID:        user.ID,
					CharacterID:   characters[i].ID,
					Name:          characters[i].Name,
					CharacterRank: characters[i].CharacterRank,
				})
				gachaResults = append(gachaResults, GachaResult{
					CharacterID: strconv.Itoa(int(characters[i].ID)),
					Name:        characters[i].Name,
				})
				break
			}
		}
	}

	err = repos.CreateUserCharacters(controller.Db, &userCharacters)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user balance now": userBalance,
		"results":          gachaResults,
	})

}

// Create new gacha pool
func (controller *Controller) CreateGachaPool(c *gin.Context) {
	var gachaCharacterOddsAll []models.GachaCharacterOdds
	err := repos.GetDefaultGachaPool(controller.Db, &gachaCharacterOddsAll)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	var idNew uint
	err = repos.GetNewestGachaID(controller.Db, &idNew)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// Increment id by 1, indicating the new gacha pool
	idNew += 1
	for i, gachaCharacterOdds := range gachaCharacterOddsAll {
		gachaCharacterOdds.GachaID = idNew
		gachaCharacterOddsAll[i] = gachaCharacterOdds
	}
	err = repos.CreateGachaCharacterOddsAll(controller.Db, &gachaCharacterOddsAll)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gachaCharacterOddsAll)
}

// Update Odds
func (controller *Controller) UpdateOdds(c *gin.Context) {
	var gachaCharacterOdds models.GachaCharacterOdds
	err := c.BindJSON(&gachaCharacterOdds)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	err = repos.UpdateOdds(controller.Db, &gachaCharacterOdds)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gachaCharacterOdds)
}
