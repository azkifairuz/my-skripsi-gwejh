export const config = {
  grpcPort: process.env.GRPC_PORT || "50052",
  financeAddr: process.env.FINANCE_SERVICE_ADDR || "finance-service:50053",
  nlpAddr: process.env.NLP_SERVICE_ADDR || "nlp-service:50054",
  ocrAddr: process.env.OCR_SERVICE_ADDR || "ocr-service:50051",
  ocrLanguage: process.env.OCR_LANGUAGE || "ind+eng",
  protoDir: process.env.PROTO_DIR || "/app/proto",
  puppeteerExecutablePath: process.env.PUPPETEER_EXECUTABLE_PATH || "/usr/bin/chromium"
};
