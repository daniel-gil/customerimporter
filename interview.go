// Package customerimporter reads from the given customers.csv file and returns a
// sorted (data structure of your choice) of email domains along with the number
// of customers with e-mail addresses for each domain.  Any errors should be
// logged (or handled). Performance matters (this is only ~3k lines, but *could*
// be 1m lines or run on a small machine).
package customerimporter

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

const defaultCsvFileName = "customers.csv"
const defaultEmailColumnName = "email"
const numExpectedEmailComponents = 2 // the 2 components are username and mail domain: username@domain.com
const emailDomainIndex = 1

// CustomerImporter interface exposes the methods for importing emails
type CustomerImporter interface {

	// SetCsvFileName sets a different CSV filename than the default one.
	SetCsvFileName(fileName string)

	// SetEmailColumnName overwrites the default email column name ('email')
	SetEmailColumnName(columnName string)

	// Import imports the emails from a csv.Reader passed as parameter, group them by email domain and return a list of email domains including a counter
	Import(reader *csv.Reader) ([]EmailGroup, error)

	// ImportFile imports the emails from the file 'customers.csv', group them by email domain and return a list of email domains including a counter
	ImportFile() ([]EmailGroup, error)
}

type customerImporter struct {
	CsvFileName     string
	EmailColumnName string
}

// New creates a new instance of customerImporter returning the interface CustomerImporter
func New() CustomerImporter {
	return &customerImporter{
		CsvFileName:     defaultCsvFileName,
		EmailColumnName: defaultEmailColumnName,
	}
}

func (ci *customerImporter) SetCsvFileName(fileName string) {
	ci.CsvFileName = fileName
}

func (ci *customerImporter) SetEmailColumnName(columnName string) {
	ci.EmailColumnName = columnName
}

func (ci *customerImporter) Import(reader *csv.Reader) ([]EmailGroup, error) {
	// extract the column names from the first line of the file
	emailColIndex, err := ci.parseHeader(reader)
	if err != nil {
		return nil, err
	}

	// lineNumber is set to the value 2 because we already have processed the first line (header)
	lineNumber := 2

	// initialize a map with the email domain as key and a counter as value
	custGrp := make(map[string]int)

	// read the file by lines
	for {
		line, err := reader.Read()
		if err == io.EOF {
			// here we have arrived to the end of file (EOF)
			break
		} else if err != nil {
			// here we found an error while reading from the file
			return nil, err
		}
		err = ci.processLine(line, emailColIndex, custGrp, lineNumber)
		if err != nil {
			// errors found, log the problem and skip this line
			log.Printf("skip line %d: %v", lineNumber, err)
		}
		lineNumber++
	}

	// sort the email domains list by name
	results := ci.sortResults(custGrp)

	return results, nil
}

func (ci *customerImporter) ImportFile() ([]EmailGroup, error) {
	// first open the file
	csvFile, err := os.Open(ci.CsvFileName)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %v", err)
	}
	defer csvFile.Close()

	// then creates a reader to access the file content
	reader := csv.NewReader(bufio.NewReader(csvFile))

	// finally call the Import function passing the Reader interface to retrieve a list of emails grouped by email domain
	customers, err := ci.Import(reader)
	if err != nil {
		return nil, err
	}
	return customers, nil
}

// read the header (first line) to determine in which column is placed the email column
func (ci *customerImporter) parseHeader(reader *csv.Reader) (int, error) {
	if reader == nil {
		return -1, fmt.Errorf("invalid reader")
	}

	// read the first line
	header, err := reader.Read()
	if err != nil {
		return -1, err
	}

	// loop through each element of the header searching for the 'email' column
	emailColIndex := 0
	for _, colName := range header {
		if colName == ci.EmailColumnName {
			return emailColIndex, nil
		}
		emailColIndex++
	}

	// if we arrive here, we haven't found the column name 'email'
	errorMsg := fmt.Sprintf("missing column name '%s' in the first line of the file", ci.EmailColumnName)
	log.Printf(errorMsg)
	return -1, fmt.Errorf(errorMsg)
}

func (ci *customerImporter) sortResults(custGrp map[string]int) []EmailGroup {
	// sort the elements alphabetically by email domain (it is, the key of the map)
	results := []EmailGroup{}
	keys := make([]string, 0, len(custGrp))
	for k := range custGrp {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// build the response with the sorted elements
	for _, domain := range keys {
		results = append(results, EmailGroup{
			EmailDomain: domain,
			Counter:     custGrp[domain],
		})
	}
	return results
}

func (ci *customerImporter) processLine(line []string, emailColIndex int, custGrp map[string]int, lineNumber int) error {
	// check that the line has the minimal number of elements before being processed
	if len(line) < emailColIndex {
		return fmt.Errorf("email field not found at line %d", lineNumber)
	}

	// extract the components (username and domain) from the email field
	emailField := line[emailColIndex]
	components := strings.Split(emailField, "@")

	// check that is a valid email (containing an '@' symbol)
	if len(components) == numExpectedEmailComponents {
		mailDomain := strings.ToLower(components[emailDomainIndex])

		// update the map, if the element does not exists will be created
		custGrp[mailDomain]++
		return nil
	}
	return fmt.Errorf("invalid email format")
}
