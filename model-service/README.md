# Model Service

Python-based machine learning and statistical modeling service using PyMC for Bayesian inference.

## Purpose

This service provides:
- **Bayesian hierarchical models** for health predictions (sleep quality, energy levels, recovery)
- **Time-lagged correlation analysis** to identify causal relationships
- **Experiment analysis** using Bayesian methods
- **Insight generation** from personal health data

## Tech Stack

- **FastAPI**: Web framework for API endpoints
- **PyMC**: Probabilistic programming for Bayesian models
- **NumPy/Pandas**: Data manipulation
- **Arviz**: Bayesian model diagnostics and visualization
- **PostgreSQL**: Database connection for fetching data

## Setup

### Install Dependencies

```bash
cd model-service
pip install -r requirements.txt
```

Or using virtual environment:
```bash
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
pip install -r requirements.txt
```

### Run Development Server

```bash
cd model-service
python -m uvicorn app.main:app --reload --host 0.0.0.0 --port 8084
```

Or:
```bash
python app/main.py
```

## API Endpoints

### Health Check
```
GET /health
```

### Sleep Quality Prediction
```
POST /models/sleep-quality/predict
{
  "user_id": "user-123",
  "features": {
    "hrv": 65.0,
    "exercise_duration": 45.0,
    "meal_timing": 3.5,
    "supplement_taken": 1,
    "sleep_quality_lag1": 75.0
  }
}
```

### Correlation Analysis
```
POST /models/correlations
{
  "user_id": "user-123",
  "target_metric": "sleep_quality",
  "start_date": "2026-01-01",
  "end_date": "2026-01-31"
}
```

### Insights
```
GET /models/insights/{user_id}
```

## Models

### Sleep Quality Model (`app/models/sleep_quality.py`)

Hierarchical Bayesian model that predicts sleep quality based on:
- HRV (Heart Rate Variability)
- Exercise duration and timing
- Meal timing
- Supplement intake
- Autoregressive component (previous day's sleep)

**Model Structure**:
```python
sleep_quality ~ Normal(μ, σ)
μ = β₀ + β₁·HRV + β₂·exercise + β₃·meal_timing + β₄·supplement + ρ·sleep_lag1
```

### Future Models

- Energy Level Prediction
- Workout Performance Prediction
- Recovery Time Estimation

## Notebooks

The `notebooks/` directory contains Jupyter notebooks for:
- Exploratory data analysis
- Model prototyping
- Experiment design
- Visualization

## Testing

```bash
pytest
```

## Docker

Build the container:
```bash
docker build -t health-assistant-model-service .
```

Run:
```bash
docker run -p 8084:8084 health-assistant-model-service
```

## Next Steps

- [ ] Implement actual PyMC models (currently placeholders)
- [ ] Add database connection for fetching user data
- [ ] Implement posterior predictive sampling
- [ ] Add model versioning and persistence
- [ ] Create experiment analysis endpoints
- [ ] Add comprehensive tests
