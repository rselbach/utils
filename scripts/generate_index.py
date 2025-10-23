#!/usr/bin/env python3
"""Generate index.html from utility metadata files."""

import json
import os
import sys
from datetime import datetime
from pathlib import Path
from typing import List, Dict


def find_utilities(base_path: Path) -> List[Dict]:
    """
    find all directories containing util.json and extract metadata.

    Args:
        base_path: root directory to search from

    Returns:
        list of utility metadata dictionaries
    """
    utilities = []

    # scan all subdirectories for util.json
    for item in base_path.iterdir():
        if not item.is_dir():
            continue

        # skip hidden directories and scripts/workflow directories
        if item.name.startswith('.') or item.name in ['scripts', 'node_modules']:
            continue

        util_json = item / 'util.json'
        if not util_json.exists():
            continue

        try:
            with open(util_json, 'r') as f:
                metadata = json.load(f)

            # get last modified time of the directory
            last_modified = get_last_modified(item)

            # add computed fields
            metadata['path'] = item.name
            metadata['last_modified'] = last_modified

            utilities.append(metadata)
        except (json.JSONDecodeError, KeyError, IOError) as e:
            print(f"Warning: Failed to process {util_json}: {e}", file=sys.stderr)
            continue

    return utilities


def get_last_modified(directory: Path) -> str:
    """
    get the most recent modification time of any file in the directory.

    Args:
        directory: path to check

    Returns:
        formatted date string
    """
    latest_time = 0

    for root, _, files in os.walk(directory):
        for file in files:
            file_path = Path(root) / file
            try:
                mtime = file_path.stat().st_mtime
                latest_time = max(latest_time, mtime)
            except OSError:
                continue

    if latest_time == 0:
        latest_time = directory.stat().st_mtime

    return datetime.fromtimestamp(latest_time).strftime('%Y-%m-%d')


def generate_utility_html(utilities: List[Dict]) -> str:
    """
    generate HTML for all utilities.

    Args:
        utilities: list of utility metadata

    Returns:
        HTML string
    """
    if not utilities:
        return '<div class="no-utilities">No utilities available yet.</div>'

    # sort by name
    utilities.sort(key=lambda u: u.get('name', '').lower())

    html_parts = []
    for util in utilities:
        name = util.get('name', 'Unnamed Utility')
        description = util.get('description', 'No description available.')
        path = util.get('path', '')
        last_modified = util.get('last_modified', 'Unknown')

        html = f'''            <div class="utility-card">
                <h2><a href="/{path}/">{name}</a></h2>
                <p class="utility-description">{description}</p>
                <div class="utility-meta">Last updated: {last_modified}</div>
            </div>'''
        html_parts.append(html)

    return '\n'.join(html_parts)


def generate_index(base_path: Path, template_path: Path, output_path: Path) -> None:
    """
    generate index.html from template and utility metadata.

    Args:
        base_path: root directory containing utilities
        template_path: path to HTML template
        output_path: where to write index.html
    """
    # find all utilities
    utilities = find_utilities(base_path)
    print(f"Found {len(utilities)} utilities")

    # generate HTML for utilities
    utilities_html = generate_utility_html(utilities)

    # read template
    with open(template_path, 'r') as f:
        template = f.read()

    # replace placeholder
    output = template.replace('<!-- UTILITIES_PLACEHOLDER -->', utilities_html)

    # write output
    with open(output_path, 'w') as f:
        f.write(output)

    print(f"Generated {output_path}")


def main():
    """main entry point."""
    # get base directory (repo root)
    base_path = Path(__file__).parent.parent
    template_path = base_path / 'index-template.html'
    output_path = base_path / 'index.html'

    if not template_path.exists():
        print(f"Error: Template not found at {template_path}", file=sys.stderr)
        sys.exit(1)

    generate_index(base_path, template_path, output_path)


if __name__ == '__main__':
    main()
