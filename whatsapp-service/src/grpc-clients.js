import { config } from "./config.js";
import { financeProto, grpc, nlpProto, ocrProto } from "./proto.js";

const credentials = grpc.credentials.createInsecure();

export const financeClient = new financeProto.FinanceService(config.financeAddr, credentials);
export const nlpClient = new nlpProto.NlpService(config.nlpAddr, credentials);
export const ocrClient = new ocrProto.OcrService(config.ocrAddr, credentials);

export function callUnary(client, method, request) {
  return new Promise((resolve, reject) => {
    client[method](request, (error, response) => {
      if (error) {
        reject(error);
        return;
      }

      resolve(response);
    });
  });
}
