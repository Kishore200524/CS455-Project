export type UserRole = "student" | "moderator" | "administrator";

export interface AuthUser {
  id: string;
  email: string;
  role: UserRole;
  createdAt: string;
}

export interface AuthSession {
  accessToken: string;
  tokenType: "Bearer";
  expiresAt: string;
  user: AuthUser;
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null;
}

function isAuthUser(value: unknown): value is AuthUser {
  return (
    isRecord(value) &&
    typeof value.id === "string" &&
    typeof value.email === "string" &&
    (value.role === "student" ||
      value.role === "moderator" ||
      value.role === "administrator") &&
    typeof value.createdAt === "string"
  );
}

function parseSession(value: unknown): AuthSession {
  if (
    !isRecord(value) ||
    typeof value.accessToken !== "string" ||
    value.tokenType !== "Bearer" ||
    typeof value.expiresAt !== "string" ||
    !isAuthUser(value.user)
  ) {
    throw new Error("The server returned an unexpected authentication response.");
  }
  return {
    accessToken: value.accessToken,
    tokenType: "Bearer",
    expiresAt: value.expiresAt,
    user: value.user,
  };
}

async function requestJSON(
  path: string,
  body: unknown,
  token?: string,
): Promise<unknown> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const response = await fetch(path, {
    method: "POST",
    headers,
    body: JSON.stringify(body),
  });
  const result: unknown = await response.json();

  if (!response.ok) {
    const message =
      isRecord(result) &&
      "error" in result &&
      typeof result.error === "string"
        ? result.error
        : "The request could not be completed.";
    throw new Error(message);
  }
  return result;
}

export async function registerStudent(email: string, password: string) {
  return parseSession(
    await requestJSON("/api/v1/auth/register", { email, password }),
  );
}

export async function login(email: string, password: string) {
  return parseSession(
    await requestJSON("/api/v1/auth/login", { email, password }),
  );
}

export async function provisionModerator(
  token: string,
  email: string,
  password: string,
) {
  const user = await requestJSON(
    "/api/v1/admin/moderators",
    { email, password },
    token,
  );
  if (!isAuthUser(user) || user.role !== "moderator") {
    throw new Error("The server returned an unexpected moderator response.");
  }
  return user;
}

export async function logout(token: string): Promise<void> {
  const response = await fetch("/api/v1/auth/logout", {
    method: "POST",
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!response.ok) {
    throw new Error("Could not end the session on the server.");
  }
}
