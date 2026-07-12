import { financeClient, nlpClient, ocrClient, callUnary } from "./grpc-clients.js";
import { grpc } from "./proto.js";

export async function pingDependencies() {
  const [finance, nlp, ocr] = await Promise.all([
    callUnary(financeClient, "Ping", { source: "whatsapp-service" }),
    callUnary(nlpClient, "Ping", { source: "whatsapp-service" }),
    callUnary(ocrClient, "Ping", { source: "whatsapp-service" })
  ]);

  return {
    finance: finance.message,
    nlp: nlp.message,
    ocr: ocr.message
  };
}

export async function registerByWhatsappNumber(whatsappNumber, username) {
  return callUnary(financeClient, "RegisterByWhatsappNumber", {
    whatsapp_number: whatsappNumber,
    username
  });
}

export async function loginByWhatsappNumber(whatsappNumber) {
  return callUnary(financeClient, "LoginByWhatsappNumber", {
    whatsapp_number: whatsappNumber
  });
}

export async function extractTransaction(text) {
  return callUnary(nlpClient, "ExtractTransaction", { text });
}

export async function extractText(image, mimeType, language) {
  return callUnary(ocrClient, "ExtractText", {
    image,
    mime_type: mimeType,
    language
  });
}

export async function resolveCategoryByName(categoryName) {
  return callUnary(financeClient, "ResolveCategoryByName", {
    category_name: categoryName || "lainnya"
  });
}

export async function createTransactionFromNlp(account, nlpResult, options = {}) {
  const category = await resolveCategoryByName(nlpResult.category);
  const amount = Number(nlpResult.amount || 0);
  const type = nlpResult.type || "expense";
  const description = String(nlpResult.description || options.fallbackName || "Transaksi").trim();

  return callUnary(financeClient, "CreateTransaction", {
    account_id: account.account_id,
    wallet_id: account.primary_wallet_id,
    category_id: category.category_id,
    type,
    amount,
    name: description || "Transaksi",
    is_ai_generated: true,
    receipt_image_url: options.receiptImageUrl || "",
    report_date: options.reportDate || new Date().toISOString()
  }).then((transaction) => ({
    ...transaction,
    categoryName: category.name
  }));
}

export async function listTransactionHistory(accountId, limit = 10) {
  return callUnary(financeClient, "ListTransactionHistory", {
    account_id: accountId,
    limit,
    offset: 0,
    type: "",
    from_date: "",
    to_date: ""
  });
}

export function isNotFound(error) {
  return error && error.code === grpc.status.NOT_FOUND;
}
