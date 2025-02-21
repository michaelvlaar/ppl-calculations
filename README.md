[![Live Demo](https://img.shields.io/badge/demo-live-brightgreen)](https://acm.vlaar.it/)
[![Conventional Commits](https://img.shields.io/badge/commits-conventional-blue)](https://www.conventionalcommits.org/)
[![License](https://img.shields.io/badge/license-MIT-green)](./LICENSE)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fmichaelvlaar%2Fppl-calculations.svg?type=shield&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2Fmichaelvlaar%2Fppl-calculations?ref=badge_shield&issueType=license)

# PPL Calculations

PPL Calculations is an open-source project designed to simplify the calculation and management of weight and balance for
an aeroclub fleet. By integrating a reliable Go-based HTTP backend with a responsive HTMX frontend, the application
enables pilots to perform quick and accurate weight and balance calculations. The project aims to enhance operational
safety and efficiency within aeroclubs and welcomes community contributions.

## TODOS 

- ⬜ Update PDF latex engine to https://github.com/go-pdf/fpdf

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

## Development Environment

This repository includes a Docker container that simplifies development by automating key tasks. The environment automatically refreshes the Golang backend, processes Tailwind CSS, and regenerates templates using [air](https://github.com/cosmtrek/air).

### Getting Started

To launch the development environment, simply run:

```bash
docker compose up
```

Once started, any changes to your code or configuration will trigger an automatic refresh, ensuring a smooth and efficient development workflow.

## Changelog

For a detailed list of changes, see the [CHANGELOG](CHANGELOG.md).

## Minification

```bash
npx tailwindcss -i ./assets/css/style.css -o ./assets/css/style.min.css -m
```

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
- **CSS / Minification:** tailwindcss cli
- **Changelog Generation:** git-chglog
- **Deployment:** Kubernetes with Kustomize
- **Design Patterns:** CQRS and Domain-Driven Design (DDD)

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
