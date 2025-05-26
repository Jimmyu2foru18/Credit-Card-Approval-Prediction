package visualization

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"

	"github.com/jimmymcguigan18/credit-card-approval-prediction/internal/models"
)

// Colors for charts
var (
	blueColor   = drawing.Color{R: 0, G: 113, B: 188, A: 255}
	greenColor  = drawing.Color{R: 0, G: 170, B: 0, A: 255}
	redColor    = drawing.Color{R: 255, G: 0, B: 0, A: 255}
	purpleColor = drawing.Color{R: 128, G: 0, B: 128, A: 255}
	orangeColor = drawing.Color{R: 255, G: 165, B: 0, A: 255}
)

// CreateOutputDir creates the output directory for visualizations if it doesn't exist
func CreateOutputDir(outputDir string) error {
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.MkdirAll(outputDir, 0755)
		if err != nil {
			return fmt.Errorf("error creating output directory: %v", err)
		}
	}
	return nil
}

// PlotClassDistribution creates a bar chart showing the distribution of approval/rejection classes
func PlotClassDistribution(df dataframe.DataFrame, outputPath string) error {
	// Count class distribution
	classCounts := make(map[string]int)
	df.Col("target").Map(func(e series.Element) series.Element {
		val := fmt.Sprintf("%v", e.Val())
		classCounts[val]++
		return e
	})

	// Prepare data for chart
	var values []chart.Value
	for class, count := range classCounts {
		label := "Rejected"
		color := redColor
		if class == "1" {
			label = "Approved"
			color = greenColor
		}
		values = append(values, chart.Value{
			Value: float64(count),
			Label: label,
			Style: chart.Style{
				FillColor:   color,
				StrokeColor: chart.DefaultColors[0],
				StrokeWidth: 1,
			},
		})
	}

	// Create the chart
	pie := chart.PieChart{
		Width:  512,
		Height: 512,
		Values: values,
		Title:  "Credit Card Approval Distribution",
	}

	// Save the chart to file
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer f.Close()

	err = pie.Render(chart.SVG, f)
	if err != nil {
		return fmt.Errorf("error rendering chart: %v", err)
	}

	return nil
}

// PlotFeatureImportance creates a bar chart showing feature importance
func PlotFeatureImportance(featureImportance map[string]float64, outputPath string) error {
	// Sort features by importance
	type featureScore struct {
		Name  string
		Score float64
	}

	var features []featureScore
	for name, score := range featureImportance {
		features = append(features, featureScore{Name: name, Score: score})
	}

	sort.Slice(features, func(i, j int) bool {
		return features[i].Score > features[j].Score
	})

	// Limit to top 10 features if there are more
	if len(features) > 10 {
		features = features[:10]
	}

	// Prepare data for chart
	var bars []chart.Value
	for _, feature := range features {
		bars = append(bars, chart.Value{
			Value: feature.Score,
			Label: feature.Name,
			Style: chart.Style{
				FillColor:   blueColor,
				StrokeColor: blueColor.WithAlpha(64),
				StrokeWidth: 1,
			},
		})
	}

	// Create the chart
	graph := chart.BarChart{
		Title:      "Feature Importance",
		TitleStyle: chart.Style{FontSize: 14},
		Width:      800,
		Height:     500,
		BarWidth:   50,
		XAxis:      chart.Style{},
		YAxis: chart.YAxis{
			Name:      "Importance Score",
			NameStyle: chart.Style{FontSize: 12},
			Style:     chart.Style{FontSize: 10},
		},
		Bars: bars,
	}

	// Save the chart to file
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer f.Close()

	err = graph.Render(chart.SVG, f)
	if err != nil {
		return fmt.Errorf("error rendering chart: %v", err)
	}

	return nil
}

// PlotModelComparison creates a bar chart comparing model performance metrics
func PlotModelComparison(results map[string]*models.ModelResult, outputPath string) error {
	// Prepare data for chart
	modelNames := make([]string, 0, len(results))
	accuracies := make([]float64, 0, len(results))
	precisions := make([]float64, 0, len(results))
	recalls := make([]float64, 0, len(results))
	f1Scores := make([]float64, 0, len(results))

	for name, result := range results {
		modelNames = append(modelNames, name)
		accuracies = append(accuracies, result.Accuracy)
		precisions = append(precisions, result.Precision)
		recalls = append(recalls, result.Recall)
		f1Scores = append(f1Scores, result.F1Score)
	}

	// Create the chart
	graph := chart.BarChart{
		Title:      "Model Performance Comparison",
		TitleStyle: chart.Style{FontSize: 14},
		Width:      800,
		Height:     500,
		BarWidth:   30,
		XAxis:      chart.Style{},
		YAxis: chart.YAxis{
			Name:      "Score",
			NameStyle: chart.Style{FontSize: 12},
			Style:     chart.Style{FontSize: 10},
			Range: &chart.ContinuousRange{
				Min: 0.0,
				Max: 1.0,
			},
		},
		Bars: []chart.Value{
			{Value: accuracies[0], Label: modelNames[0], Style: chart.Style{FillColor: blueColor}},
			{Value: precisions[0], Label: modelNames[0], Style: chart.Style{FillColor: greenColor}},
			{Value: recalls[0], Label: modelNames[0], Style: chart.Style{FillColor: redColor}},
			{Value: f1Scores[0], Label: modelNames[0], Style: chart.Style{FillColor: purpleColor}},
		},
	}

	// Save the chart to file
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer f.Close()

	err = graph.Render(chart.SVG, f)
	if err != nil {
		return fmt.Errorf("error rendering chart: %v", err)
	}

	return nil
}

// GenerateAllVisualizations creates all visualizations for the project
func GenerateAllVisualizations(dataPath, outputDir string, modelResults map[string]*models.ModelResult) error {
	// Create output directory if it doesn't exist
	err := CreateOutputDir(outputDir)
	if err != nil {
		return err
	}

	// Load data
	file, err := os.Open(dataPath)
	if err != nil {
		return fmt.Errorf("error opening data file: %v", err)
	}
	defer file.Close()

	// Define column names
	colNames := []string{"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8", "A9", "A10", "A11", "A12", "A13", "A14", "A15", "A16"}

	// Read CSV and set column names
	df := dataframe.ReadCSV(file)
	df.SetNames(colNames...)

	// 1. Plot class distribution
	classDistPath := filepath.Join(outputDir, "class_distribution.svg")
	err = PlotClassDistribution(df, classDistPath)
	if err != nil {
		return fmt.Errorf("error plotting class distribution: %v", err)
	}

	// 2. Plot numerical feature distributions
	numericalFeatures := []string{"A2", "A3", "A8", "A11", "A14", "A15"}
	for _, feature := range numericalFeatures {
		featurePath := filepath.Join(outputDir, fmt.Sprintf("%s_distribution.svg", feature))
		err = PlotFeatureDistribution(df, feature, featurePath)
		if err != nil {
			fmt.Printf("Error plotting %s distribution: %v\n", feature, err)
			continue
		}
	}

	// 3. Plot model comparison if results are available
	if len(modelResults) > 0 {
		modelCompPath := filepath.Join(outputDir, "model_comparison.svg")
		err = PlotModelComparison(modelResults, modelCompPath)
		if err != nil {
			return fmt.Errorf("error plotting model comparison: %v", err)
		}
	}

	// 4. Plot feature importance (mock data for now)
	// In a real implementation, this would come from model analysis
	mockFeatureImportance := map[string]float64{
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

	featureImpPath := filepath.Join(outputDir, "feature_importance.svg")
	err = PlotFeatureImportance(mockFeatureImportance, featureImpPath)
	if err != nil {
		return fmt.Errorf("error plotting feature importance: %v", err)
	}

	return nil
}

// PlotFeatureDistribution creates a histogram showing the distribution of a numeric feature
func PlotFeatureDistribution(df dataframe.DataFrame, feature string, outputPath string) error {
	// Extract values from dataframe
	values := make([]float64, 0, df.Nrow())
	df.Col(feature).Map(func(e series.Element) series.Element {
		val, ok := e.Val().(float64)
		if ok {
			values = append(values, val)
		}
		return e
	})

	// Create bins for histogram
	if len(values) == 0 {
		return fmt.Errorf("no valid numeric values found for feature %s", feature)
	}

	// Find min and max values
	min, max := values[0], values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	// Create 10 bins
	numBins := 10
	binWidth := (max - min) / float64(numBins)
	if binWidth == 0 { // Handle case where all values are the same
		binWidth = 1
	}

	// Count values in each bin
	binCounts := make([]int, numBins)
	binLabels := make([]string, numBins)

	for _, v := range values {
		// Calculate bin index
		binIndex := int((v - min) / binWidth)
		// Handle edge case for max value
		if binIndex >= numBins {
			binIndex = numBins - 1
		}
		binCounts[binIndex]++
	}

	// Create labels for bins
	for i := 0; i < numBins; i++ {
		lowerBound := min + float64(i)*binWidth
		upperBound := min + float64(i+1)*binWidth
		binLabels[i] = fmt.Sprintf("%.2f-%.2f", lowerBound, upperBound)
	}

	// Create bars for the histogram
	var bars []chart.Value
	for i := 0; i < numBins; i++ {
		bars = append(bars, chart.Value{
			Value: float64(binCounts[i]),
			Label: binLabels[i],
			Style: chart.Style{
				FillColor:   blueColor.WithAlpha(180),
				StrokeColor: blueColor,
				StrokeWidth: 1,
			},
		})
	}

	// Create the chart
	histogram := chart.BarChart{
		Title:      fmt.Sprintf("Distribution of %s", feature),
		TitleStyle: chart.Style{FontSize: 14},
		Width:      600,
		Height:     400,
		BarWidth:   30,
		BarSpacing: 0, // Set to 0 to make bars touch like in a histogram
		XAxis:      chart.Style{},
		YAxis: chart.YAxis{
			Name:      "Frequency",
			NameStyle: chart.Style{FontSize: 12},
			Style:     chart.Style{FontSize: 10},
		},
		Bars: bars,
	}

	// Save the chart to file
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer f.Close()

	err = histogram.Render(chart.SVG, f)
	if err != nil {
		return fmt.Errorf("error rendering chart: %v", err)
	}

	return nil
}
