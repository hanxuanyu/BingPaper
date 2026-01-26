package token

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"BingDailyImage/internal/config"
	"BingDailyImage/internal/model"
	"BingDailyImage/internal/repo"

	"golang.org/x/crypto/bcrypt"
)

func GenerateTokenString() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func CreateToken(name string, expiresAt time.Time) (*model.Token, error) {
	tString := GenerateTokenString()
	t := &model.Token{
		Token:     tString,
		Name:      name,
		ExpiresAt: expiresAt,
	}
	if err := repo.DB.Create(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func ValidateToken(tokenStr string) (*model.Token, error) {
	var t model.Token
	if err := repo.DB.Where("token = ? AND disabled = ?", tokenStr, false).First(&t).Error; err != nil {
		return nil, err
	}
	if time.Now().After(t.ExpiresAt) {
		return nil, errors.New("token expired")
	}
	return &t, nil
}

func Login(password string) (*model.Token, error) {
	cfg := config.GetConfig()
	err := bcrypt.CompareHashAndPassword([]byte(cfg.Admin.PasswordBcrypt), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	ttl := config.GetTokenTTL()
	return CreateToken("login-token", time.Now().Add(ttl))
}

func ListTokens() ([]model.Token, error) {
	var tokens []model.Token
	err := repo.DB.Order("id desc").Find(&tokens).Error
	return tokens, err
}

func UpdateToken(id uint, disabled bool) error {
	return repo.DB.Model(&model.Token{}).Where("id = ?", id).Update("disabled", disabled).Error
}

func DeleteToken(id uint) error {
	return repo.DB.Delete(&model.Token{}, id).Error
}
