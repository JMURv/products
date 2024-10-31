package dbutils

import (
	"database/sql"
	"errors"
	"fmt"
	conf "github.com/JMURv/par-pro/products/pkg/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"path/filepath"
	"strings"
)

func ApplyMigrations(db *sql.DB, conf *conf.DBConfig) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	path, _ := filepath.Abs("db/migration")
	path = filepath.ToSlash(path)

	m, err := migrate.NewWithDatabaseInstance("file://"+path, conf.Database, driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && errors.Is(err, migrate.ErrNoChange) {
		zap.L().Info("No migrations to apply")
		return nil
	} else if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	zap.L().Info("Applied migrations")
	return nil
}

func FilterItems(q *strings.Builder, args []any, filters map[string]any) []any {
	conds := make([]string, 0, len(filters))

	for key, value := range filters {
		switch key {
		case "min_price":
			conds = append(conds, "i.price >= ?")
			args = append(args, value)
		case "max_price":
			conds = append(conds, "i.price <= ?")
			args = append(args, value)
		default:
			switch v := value.(type) {
			case map[string]any:
				if attrMin, minOK := v["min"].(string); minOK {
					conds = append(conds, "(ia.name = ? AND ia.value >= ?)")
					args = append(args, key, attrMin)
				}

				if attrMax, maxOK := v["max"].(string); maxOK {
					conds = append(conds, "(ia.name = ? AND ia.value <= ?)")
					args = append(args, key, attrMax)
				}

			case []string:
				args = append(args, key)

				placeholders := make([]string, len(v))
				for i := range v {
					placeholders[i] = "?"
					args = append(args, v[i])
				}

				conds = append(
					conds,
					fmt.Sprintf("(ia.name = ? AND ia.value IN (%s))", strings.Join(placeholders, ", ")),
				)
			}
		}
	}

	if len(conds) > 0 {
		q.WriteString(" AND ")
		q.WriteString(strings.Join(conds, " AND "))
	}

	return args
}

//func FilterOrders(query *gorm.DB, filters map[string]any) *gorm.DB {
//	for key, value := range filters {
//		switch key {
//		case "min_price":
//			query = query.Where("total_amount >= ?", value)
//			continue
//		case "max_price":
//			query = query.Where("total_amount <= ?", value)
//			continue
//		}
//
//		switch v := value.(type) {
//		case []string:
//			query = query.Where(key+" IN (?)", v)
//		}
//
//	}
//	return query
//}
