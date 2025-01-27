from pydantic import BaseModel
from sentence_transformers import SentenceTransformer


class EmbeddingsService:
    def __init__(self, model):
        self.model_name = model
        self.model = SentenceTransformer(model, cache_folder="./.cache")

    def get_embeddings(self, content: str) -> [float]:
        encoded = self.model.encode(content)
        return encoded.tolist()

    def get_size(self):
        return len(self.model.encode("test"))

class EmbeddingsRequest(BaseModel):
    content: str
    model: str
    api_key: str