[![Live Demo](https://img.shields.io/badge/demo-live-brightgreen)](https://acm.vlaar.it/)
[![Conventional Commits](https://img.shields.io/badge/commits-conventional-blue)](https://www.conventionalcommits.org/)
[![License](https://img.shields.io/badge/license-MIT-green)](./LICENSE)

# PPL Calculations

PPL Calculations is an open-source project designed to simplify the calculation and management of weight and balance for an aeroclub fleet. By integrating a reliable Go-based HTTP backend with a responsive HTMX frontend, the application enables pilots to perform quick and accurate weight and balance calculations. The project aims to enhance operational safety and efficiency within aeroclubs and welcomes community contributions.

## Features

### ✅ Completed
- ✅ Save calculations as PDF on device with date and time
- ✅ Display calculations table
- ✅ Automatic take-off and landing distance calculations
- ✅ Seat position indicator to select proper arm momentum
- ✅ Fuel planning and endurance calculations based on trip and alternate distance
- ✅ Save form input values

### ⬜ Upcoming
- ⬜ Calculation information as title / subtitle when sharing via e.g. WhatsApp
- ⬜ Server-side error messages for form validation (currently not functioning properly on Firefox)
- ⬜ Support for Negative PA in performance calculations

## Installation

### Prerequisites
- [Docker](https://www.docker.com/get-started) installed on your machine
- Kubernetes cluster (optional, for deployment using Kustomize and Flux)
- **Dependencies:**
    - `xelatex`: For generating PDFs from LaTeX.
    - `librsvg2`: For rendering SVG images.
    - **Roboto Font**: Required for consistent typography in PDFs.
     
### Environment Variables

The application requires the following environment variables to function correctly. Below are example values:

```env
CSRF_KEY=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
PORT=8080
SESSION_KEY=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
TMP_PATH=/tmp/
```

### Variable Descriptions:
- **CSRF_KEY:** Used to encrypt and sign cookies and other secure parts of the application.
- **PORT:** The port on which the HTTP server runs.
- **SESSION_KEY:** Used for encrypting and signing session data.
- **TMP_PATH:** Temporary directory path used to store files temporarily. This directory does not need to be persistent.

## Running with Docker
Clone the Repository:

```bash
git clone https://github.com/michaelvlaar/ppl-calculations.git
cd ppl-calculations
```
Build the Docker Image:

```bash
docker build -t ppl-calculations:latest .
```
Run the Container:

```bash
docker run -d \
  -p 8080:8080 \
  -e CSRF_KEY=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX\
  -e PORT=8080 \
  -e SESSION_KEY=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX \
  -e TMP_PATH=/tmp/ \
  ppl-calculations:latest
```

## Usage
Once the application is running, navigate to http://localhost:8080 (or the appropriate port) to access the frontend interface. From there, you can:

- Input aircraft details and load configurations.
- Perform weight and balance calculations.
- Generate and download PDFs of your calculations.
- Save and manage calculation records for future reference.

## Changelog
For a detailed list of changes, see the [CHANGELOG](CHANGELOG.md).

## Contributing

PPL Calculations is open to contributions from the community. To contribute:

1. Fork the Repository
2. Create a Feature Branch
```bash
git checkout -b feature/your-feature-name
```
3. Commit Your Changes Using Conventional Commits
```bash
git commit -m "feat: add new feature"
```
4. Push to Your Fork
```bash
git push origin feature/your-feature-name
```
5. Open a Pull Request

Please ensure your commits follow the [Conventional Commits](https://www.conventionalcommits.org/) standard.

## Technologies Used

- **Backend:** Go (HTTP server)
- **Frontend:** HTMX
- **Changelog Generation:** git-chglog
- **Deployment:** Kubernetes with Kustomize
- **Design Patterns:** CQRS and Domain-Driven Design (DDD)

## License
This project is licensed under the MIT License. See the [LICENSE](LISENCE) file for details.
