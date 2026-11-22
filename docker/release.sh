#!/usr/bin/env bash
set -e

echo "-> Checking Git repository status..."

# 1. Check for uncommitted changes
if [ -n "$(git status --porcelain)" ]; then
  echo ""
  echo "❌ ERROR: You have uncommitted changes in the working directory!"
  echo "Commit or stash (git stash) them before creating a tag."
  exit 1
fi

# 2. Check for unpushed commits on the current branch
CURRENT_BRANCH=$(git branch --show-current)
if [ -z "$CURRENT_BRANCH" ]; then
  echo "❌ ERROR: You appear to be in a 'detached HEAD' state. Switch to the appropriate branch."
  exit 1
fi

# Make sure we have up-to-date remote information
git fetch origin "$CURRENT_BRANCH" &>/dev/null || true

# Check if the local branch is ahead of remote
UNPUSHED_COMMITS=$(git log "origin/$CURRENT_BRANCH..$CURRENT_BRANCH" --oneline)
if [ -n "$UNPUSHED_COMMITS" ]; then
  echo ""
  echo "❌ ERROR: You have unpushed commits on branch '$CURRENT_BRANCH':"
  echo "$UNPUSHED_COMMITS"
  echo ""
  echo "Run 'git push origin $CURRENT_BRANCH' before running the release script!"
  exit 1
fi

echo "✅ Git state is clean and synchronized with remote."
echo ""

# 3. Get the latest tag matching the vX.Y.Z format
LATEST_TAG=$(git tag -l "v[0-9]*.[0-9]*.[0-9]*" --sort=-v:refname | head -n 1)

if [ -z "$LATEST_TAG" ]; then
  echo "No tags matching the vX.Y.Z format found in this repository."
  echo "Setting base tag to v0.0.0"
  LATEST_TAG="v0.0.0"
fi

echo "Latest tag found: $LATEST_TAG"

# 4. Extract X, Y, Z using regex
if [[ "$LATEST_TAG" =~ ^v([0-9]+)\.([0-9]+)\.([0-9]+)$ ]]; then
  MAJOR="${BASH_REMATCH[1]}"
  MINOR="${BASH_REMATCH[2]}"
  PATCH="${BASH_REMATCH[3]}"
else
  echo "Error: Latest tag ($LATEST_TAG) does not match the vX.Y.Z pattern"
  exit 1
fi

# 5. Calculate next tag proposals
NEXT_PATCH="v${MAJOR}.${MINOR}.$((PATCH + 1))"
NEXT_MINOR="v${MAJOR}.$((MINOR + 1)).0"
NEXT_MAJOR="v$((MAJOR + 1)).0.0"

echo ""
echo "Select the version to release:"
echo "  1) Patch (bug fixes):          $NEXT_PATCH"
echo "  2) Minor (new features):       $NEXT_MINOR"
echo "  3) Major (breaking changes):   $NEXT_MAJOR"
echo ""
read -rp "Enter a number (1, 2 or 3): " CHOICE

case "$CHOICE" in
  1)
    SELECTED_TAG="$NEXT_PATCH"
    ;;
  2)
    SELECTED_TAG="$NEXT_MINOR"
    ;;
  3)
    SELECTED_TAG="$NEXT_MAJOR"
    ;;
  *)
    echo "Invalid choice. Cancelled."
    exit 1
    ;;
esac

echo ""
echo "Selected new tag: $SELECTED_TAG"
read -rp "Are you sure you want to create and push this tag? (y/N): " CONFIRM

if [[ "$CONFIRM" =~ ^[Yy]$ ]]; then
  # Create local tag
  git tag "$SELECTED_TAG"

  # Push tag to remote repository (origin)
  git push origin "$SELECTED_TAG"

  echo "Success! Tag $SELECTED_TAG has been created and pushed to origin."
else
  echo "Cancelled by user."
fi
