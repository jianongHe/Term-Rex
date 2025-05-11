#!/bin/bash

# Create bin directory
mkdir -p bin

# Compile for current platform as a test version
echo "Compiling test version for current platform..."
go build -o bin/term-rex main.go

echo "Compilation complete!"

# Note: Cross-platform compilation may require resolving dependencies
echo "Note: For cross-platform compilation, you may need to resolve audio library dependencies."
echo "It's recommended to use GitHub Actions or Docker for cross-platform builds."
