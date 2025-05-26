package evaluation

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/jimmymcguigan18/credit-card-approval-prediction/internal/models"
)

// ModelEvaluation contains evaluation metrics for all models
type ModelEvaluation struct {
	Results map[string]*models.ModelResult
}

// NewModelEvaluation creates a new ModelEvaluation instance
func NewModelEvaluation() *ModelEvaluation {
	return &ModelEvaluation{
		Results: make(map[string]*models.ModelResult),
	}
}

// AddResult adds a model result to the evaluation
func (me *ModelEvaluation) AddResult(result *models.ModelResult) {
	me.Results[result.ModelName] = result
}

// GetBestModel returns the name of the best performing model based on F1 score
func (me *ModelEvaluation) GetBestModel() string {
	bestScore := -1.0
	bestModel := ""

	for name, result := range me.Results {
		if result.F1Score > bestScore {
			bestScore = result.F1Score
			bestModel = name
		}
	}

	return bestModel
}

// PrintResults prints the evaluation results to the console
func (me *ModelEvaluation) PrintResults() {
	fmt.Println("\nModel Evaluation Results:")
	fmt.Println("=========================")

	// Print header
	fmt.Printf("%-20s %-10s %-10s %-10s %-10s\n", "Model", "Accuracy", "Precision", "Recall", "F1 Score")
	fmt.Println("------------------------------------------------------------")

	// Print results for each model
	for name, result := range me.Results {
		fmt.Printf("%-20s %-10.4f %-10.4f %-10.4f %-10.4f\n",
			name, result.Accuracy, result.Precision, result.Recall, result.F1Score)
	}

	// Print best model
	bestModel := me.GetBestModel()
	if bestModel != "" {
		fmt.Println("\nBest Model (by F1 Score):")
		fmt.Printf("%-20s %-10.4f %-10.4f %-10.4f %-10.4f\n",
			bestModel,
			me.Results[bestModel].Accuracy,
			me.Results[bestModel].Precision,
			me.Results[bestModel].Recall,
			me.Results[bestModel].F1Score)
	}
}

// SaveResultsToCSV saves the evaluation results to a CSV file
func (me *ModelEvaluation) SaveResultsToCSV(outputPath string) error {
	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Model", "Accuracy", "Precision", "Recall", "F1 Score"}
	err = writer.Write(header)
	if err != nil {
		return fmt.Errorf("error writing header: %v", err)
	}

	// Write results for each model
	for name, result := range me.Results {
		row := []string{
			name,
			strconv.FormatFloat(result.Accuracy, 'f', 4, 64),
			strconv.FormatFloat(result.Precision, 'f', 4, 64),
			strconv.FormatFloat(result.Recall, 'f', 4, 64),
			strconv.FormatFloat(result.F1Score, 'f', 4, 64),
		}

		err = writer.Write(row)
		if err != nil {
			return fmt.Errorf("error writing row: %v", err)
		}
	}

	return nil
}

// AnalyzeFeatureImportance analyzes feature importance from model results
// This is a placeholder function that would be implemented with actual model-specific
// feature importance extraction in a real application
func (me *ModelEvaluation) AnalyzeFeatureImportance() map[string]float64 {
	// In a real implementation, this would extract feature importance from models
	// For now, return mock data
	return map[string]float64{
		"A2":  0.15,
		"A3":  0.12,
		"A8":  0.18,
		"A11": 0.09,
		"A14": 0.14,
		"A15": 0.11,
		"A1":  0.07,
		"A4":  0.06,
		"A5":  0.05,
		"A6":  0.03,
	}
}

// SaveConfusionMatrices saves confusion matrices for all models to CSV files
func (me *ModelEvaluation) SaveConfusionMatrices(outputDir string) error {
	// Create output directory if it doesn't exist
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.MkdirAll(outputDir, 0755)
		if err != nil {
			return fmt.Errorf("error creating output directory: %v", err)
		}
	}

	// Save confusion matrix for each model
	for name, result := range me.Results {
		// Create output file
		filePath := fmt.Sprintf("%s/%s_confusion_matrix.csv", outputDir, name)
		file, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("error creating output file: %v", err)
		}

		// Create CSV writer
		writer := csv.NewWriter(file)

		// Get classes
		classes := make([]string, 0, len(result.ConfMatrix))
		for class := range result.ConfMatrix {
			classes = append(classes, class)
		}

		// Write header
		header := append([]string{"Actual/Predicted"}, classes...)
		err = writer.Write(header)
		if err != nil {
			file.Close()
			return fmt.Errorf("error writing header: %v", err)
		}

		// Write confusion matrix
		for _, actualClass := range classes {
			row := []string{actualClass}
			for _, predictedClass := range classes {
				count := result.ConfMatrix[actualClass][predictedClass]
				row = append(row, strconv.Itoa(count))
			}

			err = writer.Write(row)
			if err != nil {
				file.Close()
				return fmt.Errorf("error writing row: %v", err)
			}
		}

		writer.Flush()
		file.Close()
	}

	return nil
}
