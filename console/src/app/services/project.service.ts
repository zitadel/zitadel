import { inject, Injectable } from '@angular/core';
import { GrpcService } from './grpc.service';
import { MessageInitShape } from '@bufbuild/protobuf';
import {
  AddProjectRoleRequestSchema,
  CreateProjectRequestSchema,
} from '@zitadel/proto/zitadel/project/v2beta/project_service_pb';
import { mutationOptions } from '@tanstack/angular-query-experimental';
import { ManagementService } from './mgmt.service';
import { StorageKey, StorageLocation, StorageService } from './storage.service';

type CreateProjectRequest = Omit<
  Exclude<MessageInitShape<typeof CreateProjectRequestSchema>, { ['$typeName']: string }>,
  'organizationId'
>;

@Injectable({
  providedIn: 'root',
})
export class ProjectService {
  private readonly grpcService = inject(GrpcService);
  private readonly storageService = inject(StorageService);
  private readonly mgmtService = inject(ManagementService);

  private async createProject(request: CreateProjectRequest) {
    const organizationId = this.storageService.getItem(StorageKey.organizationId, StorageLocation.session) ?? undefined;
    const result = await this.grpcService.project.createProject({ ...request, organizationId });

    try {
      const { project } = await this.mgmtService.getProjectByID(result.id);
      if (!project) {
        return result;
      }
      this.mgmtService.ownedProjects.next([...this.mgmtService.ownedProjects.value, project]);
    } catch {
      //ignore
    }

    return result;
  }

  public createProjectMutationOptions = () =>
    mutationOptions({
      mutationKey: ['project', 'create'],
      mutationFn: (req: CreateProjectRequest) => this.createProject(req),
    });

  private addProjectRole(request: MessageInitShape<typeof AddProjectRoleRequestSchema>) {
    return this.grpcService.project.addProjectRole(request);
  }

  public addProjectRoleMutationsOptions = () =>
    mutationOptions({
      mutationKey: ['project', 'addRole'],
      mutationFn: (req: MessageInitShape<typeof AddProjectRoleRequestSchema>) => this.addProjectRole(req),
    });
}
