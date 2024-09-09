package main

import (
	"fmt"

	apisearcher "github.com/keyopto/kscraper/internal/apiSearcher"
	argModule "github.com/keyopto/kscraper/internal/arg"
	"github.com/keyopto/kscraper/internal/logger"
	"github.com/keyopto/kscraper/internal/types"
	"github.com/sirupsen/logrus"
)

func main() {
	logger.Logger = *logrus.New()

	logger.Logger.SetLevel(logrus.InfoLevel)

	var arg types.ArgumentCommand
	err := argModule.ParseArgs(&arg)
	if err != nil {
		fmt.Println(err)
		return
	}

	listCouldntFetch := apisearcher.ApiSearcher(arg)

	for _, errorFetch := range listCouldntFetch {
		fmt.Println("Link : \033[34m" + errorFetch.Address + "\033[0m with error \033[31m" + errorFetch.Error.Error() + "\033[0m")
	}
}
