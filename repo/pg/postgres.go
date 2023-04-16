package pg

import (
	"context"
	"fmt"
	"randomusergen/domain"
	"randomusergen/repo"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type db struct {
	con *pgxpool.Pool
}

func New(url string) (repo.UserRepo, error) {
	ctx := context.Background()

	p, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}

	if err = p.Ping(ctx); err != nil {
		return nil, err
	}

	return &db{
		con: p,
	}, nil
}

func (d *db) Close() error {
	d.con.Close()
	return nil
}

func (d *db) SaveAll(ctx context.Context, users []*domain.UserData) (int, error) {
	if len(users) == 0 {
		return 0, fmt.Errorf("no users to insert")
	}

	tx, err := d.con.Begin(ctx)
	if err != nil {
		return 0, err
	}

	// transaction handling
	err = func() error {
		ids := make([]int, 0, len(users))
		// save users
		b := &pgx.Batch{}
		for _, u := range users {
			b.Queue(
				`INSERT INTO users (
				title, first_name, last_name, gender, email, phone, cell, nat)
				VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
				u.Name.Title,
				u.Name.First,
				u.Name.Last,
				u.Gender,
				u.Email,
				u.Phone,
				u.Cell,
				u.Nat,
			).QueryRow(func(row pgx.Row) error {
				var id int
				err := row.Scan(&id)
				if err != nil {
					return err
				}
				ids = append(ids, id)
				return nil
			})
		}

		err := tx.SendBatch(ctx, b).Close()
		if err != nil {
			return err
		}

		if len(ids) != len(users) {
			return fmt.Errorf("not correct user_ids number returned")
		}

		// save logins
		b = &pgx.Batch{}
		for i, u := range users {
			if u.Login == nil {
				continue
			}
			b.Queue(
				`INSERT INTO logins (
					uuid, username, password, salt, md5, sha1, sha256, user_id)
					VALUES($1, $2, $3, $4, $5, $6, $7, $8)`,
				u.Login.UUID,
				u.Login.Username,
				u.Login.Password,
				u.Login.Salt,
				u.Login.Md5,
				u.Login.Sha1,
				u.Login.Sha256,
				ids[i],
			)
		}

		err = tx.SendBatch(ctx, b).Close()
		if err != nil {
			return err
		}

		// save dob
		b = &pgx.Batch{}
		for i, u := range users {
			if u.Dob == nil {
				continue
			}
			b.Queue(
				`INSERT INTO dob (
					date, age, user_id)
					VALUES($1, $2, $3)`,
				u.Dob.Date,
				u.Dob.Age,
				ids[i],
			)
		}

		err = tx.SendBatch(ctx, b).Close()
		if err != nil {
			return err
		}

		// save registered
		b = &pgx.Batch{}
		for i, u := range users {
			if u.Registered == nil {
				continue
			}
			b.Queue(
				`INSERT INTO registered (
					date, age, user_id)
					VALUES($1, $2, $3)`,
				u.Registered.Date,
				u.Registered.Age,
				ids[i],
			)
		}

		err = tx.SendBatch(ctx, b).Close()
		if err != nil {
			return err
		}

		// save id
		b = &pgx.Batch{}
		for i, u := range users {
			if u.ID == nil {
				continue
			}
			b.Queue(
				`INSERT INTO id (
					name, value, user_id)
					VALUES($1, $2, $3)`,
				u.ID.Name,
				u.ID.Value,
				ids[i],
			)
		}

		err = tx.SendBatch(ctx, b).Close()
		if err != nil {
			return err
		}

		// save picture
		b = &pgx.Batch{}
		for i, u := range users {
			if u.Dob == nil {
				continue
			}
			b.Queue(
				`INSERT INTO picture (
					large, medium, thumbnail, user_id)
					VALUES($1, $2, $3, $4)`,
				u.Picture.Large,
				u.Picture.Medium,
				u.Picture.Thumbnail,
				ids[i],
			)
		}

		err = tx.SendBatch(ctx, b).Close()
		if err != nil {
			return err
		}

		// save locations
		locs := make(map[int]int)

		b = &pgx.Batch{}
		for i, u := range users {
			if u.Location == nil {
				continue
			}
			i := i
			b.Queue(
				`INSERT INTO locations (
				city, state, country, postcode, user_id)
				VALUES($1, $2, $3, $4, $5) RETURNING id`,
				u.Location.City,
				u.Location.State,
				u.Location.Country,
				u.Location.Postcode,
				ids[i],
			).QueryRow(func(row pgx.Row) error {
				var id int
				err := row.Scan(&id)
				if err != nil {
					return err
				}
				locs[i] = id
				return nil
			})
		}

		err = tx.SendBatch(ctx, b).Close()
		if err != nil {
			return err
		}

		if len(locs) == 0 {
			return nil
		}

		// save streets
		b = &pgx.Batch{}
		for uid, lid := range locs {
			uid := uid
			street := users[uid].Location.Street
			if street == nil {
				continue
			}
			b.Queue(
				`INSERT INTO streets (
			number, name, location_id)
			VALUES($1, $2, $3)`,
				street.Number,
				street.Name,
				lid,
			)
		}

		err = tx.SendBatch(ctx, b).Close()
		if err != nil {
			return err
		}

		// save coordinates
		b = &pgx.Batch{}
		for uid, lid := range locs {
			uid := uid
			coordinates := users[uid].Location.Coordinates
			if coordinates == nil {
				continue
			}
			b.Queue(
				`INSERT INTO coordinates (
			latitude, longitude, location_id)
			VALUES($1, $2, $3)`,
				coordinates.Latitude,
				coordinates.Longitude,
				lid,
			)
		}

		err = tx.SendBatch(ctx, b).Close()
		if err != nil {
			return err
		}

		// save timezones
		b = &pgx.Batch{}
		for uid, lid := range locs {
			uid := uid
			timezone := users[uid].Location.Timezone
			if timezone == nil {
				continue
			}
			b.Queue(
				`INSERT INTO timezones (
			offset_time, description, location_id)
			VALUES($1, $2, $3)`,
				timezone.Offset,
				timezone.Description,
				lid,
			)
		}

		err = tx.SendBatch(ctx, b).Close()
		if err != nil {
			return err
		}

		return nil
	}()

	if err != nil {
		e := tx.Rollback(ctx)
		if e != nil {
			fmt.Println("db.rollback ", e)
		}
		return 0, err
	}

	e := tx.Commit(ctx)
	if e != nil {
		fmt.Println("db.commit ", e)
	}

	return len(users), nil
}
