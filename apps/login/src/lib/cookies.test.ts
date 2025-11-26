import { describe, expect, it, vi, beforeEach, afterEach } from "vitest";
import {
  Cookie,
  addSessionToCookie,
  updateSessionCookie,
  removeSessionFromCookie,
  getMostRecentSessionCookie,
  getSessionCookieById,
  getSessionCookieByLoginName,
  getAllSessionCookieIds,
  getAllSessions,
  getMostRecentCookieWithLoginname,
  setLanguageCookie,
} from "./cookies";

// Mock dependencies
vi.mock("next/headers", () => ({
  cookies: vi.fn(),
}));

vi.mock("@zitadel/client", () => ({
  timestampDate: vi.fn((ts: any) => new Date(Number(ts.seconds) * 1000)),
  timestampFromMs: vi.fn((ms: number) => ({ seconds: BigInt(Math.floor(ms / 1000)) }) as any),
}));

import { cookies } from "next/headers";
import { timestampDate, timestampFromMs } from "@zitadel/client";

describe("cookies", () => {
  let mockCookies: any;

  beforeEach(() => {
    vi.clearAllMocks();
    mockCookies = {
      get: vi.fn(),
      set: vi.fn(),
    };
    vi.mocked(cookies).mockResolvedValue(mockCookies);
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("setLanguageCookie", () => {
    it("should set language cookie with correct parameters", async () => {
      await setLanguageCookie("en");

      expect(mockCookies.set).toHaveBeenCalledWith({
        name: "NEXT_LOCALE",
        value: "en",
        httpOnly: true,
        path: "/",
      });
    });
  });

  describe("addSessionToCookie", () => {
    const mockSession: Cookie = {
      id: "session-1",
      token: "token-1",
      loginName: "user@example.com",
      organization: "org-1",
      creationTs: "1700000000000",
      expirationTs: "1800000000000",
      changeTs: "1700000000000",
    };

    it("should add new session to empty cookie", async () => {
      mockCookies.get.mockReturnValue(undefined);

      await addSessionToCookie({ session: mockSession });

      expect(mockCookies.set).toHaveBeenCalledWith(
        expect.objectContaining({
          name: "sessions",
          value: JSON.stringify([mockSession]),
        }),
      );
    });

    it("should prepend new session to existing sessions", async () => {
      const existingSession: Cookie = {
        id: "session-0",
        token: "token-0",
        loginName: "other@example.com",
        creationTs: "1600000000000",
        expirationTs: "1700000000000",
        changeTs: "1600000000000",
      };

      mockCookies.get.mockReturnValue({
        value: JSON.stringify([existingSession]),
      });

      await addSessionToCookie({ session: mockSession });

      const expectedSessions = [mockSession, existingSession];
      expect(mockCookies.set).toHaveBeenCalledWith(
        expect.objectContaining({
          value: JSON.stringify(expectedSessions),
        }),
      );
    });

    it("should update existing session with same loginName", async () => {
      const existingSession: Cookie = {
        id: "session-old",
        token: "token-old",
        loginName: "user@example.com",
        creationTs: "1600000000000",
        expirationTs: "1700000000000",
        changeTs: "1600000000000",
      };

      mockCookies.get.mockReturnValue({
        value: JSON.stringify([existingSession]),
      });

      const updatedSession: Cookie = {
        ...mockSession,
        id: "session-new",
        token: "token-new",
      };

      await addSessionToCookie({ session: updatedSession });

      const setCall = mockCookies.set.mock.calls[0][0];
      const sessions = JSON.parse(setCall.value);

      expect(sessions).toHaveLength(1);
      expect(sessions[0].id).toBe("session-new");
      expect(sessions[0].token).toBe("token-new");
    });

    it("should handle cookie overflow by replacing oldest session", async () => {
      // Create many sessions to exceed MAX_COOKIE_SIZE (2048)
      const manySessions: Cookie[] = Array.from({ length: 10 }, (_, i) => ({
        id: `session-${i}`,
        token: `token-${i}-very-long-token-to-increase-size-padding-padding-padding`,
        loginName: `user${i}@example.com`,
        creationTs: `${1600000000000 + i * 1000}`,
        expirationTs: `${1700000000000 + i * 1000}`,
        changeTs: `${1600000000000 + i * 1000}`,
      }));

      mockCookies.get.mockReturnValue({
        value: JSON.stringify(manySessions),
      });

      const consoleSpy = vi.spyOn(console, "log").mockImplementation(() => {});

      await addSessionToCookie({ session: mockSession });

      expect(consoleSpy).toHaveBeenCalledWith("WARNING COOKIE OVERFLOW");

      const setCall = mockCookies.set.mock.calls[0][0];
      const sessions = JSON.parse(setCall.value);

      // Should have new session and all but the first old session
      expect(sessions[0]).toEqual(mockSession);
      expect(sessions).not.toContain(manySessions[0]);

      consoleSpy.mockRestore();
    });

    it("should cleanup expired sessions when cleanup is true", async () => {
      const now = Date.now();
      const expiredSession: Cookie = {
        id: "session-expired",
        token: "token-expired",
        loginName: "expired@example.com",
        creationTs: `${now - 10000000}`,
        expirationTs: `${now - 1000}`, // Expired
        changeTs: `${now - 10000000}`,
      };

      const validSession: Cookie = {
        id: "session-valid",
        token: "token-valid",
        loginName: "valid@example.com",
        creationTs: `${now}`,
        expirationTs: `${now + 10000000}`, // Not expired
        changeTs: `${now}`,
      };

      mockCookies.get.mockReturnValue({
        value: JSON.stringify([expiredSession, validSession]),
      });

      vi.mocked(timestampDate).mockImplementation((ts: any) => new Date(Number(ts.seconds) * 1000));
      vi.mocked(timestampFromMs).mockImplementation((ms: number) => ({ seconds: BigInt(Math.floor(ms / 1000)) }) as any);

      await addSessionToCookie({ session: mockSession, cleanup: true });

      const setCall = mockCookies.set.mock.calls[0][0];
      const sessions = JSON.parse(setCall.value);

      // Should not include expired session
      const hasExpired = sessions.some((s: Cookie) => s.id === "session-expired");
      expect(hasExpired).toBe(false);
    });

    it("should set sameSite to 'none' when iFrameEnabled is true", async () => {
      mockCookies.get.mockReturnValue(undefined);

      await addSessionToCookie({ session: mockSession, iFrameEnabled: true });

      expect(mockCookies.set).toHaveBeenCalledWith(
        expect.objectContaining({
          sameSite: "none",
        }),
      );
    });

    it("should set sameSite to 'lax' when iFrameEnabled is false or undefined", async () => {
      mockCookies.get.mockReturnValue(undefined);

      await addSessionToCookie({ session: mockSession, iFrameEnabled: false });

      expect(mockCookies.set).toHaveBeenCalledWith(
        expect.objectContaining({
          sameSite: "lax",
        }),
      );

      mockCookies.set.mockClear();
      await addSessionToCookie({ session: mockSession });

      expect(mockCookies.set).toHaveBeenCalledWith(
        expect.objectContaining({
          sameSite: "lax",
        }),
      );
    });
  });

  describe("updateSessionCookie", () => {
    const mockSession: Cookie = {
      id: "session-1",
      token: "token-1",
      loginName: "user@example.com",
      creationTs: "1700000000000",
      expirationTs: "1800000000000",
      changeTs: "1700000000000",
    };

    it("should update existing session by id", async () => {
      mockCookies.get.mockReturnValue({
        value: JSON.stringify([mockSession]),
      });

      const updatedSession: Cookie = {
        ...mockSession,
        token: "new-token",
        changeTs: "1700000001000",
      };

      await updateSessionCookie({
        id: "session-1",
        session: updatedSession,
      });

      const setCall = mockCookies.set.mock.calls[0][0];
      const sessions = JSON.parse(setCall.value);

      expect(sessions[0].token).toBe("new-token");
      expect(sessions[0].changeTs).toBe("1700000001000");
    });

    it("should throw error if session id not found", async () => {
      mockCookies.get.mockReturnValue({
        value: JSON.stringify([mockSession]),
      });

      await expect(
        updateSessionCookie({
          id: "non-existent-id",
          session: mockSession,
        }),
      ).rejects.toThrow("updateSessionCookie<T>: session id not found");
    });

    it("should cleanup expired sessions when cleanup is true", async () => {
      const now = Date.now();
      const expiredSession: Cookie = {
        id: "session-expired",
        token: "token-expired",
        loginName: "expired@example.com",
        creationTs: `${now - 10000000}`,
        expirationTs: `${now - 1000}`,
        changeTs: `${now - 10000000}`,
      };

      mockCookies.get.mockReturnValue({
        value: JSON.stringify([mockSession, expiredSession]),
      });

      vi.mocked(timestampDate).mockImplementation((ts: any) => new Date(Number(ts.seconds) * 1000));
      vi.mocked(timestampFromMs).mockImplementation((ms: number) => ({ seconds: BigInt(Math.floor(ms / 1000)) }) as any);

      await updateSessionCookie({
        id: "session-1",
        session: { ...mockSession, token: "updated" },
        cleanup: true,
      });

      const setCall = mockCookies.set.mock.calls[0][0];
      const sessions = JSON.parse(setCall.value);

      const hasExpired = sessions.some((s: Cookie) => s.id === "session-expired");
      expect(hasExpired).toBe(false);
    });

    it("should respect iFrameEnabled parameter", async () => {
      mockCookies.get.mockReturnValue({
        value: JSON.stringify([mockSession]),
      });

      await updateSessionCookie({
        id: "session-1",
        session: mockSession,
        iFrameEnabled: true,
      });

      expect(mockCookies.set).toHaveBeenCalledWith(
        expect.objectContaining({
          sameSite: "none",
        }),
      );
    });
  });

  describe("removeSessionFromCookie", () => {
    const session1: Cookie = {
      id: "session-1",
      token: "token-1",
      loginName: "user1@example.com",
      creationTs: "1700000000000",
      expirationTs: "1800000000000",
      changeTs: "1700000000000",
    };

    const session2: Cookie = {
      id: "session-2",
      token: "token-2",
      loginName: "user2@example.com",
      creationTs: "1700000001000",
      expirationTs: "1800000001000",
      changeTs: "1700000001000",
    };

    it("should remove session by id", async () => {
      mockCookies.get.mockReturnValue({
        value: JSON.stringify([session1, session2]),
      });

      await removeSessionFromCookie({ session: session1 });

      const setCall = mockCookies.set.mock.calls[0][0];
      const sessions = JSON.parse(setCall.value);

      expect(sessions).toHaveLength(1);
      expect(sessions[0].id).toBe("session-2");
    });

    it("should handle removing non-existent session gracefully", async () => {
      mockCookies.get.mockReturnValue({
        value: JSON.stringify([session1]),
      });

      const nonExistentSession: Cookie = {
        ...session2,
        id: "non-existent",
      };

      await removeSessionFromCookie({ session: nonExistentSession });

      const setCall = mockCookies.set.mock.calls[0][0];
      const sessions = JSON.parse(setCall.value);

      expect(sessions).toHaveLength(1);
      expect(sessions[0].id).toBe("session-1");
    });

    it("should cleanup expired sessions when cleanup is true", async () => {
      const now = Date.now();
      const expiredSession: Cookie = {
        id: "session-expired",
        token: "token-expired",
        loginName: "expired@example.com",
        creationTs: `${now - 10000000}`,
        expirationTs: `${now - 1000}`,
        changeTs: `${now - 10000000}`,
      };

      mockCookies.get.mockReturnValue({
        value: JSON.stringify([session1, session2, expiredSession]),
      });

      vi.mocked(timestampDate).mockImplementation((ts: any) => new Date(Number(ts.seconds) * 1000));
      vi.mocked(timestampFromMs).mockImplementation((ms: number) => ({ seconds: BigInt(Math.floor(ms / 1000)) }) as any);

      await removeSessionFromCookie({ session: session1, cleanup: true });

      const setCall = mockCookies.set.mock.calls[0][0];
      const sessions = JSON.parse(setCall.value);

      expect(sessions).toHaveLength(1);
      expect(sessions[0].id).toBe("session-2");
    });

    it("should respect iFrameEnabled parameter", async () => {
      mockCookies.get.mockReturnValue({
        value: JSON.stringify([session1, session2]),
      });

      await removeSessionFromCookie({
        session: session1,
        iFrameEnabled: true,
      });

      expect(mockCookies.set).toHaveBeenCalledWith(
        expect.objectContaining({
          sameSite: "none",
        }),
      );
    });
  });

  describe("getMostRecentSessionCookie", () => {
    it("should return session with most recent changeTs", async () => {
      const session1: Cookie = {
        id: "session-1",
        token: "token-1",
        loginName: "user1@example.com",
        creationTs: "1700000000000",
        expirationTs: "1800000000000",
        changeTs: "1700000000000",
      };

      const session2: Cookie = {
        id: "session-2",
        token: "token-2",
        loginName: "user2@example.com",
        creationTs: "1700000001000",
        expirationTs: "1800000001000",
        changeTs: "1700000005000", // Most recent
      };

      mockCookies.get.mockReturnValue({
        value: JSON.stringify([session1, session2]),
      });

      const result = await getMostRecentSessionCookie();

      expect(result.id).toBe("session-2");
    });

    it("should reject when no session cookie exists", async () => {
      mockCookies.get.mockReturnValue(undefined);

      await expect(getMostRecentSessionCookie()).rejects.toBe("no session cookie found");
    });

    it("should handle single session", async () => {
      const session: Cookie = {
        id: "session-1",
        token: "token-1",
        loginName: "user@example.com",
        creationTs: "1700000000000",
        expirationTs: "1800000000000",
        changeTs: "1700000000000",
      };

      mockCookies.get.mockReturnValue({
        value: JSON.stringify([session]),
      });

      const result = await getMostRecentSessionCookie();

      expect(result.id).toBe("session-1");
    });
  });

  describe("getSessionCookieById", () => {
    const session1: Cookie = {
      id: "session-1",
      token: "token-1",
      loginName: "user1@example.com",
      organization: "org-1",
      creationTs: "1700000000000",
      expirationTs: "1800000000000",
      changeTs: "1700000000000",
    };

    const session2: Cookie = {
      id: "session-2",
      token: "token-2",
      loginName: "user2@example.com",
      organization: "org-2",
      creationTs: "1700000001000",
      expirationTs: "1800000001000",
      changeTs: "1700000001000",
    };

    it("should find session by id", async () => {
      mockCookies.get.mockReturnValue({
        value: JSON.stringify([session1, session2]),
      });

      const result = await getSessionCookieById({ sessionId: "session-1" });

      expect(result.id).toBe("session-1");
      expect(result.loginName).toBe("user1@example.com");
    });

    it("should filter by organization when provided", async () => {
      mockCookies.get.mockReturnValue({
        value: JSON.stringify([session1, session2]),
      });

      const result = await getSessionCookieById({
        sessionId: "session-2",
        organization: "org-2",
      });

      expect(result.id).toBe("session-2");
      expect(result.organization).toBe("org-2");
    });

    it("should reject if session not found", async () => {
      mockCookies.get.mockReturnValue({
        value: JSON.stringify([session1]),
      });

      await expect(getSessionCookieById({ sessionId: "non-existent" })).rejects.toBeTypeOf("undefined");
    });

    it("should reject if organization doesn't match", async () => {
      mockCookies.get.mockReturnValue({
        value: JSON.stringify([session1]),
      });

      await expect(
        getSessionCookieById({
          sessionId: "session-1",
          organization: "wrong-org",
        }),
      ).rejects.toBeTypeOf("undefined");
    });

    it("should reject when no session cookie exists", async () => {
      mockCookies.get.mockReturnValue(undefined);

      await expect(getSessionCookieById({ sessionId: "session-1" })).rejects.toBeTypeOf("undefined");
    });
  });

  describe("getSessionCookieByLoginName", () => {
    const session1: Cookie = {
      id: "session-1",
      token: "token-1",
      loginName: "user1@example.com",
      organization: "org-1",
      creationTs: "1700000000000",
      expirationTs: "1800000000000",
      changeTs: "1700000000000",
    };

    const session2: Cookie = {
      id: "session-2",
      token: "token-2",
      loginName: "user2@example.com",
      organization: "org-2",
      creationTs: "1700000001000",
      expirationTs: "1800000001000",
      changeTs: "1700000001000",
    };

    it("should find session by loginName", async () => {
      mockCookies.get.mockReturnValue({
        value: JSON.stringify([session1, session2]),
      });

      const result = await getSessionCookieByLoginName({
        loginName: "user1@example.com",
      });

      expect(result.id).toBe("session-1");
      expect(result.loginName).toBe("user1@example.com");
    });

    it("should filter by organization when provided", async () => {
      mockCookies.get.mockReturnValue({
        value: JSON.stringify([session1, session2]),
      });

      const result = await getSessionCookieByLoginName({
        loginName: "user2@example.com",
        organization: "org-2",
      });

      expect(result.id).toBe("session-2");
      expect(result.organization).toBe("org-2");
    });

    it("should reject if session not found", async () => {
      mockCookies.get.mockReturnValue({
        value: JSON.stringify([session1]),
      });

      await expect(getSessionCookieByLoginName({ loginName: "nonexistent@example.com" })).rejects.toBe(
        "no cookie found with loginName: nonexistent@example.com",
      );
    });

    it("should reject when no session cookie exists", async () => {
      mockCookies.get.mockReturnValue(undefined);

      await expect(getSessionCookieByLoginName({ loginName: "user@example.com" })).rejects.toBe("no session cookie found");
    });
  });

  describe("getAllSessionCookieIds", () => {
    it("should return all session IDs", async () => {
      const sessions: Cookie[] = [
        {
          id: "session-1",
          token: "token-1",
          loginName: "user1@example.com",
          creationTs: "1700000000000",
          expirationTs: "1800000000000",
          changeTs: "1700000000000",
        },
        {
          id: "session-2",
          token: "token-2",
          loginName: "user2@example.com",
          creationTs: "1700000001000",
          expirationTs: "1800000001000",
          changeTs: "1700000001000",
        },
      ];

      mockCookies.get.mockReturnValue({
        value: JSON.stringify(sessions),
      });

      const result = await getAllSessionCookieIds();

      expect(result).toEqual(["session-1", "session-2"]);
    });

    it("should filter expired sessions when cleanup is true", async () => {
      const now = Date.now();
      const sessions: Cookie[] = [
        {
          id: "session-valid",
          token: "token-valid",
          loginName: "valid@example.com",
          creationTs: `${now}`,
          expirationTs: `${now + 10000000}`,
          changeTs: `${now}`,
        },
        {
          id: "session-expired",
          token: "token-expired",
          loginName: "expired@example.com",
          creationTs: `${now - 10000000}`,
          expirationTs: `${now - 1000}`,
          changeTs: `${now - 10000000}`,
        },
      ];

      mockCookies.get.mockReturnValue({
        value: JSON.stringify(sessions),
      });

      vi.mocked(timestampDate).mockImplementation((ts: any) => new Date(Number(ts.seconds) * 1000));
      vi.mocked(timestampFromMs).mockImplementation((ms: number) => ({ seconds: BigInt(Math.floor(ms / 1000)) }) as any);

      const result = await getAllSessionCookieIds(true);

      expect(result).toEqual(["session-valid"]);
    });

    it("should return empty array when no sessions exist", async () => {
      mockCookies.get.mockReturnValue(undefined);

      const result = await getAllSessionCookieIds();

      expect(result).toEqual([]);
    });
  });

  describe("getAllSessions", () => {
    it("should return all sessions", async () => {
      const sessions: Cookie[] = [
        {
          id: "session-1",
          token: "token-1",
          loginName: "user1@example.com",
          creationTs: "1700000000000",
          expirationTs: "1800000000000",
          changeTs: "1700000000000",
        },
        {
          id: "session-2",
          token: "token-2",
          loginName: "user2@example.com",
          creationTs: "1700000001000",
          expirationTs: "1800000001000",
          changeTs: "1700000001000",
        },
      ];

      mockCookies.get.mockReturnValue({
        value: JSON.stringify(sessions),
      });

      const result = await getAllSessions();

      expect(result).toEqual(sessions);
    });

    it("should filter expired sessions when cleanup is true", async () => {
      const now = Date.now();
      const validSession: Cookie = {
        id: "session-valid",
        token: "token-valid",
        loginName: "valid@example.com",
        creationTs: `${now}`,
        expirationTs: `${now + 10000000}`,
        changeTs: `${now}`,
      };

      const expiredSession: Cookie = {
        id: "session-expired",
        token: "token-expired",
        loginName: "expired@example.com",
        creationTs: `${now - 10000000}`,
        expirationTs: `${now - 1000}`,
        changeTs: `${now - 10000000}`,
      };

      mockCookies.get.mockReturnValue({
        value: JSON.stringify([validSession, expiredSession]),
      });

      vi.mocked(timestampDate).mockImplementation((ts: any) => new Date(Number(ts.seconds) * 1000));
      vi.mocked(timestampFromMs).mockImplementation((ms: number) => ({ seconds: BigInt(Math.floor(ms / 1000)) }) as any);

      const result = await getAllSessions(true);

      expect(result).toHaveLength(1);
      expect(result[0].id).toBe("session-valid");
    });

    it("should return empty array when no sessions exist", async () => {
      mockCookies.get.mockReturnValue(undefined);
      const consoleSpy = vi.spyOn(console, "log").mockImplementation(() => {});

      const result = await getAllSessions();

      expect(result).toEqual([]);
      expect(consoleSpy).toHaveBeenCalledWith("getAllSessions: No session cookie found, returning empty array");

      consoleSpy.mockRestore();
    });
  });

  describe("getMostRecentCookieWithLoginname", () => {
    it("should return most recent session for loginName", async () => {
      const session1: Cookie = {
        id: "session-1",
        token: "token-1",
        loginName: "user@example.com",
        creationTs: "1700000000000",
        expirationTs: "1800000000000",
        changeTs: "1700000000000",
      };

      const session2: Cookie = {
        id: "session-2",
        token: "token-2",
        loginName: "user@example.com",
        creationTs: "1700000001000",
        expirationTs: "1800000001000",
        changeTs: "1700000005000", // Most recent
      };

      mockCookies.get.mockReturnValue({
        value: JSON.stringify([session1, session2]),
      });

      const result = await getMostRecentCookieWithLoginname({
        loginName: "user@example.com",
      });

      expect(result.id).toBe("session-2");
    });

    it("should filter by organization when provided", async () => {
      const session1: Cookie = {
        id: "session-1",
        token: "token-1",
        loginName: "user@example.com",
        organization: "org-1",
        creationTs: "1700000000000",
        expirationTs: "1800000000000",
        changeTs: "1700000005000",
      };

      const session2: Cookie = {
        id: "session-2",
        token: "token-2",
        loginName: "user@example.com",
        organization: "org-2",
        creationTs: "1700000001000",
        expirationTs: "1800000001000",
        changeTs: "1700000010000", // More recent but different org
      };

      mockCookies.get.mockReturnValue({
        value: JSON.stringify([session1, session2]),
      });

      const result = await getMostRecentCookieWithLoginname({
        loginName: "user@example.com",
        organization: "org-1",
      });

      expect(result.id).toBe("session-1");
      expect(result.organization).toBe("org-1");
    });

    it("should return most recent session when no loginName provided", async () => {
      const session1: Cookie = {
        id: "session-1",
        token: "token-1",
        loginName: "user1@example.com",
        creationTs: "1700000000000",
        expirationTs: "1800000000000",
        changeTs: "1700000000000",
      };

      const session2: Cookie = {
        id: "session-2",
        token: "token-2",
        loginName: "user2@example.com",
        creationTs: "1700000001000",
        expirationTs: "1800000001000",
        changeTs: "1700000005000",
      };

      mockCookies.get.mockReturnValue({
        value: JSON.stringify([session1, session2]),
      });

      const result = await getMostRecentCookieWithLoginname({});

      expect(result.id).toBe("session-2");
    });

    it("should reject when no matching session found", async () => {
      const session: Cookie = {
        id: "session-1",
        token: "token-1",
        loginName: "other@example.com",
        creationTs: "1700000000000",
        expirationTs: "1800000000000",
        changeTs: "1700000000000",
      };

      mockCookies.get.mockReturnValue({
        value: JSON.stringify([session]),
      });

      await expect(getMostRecentCookieWithLoginname({ loginName: "user@example.com" })).rejects.toBe(
        "Could not get the context or retrieve a session",
      );
    });

    it("should reject when no session cookie exists", async () => {
      mockCookies.get.mockReturnValue(undefined);

      await expect(getMostRecentCookieWithLoginname({ loginName: "user@example.com" })).rejects.toBe(
        "Could not read session cookie",
      );
    });
  });
});
