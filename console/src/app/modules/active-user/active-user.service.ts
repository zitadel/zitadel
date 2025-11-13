import { Injectable } from '@angular/core';
import { create, MessageInitShape } from '@bufbuild/protobuf';
import {
  ActiveUserRequestSchema,
  ActiveUserRequest as ActiveUserRequestGrpc,
  ActiveUserEntry,
} from '@zitadel/proto/zitadel/analytics/v2beta/active_user_service_pb';
import { queryOptions, skipToken } from '@tanstack/angular-query-experimental';
import { ActiveUserGrpcProviderService } from './active-user-grpc-provider.service';
import { TimestampSchema } from '@bufbuild/protobuf/wkt';

export type PrecisionType = NonNullable<ActiveUserRequestGrpc['precision']['case']>;

type ActiveUserRequest = {
  precision: PrecisionType;
  startingDateInclusive: Date;
  endingDateInclusive: Date;
};

@Injectable()
export class ActiveUserService {
  private readonly activeUserService = this.activeUserGrpcProviderService.getClient();

  constructor(private readonly activeUserGrpcProviderService: ActiveUserGrpcProviderService) {}

  public getActiveUser(request?: ActiveUserRequest) {
    const req = request
      ? ({
          startingDateInclusive: dateToTimestamp(request.startingDateInclusive),
          endingDateInclusive: dateToTimestamp(request.endingDateInclusive),
          precision: {
            case: request.precision,
            value: {},
          },
        } satisfies MessageInitShape<typeof ActiveUserRequestSchema>)
      : undefined;

    // bigint cannot be used in query keys, so we convert to string
    const queryKey =
      req !== undefined
        ? {
            ...req,
            startingDateInclusive: {
              seconds: req.startingDateInclusive.seconds.toString(),
              nanos: req.startingDateInclusive.nanos,
            },
            endingDateInclusive: {
              seconds: req.endingDateInclusive.seconds.toString(),
              nanos: req.endingDateInclusive.nanos,
            },
          }
        : undefined;

    return queryOptions({
      queryKey: ['activeUser', queryKey],
      queryFn: req ? () => this.activeUserService.listActiveUsers(req).then((res) => res.entries) : skipToken,
    });
  }
}

export function averageActiveUserEntries(entries: ActiveUserEntry[]): bigint {
  // to avoid diving by zero
  if (entries.length < 1) {
    return BigInt(0);
  }

  const sum = entries.map((entry) => entry.value).reduce((acc, curr) => acc + curr, BigInt(0));
  return sum / BigInt(entries.length);
}

export function dateToTimestamp(date: Date) {
  const millis = date.getTime();
  const seconds = Math.floor(millis / 1000);
  const nanos = (millis % 1000) * 1000 * 1000;

  return create(TimestampSchema, { seconds: BigInt(seconds), nanos });
}
