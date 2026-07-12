import path from "node:path";
import grpc from "@grpc/grpc-js";
import protoLoader from "@grpc/proto-loader";
import { config } from "./config.js";

const loaderOptions = {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true
};

function loadProto(relativePath) {
  const definition = protoLoader.loadSync(path.join(config.protoDir, relativePath), {
    ...loaderOptions,
    includeDirs: [config.protoDir]
  });

  return grpc.loadPackageDefinition(definition);
}

export const financeProto = loadProto("finance/v1/finance.proto").finance.v1;
export const nlpProto = loadProto("nlp/v1/nlp.proto").nlp.v1;
export const ocrProto = loadProto("ocr/v1/ocr.proto").ocr.v1;
export const whatsappProto = loadProto("whatsapp/v1/whatsapp.proto").whatsapp.v1;
export { grpc };
