module github.com/zzpu/ums

go 1.14

// See https://github.com/markbates/pkger/pull/112
replace github.com/markbates/pkger => github.com/falafeljan/pkger v0.17.1-0.20200722132747-95726f5b9b9b

require (
	github.com/Masterminds/sprig/v3 v3.0.0
	github.com/alecthomas/units v0.0.0-20190924025748-f65c72e2690d // indirect
	github.com/arbovm/levenshtein v0.0.0-20160628152529-48b4e1c0c4d0
	github.com/armon/go-metrics v0.3.3 // indirect
	github.com/bxcodec/faker/v3 v3.3.1
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/davidrjonas/semver-cli v0.0.0-20190116233701-ee19a9a0dda6
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/ghodss/yaml v1.0.0
	github.com/go-errors/errors v1.0.1
	github.com/go-openapi/errors v0.19.6
	github.com/go-openapi/runtime v0.19.20
	github.com/go-openapi/strfmt v0.19.5
	github.com/go-openapi/swag v0.19.9
	github.com/go-openapi/validate v0.19.10
	github.com/go-playground/locales v0.12.1 // indirect
	github.com/go-playground/universal-translator v0.16.0 // indirect
	github.com/go-swagger/go-swagger v0.25.0
	github.com/gobuffalo/fizz v1.13.1-0.20200903094245-046abeb7de46
	github.com/gobuffalo/httptest v1.0.2
	github.com/gobuffalo/packr/v2 v2.8.0 // indirect
	github.com/gobuffalo/pop/v5 v5.2.4
	github.com/gobuffalo/uuid v2.0.5+incompatible
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/golang/gddo v0.0.0-20190904175337-72a348e765d2
	github.com/golang/mock v1.3.1
	github.com/google/go-github/v27 v27.0.1
	github.com/google/go-jsonnet v0.15.1-0.20200415122941-8a0084e64395
	github.com/google/uuid v1.1.1
	github.com/gorilla/context v1.1.1
	github.com/gorilla/sessions v1.1.3
	github.com/hashicorp/consul/api v1.5.0
	github.com/hashicorp/go-hclog v0.14.1 // indirect
	github.com/hashicorp/go-immutable-radix v1.2.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4
	github.com/hashicorp/serf v0.9.2 // indirect
	github.com/imdario/mergo v0.3.7
	github.com/jteeuwen/go-bindata v3.0.7+incompatible
	github.com/julienschmidt/httprouter v1.2.0
	github.com/justinas/nosurf v1.1.1
	github.com/leodido/go-urn v1.1.0 // indirect
	github.com/markbates/pkger v0.17.0
	github.com/mattn/goveralls v0.0.5
	github.com/mikefarah/yq v1.15.0
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7 // indirect
	github.com/ory/analytics-go/v4 v4.0.0
	github.com/ory/cli v0.0.28-0.20200903094455-e203529bf12b
	github.com/ory/dockertest v3.3.5+incompatible
	github.com/ory/dockertest/v3 v3.5.4
	github.com/ory/go-acc v0.1.0
	github.com/ory/go-convenience v0.1.0
	github.com/ory/graceful v0.1.1
	github.com/ory/herodot v0.9.0
	github.com/ory/jsonschema/v3 v3.0.1
	github.com/ory/mail/v3 v3.0.0
	github.com/ory/viper v1.7.5
	github.com/ory/x v0.0.148
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.4.0
	github.com/prometheus/common v0.9.1
	github.com/shurcooL/go v0.0.0-20180423040247-9e1955d9fb6e
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/sqs/goreturns v0.0.0-20181028201513-538ac6014518
	github.com/stretchr/testify v1.6.1
	github.com/tidwall/gjson v1.3.5
	github.com/tidwall/sjson v1.0.4
	github.com/urfave/negroni v1.0.0
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/tools v0.0.0-20200717024301-6ddee64345a6
	gopkg.in/go-playground/validator.v9 v9.28.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)
