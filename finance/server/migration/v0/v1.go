package v0

import (
	"fmt"
	"github.com/akmalfairuz/finance/module/database"
)

type V1 struct{}

func (V1) Version() uint {
	return 1
}

func (v V1) Description() string {
	return "Initial Migration"
}

func (v V1) Up(db *database.DB) error {
	if _, err := db.Exec(`
CREATE TABLE users (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    displayName VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    passwordHash BINARY(60) NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0,
    pinHash VARCHAR(255),
    kycStatus TINYINT DEFAULT 0,
    role TINYINT DEFAULT 0,
    createdAt BIGINT NOT NULL
) ENGINE = InnoDB`); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	if _, err := db.Exec(`
CREATE TABLE transactions (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    userId BIGINT NOT NULL,
    type TINYINT NOT NULL,
    extraData VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    beforeBalance BIGINT NOT NULL,
    afterBalance BIGINT NOT NULL,
    amount BIGINT NOT NULL,
    createdAt BIGINT NOT NULL,
    FOREIGN KEY (userId) REFERENCES users (id)
) ENGINE = InnoDB`); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	if _, err := db.Exec(`
CREATE TABLE otp (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    type TINYINT NOT NULL,
    scope TINYINT NOT NULL,
    contact VARCHAR(255) NOT NULL,
    code VARCHAR(255) NOT NULL,
    token VARCHAR(32) NOT NULL,
    failAttempt BIGINT NOT NULL DEFAULT 0,
    retryAttempt BIGINT NOT NULL DEFAULT 0,
    lastRetryAt BIGINT NOT NULL DEFAULT 0,
    usedAt BIGINT,
    createdAt BIGINT NOT NULL,
    expiredAt BIGINT NOT NULL
) ENGINE = InnoDB`); err != nil {
		return fmt.Errorf("failed to create otp table: %w", err)
	}

	if _, err := db.Exec(`
CREATE TABLE tokens (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    value CHAR(32) NOT NULL,
    role TINYINT NOT NULL,
    userId BIGINT NOT NULL,
    deviceFcmToken VARCHAR(255),
    createdAt BIGINT NOT NULL,
    expiredAt BIGINT NOT NULL,
    FOREIGN KEY (userId) REFERENCES users (id)
) ENGINE = InnoDB`); err != nil {
		return fmt.Errorf("failed to create tokens table: %w", err)
	}

	if _, err := db.Exec(`
CREATE TABLE resetPasswords (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    userId BIGINT NOT NULL,
    token CHAR(32) NOT NULL,
    resetsAt BIGINT,
    createdAt BIGINT NOT NULL,
    expiredAt BIGINT NOT NULL,
    FOREIGN KEY (userId) REFERENCES users (id)
) ENGINE = InnoDB`); err != nil {
		return fmt.Errorf("failed to create resetPasswords table: %w", err)
	}

	if _, err := db.Exec(`
CREATE TABLE depositRequests (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    userId BIGINT NOT NULL,
    transactionId BIGINT,
    paymentMethod TINYINT NOT NULL,
    externalPaymentId VARCHAR(255) NOT NULL,
    externalPaymentData TEXT NOT NULL,
    amount BIGINT NOT NULL,
    fee BIGINT NOT NULL,
    totalAmount BIGINT NOT NULL,
    description TEXT NOT NULL,
    status TINYINT NOT NULL,
    createdAt BIGINT NOT NULL,
    expiredAt BIGINT NOT NULL,
    INDEX (externalPaymentId),
    FOREIGN KEY (transactionId) REFERENCES transactions (id),
    FOREIGN KEY (userId) REFERENCES users (id)
) ENGINE = InnoDB`); err != nil {
		return fmt.Errorf("failed to create depositRequests table: %w", err)
	}

	if _, err := db.Exec(`
CREATE TABLE depositRequestTracks (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    depositRequestId BIGINT NOT NULL,
    description TEXT NOT NULL,
    newStatus TINYINT,
    createdAt BIGINT NOT NULL,
    FOREIGN KEY (depositRequestId) REFERENCES depositRequests (id)
) ENGINE = InnoDB`); err != nil {
		return fmt.Errorf("failed to create depositRequestTracks table: %w", err)
	}

	if _, err := db.Exec(`
CREATE TABLE withdrawRequests (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    userId BIGINT NOT NULL,
    transactionId BIGINT,
    amount BIGINT NOT NULL,
    fee BIGINT NOT NULL,
    status TINYINT NOT NULL,
    bankType TINYINT NOT NULL,
    bankAccountNumber VARCHAR(255) NOT NULL,
    createdAt BIGINT NOT NULL,
    FOREIGN KEY (userId) REFERENCES users (id)
) ENGINE = InnoDB`); err != nil {
		return fmt.Errorf("failed to create withdrawRequests table: %w", err)
	}

	if _, err := db.Exec(`
CREATE TABLE withdrawRequestTracks (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    withdrawRequestId BIGINT NOT NULL,
    description TEXT NOT NULL,
    newStatus TINYINT,
    createdAt BIGINT NOT NULL,
    FOREIGN KEY (withdrawRequestId) REFERENCES withdrawRequests (id)
) ENGINE = InnoDB`); err != nil {
		return fmt.Errorf("failed to create withdrawRequestTracks table: %w", err)
	}

	if _, err := db.Exec(`
CREATE TABLE productCategories (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    parentId BIGINT,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    productType TINYINT NOT NULL,
    dedicated BOOLEAN NOT NULL,
    iconUrl VARCHAR(255),
FOREIGN KEY (parentId) REFERENCES productCategories (id) ON DELETE CASCADE
) ENGINE = InnoDB`); err != nil {
		return fmt.Errorf("failed to create productCategories table: %w", err)
	}

	if _, err := db.Exec(`
CREATE TABLE productDestinations (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) DEFAULT '',
    format VARCHAR(255) DEFAULT '',
    checkerId VARCHAR(255) DEFAULT '',
    fields TEXT NOT NULL
) ENGINE = InnoDB`); err != nil {
		return fmt.Errorf("failed to create productDestinations table: %w", err)
	}

	if _, err := db.Exec(`
CREATE TABLE products (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    categoryId BIGINT NOT NULL,
    sku VARCHAR(255) NOT NULL,
    externalSku VARCHAR(255) NOT NULL,
    provider VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    type TINYINT NOT NULL,
    imageUrl VARCHAR(255) NOT NULL,
    cutOffStart VARCHAR(255) NOT NULL,
    cutOffEnd VARCHAR(255) NOT NULL,
    destinationType BIGINT NOT NULL,
    isAvailable TINYINT NOT NULL,
    price BIGINT NOT NULL,
    maxWholesalePrice BIGINT NOT NULL,
    wholesalePrice BIGINT NOT NULL,
    beforeDiscountPrice BIGINT NOT NULL,
    createdAt BIGINT NOT NULL,
    updatedAt BIGINT NOT NULL,
    FOREIGN KEY (categoryId) REFERENCES productCategories (id),
    FOREIGN KEY (destinationType) REFERENCES productDestinations (id),
    INDEX (sku)
) ENGINE = InnoDB`); err != nil {
		return fmt.Errorf("failed to create products table: %w", err)
	}

	if _, err := db.Exec(`
CREATE TABLE purchases (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    userId BIGINT,
    refId VARCHAR(255) NOT NULL,
    contactEmail VARCHAR(255) NOT NULL,
    productId BIGINT NOT NULL,
    productName VARCHAR(255) NOT NULL,
    productType TINYINT NOT NULL,
    productDescription TEXT NOT NULL,
    productSku VARCHAR(255) NOT NULL,
    productCategoryId BIGINT NOT NULL,
    productCategoryName VARCHAR(255) NOT NULL,
    providerInfo TEXT NOT NULL,
    price BIGINT NOT NULL,
    totalBill BIGINT NOT NULL DEFAULT 0,
    totalBillFee BIGINT NOT NULL DEFAULT 0,
    fee BIGINT NOT NULL,
    wholesalePrice BIGINT NOT NULL,
    destinationId BIGINT NOT NULL,
    destination TEXT NOT NULL,
    readableDestination TEXT NOT NULL,
    detailedDestination TEXT NOT NULL,
    proof TEXT,
    paymentMethod TINYINT,
    paymentRefId VARCHAR(255) NOT NULL,
    paymentData TEXT NOT NULL,
    paymentExpiredAt BIGINT NOT NULL,
    refunded BOOLEAN NOT NULL DEFAULT false,
    extraData TEXT NOT NULL,
    billData TEXT NOT NULL,
    internalData TEXT NOT NULL,
    status TINYINT NOT NULL,
    createdAt BIGINT NOT NULL,
    UNIQUE (refId),
    INDEX (paymentRefId),
    FOREIGN KEY (userId) REFERENCES users (id)
) ENGINE = InnoDB`); err != nil {
		return fmt.Errorf("failed to create purchases table: %w", err)
	}

	if _, err := db.Exec(`
CREATE TABLE billPrePurchases (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY ,
    userId BIGINT,
    refId VARCHAR(255) NOT NULL,
    contactEmail VARCHAR(255) NOT NULL,
    productId BIGINT NOT NULL,
    destination TEXT NOT NULL,
    destinationId BIGINT NOT NULL,
    detailedDestination TEXT NOT NULL,
    readableDestination TEXT NOT NULL,
    billAmount BIGINT NOT NULL,
    price BIGINT NOT NULL,
    wholesalePrice BIGINT NOT NULL,
    billData TEXT NOT NULL,
    extraData TEXT NOT NULL,
    createdAt BIGINT NOT NULL,
    expiredAt BIGINT NOT NULL,
    UNIQUE (refId),
    FOREIGN KEY (userId) REFERENCES users (id),
    FOREIGN KEY (productId) REFERENCES products (id) ON DELETE NO ACTION
) ENGINE = InnoDB;`); err != nil {
		return fmt.Errorf("failed to create billPrePurchases table")
	}

	if _, err := db.Exec(`
CREATE TABLE purchaseTracks (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    purchaseId BIGINT NOT NULL,
    description TEXT NOT NULL,
    newStatus TINYINT,
    createdAt BIGINT NOT NULL,
    FOREIGN KEY (purchaseId) REFERENCES purchases (id)
) ENGINE = InnoDB`); err != nil {
		return fmt.Errorf("failed to create purchaseTracks table: %w", err)
	}

	if _, err := db.Exec(`
CREATE TABLE pulsaCategories (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    prefix VARCHAR(255) NOT NULL,
    categoryId BIGINT NOT NULL
) ENGINE = InnoDB`); err != nil {
		return fmt.Errorf("failed to create pulsaCategories table: %w", err)
	}

	if _, err := db.Exec(`
CREATE TABLE meta (
    name VARCHAR(255) NOT NULL PRIMARY KEY,
    value MEDIUMTEXT NOT NULL
) ENGINE = InnoDB`); err != nil {
		return fmt.Errorf("failed to create meta table: %w", err)
	}

	if _, err := db.Exec(`
CREATE TABLE notifications (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    userId BIGINT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    hasRead BOOLEAN DEFAULT false,
    type BIGINT NOT NULL,
    data TEXT NOT NULL,
    createdAt BIGINT NOT NULL,
    FOREIGN KEY (userId) REFERENCES users(id)
) ENGINE = InnoDB`); err != nil {
		return fmt.Errorf("failed to create notifications table: %w", err)
	}

	if _, err := db.Exec(`
CREATE TABLE transferToUser (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    senderUserId BIGINT NOT NULL,
    senderUsername VARCHAR(255) NOT NULL,
    senderName VARCHAR(255) NOT NULL,
    senderEmail VARCHAR(255) NOT NULL,
    senderTransactionId BIGINT NOT NULL,
    receiverUserId BIGINT NOT NULL,
    receiverUsername VARCHAR(255) NOT NULL,
    receiverName VARCHAR(255) NOT NULL,
    receiverEmail VARCHAR(255) NOT NULL,
    receiverTransactionId BIGINT NOT NULL,
    amount BIGINT NOT NULL,
    purpose TINYINT NOT NULL,
    note VARCHAR(255),
    createdAt BIGINT NOT NULL,
    FOREIGN KEY (senderUserId) REFERENCES users(id),
    FOREIGN KEY (receiverUserId) REFERENCES users(id),
    FOREIGN KEY (senderTransactionId) REFERENCES transactions(id),
    FOREIGN KEY (receiverTransactionId) REFERENCES transactions(id)
) ENGINE = InnoDB`); err != nil {
		return fmt.Errorf("failed to create transferToUser table: %w", err)
	}

	if _, err := db.Exec(`
CREATE TABLE vouchers (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    userId BIGINT NOT NULL,
    type VARCHAR(255) NOT NULL,
    expiredAt BIGINT NOT NULL,
    createdAt BIGINT NOT NULL
) ENGINE = InnoDB`); err != nil {
		return fmt.Errorf("failed to create vouchers table: %w", err)
	}

	if _, err := db.Exec(`
CREATE TABLE supportTickets (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    userId BIGINT NOT NULL,
    category VARCHAR(255) NOT NULL,
    status TINYINT DEFAULT 0,
    createdAt BIGINT NOT NULL,
    updatedAt BIGINT NOT NULL,
    FOREIGN KEY (userId) REFERENCES users(id)
) ENGINE = InnoDB`); err != nil {
		return fmt.Errorf("failed to create supportTickets table: %w", err)
	}

	if _, err := db.Exec(`
CREATE TABLE supportTicketMessages (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ticketId BIGINT NOT NULL,
    author VARCHAR(255) NOT NULL,
    role TINYINT NOT NULL,
    message TEXT NOT NULL,
    attachments TEXT NOT NULL,
    createdAt BIGINT NOT NULL,
    FOREIGN KEY (ticketId) REFERENCES supportTickets(id)
) ENGINE = InnoDB`); err != nil {
		return fmt.Errorf("failed to create supportTicketMessages table: %w", err)
	}

	if _, err := db.Exec(`
CREATE TABLE userKyc (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    userId BIGINT NOT NULL,
    documentId VARCHAR(255) NOT NULL,
    fullName VARCHAR(255) NOT NULL,
    documentFileId VARCHAR(255) NOT NULL,
    status TINYINT NOT NULL DEFAULT 0,
    createdAt BIGINT NOT NULL
) ENGINE = InnoDB`); err != nil {
		return err
	}

	return nil
}

func (v V1) Down(db *database.DB) error {
	// TODO
	return nil
}
