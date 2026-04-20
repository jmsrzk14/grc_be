package data

import (
	"database/sql"
	"grc_be/internal/conf"
	"grc_be/internal/data/schema"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// Data holds the database connection.
type Data struct {
	db *gorm.DB
}

// NewData membuat koneksi database dan menjalankan migrasi otomatis menggunakan Goose.
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	helper := log.NewHelper(logger)

	// 1. Inisialisasi GORM
	db, err := gorm.Open(postgres.Open(c.Database.Source), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		return nil, nil, err
	}

	// 2. Jalankan Goose Migrations menggunakan SQL driver dari GORM
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}
	
	goose.SetBaseFS(schema.Migrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return nil, nil, err
	}
	
	// Menjalankan semua file .sql di folder internal/data/schema/migrations
	if err := goose.Up(sqlDB, "migrations"); err != nil {
		helper.Errorf("failed to apply migrations: %v", err)
		return nil, nil, err
	}

	// 3. Jalankan Seed Data dari folder seeds
	if err := runSeeds(sqlDB, helper); err != nil {
		helper.Errorf("failed to run seeds: %v", err)
		return nil, nil, err
	}

	helper.Info("database connected, migrations and seeds applied successfully")

	cleanup := func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
		helper.Info("database connection closed")
	}

	return &Data{db: db}, cleanup, nil
}

func runSeeds(db *sql.DB, helper *log.Helper) error {
	entries, err := schema.Migrations.ReadDir("seeds")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := "seeds/" + entry.Name()
		content, err := schema.Migrations.ReadFile(path)
		if err != nil {
			return err
		}

		helper.Infof("running seed file: %s", entry.Name())
		if _, err := db.Exec(string(content)); err != nil {
			return err
		}
	}
	return nil
}
