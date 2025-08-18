#!/bin/bash

# Test script for Kafka Producer
# This script tests the Python producer with minimal configuration

set -e

echo "🧪 Testing Kafka Producer Scripts"
echo "=================================="

# Check if scripts directory exists
if [ ! -d "scripts" ]; then
    echo "❌ Scripts directory not found. Run this from the project root."
    exit 1
fi

# Test 1: Check if Python script exists
echo "🐍 Testing Python producer..."
cd scripts
if [ -f "kafka_producer.py" ]; then
    echo "✅ Python producer script exists"
else
    echo "❌ Python producer script not found"
    exit 1
fi

# Test 2: Check if Python script has required dependencies
echo "📦 Testing Python producer dependencies..."
if python3 -c "import kafka" 2>/dev/null; then
    echo "✅ Python kafka-python package is available"
else
    echo "⚠️  Python kafka-python package not found"
    echo "   Install with: pip install kafka-python"
fi

# Test 3: Check if shell runner is executable
echo "🔧 Testing shell runner..."
if [ -x "run_producer.sh" ]; then
    echo "✅ Shell runner is executable"
else
    echo "⚠️  Shell runner is not executable"
    echo "   Make executable with: chmod +x run_producer.sh"
fi

# Test 4: Test help output
echo "📖 Testing help output..."
if ./run_producer.sh --help >/dev/null 2>&1; then
    echo "✅ Help command works"
else
    echo "❌ Help command failed"
fi

echo ""
echo "🎉 Basic tests completed!"
echo ""
echo "To run the producer:"
echo "  ./scripts/run_producer.sh"
echo ""
echo "To see all options:"
echo "  ./scripts/run_producer.sh --help"
echo ""
echo "Note: Make sure Kafka is running before testing with real data."
