import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { PatchService } from "../gen/tft/v1/patch_pb";
import { PlayerService } from "../gen/tft/v1/player_pb";

export const transport = createConnectTransport({
  baseUrl: "http://localhost:8080",
});

export const patchClient = createClient(PatchService, transport);
export const playerClient = createClient(PlayerService, transport);
