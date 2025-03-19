package repository

import (
	"context"
	"errors"
	"log"
)

const userByIDQuery = `SELECT id, username FROM users WHERE id = $1`

func (q *querier) UserByID(ctx context.Context, id string) (res *User, err error) {
	log.Println("sql.user.byID")
	row := q.client.QueryRow(ctx, userByIDQuery, id)
	var user User
	if err := row.Scan(&user.ID, &user.Username); err != nil {
		return nil, err
	}
	return &user, nil
}

const emailVerificationCodeStmt = `SELECT created_at, expiry,code FROM email_verification_codes WHERE user_id = $1`

func (q *querier) EmailVerificationCode(ctx context.Context, userID string) (res *EmailVerificationCode, err error) {
	log.Println("sql.user.emailVerificationCode")

	res = new(EmailVerificationCode)
	err = q.client.QueryRow(ctx, emailVerificationCodeStmt, userID).
		Scan(
			&res.CreatedAt,
			&res.Expiry,
			&res.Code,
		)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (e *executor) CreateUser(ctx context.Context, user *User) (res *User, err error) {
	log.Println("sql.user.create")
	err = e.client.Exec(ctx, "INSERT INTO users (id, username) VALUES ($1, $2)", user.ID, user.Username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (e *executor) EmailVerificationFailed(ctx context.Context, userID string) error {
	return errors.New("not implemented")
}

func (e *executor) EmailVerificationSucceeded(ctx context.Context, userID string) error {
	return errors.New("not implemented")
}
