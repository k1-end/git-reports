# Git Reports 🚀

**Visualize Git Repository Activity Like Never Before!**


Git Reports is a command-line tool written in Go that helps you analyze and visualize key metrics from your Git repositories. Whether you're a developer, team lead, or open-source maintainer, this tool provides actionable insights into commit patterns, developer activity, and more.

---

### Sample Output
![Heatmap](https://i.imgur.com/UsvrfvA.png)

![Commits Per Developer](https://imgur.com/HKATglQ.png)

![Commits Per Hour Of Day](https://imgur.com/26rQknn.png)

![File Types & Merge Commits per year](https://imgur.com/uY0kJQ8.png)

---

## ✨ Features

- **Heatmap of Commits**: Visualize commit activity over time.
- **Commits Per Developer**: See how much each contributor has contributed.
- **Commits Per Hour**: Analyze productivity patterns throughout the day.
- **Merge Commits Per Year**: Track merge activity trends over the years.
- **File Type Analysis**: Understand which file types are most frequently changed.
- **Date Range Filtering**: Analyze commits within a specific date range.
- **HTML Output**: Generate reports in HTML format for easy sharing.

---

## 🛠️ Installation

### Option 1: Download Pre-built Binaries
1. Visit the [releases page](https://github.com/k1-end/git-reports/releases)
2. Download the appropriate binary for your operating system (Windows, macOS, or Linux)
3. Make the binary executable (on Unix-like systems):
   ```bash
   chmod +x git-reports
   ```
4. Run the tool:
   ```bash
   ./git-reports
   ```

### Option 2: Build from Source
1. Make sure you have Go installed (version 1.16 or higher).
2. Clone this repository:
   ```bash
   git clone git@github.com:k1-end/git-reports.git
   ```
3. Build the project:
   ```bash
   cd git-reports
   go build -o git-reports
   ```
4. Run the tool:
   ```bash
   ./git-reports
   ```

---

## 🚀 Usage

### Basic Usage
To generate reports for a Git repository, simply run:
```bash
./git-reports
```
By default, it will analyze the current directory.

### Specify Repository Path
To analyze a specific Git repository, use the `--path` flag:
```bash
./git-reports --path /path/to/your/git/repo
```

### Choose Output Format
You can choose between console (default) and HTML output using the `--printer` flag:
```bash
./git-reports --printer html
```

### Save Output to File
To save the report to a file, use the `--output` flag:
```bash
./git-reports --output report.html
```

### Filter by Developer
To analyze commits by a specific developer, use the `--dev` flag:
```bash
./git-reports --dev developer@example.com
```

### Filter by Date Range
To analyze commits within a specific date range, use the `--from` and `--to` flags (format: `YYYY-MM-DD`):
```bash
./git-reports --from 2023-01-01 --to 2023-12-31
```

### Combine Options
You can combine multiple options:
```bash
./git-reports --path /path/to/repo --dev developer@example.com --from 2023-01-01 --to 2023-12-31 --printer html --output report.html
```

### Check Version
To check the version of Git Reports:
```bash
./git-reports --version
```

### Use console printer with a pager (less)
In case you want to use the console printer with a pager, you can use the `less` command with -R option to allow ANSI colors:
```bash
./git-reports --printer console | less -R
```

---

## 🧑‍💻 Why Use Git Reports?

- **Insightful Metrics**: Gain a deeper understanding of your repository's activity.
- **Easy to Use**: Simple CLI interface with no dependencies other than Go.
- **Customizable**: Filter reports by developer, time, or file type.
- **Date Range Filtering**: Analyze commits within specific time periods.
- **Open Source**: Free to use, modify, and contribute to.

---

## 🤝 Contributing

We welcome contributions! If you'd like to improve Git Reports, please follow these steps:
1. Fork the repository.
2. Create a new branch (`git checkout -b feature/your-feature`).
3. Commit your changes (`git commit -m 'Add some feature'`).
4. Push to the branch (`git push origin feature/your-feature`).
5. Open a pull request.

---

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## ⭐ Star This Project

If you find Git Reports useful, please consider giving it a star on GitHub! Your support helps us grow and improve the tool.

[![Star on GitHub](https://img.shields.io/github/stars/k1-end/git-reports?style=social)](https://github.com/k1-end/git-reports)

---

Made with ❤️ by [K1-end]. Happy coding! 🎉
