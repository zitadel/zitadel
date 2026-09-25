import { create } from "@zitadel/client";
import { GetUserByIDResponseSchema } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { getUserAvatarUrl } from "./user-avatar";
import { getUserByID } from "./zitadel";

vi.mock("./zitadel", () => ({ getUserByID: vi.fn() }));

const serviceConfig = { baseUrl: "https://idp.example.com" };

describe("getUserAvatarUrl", () => {
  beforeEach(() => vi.resetAllMocks());

  it("loads the human profile image by the session user ID", async () => {
    vi.mocked(getUserByID).mockResolvedValue(
      create(GetUserByIDResponseSchema, {
        user: { type: { case: "human", value: { profile: { avatarUrl: "https://idp.example.com/assets/v1/photo" } } } },
      }),
    );
    expect(await getUserAvatarUrl({ serviceConfig, userId: "user-1" })).toBe("https://idp.example.com/assets/v1/photo");
    expect(getUserByID).toHaveBeenCalledWith({ serviceConfig, userId: "user-1" });
  });

  it("does not look up unidentified users", async () => {
    expect(await getUserAvatarUrl({ serviceConfig })).toBeUndefined();
    expect(getUserByID).not.toHaveBeenCalled();
  });

  it("returns no image for a machine user", async () => {
    vi.mocked(getUserByID).mockResolvedValue(
      create(GetUserByIDResponseSchema, {
        user: { type: { case: "machine", value: {} } },
      }),
    );
    expect(await getUserAvatarUrl({ serviceConfig, userId: "machine-1" })).toBeUndefined();
  });

  it("returns no image for a missing profile", async () => {
    vi.mocked(getUserByID).mockResolvedValue(create(GetUserByIDResponseSchema));
    expect(await getUserAvatarUrl({ serviceConfig, userId: "deleted" })).toBeUndefined();
  });

  it("keeps login available when the profile lookup fails", async () => {
    vi.mocked(getUserByID).mockRejectedValue(new Error("unavailable"));
    expect(await getUserAvatarUrl({ serviceConfig, userId: "user-1" })).toBeUndefined();
  });
});
