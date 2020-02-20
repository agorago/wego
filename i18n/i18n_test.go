package i18n_test

import (
	"context"
	"fmt"

	bplusc "github.com/MenaEnergyVentures/bplus/context"
	"github.com/MenaEnergyVentures/bplus/i18n"
)

func init() {
	i18n.InitConfig("test-configs")
}

func ExampleEnglish() {
	ctx := context.Background()
	ctx = bplusc.Add(ctx, "lang", "en-US")
	ctx = bplusc.Add(ctx, "Accept-Language", "en-US")

	fmt.Println(i18n.Translate(ctx, "good_morning", nil))
	fmt.Println(i18n.Translate(ctx, "good_afternoon", nil))
	fmt.Println(i18n.Translate(ctx, "good_evening", nil))
	// Output:
	// Good Morning
	// Good Afternoon
	// Good Evening
}

func ExampleSpanish() {
	ctx := context.Background()
	ctx = bplusc.Add(ctx, "lang", "es")

	fmt.Println(i18n.Translate(ctx, "good_morning", nil))
	fmt.Println(i18n.Translate(ctx, "good_afternoon", nil))
	fmt.Println(i18n.Translate(ctx, "good_evening", nil))
	// Output:
	// Buenos Días
	// Buenas Tardes
	// Buenas Noches
}

func ExampleSpanishWithParam() {
	ctx := context.Background()
	ctx = bplusc.Add(ctx, "lang", "es")

	fmt.Println(i18n.Translate(ctx, "goodbye", map[string]interface{}{
		"Name": "Gopher",
	}))

	// Output:
	// Adios Gopher
}

func ExampleSpanishWithParamOverride() {
	ctx := context.Background()
	ctx = bplusc.Add(ctx, "lang", "es")
	ctx = bplusc.Add(ctx, "Accept-Language", "en-US")

	fmt.Println(i18n.Translate(ctx, "goodbye", map[string]interface{}{
		"Name": "Gopher",
	}))

	// Output:
	// Adios Gopher
}

func ExampleEnglishWithParamAcceptLanguage() {
	ctx := context.Background()

	ctx = bplusc.Add(ctx, "Accept-Language", "en-US")

	fmt.Println(i18n.Translate(ctx, "goodbye", map[string]interface{}{
		"Name": "Gopher",
	}))

	// Output:
	// Good Bye Gopher
}

func ExampleEnglishWithParamUseDefault() {
	ctx := context.Background()

	fmt.Println(i18n.Translate(ctx, "goodbye", map[string]interface{}{
		"Name": "Gopher",
	}))

	// Output:
	// Good Bye Gopher
}
