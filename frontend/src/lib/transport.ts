import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { PatchService } from "@/gen/tft/v1/patch_pb";
import { PlayerService } from "@/gen/tft/v1/player_pb";
import { AuthService } from "@/gen/tft/v1/auth_pb";
import { SimulationService } from "@/gen/tft/v1/simulation_pb";

export const transport = createConnectTransport({
  baseUrl: "http://localhost:8080",
  interceptors: [
    (next) => async (req) => {
      const token = localStorage.getItem("tft-oracle-auth");
      if (token) {
        try {
          const parsed = JSON.parse(token);
          if (parsed?.state?.token) {
            req.header.set("Authorization", `Bearer ${parsed.state.token}`);
          }
        } catch {
          // ignore malformed storage
        }
      }
      return next(req);
    },
  ],
});

export const patchClient = createClient(PatchService, transport);
export const playerClient = createClient(PlayerService, transport);
export const authClient = createClient(AuthService, transport);
export const simulationClient = createClient(SimulationService, transport);
