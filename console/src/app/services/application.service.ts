import { inject, Injectable } from '@angular/core';
import { GrpcService } from './grpc.service';
import { mutationOptions } from '@tanstack/angular-query-experimental';
import { MessageInitShape } from '@bufbuild/protobuf';
import { CreateApplicationRequestSchema, CreateApplicationResponse } from '@zitadel/proto/zitadel/app/v2beta/app_service_pb';

type CreationRequestTypeCase = Extract<
  MessageInitShape<typeof CreateApplicationRequestSchema>['creationRequestType'],
  { case: string }
>['case'];

type CreationResponseTypeCase<T extends CreationRequestTypeCase = CreationRequestTypeCase> = T extends `${infer C}Request`
  ? `${C}Response`
  : never;

// MessageInitShape creates a union type with one type where everything is required and one where everything is optional
// because the optional one is a superset of the required one we can safely exclude it
// this is needed because otherwise typescript gets confused somehow
type CreateApplicationRequest = Exclude<MessageInitShape<typeof CreateApplicationRequestSchema>, { $typeName: string }>;
type ExtractCreateApplicationRequest<T extends CreationRequestTypeCase = CreationRequestTypeCase> = Omit<
  CreateApplicationRequest,
  'creationRequestType'
> & {
  // for some reason typescript needs the double Extract first going to a string and then finally to T
  creationRequestType: Extract<Extract<CreateApplicationRequest['creationRequestType'], { case: string }>, { case: T }>;
};

type ExtractCreateApplicationResponse<T extends CreationResponseTypeCase = CreationResponseTypeCase> = Omit<
  CreateApplicationResponse,
  'creationResponseType'
> & {
  creationResponseType: Extract<Extract<CreateApplicationResponse['creationResponseType'], { case: string }>, { case: T }>;
};

@Injectable({
  providedIn: 'root',
})
export class ApplicationService {
  private readonly grpcService = inject(GrpcService);

  // this method is a generic version over the generated grpc createApplication method
  // which adds type safety for the applicationType oneof field
  private async createApplication<T extends CreationRequestTypeCase>(
    request: ExtractCreateApplicationRequest<T>,
  ): Promise<ExtractCreateApplicationResponse<CreationResponseTypeCase<T>>> {
    const { creationResponseType, ...response } = await this.grpcService.application.createApplication(request);

    if (!creationResponseType.case) {
      throw new Error('Application type is undefined in the response');
    }

    const responseWithCreationResponseType = { ...response, creationResponseType };
    if (responseIsSameCreationTypeAsRequest(responseWithCreationResponseType, request)) {
      return responseWithCreationResponseType;
    }

    throw new Error(
      `Mismatched application type in response. Expected: ${request.creationRequestType.case}, Received: ${creationResponseType.case}`,
    );
  }

  public createApplicationMutationOptions = <T extends CreationRequestTypeCase>() =>
    mutationOptions({
      mutationKey: ['createApplication'],
      mutationFn: (req: ExtractCreateApplicationRequest<T>) => this.createApplication(req),
    });
}

// typescript can't do narrowing on generics, so we need to use a type guard
// be careful when changing typeguards as they circumvent some of typescripts type safety
function responseIsSameCreationTypeAsRequest<T extends CreationRequestTypeCase>(
  response: ExtractCreateApplicationResponse,
  request: ExtractCreateApplicationRequest<T>,
): response is ExtractCreateApplicationResponse<CreationResponseTypeCase<T>> {
  response.creationResponseType;
  const caze = request.creationRequestType.case.replace('Request', '');
  return response.creationResponseType.case.includes(caze);
}
