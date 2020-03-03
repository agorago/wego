package i18n

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"gitlab.intelligentb.com/devops/bplus/config"
	bplusc "gitlab.intelligentb.com/devops/bplus/context"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *goi18n.Bundle

func init() {
	bundle = goi18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	configPath := config.GetConfigPath() + "/bundles"
	InitConfig(configPath)
}

// InitConfig - allow end users to set alternate config paths.
// By default we use config path defined in an environment variable
func InitConfig(cpath string) {
	filepath.Walk(cpath, func(s string, info os.FileInfo, err error) error {
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

		bundle.MustParseMessageFileBytes(buf, lang+".toml")

		return nil
	})
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

// getLocalizer - look for both lang and accept headers.
// lang overrides accept (if it exists)
// If both headers dont exist then default it to the default language of the bundle
func getLocalizer(ctx context.Context) *goi18n.Localizer {
	lang, oklang := bplusc.Value(ctx, "lang").(string)
	accept, okaccept := bplusc.Value(ctx, "Accept-Language").(string)
	defaultLanguage := config.GetDefaultLanguage()
	if oklang {
		if okaccept {
			return i18n.NewLocalizer(bundle, lang, accept, defaultLanguage)
		}
		return i18n.NewLocalizer(bundle, lang, defaultLanguage)
	}
	if okaccept {
		return i18n.NewLocalizer(bundle, accept, defaultLanguage)
	}
	return i18n.NewLocalizer(bundle, defaultLanguage)
}

// Translate - translate a string into the language as specified in the request
// if lang is not specified then use the default language
func Translate(ctx context.Context, s string, m map[string]interface{}) string {
	localizer := getLocalizer(ctx)

	t, err := localizer.Localize(
		&i18n.LocalizeConfig{
			TemplateData: m,
			MessageID:    s,
		},
	)
	// If you cannot translate s just return s and log an error
	if err != nil {
		fmt.Fprintf(os.Stderr, "i18n: Missing message resource %s\n", s)
		return s
	}
	return t
}
