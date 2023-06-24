package core

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
)

// import a CSV code file
func ReadCSVFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func ReadTextFile(filePath string) []SnomedDescription {
	// we'll use a scanner to parse each lines of the files
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	descriptions := []SnomedDescription{}

	// skip the first line, as it contains headers
	scanner.Scan()

	// iterate on each lines
	for scanner.Scan() {
		cols := strings.Split(scanner.Text(), "\t")
		id, err := strconv.Atoi(cols[0])
		if err != nil {
			log.Fatal(err)
		}
		effectiveTime, err := strconv.Atoi(cols[1])
		if err != nil {
			log.Fatal(err)
		}
		active, err := strconv.Atoi(cols[2])
		if err != nil {
			log.Fatal(err)
		}
		moduleId, err := strconv.ParseUint(cols[3], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		conceptId, err := strconv.Atoi(cols[4])
		if err != nil {
			log.Fatal(err)
		}
		languageCode := cols[5]
		typeId, err := strconv.ParseUint(cols[6], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		term := cols[7]
		caseSignificanceId, err := strconv.ParseUint(cols[8], 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		// create a new description record
		description := &SnomedDescription{
			Id:                 id,
			EffectiveTime:      effectiveTime,
			Active:             active == 1,
			ModuleId:           moduleId,
			ConceptId:          conceptId,
			LanguageCode:       languageCode,
			TypeId:             typeId,
			Term:               term,
			CaseSignificanceId: caseSignificanceId,
		}
		log.Println("Concept ID: " + strconv.Itoa(description.ConceptId) + " Term: " + description.Term)

		// add description to slice
		descriptions = append(descriptions, *description)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return descriptions
}
