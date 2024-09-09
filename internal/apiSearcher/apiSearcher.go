package apisearcher

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"golang.org/x/net/html"

	"github.com/keyopto/kscraper/internal/logger"
	"github.com/keyopto/kscraper/internal/types"
)

type ErrorsAPIStruct struct {
	Address string
	Error   error
}

type apiSearchParamsStruct struct {
	waitGroup                sync.WaitGroup
	mutex                    sync.Mutex
	baseHttpAddress          string
	errorsEncoutered         []ErrorsAPIStruct
	addressesToSearch        []string
	addressesAlreadySearched []string
}

func ApiSearcher(argument types.ArgumentCommand) []ErrorsAPIStruct {
	var searchParams apiSearchParamsStruct

	parsedURL := argument.HttpAddress

	searchParams.baseHttpAddress = getBaseAddressFromUrl(parsedURL)
	searchParams.addressesToSearch = []string{parsedURL.String()}

	for len(searchParams.addressesToSearch) != 0 {
		// We make another list to be able to refill addressesToSearch with async calls
		toSearchList := make([]string, len(searchParams.addressesToSearch))
		copy(toSearchList, searchParams.addressesToSearch)
		searchParams.addressesToSearch = searchParams.addressesToSearch[:0]

		for _, addressToSearch := range toSearchList {
			searchParams.waitGroup.Add(1)
			go httpSearch(&searchParams, addressToSearch)
		}
		searchParams.waitGroup.Wait()
	}

	return searchParams.errorsEncoutered
}

func getBaseAddressFromUrl(url *url.URL) string {
	return url.Scheme + "://" + url.Host
}

func httpSearch(searchParams *apiSearchParamsStruct, httpAddress string) {
	defer searchParams.waitGroup.Done()

	logger.Logger.Info("Test : " + httpAddress)

	response, errRequest := http.Get(httpAddress)
	if errRequest != nil {
		logger.Logger.Debug("Error for Get Request on " + httpAddress)
		searchParams.errorsEncoutered = append(searchParams.errorsEncoutered, ErrorsAPIStruct{httpAddress, errRequest})
		return
	}

	if response.StatusCode >= 300 {
		searchParams.errorsEncoutered = append(searchParams.errorsEncoutered, ErrorsAPIStruct{httpAddress, errors.New(fmt.Sprintf("Error %d", response.StatusCode))})
	}

	// we do not want to go to other websites than the base one
	if !strings.HasPrefix(httpAddress, searchParams.baseHttpAddress) {
		return
	}

	parsedHtml, errParsing := html.Parse(response.Body)

	if errParsing != nil {
		fmt.Println("Error : Could not parse html")
		return
	}

	linksInPage := searchLinkAndAddItToWhatNeedsToBeSearched(searchParams, parsedHtml)

	searchParams.mutex.Lock()
	addAddressesSearched(searchParams, linksInPage)
	searchParams.mutex.Unlock()
}

func addAddressesSearched(searchParams *apiSearchParamsStruct, linksToAdd []string) {
	for _, link := range linksToAdd {
		parsedURL, err := url.Parse(link)
		if err != nil {
			fmt.Println("Error parsing URL:", err)
		}

		baseURL := fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, parsedURL.Path)

		if !contains(searchParams.addressesAlreadySearched, baseURL) {
			searchParams.addressesToSearch = append(searchParams.addressesToSearch, baseURL)
			searchParams.addressesAlreadySearched = append(searchParams.addressesAlreadySearched, baseURL)
		}
	}
}

func searchLinkAndAddItToWhatNeedsToBeSearched(searchParams *apiSearchParamsStruct, n *html.Node) []string {
	var toReturn []string

	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				link := attr.Val
				if strings.HasPrefix(link, "/") {
					link = searchParams.baseHttpAddress + link
				}
				toReturn = append(toReturn, link)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		toReturn = append(toReturn, searchLinkAndAddItToWhatNeedsToBeSearched(searchParams, c)...)
	}

	return toReturn
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
