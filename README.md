# Credit Card Approval Prediction

## Overview

This project implements a machine learning model to predict whether a credit card application will be approved or rejected based on applicant information. It uses the UCI Credit Card Approval dataset and implements various classification algorithms to identify patterns in approval decisions.

## Dataset

The dataset contains 690 credit card applications with 15 attributes plus a class attribute indicating approval status. All attribute names and values have been anonymized to protect confidentiality. The dataset includes:

- A mix of continuous and categorical features
- 44.5% approved applications and 55.5% rejected applications
- 5% of cases have missing values

## Features

- Data preprocessing and cleaning
- Exploratory data analysis with visualizations
- Implementation of multiple classification algorithms
- Model evaluation and comparison
- Feature importance analysis

## Project Structure

```
Credit Card Approval Prediction/
├── cmd/
│   └── main.go   
├── data/
│   ├── raw/
│   └── processed/
├── internal/
│   ├── models/
│   ├── preprocessing/
│   ├── evaluation/
│   └── visualization/
├── notebooks/
├── tests/
├── go.mod     
├── go.sum 
├── README.md  
└── project_proposal.md 
```

## Requirements

- Go
- Required Go packages
  - github.com/go-gota/gota
  - github.com/wcharczuk/go-chart/v2
  - gonum.org/v1/gonum

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/jimmymcguigan18/credit-card-approval-prediction.git
   cd credit-card-approval-prediction
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

## Usage

1. Run the main application:
   ```bash
   go run cmd/main.go
   ```

2. For specific tasks:
   ```bash
   # Preprocess data only
   go run cmd/main.go --preprocess
   
   # Train models only
   go run cmd/main.go --train
   
   # Evaluate models only
   go run cmd/main.go --evaluate
   
   # Generate visualizations
   go run cmd/main.go --visualize
   ```

## Model Performance

*Note: This section will be updated after model implementation and evaluation.*

Preliminary results show:

| Model | Accuracy | Precision | Recall | F1 Score |
|-------|----------|-----------|--------|----------|
| Logistic Regression | TBD | TBD | TBD | TBD |
| Random Forest | TBD | TBD | TBD | TBD |
| Gradient Boosting | TBD | TBD | TBD | TBD |
| SVM | TBD | TBD | TBD | TBD |
| Neural Network | TBD | TBD | TBD | TBD |

---
