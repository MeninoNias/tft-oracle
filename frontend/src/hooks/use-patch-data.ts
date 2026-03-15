import { useQuery } from "@tanstack/react-query";
import { patchClient } from "../lib/transport";

export function usePatchData(setNumber = 0) {
  return useQuery({
    queryKey: ["patchData", setNumber],
    queryFn: () => patchClient.getPatchData({ setNumber }),
  });
}
