from fastapi import FastAPI
from fastapi.responses import HTMLResponse

app = FastAPI()


@app.get("/")
def get_home():
    return HTMLResponse("<h1>FastApi backend</h1>")


@app.get("/health")
def get_health():
    return "healthy"
