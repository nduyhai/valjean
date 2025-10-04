# go-module

[![Go](https://img.shields.io/badge/go-1.24+-blue)](https://go.dev/)
[![License](https://img.shields.io/github/license/nduyhai/go-module)](LICENSE)

A GitHub template repository for bootstrapping a new Go project with a clean, idiomatic layout.

## Features

- ✅ Linter config (`golangci-lint`)
- ✅ GitHub actions
- ✅ Basic Makefile
- ✅ MIT License

## Getting Started

### 📦 Create a New Project

Click the **[Use this template](https://github.com/your-org/go-module/generate)** button to generate a new repository based on this template.

### 🛠️ Customize

After creating your repo, follow these steps:

```bash
# Clone your new project
git clone https://github.com/your-username/your-project-name
cd your-project-name

# Update module path
go mod edit -module github.com/your-username/your-project-name

# Tidy up dependencies
go mod tidy
```
Edit the README.md, package names, and other placeholders as needed.

### 🏃 Run the Project
```shell
make run
```

