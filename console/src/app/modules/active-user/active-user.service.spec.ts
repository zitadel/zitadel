import { beforeEach, describe, expect } from "vitest";
import { ActiveUserEntrySchema } from "@zitadel/proto/zitadel/analytics/v2beta/active_user_service_pb";
import { create } from "@bufbuild/protobuf";
import {
  ActiveUserService,
  averageActiveUserEntries,
  dateToTimestamp,
} from "@/modules/active-user/active-user.service";
import { skipToken } from "@tanstack/angular-query-experimental";
import {
  ActiveUserGrpcMockProviderService,
  ActiveUserGrpcProviderService,
} from "@/modules/active-user/active-user-grpc-provider.service";
import { TestBed } from "@angular/core/testing";
import { GrpcService } from "@/services/grpc.service";

describe("#averageActiveUsers", () => {
  test("#averageActiveUsers should return correct average for 3 entries", async () => {
    const testData = [
      create(ActiveUserEntrySchema, {
        value: BigInt(100),
      }),
      create(ActiveUserEntrySchema, {
        value: BigInt(200),
      }),
      create(ActiveUserEntrySchema, {
        value: BigInt(150),
      }),
    ];

    expect(averageActiveUserEntries(testData)).toBe(BigInt(150));
  });

  test("#averageActiveUsers should return correct average for an empty array", async () => {
    expect(averageActiveUserEntries([])).toBe(BigInt(0));
  });
});

describe("#dateToTimestamp", () => {
  test("#dateToTimestamp should convert Date to Timestamp correctly", () => {
    // Test with a specific date: January 1, 2024 12:00:00 UTC
    const testDate = new Date("2024-01-01T12:00:00.000Z");
    const timestamp = dateToTimestamp(testDate);

    // Expected: 1704110400 seconds (January 1, 2024 12:00:00 UTC)
    expect(timestamp.seconds).toBe(BigInt(1704110400));
    expect(timestamp.nanos).toBe(0);
  });

  test("#dateToTimestamp should handle milliseconds correctly", () => {
    // Test with a date that has milliseconds: January 1, 2024 12:00:00.123 UTC
    const testDate = new Date("2024-01-01T12:00:00.123Z");
    const timestamp = dateToTimestamp(testDate);

    // Expected: 1704110400 seconds + 123000000 nanoseconds
    expect(timestamp.seconds).toBe(BigInt(1704110400));
    expect(timestamp.nanos).toBe(123000000);
  });

  test("#dateToTimestamp should handle edge case with maximum milliseconds", () => {
    // Test with a date that has 999 milliseconds
    const testDate = new Date("2024-01-01T12:00:00.999Z");
    const timestamp = dateToTimestamp(testDate);

    // Expected: 1704110400 seconds + 999000000 nanoseconds
    expect(timestamp.seconds).toBe(BigInt(1704110400));
    expect(timestamp.nanos).toBe(999000000);
  });

  test("#dateToTimestamp should handle epoch time correctly", () => {
    // Test with epoch time: January 1, 1970 00:00:00 UTC
    const testDate = new Date("1970-01-01T00:00:00.000Z");
    const timestamp = dateToTimestamp(testDate);

    expect(timestamp.seconds).toBe(BigInt(0));
    expect(timestamp.nanos).toBe(0);
  });
});

describe("ActiveUserService", () => {
  let activeUserService: ActiveUserService;

  beforeEach(async () => {
    TestBed.configureTestingModule({
      providers: [
        {
          provide: ActiveUserGrpcProviderService,
          useClass: ActiveUserGrpcMockProviderService,
        },
        { provide: GrpcService, useValue: {} },
        ActiveUserService,
      ],
    });
    activeUserService = TestBed.inject(ActiveUserService);
  });

  test("#getActiveUser should return empty params for undefined request", () => {
    // Act
    const options = activeUserService.getActiveUser();

    // Assert
    expect(options.queryKey).toStrictEqual(["activeUser", undefined]);
    expect(options.queryFn).toBe(skipToken);
  });

  test("#getActiveUser should return correct params for defined request", () => {
    // Arrange
    const startingDateInclusive = new Date();
    startingDateInclusive.setDate(-10);
    startingDateInclusive.setHours(0, 0, 0, 0);
    const startingDateTimestamp = dateToTimestamp(startingDateInclusive);

    const endingDateInclusive = new Date();
    endingDateInclusive.setHours(0, 0, 0, 0);
    const endingDateTimestamp = dateToTimestamp(endingDateInclusive);

    // Act
    const options = activeUserService.getActiveUser({
      precision: "dailyPrecision",
      startingDateInclusive,
      endingDateInclusive,
    });

    // Assert
    expect(options.queryKey).toStrictEqual([
      "activeUser",
      {
        startingDateInclusive: {
          seconds: startingDateTimestamp.seconds.toString(),
          nanos: startingDateTimestamp.nanos,
        },
        endingDateInclusive: {
          seconds: endingDateTimestamp.seconds.toString(),
          nanos: endingDateTimestamp.nanos,
        },
        precision: {
          case: "dailyPrecision",
          value: {},
        },
      },
    ]);
  });
});
