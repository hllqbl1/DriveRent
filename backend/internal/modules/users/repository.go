package users

import (
	"context"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user User,) error {

	query := `
		INSERT INTO users
		(name, email, password, phone, role)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.Name,
		user.Email,
		user.Password,
		user.Phone,
		user.Role,
	)

	if err != nil {
		return err
	}

	return nil
}

func(r *UserRepository) DeleteUser(ctx context.Context, id int64) error {
    query := `DELETE FROM users WHERE id = $1`

    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return sql.ErrNoRows
    }

    return nil
}


func (r *UserRepository) UpdateUser(ctx context.Context, user User) error {
    query := `
        UPDATE users
        SET name = $1,
            email = $2,
            phone = $3
        WHERE id = $4
    `

    result, err := r.db.ExecContext(
        ctx,
        query,
        user.Name,
        user.Email,
        user.Phone,
        user.ID,
    )

    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return sql.ErrNoRows
    }

    return nil
}


func (r *UserRepository) GetAllUsers(ctx context.Context) ([]User, error) {
	query := ` SELECT id, name, email, phone, role, created_at, updated_at FROM users ORDER BY id DESC`

	rows, err := r.db.QueryContext(ctx,query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Phone,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return  nil, err
	}

	return  users, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, id int) (*User, error ) {
	query := `
    SELECT ID, name, email, phone, role, created_at, updated_at
    FROM users
    WHERE id = $1`

    var user User

    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID,
        &user.Name,
        &user.Email,
        &user.Phone,
        &user.Role,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    if err != nil {
        return nil, err
    }

    return &user, nil

}


func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {

    query := `
        SELECT id, name, email, password, phone, role, created_at, updated_at
        FROM users
        WHERE email = $1
    `

    var user User

    err := r.db.QueryRowContext(
        ctx,
        query,
        email,
    ).Scan(
        &user.ID,
        &user.Name,
        &user.Email,
        &user.Password,
        &user.Phone,
        &user.Role,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    if err != nil {
        return nil, err
    }

    return &user, nil
}
