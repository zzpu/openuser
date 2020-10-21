package identity

import (
	"database/sql/driver"
	"encoding/json"
	"sync"
	"time"

	"github.com/ory/herodot"
	"github.com/ory/x/sqlxx"

	"github.com/zzpu/ums/driver/configuration"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/zzpu/ums/x"
)

type (
	// Identity represents an ORY Kratos identity
	//
	// An identity can be a real human, a service, an IoT device - everything that
	// can be described as an "actor" in a system.
	//
	// swagger:model identity
	Identity struct {
		l *sync.RWMutex `db:"-" faker:"-"`

		// ID is a unique identifier chosen by you. It can be a URN (e.g. "arn:aws:iam::123456789012"),
		// a stringified integer (e.g. "123456789012"), a uuid (e.g. "9f425a8d-7efc-4768-8f23-7647a74fdf13"). It is up to you
		// to pick a format you'd like. It is discouraged to use a personally identifiable value here, like the username
		// or the email, as this field is immutable.
		//
		// required: true
		ID uuid.UUID `json:"id" faker:"-" db:"id"`

		// Credentials represents all credentials that can be used for authenticating this identity.
		Credentials map[CredentialsType]Credentials `json:"-" faker:"-" db:"-"`

		// SchemaID is the ID of the JSON Schema to be used for validating the identity's traits.
		//
		// required: true
		SchemaID string `json:"schema_id" faker:"-" db:"schema_id"`

		// SchemaURL is the URL of the endpoint where the identity's traits schema can be fetched from.
		//
		// format: url
		// required: true
		SchemaURL string `json:"schema_url" faker:"-" db:"-"`

		// Traits represent an identity's traits. The identity is able to create, modify, and delete traits
		// in a self-service manner. The input will always be validated against the JSON Schema defined
		// in `schema_url`.
		//
		// required: true
		Traits Traits `json:"traits" faker:"-" db:"traits"`

		// VerifiableAddresses contains all the addresses that can be verified by the user.
		//
		// Extensions:
		// ---
		// x-omitempty: true
		// ---
		VerifiableAddresses []VerifiableAddress `json:"verifiable_addresses,omitempty" faker:"-" has_many:"identity_verifiable_addresses" fk_id:"identity_id"`

		// RecoveryAddresses contains all the addresses that can be used to recover an identity.
		//
		// Extensions:
		// ---
		// x-omitempty: true
		// ---
		RecoveryAddresses []RecoveryAddress `json:"recovery_addresses,omitempty" faker:"-" has_many:"identity_recovery_addresses" fk_id:"identity_id"`

		// CredentialsCollection is a helper struct field for gobuffalo.pop.
		CredentialsCollection CredentialsCollection `json:"-" faker:"-" has_many:"identity_credentials" fk_id:"identity_id"`

		// CreatedAt is a helper struct field for gobuffalo.pop.
		CreatedAt time.Time `json:"-" db:"created_at"`

		// UpdatedAt is a helper struct field for gobuffalo.pop.
		UpdatedAt time.Time `json:"-" db:"updated_at"`
	}
	Traits json.RawMessage
)

func (t *Traits) Scan(value interface{}) error {
	return sqlxx.JSONScan(t, value)
}

func (t Traits) Value() (driver.Value, error) {
	return sqlxx.JSONValue(t)
}

func (t *Traits) String() string {
	return string(*t)
}

// MarshalJSON returns m as the JSON encoding of m.
func (t Traits) MarshalJSON() ([]byte, error) {
	if t == nil {
		return []byte("null"), nil
	}
	return t, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (t *Traits) UnmarshalJSON(data []byte) error {
	if t == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*t = append((*t)[0:0], data...)
	return nil
}

func (i Identity) TableName() string {
	return "identities"
}

func (i *Identity) lock() *sync.RWMutex {
	if i.l == nil {
		i.l = new(sync.RWMutex)
	}
	return i.l
}

func (i *Identity) SetSecurityAnswers(answers map[string]string) {
	i.lock().Lock()
	defer i.lock().Unlock()

}

func (i *Identity) SetCredentials(t CredentialsType, c Credentials) {
	i.lock().Lock()
	defer i.lock().Unlock()
	if i.Credentials == nil {
		i.Credentials = make(map[CredentialsType]Credentials)
	}

	c.Type = t
	i.Credentials[t] = c
}

func (i *Identity) GetCredentials(t CredentialsType) (*Credentials, bool) {
	i.lock().RLock()
	defer i.lock().RUnlock()

	if c, ok := i.Credentials[t]; ok {
		return &c, true
	}

	return nil, false
}

func (i *Identity) ParseCredentials(t CredentialsType, config interface{}) (*Credentials, error) {
	i.lock().RLock()
	defer i.lock().RUnlock()

	if c, ok := i.Credentials[t]; ok {
		if err := json.Unmarshal(c.Config, config); err != nil {
			return nil, errors.WithStack(err)
		}
		return &c, nil
	}

	return nil, herodot.ErrNotFound.WithReasonf("identity does not have credential type %s", t)
}

func (i *Identity) CopyWithoutCredentials() *Identity {
	var ii = *i
	ii.Credentials = nil
	return &ii
}

func NewIdentity(traitsSchemaID string) *Identity {
	if traitsSchemaID == "" {
		traitsSchemaID = configuration.DefaultIdentityTraitsSchemaID
	}

	return &Identity{
		ID:                  x.NewUUID(),
		Credentials:         map[CredentialsType]Credentials{},
		Traits:              Traits("{}"),
		SchemaID:            traitsSchemaID,
		VerifiableAddresses: []VerifiableAddress{},
		l:                   new(sync.RWMutex),
	}
}
