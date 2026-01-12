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

## The Technical Vision

### Proving Two Hypotheses

**Product Hypothesis**: N=1 causal health optimization is both possible and valuable with modern AI/ML tools.

**Engineering Hypothesis**: GenAI dramatically lowers the barrier to building in adjacent domains. An experienced backend engineer can ship a production-quality Flutter app, train Bayesian models, and integrate LLMs - solo, in months not years.

### Tech Stack

- **Frontend**: Flutter (mobile + web, single codebase)
- **Backend**: Go microservices (API gateway, data service, experiment engine, ingestion)
- **ML/Models**: Python (PyMC for Bayesian models, FastAPI for serving)
- **Database**: PostgreSQL + TimescaleDB (time-series optimization)
- **Storage**: AWS S3 (meal photos, lab PDFs)
- **AI/LLM**: Claude/GPT-4V for meal analysis, health report parsing, reasoning
- **Infra**: Docker, AWS (ECS, RDS, S3)

### Architecture Philosophy

- **Microservices**: Clean separation of concerns (data, models, experiments, ingestion)
- **API-first**: REST endpoints for all interactions
- **Time-series native**: TimescaleDB for efficient time-based queries
- **Ultra hands-on**: Full access to raw data, model diagnostics, experiment design
- **Statistically rigorous**: Bayesian methods, credible intervals, proper causal inference

See [highleveldesign.md](./highleveldesign.md) for full architecture details.

---

## Project Documentation

- **[idea.md](./idea.md)**: Full product vision, problem statement, use cases, success criteria
- **[highleveldesign.md](./highleveldesign.md)**: System architecture, tech stack, data model, API design
- **[project-plan.md](./project-plan.md)**: 6-8 month roadmap with realistic milestones (4-5 hours/week constraint)

---

## Current Status

**Phase**: Planning & Documentation Complete ✅

**Next Steps**:
- M1: Foundation (Weeks 1-4) - Set up local dev environment, database, basic API
- Begin Garmin API integration
- Start collecting personal data

**Timeline**: MVP in 6-8 months (targeting functional system by mid-2026)

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

## Why Open Source This? (Future)

Once the MVP is proven, the plan is to open source components:

- **Experiment engine**: Reusable library for designing and analyzing n=1 experiments
- **Bayesian health models**: Templates for common health predictions (sleep, energy, recovery)
- **LLM integration patterns**: How to use AI for health data parsing and reasoning

The world needs more tools for rigorous self-experimentation. If this works for me, the components can help others.

But first: build it, use it, prove it works.

---

## Getting Started (Coming Soon)

Setup instructions will be added once M1 is complete. For now, see [project-plan.md](./project-plan.md) for the development roadmap.

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

## License

TBD (likely MIT once open sourced)

---

## Contact

This is a personal project by an experienced engineer exploring the intersection of health optimization, causal inference, and GenAI-accelerated development.

If this resonates with you or you're working on similar problems, reach out. Always interested in discussing n=1 experimentation, Bayesian methods, or the future of personalized health.

---

**Last Updated**: January 2026
**Status**: Planning Complete, Ready to Build 🚀
