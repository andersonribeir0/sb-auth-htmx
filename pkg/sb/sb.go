package sb

import (
	"errors"
	"os"

	"github.com/nedpals/supabase-go"
)

var (
	Client   *supabase.Client
	sbHost   string
	sbSecret string
)

func Init() error {
	sbHost = os.Getenv("SUPABASE_URL")
	if sbHost == "" {
		return errors.New("supabase host is required")
	}
	sbSecret = os.Getenv("SUPABASE_SECRET")
	if sbSecret == "" {
		return errors.New("supabase secret is required")
	}

	Client = supabase.CreateClient(sbHost, sbSecret)

	return nil
}
