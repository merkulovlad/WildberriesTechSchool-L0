#!/bin/bash

# Kafka Producer Runner Script
# This script provides easy ways to run the Kafka producer with different configurations

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default configuration
DEFAULT_BROKERS="kafka:29092"
DEFAULT_TOPIC="orders"
DEFAULT_COUNT=20
DEFAULT_DELAY=1

# Function to print usage
print_usage() {
    echo -e "${BLUE}Kafka Producer Runner Script${NC}"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "OPTIONS:"
    echo "  -b, --brokers BROKERS    Kafka brokers (default: $DEFAULT_BROKERS)"
    echo "  -t, --topic TOPIC        Kafka topic (default: $DEFAULT_TOPIC)"
    echo "  -c, --count COUNT        Number of messages (default: $DEFAULT_COUNT)"
    echo "  -d, --delay DELAY        Delay between messages in seconds (default: $DEFAULT_DELAY)"
    echo "  -h, --help               Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                        # Run with defaults (20 messages, 1s delay)"
    echo "  $0 -c 100 -d 0.5         # Send 100 messages with 0.5s delay"
    echo "  $0 -b kafka:9092 -t test # Use custom broker and topic"
    echo ""
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check dependencies
check_dependencies() {
    if ! command_exists python3; then
        echo -e "${RED}Error: python3 is not installed${NC}"
        exit 1
    fi
    
    if ! python3 -c "import kafka" 2>/dev/null; then
        echo -e "${YELLOW}Warning: kafka-python package not found${NC}"
        echo -e "${YELLOW}Installing dependencies...${NC}"
        pip3 install -r "$(dirname "$0")/requirements.txt"
    fi
}

# Function to run Python producer
run_python_producer() {
    local brokers=$1
    local topic=$2
    local count=$3
    local delay=$4
    
    echo -e "${GREEN}Running Python Kafka Producer...${NC}"
    
    # Set environment variables
    export KAFKA_BROKERS="$brokers"
    export KAFKA_TOPIC="$topic"
    export MESSAGE_COUNT="$count"
    export MESSAGE_DELAY="$delay"
    
    # Run the Python script
    python3 "$(dirname "$0")/kafka_producer.py"
}

# Parse command line arguments
BROKERS=$DEFAULT_BROKERS
TOPIC=$DEFAULT_TOPIC
COUNT=$DEFAULT_COUNT
DELAY=$DEFAULT_DELAY

while [[ $# -gt 0 ]]; do
    case $1 in
        -b|--brokers)
            BROKERS="$2"
            shift 2
            ;;
        -t|--topic)
            TOPIC="$2"
            shift 2
            ;;
        -c|--count)
            COUNT="$2"
            shift 2
            ;;
        -d|--delay)
            DELAY="$2"
            shift 2
            ;;
        -h|--help)
            print_usage
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            print_usage
            exit 1
            ;;
    esac
done

# Validate inputs
if ! [[ "$COUNT" =~ ^[0-9]+$ ]] || [ "$COUNT" -lt 1 ]; then
    echo -e "${RED}Error: Count must be a positive integer${NC}"
    exit 1
fi

if ! [[ "$DELAY" =~ ^[0-9]+\.?[0-9]*$ ]] || [ "$(echo "$DELAY < 0" | bc -l)" -eq 1 ]; then
    echo -e "${RED}Error: Delay must be a non-negative number${NC}"
    exit 1
fi

# Check dependencies
check_dependencies

# Print configuration
echo -e "${BLUE}Configuration:${NC}"
echo "  Script Type: Python"
echo "  Kafka Brokers: $BROKERS"
echo "  Topic: $TOPIC"
echo "  Message Count: $COUNT"
echo "  Delay: ${DELAY}s"
echo ""

# Run the Python producer
run_python_producer "$BROKERS" "$TOPIC" "$COUNT" "$DELAY"
