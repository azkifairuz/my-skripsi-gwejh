FROM python:3.13-slim

WORKDIR /app

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        protobuf-compiler \
        tesseract-ocr \
        tesseract-ocr-eng \
        tesseract-ocr-ind \
    && rm -rf /var/lib/apt/lists/*

COPY ocr-service/requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

COPY ocr-service/app ./app
COPY proto/ocr/v1/ocr.proto ./proto/ocr/v1/ocr.proto

RUN mkdir -p generated/proto/ocr/v1 \
    && touch generated/__init__.py generated/proto/__init__.py generated/proto/ocr/__init__.py generated/proto/ocr/v1/__init__.py \
    && python -m grpc_tools.protoc \
        -I proto \
        --python_out=generated/proto \
        --grpc_python_out=generated/proto \
        ocr/v1/ocr.proto \
    && sed -i 's/^from ocr.v1 import ocr_pb2/from generated.proto.ocr.v1 import ocr_pb2/' generated/proto/ocr/v1/ocr_pb2_grpc.py

ENTRYPOINT ["python", "-m", "app.server"]
