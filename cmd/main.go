package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jimmymcguigan18/credit-card-approval-prediction/internal/evaluation"
	"github.com/jimmymcguigan18/credit-card-approval-prediction/internal/models"
	"github.com/jimmymcguigan18/credit-card-approval-prediction/internal/preprocessing"
	"github.com/jimmymcguigan18/credit-card-approval-prediction/internal/visualization"
)

func main() {
	// Define command line flags
	preprocessPtr := flag.Bool("preprocess", false, "Run data preprocessing")
	trainPtr := flag.Bool("train", false, "Train models")
	evaluatePtr := flag.Bool("evaluate", false, "Evaluate models")
	visualizePtr := flag.Bool("visualize", false, "Generate visualizations")
	flag.Parse()

	// Get project root directory
	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error getting executable path: %v\n", err)
		os.Exit(1)
	}

	// Assuming the executable is in the project root or built into it
	projectRoot := filepath.Dir(execPath)

	// If no flags are specified, run all steps
	runAll := !*preprocessPtr && !*trainPtr && !*evaluatePtr && !*visualizePtr

	// Define file paths
	rawDataPath := filepath.Join(projectRoot, "data", "raw", "crx.data")
	trainDataPath := filepath.Join(projectRoot, "data", "processed", "train.csv")
	testDataPath := filepath.Join(projectRoot, "data", "processed", "test.csv")
	modelEvalPath := filepath.Join(projectRoot, "data", "processed", "model_evaluation.csv")
	visualizationDir := filepath.Join(projectRoot, "data", "processed", "visualizations")
	confusionMatrixDir := filepath.Join(projectRoot, "data", "processed", "confusion_matrices")

	// Initialize evaluation object
	modelEval := evaluation.NewModelEvaluation()

	// Run the pipeline steps based on flags
	if *preprocessPtr || runAll {
		fmt.Println("Running preprocessing...")
		// Implement preprocessing
		data, err := preprocessing.LoadData(rawDataPath)
		if err != nil {
			fmt.Printf("Error loading data: %v\n", err)
			os.Exit(1)
		}

		// Handle missing values
		data.HandleMissingValues()

		// Encode categorical variables
		if err := data.EncodeCategoricalFeatures(); err != nil {
			fmt.Printf("Error encoding categorical features: %v\n", err)
			os.Exit(1)
		}

		// Convert target variable
		if err := data.ConvertTargetVariable(); err != nil {
			fmt.Printf("Error converting target variable: %v\n", err)
			os.Exit(1)
		}

		// Normalize numerical features
		data.NormalizeFeatures()

		// Save processed data
		err = data.SaveProcessedData(trainDataPath, testDataPath)
		if err != nil {
			fmt.Printf("Error saving processed data: %v\n", err)
			os.Exit(1)
		}

		// Split data into train and test sets (already handled in SaveProcessedData)
		if err != nil {
			fmt.Printf("Error splitting data: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Preprocessing completed successfully!")
	}

	if *trainPtr || runAll {
		fmt.Println("Training models...")
		// Implement model training
		modelResults, err := models.TrainAllModels(trainDataPath, testDataPath)
		if err != nil {
			fmt.Printf("Error training models: %v\n", err)
			os.Exit(1)
		}

		// Add results to evaluation
		for _, result := range modelResults {
			modelEval.AddResult(result)
		}

		fmt.Println("Model training completed successfully!")
	}

	if *evaluatePtr || runAll {
		fmt.Println("Evaluating models...")
		// Implement model evaluation
		modelEval.PrintResults()

		// Save evaluation results
		err := modelEval.SaveResultsToCSV(modelEvalPath)
		if err != nil {
			fmt.Printf("Error saving evaluation results: %v\n", err)
			os.Exit(1)
		}

		// Save confusion matrices
		err = modelEval.SaveConfusionMatrices(confusionMatrixDir)
		if err != nil {
			fmt.Printf("Error saving confusion matrices: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Model evaluation completed successfully!")
	}

	if *visualizePtr || runAll {
		fmt.Println("Generating visualizations...")
		// Create visualization directory if it doesn't exist
		if err := visualization.CreateOutputDir(visualizationDir); err != nil {
			fmt.Printf("Error creating visualization directory: %v\n", err)
			os.Exit(1)
		}

		// Generate all visualizations
		err := visualization.GenerateAllVisualizations(
			trainDataPath,
			visualizationDir,
			modelEval.Results,
		)
		if err != nil {
			fmt.Printf("Error generating visualizations: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Visualization generation completed successfully!")
	}

	fmt.Println("Pipeline completed successfully!")
}
