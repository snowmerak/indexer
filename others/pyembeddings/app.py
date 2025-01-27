from concurrent import futures

from embeddings_service import EmbeddingsService
from gen.rpc.embeddings_pb2_grpc import add_EmbeddingsServiceServicer_to_server


embedding_model = "lemon-mint/gte-modernbert-base-code-3"

def serve(grpc=None):
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    add_EmbeddingsServiceServicer_to_server(EmbeddingsService(embedding_model), server)
    server.add_insecure_port("[::]:50051")
    server.start()
    server.wait_for_termination()