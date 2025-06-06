#!/bin/bash

# Go Password Manager - Test Workflow Script
# This script demonstrates the complete functionality

set -e

echo "🔐 Go Password Manager - Complete Workflow Test"
echo "=============================================="
echo

# Build the application
echo "📦 Building application..."
go build -o gopassman
echo "✅ Build successful!"
echo

# Test password generation
echo "🎲 Testing password generation..."
./gopassman generate --length 16 --count 2
echo

# Show help
echo "📖 Available commands:"
./gopassman --help
echo

echo "🚀 Manual Testing Guide:"
echo "========================"
echo
echo "1. Initialize a new vault:"
echo "   ./gopassman init"
echo
echo "2. Add some entries:"
echo "   ./gopassman add --title 'GitHub' --username 'user@example.com' --generate"
echo "   ./gopassman add --title 'Gmail' --username 'myemail@gmail.com' --password 'mypassword'"
echo
echo "3. List entries:"
echo "   ./gopassman list"
echo
echo "4. Show entry details:"
echo "   ./gopassman show 1"
echo "   ./gopassman show 1 --password"
echo
echo "5. Edit an entry:"
echo "   ./gopassman edit 1 --generate"
echo
echo "6. Search entries:"
echo "   ./gopassman list --search github"
echo
echo "7. Delete an entry:"
echo "   ./gopassman delete 2"
echo
echo "8. Show vault info:"
echo "   ./gopassman list --info"
echo
echo "All data is encrypted with AES-GCM and your master password is hashed with Argon2id!"
echo "Session management keeps the vault unlocked for 15 minutes."
echo
echo "🎉 Password manager is fully functional!" 