# Personal Health Assistant

> **Building a personalized health optimization system that discovers causal relationships and suggests evidence-based interventions for N=1.**

## The Rationale: Why Now?

### The GenAI Inflection Point

We're at a unique moment in technology where **personalized AI applications for N=1 use cases are becoming reality**. Historically, building sophisticated health tracking and analysis systems required large teams and significant resources. The advent of generative AI changes everything:

1. **Perception as a Service**: Computer vision models can analyze meal photos and extract nutritional information. LLMs can parse complex medical PDFs. What once required specialized ML expertise is now an API call.

2. **Reasoning Augmentation**: AI can help design experiments, interpret statistical results, and generate insights. The barrier to building intelligent systems has collapsed.

3. **Cross-Domain Velocity**: GenAI enables experienced engineers to rapidly build in adjacent domains. This project is a proof point: a 20+ year backend engineer building a full-stack Flutter app, Bayesian models, and ML pipelines - all with AI as a force multiplier.

**Hypothesis**: In the GenAI era, solo developers can build production-quality personalized health systems that were previously impossible outside research labs or well-funded startups.

This project aims to prove that hypothesis.

### The Personal Motivation

As someone who:
- Goes to the gym regularly
- Tries to eat right and optimize nutrition
- Tracks fitness data via Garmin
- Takes supplements but doesn't know what actually works
- Wants to feel better, perform better, and understand my own physiology

I have all the ingredients for a personalized health optimization system **except the system itself**.

Generic health advice doesn't work. Population studies don't account for my unique genetics, lifestyle, and goals. I need answers to questions like:

- Does creatine actually improve my recovery, or is it placebo?
- What meal timing maximizes my energy and focus?
- Which supplements move the needle on my biomarkers?
- How does sleep quality really affect my workout performance?

The only way to answer these questions is through **rigorous n=1 experimentation with proper causal inference**.

So I'm building it.

---

## What Is This Project?

A **fully personalized health assistant** that:

1. **Aggregates** all my health data:
   - Wearable metrics (Garmin: sleep, HRV, activity, stress)
   - Lab results (blood panels, biomarkers)
   - Subjective feelings (energy, mood, focus)
   - Nutrition (meal photos → macros via AI)
   - Supplements (what I take, when, compliance)

2. **Discovers** causal relationships:
   - Time-series Bayesian models to identify what predicts outcomes
   - Move beyond correlation to causation through experiments
   - Builds confidence with hierarchical models and uncertainty quantification

3. **Suggests** evidence-based interventions:
   - Proposes experiments: "Test creatine for 4 weeks, measure recovery metrics"
   - Tracks compliance and outcomes
   - Analyzes results with proper statistics
   - Updates beliefs as data accumulates

4. **Optimizes** continuously:
   - Multi-armed bandit approaches for supplement stacks
   - Adaptive experiment design
   - Long-term tracking of what works for me

This is not another passive tracking app. This is an **active experimentation platform** for self-optimization.

---

## Project Structure

```
health-assistant/
├── docs/                        # Documentation
│   ├── README.md               # Detailed project vision (this file moved)
│   ├── idea.md                 # Product vision and problem statement
│   ├── highleveldesign.md      # System architecture and tech stack
│   └── project-plan.md         # 6-8 month roadmap
│
├── backend/                     # Go backend (single module, multiple binaries)
│   ├── cmd/                    # Service entry points
│   │   ├── api-gateway/        # Main API gateway (port 8080)
│   │   ├── data-service/       # CRUD operations (port 8081)
│   │   ├── experiment-service/ # Experiment engine (port 8082)
│   │   └── ingestion-service/  # Garmin sync, photo processing (port 8083)
│   ├── internal/               # Private application code
│   │   ├── api/               # HTTP handlers
│   │   ├── db/                # Database layer
│   │   ├── auth/              # JWT authentication
│   │   ├── models/            # Domain models
│   │   ├── garmin/            # Garmin API client
│   │   ├── llm/               # LLM integrations
│   │   └── config/            # Configuration
│   └── go.mod                 # Go module definition
│
├── model-service/              # Python ML service
│   ├── app/
│   │   ├── main.py            # FastAPI application
│   │   ├── models/            # PyMC Bayesian models
│   │   └── api/               # API routes
│   ├── notebooks/             # Jupyter notebooks for exploration
│   ├── requirements.txt
│   └── Dockerfile
│
├── app/                        # Flutter application
│   └── health_assistant/      # Flutter project (to be created)
│       ├── lib/               # Dart source code
│       ├── test/              # Tests
│       └── pubspec.yaml       # Dependencies
│
├── infra/                      # Infrastructure
│   ├── docker-compose.yml     # Local development stack
│   ├── docker/                # Additional Docker configs
│   └── terraform/             # AWS infrastructure (future)
│
├── scripts/                    # Helper scripts
│   ├── db/
│   │   ├── init.sql           # Database schema
│   │   ├── migrations/        # SQL migrations
│   │   └── seed.sql           # Sample data
│   └── bin/                   # Utility scripts
│
├── .gitignore
└── README.md                   # This file
```

---

## Tech Stack

### Backend
- **Language**: Go 1.22+
- **Architecture**: Single module with multiple service binaries
- **Database**: PostgreSQL + TimescaleDB (time-series optimization)
- **Storage**: AWS S3 (meal photos, lab PDFs)
- **Auth**: JWT tokens

### Model Service
- **Language**: Python 3.11+
- **Framework**: FastAPI
- **ML**: PyMC (Bayesian models), NumPy, Pandas
- **Visualization**: Matplotlib, Plotly

### Frontend
- **Framework**: Flutter (Dart)
- **Platforms**: Mobile (iOS, Android) + Web
- **State Management**: Riverpod or Bloc

### Infrastructure
- **Containerization**: Docker, Docker Compose
- **Cloud**: AWS (ECS, RDS, S3)
- **Local Development**: PostgreSQL + TimescaleDB, MinIO (S3-compatible)

---

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.22+ (for backend development)
- Python 3.11+ (for model service development)
- Flutter SDK (for mobile app development)

### Local Development Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/satishThakur/health-assistant.git
   cd health-assistant
   ```

2. **Set up environment variables**
   ```bash
   cd infra
   cp .env.example .env
   # Edit .env with your API keys (Garmin, OpenAI, etc.)
   ```

3. **Start the infrastructure**
   ```bash
   cd infra
   docker-compose up -d
   ```

   This will start:
   - PostgreSQL + TimescaleDB (port 5432)
   - MinIO (S3-compatible storage, ports 9000, 9001)
   - API Gateway (port 8080)
   - Data Service (port 8081)
   - Experiment Service (port 8082)
   - Ingestion Service (port 8083)
   - Model Service (port 8084)

4. **Verify services are running**
   ```bash
   curl http://localhost:8080/health  # API Gateway
   curl http://localhost:8084/health  # Model Service
   ```

5. **Access MinIO Console**
   - URL: http://localhost:9001
   - Username: `minioadmin`
   - Password: `minioadmin`

### Running Services Individually (for development)

**Backend Services (Go)**:
```bash
cd backend

# Run a specific service
go run ./cmd/api-gateway
go run ./cmd/data-service
go run ./cmd/experiment-service
go run ./cmd/ingestion-service
```

**Model Service (Python)**:
```bash
cd model-service
pip install -r requirements.txt
python app/main.py
```

**Flutter App**:
```bash
cd app/health_assistant  # After running flutter create
flutter pub get
flutter run
```

---

## Documentation

- **[docs/idea.md](./docs/idea.md)**: Full product vision, problem statement, use cases
- **[docs/highleveldesign.md](./docs/highleveldesign.md)**: System architecture, data model, API design
- **[docs/project-plan.md](./docs/project-plan.md)**: 6-8 month roadmap with milestones
- **[backend/README.md](./backend/README.md)**: Backend service documentation
- **[model-service/README.md](./model-service/README.md)**: Model service documentation
- **[app/README.md](./app/README.md)**: Flutter app documentation

---

## Current Status

**Phase**: Foundation Setup Complete ✅

### Completed
- ✅ Project structure defined
- ✅ Backend skeleton (Go services)
- ✅ Model service skeleton (Python/FastAPI)
- ✅ Database schema (PostgreSQL + TimescaleDB)
- ✅ Docker Compose for local development
- ✅ Domain models and configuration
- ✅ Sample data seed scripts

### Next Steps
- [ ] Initialize Flutter app
- [ ] Implement database connection in backend
- [ ] Add JWT authentication
- [ ] Build first API endpoints
- [ ] Integrate Garmin API
- [ ] Add LLM integration for meal analysis
- [ ] Implement first Bayesian model (sleep quality)

See [docs/project-plan.md](./docs/project-plan.md) for detailed roadmap.

---

## The First Experiment

To validate the entire system, the first experiment will be:

**"Does creatine and/or whey protein supplementation improve workout performance and recovery?"**

**Metrics to track**:
- Workout performance (weight, reps, perceived exertion)
- Recovery (HRV normalization time, muscle soreness, body battery)
- Subjective energy and physical feeling scores

**Design**: 12-week factorial experiment with proper controls and washout periods

**Expected outcome**: Bayesian posterior distributions showing effect sizes with uncertainty. Finally know if creatine actually works for me.

---

## Tech Hypothesis Tracker

As I build this project, I'll track how GenAI accelerates development in domains where I have less experience:

| Domain | Prior Experience | GenAI Accelerator | Outcome |
|--------|------------------|-------------------|---------|
| **Flutter/Dart** | Minimal (backend-focused) | Claude Code for UI components, state management | TBD |
| **Bayesian Modeling** | Learning in progress | LLM for PyMC code, model debugging | TBD |
| **Time-series ML** | Basic understanding | AI for feature engineering, model selection | TBD |
| **AWS Deployment** | Experienced, but rusty | AI for Terraform, troubleshooting | TBD |
| **UI/UX Design** | Weak | AI for design suggestions, Flutter widgets | TBD |

**Expected Result**: 3-5x faster development in adjacent domains compared to traditional learning curve.

Will update this table as the project progresses.

---

## Contributing

This is currently a personal project for learning and experimentation. Once the MVP is proven, components may be open sourced for community contribution.

---

## License

TBD (likely MIT once open sourced)

---

## Contact

This is a personal project by an experienced engineer exploring the intersection of health optimization, causal inference, and GenAI-accelerated development.

If this resonates with you or you're working on similar problems, feel free to open an issue or reach out.

---

**Last Updated**: January 2026
**Status**: Foundation Complete, Ready to Build 🚀
