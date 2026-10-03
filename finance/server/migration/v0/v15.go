package v0

import (
	"github.com/akmalfairuz/finance/module/database"
)

type V15 struct{}

func (V15) Version() uint {
	return 15
}

func (v V15) Description() string {
	return "Create analyticalReports table"
}

func (v V15) Up(db *database.DB) error {
	if _, err := db.Exec(`CREATE TABLE analyticalReports (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    profit BIGINT NOT NULL,
    newUserCount BIGINT NOT NULL,
    activeUsers BIGINT NOT NULL,
    totalDeposit BIGINT NOT NULL,
    depositCount BIGINT NOT NULL,
    totalUserBalance BIGINT NOT NULL,
    totalSuccessfulPurchasePrice BIGINT NOT NULL,
    totalSuccessfulPurchaseWholesalePrice BIGINT NOT NULL,
    successfulPurchaseCount BIGINT NOT NULL,
    pendingPurchaseCount BIGINT NOT NULL,
    failedPurchaseCount BIGINT NOT NULL,
    reportDateStart BIGINT NOT NULL,
    reportDateEnd BIGINT NOT NULL,
    createdAt BIGINT NOT NULL
) ENGINE=InnoDB;`); err != nil {
		return err
	}
	return nil
}

func (v V15) Down(db *database.DB) error {
	if _, err := db.Exec("DROP TABLE analyticalReports"); err != nil {
		return err
	}
	return nil
}
