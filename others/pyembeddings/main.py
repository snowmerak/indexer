from typing import Union
from fastapi import FastAPI
from embeddings_service import EmbeddingsRequest, EmbeddingsService

app = FastAPI()


@app.get("/")
def read_root():
    return {"status": "ok"}


embeddingsService = EmbeddingsService("lemon-mint/gte-modernbert-base-code-3")

@app.post("/embed")
def read_item(request: EmbeddingsRequest):
    return {"embeddings": embeddingsService.get_embeddings(request.content)}

@app.get("/size")
def read_item():
    return {"size": embeddingsService.get_size()}