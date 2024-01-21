import json
import os
import pika
import time
from pymongo import MongoClient
from lib.model import get_all_available_models
from lib.train import train_model
import logging

logger = logging.getLogger("Consumer")
logging.basicConfig(level=logging.INFO)


def seed_models_in_db():
    logger.info("Start seeding mongodb.")
    db_client = MongoClient(os.getenv("MONGO_URI"))
    db = db_client["tracking-detector"]
    collection = db[os.getenv("MODELS_COLLECTION")]
    for model in get_all_available_models():
        collection.update_one(
            {"name": model.model_name},
            {"$set": {"name": model.model_name,
                      "description": model.model_desc, "dims": model.input_dims}},
            upsert=True
        )
    logger.info("Finished seeding mongodb.")


def callback(ch, method, properties, body):
    # Deserialize the body of the message
    job_data = json.loads(body.decode('utf-8'))
    function_name = job_data['functionName']
    args = job_data['args']

    if function_name == 'train_model':
        train_model(*args)

    ch.basic_ack(delivery_tag=method.delivery_tag)


def main():
    seed_models_in_db()
    time.sleep(30)

    connection = pika.BlockingConnection(
        pika.ConnectionParameters(
            host='rabbitmq', port=5672, credentials=pika.PlainCredentials('guest', 'guest'))
    )
    channel = connection.channel()

    channel.queue_declare(queue='training', durable=True)

    channel.basic_consume(queue='training', on_message_callback=callback)

    channel.start_consuming()


if __name__ == "__main__":
    main()