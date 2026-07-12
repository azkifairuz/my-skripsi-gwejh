import {
  createTransactionFromNlp,
  extractText,
  extractTransaction,
  isNotFound,
  listTransactionHistory,
  loginByWhatsappNumber,
  registerByWhatsappNumber
} from "./services.js";
import { normalizeWhatsappNumber, summarizeTransaction, toRupiah } from "./format.js";
import { config } from "./config.js";

export async function handleMessage(message) {
  if (message.fromMe || message.from.endsWith("@g.us")) {
    return;
  }

  const whatsappNumber = normalizeWhatsappNumber(message.from);
  const body = String(message.body || "").trim();

  try {
    if (body.startsWith("/register")) {
      await handleRegister(message, whatsappNumber, body);
      return;
    }

    if (body === "/login") {
      await handleLogin(message, whatsappNumber);
      return;
    }

    if (body.startsWith("/history")) {
      await handleHistory(message, whatsappNumber, body);
      return;
    }

    if (message.hasMedia) {
      await handleReceiptImage(message, whatsappNumber);
      return;
    }

    if (body.length > 0) {
      await handleTextTransaction(message, whatsappNumber, body);
      return;
    }
  } catch (error) {
    console.error("message handling failed", error);
    await message.reply("Maaf, proses gagal. Coba beberapa saat lagi.");
  }
}

async function handleRegister(message, whatsappNumber, body) {
  const username = body.replace(/^\/register/i, "").trim();
  if (!username) {
    await message.reply("Format: /register Nama Kamu");
    return;
  }

  const account = await registerByWhatsappNumber(whatsappNumber, username);
  await message.reply(account.created
    ? `Registrasi berhasil. Halo ${account.username || username}.`
    : `Nomor ini sudah terdaftar. Halo ${account.username || username}.`);
}

async function handleLogin(message, whatsappNumber) {
  try {
    const account = await loginByWhatsappNumber(whatsappNumber);
    await message.reply(`Login berhasil. Account: ${account.account_id}`);
  } catch (error) {
    if (isNotFound(error)) {
      await message.reply("Nomor belum terdaftar. Ketik /register Nama Kamu");
      return;
    }
    throw error;
  }
}

async function handleHistory(message, whatsappNumber, body) {
  const account = await requireRegisteredAccount(message, whatsappNumber);
  if (!account) {
    return;
  }

  const limit = parseLimit(body);
  const history = await listTransactionHistory(account.account_id, limit);
  if (!history.transactions || history.transactions.length === 0) {
    await message.reply("Belum ada transaksi.");
    return;
  }

  const lines = history.transactions.map((item, index) => {
    return `${index + 1}. ${item.name} - ${toRupiah(item.amount)} (${item.type}, ${item.category_name})`;
  });
  await message.reply(["History transaksi:", ...lines].join("\n"));
}

async function handleTextTransaction(message, whatsappNumber, text) {
  const account = await requireRegisteredAccount(message, whatsappNumber);
  if (!account) {
    return;
  }

  const nlpResult = await extractTransaction(text);
  const transaction = await createTransactionFromNlp(account, nlpResult, {
    fallbackName: text,
    reportDate: new Date().toISOString()
  });

  await message.reply(summarizeTransaction(transaction));
}

async function handleReceiptImage(message, whatsappNumber) {
  const account = await requireRegisteredAccount(message, whatsappNumber);
  if (!account) {
    return;
  }

  const media = await message.downloadMedia();
  if (!media || !media.mimetype || !media.mimetype.startsWith("image/")) {
    await message.reply("File harus berupa gambar struk.");
    return;
  }

  const image = Buffer.from(media.data, "base64");
  const ocrResult = await extractText(image, media.mimetype, config.ocrLanguage);
  if (!ocrResult.text || !ocrResult.text.trim()) {
    await message.reply("Teks struk tidak terbaca. Coba foto yang lebih jelas.");
    return;
  }

  const nlpResult = await extractTransaction(ocrResult.text);
  const transaction = await createTransactionFromNlp(account, nlpResult, {
    fallbackName: ocrResult.text.slice(0, 80),
    reportDate: new Date().toISOString()
  });

  const confidence = Number(ocrResult.confidence || 0).toFixed(1);
  await message.reply([
    `Struk terbaca. Confidence OCR: ${confidence}%`,
    summarizeTransaction(transaction)
  ].join("\n\n"));
}

async function requireRegisteredAccount(message, whatsappNumber) {
  try {
    return await loginByWhatsappNumber(whatsappNumber);
  } catch (error) {
    if (isNotFound(error)) {
      await message.reply("Nomor belum terdaftar. Ketik /register Nama Kamu");
      return null;
    }
    throw error;
  }
}

function parseLimit(body) {
  const [, rawLimit] = body.split(/\s+/);
  const limit = Number.parseInt(rawLimit, 10);
  if (Number.isNaN(limit) || limit <= 0) {
    return 10;
  }
  return Math.min(limit, 20);
}
