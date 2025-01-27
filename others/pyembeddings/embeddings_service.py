from gen.rpc.embeddings_pb2 import GetEmbeddingsResponse, GetEmbeddingsRequest
from gen.rpc.embeddings_pb2_grpc import EmbeddingsServiceServicer
from sentence_transformers import SentenceTransformer


class EmbeddingsService(EmbeddingsServiceServicer):
    def __init__(self, model):
        self.model_name = model
        self.model = model = SentenceTransformer(model, cache_folder="./.cache")

    def GetEmbeddings(self, request: GetEmbeddingsRequest, context):
        encoded = self.model.encode(request.contents)
        return GetEmbeddingsResponse(embeddings=encoded[0])
