package main

import (
	"fmt"
	"strings"
)

func main() {
	fileWithPath := "/uploads/english_text.txt.comp"
	fmt.Println(getOnlyFileNameWithExtension(fileWithPath))
}

func getOnlyFileNameWithExtension(fileWithPath string) string {

	splitAtDots := strings.Split(fileWithPath, ".")

	correctExtension := splitAtDots[0] + "." + splitAtDots[1]

	splitAtSlashes := strings.Split(correctExtension, "/")

	return splitAtSlashes[2]
}
