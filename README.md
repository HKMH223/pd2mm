# pd2mm
A command-line Payday 2 mod manager.

This is tool will require trial and error to get the effect you want. If you do not feel comfortable using this, please use [Myth Mod Manager](https://github.com/Wolfmyths/Myth-Mod-Manager)

> [!Warning]
> #### Usage: [USAGE.md](./USAGE.md)
> #### Under Development
> ###### There is no stable version.

## Privacy
`pd2mm` is an open source project. Your commit credentials as author of a commit will be visible by anyone. Please make sure you understand this before submitting a PR.
Feel free to use a "fake" username and email on your commits by using the following commands:
```bash
git config --local user.name "USERNAME"
git config --local user.email "USERNAME@SOMETHING.com"
```

## Requirements (Building)
- Go 1.24.3 or later.
- GNU Make 4.4.1 or later.

### Requirements (Development)
- All of the above requirements.
- golangci-lint 
    - `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
- gofumpt
    - `go install mvdan.cc/gofumpt@latest`
- Optionally: deadcode
    - `go install golang.org/x/tools/cmd/deadcode@latest`

## Building
- Run `make all`.
- Run `pd2mm`. Use `pd2mm -h` to view a list of commands.

## Platforms

|        | Windows|Linux (Untested)|Mac OS (Untested)|
|--------|--------|----------------|-----------------|
| x86-64 | ✅ | ❌ | ❌ |
| x86    | ❌ | ❌ | ❌ |
| ARM64  | ❌ | ❌ | ❌ |

## Contribution Guidelines
If you would like to contribute to `pd2mm` please take the time to carefully read the guidelines below.

### Commit Workflow
- Read through [FORMATTING.md](./FORMATTING.md).
- Run `make lint` and ensure ALL diagnostics are fixed.
- Run `make fmt` to ensure consistent formatting.
- Create concise, descriptive commit messages to summarize your changes.
    - Optionally: use `git cz` with the [Commitizen CLI](https://github.com/commitizen/cz-cli#conventional-commit-messages-as-a-global-utility) to prepare commit messages.
- Provide *at least* one short sentence or paragraph in your commit message body to describe your thought process for the changes being committed.

### Pull Requests (PRs) should only contain one feature or fix.
It is very difficult to review pull requests which touch multiple unrelated features and parts of the codebase.

Please do not submit pull requests like this; you will be asked to separate them into smaller PRs that deal only with one feature or bug fix at a time.

### Codebase refactors must have prior approval.
Refactors to the structure of the codebase are not taken lightly and require prior discussion and approval.

Please do not start refactoring the codebase with the expectation of having your changes integrated until you receive an explicit approval or a request to do so.

Similarly, when implementing features and bug fixes, please stick to the structure of the codebase as much as possible and do not take this as an opportunity to do some "refactoring along the way".

It is extremely difficult to review PRs for features and bug fixes if they are lost in sweeping changes to the structure of the codebase.

# License
See [LICENSE](./LICENSE) file.