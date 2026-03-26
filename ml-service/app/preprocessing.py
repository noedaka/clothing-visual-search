import io
from PIL import Image
import torch
import torchvision.transforms as transforms

preprocess = transforms.Compose([
    transforms.Resize(256),
    transforms.CenterCrop(224),
    transforms.ToTensor(),
    transforms.Normalize(mean=[0.485, 0.456, 0.406],
                         std=[0.229, 0.224, 0.225])
])

def load_image(image_bytes: bytes) -> torch.Tensor:
    try:
        img = Image.open(io.BytesIO(image_bytes)).convert('RGB')
    except Exception as e:
        raise ValueError(f"Invalid image data: {e}")
    return preprocess(img).unsqueeze(0)