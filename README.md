# git-lint
A linter for enforcing git repository rules and conventions.

## Installation

Checkout [releases](https://github.com/troykinsella/git-lint/releases) and download the appropriate binary for your system.
Put the binary in a convenient place, such as `/usr/local/bin/git-lint`.

Or, run the handy dandy install script:
(Note: go read the script and understand what you're running before trusting it)
```bash
export PREFIX=~ # install into ~/bin
wget -q -O - https://raw.githubusercontent.com/troykinsella/git-lint/master/install.sh | bash
```

Or, run these commands to download and install:
```bash
VERSION=0.0.1
OS=darwin # or linux
curl -SL -o /usr/local/bin/git-lint https://github.com/troykinsella/git-lint/releases/download/v${VERSION}/git-lint_${OS}_amd64
chmod +x /usr/local/bin/git-lint
```

Or, for [Go lang](https://golang.org/doc/code.html) projects, from your `GOPATH`:
```bash
go get github.com/troykinsella/git-lint
```

Lastly, test the installation:
```bash
git-lint -h
```

## Usage

### Requirements

As `git-lint` uses [src-d/go-git](https://github.com/src-d/go-git),
it does not rely on `git` being installed.

### Configuration

As the kinds of rules that this linter enforces, unlike most other linters,
are not based on best practices or basic prevention of common errors, but are
based on organization-/project-specific rules, NO checks are performed by default.
`git-lint` must be told what you want it to enforce by activating checks in a 
configuration file.  

The configuration file is in YAML, and defaults to the path `.git-lint.yml`,
although an alternate path can be specified by passing `-c` (or `--config`) to
`git-lint`.

### Configuration File Format

```yaml
# .git-lint.yml
rules:
  <rule-name>:
    <rule-options>
```

### Full Configuration Example

```yaml
# .git-lint.yml
rules:
  branch_count:
    max: 100
    warn: true
  branch_last_commit:
    max_duration: 3M
  branch_name:
    allow: true
    patterns:
    - master
    - develop
    - feature/.*
    - release/v.*
  branch_singleton:
    singletons:
    - release/v.*
  git_ignore:
    entries:
    - tmp
  git_keep:
    directories:
    - bin
  tag_name:
    allow: true
    patterns:
    - v.*
```

### Common Rule Options

* `warn`: Optional. Boolean. Default: `false`. When set to `true` 
  and when the rule check fails, the overall `git-lint` execution
  will not be failed due to this check. Instead, only warning 
  messages will be emitted.

### Time Duration Strings

`git-lint` uses [sloppy_duration](https://github.com/troykinsella/sloppy-duration)
for parsing strings that represent time durations.
It's "sloppy" because it sacrifices precision for conciseness. For example,
the amount of time in a month depends on the month, but `sloppy_duration` 
takes a close-enough guess by using "one year / 12".

Here's a summary of supported syntax:

| Example value | Meaning |
| ------------- | ------- |
| 2s            | 2 seconds |
| 5m            | 5 minutes |
| 10h           | 10 hours |
| 3d            | 3 days (3 x 24 hours) |
| 2w            | 2 weeks (14 x 1 day) |
| 6M            | 6 months (1 year / 2) |
| 1y            | 1 year (365 days) |

### Running `git-lint`

`git-lint` operates on a locally cloned repository directory.
Run it like so:

```bash
git clone git@github.com/some/repo.git
cd repo
git-lint
```

### Running `git-lint` with Docker

```bash
docker pull troykinsella/git-lint

git clone git@github.com/some/repo.git
cd repo

docker run -it --rm -v $PWD:/repo -w /repo troykinsella/git-lint git-lint
```

### Rules

#### `GL001` - `branch_count`

Ensures a maximum number of permitted branches in a repository.

Options:
* `max`: Required. Integer. The maximum number of branches to allow.

Example:
```yaml
rules:
  branch_count:
    max: 100
```

#### `GL002` - `branch_last_commit`

Fails when the last commit for a branch is older than a specified time duration.

Options:
* `max_duration`: Required. String. A time duration string that 
  defines the maximum time span that can elapse since the last commit. 

Example:
```yaml
# Fail when a branch hasn't been updated in the last 6 months
rules:
  branch_last_commit:
    max_duration: 6M
```

#### `GL003` - `branch_name`

Specifies a whitelist or a blacklist of acceptable branch name patterns.

Options:
* `allow`: Optional. Boolean. Default: `false`. When `true`, allow only the
  specified branch name patterns to exist, and when `false` disallow
  the specified branch name patterns.
* `patterns`: Required. List of regular expressions that match branch names.

Example:
```yaml
# Fail a branch called "feechure-BUG-32", for example
rules:
  branch_name:
    allow: true
    patterns:
    - master
    - develop
    - feature/.*
    - release/v.*
```

#### `GL004` - `branch_singleton`

Ensure there is only one of a given branch name for each of the specified patterns.

Options:
* `singletons`: Required. List of regular expressions that match branch names.
  When more than one match is found, the rule fails.

Example:
```yaml
# Ensure there is only one release/* branch at a time.
rules:
  branch_singleton:
    singletons:
    - release/v.*
```

#### `GL005` - `git_ignore`

Ensure `.gitignore` entries exist.

Options:
* `entries`: Required. A list of strings, which are required to appear
  in `.gitignore`.

Example:
```yaml
rules:
  git_ignore:
    entries:
    - tmp
```

#### `GL006` - `git_keep`

Ensure a directory is kept by git, whether by containing a `.gitkeep`-like
file, or by simply being populated with files.

Options:
* `directories`: Required. A list of directories that should exist.

Example:
```yaml
rules:
  git_keep:
    directories:
    - bin
```

#### `GL007` - `tag_name`

Specifies a whitelist or a blacklist of acceptable tag name patterns.

Options:
* `allow`: Optional. Boolean. Default: `false`. When `true`, allow only the
  specified tag name patterns to exist, and when `false` disallow
  the specified tag name patterns.
* `patterns`: Required. List of regular expressions that match tag names.

Example:
```yaml
# Fail a tag called "1.2.3", for example, demanding a "v" prefix
rules:
  tag_name:
    allow: true
    patterns:
    - v.*
```

## License

MIT © Troy Kinsella
