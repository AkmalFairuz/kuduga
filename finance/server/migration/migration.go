package migration

import (
	"errors"
	"flag"
	"fmt"
	"github.com/akmalfairuz/finance/module/database"
	"github.com/akmalfairuz/finance/server/migration/v0"
	"github.com/sirupsen/logrus"
	"sort"
	"time"
)

type Migration interface {
	Version() uint
	Description() string
	Up(db *database.DB) error
	Down(db *database.DB) error
}

var migrations = []Migration{
	v0.V1{},
	v0.V2{},
	v0.V3{},
	v0.V4{},
	v0.V5{},
	v0.V6{},
	v0.V7{},
	v0.V8{},
	v0.V9{},
	v0.V10{},
	v0.V11{},
	v0.V12{},
	v0.V13{},
	v0.V14{},
	v0.V15{},
	v0.V16{},
	v0.V17{},
}

const migrationTrackSchema = `
CREATE TABLE migrationTrack (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    fromVersion INT UNSIGNED NOT NULL,
    targetVersion INT UNSIGNED NOT NULL,
    createdAt BIGINT NOT NULL
) ENGINE = InnoDB`

type Track struct {
	ID            int   `db:"id"`
	FromVersion   uint  `db:"fromVersion"`
	TargetVersion uint  `db:"targetVersion"`
	CreatedAt     int64 `db:"createdAt"`
}

func currentVersion(db *database.DB) (uint, error) {
	var lastMigration Track

	exists, err := db.Exists("SHOW TABLES LIKE 'migrationTrack'")
	if err != nil {
		return 0, err
	}
	if !exists {
		_, err := db.Exec(migrationTrackSchema)
		return 0, err
	}

	if err := db.Get(&lastMigration, "SELECT * FROM migrationTrack ORDER BY createdAt DESC LIMIT 1"); err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return lastMigration.TargetVersion, nil
}

func Up(log *logrus.Logger, db *database.DB) error {
	cur, err := currentVersion(db)
	if err != nil {
		return err
	}

	previousVersion := cur

	allowMigration := flag.Bool("migrate", false, "Allow Database Migration")
	flag.Parse()

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version() > migrations[j].Version()
	})

	for _, m := range migrations {
		if previousVersion+1 != m.Version() {
			continue
		}
		if !*allowMigration {
			log.Fatalf("migration is not allowed, try add -migrate flag")
		}
		log.Infof("migrating to v%v...", m.Version())

		if err := m.Up(db); err != nil {
			return fmt.Errorf("failed to migrate v%v: %w", m.Version(), err)
		}

		if _, err := db.NamedExec(
			"INSERT INTO migrationTrack (fromVersion, targetVersion, createdAt) VALUES (:fromVersion, :targetVersion, :createdAt)",
			&Track{
				FromVersion:   previousVersion,
				TargetVersion: m.Version(),
				CreatedAt:     time.Now().Unix(),
			}); err != nil {
			return err
		}

		log.Infof("successfully migration to v%v!", m.Version())
		previousVersion = m.Version()
	}
	return nil
}
