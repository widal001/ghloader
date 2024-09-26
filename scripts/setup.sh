#! /bin/bash
# Runs setup commands, like saving the GitHub token to a local .env file

# Saves GitHub token to .env file
echo "Enter your GitHub token"
read token
echo "GITHUB_TOKEN=${token}" > .env
echo "Saved your GitHub token to a .env file in the current directory"
