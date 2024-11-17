#!/bin/bash

main() {
    add_prefix
    add_msg
}

add_prefix() {
    # Prompt for commit type
    while true; do
        echo "Enter commit type (prefix):"
        echo "1. feat"
        echo "2. fix"
        echo "3. docs"
        echo "4. chores"
        echo "Choose a number from 1-4: "
        read -r num
        case $num in
            1) type="feat"; break;;
            2) type="fix"; break;;
            3) type="docs"; break;;
            4) type="chores"; break;;
            *) echo "Please choose again";;
        esac
    done
}

add_msg() {
    while true; do
        # Prompt for commit message
            echo "Enter commit message (excluding prefix): "
            read -r message

            # Check if message is empty
            if [ -z "$message" ]; then
                echo "Commit message cannot be empty."
                continue
            fi

            # Commit the changes
            if git commit -S -m "$type: $message"; then
                echo "Commit successful."
                break
            else
                echo "Error committing changes. Please check the details and try again."
                return 1
            fi
    done
}

main
