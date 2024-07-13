package helpers

import (
	"context"
	"github.com/mailgun/mailgun-go/v4"
	"os"
	"time"
)

func SendEmail(from string, to string, subject string, text string) error {
	mg := mailgun.NewMailgun(os.Getenv("MG_DOMAIN"), os.Getenv("MG_API_KEY"))
	message := mg.NewMessage(from, subject, text, to)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err := mg.Send(ctx, message)
	return err
}
