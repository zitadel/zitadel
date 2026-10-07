import { Injectable } from '@angular/core';
import { QueryClient, queryOptions } from '@tanstack/angular-query-experimental';
import { GrpcService } from './grpc.service';
import { UserService } from './user.service';

@Injectable({
  providedIn: 'root',
})
export class NewSettingsService {
  constructor(
    private readonly grpcService: GrpcService,
    private readonly userService: UserService,
    private readonly queryClient: QueryClient,
  ) {}

  // The general settings of the instance, readable by organization admins as well.
  public getGeneralSettingsQueryOptions() {
    return queryOptions({
      queryKey: [this.userService.userId(), 'SettingsService', 'generalSettings'],
      queryFn: ({ signal }) => this.grpcService.settingsNew.getGeneralSettings({}, { signal }),
    });
  }

  // Reloads the general settings, e.g. after the active email provider changed.
  public invalidateGeneralSettings() {
    return this.queryClient.invalidateQueries({
      queryKey: [this.userService.userId(), 'SettingsService', 'generalSettings'],
    });
  }
}
