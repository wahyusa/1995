// backdate.go
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
    repoDir := flag.String("repo", "./repo", "path to your Git repository")
    startStr := flag.String("start", "", "start date (inclusive), format YYYY-MM-DD")
    endStr   := flag.String("end",   "", "end date (inclusive), format YYYY-MM-DD")
    msg      := flag.String("msg",   "Automated backdate commit", "commit message")
    flag.Parse()

    if *startStr == "" || *endStr == "" {
        fmt.Println("⚠️  You must supply both --start and --end dates")
        os.Exit(1)
    }

    // parse dates
    start, err := time.Parse("2006-01-02", *startStr)
    if err != nil {
        panic(err)
    }
    end, err := time.Parse("2006-01-02", *endStr)
    if err != nil {
        panic(err)
    }

    // ensure repo exists
    if _, err := os.Stat(*repoDir); os.IsNotExist(err) {
        fmt.Printf("❌ Repo not found: %s\n", *repoDir)
        os.Exit(1)
    }

    // iterate days
    for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
        datestr := d.Format(time.RFC3339)
        env := os.Environ()
        env = append(env, "GIT_AUTHOR_DATE="+datestr, "GIT_COMMITTER_DATE="+datestr)

        // make a trivial change: append the date to a file
        file := filepath.Join(*repoDir, "backdate.log")
        f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
        if err != nil {
            panic(err)
        }
        fmt.Fprintln(f, datestr)
        f.Close()

        // git add
        cmd := exec.Command("git", "-C", *repoDir, "add", "backdate.log")
        cmd.Env = env
        if out, err := cmd.CombinedOutput(); err != nil {
            fmt.Printf("git add error: %s\n%s\n", err, out)
            os.Exit(1)
        }

        // git commit -S
        commitMsg := fmt.Sprintf("%s: %s", *msg, datestr[:10])
        cmd = exec.Command("git", "-C", *repoDir,
            "commit", "-S", "-m", commitMsg,
        )
        cmd.Env = env
        if out, err := cmd.CombinedOutput(); err != nil {
            fmt.Printf("git commit error: %s\n%s\n", err, out)
            os.Exit(1)
        }

        fmt.Printf("✔ Committed %s\n", datestr[:10])
    }

    fmt.Println("✅ All done! You can now `git push` your back‑dated commits.")
}

