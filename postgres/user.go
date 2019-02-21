package postgres

import (
	"encoding/json"

	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/shjp/shjp-core"
	"github.com/shjp/shjp-core/model"
)

// UserQueryStrategy implements QueryStrategy for users
type UserQueryStrategy struct {
	*pg.DB
}

type user struct {
	model.User

	tableName struct{} `sql:"select:users_full"`
}

// ModelName outputs this model's name
func (s *UserQueryStrategy) ModelName() string {
	return "user"
}

// GetAll returns all users
func (s *UserQueryStrategy) GetAll() ([]core.Model, error) {
	users := make([]*user, 0)

	if err := s.DB.Model(&users).Select(); err != nil {
		return nil, err
	}

	result := make([]core.Model, len(users))
	for i, u := range users {
		result[i] = core.Model(u)
	}

	return result, nil
}

// GetOne returns one user
func (s *UserQueryStrategy) GetOne(id string) (core.Model, error) {
	var u user
	u.ID = id
	if err := s.DB.Model(&u).WherePK().First(); err != nil {
		return nil, err
	}

	return &u, nil
}

// Search finds all users meeting the criteria given by the payload
func (s *UserQueryStrategy) Search(payload []byte) ([]core.Model, error) {
	var params user
	if err := json.Unmarshal(payload, &params); err != nil {
		return nil, errors.Wrap(err, "Error deserializing payload")
	}

	us := make([]*user, 0)

	query := s.DB.Model(&us)

	if params.ID != "" {
		query = query.Where("id = ?", params.ID)
	}

	if params.Name != nil {
		query = query.Where("name ilike ?", "%"+*params.Name+"%")
	}

	if params.Email != nil {
		query = query.Where("email = ?", *params.Email)
	}

	if params.AccountType != nil {
		query = query.Where("account_type = ?", *params.AccountType)
	}

	if err := query.Select(); err != nil {
		return nil, err
	}

	result := make([]core.Model, 0)

	for _, u := range us {
		if params.Password != nil {
			if err := bcrypt.CompareHashAndPassword([]byte(*u.AccountSecret), []byte(*params.Password)); err != nil {
				continue
			}
		}
		result = append(result, core.Model(u))
	}

	return result, nil
}

// Upsert upserts a user
func (s *UserQueryStrategy) Upsert(m core.Model) error {
	u := m.(*core.User)
	if err := populateAccountSecret(u); err != nil {
		return errors.Wrap(err, "Error while performing transformation before upsert")
	}
	_, err := s.DB.Model(m).
		OnConflict("(id) DO UPDATE").
		Set(`(
			name,
			email,
			baptismal_name,
			birthday,
			feastday,
			last_active,
			account_type,
			account_secret
		) = (
			?name,
			?email,
			?baptismal_name,
			?birthday,
			?feastday,
			?last_active,
			?account_type,
			?account_secret)`).
		Insert(u)
	return err
}

// UpsertRelationship upserts a user relationship
func (s *UserQueryStrategy) UpsertRelationship(e core.Entity, relation string) error {
	return nil
}

// populateAccountSecret populates the account secret from the given
// user's secret seed and account type
func populateAccountSecret(u *core.User) error {
	accountSecret, err := getEmailAccountSecret(*u.Password)
	if err != nil {
		return err
	}
	u.AccountSecret = &accountSecret
	return nil
}

// getEmailAccountSecret produces account secret for email account from password
func getEmailAccountSecret(password string) (string, error) {
	hashByte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.Wrap(err, "failed generating hash")
	}
	return string(hashByte), nil
}
