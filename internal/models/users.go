package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	Id             int       `db:"id"`
	Name           string    `db:"name"`
	Email          string    `db:"email"`
	HashedPassword string    `db:"hashed_password"`
	Created        time.Time `db:"created"`
}
type UsersNoPassword struct {
	Id      int       `db:"id"`
	Name    string    `db:"name"`
	Email   string    `db:"email"`
	Created time.Time `db:"created"`
}

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
	Get(id int) (UsersNoPassword, error)
	PasswordUpdate(id int, currentPassword, newPassword string) error
}

type UserModel struct {
	Pool *pgxpool.Pool
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (name, email, hashed_password, created) 
			VALUES(@name, @email, @hashedPassword, @createdAt)`

	args := pgx.NamedArgs{
		"name":           name,
		"email":          email,
		"hashedPassword": hashedPassword,
		"createdAt":      time.Now(),
	}

	commandTag, err := m.Pool.Exec(context.Background(), query, args)

	if err != nil {
		var pgErr *pgconn.PgError

		//we basically check is this a postgre error? if not return the error
		if !errors.As(err, &pgErr) {
			return err
		}

		//23505 is postgre error code for unique value constraint violation.
		if pgErr.Code == "23505" {
			return ErrDuplicateEmail
		}

	}

	if commandTag.RowsAffected() != 1 {
		return err
	}

	fmt.Println(string(hashedPassword))
	return nil
}

// return id, err
func (m *UserModel) Authenticate(email, password string) (int, error) {
	query := `SELECT * from users where email = @email`
	args := pgx.NamedArgs{
		"email": email,
	}

	rows, err := m.Pool.Query(context.Background(), query, args)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[Users])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	return user.Id, nil
}

// return true or f, why not just return the user idk
func (m *UserModel) Exists(id int) (bool, error) {
	query := `SELECT * from users where id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}

	rows, err := m.Pool.Query(context.Background(), query, args)
	if err != nil {
		return false, err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[Users])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, ErrInvalidCredentials
		}

		return false, err
	}
	return true, nil
}

func (m *UserModel) Get(id int) (UsersNoPassword, error) {
	query := `SELECT id, name, email, created FROM users WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}

	rows, err := m.Pool.Query(context.Background(), query, args)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UsersNoPassword{}, ErrInvalidCredentials
		}

		return UsersNoPassword{}, err
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[UsersNoPassword])
	if err != nil {
		return UsersNoPassword{}, nil
	}
	return *user, nil
}

func (m *UserModel) PasswordUpdate(id int, currentPassword, newPassword string) error {
	// select the data, get the hash

	query := `SELECT hashed_password FROM users where id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}

	rows, err := m.Pool.Query(context.Background(), query, args)
	if err != nil {
		return err
	}

	oldHashed, err := pgx.CollectOneRow(rows, pgx.RowTo[string])
	if err != nil {
		return err
	}

	// compare the hash, return errInvalidcreden if fail

	err = bcrypt.CompareHashAndPassword([]byte(oldHashed), []byte(currentPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		}
		return err
	}

	// make a hash
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}

	//update password
	query = `UPDATE users SET hashed_password = @newHashed WHERE id = @id`
	args = pgx.NamedArgs{
		"id":        id,
		"newHashed": newHashedPassword,
	}

	commandTag, err := m.Pool.Exec(context.Background(), query, args)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return err
	}

	return nil
}
