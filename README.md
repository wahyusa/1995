# 1995

1995 Code

Average boring command executor cli tools build with Go

## How it works

- Pick a start and end date (`--start`, `--end`).
- For every day in that range, it appends the date to `SOMETHING.md`.
- Commits **only** that file each time, signed with `-S` (so it shows as "Verified").
- Forces the commit timestamp with `GIT_AUTHOR_DATE` and `GIT_COMMITTER_DATE`.

## Usage

```bash
go run . --repo=/path/to/your/repo --start=2025-01-01 --end=2025-06-30 --msg="chore: blabla"
```

Execute in your current working directory/CWD

```bash
go run . --start=2025-01-01 --end=2025-06-30 --msg="chore: blabla"
```

It will print each commit as it goes. Once done, you can verify with:

```bash
cd /path/to/your/repo
git log --pretty=fuller --date=iso
git log --show-signature
```

Then push:

```bash
git push origin main
```

## Requirements

- Go ≥ 1.18
- GPG key set up and configured to sign commits (`git config --global commit.gpgsign true`)

## Why `--only`?

By default, Git would also commit any other dirty files. `--only SOMETHING.md` makes sure each commit touches just that file, so your working tree stays untouched.

## Have fun

Use responsibly. Great for showing off CLI + Git knowledge, and learning how commit metadata works under the hood.
