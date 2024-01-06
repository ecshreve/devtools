#!/bin/bash

# This script automates the generation of a Git commit message based on user input.
# It uses external tools such as jq, gum, and git. 
# The script checks for necessary prerequisites, generates a commit message, 
# allows the user to edit the message, and commits the changes with the final message.

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
    echo "Please look at the changes to these files and provide the following:
    - DESC: concise description containing 40 characters or less
    - SUMMARY: a brief summary suitable for the body of a commit message explaining the changes. 
    - TYPE: Classify this change as one of the [fix, feat, test, docs]
    - SCOPE: Provide one word describing the area of the repository most affected

    Format: { \"DESC\": \"<DESC>\", \"SUMMARY\": \"<SUMMARY>\", \"TYPE\": \"<TYPE>\", \"SCOPE\": \"<SCOPE>\" }"
}

# Generate commit message based on user input
generate_commit_message() {
    echo $(git diff --cached | mods --fanciness 4 --temp .4 "$(define_prompt)" 'format as json')
}

# Parse JSON response into components of the commit message
parse_response() {
    local raw_resp=$1
    DESC=$(echo $raw_resp | jq -r '.DESC')
    SUMMARY=$(echo $raw_resp | jq -r '.SUMMARY')
    TYPE=$(echo $raw_resp | jq -r '.TYPE')
    SCOPE=$(echo $raw_resp | jq -r '.SCOPE')
}

# Edit the commit message with user input
edit_commit_message() {
    continue_editing=true
    while $continue_editing; do
        clear
        FIRST=$(echo "$TYPE($SCOPE): $DESC" | fold -w 50 -s)
        GENCOM=$(echo -e "$FIRST\n\n$SUMMARY" | fold -w 72 -s)
        gum style --border normal --margin "1" --padding "1 2" --border-foreground 212 --width 80 "$GENCOM"
        gum confirm "Use this?" && continue_editing=false && break

        CANCEL=$(echo '{{ Bold "=Cancel=" }}' | gum format -t template)

        PICK=$(gum choose --limit 1 "TYPE" "SCOPE" "DESC" "SUMMARY" $CANCEL)
        case $PICK in
            "TYPE") 
                TYPE=$(gum choose --limit 1 "fix" "feat" "test" "docs")
                ;;
            "SCOPE")
                SCOPE=$(gum input --value "$SCOPE")
                ;;
            "DESC")
                DESC=$(gum input --value "$DESC")
                ;;
            "SUMMARY")
                SUMMARY=$(gum write --width 72 --value "$SUMMARY")
                ;;
            "Cancel")
                return 1
                ;;
        esac
    done
}

# Main execution flow
main() {
    check_required_commands
    check_staged_files
    check_openai_api_key

    raw=$(generate_commit_message)
    parse_response "$raw"

    edit_commit_message
    if [ $? -eq 1 ]; then
        echo "Commit message editing canceled by the user."
        exit 1
    fi

    gum style --border double --margin "1" --padding "1 2" --border-foreground 35 --width 80 "$GENCOM"
    if gum confirm "Commit changes?"; then
        git commit -m "$GENCOM"
    else
        echo "Commit canceled by the user."
        exit 1
    fi

    exit 0
}

main 