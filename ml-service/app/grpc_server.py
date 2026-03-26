import grpc
from concurrent import futures
import logging

from . import config
from . import model as model_module
from . import preprocessing

from mlpb import ml_service_pb2, ml_service_pb2_grpc

class MLServiceServicer(ml_service_pb2_grpc.MLServiceServicer):
    def __init__(self):
        self.model = model_module.EmbeddingModel()

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
    logging.info(f"ML gRPC server listening on port {config.GRPC_PORT}")
    server.start()
    server.wait_for_termination()