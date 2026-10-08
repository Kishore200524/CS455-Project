import { useState, type FormEvent } from "react";
import { login, registerStudent, type AuthSession } from "../services/auth";

interface AuthPageProps {
  onAuthenticated: (session: AuthSession) => void;
  notice: string;
}

function AuthPage({ onAuthenticated, notice }: AuthPageProps) {
  const [isRegistering, setIsRegistering] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState("");

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setIsSubmitting(true);
    const form = new FormData(event.currentTarget);
    const email = String(form.get("email") ?? "");
    const password = String(form.get("password") ?? "");

    try {
      const session = isRegistering
        ? await registerStudent(email, password)
        : await login(email, password);
      onAuthenticated(session);
    } catch (requestError: unknown) {
      setError(
        requestError instanceof Error
          ? requestError.message
          : "Authentication failed. Please try again.",
      );
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <main className="flex min-h-screen items-center justify-center bg-mist px-6 py-12 text-ink">
      <section className="w-full max-w-md rounded-2xl border border-ink/10 bg-white p-7 shadow-sm sm:p-9">
        <a className="text-xl font-bold tracking-tight text-ocean" href="/">
          CampusEcho
        </a>
        <h1 className="mt-8 text-3xl font-bold">
          {isRegistering ? "Create your student account" : "Welcome back"}
        </h1>
        <p className="mt-2 text-sm leading-6 text-ink/60">
          {isRegistering
            ? "Register with your college email. Your account identity is kept separate from feedback."
            : "Sign in with your account. Moderator accounts are provisioned by an administrator."}
        </p>
        {notice && (
          <p
            className="mt-5 rounded-lg bg-amber-50 px-4 py-3 text-sm text-amber-900"
            role="status"
          >
            {notice}
          </p>
        )}

        <form className="mt-7 space-y-5" onSubmit={handleSubmit}>
          <label className="block">
            <span className="mb-2 block text-sm font-semibold">College email</span>
            <input
              autoComplete="email"
              className="w-full rounded-lg border border-ink/15 px-4 py-3 outline-none transition placeholder:text-ink/35 focus:border-ocean focus:ring-2 focus:ring-ocean/15"
              name="email"
              placeholder="you@college.edu"
              required
              type="email"
            />
          </label>

          <label className="block">
            <span className="mb-2 block text-sm font-semibold">Password</span>
            <input
              autoComplete={isRegistering ? "new-password" : "current-password"}
              className="w-full rounded-lg border border-ink/15 px-4 py-3 outline-none transition placeholder:text-ink/35 focus:border-ocean focus:ring-2 focus:ring-ocean/15"
              minLength={isRegistering ? 12 : undefined}
              maxLength={72}
              name="password"
              required
              type="password"
            />
            {isRegistering && (
              <span className="mt-2 block text-xs text-ink/50">
                Use 12–72 characters. This starter validates email format but
                does not yet verify enrollment against a college directory.
              </span>
            )}
          </label>

          {error && (
            <p
              className="rounded-lg bg-rose-50 px-4 py-3 text-sm text-rose-800"
              role="alert"
            >
              {error}
            </p>
          )}

          <button
            className="w-full rounded-lg bg-ocean px-5 py-3 font-semibold text-white transition hover:bg-ocean/90 focus:outline-none focus:ring-2 focus:ring-ocean focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-60"
            disabled={isSubmitting}
            type="submit"
          >
            {isSubmitting
              ? "Please wait…"
              : isRegistering
                ? "Create account"
                : "Sign in"}
          </button>
        </form>

        <p className="mt-6 text-center text-sm text-ink/60">
          {isRegistering ? "Already registered?" : "New student?"}{" "}
          <button
            className="font-semibold text-ocean hover:underline"
            onClick={() => {
              setIsRegistering(!isRegistering);
              setError("");
            }}
            type="button"
          >
            {isRegistering ? "Sign in" : "Create an account"}
          </button>
        </p>
      </section>
    </main>
  );
}

export default AuthPage;
