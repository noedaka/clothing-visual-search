import torch
import torchvision.models as models
from . import config

class EmbeddingModel:
    def __init__(self):
        self.device = config.DEVICE
        self.model = self._load_model()
        self.model.to(self.device)
        self.model.eval()

    def _load_model(self):
        if config.MODEL_NAME == "resnet50":
            model = models.resnet50(weights=models.ResNet50_Weights.IMAGENET1K_V1)
            model = torch.nn.Sequential(*list(model.children())[:-1])
        else:
            raise ValueError(f"Unsupported model: {config.MODEL_NAME}")
        
        return model

    def get_embedding(self, image_tensor: torch.Tensor) -> list:
        with torch.no_grad():
            image_tensor = image_tensor.to(self.device)
            embedding = self.model(image_tensor)
            embedding = embedding.cpu().squeeze().numpy().tolist()
            
        return embedding