import { Client } from '@connectrpc/connect';
import {
  ActiveUserEntry,
  ActiveUserEntrySchema,
  ActiveUserRequestSchema,
  ActiveUserResponse,
  ActiveUserResponseSchema,
  ActiveUserService as ActiveUserServiceGrpc,
} from '@zitadel/proto/zitadel/analytics/v2beta/active_user_service_pb';
import { create, MessageInitShape } from '@bufbuild/protobuf';
import { dateToTimestamp, PrecisionType } from './active-user.service';
import { TimestampToDatePipe } from '@/pipes/timestamp-to-date-pipe/timestamp-to-date.pipe';
import { TimestampSchema, Timestamp } from '@bufbuild/protobuf/wkt';

export class ActiveUserGrpcMockService implements Client<typeof ActiveUserServiceGrpc> {
  private isTimestamp(timestamp: MessageInitShape<typeof TimestampSchema> | undefined): timestamp is Timestamp {
    return timestamp?.seconds !== undefined && timestamp?.nanos !== undefined && timestamp.$typeName !== undefined;
  }

  public async listActiveUsers(request: MessageInitShape<typeof ActiveUserRequestSchema>): Promise<ActiveUserResponse> {
    const { precision, startingDateInclusive, endingDateInclusive } = request;

    if (!precision?.case || !this.isTimestamp(startingDateInclusive) || !this.isTimestamp(endingDateInclusive)) {
      throw new Error('Invalid request');
    }

    const timestampToDatePipe = new TimestampToDatePipe();

    return create(ActiveUserResponseSchema, {
      entries: this.generateFakeEntries(
        timestampToDatePipe.transform(startingDateInclusive),
        timestampToDatePipe.transform(endingDateInclusive),
        precision.case,
      ),
    });
  }

  private generateFakeEntries(start: Date, end: Date, precision: PrecisionType): ActiveUserEntry[] {
    const entries: ActiveUserEntry[] = [];

    // Value between 100 and 200
    let baseValue = 100 + Math.random() * 100; // base value

    const increase = (date: Date) =>
      precision === 'dailyPrecision' ? date.setDate(date.getDate() + 1) : date.setMonth(date.getMonth() + 1);

    for (let i = new Date(start); i <= end; increase(i)) {
      // Increase or decrease base value by a Maximum of 10
      const change = (Math.random() - 0.5) * 20;
      // Make sure value doesn't go below 0
      const newValue = Math.floor(Math.max(0, baseValue + change));

      entries.push(
        create(ActiveUserEntrySchema, {
          date: dateToTimestamp(i),
          value: BigInt(newValue),
        }),
      );
    }

    return entries;
  }
}
