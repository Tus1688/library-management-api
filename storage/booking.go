package storage

import (
	"database/sql"
	"errors"
	"github.com/Tus1688/library-management-api/types"
	"strconv"
	"strings"
)

func (s *PostgresStore) CreateBooking(uid *string, req *types.CreateBooking) (types.CreateId, int, types.Err) {
	tx, err := s.db.Begin()
	if err != nil {
		return types.CreateId{}, 500, types.Err{Error: "unable to create booking"}
	}
	var isBooked bool
	err = tx.QueryRow(`SELECT is_booked FROM books WHERE id = $1`, req.BookId).
		Scan(&isBooked)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return types.CreateId{}, 404, types.Err{Error: "book not found"}
		}

		if strings.Contains(err.Error(), "uuid") {
			return types.CreateId{}, 400, types.Err{Error: "invalid id"}
		}

		return types.CreateId{}, 500, types.Err{Error: "unable to create booking"}
	}

	if isBooked {
		tx.Rollback()
		return types.CreateId{}, 409, types.Err{Error: "book is already booked"}
	}

	var bookingId types.CreateId
	err = tx.QueryRow(`INSERT INTO bookings(book_id, customer_name, customer_phone, updated_by) VALUES ($1, $2, $3, $4) RETURNING id`,
		req.BookId, req.CustomerName, req.CustomerPhone, uid).Scan(&bookingId.Id)
	if err != nil {
		tx.Rollback()
		return types.CreateId{}, 500, types.Err{Error: "unable to create booking"}
	}

	_, err = tx.Exec(`UPDATE books SET is_booked = TRUE, booked_until = NOW() + INTERVAL '7 days' WHERE id = $1`,
		req.BookId)
	if err != nil {
		tx.Rollback()
		return types.CreateId{}, 500, types.Err{Error: "unable to create booking"}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return types.CreateId{}, 500, types.Err{Error: "unable to create booking"}
	}

	return bookingId, 201, types.Err{}
}

func (s *PostgresStore) ReturnBook(id *string) (int, types.Err) {
	tx, err := s.db.Begin()
	if err != nil {
		return 500, types.Err{Error: "unable to return book"}
	}

	var bookId string
	var isReturned bool
	err = tx.
		QueryRow(`SELECT bo.book_id, bo.is_returned FROM bookings bo INNER JOIN books b ON bo.book_id = b.id WHERE bo.id = $1 FOR UPDATE`, id).
		Scan(&bookId, &isReturned)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return 404, types.Err{Error: "booking not found"}
		}

		if strings.Contains(err.Error(), "uuid") {
			return 400, types.Err{Error: "invalid id"}
		}

		return 500, types.Err{Error: "unable to return book"}
	}

	if isReturned {
		tx.Rollback()
		return 409, types.Err{Error: "book is already returned"}
	}

	_, err = tx.Exec(`UPDATE bookings SET is_returned = TRUE, returned_at = NOW() WHERE id = $1`, *id)
	if err != nil {
		tx.Rollback()
		return 500, types.Err{Error: "unable to return book"}
	}

	_, err = tx.Exec(`UPDATE books SET is_booked = FALSE, booked_until = NULL WHERE id = $1`, bookId)
	if err != nil {
		tx.Rollback()
		return 500, types.Err{Error: "unable to return book"}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 500, types.Err{Error: "unable to return book"}
	}

	return 200, types.Err{}
}

func (s *PostgresStore) GetBooking(lastId, limit *int) ([]types.GetBooking, int, types.Err) {
	query := `SELECT bo.id, bo.pagination_id, b.id, b.title, b.author, bo.customer_name, bo.customer_phone,
	bo.created_at + INTERVAL '7 days', bo.created_at, bo.updated_at, e.username, COALESCE(bo.returned_at::TEXT, ''), bo.is_returned
	FROM bookings bo
	INNER JOIN books b ON bo.book_id = b.id
	INNER JOIN employees e ON bo.updated_by = e.id
	`
	var args []interface{}
	argsCount := 1

	if *lastId != 0 {
		query += `WHERE bo.pagination_id < $1 `
		args = append(args, lastId)
		argsCount++
	}

	query += `ORDER BY bo.pagination_id DESC `

	if *limit != 0 {
		query += `LIMIT $` + strconv.Itoa(argsCount)
		args = append(args, limit)
		argsCount++
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, 500, types.Err{Error: "unable to get bookings"}
	}
	defer rows.Close()

	var bookings []types.GetBooking
	for rows.Next() {
		var booking types.GetBooking
		err := rows.Scan(&booking.Id, &booking.PaginationId, &booking.BookId, &booking.BookTitle, &booking.BookAuthor,
			&booking.CustomerName, &booking.CustomerPhone, &booking.BookedUntil, &booking.CreatedAt, &booking.UpdatedAt,
			&booking.UpdatedBy, &booking.ReturnedAt, &booking.IsReturned)
		if err != nil {
			return nil, 500, types.Err{Error: "unable to get bookings"}
		}
		bookings = append(bookings, booking)
	}

	if len(bookings) == 0 {
		return nil, 404, types.Err{Error: "no bookings found"}
	}

	return bookings, 200, types.Err{}
}
