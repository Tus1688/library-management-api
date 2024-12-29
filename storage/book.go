package storage

import (
	"github.com/Tus1688/library-management-api/types"
	"strconv"
	"strings"
)

func (s *PostgresStore) GetBook(searchQuery *string, lastId, limit *int) ([]types.ListBook, int, types.Err) {
	query := `SELECT id, pagination_id, title, author, description, is_booked, COALESCE(booked_until::TEXT,''),
	created_at, updated_at FROM books`
	var args []interface{}
	argsCount := 1

	if *searchQuery != "" || *lastId != 0 {
		query += ` WHERE`
	}

	if *searchQuery != "" {
		query += ` title ILIKE '%' ||` + ` $` + strconv.Itoa(argsCount) + ` || '%'`
		args = append(args, *searchQuery)
		argsCount++
	}

	if *lastId != 0 {
		if argsCount > 1 {
			query += ` AND`
		}
		query += ` pagination_id < $` + strconv.Itoa(argsCount)
		args = append(args, *lastId)
		argsCount++
	}

	query += ` ORDER BY pagination_id DESC`

	if *limit != 0 {
		query += ` LIMIT $` + strconv.Itoa(argsCount)
		args = append(args, *limit)
		argsCount++
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, 500, types.Err{Error: "unable to get books"}
	}
	defer rows.Close()

	var books []types.ListBook
	for rows.Next() {
		var book types.ListBook
		err := rows.Scan(&book.Id, &book.PaginationId, &book.Title, &book.Author, &book.Description,
			&book.IsBooked, &book.BookedUntil, &book.CreatedAt, &book.UpdatedAt)
		if err != nil {
			return nil, 500, types.Err{Error: "unable to get books"}
		}
		books = append(books, book)
	}

	if len(books) == 0 {
		return nil, 404, types.Err{Error: "no books found"}
	}

	return books, 200, types.Err{}
}

func (s *PostgresStore) CreateBook(req *types.CreateBook) (types.CreateId, int, types.Err) {
	var id types.CreateId
	err := s.db.QueryRow(`INSERT INTO books(title, author, description) VALUES ($1, $2, $3) RETURNING id`,
		req.Title, req.Author, req.Description).Scan(&id.Id)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return types.CreateId{}, 409, types.Err{Error: "book already exists"}
		}

		return types.CreateId{}, 500, types.Err{Error: "unable to create book"}
	}

	return id, 201, types.Err{}
}

func (s *PostgresStore) DeleteBook(id *string) (int, types.Err) {
	res, err := s.db.Exec(`DELETE FROM books WHERE id = $1`, *id)
	if err != nil {
		if strings.Contains(err.Error(), "foreign key constraint") {
			return 409, types.Err{Error: "book is being used"}
		}
		if strings.Contains(err.Error(), "uuid") {
			return 400, types.Err{Error: "invalid id"}
		}

		return 500, types.Err{Error: "unable to delete book"}
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 500, types.Err{Error: "unable to delete book"}
	}

	if rowsAffected == 0 {
		return 404, types.Err{Error: "book not found"}
	}

	return 200, types.Err{}
}

func (s *PostgresStore) UpdateBook(req *types.UpdateBook) (int, types.Err) {
	res, err := s.db.Exec(`UPDATE books SET title = $1, author = $2, description = $3 WHERE id = $4`,
		req.Title, req.Author, req.Description, req.Id)
	if err != nil {
		if strings.Contains(err.Error(), "uuid") {
			return 400, types.Err{Error: "invalid id"}
		}

		if strings.Contains(err.Error(), "duplicate") {
			return 409, types.Err{Error: "book with that name already exists"}
		}

		return 500, types.Err{Error: "unable to update book"}
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 500, types.Err{Error: "unable to update book"}
	}

	if rowsAffected == 0 {
		return 404, types.Err{Error: "book not found"}
	}

	return 200, types.Err{}
}
