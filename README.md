# clothing-visual-search

**Clothing visual search engine** - a demo e‑commerce website that lets users search for products by uploading an image.  
The system uses a **ResNet50** neural network to extract feature vectors (embeddings) and finds visually similar items in the catalogue.

## Architecture

- **Frontend**: React (Vite) - upload image, display results  
- **Backend**: Go - REST API, product management, orchestration  
- **ML Service**: Python (gRPC) - ResNet50 inference, returns embeddings  
- **Vector Database**: Milvus - stores and searches embeddings  
- **Object Storage**: MinIO - stores product images  
- **Relational DB**: PostgreSQL - product metadata  

