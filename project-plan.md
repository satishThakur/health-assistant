# Personal Health Assistant - Project Plan

## Constraints & Assumptions

- **Available Time**: 4-5 hours per week (~16-20 hours/month)
- **Timeline**: 6-8 months to functional MVP
- **Solo Developer**: No external dependencies
- **Scope**: Single user (personal use) for MVP

## Milestone Overview

| Milestone | Duration | Deliverable |
|-----------|----------|-------------|
| M1: Foundation | 3-4 weeks | Local dev environment, database, basic API |
| M2: Data Ingestion | 4-5 weeks | Garmin sync, manual logging working |
| M3: Flutter MVP | 4-5 weeks | Daily logging, dashboard, data viewing |
| M4: First Model | 4-6 weeks | Sleep quality model, basic insights |
| M5: Experiment Engine | 4-5 weeks | Design, track, analyze experiments |
| M6: First Experiment | 4 weeks | Run creatine experiment end-to-end |
| M7: Polish & Deploy | 3-4 weeks | AWS deployment, monitoring, refinement |

**Total Estimated Timeline**: 26-33 weeks (~6-8 months)

---

## Detailed Breakdown

### M1: Foundation (Weeks 1-4) - 16-20 hours

**Goal**: Local development environment running with basic API and database

#### Week 1 (4-5 hours)
- [ ] Initialize Git repo with README
- [ ] Set up Docker Compose:
  - PostgreSQL + TimescaleDB
  - MinIO (local S3)
- [ ] Create database schema (events, experiments, users tables)
- [ ] Test TimescaleDB hypertable creation

#### Week 2 (4-5 hours)
- [ ] Scaffold Go API Gateway with Chi/Fiber
- [ ] Basic JWT authentication (hardcoded user for now)
- [ ] Health check endpoint
- [ ] Database connection pooling (pgx)

#### Week 3 (4-5 hours)
- [ ] Build Data Service (Go):
  - POST /events endpoint
  - GET /events endpoint with time filters
- [ ] Write database seed script with sample data
- [ ] Test queries with sample time-series data

#### Week 4 (4-5 hours)
- [ ] Set up Python Model Service scaffold:
  - FastAPI with Docker
  - Basic health check endpoint
  - Connect to PostgreSQL
- [ ] Verify inter-service communication (Go → Python)
- [ ] Document API contracts in OpenAPI/Swagger

**M1 Success Criteria**:
- ✅ Can POST an event via API and query it back
- ✅ TimescaleDB queries running fast on sample data
- ✅ Docker Compose brings up full stack
- ✅ Basic API documentation exists

---

### M2: Data Ingestion (Weeks 5-9) - 20-25 hours

**Goal**: Garmin data flowing in, manual logging working, photos to S3

#### Week 5 (4-5 hours)
- [ ] Register Garmin Developer account
- [ ] Implement OAuth2 flow for Garmin API
- [ ] Test connection, fetch sample data
- [ ] Store OAuth tokens securely

#### Week 6 (4-5 hours)
- [ ] Build Ingestion Service (Go):
  - Garmin API client
  - Fetch sleep data
  - Fetch activity data
  - Write to events table
- [ ] Implement hourly cron job (can test with manual trigger)

#### Week 7 (4-5 hours)
- [ ] Integrate LLM API (OpenAI GPT-4V or Claude):
  - Photo upload endpoint
  - Send to LLM for macro extraction
  - Parse response into structured data
- [ ] Store photo in MinIO/S3
- [ ] Create meal event with macros + photo URL

#### Week 8 (4-5 hours)
- [ ] Manual data entry endpoints:
  - POST subjective feelings
  - POST supplement log
  - PATCH meal macros (manual correction)
- [ ] Validation logic for all inputs

#### Week 9 (4-5 hours)
- [ ] Test full ingestion pipeline:
  - Garmin sync with real data
  - Photo upload and macro extraction
  - Manual logging
- [ ] Error handling and retry logic
- [ ] Basic logging/monitoring

**M2 Success Criteria**:
- ✅ Real Garmin data syncing hourly into database
- ✅ Can upload meal photo and get macros back
- ✅ Can manually log feelings and supplements
- ✅ All data visible in database with correct schema

---

### M3: Flutter MVP (Weeks 10-14) - 20-25 hours

**Goal**: Working Flutter app for daily logging and viewing data

#### Week 10 (4-5 hours)
- [ ] Flutter project setup (mobile + web target)
- [ ] Set up state management (Riverpod or Bloc)
- [ ] Configure API client (Dio) with auth interceptor
- [ ] Basic navigation structure (bottom nav bar)

#### Week 11 (4-5 hours)
- [ ] Build Daily Log screen:
  - Subjective feeling sliders (energy, mood, focus, physical)
  - Notes text field
  - Submit button
- [ ] Connect to POST /events/subjective API
- [ ] Simple success/error feedback

#### Week 12 (4-5 hours)
- [ ] Build Meal Logging screen:
  - Camera integration (image_picker)
  - Upload to API
  - Display extracted macros
  - Allow editing macros
- [ ] Loading states while LLM processes

#### Week 13 (4-5 hours)
- [ ] Build Supplement Logging screen:
  - Checklist of supplements
  - Timestamp recording
  - Mark as taken/missed
- [ ] Build Dashboard screen:
  - Today's date
  - Garmin sync status
  - Key metrics (sleep score, HRV, steps)
  - Mock visualizations (charts)

#### Week 14 (4-5 hours)
- [ ] Build Timeline/Data Explorer screen:
  - List view of all events
  - Filter by date range and type
  - Basic search
- [ ] Polish UI/UX
- [ ] Test on physical device

**M3 Success Criteria**:
- ✅ Can log subjective feelings 2x/day via app
- ✅ Can upload meal photo and see/edit macros
- ✅ Can log supplements
- ✅ Dashboard shows real Garmin data
- ✅ App feels responsive and usable

---

### M4: First Model (Weeks 15-20) - 24-30 hours

**Goal**: Sleep quality prediction model running, generating insights

#### Week 15 (4-5 hours)
- [ ] Data exploration in Jupyter:
  - Load events from database
  - Visualize time series
  - Check data quality and gaps
- [ ] Feature engineering notebook:
  - Lag features (yesterday's values)
  - Rolling averages
  - Time of day encoding

#### Week 16 (5-6 hours)
- [ ] Build first PyMC model:
  - Simple linear model for sleep quality
  - Predictors: HRV, exercise, meal timing
  - Prior specifications
- [ ] Train on existing data
- [ ] Validate with posterior predictive checks

#### Week 17 (5-6 hours)
- [ ] Implement time-series model:
  - Autoregressive component
  - Time-lagged effects
  - Hierarchical priors
- [ ] Compare models (DIC/WAIC)

#### Week 18 (4-5 hours)
- [ ] Model Service API endpoints:
  - POST /models/sleep/train (trigger retraining)
  - GET /models/sleep/predict (next-day prediction)
  - GET /models/sleep/feature-importance
- [ ] Serialize and save trained models

#### Week 19 (4-5 hours)
- [ ] Build insights generation:
  - Compute time-lagged correlations
  - Identify strongest predictors
  - Format as human-readable insights
- [ ] Endpoint: GET /insights

#### Week 20 (4-5 hours)
- [ ] Flutter: Insights screen
  - Display top correlations
  - "Your sleep quality is most affected by..."
  - Visualizations (bar charts, time series)
- [ ] Test end-to-end prediction flow

**M4 Success Criteria**:
- ✅ Sleep quality model trained on personal data
- ✅ Can get next-day sleep prediction with confidence interval
- ✅ Insights visible in Flutter app
- ✅ Model retrains weekly with new data

---

### M5: Experiment Engine (Weeks 21-25) - 20-25 hours

**Goal**: System can design, track, and analyze experiments

#### Week 21 (4-5 hours)
- [ ] Design experiment proposal logic:
  - Parse model insights for intervention candidates
  - Generate experiment specs (duration, metrics to track)
  - Store as "proposed" experiments
- [ ] Endpoint: GET /experiments/proposals

#### Week 22 (4-5 hours)
- [ ] Flutter: Experiments screen
  - List proposed experiments (swipe cards?)
  - Show hypothesis, design, duration
  - Accept/Reject buttons
- [ ] Endpoint: POST /experiments/:id/accept

#### Week 23 (4-5 hours)
- [ ] Experiment tracking logic:
  - Mark experiment as "active"
  - Check compliance (are supplements being logged?)
  - Send reminders if needed
- [ ] Flutter: Active experiment card on dashboard

#### Week 24 (5-6 hours)
- [ ] Experiment analysis (Python):
  - Fetch data for experiment period
  - Compare intervention vs baseline/control
  - Bayesian estimation of effect size
  - Credible intervals, posterior probabilities
- [ ] Endpoint: POST /experiments/:id/analyze

#### Week 25 (4-5 hours)
- [ ] Flutter: Experiment results screen
  - Show statistical outcomes
  - Visualizations (posterior distributions)
  - "Probability of positive effect: 87%"
- [ ] Mark experiment as "completed"

**M5 Success Criteria**:
- ✅ System proposes experiment based on insights
- ✅ Can accept experiment and track compliance
- ✅ Experiment analyzed with Bayesian stats
- ✅ Results visible in app with uncertainty

---

### M6: First Experiment (Weeks 26-29) - Real-world test

**Goal**: Run the creatine/protein experiment from start to finish

#### Weeks 26-29 (4 weeks of data collection)
- [ ] Accept experiment: "Creatine effect on recovery"
- [ ] Log compliance daily:
  - Creatine intake (5g post-workout)
  - Workout performance (weight, reps, RPE)
  - Recovery metrics (HRV, soreness, body battery)
- [ ] Weeks 1-2: Creatine + protein
- [ ] Weeks 3-4: Control (no creatine)

**Activities during this phase** (5-10 hours total):
- [ ] Monitor data quality
- [ ] Fix any bugs in logging flow
- [ ] Ensure Garmin sync staying reliable
- [ ] Check that experiment tracking works
- [ ] Week 4: Run analysis
- [ ] Review results, document learnings

**M6 Success Criteria**:
- ✅ Completed 4-week experiment with >90% compliance
- ✅ Analysis shows clear effect size estimates
- ✅ Insights actionable (continue/stop creatine)
- ✅ System proved end-to-end value

---

### M7: Polish & Deploy (Weeks 30-33) - 16-20 hours

**Goal**: Production deployment on AWS, monitoring, documentation

#### Week 30 (4-5 hours)
- [ ] AWS setup:
  - RDS PostgreSQL with TimescaleDB
  - S3 bucket for photos
  - ECR for Docker images
- [ ] Provision infrastructure (Terraform or manual)

#### Week 31 (4-5 hours)
- [ ] Build CI/CD pipeline:
  - GitHub Actions for tests
  - Build and push Docker images to ECR
  - Deploy to ECS Fargate
- [ ] Environment variables and secrets management

#### Week 32 (4-5 hours)
- [ ] Deploy Flutter web app:
  - Build for web
  - Host on S3 + CloudFront or Amplify
- [ ] Test production deployment end-to-end
- [ ] Fix deployment issues

#### Week 33 (4-5 hours)
- [ ] Set up monitoring:
  - CloudWatch alarms (API errors, DB connections)
  - Grafana dashboard for key metrics
- [ ] Write operational runbook
- [ ] Final documentation update
- [ ] Celebrate! 🎉

**M7 Success Criteria**:
- ✅ System running on AWS, accessible from anywhere
- ✅ Automated deployments working
- ✅ Monitoring and alerts configured
- ✅ Documentation up to date

---

## Risk Mitigation

### Technical Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Garmin API rate limits | Medium | High | Cache aggressively, batch requests, fall back to manual entry |
| LLM macro accuracy too low | Medium | Medium | Allow manual override, collect ground truth for fine-tuning |
| Not enough data for models | High | High | Design focused experiments, use strong priors, accept high uncertainty initially |
| Time overruns | High | Medium | Cut scope aggressively, focus on one experiment well |
| Burnout (solo project) | Medium | High | Keep milestones small, celebrate wins, take breaks |

### Schedule Risks

**If falling behind**:
1. **Cut scope first**: Remove "nice-to-haves" (e.g., skip Flutter web, mobile only)
2. **Simplify models**: Use simpler regression instead of full Bayesian hierarchy
3. **Manual workflows**: Skip automation, do some steps manually until working
4. **Delay deployment**: Run locally longer, deploy later

**Buffer**: 2-3 weeks built into timeline for unknowns

---

## Weekly Workflow Recommendation

Given limited time, suggest this pattern:

**Weekdays (1-2 hours, 2-3 times/week)**:
- Focused coding sessions
- Complete one small task per session
- Push code frequently

**Weekend (2-3 hour block)**:
- Larger integration work
- Testing end-to-end flows
- Planning next week's tasks

**Monthly review** (1 hour):
- Assess progress vs plan
- Adjust timeline if needed
- Celebrate milestones

---

## Definition of Done (MVP)

The project is "done" when:

1. ✅ All data sources integrated (Garmin, manual, photos)
2. ✅ At least 8 weeks of consistent data collected
3. ✅ One predictive model deployed (sleep quality)
4. ✅ One experiment completed with statistical analysis (creatine)
5. ✅ Flutter app functional for daily use
6. ✅ System running in production (AWS)
7. ✅ Actionable insights generated from personal data

At that point, the system proves the core hypothesis: **causal inference from n=1 data is possible and valuable**.

---

## Post-MVP Roadmap (Future)

Once MVP is solid, consider:

- Additional models (energy prediction, workout performance)
- Multi-armed bandit for supplement optimization
- Integration with other wearables (Oura, Apple Watch)
- Voice logging (Alexa/Google Assistant)
- PDF lab result parsing automation
- Sharing insights with doctor
- Open source components (experiment engine library)

But don't even think about these until MVP delivers value!

---

## Success Metrics

**Leading Indicators** (process):
- Consistent 4-5 hours/week logged
- Weekly commits to repo
- Milestones completed within 1 week of estimate

**Lagging Indicators** (outcome):
- 90%+ Garmin sync success rate
- Daily logging compliance >85%
- At least one significant insight discovered
- One intervention made based on data

**Ultimate Success**:
- You trust the system enough to make health decisions from it
- You use it daily without friction
- You've learned something genuinely new about your health

Good luck! 🚀
