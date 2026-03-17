import { useMutation, useQuery } from "@tanstack/react-query";
import { ConnectError, Code } from "@connectrpc/connect";
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

export function useRegisterAndLogin() {
  const setAuth = useAuthStore((s) => s.setAuth);

  return useMutation({
    mutationFn: async (params: {
      gameName: string;
      tagLine: string;
      region: string;
    }) => {
      const registerRes = await authClient.register(params);
      const loginRes = await authClient.login({
        accessKey: registerRes.accessKey,
      });
      if (loginRes.user && loginRes.sessionToken) {
        setAuth(loginRes.sessionToken, loginRes.user);
      }
      return {
        accessKey: registerRes.accessKey,
        user: loginRes.user,
      };
    },
  });
}

export function useCurrentUser() {
  const token = useAuthStore((s) => s.token);
  const logout = useAuthStore((s) => s.logout);

  return useQuery({
    queryKey: ["currentUser", token],
    queryFn: async () => {
      try {
        return await authClient.getCurrentUser({});
      } catch (err) {
        // Stale/invalid token → clear auth state so UI falls back to onboarding
        if (
          err instanceof ConnectError &&
          err.code === Code.Unauthenticated
        ) {
          logout();
        }
        throw err;
      }
    },
    enabled: !!token,
    retry: false,
    staleTime: 5 * 60 * 1000,
  });
}
