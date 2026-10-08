import { useState, type FormEvent } from "react";
import type { AuthSession, AuthUser } from "../services/auth";
import { provisionModerator } from "../services/auth";

interface AdministratorPageProps {
  session: AuthSession;
  onLogout: () => void;
}

function AdministratorPage({ session, onLogout }: AdministratorPageProps) {
  const [createdModerator, setCreatedModerator] = useState<AuthUser | null>(null);
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setCreatedModerator(null);
    setIsSubmitting(true);
    const formElement = event.currentTarget;
    const form = new FormData(formElement);

    try {
      const user = await provisionModerator(
        session.accessToken,
        String(form.get("email") ?? ""),
        String(form.get("password") ?? ""),
      );
      setCreatedModerator(user);
      formElement.reset();
    } catch (requestError: unknown) {
      setError(
        requestError instanceof Error
          ? requestError.message
          : "Moderator account could not be created.",
      );
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <main className="min-h-screen bg-mist text-ink">
      <header className="border-b border-ink/10 bg-white">
        <div className="mx-auto flex max-w-5xl items-center justify-between px-6 py-5">
          <span className="text-xl font-bold tracking-tight text-ocean">
            CampusEcho
          </span>
          <button
            className="rounded-lg border border-ink/15 px-4 py-2 text-sm font-semibold hover:bg-mist"
            onClick={onLogout}
            type="button"
          >
            Sign out
          </button>
        </div>
      </header>

      <section className="mx-auto max-w-5xl px-6 py-12">
        <p className="text-sm font-bold uppercase tracking-[0.2em] text-ocean">
          Administrator
        </p>
        <h1 className="mt-3 text-3xl font-bold">Moderator accounts</h1>
        <p className="mt-3 max-w-2xl leading-7 text-ink/65">
          Create moderator accounts for authorized staff. Moderator access is
          assigned by the server, not selected during public registration.
        </p>

        <form
          className="mt-8 max-w-xl space-y-5 rounded-2xl border border-ink/10 bg-white p-6 shadow-sm"
          onSubmit={handleSubmit}
        >
          <label className="block">
            <span className="mb-2 block text-sm font-semibold">Staff email</span>
            <input
              autoComplete="email"
              className="w-full rounded-lg border border-ink/15 px-4 py-3 outline-none focus:border-ocean focus:ring-2 focus:ring-ocean/15"
              name="email"
              required
              type="email"
            />
          </label>

          <label className="block">
            <span className="mb-2 block text-sm font-semibold">
              Initial password
            </span>
            <input
              autoComplete="new-password"
              className="w-full rounded-lg border border-ink/15 px-4 py-3 outline-none focus:border-ocean focus:ring-2 focus:ring-ocean/15"
              minLength={12}
              maxLength={72}
              name="password"
              required
              type="password"
            />
            <span className="mt-2 block text-xs text-ink/50">
              Use 12–72 characters and deliver it to the staff member securely.
            </span>
          </label>

          {error && (
            <p
              className="rounded-lg bg-rose-50 px-4 py-3 text-sm text-rose-800"
              role="alert"
            >
              {error}
            </p>
          )}
          {createdModerator && (
            <p
              className="rounded-lg bg-emerald-50 px-4 py-3 text-sm text-emerald-800"
              role="status"
            >
              Moderator account created for {createdModerator.email}.
            </p>
          )}

          <button
            className="rounded-lg bg-ocean px-5 py-3 font-semibold text-white hover:bg-ocean/90 disabled:opacity-60"
            disabled={isSubmitting}
            type="submit"
          >
            {isSubmitting ? "Creating…" : "Create moderator"}
          </button>
        </form>
      </section>
    </main>
  );
}

export default AdministratorPage;
