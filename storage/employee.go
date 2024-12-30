package storage

import (
	"github.com/Tus1688/library-management-api/types"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

// CreateEmployee creates a new employee with the given details.
// It hashes the password using bcrypt and inserts the employee into the employees table.
// If the username already exists, it returns a 409 status code.
// Parameters:
// - req: a pointer to the CreateEmployee request containing the employee details
// Returns the created employee ID, status code, and an error if the operation fails.
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

// GetEmployee retrieves a list of all employees from the database.
// It constructs a SQL query to fetch employees and returns the list of employees.
// Returns a slice of ListEmployee, status code, and an error if the operation fails.
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

// DeleteEmployee deletes an employee from the database based on the provided ID.
// It checks if the current user is trying to delete themselves and returns a 403 status code if so.
// If the employee is being used, it returns a 409 status code.
// Parameters:
// - currentUserId: a pointer to the current user's ID
// - id: a pointer to the employee ID to be deleted
// Returns the status code and an error if the operation fails.
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
