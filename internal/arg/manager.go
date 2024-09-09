package arg

import (
	"errors"
	"net/url"
	"os"

	"github.com/keyopto/kscraper/internal/types"
)

const HTTP_REGEX = `^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/|\/|\/\/)?[A-z0-9_-]*?[:]?[A-z0-9_-]*?[@]?[A-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?$`

func ParseArgs(arg *types.ArgumentCommand) error {
	if len(os.Args) < 2 {
		return errors.New("Error : you need to pass the http base address in the command")
	}

	if len(os.Args) > 2 {
		return errors.New("Error : too many arguments")
	}

	httpAddress := os.Args[1]

	parsedUrl, err := url.Parse(httpAddress)
	if err != nil {
		return errors.New("Error : The argument is not a valid http address")
	}

	arg.HttpAddress = parsedUrl

	return nil
}
