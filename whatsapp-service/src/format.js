export function normalizeWhatsappNumber(value) {
  return String(value || "")
    .replace("@c.us", "")
    .replace("@s.whatsapp.net", "")
    .replace(/[+\s-]/g, "")
    .trim();
}

export function toRupiah(value) {
  const amount = Number(value || 0);
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0
  }).format(amount);
}

export function summarizeTransaction(transaction) {
  return [
    "Transaksi tersimpan.",
    `Jenis: ${transaction.type}`,
    `Nominal: ${toRupiah(transaction.amount)}`,
    `Kategori: ${transaction.categoryName || transaction.category_name || "-"}`,
    `Nama: ${transaction.name}`,
    `ID: ${transaction.transaction_id || transaction.transactionId}`
  ].join("\n");
}
