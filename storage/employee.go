package storage

import (
	"github.com/Tus1688/library-management-api/types"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func (s *PostgresStore) CreateEmployee(req *types.CreateEmployee) (types.CreateId, int, types.Err) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return types.CreateId{}, 500, types.Err{Error: "unable to create employee"}
	}
	var userId types.CreateId
	err = s.db.QueryRow(`INSERT INTO employees(username, password) VALUES ($1, $2) RETURNING id`,
		req.Username, string(hashedPassword)).Scan(&userId.Id)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return types.CreateId{}, 409, types.Err{Error: "username already exists"}
		}

		return types.CreateId{}, 500, types.Err{Error: "unable to create employee"}
	}

	return userId, 0, types.Err{}
}

func (s *PostgresStore) GetEmployee() ([]types.ListEmployee, int, types.Err) {
	rows, err := s.db.Query(`SELECT id, username, created_at, updated_at FROM employees`)
	if err != nil {
		return nil, 500, types.Err{Error: "unable to get employees"}
	}
	defer rows.Close()

	var employees []types.ListEmployee
	for rows.Next() {
		var employee types.ListEmployee
		err := rows.Scan(&employee.Id, &employee.Username, &employee.CreatedAt, &employee.UpdatedAt)
		if err != nil {
			return nil, 500, types.Err{Error: "unable to get employees"}
		}
		employees = append(employees, employee)
	}

	return employees, 0, types.Err{}
}

func (s *PostgresStore) DeleteEmployee(currentUserId, id *string) (int, types.Err) {
	if *currentUserId == *id {
		return 403, types.Err{Error: "cannot delete yourself"}
	}

	res, err := s.db.Exec(`DELETE FROM employees WHERE id = $1`, *id)
	if err != nil {
		if strings.Contains(err.Error(), "foreign key constraint") {
			return 409, types.Err{Error: "employee is being used"}
		}
		if strings.Contains(err.Error(), "uuid") {
			return 400, types.Err{Error: "invalid id"}
		}

		return 500, types.Err{Error: "unable to delete employee"}
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 500, types.Err{Error: "unable to delete employee"}
	}

	if rowsAffected == 0 {
		return 404, types.Err{Error: "employee not found"}
	}

	return 200, types.Err{}
}
