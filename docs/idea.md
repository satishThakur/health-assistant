# Personal Health Assistant - Product Vision

## Problem Statement

Current health tracking apps are glorified data collectors. They show you charts and trends but don't help you understand **why** you feel good some days and terrible others. More critically, they don't help you **intervene** - to systematically test what actually works for **your** unique physiology.

### The Core Problems

1. **Data Silos**: Wearable data (Garmin), lab results, subjective feelings, nutrition, and supplements live in separate places
2. **Correlation vs Causation**: Apps show correlations but can't distinguish what's actually affecting your health vs coincidence
3. **Generic Advice**: Recommendations are population-level, not personalized to your n=1 data
4. **No Experimentation Framework**: No systematic way to test interventions (e.g., "Does magnesium actually improve my sleep?")

### The Opportunity

Build a **personalized health assistant** that:
- Aggregates all health data into a unified timeline
- Uses Bayesian hierarchical models to discover causal relationships in your data
- Suggests evidence-based experiments tailored to you
- Tracks intervention outcomes and updates beliefs
- Provides ultra hands-on control for power users

## Target User

**Primary**: Self (experienced engineer, quantified-self enthusiast, learning Bayesian modeling)

**Secondary**: Health-conscious individuals who want to move beyond tracking to optimization

## Data Sources

### Objective Metrics
1. **Wearable Data (Garmin)**
   - Sleep quality, stages, HRV
   - Activity: steps, workouts, intensity
   - Body battery, stress scores
   - Heart rate zones, VO2 max
   - Sync frequency: Hourly

2. **Lab Results**
   - Blood panels (CBC, metabolic, hormones, vitamins)
   - Parsed from PDFs/images via LLM
   - Stored with reference ranges

### Subjective Metrics
3. **Daily Feelings** (2x/day: morning + evening)
   - Energy level
   - Mood
   - Focus/cognition
   - Physical feeling
   - (Each on numeric scale with optional notes)

### Behavioral Data
4. **Nutrition**
   - Meal photos → LLM extracts macros
   - Manual entry option
   - Timing and composition

5. **Supplements**
   - What was taken, when
   - Track compliance vs planned regimen
   - Flag missed doses

## Use Cases

### Phase 1: Observe & Describe
- "What are my baseline patterns?"
- "When do I naturally sleep best?"
- "How variable is my HRV?"

### Phase 2: Discover Correlations
- "What predicts good sleep quality?"
- "Does meal timing affect next-day energy?"
- "What patterns exist between workouts and recovery?"

### Phase 3: Causal Experiments
- **Example**: "Does creatine supplementation improve workout performance and recovery?"
  - Track: gym performance (weight/reps), soreness ratings, HRV recovery
  - Design: 4 weeks on, 2 weeks washout, 4 weeks on
  - Analysis: Bayesian estimation of effect size with uncertainty

### Phase 4: Optimize & Intervene
- "Given my current state, what intervention would most improve sleep?"
- Multi-armed bandit for supplement stack optimization
- Personalized protocol recommendations

## Success Criteria

### MVP Success
- All data sources integrated and flowing into unified timeline
- At least one causal model running (sleep quality prediction)
- First experiment designed, executed, and analyzed
- Flutter app functional for daily logging and viewing insights

### Long-term Success
- Discovered 3+ genuine causal relationships in personal data
- Successfully optimized at least one health metric through intervention
- Models update in real-time as new data arrives
- System provides confidence intervals, not false certainty

## Non-Goals (For MVP)

- Social features, sharing, community
- Mobile app distribution (App Store/Play Store)
- Multi-user support
- Real-time alerting
- Integration with medical records/providers

## Key Risks & Mitigations

| Risk | Mitigation |
|------|------------|
| Too little data for causal inference | Start with strong priors, design focused experiments |
| Garmin API rate limits | Batch requests, cache aggressively |
| LLM hallucination on food macros | Always show confidence, allow manual override |
| Overfitting on n=1 data | Use hierarchical models, regularization, cross-validation where possible |
| Experiment fatigue | Keep experiments short (2-4 weeks), limit to 1 active at a time |

## Guiding Principles

1. **Causation over Correlation**: Build for experimentation, not just observation
2. **Uncertainty as Feature**: Show confidence intervals, admit what we don't know
3. **User in Control**: Ultra hands-on access to data, models, raw outputs
4. **Incremental Value**: Each phase delivers usable insights, not just infrastructure
5. **Statistical Rigor**: Proper Bayesian methods, no p-hacking

## First Hypothesis to Test

**"Does creatine and/or whey protein supplementation improve workout performance and recovery time?"**

**Measurable Outcomes**:
- Workout performance (weight lifted, reps, perceived exertion)
- Recovery metrics (HRV normalization time, muscle soreness ratings, body battery recovery)
- Subjective energy and physical feeling scores

**Experimental Design**:
- Factorial design: creatine (yes/no) × protein timing (post-workout vs distributed)
- 12-week total: 4 weeks per condition
- Control for: workout intensity, sleep, overall calories
