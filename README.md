# trackpoint

A lightweight CLI for tracing and diffing infrastructure state changes across deploys.

---

## Installation

```bash
go install github.com/yourorg/trackpoint@latest
```

Or build from source:

```bash
git clone https://github.com/yourorg/trackpoint.git && cd trackpoint && go build -o trackpoint .
```

---

## Usage

Capture a snapshot of your current infrastructure state before a deploy:

```bash
trackpoint snapshot --env production --out before.json
```

Run your deploy, then capture the post-deploy state and diff the two snapshots:

```bash
trackpoint snapshot --env production --out after.json
trackpoint diff before.json after.json
```

Example output:

```
~ service/api        replicas: 2 → 4
+ service/worker     (new)
- service/legacy-job (removed)
```

### Commands

| Command              | Description                              |
|----------------------|------------------------------------------|
| `snapshot`           | Capture current infrastructure state     |
| `diff <a> <b>`       | Compare two state snapshots              |
| `history`            | View a log of past snapshots             |

---

## License

MIT © yourorg