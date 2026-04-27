# 🛰️ Orbital Debris Tracker

## Description
The Orbital Debris Tracker is a high-performance, full-stack application designed to ingest, predict, and visualize the trajectories of thousands of satellites and space debris in Earth's orbit. Utilizing real-time Two-Line Element (TLE) data, the system relies on a Go-based compute engine to run complex SGP4 orbital propagation and spatial partitioning algorithms. The results are served via a Python FastAPI gateway to a highly optimized React/WebGL frontend, allowing operators to monitor conjunction risks and visualize 30,000+ orbital objects in a 3D space environment simultaneously.

## 📑 Table of Contents
- [Features](#-features)
- [Technologies Used](#-technologies-used)
- [Installation](#-installation)
- [Usage](#-usage)
- [Project Structure](#-project-structure)
- [Contributing](#-contributing)
- [License](#-license)

## 🚀 Features
* **Automated Data Ingestion:** A Go-based worker process automatically fetches and parses raw 3-line text files from Space-Track.org on a scheduled cron-like basis.
* **High-Performance Compute Engine:** Predicts orbital locations over time using SGP4 propagation algorithms. Implements 3D spatial partitioning (Octrees) to avoid O(N²) collision checks.
* **Conjunction Risk Calculation:** Calculates Probability of Collision (Pc), Miss Distance, and Time of Closest Approach (TCA) for early warning alerts.
* **FastAPI Gateway:** A robust client-facing API backend that utilizes Redis for instantaneous orbital position lookups and gRPC to communicate with the compute engine.
* **3D Orbital Command Center:** A React and TypeScript frontend that leverages WebGL (Three.js/React Three Fiber) with custom GLSL shaders and instanced meshes to render 30,000+ points smoothly.
* **Scalable Infrastructure:** Production-ready Kubernetes manifests are included to dynamically scale CPU-heavy propagation pods and handle massive datasets.

## 🛠️ Technologies Used
* **Compute Engine & Ingestion:** Go, gRPC, SGP4 Math
* **Backend API:** Python, FastAPI, Redis
* **Frontend:** React, TypeScript, Three.js, WebGL/GLSL
* **Database:** PostgreSQL (Historical Catalog), Redis (In-Memory Positions)
* **Infrastructure:** Docker, Docker Compose, Kubernetes

## ⚙️ Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/mtepenner/orbital-debris-tracker.git
   cd orbital-debris-tracker
   ```

2. (Optional) Compile protobufs if developing the Go engine:
   ```bash
   make proto
   ```

3. Spin up the entire stack locally (Databases, Go Engine, FastAPI, and React) using Docker Compose:
   ```bash
   docker-compose up -d
   ```

## 💻 Usage

* **Access the 3D Visualizer:** Open your browser and navigate to `http://localhost:3000` (or your configured port) to view the Globe Visualizer and debris cloud.
* **Monitor Conjunctions:** Use the "Conjunction Alerts" sidebar in the UI to view real-time feeds of upcoming close approaches and assess the calculated collision risks.
* **Filter Catalog Data:** Click on individual objects within the WebGL canvas to open the "Object Inspector" and view specific satellite or debris telemetry.

## 📂 Project Structure
* `/tle_ingestion`: Go worker for syncing and parsing TLE data from Space-Track.org.
* `/compute_engine`: Go/C++ engine performing rigorous SGP4 propagation and risk calculations.
* `/api_gateway`: Python/FastAPI client-facing server managing catalog routes and conjunction streams.
* `/frontend`: React/TypeScript 3D visualization dashboard using custom WebGL shaders.
* `/k8s`: Kubernetes manifests for deploying cronjobs, compute pods, databases, and APIs in production.
* `/.github/workflows`: CI/CD pipelines including unit tests for SGP4 math and deployment automation.

## 🤝 Contributing
Contributions are highly encouraged! Whether it's optimizing a GLSL shader, improving the Go spatial partitioning, or expanding the API, feel free to open a Pull Request. Please ensure that all SGP4 propagation and Octree math unit tests pass before submitting.

## 📄 License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
