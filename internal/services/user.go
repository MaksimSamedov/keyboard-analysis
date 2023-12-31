package services

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"keyboard-analysis/internal/config"
	"keyboard-analysis/internal/models"
	"keyboard-analysis/internal/transport/dto/auth"
)

type UserService struct {
	db   *gorm.DB
	conf config.Config
}

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrDatabase = errors.New("database error")
var ErrUserNotFound = errors.New("user not found")
var ErrUserExists = errors.New("user already exists")
var ErrInvalidToken = errors.New("invalid token")
var ErrTokenExpired = errors.New("token is expired")
var ErrInvalidPasswordLength = errors.New("invalid password length")
var ErrInvalidPasswordsCount = errors.New("invalid passwords count")

func NewUserService(db *gorm.DB, conf config.Config) *UserService {
	return &UserService{db: db, conf: conf}
}

func (s *UserService) Login(dto auth.UserCredentials) (*models.User, error) {
	if err := dto.Validate(); err != nil {
		return nil, err
	}
	usr, err := s.Find(dto.Login)
	if err != nil {
		return nil, err
	}
	if usr == nil {
		return nil, ErrUserNotFound
	}
	if !usr.Attempt(dto.Password) {
		return nil, ErrInvalidCredentials
	}

	return usr, nil
}

func (s *UserService) Register(dto auth.UserCredentials) (*models.User, error) {
	if err := dto.Validate(); err != nil {
		return nil, err
	}

	usr, err := s.Find(dto.Login)
	if err != nil && err != ErrUserNotFound {
		return nil, err
	}
	if usr != nil {
		return nil, ErrUserExists
	}

	usr, err = models.NewUser(dto.Login, dto.Password)
	if err != nil {
		return nil, err
	}

	if err := s.db.Save(usr).Error; err != nil {
		return nil, ErrDatabase
	}

	return usr, nil
}

func (s *UserService) Find(login string) (*models.User, error) {
	var usr *models.User
	if err := s.query().Where("login = ?", login).First(&usr).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, ErrDatabase
	}
	return usr, nil
}

func (s *UserService) RetrieveByCredentials(creds *auth.UserCredentials) (*models.User, error) {
	var usr *models.User
	if err := s.query().Where("login = ?", creds.Login).First(&usr).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, ErrDatabase
	}
	if !usr.Attempt(creds.Password) {
		return nil, ErrInvalidCredentials
	}
	return usr, nil
}

func (s *UserService) RetrieveByToken(accessToken string) (*models.User, error) {
	var token *models.AccessToken
	if err := s.query().Where("token = ?", accessToken).First(&token).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrInvalidToken
		}
		return nil, ErrDatabase
	}
	if token.IsExpired() {
		return nil, ErrTokenExpired
	}

	return &token.User, nil
}

func (s *UserService) SetUserPasswords(usr *models.User, passwords []string) error {
	if uint(len(passwords)) != s.conf.PasswordsCount {
		return ErrInvalidPasswordsCount
	}

	var pws []*models.Password
	for _, pw := range passwords {
		l := uint(len([]rune(pw)))
		if l > s.conf.PasswordMaxLength || l < s.conf.PasswordMinLength {
			return ErrInvalidPasswordLength
		}
		pws = append(pws, models.PasswordFromString(pw))
	}

	if err := usr.SetPasswords(pws); err != nil {
		return err
	}

	if err := s.db.Save(usr).Error; err != nil {
		return ErrDatabase
	}

	return nil
}

func (s *UserService) SetUserSecret(usr *models.User, secretNote string) error {
	usr.SetSecretNote(secretNote)
	if err := s.db.Save(usr).Error; err != nil {
		return ErrDatabase
	}

	return nil
}

func (s *UserService) query() *gorm.DB {
	return s.db.Preload(clause.Associations)
}
