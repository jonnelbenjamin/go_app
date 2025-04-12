
# 🌟 Go CLI Greeter Application

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-blue?logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-green" alt="License">
</p>

A feature-packed CLI application with colorful output, built with:
- 🏗️ Cobra CLI framework
- 🎨 Fatih/color for terminal styling
- 🌤️ Free weather API integration
- 😄 Random joke generator

## 🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/jonnelbenjamin/go_app.git
cd go_app

# Initialize module
go mod init github.com/jonnelbenjamin/go_app

# Add dependencies
go get github.com/spf13/cobra@latest
go get github.com/fatih/color@latest

# Install dependencies
go mod tidy

# Run directly (development)
go run . --name Alice --joke

# Or build and run (production)
go build -o greeter
./greeter --weather --city Tokyo

go run . --name Bob --uppercase

# Default city (London)
go run . --weather

# Specific city
go run . --weather --city Tokyo

# Random joke
go run . --joke

# Random fact
go run . --fact

# Number guessing game
go run . --game

# Personal greeting with weather and joke
go run . --name Charlie --weather --city Paris --joke

# Full experience
go run . --name Dana --uppercase --weather --city NewYork --joke --fact --game

go run . --help