package courier

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/ory/herodot"

	gomail "github.com/ory/mail/v3"

	"github.com/zzpu/ums/driver/configuration"
	"github.com/zzpu/ums/x"
)

type (
	smtpDependencies interface {
		PersistenceProvider
		x.LoggingProvider
	}
	Courier struct {
		Dialer *gomail.Dialer
		d      smtpDependencies
		c      configuration.Provider
		// graceful shutdown handling
		ctx      context.Context
		shutdown context.CancelFunc
	}
	Provider interface {
		Courier() *Courier
	}
)

func NewSMTP(d smtpDependencies, c configuration.Provider) *Courier {
	uri := c.CourierSMTPURL()
	password, _ := uri.User.Password()
	port, _ := strconv.ParseInt(uri.Port(), 10, 64)
	ctx, cancel := context.WithCancel(context.Background())

	var ssl bool
	var tlsConfig *tls.Config
	if uri.Scheme == "smtps" {
		ssl = true
		sslSkipVerify, _ := strconv.ParseBool(uri.Query().Get("skip_ssl_verify"))
		// #nosec G402 This is ok (and required!) because it is configurable and disabled by default.
		tlsConfig = &tls.Config{InsecureSkipVerify: sslSkipVerify, ServerName: uri.Hostname()}
	}

	return &Courier{
		d:        d,
		c:        c,
		ctx:      ctx,
		shutdown: cancel,
		Dialer: &gomail.Dialer{
			/* #nosec we need to support SMTP servers without TLS */
			TLSConfig:    tlsConfig,
			Host:         uri.Hostname(),
			Port:         int(port),
			Username:     uri.User.Username(),
			Password:     password,
			SSL:          ssl,
			Timeout:      time.Second * 10,
			RetryFailure: true,
		},
	}
}

func (m *Courier) QueueEmail(ctx context.Context, t EmailTemplate) (uuid.UUID, error) {
	body, err := t.EmailBody()
	if err != nil {
		return uuid.Nil, err
	}

	subject, err := t.EmailSubject()
	if err != nil {
		return uuid.Nil, err
	}

	recipient, err := t.EmailRecipient()
	if err != nil {
		return uuid.Nil, err
	}

	message := &Message{
		Status:    MessageStatusQueued,
		Type:      MessageTypeEmail,
		Body:      body,
		Subject:   subject,
		Recipient: recipient,
	}
	if err := m.d.CourierPersister().AddMessage(ctx, message); err != nil {
		return uuid.Nil, err
	}
	return message.ID, nil
}

func (m *Courier) Work() error {
	errChan := make(chan error)
	defer close(errChan)

	go m.watchMessages(m.ctx, errChan)

	select {
	case <-m.ctx.Done():
		if errors.Is(m.ctx.Err(), context.Canceled) {
			return nil
		}
		return m.ctx.Err()
	case err := <-errChan:
		return err
	}
}

func (m *Courier) Shutdown(ctx context.Context) error {
	m.shutdown()
	return nil
}

func (m *Courier) watchMessages(ctx context.Context, errChan chan error) {
	for {
		if err := backoff.Retry(func() error {
			if len(m.Dialer.Host) == 0 {
				return errors.WithStack(herodot.ErrInternalServerError.WithReasonf("Courier tried to deliver an email but courier.smtp_url is not set!"))
			}

			messages, err := m.d.CourierPersister().NextMessages(ctx, 10)
			if err != nil {
				if errors.Is(err, ErrQueueEmpty) {
					return nil
				}
				return err
			}
			for k := range messages {
				var msg = messages[k]

				switch msg.Type {
				case MessageTypeEmail:
					from := m.c.CourierSMTPFrom()
					gm := gomail.NewMessage()
					gm.SetHeader("From", from)
					gm.SetHeader("To", msg.Recipient)
					gm.SetHeader("Subject", msg.Subject)
					gm.SetBody("text/plain", msg.Body)
					gm.AddAlternative("text/html", msg.Body)

					if err := m.Dialer.DialAndSend(ctx, gm); err != nil {
						m.d.Logger().
							WithError(err).
							WithField("smtp_server", fmt.Sprintf("%s:%d", m.Dialer.Host, m.Dialer.Port)).
							WithField("smtp_ssl_enabled", m.Dialer.SSL).
							// WithField("email_to", msg.Recipient).
							WithField("message_from", from).
							Error("Unable to send email using SMTP connection.")
						continue
					}

					if err := m.d.CourierPersister().SetMessageStatus(ctx, msg.ID, MessageStatusSent); err != nil {
						m.d.Logger().
							WithError(err).
							WithField("message_id", msg.ID).
							Error(`Unable to set the message status to "sent".`)
						return err
					}

					m.d.Logger().
						WithField("message_id", msg.ID).
						WithField("message_type", msg.Type).
						WithField("message_subject", msg.Subject).
						Debug("Courier sent out message.")
				default:
					return errors.Errorf("received unexpected message type: %d", msg.Type)
				}
			}

			return nil
		}, backoff.NewExponentialBackOff()); err != nil {
			errChan <- err
			return
		}
		time.Sleep(time.Second)
	}
}
