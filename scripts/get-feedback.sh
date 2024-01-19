#!/bin/bash

# Check if necessary tools are installed (jq, gum, git)
check_required_commands() {
    for cmd in jq gum git; do
        if ! command -v $cmd &> /dev/null; then
            echo "Error: $cmd is not installed." >&2
            exit 1
        fi
    done
}

# Check if there are files staged for commit
check_staged_files() {
    if [ -z "$(git diff --cached --name-only)" ]; then
        echo "No files staged for commit."
        exit 1
    fi
}

# Check for OPENAI_API_KEY environment variable.
check_openai_api_key() {
    if [ -z "$OPENAI_API_KEY" ]; then
        echo "Error: OPENAI_API_KEY environment variable is not set." >&2
        exit 1
    fi
}

# Define a user prompt for input.
define_prompt() {
    echo "Please look at the changes to these files and provide general feedback on the changes, such as:
    - Are there any typos?
    - Are there any bugs?
    - Are there any improvements that can be made?
    - Are there any other comments you would like to make?
    "
}

# Get feedback on staged changes
get_feedback() {
    git diff --cached | mods --fanciness 4 --temp .4 "$(define_prompt)" | glow -
}

main() {
    check_required_commands
    check_staged_files
    check_openai_api_key

    get_feedback
}

main
