CREATE TABLE employees(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT NOT NULL UNIQUE,
    password BYTEA NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE books(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pagination_id BIGSERIAL,
    title TEXT NOT NULL UNIQUE,
    author TEXT NOT NULL,
    description TEXT NOT NULL,
    is_booked BOOLEAN DEFAULT FALSE,
    booked_until TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_books_is_booked ON books(is_booked);
CREATE INDEX idx_books_booked_until ON books(booked_until);

CREATE TABLE bookings(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pagination_id BIGSERIAL,
    book_id UUID NOT NULL,
    customer_name TEXT NOT NULL,
    customer_phone TEXT NOT NULL,
    is_returned BOOLEAN DEFAULT FALSE,
    returned_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    updated_by UUID,
    FOREIGN KEY (book_id) REFERENCES books(id),
    FOREIGN KEY (updated_by) REFERENCES employees(id)
);

CREATE INDEX idx_bookings_is_returned ON bookings(is_returned);