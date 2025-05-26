package preprocessing

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

// CreditData represents the structure of our credit card approval dataset
type CreditData struct {
	DF dataframe.DataFrame
}

// LoadData loads the credit card dataset from a CSV file
func LoadData(filepath string) (*CreditData, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Read CSV file
	reader := csv.NewReader(file)
	reader.Comma = ','
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %v", err)
	}

	// Ensure we have data
	if len(records) < 1 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	// Define column names
	colNames := []string{"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8", "A9", "A10", "A11", "A12", "A13", "A14", "A15", "A16"}

	// Create dataframe with explicit type inference
	df := dataframe.LoadRecords(
		records,
		dataframe.DetectTypes(false),
		dataframe.DefaultType(series.String),
	)

	// Set column names
	df.SetNames(colNames...)

	return &CreditData{DF: df}, nil
}

// HandleMissingValues imputes missing values in the dataset
func (cd *CreditData) HandleMissingValues() {
	// Replace '?' with NaN for all columns
	for _, colName := range cd.DF.Names() {
		cd.DF = cd.DF.Mutate(
			cd.DF.Col(colName).Map(func(e series.Element) series.Element {
				if e.IsNA() {
					return e
				}
				str, ok := e.Val().(string)
				if ok && str == "?" {
					e.Set(nil)
				}
				return e
			}),
		)
	}

	// For categorical variables, replace missing values with the most frequent value
	categoricalCols := []string{"A1", "A4", "A5", "A6", "A7", "A9", "A10", "A12", "A13"}
	for _, col := range categoricalCols {
		// Find most frequent value
		valCounts := make(map[string]int)
		// Get the column and ensure it exists
		s := cd.DF.Col(col)
		if s.Err != nil {
			continue // Skip this column if it doesn't exist
		}

		// Convert all values to strings first
		strVals := make([]interface{}, s.Len())
		for i := 0; i < s.Len(); i++ {
			e := s.Elem(i)
			if e.IsNA() {
				strVals[i] = nil
				continue
			}
			strVals[i] = fmt.Sprintf("%v", e.Val())
		}

		// Create new string series and update dataframe
		newSeries := series.New(strVals, series.String, col)
		cd.DF = cd.DF.Mutate(newSeries)

		// Count values in the new string series
		newSeries.Map(func(e series.Element) series.Element {
			if !e.IsNA() {
				if str, ok := e.Val().(string); ok {
					valCounts[str]++
				}
			}
			return e
		})

		// Find the most frequent value
		mostFreqVal := ""
		maxCount := 0
		for val, count := range valCounts {
			if count > maxCount {
				maxCount = count
				mostFreqVal = val
			}
		}

		// Replace missing values with most frequent value
		cd.DF = cd.DF.Mutate(
			series.New(
				newSeries.Map(func(e series.Element) series.Element {
					if e.IsNA() {
						e.Set(mostFreqVal)
					}
					return e
				}),
				series.String,
				col,
			),
		)
	}

	// For continuous variables, replace missing values with the mean
	continuousCols := []string{"A2", "A3", "A8", "A11", "A14", "A15"}
	for _, col := range continuousCols {
		// Get the column and ensure it exists
		s := cd.DF.Col(col)
		if s.Err != nil {
			continue // Skip this column if it doesn't exist
		}

		// First pass: convert values to float64 and calculate mean
		sum := 0.0
		count := 0
		floatVals := make([]interface{}, s.Len())

		for i := 0; i < s.Len(); i++ {
			e := s.Elem(i)
			if e.IsNA() {
				floatVals[i] = nil
				continue
			}

			// Try to convert the value to float64
			strVal := fmt.Sprintf("%v", e.Val())
			if val, err := strconv.ParseFloat(strVal, 64); err == nil {
				floatVals[i] = val
				sum += val
				count++
			} else {
				floatVals[i] = nil
			}
		}

		// Calculate mean
		mean := 0.0
		if count > 0 {
			mean = sum / float64(count)
		}

		// Second pass: fill missing values with mean
		for i := 0; i < len(floatVals); i++ {
			if floatVals[i] == nil {
				floatVals[i] = mean
			}
		}

		// Update the dataframe with the new float series
		cd.DF = cd.DF.Mutate(
			series.New(floatVals, series.Float, col),
		)
	}
}

// EncodeCategoricalFeatures converts categorical features to numerical values
func (cd *CreditData) EncodeCategoricalFeatures() error {
	// One-hot encode categorical variables
	categoricalCols := []string{"A1", "A4", "A5", "A6", "A7", "A9", "A10", "A12", "A13"}

	// Verify DataFrame is not nil
	if cd.DF.Err != nil {
		return fmt.Errorf("invalid dataframe: %v", cd.DF.Err)
	}

	for _, col := range categoricalCols {
		// Get the column and ensure it exists
		s := cd.DF.Col(col)
		if s.Err != nil {
			fmt.Printf("Warning: Column %s not found, skipping\n", col)
			continue
		}

		// Get unique values
		uniqueVals := make(map[string]bool)
		for i := 0; i < s.Len(); i++ {
			e := s.Elem(i)
			if !e.IsNA() {
				strVal := fmt.Sprintf("%v", e.Val())
				if strVal != "" {
					uniqueVals[strVal] = true
				}
			}
		}

		// Create one-hot encoded columns
		for val := range uniqueVals {
			newColName := fmt.Sprintf("%s_%s", col, val)
			oneHotVals := make([]interface{}, s.Len())

			// Fill one-hot values
			for i := 0; i < s.Len(); i++ {
				e := s.Elem(i)
				if e.IsNA() {
					oneHotVals[i] = 0
					continue
				}
				strVal := fmt.Sprintf("%v", e.Val())
				if strVal == val {
					oneHotVals[i] = 1
				} else {
					oneHotVals[i] = 0
				}
			}

			// Add the one-hot encoded column
			newSeries := series.New(oneHotVals, series.Int, newColName)
			if newSeries.Err != nil {
				return fmt.Errorf("error creating one-hot encoded column %s: %v", newColName, newSeries.Err)
			}
			cd.DF = cd.DF.Mutate(newSeries)
		}
	}
	return nil
}

// ConvertTargetVariable converts the target variable (A16) to binary (0/1)
func (cd *CreditData) ConvertTargetVariable() error {
	// Get the target column
	s := cd.DF.Col("A16")
	if s.Err != nil {
		return fmt.Errorf("error accessing target column A16: %v", s.Err)
	}

	// Convert target variable to binary (0/1)
	cd.DF = cd.DF.Mutate(
		series.New(
			s.Map(func(e series.Element) series.Element {
				if e.IsNA() {
					e.Set(0)
					return e
				}
				str := fmt.Sprintf("%v", e.Val())
				if str == "+" {
					e.Set(1)
				} else {
					e.Set(0)
				}
				return e
			}),
			series.Int,
			"A16",
		),
	)
	return nil
}

// NormalizeFeatures scales numerical features to a standard range
func (cd *CreditData) NormalizeFeatures() {
	continuousCols := []string{"A2", "A3", "A8", "A11", "A14", "A15"}
	for _, col := range continuousCols {
		// Find min and max values
		min := math.MaxFloat64
		max := -math.MaxFloat64

		// Iterate through each element to find min and max
		for i := 0; i < cd.DF.Nrow(); i++ {
			v := cd.DF.Col(col).Elem(i).Val()
			str, ok := v.(string)
			if !ok {
				continue
			}
			val, err := strconv.ParseFloat(str, 64)
			if err != nil {
				continue
			}
			if val < min {
				min = val
			}
			if val > max {
				max = val
			}
		}

		// Skip normalization if min equals max
		if min == max {
			continue
		}

		// Normalize values to [0,1] range
		// Create a temporary series to hold the normalized values
		values := make([]interface{}, cd.DF.Nrow())
		for i := 0; i < cd.DF.Nrow(); i++ {
			v := cd.DF.Col(col).Elem(i).Val()
			str, ok := v.(string)
			if !ok {
				values[i] = 0.0
				continue
			}
			val, err := strconv.ParseFloat(str, 64)
			if err != nil {
				values[i] = 0.0
				continue
			}
			normalized := (val - min) / (max - min)
			values[i] = normalized
		}

		cd.DF = cd.DF.Mutate(
			series.New(
				values,
				series.Float,
				fmt.Sprintf("%s_norm", col),
			),
		)
	}
}

// SplitTrainTest splits the data into training and testing sets
func (cd *CreditData) SplitTrainTest(testSize float64) (trainDF, testDF dataframe.DataFrame) {
	// Shuffle the data
	shuffled := cd.DF.Arrange(dataframe.Sort("target"))

	// Calculate split index
	totalRows := shuffled.Nrow()
	testRows := int(float64(totalRows) * testSize)
	trainRows := totalRows - testRows

	// Split the data
	trainDF = shuffled.Subset(series.Ints(generateRange(0, trainRows)))
	testDF = shuffled.Subset(series.Ints(generateRange(trainRows, totalRows)))

	return trainDF, testDF
}

// generateRange creates a slice of integers from start to end (exclusive)
func generateRange(start, end int) []int {
	rangeSlice := make([]int, end-start)
	for i := 0; i < end-start; i++ {
		rangeSlice[i] = start + i
	}
	return rangeSlice
}

// SaveProcessedData saves the processed data to CSV files
func (cd *CreditData) SaveProcessedData(trainPath, testPath string) error {
	// Split the data
	trainDF, testDF := cd.SplitTrainTest(0.2)

	// Save training data
	trainFile, err := os.Create(trainPath)
	if err != nil {
		return fmt.Errorf("error creating training file: %v", err)
	}
	defer trainFile.Close()

	trainWriter := csv.NewWriter(trainFile)
	defer trainWriter.Flush()

	// Write header
	if err := trainWriter.Write(trainDF.Names()); err != nil {
		return fmt.Errorf("error writing training header: %v", err)
	}

	// Write data
	for i := 0; i < trainDF.Nrow(); i++ {
		row := make([]string, trainDF.Ncol())
		for j := range trainDF.Names() {
			row[j] = fmt.Sprintf("%v", trainDF.Elem(i, j).Val())
		}
		if err := trainWriter.Write(row); err != nil {
			return fmt.Errorf("error writing training row: %v", err)
		}
	}

	// Save test data
	testFile, err := os.Create(testPath)
	if err != nil {
		return fmt.Errorf("error creating test file: %v", err)
	}
	defer testFile.Close()

	testWriter := csv.NewWriter(testFile)
	defer testWriter.Flush()

	// Write header
	if err := testWriter.Write(testDF.Names()); err != nil {
		return fmt.Errorf("error writing test header: %v", err)
	}

	// Write data
	for i := 0; i < testDF.Nrow(); i++ {
		row := make([]string, testDF.Ncol())
		for j := range testDF.Names() {
			row[j] = fmt.Sprintf("%v", testDF.Elem(i, j).Val())
		}
		if err := testWriter.Write(row); err != nil {
			return fmt.Errorf("error writing test row: %v", err)
		}
	}

	return nil
}

// PreprocessPipeline runs the complete preprocessing pipeline
func PreprocessPipeline(inputPath, trainOutputPath, testOutputPath string) error {
	// Load data
	data, err := LoadData(inputPath)
	if err != nil {
		return fmt.Errorf("error loading data: %v", err)
	}

	// Apply preprocessing steps
	data.HandleMissingValues()
	data.EncodeCategoricalFeatures()
	data.NormalizeFeatures()

	// Save processed data
	err = data.SaveProcessedData(trainOutputPath, testOutputPath)
	if err != nil {
		return fmt.Errorf("error saving processed data: %v", err)
	}

	return nil
}
