import { inject, Injectable } from '@angular/core';
import { TextQueryMethod } from 'src/app/proto/generated/zitadel/object_pb';
import {
  DisplayNameQuery,
  EmailQuery,
  SearchQuery as UserSearchQuery,
  User,
  UserNameQuery,
} from 'src/app/proto/generated/zitadel/user_pb';

import { ManagementService } from './mgmt.service';

export type SuggestionField = 'displayName' | 'userName' | 'email';

/** Kept small: the list is a hint while typing, not a browsable result set. */
const SUGGESTION_LIMIT = 10;

/**
 * Supplies existing field values for the list filters.
 *
 * Typing an exact match by hand is close to impossible with the "equals" method, so the
 * filter inputs offer what is actually stored while still accepting free text.
 */
@Injectable({
  providedIn: 'root',
})
export class UserSuggestionService {
  private readonly mgmtService = inject(ManagementService);

  public async suggest(field: SuggestionField, term: string): Promise<string[]> {
    const trimmed = term.trim();
    if (!trimmed) {
      return [];
    }

    const response = await this.mgmtService.listUsers(SUGGESTION_LIMIT, 0, [this.buildQuery(field, trimmed)]);

    const values = response.resultList
      .map((user) => UserSuggestionService.valueOf(field, user))
      .filter((value): value is string => !!value);

    // Several users can share an address, so the raw list would repeat entries.
    return Array.from(new Set(values));
  }

  private buildQuery(field: SuggestionField, term: string): UserSearchQuery {
    const query = new UserSearchQuery();
    const method = TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE;

    switch (field) {
      case 'displayName':
        const dnq = new DisplayNameQuery();
        dnq.setMethod(method);
        dnq.setDisplayName(term);
        query.setDisplayNameQuery(dnq);
        break;
      case 'email':
        const eq = new EmailQuery();
        eq.setMethod(method);
        eq.setEmailAddress(term);
        query.setEmailQuery(eq);
        break;
      case 'userName':
        const unq = new UserNameQuery();
        unq.setMethod(method);
        unq.setUserName(term);
        query.setUserNameQuery(unq);
        break;
    }

    return query;
  }

  private static valueOf(field: SuggestionField, user: User.AsObject): string | undefined {
    switch (field) {
      case 'displayName':
        // machines have no profile, their name is the closest equivalent
        return user.human?.profile?.displayName ?? user.machine?.name;
      case 'email':
        return user.human?.email?.email;
      case 'userName':
        return user.userName;
    }
  }
}
