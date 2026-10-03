package envutils

import (
	"errors"
	"fmt"
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	"os"
)

func Load[T any](envPath string, dst *T) error {
	_, err := os.Stat(envPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if err == nil {
		if err := godotenv.Load(envPath); err != nil {
			return fmt.Errorf("failed to load .env file: %w", err)
		}
	}
	if err := env.Parse(dst); err != nil {
		return fmt.Errorf("error parsing env: %w", err)
	}
	return nil
}
