import os
import torch

MODEL_NAME = "resnet50"
EMBEDDING_DIM = 2048

GRPC_PORT = os.getenv("PYTHON_PORT", "50051")

DEVICE = torch.device("cuda" if torch.cuda.is_available() else "cpu")