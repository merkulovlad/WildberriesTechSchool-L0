#!/usr/bin/env python3
"""
Kafka Producer Script for Generating Order Data

This script generates realistic order data and sends it to Kafka.
It can be configured via environment variables or command line arguments.
"""

import json
import os
import random
import time
import uuid
from datetime import datetime
from typing import List, Dict, Any
from datetime import datetime, timezone

try:
    from kafka import KafkaProducer
    from kafka.errors import KafkaError
except ImportError:
    print("Error: kafka-python package not found. Install it with: pip install kafka-python")
    exit(1)

# Configuration
class Config:
    def __init__(self):
        self.kafka_brokers = os.getenv('KAFKA_BROKERS', 'localhost:29092').split(',')
        self.topic = os.getenv('KAFKA_TOPIC', 'orders')
        self.message_count = int(os.getenv('MESSAGE_COUNT', '20'))
        self.delay = float(os.getenv('MESSAGE_DELAY', '1.0'))

# Sample data for generating realistic orders
NAMES = [
    "Иван Иванов", "Мария Петрова", "Алексей Сидоров", "Елена Козлова",
    "Дмитрий Волков", "Анна Морозова", "Сергей Соколов", "Ольга Лебедева",
    "Николай Козлов", "Татьяна Новикова", "Андрей Морозов", "Наталья Петрова",
    "Владимир Соловьев", "Екатерина Васильева", "Михаил Зайцев", "Ирина Семенова",
    "Александр Голубев", "Людмила Виноградова", "Виктор Богданов", "Галина Воробьева",
]

CITIES = [
    "Москва", "Санкт-Петербург", "Новосибирск", "Екатеринбург", "Казань",
    "Нижний Новгород", "Челябинск", "Самара", "Ростов-на-Дону", "Уфа",
    "Волгоград", "Пермь", "Воронеж", "Краснодар", "Саратов",
]

REGIONS = [
    "Московская область", "Ленинградская область", "Свердловская область",
    "Ростовская область", "Краснодарский край", "Татарстан", "Башкортостан",
    "Челябинская область", "Самарская область", "Нижегородская область",
]

DELIVERY_SERVICES = [
    "СДЭК", "Boxberry", "Почта России", "DHL", "FedEx", "UPS",
    "Яндекс.Доставка", "СберЛогистика", "ПЭК", "Деловые Линии",
]

ITEM_NAMES = [
    "Смартфон iPhone 15 Pro", "Ноутбук MacBook Air", "Наушники AirPods Pro",
    "Планшет iPad Air", "Умные часы Apple Watch", "Телевизор Samsung QLED",
    "Игровая консоль PlayStation 5", "Фотоаппарат Canon EOS R",
    "Беспроводная колонка JBL", "Электронная книга Kindle",
    "Монитор Dell UltraSharp", "Клавиатура Logitech MX Keys",
    "Мышь Logitech MX Master", "Веб-камера Logitech StreamCam",
    "Микрофон Blue Yeti", "Принтер HP LaserJet", "Сканер Epson Perfection",
    "МФУ Canon Pixma", "Внешний жесткий диск WD", "SSD накопитель Samsung",
]

BRANDS = [
    "Apple", "Samsung", "Sony", "LG", "Canon", "Nikon", "Dell", "HP",
    "Lenovo", "Asus", "Logitech", "JBL", "Bose", "Sennheiser", "WD",
    "Seagate", "Kingston", "Corsair", "Razer", "SteelSeries",
]

LOCALES = ["ru", "en", "de", "fr", "es"]

# Global counter for sequential order IDs
order_counter = 0

def generate_random_string(length: int) -> str:
    """Generate a random string of specified length."""
    import string
    return ''.join(random.choices(string.ascii_letters + string.digits, k=length))

def generate_phone_number() -> str:
    """Generate a realistic Russian phone number."""
    return f"+7{random.randint(1000000000, 9999999999)}"

def generate_email() -> str:
    """Generate a random email address."""
    domains = ["gmail.com", "yandex.ru", "mail.ru", "outlook.com"]
    username = generate_random_string(8)
    domain = random.choice(domains)
    return f"{username}@{domain}"

def generate_order_uid() -> str:
    """Generate a sequential order UID like b1, b2, b3, etc."""
    global order_counter
    order_counter += 1
    return f"b{order_counter}"

def generate_track_number() -> str:
    """Generate a track number."""
    return f"WBILMT{random.randint(100000, 999999)}"

def generate_random_item() -> Dict[str, Any]:
    """Generate a random item for an order."""
    price = random.randint(10, 110) * 1000  # 10k - 110k rubles
    sale = random.randint(0, 30)  # 0-30% sale
    total_price = price - (price * sale // 100)
    
    return {
        "chrt_id": random.randint(100000, 999999),
        "track_number": generate_track_number(),
        "price": price,
        "rid": generate_random_string(10),
        "name": random.choice(ITEM_NAMES),
        "sale": sale,
        "size": str(random.randint(1, 50)),
        "total_price": total_price,
        "nm_id": random.randint(100000, 999999),
        "brand": random.choice(BRANDS),
        "status": random.randint(1, 5)
    }

def generate_random_order() -> Dict[str, Any]:
    """Generate a complete random order."""
    now = datetime.now()
    
    # Generate random items (1-5 items)
    item_count = random.randint(1, 5)
    items = [generate_random_item() for _ in range(item_count)]
    
    # Calculate goods total
    goods_total = sum(item["total_price"] for item in items)
    
    # Generate delivery info
    delivery = {
        "name": random.choice(NAMES),
        "phone": generate_phone_number(),
        "zip": str(random.randint(100000, 999999)),
        "city": random.choice(CITIES),
        "address": f"ул. {generate_random_string(10)}, д. {random.randint(1, 100)}, кв. {random.randint(1, 100)}",
        "region": random.choice(REGIONS),
        "email": generate_email()
    }
    
    # Generate payment info
    payment = {
        "transaction": generate_random_string(20),
        "request_id": generate_random_string(15),
        "currency": "RUB",
        "provider": "wbpay",
        "amount": goods_total + random.randint(1000, 6000),  # goods + delivery
        "payment_dt": int(now.timestamp()),
        "bank": "alpha",
        "delivery_cost": random.randint(1000, 6000),
        "goods_total": goods_total,
        "custom_fee": random.randint(0, 1000)
    }
    
    # Create the complete order
    order = {
        "order_uid": generate_order_uid(),
        "track_number": generate_track_number(),
        "entry": "WBILMT",
        "delivery": delivery,
        "payment": payment,
        "items": items,
        "locale": random.choice(LOCALES),
        "internal_signature": "",
        "customer_id": f"customer_{random.randint(10000, 99999)}",
        "delivery_service": random.choice(DELIVERY_SERVICES),
        "shardkey": f"shard_{random.randint(0, 9)}",
        "sm_id": random.randint(0, 99),
        "date_created": rfc3339_now(),
        "oof_shard": f"oof_{random.randint(0, 9)}"
    }
    
    return order

def rfc3339_now() -> str:
    # RFC3339 with microseconds and trailing Z
    return datetime.now(timezone.utc).isoformat(timespec="microseconds").replace("+00:00", "Z")

def delivery_report(err, msg):
    """Callback function for message delivery reports."""
    if err is not None:
        print(f"Message delivery failed: {err}")
    else:
        print(f"Message delivered to {msg.topic} [{msg.partition}] at offset {msg.offset}")

def main():
    """Main function to run the Kafka producer."""
    global order_counter
    order_counter = 0  # Reset counter for each run
    
    config = Config()
    
    print(f"Starting Kafka producer...")
    print(f"Brokers: {config.kafka_brokers}")
    print(f"Topic: {config.topic}")
    print(f"Message count: {config.message_count}")
    print(f"Delay between messages: {config.delay} seconds")
    print(f"Order IDs will be: b1, b2, b3, ..., b{config.message_count}")
    
    # Create Kafka producer
    try:
        producer = KafkaProducer(
            bootstrap_servers=config.kafka_brokers,
            value_serializer=lambda v: json.dumps(v, ensure_ascii=False).encode('utf-8'),
            key_serializer=lambda k: k.encode('utf-8') if k else None,
            acks='all',  # Wait for all replicas to acknowledge
            retries=3,   # Retry failed sends
        )
    except Exception as e:
        print(f"Failed to create Kafka producer: {e}")
        return
    
    try:
        # Generate and send messages
        for i in range(config.message_count):
            order = generate_random_order()
            
            # Send message to Kafka
            future = producer.send(
                topic=config.topic,
                key=order['order_uid'],
                value=order
            )
            
            # Add delivery report callback
            future.add_callback(delivery_report)
            
            print(f"Sent order {order['order_uid']} ({i+1}/{config.message_count})")
            
            # Wait before sending next message (except for the last one)
            if i < config.message_count - 1:
                time.sleep(config.delay)
        
        # Wait for all messages to be delivered
        producer.flush()
        print(f"Finished sending {config.message_count} messages to Kafka")
        
    except KeyboardInterrupt:
        print("\nInterrupted by user")
    except Exception as e:
        print(f"Error occurred: {e}")
    finally:
        producer.close()

if __name__ == "__main__":
    main()
