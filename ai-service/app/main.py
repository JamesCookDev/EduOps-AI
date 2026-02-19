from fastapi import FastAPI
from pydantic import BaseModel
import uvicorn

app = FastAPI()

class MetricsPayload(BaseModel):
    container_id: str
    cpu: float
    memory: float
    restart_count: int

@app.post("/analyze")
def analyze(payload: MetricsPayload):
    return {
        "anomaly": False,
        "anomaly_score": 0.0,
        "health_modifier": 0,
        "problem_type": None,
        "explanation": None,
        "suggestion": None
    }

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8001)
