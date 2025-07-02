package main

import (
    "flag"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "time"
)

func main() {
    repoDir := flag.String("repo", ".", "path to your Git repository")
    startStr := flag.String("start", "", "start date (inclusive), format YYYY-MM-DD")
    endStr := flag.String("end", "", "end date (inclusive), format YYYY-MM-DD")
    msg := flag.String("msg", "Automated backdate commit", "commit message prefix")
    flag.Parse()

    if *startStr == "" || *endStr == "" {
        fmt.Println("⚠️  You must supply both --start and --end dates")
        os.Exit(1)
    }

    start, err := time.Parse("2006-01-02", *startStr)
    if err != nil {
        fmt.Printf("Invalid start date: %v\n", err)
        os.Exit(1)
    }
    end, err := time.Parse("2006-01-02", *endStr)
    if err != nil {
        fmt.Printf("Invalid end date: %v\n", err)
        os.Exit(1)
    }

    if stat, err := os.Stat(*repoDir); err != nil || !stat.IsDir() {
        fmt.Printf("❌ Repo not found or not a directory: %s\n", *repoDir)
        os.Exit(1)
    }

    step := 1
    if start.After(end) {
        step = -1
    }

    filePath := filepath.Join(*repoDir, "SOMETHING.md")

    // ensure file exists and is tracked by git before loop
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        f, err := os.Create(filePath)
        if err != nil {
            fmt.Printf("Could not create file: %v\n", err)
            os.Exit(1)
        }
        f.Close()

        cmd := exec.Command("git", "-C", *repoDir, "add", "SOMETHING.md")
        if out, err := cmd.CombinedOutput(); err != nil {
            fmt.Printf("git add error: %v\n%s\n", err, out)
            os.Exit(1)
        }

        cmd = exec.Command("git", "-C", *repoDir, "commit", "-m", "init SOMETHING.md")
        if out, err := cmd.CombinedOutput(); err != nil {
            fmt.Printf("git commit init error: %v\n%s\n", err, out)
            os.Exit(1)
        }
    }

    for d := start; ; d = d.AddDate(0, 0, step) {
        datestr := d.Format(time.RFC3339)

        env := os.Environ()
        env = append(env,
            "GIT_AUTHOR_DATE="+datestr,
            "GIT_COMMITTER_DATE="+datestr,
        )

        f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0o644)
        if err != nil {
            fmt.Printf("Error opening %s: %v\n", filePath, err)
            os.Exit(1)
        }
        fmt.Fprintln(f, datestr[:10])
        f.Close()

        commitMsg := fmt.Sprintf("%s: %s", *msg, datestr[:10])

        cmd := exec.Command("git", "-C", *repoDir,
            "commit", "-S", "--only", "SOMETHING.md", "-m", commitMsg,
        )
        cmd.Env = env
        if out, err := cmd.CombinedOutput(); err != nil {
            fmt.Printf("git commit error: %v\n%s\n", err, out)
            os.Exit(1)
        }

        fmt.Printf("✔ Committed %s\n", datestr[:10])

        if (step > 0 && !d.Before(end)) || (step < 0 && !d.After(end)) {
            break
        }
    }

    fmt.Println("✅ All done! Run `git log --pretty=fuller --date=iso` then `git push`.")
}
