package texts_test

import (
	"fmt"

	"github.com/ChristianSiegert/go-packages/texts"
)

func ExampleTruncate() {
	fmt.Println(texts.Truncate("Hello world", 10, "…", false))
	fmt.Println(texts.Truncate("Hello world", 10, "…", true))
	// Output:
	// Hello …
	// Hello wor…
}
