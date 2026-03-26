import logging
from .grpc_server import serve

if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    serve()