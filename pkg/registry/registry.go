package registry

import (
	"fmt"

	t "github.com/audibleblink/kh/pkg/types"
	"github.com/spf13/viper"
)

// Registry holds the Unmarshaled YAML configs where the CLI can dynamically choose which
// service to validate against based on user input.
var Registry = make(map[string]*t.KeyHack)

func Build() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("keyhacks")

	viper.AddConfigPath("/etc/")
	viper.AddConfigPath("$HOME/.config/")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	viper.Unmarshal(&Registry)
}
