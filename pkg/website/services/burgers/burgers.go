package burgers

import (
	"context"
	"errors"
	"log"

	"github.com/banch0/mux/pkg/website/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

// ErrPoolCantBeNil ...
var ErrPoolCantBeNil = errors.New("pool can't be nil")

// ServiceBurgers ...
type ServiceBurgers struct {
	pool *pgxpool.Pool
}

// NewBurgersSvc ...
func NewBurgersSvc(pool *pgxpool.Pool) *ServiceBurgers {
	if pool == nil {
		panic(ErrPoolCantBeNil)
	}
	return &ServiceBurgers{pool: pool}
}

// BurgersDelete ...
func (service *ServiceBurgers) BurgersDelete(ctx context.Context, id int) error {
	conn, err := service.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	comTag, err := conn.Exec(ctx, "DELETE FROM burgers WHERE id = $1", &id)
	if err != nil {
		log.Println("ERROR: ", err)
	}
	log.Println(comTag)
	return err
}

// BurgersList ...
func (service *ServiceBurgers) BurgersList(ctx context.Context) (list []models.Burger, err error) {
	list = make([]models.Burger, 0)
	conn, err := service.pool.Acquire(context.Background())
	if err != nil {
		return nil, err //errors.Wrap("Cant' ", err) // TODO: wrap to specific error
	}
	defer conn.Release()
	rows, err := conn.Query(context.Background(),
		"SELECT id, name, price FROM burgers WHERE removed = FALSE")
	if err != nil {
		// if err == pgx.ErrNoRows {
		// 	return nil, errors.New("no rows")
		// }
		log.Println("errors", err)
		return nil, err // TODO: wrap to specific error
	}
	defer rows.Close()

	for rows.Next() {
		item := models.Burger{}
		err := rows.Scan(&item.ID, &item.Name, &item.Price)
		if err != nil {
			return nil, err // TODO: wrap to specific error
		}
		list = append(list, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return list, nil
}

// SaveBurger ...
func (service *ServiceBurgers) SaveBurger(ctx context.Context, model models.Burger) (err error) {
	conn, err := service.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	res, err := conn.Exec(ctx,
		"INSERT INTO burgers (id, name, price, removed) VALUES ($1, $2, $3, $4)",
		&model.ID, &model.Name, &model.Price, &model.Removed)
	rows := res.RowsAffected()
	if rows == 0 {
		return
	}
	log.Println(res)
	return err
}

// GetBurgerByID ...
func (service *ServiceBurgers) GetBurgerByID(ctx context.Context, id int) (model *models.Burger, err error) {
	conn, err := service.pool.Acquire(ctx)
	if err != nil {
		log.Println("Error GetByID here: ", err)
		return nil, err
	}
	defer conn.Release()
	model = new(models.Burger)
	row := conn.QueryRow(ctx,
		"SELECT id, name, price FROM burgers WHERE id = $1;", &id)
	row.Scan(&model.ID, &model.Name, &model.Price)
	log.Println("Row: ", model)
	return model, err
}

// UpdateByID ...
func (service *ServiceBurgers) UpdateByID(ctx context.Context, model models.Burger, id int) (err error) {
	conn, err := service.pool.Acquire(ctx)
	if err != nil {
		log.Println("Error Update here: ", err)
		return err
	}
	defer conn.Release()

	comTag, err := conn.Exec(ctx,
		"UPDATE burgers SET name = $1, price = $2, removed = $3 WHERE id = $4;",
		&model.Name, &model.Price, &model.Removed, &id)
	if err != nil {
		log.Println("Errors :", err)
	}
	log.Println(comTag)
	return err
}
