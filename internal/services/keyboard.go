package services

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"keyboard-analysis/internal/config"
	"keyboard-analysis/internal/models"
	"keyboard-analysis/internal/transport/dto/process"
	"keyboard-analysis/internal/utils/flow"
)

type KeyboardService struct {
	db          *gorm.DB
	conf        config.Config
	userService *UserService
}

var ErrGotInvalidPwList = errors.New("got invalid passwords list")
var ErrWeakResult = errors.New("got weak result")
var ErrNeedMoreSamples = errors.New("need more samples")
var ErrSampleDontSatisfy = errors.New("sample don't satisfy")

func NewKeyboardService(db *gorm.DB, conf config.Config, userService *UserService) *KeyboardService {
	return &KeyboardService{
		db:          db,
		conf:        conf,
		userService: userService,
	}
}

func (s *KeyboardService) ProcessFlow(dto *process.KeyboardFlowResults) (*models.User, []*models.KeyboardFlow, error) {
	// Авторизуем пользователя
	usr, err := s.userService.RetrieveByCredentials(&dto.Auth)
	if err != nil {
		return nil, nil, err
	}

	// Спарсим и проверим полученные замеры
	flowsToSave, err := s.validateFlows(usr, dto.Flows)
	if err != nil {
		return nil, nil, err
	}

	// Сохраним в бд
	if err := s.db.Save(flowsToSave).Error; err != nil {
		return nil, nil, err
	}

	return usr, flowsToSave, nil
}

func (s *KeyboardService) validateFlows(usr *models.User, flows []process.KeyboardFlowResult) ([]*models.KeyboardFlow, error) {

	var filled = make(map[string]*models.Password)
	for _, pw := range usr.Passwords {
		filled[pw.Password] = pw
	}

	var flowsToSave []*models.KeyboardFlow
	for _, f := range flows {
		pw, ok := filled[f.Phrase]
		if !ok || pw == nil {
			return nil, ErrGotInvalidPwList
		}
		safeFlow := &models.KeyboardFlow{
			Flow:     f.Flow,
			Password: *pw,
		}
		safeFlow.RemoveInvalidEvents()
		if len(safeFlow.Flow) < len(safeFlow.Password.Password) {
			return nil, ErrWeakResult
		}
		safeFlow.TruncateTime()
		flowsToSave = append(flowsToSave, safeFlow)
	}

	if len(flowsToSave) != len(usr.Passwords) {
		return nil, ErrGotInvalidPwList
	}

	return flowsToSave, nil
}

func (s *KeyboardService) History(usr *models.User) ([]models.KeyboardFlow, error) {
	query := s.query()
	query.Where(
		"password_id IN (?)",
		s.query().
			Model(models.Password{}).
			Select("id").
			Where("user_id = (?)", usr.ID),
	)

	var res []models.KeyboardFlow
	if err := query.Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (s *KeyboardService) SingleHistory(usr *models.User, id uint) (*models.KeyboardFlow, error) {
	query := s.query().Where("id = ?", id)
	query.Where(
		"password_id IN (?)",
		s.query().
			Model(models.Password{}).
			Select("id").
			Where("user_id = (?)", usr.ID),
	)

	var res *models.KeyboardFlow
	if err := query.First(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (s *KeyboardService) GetToken(dto *process.KeyboardFlowResults) (*models.AccessToken, error) {
	// Авторизуем пользователя
	usr, err := s.userService.RetrieveByCredentials(&dto.Auth)
	if err != nil {
		return nil, err
	}

	// Спарсим и проверим полученные замеры
	validFlows, err := s.validateFlows(usr, dto.Flows)
	if err != nil {
		return nil, err
	}

	// Проверим что замеры совпадают с теми, что лежат в базе
	valid, err := s.CompareFlows(usr, validFlows)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, ErrSampleDontSatisfy
	}

	// generate new token
	token := models.NewToken(*usr, s.conf.TokenLength, s.conf.TokenLifetime)
	if err := s.query().Save(&token).Error; err != nil {
		return nil, ErrDatabase
	}

	return token, nil
}

func (s *KeyboardService) CompareFlows(usr *models.User, flows []*models.KeyboardFlow) (bool, error) {
	// распределим флоу по паролям
	pwMap := make(map[string]*models.KeyboardFlow)
	for _, pwFlow := range flows {
		pwMap[pwFlow.Password.Password] = pwFlow
	}

	// проверим, что все пароли указаны
	for _, password := range usr.Passwords {
		if val, ok := pwMap[password.Password]; !ok || val == nil {
			return false, ErrGotInvalidPwList
		}
	}

	// проверим сходство
	analyser := flow.NewAnalyser(s.conf.AnalyserProps)
	for _, password := range usr.Passwords {
		got := pwMap[password.Password]
		samples, err := s.GetSamples(password, s.conf.SamplesCompareCount)
		if err != nil {
			return false, err
		}
		if len(samples) < s.conf.MinSamples {
			return false, ErrNeedMoreSamples
		}
		analyser.AddTask(got, samples)
	}
	fits, err := analyser.Analyse()
	if err != nil {
		return false, err
	}
	return fits, nil
}

func (s *KeyboardService) GetSamples(pw *models.Password, count int) ([]*models.KeyboardFlow, error) {
	var res []*models.KeyboardFlow
	query := s.query().Where("password_id = ?", pw.ID).Order("id DESC").Limit(count)
	if err := query.Find(&res).Error; err != nil {
		return nil, ErrDatabase
	}
	return res, nil
}

func (s *KeyboardService) NeedSamples(usr *models.User) (bool, error) {
	for _, pw := range usr.Passwords {
		samples, err := s.GetSamples(pw, s.conf.MinSamples)
		if err != nil {
			return false, err
		}
		if len(samples) < s.conf.MinSamples {
			return true, nil
		}
	}
	return false, nil
}

func (s *KeyboardService) query() *gorm.DB {
	return s.db.Preload(clause.Associations)
}
