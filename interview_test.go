package customerimporter

import (
	"encoding/csv"
	"strings"
	"testing"
)

func TestImport(t *testing.T) {

	expected := make(map[string]int)
	expected["github.io"] = 2
	expected["cnet.com"] = 1
	expected["whitehouse.gov"] = 1
	expected["woothemes.com"] = 2
	expected["google.fr"] = 3

	content := `first_name,last_name,email,gender,ip_address
				Mildred,Hernandez,mhernandez0@github.io,Female,38.194.51.128
				Norma,Allen,nallen8@cnet.com,Female,168.67.162.1
				Anna,Rivera,ariverag@whitehouse.gov,Female,105.158.80.2
				Lori,Elliott,lelliotth@github.io,Female,160.108.154.74
				Wanda,Lewis,wlewisk@woothemes.com,Female,25.32.100.250
				Robert,Hunter,rhunterp@google.fr,Male,130.35.232.64
				Gregory,Ryan,gryanq@google.fr,Male,188.242.255.152
				Andrew,Morgan,amorganr@google.fr,Male,3.184.160.117
				Peter,Day,pdays@woothemes.com,Male,0.24.246.12`

	contentBadHeader := `first_name,last_name,MAIL,gender,ip_address
	Mildred,Hernandez,mhernandez0@github.io,Female,38.194.51.128
	Bonnie,Ortiz,bortiz1@cyberchimps.com,Female,197.54.209.129
	Dennis,Henry,dhenry2@hubpages.com,Male,155.75.186.217
	Justin,Hansen,jhansen3@360.cn,Male,251.166.224.119`

	contentBadFormat := `first_name,last_name,email,gender,ip_address
	Mildred,Hernandez,mhernandez0@github.io,Female,38.194.51.128
	Bonnie,Ortiz,bortiz1@cyberchimps.com,Female,197.54.209.129
	Dennis,Henry,dhenry2,Male,155.75.186.217
	Justin,Hansen,jhansen3@360.cn,Male,251.166.224.119`

	contentMissingEmail := `first_name,last_name,email,gender,ip_address
	Mildred,Hernandez,mhernandez0@github.io,Female,38.194.51.128
	Dennis,Henry
	Bonnie,Ortiz,bortiz1@cyberchimps.com,Female,197.54.209.129
	Justin,Hansen,jhansen3@360.cn,Male,251.166.224.119`

	tt := []struct {
		name            string
		readerIsNil     bool
		csvFileContent  string
		expectedResults map[string]int
		errMsg          string
	}{
		{"Positive case", false, content, expected, ""},
		{"Reader is nil", true, "", nil, "invalid reader"},
		{"Bad header", false, contentBadHeader, nil, "missing column name 'email' in the first line of the file"},
		{"Bad email format", false, contentBadFormat, nil, "skip line 4: invalid email format"},
		{"Missing email value", false, contentMissingEmail, nil, "record on line 3: wrong number of fields"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ci := New()

			var reader *csv.Reader
			if !tc.readerIsNil {
				reader = csv.NewReader(strings.NewReader(tc.csvFileContent))
			}
			customers, err := ci.Import(reader)

			if !checkError(tc.errMsg, err, t) {
				if tc.expectedResults != nil {
					if len(tc.expectedResults) != len(customers) {
						t.Errorf("expected customers length %d; got %d", len(tc.expectedResults), len(customers))
					}
					foundKey := false
					for k, v := range tc.expectedResults {
						for _, cg := range customers {
							if cg.EmailDomain == k {
								foundKey = true
								if cg.Counter != v {
									t.Errorf("expected counter %d and got %d for email domain %s", v, cg.Counter, k)
								}
								break
							}
						}
						if !foundKey {
							// here the email domain is not contained in the results
							t.Errorf("missing email domain %s", k)
						}
					}
				}
			}
		})
	}
}

func TestImportFile(t *testing.T) {

	tt := []struct {
		name               string
		csvFileName        string
		expectedNumResults int
		errMsg             string
	}{
		{"Positive case", "", 500, ""},
		{"Missing CSV file", "dummy_file.csv", -1, "unable to open file: open dummy_file.csv: no such file or directory"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ci := New()

			if tc.csvFileName != "" {
				ci.SetCsvFileName(tc.csvFileName)
			}

			customers, err := ci.ImportFile()

			if !checkError(tc.errMsg, err, t) {
				if tc.expectedNumResults != -1 {
					if tc.expectedNumResults != len(customers) {
						t.Errorf("expected customers length %d; got %d", tc.expectedNumResults, len(customers))
					}
				}
			}
		})
	}
}

func checkError(errMsg string, err error, t *testing.T) bool {
	if err != nil {
		if errMsg == "" {
			// here the testcase didn't expect any error, but found an error
			t.Errorf("unexpected error: %v", err)
		} else if errMsg != err.Error() {
			// here the testcase expected another error than the received
			t.Errorf("expected error message: %v; got: %v", errMsg, err.Error())
		}
		return true
	}
	return false
}
