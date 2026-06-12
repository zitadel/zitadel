import { Injectable } from '@angular/core';
import { MessageInitShape } from '@bufbuild/protobuf';
import {
  AddUsersToGroupRequestSchema,
  AddUsersToGroupResponse,
  CreateGroupRequestSchema,
  CreateGroupResponse,
  DeleteGroupRequestSchema,
  DeleteGroupResponse,
  GetGroupRequestSchema,
  GetGroupResponse,
  ListGroupsRequestSchema,
  ListGroupsResponse,
  ListGroupUsersRequestSchema,
  ListGroupUsersResponse,
  RemoveUsersFromGroupRequestSchema,
  RemoveUsersFromGroupResponse,
  UpdateGroupRequestSchema,
  UpdateGroupResponse,
} from '@zitadel/proto/zitadel/group/v2/group_service_pb';

import { GrpcService } from './grpc.service';

@Injectable({
  providedIn: 'root',
})
export class GroupService {
  constructor(private readonly grpcService: GrpcService) {}

  public getGroup(req: MessageInitShape<typeof GetGroupRequestSchema>): Promise<GetGroupResponse> {
    return this.grpcService.group.getGroup(req);
  }

  public listGroups(req: MessageInitShape<typeof ListGroupsRequestSchema>): Promise<ListGroupsResponse> {
    return this.grpcService.group.listGroups(req);
  }

  public createGroup(req: MessageInitShape<typeof CreateGroupRequestSchema>): Promise<CreateGroupResponse> {
    return this.grpcService.group.createGroup(req);
  }

  public updateGroup(req: MessageInitShape<typeof UpdateGroupRequestSchema>): Promise<UpdateGroupResponse> {
    return this.grpcService.group.updateGroup(req);
  }

  public deleteGroup(req: MessageInitShape<typeof DeleteGroupRequestSchema>): Promise<DeleteGroupResponse> {
    return this.grpcService.group.deleteGroup(req);
  }

  public listGroupUsers(req: MessageInitShape<typeof ListGroupUsersRequestSchema>): Promise<ListGroupUsersResponse> {
    return this.grpcService.group.listGroupUsers(req);
  }

  public addUsersToGroup(req: MessageInitShape<typeof AddUsersToGroupRequestSchema>): Promise<AddUsersToGroupResponse> {
    return this.grpcService.group.addUsersToGroup(req);
  }

  public removeUsersFromGroup(
    req: MessageInitShape<typeof RemoveUsersFromGroupRequestSchema>,
  ): Promise<RemoveUsersFromGroupResponse> {
    return this.grpcService.group.removeUsersFromGroup(req);
  }
}
