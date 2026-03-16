import { useMutation, useQuery } from "@tanstack/react-query";
import { authClient } from "@/lib/transport";
import { useAuthStore } from "@/stores/auth-store";

export function useRegister() {
  return useMutation({
    mutationFn: async (params: {
      gameName: string;
      tagLine: string;
      region: string;
    }) => {
      const res = await authClient.register(params);
      return res;
    },
  });
}

export function useLogin() {
  const setAuth = useAuthStore((s) => s.setAuth);

  return useMutation({
    mutationFn: async (accessKey: string) => {
      const res = await authClient.login({ accessKey });
      if (res.user && res.sessionToken) {
        setAuth(res.sessionToken, res.user);
      }
      return res;
    },
  });
}

export function useCurrentUser() {
  const token = useAuthStore((s) => s.token);

  return useQuery({
    queryKey: ["currentUser", token],
    queryFn: () => authClient.getCurrentUser({}),
    enabled: !!token,
    retry: false,
    staleTime: 5 * 60 * 1000,
  });
}
