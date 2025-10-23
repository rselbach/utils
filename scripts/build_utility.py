#!/usr/bin/env python3
"""Build a single utility based on its metadata."""

import json
import subprocess
import sys
from pathlib import Path


def build_utility(utility_dir: Path) -> bool:
    """
    build a utility if it has a build command.

    Args:
        utility_dir: path to utility directory

    Returns:
        True if build succeeded or no build needed, False if build failed
    """
    util_json = utility_dir / 'util.json'

    if not util_json.exists():
        print(f"Error: {util_json} not found", file=sys.stderr)
        return False

    try:
        with open(util_json, 'r') as f:
            metadata = json.load(f)
    except (json.JSONDecodeError, IOError) as e:
        print(f"Error: Failed to read {util_json}: {e}", file=sys.stderr)
        return False

    build_cmd = metadata.get('build')

    if not build_cmd:
        print(f"No build command for {utility_dir.name}, skipping")
        return True

    print(f"Building {utility_dir.name}...")
    print(f"Running: {build_cmd}")

    try:
        result = subprocess.run(
            build_cmd,
            shell=True,
            cwd=utility_dir,
            check=True,
            capture_output=False,
            text=True
        )
        print(f"Build succeeded for {utility_dir.name}")
        return True
    except subprocess.CalledProcessError as e:
        print(f"Build failed for {utility_dir.name}", file=sys.stderr)
        return False


def main():
    """main entry point."""
    if len(sys.argv) != 2:
        print("Usage: build_utility.py <utility_directory>", file=sys.stderr)
        sys.exit(1)

    utility_dir = Path(sys.argv[1])

    if not utility_dir.is_dir():
        print(f"Error: {utility_dir} is not a directory", file=sys.stderr)
        sys.exit(1)

    success = build_utility(utility_dir)
    sys.exit(0 if success else 1)


if __name__ == '__main__':
    main()
