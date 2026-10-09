import { useState } from "react";
import AdministratorPage from "./pages/AdministratorPage";
import AuthPage from "./pages/AuthPage";
import StudentDashboardPage from "./pages/StudentDashboardPage";
import { logout, type AuthSession } from "./services/auth";

function App() {
  const [session, setSession] = useState<AuthSession | null>(null);
  const [sessionError, setSessionError] = useState("");

  async function handleLogout() {
    if (session) {
      try {
        await logout(session.accessToken ?? "");
        setSession(null);
        setSessionError("");
      } catch (error: unknown) {
        setSessionError(
          error instanceof Error
            ? error.message
            : "Could not sign out. Please try again.",
        );
        setSession(null);
      }
    }
  }

  if (!session) {
    return (
      <AuthPage
        notice={sessionError}
        onAuthenticated={(nextSession) => {
          setSession(nextSession);
          setSessionError("");
        }}
      />
    );
  }
  if (session.user.role === "student") {
    return <StudentDashboardPage session={session} onLogout={handleLogout} />;
  }
  if (session.user.role === "administrator") {
    return <AdministratorPage session={session} onLogout={handleLogout} />;
  }
  return (
    <main className="flex min-h-screen items-center justify-center bg-mist px-6 text-ink">
      <section className="max-w-lg rounded-2xl border border-ink/10 bg-white p-8 text-center shadow-sm">
        <p className="text-sm font-bold uppercase tracking-[0.2em] text-ocean">
          Moderator
        </p>
        <h1 className="mt-3 text-2xl font-bold">Moderator dashboard pending</h1>
        <p className="mt-3 leading-7 text-ink/65">
          Your moderator account is authenticated. The moderation queue and
          ticket workflows have not been implemented yet.
        </p>
        <button
          className="mt-6 rounded-lg bg-ocean px-5 py-3 font-semibold text-white hover:bg-ocean/90"
          onClick={handleLogout}
          type="button"
        >
          Sign out
        </button>
      </section>
    </main>
  );
}

export default App;
