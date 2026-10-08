import { useEffect, useState } from "react";

export type ApiHealth = "checking" | "online" | "offline";

export function useApiHealth(): ApiHealth {
  const [health, setHealth] = useState<ApiHealth>("checking");

  useEffect(() => {
    const controller = new AbortController();

    fetch("/api/v1/health", { signal: controller.signal })
      .then((response) => {
        if (!response.ok) {
          throw new Error("API health check failed");
        }
        setHealth("online");
      })
      .catch((error: unknown) => {
        if (error instanceof Error && error.name === "AbortError") {
          return;
        }
        setHealth("offline");
      });

    return () => controller.abort();
  }, []);

  return health;
}
