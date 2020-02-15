package i18n

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/MenaEnergyVentures/bplus/config"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *goi18n.Bundle

func init() {
	bundle = goi18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	configPath := config.GetConfigPath() + "/bundles"
	fmt.Printf("config path is %s\n", configPath)
	filepath.Walk(configPath, func(s string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		ind := strings.Index(s, "bundles/")
		if ind == -1 {
			return nil
		}
		a := substring(s, ind+len("bundles/"))
		i := strings.Index(a, "/")
		if i == -1 {
			return nil // discard if stuff is not available under respective bundle folder
		}
		lang := substring(a, 0, i)
		buf, err := ioutil.ReadFile(s)
		fmt.Printf("Registering %s as lang %s\n", buf, lang)
		bundle.MustParseMessageFileBytes(buf, lang+".toml")
		fmt.Println(s)
		return nil
	})

	// bundle.MustLoadMessageFile(configPath + "/bundles/en")
}

// returns a sub string of s starting from ind[0] and ends with an optional ind[1]. If
// ind[1] is not specified then it goes up to end of string
func substring(s string, ind ...int) string {
	ret := []rune(s)
	if len(ind) == 1 {
		ret = ret[ind[0]:]
	} else {
		ret = ret[ind[0]:ind[1]]
	}
	return string(ret)
}

// Translate - translate a string into the language as specified in the request or if not
// specified to the default language
func Translate(ctx context.Context, s string, m map[string]interface{}) string {
	return ""
}
