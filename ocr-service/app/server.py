from concurrent import futures
import os
import subprocess
import tempfile

import grpc

from generated.proto.ocr.v1 import ocr_pb2, ocr_pb2_grpc


class OcrService(ocr_pb2_grpc.OcrServiceServicer):
    def HealthCheck(self, request, context):
        return ocr_pb2.HealthCheckResponse(status="ok")

    def Ping(self, request, context):
        return ocr_pb2.PingResponse(message="pong")

    def ExtractText(self, request, context):
        if not request.image:
            context.abort(grpc.StatusCode.INVALID_ARGUMENT, "image is required")

        language = request.language.strip() or os.getenv("TESSERACT_LANG", "ind+eng")

        try:
            text, confidence = extract_text(request.image, language)
        except subprocess.CalledProcessError as exc:
            message = exc.stderr.decode("utf-8", errors="replace").strip()
            context.abort(grpc.StatusCode.INTERNAL, f"tesseract failed: {message}")
        except Exception as exc:
            context.abort(grpc.StatusCode.INTERNAL, f"ocr failed: {exc}")

        return ocr_pb2.ExtractTextResponse(text=text, confidence=confidence)


def extract_text(image_bytes, language):
    with tempfile.NamedTemporaryFile(suffix=".image") as image_file:
        image_file.write(image_bytes)
        image_file.flush()

        tsv = subprocess.run(
            ["tesseract", image_file.name, "stdout", "-l", language, "--psm", "6", "tsv"],
            check=True,
            capture_output=True,
        ).stdout.decode("utf-8", errors="replace")

    return parse_tesseract_tsv(tsv)


def parse_tesseract_tsv(tsv):
    words = []
    confidences = []

    for line in tsv.splitlines()[1:]:
        columns = line.split("\t")
        if len(columns) < 12:
            continue

        text = columns[11].strip()
        if not text:
            continue

        words.append(text)

        try:
            confidence = float(columns[10])
        except ValueError:
            continue

        if confidence >= 0:
            confidences.append(confidence)

    average_confidence = sum(confidences) / len(confidences) if confidences else 0.0
    return " ".join(words), average_confidence


def serve():
    port = os.getenv("GRPC_PORT", "50051")
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=4))
    ocr_pb2_grpc.add_OcrServiceServicer_to_server(OcrService(), server)
    server.add_insecure_port(f"[::]:{port}")
    server.start()
    print(f"ocr gRPC server listening on port {port}", flush=True)
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
