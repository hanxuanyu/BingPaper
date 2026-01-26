package token

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"BingPaper/internal/config"
	"BingPaper/internal/model"
	"BingPaper/internal/repo"

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
	expiresAt := time.Now().Add(ttl)
	name := "login-token"

	// 如果已存在同名 token，则刷新时间并返回
	var t model.Token
	if err := repo.DB.Where("name = ?", name).First(&t).Error; err == nil {
		t.ExpiresAt = expiresAt
		t.Disabled = false
		if err := repo.DB.Save(&t).Error; err != nil {
			return nil, err
		}
		return &t, nil
	}

	return CreateToken(name, expiresAt)
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
