# Customer Importer

## Introduction
This package reads customer records and group them by email domain calculating the number of customers with the same domain.

## Initialization

### New
Creates an instance of an object that implements the CustomerImporter interface.

```go
func New() CustomerImporter 
```

## CustomerImporter interface
This interface exposes the following public methods: `SetCsvFileName`, `SetEmailColumnName`, `Import` and `ImportFile`.

```go 
type CustomerImporter interface {
	SetCsvFileName(fileName string)
	SetEmailColumnName(columnName string)
	Import(reader *csv.Reader) ([]EmailGroup, error)
	ImportFile() ([]EmailGroup, error)
}
```

### SetCsvFileName
Changes the default CSV filename `customers.csv` to a new one.

``` go
func (ci *customerImporter) SetCsvFileName(fileName string)
```

### SetEmailColumnName
Changes the default email column name `email` to a new one.

``` go
func (ci *customerImporter) SetEmailColumnName(columnName string)
```


### Import
This method reads the customer records from the reader passed as input parameter, groups them by email domain and returns a list of email domains  sorted alphabetically including a counter.

``` go
func (ci *customerImporter) Import(reader *csv.Reader) ([]EmailGroup, error)
```

### ImportFile
This method reads from the given file (by default *customers.csv*), groups them by email domain and returns a list of email domains sorted alphabetically including a counter.

``` go
func (ci *customerImporter) ImportFile() ([]EmailGroup, error)
```

### EmailGroup
EmailGroup is the type used in the import responses, it is an structure formed by an email domain and a counter (how many email addresses have the same email domain).

``` go
type EmailGroup struct {
	EmailDomain string
	Counter     int
}
```

## CSV Example: `customers.csv`
The file `customers.csv` contains data samples for testing purposes.

Here is an overview of the first rows:

``` csv
first_name,last_name,email,gender,ip_address
Mildred,Hernandez,mhernandez0@github.io,Female,38.194.51.128
Bonnie,Ortiz,bortiz1@cyberchimps.com,Female,197.54.209.129
Dennis,Henry,dhenry2@hubpages.com,Male,155.75.186.217
```

Notice that the first row of the file is expected to be the column names separated by commas. By default, the package will search for the column with name `email`, although this value can be changed by using the `SetEmailColumnName` method.



## Usage
To use this package, first of all we need to instanciate a CustomerImporter by calling the function `New`. After this, we can either import customer records with the following methods: `Import` or  `ImportFile`.

### Import (from a file)

```go
// open the file and create a reader
csvFile, err := os.Open("customers.csv")
if err != nil {
    return nil, fmt.Errorf("unable to open file: %v", err)
}
reader := csv.NewReader(bufio.NewReader(csvFile))

// instanciate the customer importer
customerImporter := New()

// import from a reader
customers, err := customerImporter.Import(reader)
if err != nil {
    return nil, err
}
```


### Import (from a string)

```go
// create a reader from a string
text :=  `first_name,last_name,email,gender,ip_address
          Mildred,Hernandez,mhernandez0@github.io,Female,38.194.51.128
          Norma,Allen,nallen8@cnet.com,Female,168.67.162.1`
reader := csv.NewReader(strings.NewReader(text))

// instanciate the customer importer
customerImporter := New()

// import from a reader
customers, err := customerImporter.Import(reader)
if err != nil {
    return nil, err
}
```


### ImportFile

```go
// instanciate the customer importer
customerImporter := New()

// import data from default file ('customers.csv')
customers, err := customerImporter.ImportFile()
if err != nil {
    return nil, err
}
```

### ImportFile (from custom file)

```go
// instanciate the customer importer
customerImporter := New()

// change CSV filename 
customerImporter.SetCsvFileName("myCustomFile.csv")

// import data from previous defined file
customers, err := customerImporter.ImportFile()
if err != nil {
    return nil, err
}
```


### Change default email column name

```go
// create a reader from a string
text :=  `first_name,last_name,my_custom_email_column_name,gender,ip_address
          Mildred,Hernandez,mhernandez0@github.io,Female,38.194.51.128
          Norma,Allen,nallen8@cnet.com,Female,168.67.162.1`
reader := csv.NewReader(strings.NewReader(text))

// instanciate the customer importer
customerImporter := New()

// change the expected email column name
customerImporter.SetEmailColumnName("my_custom_email_column_name")

// import from a reader
customers, err := customerImporter.Import(reader)
if err != nil {
    return nil, err
}
```

## Testing
The file `interview_test.go` contains the test cases for the `Import` and `ImportFile` methods.

### Run tests
```bash
go test -v
```