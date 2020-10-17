package identities

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/markbates/pkger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"github.com/zzpu/openuser/internal/httpclient/client/public"

	"github.com/zzpu/openuser/internal/clihelpers"

	"github.com/ory/jsonschema/v3"
	"github.com/ory/x/viperx"
	"github.com/zzpu/openuser/cmd/cliclient"
)

var validateCmd = &cobra.Command{
	Use:   "validate <file.json [file-2.json [file-3.json] ...]>",
	Short: "Validate local identity files",
	Long: `This command allows validation of identity files.
It validates against the payload of the API and the identity schema as configured in Kratos.
Identities can be supplied via STD_IN or JSON files containing a single or an array of identities.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := cliclient.NewClient(cmd)

		is, err := readIdentities(cmd, args)
		if err != nil {
			return err
		}

		for src, i := range is {
			err = validateIdentity(cmd, src, i, c.Public.GetSchema)
			if err != nil {
				return err
			}
		}

		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "All identity files are valid.")
		return nil
	},
}

var schemas = make(map[string]*jsonschema.Schema)

const createIdentityPath = "api.swagger.json#/definitions/CreateIdentity"

type schemaGetter = func(params *public.GetSchemaParams) (*public.GetSchemaOK, error)

// validateIdentity validates the json payload fc against
// 1. the swagger payload definition and
// 2. the remote custom identity schema.
func validateIdentity(cmd *cobra.Command, src, i string, getRemoteSchema schemaGetter) error {
	swaggerSchema, ok := schemas[createIdentityPath]
	if !ok {
		// get swagger schema
		sf, err := pkger.Open("/.schema/api.swagger.json")
		if err != nil {
			return errors.Wrap(err, "Could not open swagger schema. This is an error with the binary you use and should be reported. Thanks ;)")
		}

		// add swagger schema
		schemaCompiler := jsonschema.NewCompiler()
		err = schemaCompiler.AddResource("api.swagger.json", sf)
		if err != nil {
			return errors.Wrap(err, "Could not add swagger schema to the schema compiler. This is an error with the binary you use and should be reported. Thanks ;)")
		}

		// compile swagger payload definition
		swaggerSchema, err = schemaCompiler.Compile(createIdentityPath)
		if err != nil {
			return errors.Wrap(err, "Could not compile the identity schema. This is an error with the binary you use and should be reported. Thanks ;)")
		}
		// force additional properties to false because swagger does not render this
		swaggerSchema.AdditionalProperties = false
		schemas[createIdentityPath] = swaggerSchema
	}

	// validate against swagger definition
	var foundValidationErrors bool
	err := swaggerSchema.Validate(bytes.NewBufferString(i))
	if err != nil {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "%s: not valid\n", src)
		viperx.PrintHumanReadableValidationErrors(cmd.ErrOrStderr(), err)
		foundValidationErrors = true
	}

	// get custom identity schema id
	sid := gjson.Get(i, "schema_id")
	if !sid.Exists() {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), `%s: Expected key "schema_id" to be defined in identity file`, src)
		return clihelpers.FailSilently(cmd)
	}

	customSchema, ok := schemas[sid.String()]
	if !ok {
		// get custom identity schema
		ts, err := getRemoteSchema(&public.GetSchemaParams{ID: sid.String(), Context: context.Background()})
		if err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "%s: Could not fetch schema with ID \"%s\": %s\n", src, sid.String(), err)
			return clihelpers.FailSilently(cmd)
		}
		sf, err := json.Marshal(ts.Payload)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("%s: Could not marshal the traits schema. This usually means there is a problem with your upstream service as it served an invalid response.", src))
		}

		// compile custom identity schema
		customSchema, err = jsonschema.CompileString("identity_traits.schema.json", string(sf))
		if err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "%s: Could not compile the traits schema: %s\n", src, err)
			return clihelpers.FailSilently(cmd)
		}
		schemas[sid.String()] = customSchema
	}

	// validate against custom identity schema
	err = customSchema.Validate(bytes.NewBufferString(i))
	if err != nil {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "%s: not valid\n", src)
		viperx.PrintHumanReadableValidationErrors(cmd.ErrOrStderr(), err)
		foundValidationErrors = true
	}

	if foundValidationErrors {
		return clihelpers.FailSilently(cmd)
	}
	return nil
}
