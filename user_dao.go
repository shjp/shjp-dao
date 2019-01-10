package dao

import (
	"encoding/json"

	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/shjp/shjp-core"
	"github.com/shjp/shjp-core/model"
)

type userDAO struct {
	DB *pg.DB
}

type user struct {
	model.User

	tableName struct{} `sql:"select:users_full"`
}

// GetAll returns all users
func (o *userDAO) GetAll() ([]core.Model, error) {
	users := make([]*user, 0)

	if err := o.DB.Model(&users).Select(); err != nil {
		return nil, err
	}

	result := make([]core.Model, len(users))
	for i, u := range users {
		result[i] = core.Model(u)
	}

	return result, nil
}

// GetOne returns one user
func (o *userDAO) GetOne(id string) (core.Model, error) {
	var u user
	var err error
	u.ID = id
	if err := o.DB.Model(&u).First(); err != nil {
		return nil, err
	}

	return &u, err
}

// Search finds all users meeting the criteria given by the payload
func (o *userDAO) Search(payload []byte) ([]core.Model, error) {
	var params user
	if err := json.Unmarshal(payload, &params); err != nil {
		return nil, errors.Wrap(err, "Error deserializing payload")
	}

	us := make([]*user, 0)

	query := o.DB.Model(&us)

	if params.ID != "" {
		query = query.Where("id = ?", params.ID)
	}

	if params.Name != nil {
		query = query.Where("name ilike ?", "%"+*params.Name+"%")
	}

	if params.Email != nil {
		query = query.Where("email = ?", *params.Email)
	}

	if params.Password != nil {
		if err := populateAccountSecret(&params.User.User); err != nil {
			return nil, errors.Wrap(err, "Error populating account secret")
		}
		query = query.Where("account_secret = ?", *params.AccountSecret)
	}

	if params.AccountType != nil {
		query = query.Where("account_type = ?", *params.AccountType)
	}

	if err := query.Select(); err != nil {
		return nil, err
	}

	result := make([]core.Model, len(us))
	for i, u := range us {
		result[i] = core.Model(u)
	}

	return result, nil
}

// Upsert upserts a user
func (o *userDAO) Upsert(m core.Model) error {
	u := m.(*core.User)
	if err := populateAccountSecret(u); err != nil {
		return errors.Wrap(err, "Error while performing transformation before upsert")
	}
	return o.DB.Insert(u)
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
