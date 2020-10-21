package link_test

import (
	"github.com/ory/viper"

	"github.com/zzpu/ums/driver/configuration"
	"github.com/zzpu/ums/identity"
	"github.com/zzpu/ums/selfservice/flow/recovery"
)

func initViper() {
	viper.Set(configuration.ViperKeyDefaultIdentitySchemaURL, "file://./stub/default.schema.json")
	viper.Set(configuration.ViperKeySelfServiceBrowserDefaultReturnTo, "https://www.ory.sh")
	viper.Set(configuration.ViperKeySelfServiceStrategyConfig+"."+identity.CredentialsTypePassword.String()+".enabled", true)
	viper.Set(configuration.ViperKeySelfServiceStrategyConfig+"."+recovery.StrategyRecoveryLinkName+".enabled", true)
	viper.Set(configuration.ViperKeySelfServiceRecoveryEnabled, true)
	viper.Set(configuration.ViperKeySelfServiceVerificationEnabled, true)
}
