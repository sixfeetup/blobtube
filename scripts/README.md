# Scripts

This directory contains utility scripts for the BlobTube project.

## create-github-issues.sh

Creates all GitHub infrastructure for the project:
- 7 milestones (one per phase)
- 15 labels (phases, epics, priority, status)
- 28 issues with full details, dependencies, and estimates

### Prerequisites

1. Install GitHub CLI (`gh`):
   ```bash
   # Already available in nix-shell
   nix-shell -p gh
   ```

2. Authenticate with GitHub:
   ```bash
   nix-shell -p gh --run 'gh auth login'
   ```
   Follow the prompts to authenticate with your GitHub account.

### Usage

```bash
# Run from repository root
nix-shell -p gh --run './scripts/create-github-issues.sh'
```

The script will:
1. Create 7 milestones for each development phase
2. Create labels for organizing issues (phases, epics, priorities)
3. Create all 28 tickets with:
   - Full descriptions
   - Task checklists
   - Dependencies
   - Time estimates
   - Links to relevant documentation (ADRs, user stories)
   - Appropriate labels and milestones

### After Running

View the results:
- All issues: https://github.com/sixfeetup/blobtube/issues
- Milestones: https://github.com/sixfeetup/blobtube/milestones
- Project board: Create manually or use GitHub's auto-triage features

### Troubleshooting

**"Not authenticated" error:**
```bash
nix-shell -p gh --run 'gh auth login'
```

**Permission denied:**
```bash
chmod +x scripts/create-github-issues.sh
```

**Script fails partway through:**
The script is idempotent for labels (won't fail if labels exist). For milestones and issues, you may need to clean up partial creates manually before re-running.
