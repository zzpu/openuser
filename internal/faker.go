package internal

import (
	"math/rand"
	"net/http"
	"reflect"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/pkg/errors"

	"github.com/ory/x/randx"

	"github.com/zzpu/openuser/identity"
	"github.com/zzpu/openuser/selfservice/flow"
	"github.com/zzpu/openuser/selfservice/flow/login"
	"github.com/zzpu/openuser/selfservice/flow/recovery"
	"github.com/zzpu/openuser/selfservice/flow/registration"
	"github.com/zzpu/openuser/selfservice/flow/settings"
	"github.com/zzpu/openuser/selfservice/flow/verification"
	"github.com/zzpu/openuser/selfservice/form"
	"github.com/zzpu/openuser/x"
)

func RegisterFakes() {
	_ = faker.SetRandomMapAndSliceSize(4)

	if err := faker.AddProvider("birthdate", func(v reflect.Value) (interface{}, error) {
		return time.Now().Add(time.Duration(rand.Int())).Round(time.Second).UTC(), nil
	}); err != nil {
		panic(err)
	}

	if err := faker.AddProvider("time_types", func(v reflect.Value) (interface{}, error) {
		es := make([]time.Time, rand.Intn(5))
		for k := range es {
			es[k] = time.Now().Add(time.Duration(rand.Int())).Round(time.Second).UTC()
		}
		return es, nil
	}); err != nil {
		panic(err)
	}

	if err := faker.AddProvider("http_header", func(v reflect.Value) (interface{}, error) {
		headers := http.Header{}
		for i := 0; i <= rand.Intn(5); i++ {
			values := make([]string, rand.Intn(4)+1)
			for k := range values {
				values[k] = randx.MustString(8, randx.AlphaNum)
			}
			headers[randx.MustString(8, randx.AlphaNum)] = values
		}

		return headers, nil
	}); err != nil {
		panic(err)
	}

	if err := faker.AddProvider("http_method", func(v reflect.Value) (interface{}, error) {
		methods := []string{"POST", "PUT", "GET", "PATCH"}
		return methods[rand.Intn(len(methods))], nil
	}); err != nil {
		panic(err)
	}

	if err := faker.AddProvider("identity_credentials_type", func(v reflect.Value) (interface{}, error) {
		methods := []identity.CredentialsType{identity.CredentialsTypePassword, identity.CredentialsTypePassword}
		return string(methods[rand.Intn(len(methods))]), nil
	}); err != nil {
		panic(err)
	}

	if err := faker.AddProvider("string", func(v reflect.Value) (interface{}, error) {
		return randx.MustString(25, randx.AlphaNum), nil
	}); err != nil {
		panic(err)
	}

	if err := faker.AddProvider("time_type", func(v reflect.Value) (interface{}, error) {
		return time.Now().Add(time.Duration(rand.Int())).Round(time.Second).UTC(), nil
	}); err != nil {
		panic(err)
	}

	if err := faker.AddProvider("login_flow_methods", func(v reflect.Value) (interface{}, error) {
		var methods = make(map[identity.CredentialsType]*login.FlowMethod)
		for _, ct := range []identity.CredentialsType{identity.CredentialsTypePassword, identity.CredentialsTypeOIDC} {
			var f form.HTMLForm
			if err := faker.FakeData(&f); err != nil {
				return nil, err
			}
			methods[ct] = &login.FlowMethod{
				Method: ct,
				Config: &login.FlowMethodConfig{FlowMethodConfigurator: &f},
			}

		}
		return methods, nil
	}); err != nil {
		panic(err)
	}

	if err := faker.AddProvider("registration_flow_methods", func(v reflect.Value) (interface{}, error) {
		var methods = make(map[identity.CredentialsType]*registration.FlowMethod)
		for _, ct := range []identity.CredentialsType{identity.CredentialsTypePassword, identity.CredentialsTypeOIDC} {
			var f form.HTMLForm
			if err := faker.FakeData(&f); err != nil {
				return nil, errors.WithStack(err)
			}
			methods[ct] = &registration.FlowMethod{
				Method: ct,
				Config: &registration.FlowMethodConfig{FlowMethodConfigurator: &f},
			}
		}
		return methods, nil
	}); err != nil {
		panic(err)
	}

	if err := faker.AddProvider("settings_flow_methods", func(v reflect.Value) (interface{}, error) {
		var methods = make(map[string]*settings.FlowMethod)
		for _, ct := range []string{settings.StrategyProfile, string(identity.CredentialsTypePassword), string(identity.CredentialsTypeOIDC)} {
			var f form.HTMLForm
			if err := faker.FakeData(&f); err != nil {
				return nil, err
			}
			methods[ct] = &settings.FlowMethod{
				Method: ct,
				Config: &settings.FlowMethodConfig{FlowMethodConfigurator: &f},
			}
		}
		return methods, nil
	}); err != nil {
		panic(err)
	}

	if err := faker.AddProvider("recovery_flow_methods", func(v reflect.Value) (interface{}, error) {
		var methods = make(map[string]*recovery.FlowMethod)
		for _, ct := range []string{recovery.StrategyRecoveryLinkName} {
			var f form.HTMLForm
			if err := faker.FakeData(&f); err != nil {
				return nil, err
			}
			methods[ct] = &recovery.FlowMethod{
				Method: ct,
				Config: &recovery.FlowMethodConfig{FlowMethodConfigurator: &f},
			}
		}
		return methods, nil
	}); err != nil {
		panic(err)
	}

	if err := faker.AddProvider("verification_flow_methods", func(v reflect.Value) (interface{}, error) {
		var methods = make(map[string]*verification.FlowMethod)
		for _, ct := range []string{verification.StrategyVerificationLinkName} {
			var f form.HTMLForm
			if err := faker.FakeData(&f); err != nil {
				return nil, err
			}
			methods[ct] = &verification.FlowMethod{
				Method: ct,
				Config: &verification.FlowMethodConfig{FlowMethodConfigurator: &f},
			}
		}
		return methods, nil
	}); err != nil {
		panic(err)
	}

	if err := faker.AddProvider("uuid", func(v reflect.Value) (interface{}, error) {
		return x.NewUUID(), nil
	}); err != nil {
		panic(err)
	}

	if err := faker.AddProvider("identity", func(v reflect.Value) (interface{}, error) {
		var i identity.Identity
		return &i, faker.FakeData(&i)
	}); err != nil {
		panic(err)
	}

	if err := faker.AddProvider("flow_type", func(v reflect.Value) (interface{}, error) {
		if rand.Intn(2) == 0 {
			return flow.TypeAPI, nil
		}
		return flow.TypeBrowser, nil
	}); err != nil {
		panic(err)
	}
}
