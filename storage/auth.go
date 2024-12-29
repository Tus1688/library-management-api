package storage

import (
	"database/sql"
	"errors"
	"github.com/Tus1688/library-management-api/types"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (s *PostgresStore) InitAdmin(username, password *string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`INSERT INTO employees(username, password) VALUES ($1, $2) ON CONFLICT (username) DO
	UPDATE SET password = $2`, *username, string(hashedPassword))
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) Login(req *types.LoginRequest) (string, int, types.Err) {
	var hashedPassword, userId string
	err := s.db.QueryRow(`SELECT id, password FROM employees WHERE username = $1`, req.Username).
		Scan(&userId, &hashedPassword)
	if err != nil {
		// random sleep to simulate query / bcrypt time
		duration := time.Duration(100+(time.Now().UnixNano()%400)) * time.Millisecond
		time.Sleep(duration)

		if errors.Is(err, sql.ErrNoRows) {
			return "", 401, types.Err{Error: "invalid username or password"}
		}

		return "", 500, types.Err{Error: "error logging in"}
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		// random sleep to simulate query / bcrypt time
		duration := time.Duration(100+(time.Now().UnixNano()%400)) * time.Millisecond
		time.Sleep(duration)
		return "", 401, types.Err{Error: "invalid username or password"}
	}

	return userId, 200, types.Err{}
}
