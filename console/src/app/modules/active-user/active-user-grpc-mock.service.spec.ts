import { ActiveUserGrpcMockService } from './active-user-grpc-mock.service';
import { dateToTimestamp } from './active-user.service';
import moment from 'moment';
import { TimestampToDatePipe } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date.pipe';

describe('ActiveUserGrpcMockService', () => {
  let service: ActiveUserGrpcMockService;
  let timestampToDate = new TimestampToDatePipe();

  beforeEach(() => {
    service = new ActiveUserGrpcMockService();
  });

  it('#listActiveUsers daily precision should return random data in range', async () => {
    const { entries } = await service.listActiveUsers({
      precision: {
        case: 'dailyPrecision',
        value: {},
      },
      startingDateInclusive: dateToTimestamp(new Date('2023-01-01')),
      endingDateInclusive: dateToTimestamp(new Date('2023-01-10')),
    });

    await expect(entries.length).toBe(10);
    for (let i = 1; i < entries.length; i++) {
      const previousEntry = timestampToDate.transform(entries[i - 1].date!);
      const entry = timestampToDate.transform(entries[i].date!);

      const diff = moment(entry).diff(moment(previousEntry), 'days');
      await expect(diff).toBe(1);
    }
  });

  it('#listActiveUsers monthly precision should return random data in range', async () => {
    const { entries } = await service.listActiveUsers({
      precision: {
        case: 'monthlyPrecision',
        value: {},
      },
      startingDateInclusive: dateToTimestamp(new Date('2023-01-01')),
      endingDateInclusive: dateToTimestamp(new Date('2023-10-01')),
    });

    await expect(entries.length).toBe(10);
    for (let i = 1; i < entries.length; i++) {
      const previousEntry = timestampToDate.transform(entries[i - 1].date!);
      const entry = timestampToDate.transform(entries[i].date!);

      const diff = moment(entry).diff(moment(previousEntry), 'months');
      await expect(diff).toBe(1);
    }
  });

  it('#listActiveUsers should throw error on invalid request', async () => {
    await expectAsync(service.listActiveUsers({})).toBeRejectedWith(new Error('Invalid request'));
    await expectAsync(
      service.listActiveUsers({
        precision: {
          case: 'dailyPrecision',
          value: {},
        },
      }),
    ).toBeRejectedWith(new Error('Invalid request'));
    await expectAsync(
      service.listActiveUsers({
        precision: {
          case: 'dailyPrecision',
          value: {},
        },
        startingDateInclusive: {
          seconds: BigInt(9),
        },
      }),
    ).toBeRejectedWith(new Error('Invalid request'));
    await expectAsync(
      service.listActiveUsers({
        precision: {
          case: 'dailyPrecision',
          value: {},
        },
        endingDateInclusive: {
          seconds: BigInt(9),
        },
      }),
    ).toBeRejectedWith(new Error('Invalid request'));
  });
});
