import gzip
import os
import random
import shutil
import string
import numpy as np
import pandas as pd
from sklearn.model_selection import train_test_split
from sklearn.metrics import f1_score
from tensorflow import keras
from keras.models import Sequential
from keras.layers import Dense, Activation, Dropout, Embedding, Flatten, Masking, ReLU
from minio import Minio
from pymongo import MongoClient
from datetime import datetime
import tensorflowjs as tfjs
import tarfile
import logging

logger = logging.getLogger("Model")
logging.basicConfig(level=logging.INFO)


class Model:
    def __init__(self, model_name: str, model_desc: str, input_dims):
        self.model_name = model_name
        self.model_desc = model_desc
        self.input_dims = input_dims
        pass

    def compiled_model(self):
        pass

    def verify_bucket_exists(self, client, bucket_name):
        if client.bucket_exists(bucket_name):
            return
        client.make_bucket(bucket_name, "eu-central-1")
        client.set_bucket_versioning(bucket_name, {"Status": "Enabled"})

    def train(self, dataset, batch_size, epochs):
        try:
            logger.info("Start connecting to minio service.")
            client = Minio(
                os.getenv("MINIO_URI"),
                access_key=os.getenv("MINIO_ACCESS_KEY"),
                secret_key=os.getenv("MINIO_PRIVATE_KEY"),
                secure=False,
            )
            logger.info("Connection successful")
            export_bucket = os.getenv("EXPORT_BUCKET_NAME")
            model_bucket = os.getenv("MODEL_BUCKET_NAME")
            logger.info("Start verifing bucket")
            self.verify_bucket_exists(client, export_bucket)
            self.verify_bucket_exists(client, model_bucket)
            logger.info("Verification successful")
            logger.info("Start loading training data into dataframe.")

            file_data = client.get_object(export_bucket, dataset + ".csv.gz")
            gzip_file = gzip.GzipFile(fileobj=file_data)
            # csv_content = gzip_file.read().decode("utf-8")

            df = pd.read_csv(gzip_file)
            logger.info("Finished loading training data into dataframe.")

            X = df.iloc[:, 0:(self.input_dims[0] + 1)].values
            y = df.iloc[:, -1].values

            X_train, X_test, y_train, y_test = train_test_split(
                X, y, test_size=0.33, random_state=42
            )

            model = self.compiled_model()

            history = model.fit(
                X_train,
                y_train,
                batch_size=batch_size,
                epochs=epochs,
                validation_data=(X_test, y_test),
            )

            pred_train = np.rint(model.predict(X_train))
            f1_train = f1_score(y_train, pred_train)

            pred_test = np.rint(model.predict(X_test))
            f1_test = f1_score(y_test, pred_test)
            logger.info("Finished Training")

            logger.info("Store Training Results to MongoDb")
            db_client = MongoClient(os.getenv("MONGO_URI"))
            db = db_client["tracking-detector"]
            collection = db[os.getenv("TRAINING_RUNS_COLLECTION")]
            collection.insert_one(
                {
                    "name": self.model_name,
                    "dataSet": dataset,
                    "time": datetime.now().isoformat(),
                    "f1Train": f1_train,
                    "f1Test": f1_test,
                    "trainingHistory": history.history,
                    "batchSize": batch_size,
                    "epochs": epochs
                }
            )
            logger.info("Inserted Training Results")
            logger.info(f"Starting export of Model '{self.model_name}'")
            tempName = generate_temp_filename()
            tfjs.converters.save_keras_model(model, tempName)

            logger.info("Start uploading raw model to minio")
            for root, _, files in os.walk("./" + tempName):
                for file in files:
                    file_path = os.path.join(root, file)
                    client.fput_object(
                        model_bucket,
                        self.model_name +"/" + dataset + "/" + file,
                        file_path,
                    )
            tar = tarfile.open(f"{self.model_name}_{dataset}.tar.gz", "w:gz")
            tar.add(tempName, arcname=self.model_name)
            tar.close()
            client.fput_object(
                model_bucket,
                self.model_name + "/" + f"{self.model_name}_{dataset}.tar.gz",
                f"./{self.model_name}_{dataset}.tar.gz",
            )
            logger.info("Uploaded files to minio")
            shutil.rmtree(f"./{tempName}", ignore_errors=False, onerror=None)
            os.remove(f"./{self.model_name}_{dataset}.tar.gz")
            logger.info("Cleaned up folders")
        except Exception as e:
            print(e)
            logger.error(
                f"An error occurred while trying to train the model '{self.model_name}'.",
                e,
            )


def generate_temp_filename() -> str:
    return "".join(random.choices(string.ascii_uppercase + string.digits, k=10))


class ReferrerBasedUrlClassifier(Model):
    def __init__(self):
        super().__init__(
            "ReferrerBasedUrlClassifier",
            "This model determines a tracking request by encoding 200 chars of the url and using the request type, request method and the frametype.",
            [203, 1],
        )

    def compiled_model(self):
        model = Sequential()
        model.add(Embedding(90, 32, input_length=204, mask_zero=True))
        model.add(Flatten())
        model.add(Dense(512, input_shape=(6528,)))
        model.add(Dropout(0.5))
        model.add(ReLU())
        model.add(Dense(256, input_shape=(512,)))
        model.add(Dropout(0.5))
        model.add(ReLU())
        model.add(Dense(128, input_shape=(256,)))
        model.add(Dropout(0.5))
        model.add(ReLU())
        model.add(Dense(1, input_shape=(128,)))
        model.add(Activation("sigmoid"))

        model.compile(
            loss="mse",
            optimizer=keras.optimizers.RMSprop(learning_rate=1e-2),
            metrics=["accuracy"],
        )
        return model


def get_all_available_models():
    refBased = ReferrerBasedUrlClassifier()
    models = [refBased]
    return models