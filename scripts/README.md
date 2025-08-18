# Kafka Producer Scripts

This directory contains scripts to generate realistic order data and send it to Kafka for testing and development purposes.

## Overview

The scripts generate realistic e-commerce order data that matches your existing data models and sends it to Kafka topics. This is useful for:

- Testing your Kafka consumer
- Load testing your system
- Development and debugging
- Demonstrating the system functionality

## Available Scripts

### 1. Python Producer (`kafka_producer.py`)

- **Language**: Python 3
- **Dependencies**: `kafka-python` package
- **Performance**: Good for development and testing
- **Best for**: Quick testing, development, data analysis

### 2. Shell Runner (`run_producer.sh`)

- **Purpose**: Easy execution with configuration options
- **Features**: Command-line arguments, dependency checking, colored output
- **Best for**: Daily use, automation, CI/CD

## Quick Start

### Prerequisites

1. **Kafka running**: Make sure your Kafka cluster is running
2. **Topic exists**: Ensure the target topic exists (or enable auto-creation)

### Using the Shell Runner (Recommended)

```bash
# Make the script executable
chmod +x scripts/run_producer.sh

# Run with defaults (20 messages, 1s delay)
./scripts/run_producer.sh

# Send 100 messages with 0.5s delay
./scripts/run_producer.sh -c 100 -d 0.5

# Use custom Kafka configuration
./scripts/run_producer.sh -b kafka:9092 -t test-orders -c 50
```

### Direct Execution

#### Python Version

```bash
cd scripts
pip install -r requirements.txt
python3 kafka_producer.py
```

## Configuration

### Environment Variables

| Variable        | Default          | Description                              |
| --------------- | ---------------- | ---------------------------------------- |
| `KAFKA_BROKERS` | `localhost:9092` | Kafka broker addresses (comma-separated) |
| `KAFKA_TOPIC`   | `orders`         | Kafka topic name                         |
| `MESSAGE_COUNT` | `20`             | Number of messages to send               |
| `MESSAGE_DELAY` | `1.0`            | Delay between messages in seconds        |

### Configuration File

1. Copy the example configuration:

   ```bash
   cp scripts/config.env.example scripts/config.env
   ```

2. Modify `scripts/config.env` with your settings

3. Source the configuration:

   ```bash
   source scripts/config.env
   ```

## Command Line Options

The shell runner supports these options:

```bash
./scripts/run_producer.sh [OPTIONS]

Options:
  -b, --brokers BROKERS    Kafka brokers (default: localhost:9092)
  -t, --topic TOPIC        Kafka topic (default: orders)
  -c, --count COUNT        Number of messages (default: 20)
  -d, --delay DELAY        Delay between messages in seconds (default: 1)
  -h, --help               Show help message
```

## Examples

### Basic Usage

```bash
# Send 20 orders using Python
./scripts/run_producer.sh

# Send 50 orders
./scripts/run_producer.sh -c 50
```

### Load Testing

```bash
# Send 1000 orders with minimal delay
./scripts/run_producer.sh -c 1000 -d 0.1

# Send 500 orders with 2s delay
./scripts/run_producer.sh -c 500 -d 2
```

### Different Environments

```bash
# Development
./scripts/run_producer.sh -b localhost:9092 -t orders-dev

# Docker environment
./scripts/run_producer.sh -b kafka:9092 -t orders-prod

# Multiple brokers
./scripts/run_producer.sh -b "kafka1:9092,kafka2:9092" -t orders
```

## Generated Data Structure

The scripts generate realistic order data that matches your existing models:

### Order Structure

- **OrderUID**: Unique identifier
- **TrackNumber**: Shipping tracking number
- **Delivery**: Customer delivery information
- **Payment**: Payment details and amounts
- **Items**: 1-5 random products with realistic pricing
- **Metadata**: Timestamps, locale, customer ID, etc.

### Sample Data

- **Names**: Russian names for realistic delivery addresses
- **Cities**: Major Russian cities
- **Products**: Tech products (phones, laptops, accessories)
- **Pricing**: Realistic prices in rubles with sales
- **Addresses**: Russian-style addresses

## Docker Integration

The producer is now integrated into your Docker setup:

```bash
# Run the entire stack including producer
docker-compose up -d

# Run just the producer service
docker-compose up kafka-producer

# View producer logs
docker-compose logs kafka-producer
```

The producer service is configured to send 20 messages by default and will automatically stop after completion.

## Troubleshooting

### Common Issues

1. **Connection refused**

   - Check if Kafka is running
   - Verify broker address and port
   - Ensure firewall allows connections

2. **Topic not found**

   - Check if topic exists
   - Enable auto-creation: `KAFKA_AUTO_CREATE_TOPICS_ENABLE=true`

3. **Python import errors**

   - Install dependencies: `pip install -r requirements.txt`
   - Check Python version (requires Python 3.6+)

### Debug Mode

For debugging, you can modify the scripts to:

- Add more verbose logging
- Reduce message count for testing
- Use longer delays between messages

## Performance Considerations

- **Python version**: Good for hundreds of messages per second
- **Network**: Ensure sufficient bandwidth for high-volume testing
- **Kafka**: Monitor broker performance during load testing

## Integration with Your System

The generated data is compatible with your existing:

- Kafka consumer (`internal/kafka/consumer.go`)
- Order service (`internal/service/order/`)
- Database models and migrations
- API endpoints

## Contributing

To add new data types or modify the generation logic:

1. Update the sample data arrays in the scripts
2. Modify the generation functions
3. Test with your consumer
4. Update this documentation

## License

These scripts are part of your project and follow the same license terms.
