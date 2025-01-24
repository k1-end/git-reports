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

---

## 🛠️ Installation

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
   ./git-reports /path/to/your/git/repo
   ```

---

## 🚀 Usage

### Basic Usage
To generate reports for a Git repository, simply run:
```bash
./git-reports /path/to/your/git/repo
```

### Filter by Developer
To analyze commits by a specific developer, use the `--dev` flag:
```bash
./git-reports /path/to/your/git/repo --dev developer@example.com
```

### Filter by Date Range
To analyze commits within a specific date range, use the `--from` and `--to` flags (format: `YYYY-MM-DD`):
```bash
./git-reports /path/to/your/git/repo --from 2023-01-01 --to 2023-12-31
```

### Combine Filters
You can combine filters to analyze commits by a specific developer within a date range:
```bash
./git-reports /path/to/your/git/repo --dev developer@example.com --from 2023-01-01 --to 2023-12-31
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

Made with ❤️ by [Your Name]. Happy coding! 🎉
