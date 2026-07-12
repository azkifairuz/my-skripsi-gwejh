import qrcode from "qrcode-terminal";
import pkg from "whatsapp-web.js";
import { config } from "./config.js";
import { startGrpcServer } from "./grpc-server.js";
import { handleMessage } from "./message-handler.js";

const { Client, LocalAuth } = pkg;

startGrpcServer();

const client = new Client({
  authStrategy: new LocalAuth({
    dataPath: "/app/.wwebjs_auth"
  }),
  puppeteer: {
    executablePath: config.puppeteerExecutablePath,
    headless: true,
    args: [
      "--no-sandbox",
      "--disable-setuid-sandbox",
      "--disable-dev-shm-usage",
      "--disable-accelerated-2d-canvas",
      "--no-first-run",
      "--no-zygote",
      "--disable-gpu"
    ]
  }
});

client.on("qr", (qr) => {
  console.log("Scan this QR with WhatsApp:");
  qrcode.generate(qr, { small: true });
});

client.on("ready", () => {
  console.log("whatsapp client is ready");
});

client.on("authenticated", () => {
  console.log("whatsapp client authenticated");
});

client.on("auth_failure", (message) => {
  console.error("whatsapp auth failure", message);
});

client.on("message", handleMessage);

client.initialize();
