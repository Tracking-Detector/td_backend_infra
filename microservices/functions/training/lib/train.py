from lib.model import get_all_available_models
def train_model(modelName: str, dataSetName: str):
    print(f"Training model: {modelName} with data: {dataSetName}")
    for model in get_all_available_models():
        if model.model_name == modelName:
            model.train(dataSetName, 128, 10)