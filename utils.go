// csvParser.go

package main

import (
    "fmt"
	"encoding/csv"
    "os"
    "io"
    "strconv"
    "regexp"
    "strings"
    "net/url"
    "log"
    "net/http"
)

//This section represents the CSV parser methods and utils 
//to process CSV files with the parcelUser table as destination
//Since its only for testing and representation purpose, this 
//method wasnt worked to be as much as generic as it can be

//This method holds the logic to evaluate path received
//In case of external URL it will download the csv file into the localCsv directory
//and after process it, the csv will be erased to avoid garbage
func loadUsersFromCsv(a *App, csvPath string, csvLocation string) (int, error) {

	var csvFileName string = ""

	if csvPath != "" && evaluateRegexOnString(`http[s]://[A-Za-z./-_].csv$`,csvPath) == false{
		fmt.Println("Invalid path " , csvPath)
		return 0, nil
	}


	if csvLocation == "remote" {
		fmt.Println("Running under downloadable file " , csvPath)

	    // Build fileName from fullPath
	    fileURL, err := url.Parse(csvPath)
	    if err != nil {
	        log.Fatal(err)
	    }
	    path := fileURL.Path
	    segments := strings.Split(path, "/")
	    fileName := segments[len(segments)-1]
	 
	    // Create blank file
	    file, err := os.Create("localCsv/"+fileName)
	    if err != nil {
	        log.Fatal(err)
	    }
	    client := http.Client{
	        CheckRedirect: func(r *http.Request, via []*http.Request) error {
	            r.URL.Opaque = r.URL.Path
	            return nil
	        },
	    }
	    // Put content on file
	    resp, err := client.Get(csvPath)
	    if err != nil {
	        log.Fatal(err)
	    }
	    defer resp.Body.Close()
	 
	    size, err := io.Copy(file, resp.Body)
	 
	    defer file.Close()
	 
	    fmt.Printf("Downloaded a file %s with size %d bytes", fileName, size)

	    csvFileName = fileName

	}

	if csvFileName == "" {
		csvFileName = "localCsv/test_file.csv"
	}

	recordFile, err := os.Open(csvFileName)
	if err != nil {
		fmt.Println("An error encountered ::", err)
		return 0, err
	}

	count, err := computeCustomerDataFromCsv(a,recordFile)

	recordFile.Close()

	if csvLocation == "remote"{
		fErr := os.Remove(csvFileName)
	    if fErr != nil {
	        log.Fatal(fErr)
	    }
	}

	if err != nil{
		fmt.Println("Exec computeCustomerDataFromCsv with errors")
		return count, err
	}


	return count, nil
		
}

//This method will parse the content to create a set of users, 
//leaving a structure ready to be executed on model layer
func computeCustomerDataFromCsv(a *App, recordFile *os.File) (int,error){

	var setOfUsers = []userDataType{}
	var customer userDataType
	var count int = 0

	reader := csv.NewReader(recordFile)

	_, err := reader.Read()
	if err != nil {
		fmt.Println("An error encountered ::", err)
		return 0 , err
	}

	for i:= 0;; i++ {
		record, err := reader.Read()
		if err == io.EOF {
			break // reached end of the file
		} else if err != nil {
			fmt.Println("An error encountered ::", err)
			return 0 , err
		}
		count ++

		customer.ID, _ = strconv.Atoi(record[0])
		customer.Email = record[1]
		customer.Phone = record[2]
		valueOfParcelWeight, _ := strconv.ParseFloat(strings.TrimSpace(record[3]), 64)
		customer.ParcelWeight = valueOfParcelWeight

		if evaluateRegexOnString(`237\ ?[2368]\d{7,8}$`,record[2]) == true{
			customer.Country = 1
			setOfUsers = append(setOfUsers, customer)
			continue
		}
		if evaluateRegexOnString(`251\ ?[1-59]\d{8}$`,record[2]) == true{
			customer.Country = 2
			setOfUsers = append(setOfUsers, customer)
			continue
		}
		if evaluateRegexOnString(`212\ ?[5-9]\d{8}$`,record[2]) == true{
			customer.Country = 3
			setOfUsers = append(setOfUsers, customer)
			continue
		}
		if evaluateRegexOnString(`237\ ?[2368]\d{7,8}$`,record[2]) == true{
			customer.Country = 4
			setOfUsers = append(setOfUsers, customer)
			continue
		}
		if evaluateRegexOnString(`256\ ?\d{9}$`,record[2]) == true{
			customer.Country = 5
			setOfUsers = append(setOfUsers, customer)
			continue
		}
		customer.Country = 0
		setOfUsers = append(setOfUsers, customer)
	}

	countReg, err := createSetOfUsersByBatch(a.DB, setOfUsers)
	if err != nil{
		fmt.Println("Exec createSetOfUsersByBatch with errors")
		return 0 , err
	}

	return countReg , nil
}

//This method does the evaluation between a string and a regular expression provided
func evaluateRegexOnString(regex string, data string) (bool){
	var expression = regexp.MustCompile(regex)
	match := expression.MatchString(data)
	return match
}