# 🔮 GitHub Horoscope CLI

A fun and cosmic command-line tool that gives you a **daily horoscope** based on your GitHub activity — starring your open issues, PRs, and a dash of zodiac wisdom ✨

---

## 🌟 Features

- 🔍 Automatically fetches open issues and PRs from a repo
- 🌠 Assigns you a daily zodiac sign
- 🧘 Delivers a randomly selected GitHub-themed fortune
- 📋 Outputs everything in a clean, tabular format

---

## 📦 Installation

### ⚡ Option 1: Run from source (requires Go)

```bash
git clone https://github.com/AbhishekLaddha54/github-horoscope
cd github-horoscope
go mod tidy
go run main.go --repo owner/repo-name
