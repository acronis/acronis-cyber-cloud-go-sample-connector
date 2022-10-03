// Copyright (c) 2021 Acronis International GmbH
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// The `storage` package initialize database and provides data access methods
package storage

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/config"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/models"
)

// Create connection based on provided database details
func SetupDB(dbConfig *config.Config) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		dbConfig.DB.Host, dbConfig.DB.Username, dbConfig.DB.Password, dbConfig.DB.Database, dbConfig.DB.Port)

	gormLogger := gormlogger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), gormlogger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      gormlogger.Silent,
		Colorful:      true,
	})
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger})

	if migrateErr := db.AutoMigrate(
		&models.User{},
		&models.Tenant{},
		&models.OfferingItem{},
		&models.AccessPolicy{},
		&models.Usage{}); migrateErr != nil {
		log.Fatal(migrateErr)
	}

	if err != nil {
		log.Fatal(err)
	}

	config.DBConn = db
}
