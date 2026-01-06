import { ActiveUserGrpcMockService } from "./active-user-grpc-mock.service";
import { dateToTimestamp } from "./active-user.service";
import moment from "moment";
import { TimestampToDatePipe } from "@/pipes/timestamp-to-date-pipe/timestamp-to-date.pipe";
import { expect } from "vitest";

describe("ActiveUserGrpcMockService", () => {
  let service: ActiveUserGrpcMockService;
  const timestampToDate = new TimestampToDatePipe();

  beforeEach(() => {
    service = new ActiveUserGrpcMockService();
  });

  it("#listActiveUsers daily precision should return random data in range", async () => {
    const { entries } = await service.listActiveUsers({
      precision: {
        case: "dailyPrecision",
        value: {},
      },
      startingDateInclusive: dateToTimestamp(new Date("2023-01-01")),
      endingDateInclusive: dateToTimestamp(new Date("2023-01-10")),
    });

    expect(entries.length).toBe(10);
    for (let i = 1; i < entries.length; i++) {
      const previousEntry = timestampToDate.transform(entries[i - 1].date!);
      const entry = timestampToDate.transform(entries[i].date!);

      const diff = moment(entry).diff(moment(previousEntry), "days");
      expect(diff).toBe(1);
    }
  });

  it("#listActiveUsers monthly precision should return random data in range", async () => {
    const { entries } = await service.listActiveUsers({
      precision: {
        case: "monthlyPrecision",
        value: {},
      },
      startingDateInclusive: dateToTimestamp(new Date("2023-01-01")),
      endingDateInclusive: dateToTimestamp(new Date("2023-10-01")),
    });

    expect(entries.length).toBe(10);
    for (let i = 1; i < entries.length; i++) {
      const previousEntry = timestampToDate.transform(entries[i - 1].date!);
      const entry = timestampToDate.transform(entries[i].date!);

      const diff = moment(entry).diff(moment(previousEntry), "months");
      expect(diff).toBe(1);
    }
  });

  it("#listActiveUsers should throw error on invalid request", async () => {
    await expect(service.listActiveUsers({})).rejects.toThrowError(
      new Error("Invalid request")
    );
    await expect(
      service.listActiveUsers({
        precision: {
          case: "dailyPrecision",
          value: {},
        },
      })
    ).rejects.toThrowError(new Error("Invalid request"));
    await expect(
      service.listActiveUsers({
        precision: {
          case: "dailyPrecision",
          value: {},
        },
        startingDateInclusive: {
          seconds: BigInt(9),
        },
      })
    ).rejects.toThrowError(new Error("Invalid request"));
    await expect(
      service.listActiveUsers({
        precision: {
          case: "dailyPrecision",
          value: {},
        },
        endingDateInclusive: {
          seconds: BigInt(9),
        },
      })
    ).rejects.toThrowError(new Error("Invalid request"));
  });
});
