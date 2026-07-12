import { pingDependencies } from "./services.js";
import { config } from "./config.js";
import { grpc, whatsappProto } from "./proto.js";

export function startGrpcServer() {
  const server = new grpc.Server();

  server.addService(whatsappProto.WhatsappService.service, {
    HealthCheck: (_call, callback) => {
      callback(null, { status: "ok" });
    },
    PingDependencies: async (_call, callback) => {
      try {
        callback(null, await pingDependencies());
      } catch (error) {
        callback({
          code: grpc.status.UNAVAILABLE,
          message: `ping dependencies failed: ${error.message}`
        });
      }
    }
  });

  server.bindAsync(`0.0.0.0:${config.grpcPort}`, grpc.ServerCredentials.createInsecure(), (error) => {
    if (error) {
      console.error("failed to bind whatsapp gRPC server", error);
      process.exit(1);
    }

    server.start();
    console.log(`whatsapp gRPC server listening on port ${config.grpcPort}`);
  });
}
