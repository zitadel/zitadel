import { Injectable } from '@angular/core';
import { MessageInitShape } from '@bufbuild/protobuf';
import {
  AddUsersToGroupRequestSchema,
  AddUsersToGroupResponse,
  CreateGroupGrantRequestSchema,
  CreateGroupGrantResponse,
  CreateGroupRequestSchema,
  CreateGroupResponse,
  DeleteGroupGrantRequestSchema,
  DeleteGroupGrantResponse,
  DeleteGroupRequestSchema,
  DeleteGroupResponse,
  GetGroupRequestSchema,
  GetGroupResponse,
  ListGroupGrantsRequestSchema,
  ListGroupGrantsResponse,
  ListGroupsRequestSchema,
  ListGroupsResponse,
  ListGroupUsersRequestSchema,
  ListGroupUsersResponse,
  RemoveUsersFromGroupRequestSchema,
  RemoveUsersFromGroupResponse,
  UpdateGroupGrantRequestSchema,
  UpdateGroupGrantResponse,
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

  public createGroupGrant(req: MessageInitShape<typeof CreateGroupGrantRequestSchema>): Promise<CreateGroupGrantResponse> {
    return this.grpcService.group.createGroupGrant(req);
  }

  public updateGroupGrant(req: MessageInitShape<typeof UpdateGroupGrantRequestSchema>): Promise<UpdateGroupGrantResponse> {
    return this.grpcService.group.updateGroupGrant(req);
  }

  public deleteGroupGrant(req: MessageInitShape<typeof DeleteGroupGrantRequestSchema>): Promise<DeleteGroupGrantResponse> {
    return this.grpcService.group.deleteGroupGrant(req);
  }

  public listGroupGrants(req: MessageInitShape<typeof ListGroupGrantsRequestSchema>): Promise<ListGroupGrantsResponse> {
    return this.grpcService.group.listGroupGrants(req);
  }
}
