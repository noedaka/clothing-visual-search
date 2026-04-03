import grpc
from concurrent import futures
import logging
import sys

from . import config
from . import model as model_module
from . import preprocessing
from mlpb import ml_service_pb2, ml_service_pb2_grpc


class MLServiceServicer(ml_service_pb2_grpc.MLServiceServicer):
    def __init__(self):
        logging.info(f"Loading model {config.MODEL_NAME} on device {config.DEVICE}...")
        self.model = model_module.EmbeddingModel()
        logging.info("Model loaded successfully!")

    def GetEmbedding(self, request, context):
        try:
            image_tensor = preprocessing.load_image(request.image_data)
            embedding = self.model.get_embedding(image_tensor)
            return ml_service_pb2.EmbeddingResponse(embedding=embedding)
        except Exception as e:
            logging.error(f"Error processing request: {e}")
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(f"Invalid image: {str(e)}")
            return ml_service_pb2.EmbeddingResponse()


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    ml_service_pb2_grpc.add_MLServiceServicer_to_server(MLServiceServicer(), server)
    server.add_insecure_port(f'[::]:{config.GRPC_PORT}')
    
    logging.info(f"ML gRPC server starting on port {config.GRPC_PORT}...")
    server.start()
    logging.info(f"ML gRPC server listening on port {config.GRPC_PORT}")
    
    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        server.stop(0)
        logging.info("Server stopped gracefully")


if __name__ == "__main__":
    logging.basicConfig(
        level=logging.INFO,
        format='%(asctime)s - %(levelname)s - %(message)s'
    )
    serve()