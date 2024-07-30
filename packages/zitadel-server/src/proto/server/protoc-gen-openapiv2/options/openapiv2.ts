/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { Value } from "../../google/protobuf/struct";

export const protobufPackage = "grpc.gateway.protoc_gen_openapiv2.options";

/**
 * Scheme describes the schemes supported by the OpenAPI Swagger
 * and Operation objects.
 */
export enum Scheme {
  UNKNOWN = 0,
  HTTP = 1,
  HTTPS = 2,
  WS = 3,
  WSS = 4,
  UNRECOGNIZED = -1,
}

export function schemeFromJSON(object: any): Scheme {
  switch (object) {
    case 0:
    case "UNKNOWN":
      return Scheme.UNKNOWN;
    case 1:
    case "HTTP":
      return Scheme.HTTP;
    case 2:
    case "HTTPS":
      return Scheme.HTTPS;
    case 3:
    case "WS":
      return Scheme.WS;
    case 4:
    case "WSS":
      return Scheme.WSS;
    case -1:
    case "UNRECOGNIZED":
    default:
      return Scheme.UNRECOGNIZED;
  }
}

export function schemeToJSON(object: Scheme): string {
  switch (object) {
    case Scheme.UNKNOWN:
      return "UNKNOWN";
    case Scheme.HTTP:
      return "HTTP";
    case Scheme.HTTPS:
      return "HTTPS";
    case Scheme.WS:
      return "WS";
    case Scheme.WSS:
      return "WSS";
    case Scheme.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

/**
 * `Swagger` is a representation of OpenAPI v2 specification's Swagger object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#swaggerObject
 *
 * Example:
 *
 *  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
 *    info: {
 *      title: "Echo API";
 *      version: "1.0";
 *      description: "";
 *      contact: {
 *        name: "gRPC-Gateway project";
 *        url: "https://github.com/grpc-ecosystem/grpc-gateway";
 *        email: "none@example.com";
 *      };
 *      license: {
 *        name: "BSD 3-Clause License";
 *        url: "https://github.com/grpc-ecosystem/grpc-gateway/blob/main/LICENSE.txt";
 *      };
 *    };
 *    schemes: HTTPS;
 *    consumes: "application/json";
 *    produces: "application/json";
 *  };
 */
export interface Swagger {
  /**
   * Specifies the OpenAPI Specification version being used. It can be
   * used by the OpenAPI UI and other clients to interpret the API listing. The
   * value MUST be "2.0".
   */
  swagger: string;
  /**
   * Provides metadata about the API. The metadata can be used by the
   * clients if needed.
   */
  info:
    | Info
    | undefined;
  /**
   * The host (name or ip) serving the API. This MUST be the host only and does
   * not include the scheme nor sub-paths. It MAY include a port. If the host is
   * not included, the host serving the documentation is to be used (including
   * the port). The host does not support path templating.
   */
  host: string;
  /**
   * The base path on which the API is served, which is relative to the host. If
   * it is not included, the API is served directly under the host. The value
   * MUST start with a leading slash (/). The basePath does not support path
   * templating.
   * Note that using `base_path` does not change the endpoint paths that are
   * generated in the resulting OpenAPI file. If you wish to use `base_path`
   * with relatively generated OpenAPI paths, the `base_path` prefix must be
   * manually removed from your `google.api.http` paths and your code changed to
   * serve the API from the `base_path`.
   */
  basePath: string;
  /**
   * The transfer protocol of the API. Values MUST be from the list: "http",
   * "https", "ws", "wss". If the schemes is not included, the default scheme to
   * be used is the one used to access the OpenAPI definition itself.
   */
  schemes: Scheme[];
  /**
   * A list of MIME types the APIs can consume. This is global to all APIs but
   * can be overridden on specific API calls. Value MUST be as described under
   * Mime Types.
   */
  consumes: string[];
  /**
   * A list of MIME types the APIs can produce. This is global to all APIs but
   * can be overridden on specific API calls. Value MUST be as described under
   * Mime Types.
   */
  produces: string[];
  /**
   * An object to hold responses that can be used across operations. This
   * property does not define global responses for all operations.
   */
  responses: { [key: string]: Response };
  /** Security scheme definitions that can be used across the specification. */
  securityDefinitions:
    | SecurityDefinitions
    | undefined;
  /**
   * A declaration of which security schemes are applied for the API as a whole.
   * The list of values describes alternative security schemes that can be used
   * (that is, there is a logical OR between the security requirements).
   * Individual operations can override this definition.
   */
  security: SecurityRequirement[];
  /**
   * A list of tags for API documentation control. Tags can be used for logical
   * grouping of operations by resources or any other qualifier.
   */
  tags: Tag[];
  /** Additional external documentation. */
  externalDocs:
    | ExternalDocumentation
    | undefined;
  /**
   * Custom properties that start with "x-" such as "x-foo" used to describe
   * extra functionality that is not covered by the standard OpenAPI Specification.
   * See: https://swagger.io/docs/specification/2-0/swagger-extensions/
   */
  extensions: { [key: string]: any | undefined };
}

export interface Swagger_ResponsesEntry {
  key: string;
  value: Response | undefined;
}

export interface Swagger_ExtensionsEntry {
  key: string;
  value: any | undefined;
}

/**
 * `Operation` is a representation of OpenAPI v2 specification's Operation object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#operationObject
 *
 * Example:
 *
 *  service EchoService {
 *    rpc Echo(SimpleMessage) returns (SimpleMessage) {
 *      option (google.api.http) = {
 *        get: "/v1/example/echo/{id}"
 *      };
 *
 *      option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
 *        summary: "Get a message.";
 *        operation_id: "getMessage";
 *        tags: "echo";
 *        responses: {
 *          key: "200"
 *            value: {
 *            description: "OK";
 *          }
 *        }
 *      };
 *    }
 *  }
 */
export interface Operation {
  /**
   * A list of tags for API documentation control. Tags can be used for logical
   * grouping of operations by resources or any other qualifier.
   */
  tags: string[];
  /**
   * A short summary of what the operation does. For maximum readability in the
   * swagger-ui, this field SHOULD be less than 120 characters.
   */
  summary: string;
  /**
   * A verbose explanation of the operation behavior. GFM syntax can be used for
   * rich text representation.
   */
  description: string;
  /** Additional external documentation for this operation. */
  externalDocs:
    | ExternalDocumentation
    | undefined;
  /**
   * Unique string used to identify the operation. The id MUST be unique among
   * all operations described in the API. Tools and libraries MAY use the
   * operationId to uniquely identify an operation, therefore, it is recommended
   * to follow common programming naming conventions.
   */
  operationId: string;
  /**
   * A list of MIME types the operation can consume. This overrides the consumes
   * definition at the OpenAPI Object. An empty value MAY be used to clear the
   * global definition. Value MUST be as described under Mime Types.
   */
  consumes: string[];
  /**
   * A list of MIME types the operation can produce. This overrides the produces
   * definition at the OpenAPI Object. An empty value MAY be used to clear the
   * global definition. Value MUST be as described under Mime Types.
   */
  produces: string[];
  /**
   * The list of possible responses as they are returned from executing this
   * operation.
   */
  responses: { [key: string]: Response };
  /**
   * The transfer protocol for the operation. Values MUST be from the list:
   * "http", "https", "ws", "wss". The value overrides the OpenAPI Object
   * schemes definition.
   */
  schemes: Scheme[];
  /**
   * Declares this operation to be deprecated. Usage of the declared operation
   * should be refrained. Default value is false.
   */
  deprecated: boolean;
  /**
   * A declaration of which security schemes are applied for this operation. The
   * list of values describes alternative security schemes that can be used
   * (that is, there is a logical OR between the security requirements). This
   * definition overrides any declared top-level security. To remove a top-level
   * security declaration, an empty array can be used.
   */
  security: SecurityRequirement[];
  /**
   * Custom properties that start with "x-" such as "x-foo" used to describe
   * extra functionality that is not covered by the standard OpenAPI Specification.
   * See: https://swagger.io/docs/specification/2-0/swagger-extensions/
   */
  extensions: { [key: string]: any | undefined };
  /**
   * Custom parameters such as HTTP request headers.
   * See: https://swagger.io/docs/specification/2-0/describing-parameters/
   * and https://swagger.io/specification/v2/#parameter-object.
   */
  parameters: Parameters | undefined;
}

export interface Operation_ResponsesEntry {
  key: string;
  value: Response | undefined;
}

export interface Operation_ExtensionsEntry {
  key: string;
  value: any | undefined;
}

/**
 * `Parameters` is a representation of OpenAPI v2 specification's parameters object.
 * Note: This technically breaks compatibility with the OpenAPI 2 definition structure as we only
 * allow header parameters to be set here since we do not want users specifying custom non-header
 * parameters beyond those inferred from the Protobuf schema.
 * See: https://swagger.io/specification/v2/#parameter-object
 */
export interface Parameters {
  /**
   * `Headers` is one or more HTTP header parameter.
   * See: https://swagger.io/docs/specification/2-0/describing-parameters/#header-parameters
   */
  headers: HeaderParameter[];
}

/**
 * `HeaderParameter` a HTTP header parameter.
 * See: https://swagger.io/specification/v2/#parameter-object
 */
export interface HeaderParameter {
  /** `Name` is the header name. */
  name: string;
  /** `Description` is a short description of the header. */
  description: string;
  /**
   * `Type` is the type of the object. The value MUST be one of "string", "number", "integer", or "boolean". The "array" type is not supported.
   * See: https://swagger.io/specification/v2/#parameterType.
   */
  type: HeaderParameter_Type;
  /** `Format` The extending format for the previously mentioned type. */
  format: string;
  /** `Required` indicates if the header is optional */
  required: boolean;
}

/**
 * `Type` is a a supported HTTP header type.
 * See https://swagger.io/specification/v2/#parameterType.
 */
export enum HeaderParameter_Type {
  UNKNOWN = 0,
  STRING = 1,
  NUMBER = 2,
  INTEGER = 3,
  BOOLEAN = 4,
  UNRECOGNIZED = -1,
}

export function headerParameter_TypeFromJSON(object: any): HeaderParameter_Type {
  switch (object) {
    case 0:
    case "UNKNOWN":
      return HeaderParameter_Type.UNKNOWN;
    case 1:
    case "STRING":
      return HeaderParameter_Type.STRING;
    case 2:
    case "NUMBER":
      return HeaderParameter_Type.NUMBER;
    case 3:
    case "INTEGER":
      return HeaderParameter_Type.INTEGER;
    case 4:
    case "BOOLEAN":
      return HeaderParameter_Type.BOOLEAN;
    case -1:
    case "UNRECOGNIZED":
    default:
      return HeaderParameter_Type.UNRECOGNIZED;
  }
}

export function headerParameter_TypeToJSON(object: HeaderParameter_Type): string {
  switch (object) {
    case HeaderParameter_Type.UNKNOWN:
      return "UNKNOWN";
    case HeaderParameter_Type.STRING:
      return "STRING";
    case HeaderParameter_Type.NUMBER:
      return "NUMBER";
    case HeaderParameter_Type.INTEGER:
      return "INTEGER";
    case HeaderParameter_Type.BOOLEAN:
      return "BOOLEAN";
    case HeaderParameter_Type.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

/**
 * `Header` is a representation of OpenAPI v2 specification's Header object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#headerObject
 */
export interface Header {
  /** `Description` is a short description of the header. */
  description: string;
  /** The type of the object. The value MUST be one of "string", "number", "integer", or "boolean". The "array" type is not supported. */
  type: string;
  /** `Format` The extending format for the previously mentioned type. */
  format: string;
  /**
   * `Default` Declares the value of the header that the server will use if none is provided.
   * See: https://tools.ietf.org/html/draft-fge-json-schema-validation-00#section-6.2.
   * Unlike JSON Schema this value MUST conform to the defined type for the header.
   */
  default: string;
  /** 'Pattern' See https://tools.ietf.org/html/draft-fge-json-schema-validation-00#section-5.2.3. */
  pattern: string;
}

/**
 * `Response` is a representation of OpenAPI v2 specification's Response object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#responseObject
 */
export interface Response {
  /**
   * `Description` is a short description of the response.
   * GFM syntax can be used for rich text representation.
   */
  description: string;
  /**
   * `Schema` optionally defines the structure of the response.
   * If `Schema` is not provided, it means there is no content to the response.
   */
  schema:
    | Schema
    | undefined;
  /**
   * `Headers` A list of headers that are sent with the response.
   * `Header` name is expected to be a string in the canonical format of the MIME header key
   * See: https://golang.org/pkg/net/textproto/#CanonicalMIMEHeaderKey
   */
  headers: { [key: string]: Header };
  /**
   * `Examples` gives per-mimetype response examples.
   * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#example-object
   */
  examples: { [key: string]: string };
  /**
   * Custom properties that start with "x-" such as "x-foo" used to describe
   * extra functionality that is not covered by the standard OpenAPI Specification.
   * See: https://swagger.io/docs/specification/2-0/swagger-extensions/
   */
  extensions: { [key: string]: any | undefined };
}

export interface Response_HeadersEntry {
  key: string;
  value: Header | undefined;
}

export interface Response_ExamplesEntry {
  key: string;
  value: string;
}

export interface Response_ExtensionsEntry {
  key: string;
  value: any | undefined;
}

/**
 * `Info` is a representation of OpenAPI v2 specification's Info object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#infoObject
 *
 * Example:
 *
 *  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
 *    info: {
 *      title: "Echo API";
 *      version: "1.0";
 *      description: "";
 *      contact: {
 *        name: "gRPC-Gateway project";
 *        url: "https://github.com/grpc-ecosystem/grpc-gateway";
 *        email: "none@example.com";
 *      };
 *      license: {
 *        name: "BSD 3-Clause License";
 *        url: "https://github.com/grpc-ecosystem/grpc-gateway/blob/main/LICENSE.txt";
 *      };
 *    };
 *    ...
 *  };
 */
export interface Info {
  /** The title of the application. */
  title: string;
  /**
   * A short description of the application. GFM syntax can be used for rich
   * text representation.
   */
  description: string;
  /** The Terms of Service for the API. */
  termsOfService: string;
  /** The contact information for the exposed API. */
  contact:
    | Contact
    | undefined;
  /** The license information for the exposed API. */
  license:
    | License
    | undefined;
  /**
   * Provides the version of the application API (not to be confused
   * with the specification version).
   */
  version: string;
  /**
   * Custom properties that start with "x-" such as "x-foo" used to describe
   * extra functionality that is not covered by the standard OpenAPI Specification.
   * See: https://swagger.io/docs/specification/2-0/swagger-extensions/
   */
  extensions: { [key: string]: any | undefined };
}

export interface Info_ExtensionsEntry {
  key: string;
  value: any | undefined;
}

/**
 * `Contact` is a representation of OpenAPI v2 specification's Contact object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#contactObject
 *
 * Example:
 *
 *  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
 *    info: {
 *      ...
 *      contact: {
 *        name: "gRPC-Gateway project";
 *        url: "https://github.com/grpc-ecosystem/grpc-gateway";
 *        email: "none@example.com";
 *      };
 *      ...
 *    };
 *    ...
 *  };
 */
export interface Contact {
  /** The identifying name of the contact person/organization. */
  name: string;
  /**
   * The URL pointing to the contact information. MUST be in the format of a
   * URL.
   */
  url: string;
  /**
   * The email address of the contact person/organization. MUST be in the format
   * of an email address.
   */
  email: string;
}

/**
 * `License` is a representation of OpenAPI v2 specification's License object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#licenseObject
 *
 * Example:
 *
 *  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
 *    info: {
 *      ...
 *      license: {
 *        name: "BSD 3-Clause License";
 *        url: "https://github.com/grpc-ecosystem/grpc-gateway/blob/main/LICENSE.txt";
 *      };
 *      ...
 *    };
 *    ...
 *  };
 */
export interface License {
  /** The license name used for the API. */
  name: string;
  /** A URL to the license used for the API. MUST be in the format of a URL. */
  url: string;
}

/**
 * `ExternalDocumentation` is a representation of OpenAPI v2 specification's
 * ExternalDocumentation object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#externalDocumentationObject
 *
 * Example:
 *
 *  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
 *    ...
 *    external_docs: {
 *      description: "More about gRPC-Gateway";
 *      url: "https://github.com/grpc-ecosystem/grpc-gateway";
 *    }
 *    ...
 *  };
 */
export interface ExternalDocumentation {
  /**
   * A short description of the target documentation. GFM syntax can be used for
   * rich text representation.
   */
  description: string;
  /**
   * The URL for the target documentation. Value MUST be in the format
   * of a URL.
   */
  url: string;
}

/**
 * `Schema` is a representation of OpenAPI v2 specification's Schema object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#schemaObject
 */
export interface Schema {
  jsonSchema:
    | JSONSchema
    | undefined;
  /**
   * Adds support for polymorphism. The discriminator is the schema property
   * name that is used to differentiate between other schema that inherit this
   * schema. The property name used MUST be defined at this schema and it MUST
   * be in the required property list. When used, the value MUST be the name of
   * this schema or any schema that inherits it.
   */
  discriminator: string;
  /**
   * Relevant only for Schema "properties" definitions. Declares the property as
   * "read only". This means that it MAY be sent as part of a response but MUST
   * NOT be sent as part of the request. Properties marked as readOnly being
   * true SHOULD NOT be in the required list of the defined schema. Default
   * value is false.
   */
  readOnly: boolean;
  /** Additional external documentation for this schema. */
  externalDocs:
    | ExternalDocumentation
    | undefined;
  /**
   * A free-form property to include an example of an instance for this schema in JSON.
   * This is copied verbatim to the output.
   */
  example: string;
}

/**
 * `JSONSchema` represents properties from JSON Schema taken, and as used, in
 * the OpenAPI v2 spec.
 *
 * This includes changes made by OpenAPI v2.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#schemaObject
 *
 * See also: https://cswr.github.io/JsonSchema/spec/basic_types/,
 * https://github.com/json-schema-org/json-schema-spec/blob/master/schema.json
 *
 * Example:
 *
 *  message SimpleMessage {
 *    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
 *      json_schema: {
 *        title: "SimpleMessage"
 *        description: "A simple message."
 *        required: ["id"]
 *      }
 *    };
 *
 *    // Id represents the message identifier.
 *    string id = 1; [
 *        (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = {
 *          description: "The unique identifier of the simple message."
 *        }];
 *  }
 */
export interface JSONSchema {
  /**
   * Ref is used to define an external reference to include in the message.
   * This could be a fully qualified proto message reference, and that type must
   * be imported into the protofile. If no message is identified, the Ref will
   * be used verbatim in the output.
   * For example:
   *  `ref: ".google.protobuf.Timestamp"`.
   */
  ref: string;
  /** The title of the schema. */
  title: string;
  /** A short description of the schema. */
  description: string;
  default: string;
  readOnly: boolean;
  /**
   * A free-form property to include a JSON example of this field. This is copied
   * verbatim to the output swagger.json. Quotes must be escaped.
   * This property is the same for 2.0 and 3.0.0 https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/3.0.0.md#schemaObject  https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#schemaObject
   */
  example: string;
  multipleOf: number;
  /**
   * Maximum represents an inclusive upper limit for a numeric instance. The
   * value of MUST be a number,
   */
  maximum: number;
  exclusiveMaximum: boolean;
  /**
   * minimum represents an inclusive lower limit for a numeric instance. The
   * value of MUST be a number,
   */
  minimum: number;
  exclusiveMinimum: boolean;
  maxLength: number;
  minLength: number;
  pattern: string;
  maxItems: number;
  minItems: number;
  uniqueItems: boolean;
  maxProperties: number;
  minProperties: number;
  required: string[];
  /** Items in 'array' must be unique. */
  array: string[];
  type: JSONSchema_JSONSchemaSimpleTypes[];
  /** `Format` */
  format: string;
  /** Items in `enum` must be unique https://tools.ietf.org/html/draft-fge-json-schema-validation-00#section-5.5.1 */
  enum: string[];
  /** Additional field level properties used when generating the OpenAPI v2 file. */
  fieldConfiguration:
    | JSONSchema_FieldConfiguration
    | undefined;
  /**
   * Custom properties that start with "x-" such as "x-foo" used to describe
   * extra functionality that is not covered by the standard OpenAPI Specification.
   * See: https://swagger.io/docs/specification/2-0/swagger-extensions/
   */
  extensions: { [key: string]: any | undefined };
}

export enum JSONSchema_JSONSchemaSimpleTypes {
  UNKNOWN = 0,
  ARRAY = 1,
  BOOLEAN = 2,
  INTEGER = 3,
  NULL = 4,
  NUMBER = 5,
  OBJECT = 6,
  STRING = 7,
  UNRECOGNIZED = -1,
}

export function jSONSchema_JSONSchemaSimpleTypesFromJSON(object: any): JSONSchema_JSONSchemaSimpleTypes {
  switch (object) {
    case 0:
    case "UNKNOWN":
      return JSONSchema_JSONSchemaSimpleTypes.UNKNOWN;
    case 1:
    case "ARRAY":
      return JSONSchema_JSONSchemaSimpleTypes.ARRAY;
    case 2:
    case "BOOLEAN":
      return JSONSchema_JSONSchemaSimpleTypes.BOOLEAN;
    case 3:
    case "INTEGER":
      return JSONSchema_JSONSchemaSimpleTypes.INTEGER;
    case 4:
    case "NULL":
      return JSONSchema_JSONSchemaSimpleTypes.NULL;
    case 5:
    case "NUMBER":
      return JSONSchema_JSONSchemaSimpleTypes.NUMBER;
    case 6:
    case "OBJECT":
      return JSONSchema_JSONSchemaSimpleTypes.OBJECT;
    case 7:
    case "STRING":
      return JSONSchema_JSONSchemaSimpleTypes.STRING;
    case -1:
    case "UNRECOGNIZED":
    default:
      return JSONSchema_JSONSchemaSimpleTypes.UNRECOGNIZED;
  }
}

export function jSONSchema_JSONSchemaSimpleTypesToJSON(object: JSONSchema_JSONSchemaSimpleTypes): string {
  switch (object) {
    case JSONSchema_JSONSchemaSimpleTypes.UNKNOWN:
      return "UNKNOWN";
    case JSONSchema_JSONSchemaSimpleTypes.ARRAY:
      return "ARRAY";
    case JSONSchema_JSONSchemaSimpleTypes.BOOLEAN:
      return "BOOLEAN";
    case JSONSchema_JSONSchemaSimpleTypes.INTEGER:
      return "INTEGER";
    case JSONSchema_JSONSchemaSimpleTypes.NULL:
      return "NULL";
    case JSONSchema_JSONSchemaSimpleTypes.NUMBER:
      return "NUMBER";
    case JSONSchema_JSONSchemaSimpleTypes.OBJECT:
      return "OBJECT";
    case JSONSchema_JSONSchemaSimpleTypes.STRING:
      return "STRING";
    case JSONSchema_JSONSchemaSimpleTypes.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

/**
 * 'FieldConfiguration' provides additional field level properties used when generating the OpenAPI v2 file.
 * These properties are not defined by OpenAPIv2, but they are used to control the generation.
 */
export interface JSONSchema_FieldConfiguration {
  /**
   * Alternative parameter name when used as path parameter. If set, this will
   * be used as the complete parameter name when this field is used as a path
   * parameter. Use this to avoid having auto generated path parameter names
   * for overlapping paths.
   */
  pathParamName: string;
}

export interface JSONSchema_ExtensionsEntry {
  key: string;
  value: any | undefined;
}

/**
 * `Tag` is a representation of OpenAPI v2 specification's Tag object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#tagObject
 */
export interface Tag {
  /**
   * The name of the tag. Use it to allow override of the name of a
   * global Tag object, then use that name to reference the tag throughout the
   * OpenAPI file.
   */
  name: string;
  /**
   * A short description for the tag. GFM syntax can be used for rich text
   * representation.
   */
  description: string;
  /** Additional external documentation for this tag. */
  externalDocs:
    | ExternalDocumentation
    | undefined;
  /**
   * Custom properties that start with "x-" such as "x-foo" used to describe
   * extra functionality that is not covered by the standard OpenAPI Specification.
   * See: https://swagger.io/docs/specification/2-0/swagger-extensions/
   */
  extensions: { [key: string]: any | undefined };
}

export interface Tag_ExtensionsEntry {
  key: string;
  value: any | undefined;
}

/**
 * `SecurityDefinitions` is a representation of OpenAPI v2 specification's
 * Security Definitions object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#securityDefinitionsObject
 *
 * A declaration of the security schemes available to be used in the
 * specification. This does not enforce the security schemes on the operations
 * and only serves to provide the relevant details for each scheme.
 */
export interface SecurityDefinitions {
  /**
   * A single security scheme definition, mapping a "name" to the scheme it
   * defines.
   */
  security: { [key: string]: SecurityScheme };
}

export interface SecurityDefinitions_SecurityEntry {
  key: string;
  value: SecurityScheme | undefined;
}

/**
 * `SecurityScheme` is a representation of OpenAPI v2 specification's
 * Security Scheme object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#securitySchemeObject
 *
 * Allows the definition of a security scheme that can be used by the
 * operations. Supported schemes are basic authentication, an API key (either as
 * a header or as a query parameter) and OAuth2's common flows (implicit,
 * password, application and access code).
 */
export interface SecurityScheme {
  /**
   * The type of the security scheme. Valid values are "basic",
   * "apiKey" or "oauth2".
   */
  type: SecurityScheme_Type;
  /** A short description for security scheme. */
  description: string;
  /**
   * The name of the header or query parameter to be used.
   * Valid for apiKey.
   */
  name: string;
  /**
   * The location of the API key. Valid values are "query" or
   * "header".
   * Valid for apiKey.
   */
  in: SecurityScheme_In;
  /**
   * The flow used by the OAuth2 security scheme. Valid values are
   * "implicit", "password", "application" or "accessCode".
   * Valid for oauth2.
   */
  flow: SecurityScheme_Flow;
  /**
   * The authorization URL to be used for this flow. This SHOULD be in
   * the form of a URL.
   * Valid for oauth2/implicit and oauth2/accessCode.
   */
  authorizationUrl: string;
  /**
   * The token URL to be used for this flow. This SHOULD be in the
   * form of a URL.
   * Valid for oauth2/password, oauth2/application and oauth2/accessCode.
   */
  tokenUrl: string;
  /**
   * The available scopes for the OAuth2 security scheme.
   * Valid for oauth2.
   */
  scopes:
    | Scopes
    | undefined;
  /**
   * Custom properties that start with "x-" such as "x-foo" used to describe
   * extra functionality that is not covered by the standard OpenAPI Specification.
   * See: https://swagger.io/docs/specification/2-0/swagger-extensions/
   */
  extensions: { [key: string]: any | undefined };
}

/**
 * The type of the security scheme. Valid values are "basic",
 * "apiKey" or "oauth2".
 */
export enum SecurityScheme_Type {
  TYPE_INVALID = 0,
  TYPE_BASIC = 1,
  TYPE_API_KEY = 2,
  TYPE_OAUTH2 = 3,
  UNRECOGNIZED = -1,
}

export function securityScheme_TypeFromJSON(object: any): SecurityScheme_Type {
  switch (object) {
    case 0:
    case "TYPE_INVALID":
      return SecurityScheme_Type.TYPE_INVALID;
    case 1:
    case "TYPE_BASIC":
      return SecurityScheme_Type.TYPE_BASIC;
    case 2:
    case "TYPE_API_KEY":
      return SecurityScheme_Type.TYPE_API_KEY;
    case 3:
    case "TYPE_OAUTH2":
      return SecurityScheme_Type.TYPE_OAUTH2;
    case -1:
    case "UNRECOGNIZED":
    default:
      return SecurityScheme_Type.UNRECOGNIZED;
  }
}

export function securityScheme_TypeToJSON(object: SecurityScheme_Type): string {
  switch (object) {
    case SecurityScheme_Type.TYPE_INVALID:
      return "TYPE_INVALID";
    case SecurityScheme_Type.TYPE_BASIC:
      return "TYPE_BASIC";
    case SecurityScheme_Type.TYPE_API_KEY:
      return "TYPE_API_KEY";
    case SecurityScheme_Type.TYPE_OAUTH2:
      return "TYPE_OAUTH2";
    case SecurityScheme_Type.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

/** The location of the API key. Valid values are "query" or "header". */
export enum SecurityScheme_In {
  IN_INVALID = 0,
  IN_QUERY = 1,
  IN_HEADER = 2,
  UNRECOGNIZED = -1,
}

export function securityScheme_InFromJSON(object: any): SecurityScheme_In {
  switch (object) {
    case 0:
    case "IN_INVALID":
      return SecurityScheme_In.IN_INVALID;
    case 1:
    case "IN_QUERY":
      return SecurityScheme_In.IN_QUERY;
    case 2:
    case "IN_HEADER":
      return SecurityScheme_In.IN_HEADER;
    case -1:
    case "UNRECOGNIZED":
    default:
      return SecurityScheme_In.UNRECOGNIZED;
  }
}

export function securityScheme_InToJSON(object: SecurityScheme_In): string {
  switch (object) {
    case SecurityScheme_In.IN_INVALID:
      return "IN_INVALID";
    case SecurityScheme_In.IN_QUERY:
      return "IN_QUERY";
    case SecurityScheme_In.IN_HEADER:
      return "IN_HEADER";
    case SecurityScheme_In.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

/**
 * The flow used by the OAuth2 security scheme. Valid values are
 * "implicit", "password", "application" or "accessCode".
 */
export enum SecurityScheme_Flow {
  FLOW_INVALID = 0,
  FLOW_IMPLICIT = 1,
  FLOW_PASSWORD = 2,
  FLOW_APPLICATION = 3,
  FLOW_ACCESS_CODE = 4,
  UNRECOGNIZED = -1,
}

export function securityScheme_FlowFromJSON(object: any): SecurityScheme_Flow {
  switch (object) {
    case 0:
    case "FLOW_INVALID":
      return SecurityScheme_Flow.FLOW_INVALID;
    case 1:
    case "FLOW_IMPLICIT":
      return SecurityScheme_Flow.FLOW_IMPLICIT;
    case 2:
    case "FLOW_PASSWORD":
      return SecurityScheme_Flow.FLOW_PASSWORD;
    case 3:
    case "FLOW_APPLICATION":
      return SecurityScheme_Flow.FLOW_APPLICATION;
    case 4:
    case "FLOW_ACCESS_CODE":
      return SecurityScheme_Flow.FLOW_ACCESS_CODE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return SecurityScheme_Flow.UNRECOGNIZED;
  }
}

export function securityScheme_FlowToJSON(object: SecurityScheme_Flow): string {
  switch (object) {
    case SecurityScheme_Flow.FLOW_INVALID:
      return "FLOW_INVALID";
    case SecurityScheme_Flow.FLOW_IMPLICIT:
      return "FLOW_IMPLICIT";
    case SecurityScheme_Flow.FLOW_PASSWORD:
      return "FLOW_PASSWORD";
    case SecurityScheme_Flow.FLOW_APPLICATION:
      return "FLOW_APPLICATION";
    case SecurityScheme_Flow.FLOW_ACCESS_CODE:
      return "FLOW_ACCESS_CODE";
    case SecurityScheme_Flow.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface SecurityScheme_ExtensionsEntry {
  key: string;
  value: any | undefined;
}

/**
 * `SecurityRequirement` is a representation of OpenAPI v2 specification's
 * Security Requirement object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#securityRequirementObject
 *
 * Lists the required security schemes to execute this operation. The object can
 * have multiple security schemes declared in it which are all required (that
 * is, there is a logical AND between the schemes).
 *
 * The name used for each property MUST correspond to a security scheme
 * declared in the Security Definitions.
 */
export interface SecurityRequirement {
  /**
   * Each name must correspond to a security scheme which is declared in
   * the Security Definitions. If the security scheme is of type "oauth2",
   * then the value is a list of scope names required for the execution.
   * For other security scheme types, the array MUST be empty.
   */
  securityRequirement: { [key: string]: SecurityRequirement_SecurityRequirementValue };
}

/**
 * If the security scheme is of type "oauth2", then the value is a list of
 * scope names required for the execution. For other security scheme types,
 * the array MUST be empty.
 */
export interface SecurityRequirement_SecurityRequirementValue {
  scope: string[];
}

export interface SecurityRequirement_SecurityRequirementEntry {
  key: string;
  value: SecurityRequirement_SecurityRequirementValue | undefined;
}

/**
 * `Scopes` is a representation of OpenAPI v2 specification's Scopes object.
 *
 * See: https://github.com/OAI/OpenAPI-Specification/blob/3.0.0/versions/2.0.md#scopesObject
 *
 * Lists the available scopes for an OAuth2 security scheme.
 */
export interface Scopes {
  /**
   * Maps between a name of a scope to a short description of it (as the value
   * of the property).
   */
  scope: { [key: string]: string };
}

export interface Scopes_ScopeEntry {
  key: string;
  value: string;
}

function createBaseSwagger(): Swagger {
  return {
    swagger: "",
    info: undefined,
    host: "",
    basePath: "",
    schemes: [],
    consumes: [],
    produces: [],
    responses: {},
    securityDefinitions: undefined,
    security: [],
    tags: [],
    externalDocs: undefined,
    extensions: {},
  };
}

export const Swagger = {
  encode(message: Swagger, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.swagger !== "") {
      writer.uint32(10).string(message.swagger);
    }
    if (message.info !== undefined) {
      Info.encode(message.info, writer.uint32(18).fork()).ldelim();
    }
    if (message.host !== "") {
      writer.uint32(26).string(message.host);
    }
    if (message.basePath !== "") {
      writer.uint32(34).string(message.basePath);
    }
    writer.uint32(42).fork();
    for (const v of message.schemes) {
      writer.int32(v);
    }
    writer.ldelim();
    for (const v of message.consumes) {
      writer.uint32(50).string(v!);
    }
    for (const v of message.produces) {
      writer.uint32(58).string(v!);
    }
    Object.entries(message.responses).forEach(([key, value]) => {
      Swagger_ResponsesEntry.encode({ key: key as any, value }, writer.uint32(82).fork()).ldelim();
    });
    if (message.securityDefinitions !== undefined) {
      SecurityDefinitions.encode(message.securityDefinitions, writer.uint32(90).fork()).ldelim();
    }
    for (const v of message.security) {
      SecurityRequirement.encode(v!, writer.uint32(98).fork()).ldelim();
    }
    for (const v of message.tags) {
      Tag.encode(v!, writer.uint32(106).fork()).ldelim();
    }
    if (message.externalDocs !== undefined) {
      ExternalDocumentation.encode(message.externalDocs, writer.uint32(114).fork()).ldelim();
    }
    Object.entries(message.extensions).forEach(([key, value]) => {
      if (value !== undefined) {
        Swagger_ExtensionsEntry.encode({ key: key as any, value }, writer.uint32(122).fork()).ldelim();
      }
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Swagger {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSwagger();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.swagger = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.info = Info.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.host = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.basePath = reader.string();
          continue;
        case 5:
          if (tag == 40) {
            message.schemes.push(reader.int32() as any);
            continue;
          }

          if (tag == 42) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.schemes.push(reader.int32() as any);
            }

            continue;
          }

          break;
        case 6:
          if (tag != 50) {
            break;
          }

          message.consumes.push(reader.string());
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.produces.push(reader.string());
          continue;
        case 10:
          if (tag != 82) {
            break;
          }

          const entry10 = Swagger_ResponsesEntry.decode(reader, reader.uint32());
          if (entry10.value !== undefined) {
            message.responses[entry10.key] = entry10.value;
          }
          continue;
        case 11:
          if (tag != 90) {
            break;
          }

          message.securityDefinitions = SecurityDefinitions.decode(reader, reader.uint32());
          continue;
        case 12:
          if (tag != 98) {
            break;
          }

          message.security.push(SecurityRequirement.decode(reader, reader.uint32()));
          continue;
        case 13:
          if (tag != 106) {
            break;
          }

          message.tags.push(Tag.decode(reader, reader.uint32()));
          continue;
        case 14:
          if (tag != 114) {
            break;
          }

          message.externalDocs = ExternalDocumentation.decode(reader, reader.uint32());
          continue;
        case 15:
          if (tag != 122) {
            break;
          }

          const entry15 = Swagger_ExtensionsEntry.decode(reader, reader.uint32());
          if (entry15.value !== undefined) {
            message.extensions[entry15.key] = entry15.value;
          }
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Swagger {
    return {
      swagger: isSet(object.swagger) ? String(object.swagger) : "",
      info: isSet(object.info) ? Info.fromJSON(object.info) : undefined,
      host: isSet(object.host) ? String(object.host) : "",
      basePath: isSet(object.basePath) ? String(object.basePath) : "",
      schemes: Array.isArray(object?.schemes) ? object.schemes.map((e: any) => schemeFromJSON(e)) : [],
      consumes: Array.isArray(object?.consumes) ? object.consumes.map((e: any) => String(e)) : [],
      produces: Array.isArray(object?.produces) ? object.produces.map((e: any) => String(e)) : [],
      responses: isObject(object.responses)
        ? Object.entries(object.responses).reduce<{ [key: string]: Response }>((acc, [key, value]) => {
          acc[key] = Response.fromJSON(value);
          return acc;
        }, {})
        : {},
      securityDefinitions: isSet(object.securityDefinitions)
        ? SecurityDefinitions.fromJSON(object.securityDefinitions)
        : undefined,
      security: Array.isArray(object?.security) ? object.security.map((e: any) => SecurityRequirement.fromJSON(e)) : [],
      tags: Array.isArray(object?.tags) ? object.tags.map((e: any) => Tag.fromJSON(e)) : [],
      externalDocs: isSet(object.externalDocs) ? ExternalDocumentation.fromJSON(object.externalDocs) : undefined,
      extensions: isObject(object.extensions)
        ? Object.entries(object.extensions).reduce<{ [key: string]: any | undefined }>((acc, [key, value]) => {
          acc[key] = value as any | undefined;
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: Swagger): unknown {
    const obj: any = {};
    message.swagger !== undefined && (obj.swagger = message.swagger);
    message.info !== undefined && (obj.info = message.info ? Info.toJSON(message.info) : undefined);
    message.host !== undefined && (obj.host = message.host);
    message.basePath !== undefined && (obj.basePath = message.basePath);
    if (message.schemes) {
      obj.schemes = message.schemes.map((e) => schemeToJSON(e));
    } else {
      obj.schemes = [];
    }
    if (message.consumes) {
      obj.consumes = message.consumes.map((e) => e);
    } else {
      obj.consumes = [];
    }
    if (message.produces) {
      obj.produces = message.produces.map((e) => e);
    } else {
      obj.produces = [];
    }
    obj.responses = {};
    if (message.responses) {
      Object.entries(message.responses).forEach(([k, v]) => {
        obj.responses[k] = Response.toJSON(v);
      });
    }
    message.securityDefinitions !== undefined && (obj.securityDefinitions = message.securityDefinitions
      ? SecurityDefinitions.toJSON(message.securityDefinitions)
      : undefined);
    if (message.security) {
      obj.security = message.security.map((e) => e ? SecurityRequirement.toJSON(e) : undefined);
    } else {
      obj.security = [];
    }
    if (message.tags) {
      obj.tags = message.tags.map((e) => e ? Tag.toJSON(e) : undefined);
    } else {
      obj.tags = [];
    }
    message.externalDocs !== undefined &&
      (obj.externalDocs = message.externalDocs ? ExternalDocumentation.toJSON(message.externalDocs) : undefined);
    obj.extensions = {};
    if (message.extensions) {
      Object.entries(message.extensions).forEach(([k, v]) => {
        obj.extensions[k] = v;
      });
    }
    return obj;
  },

  create(base?: DeepPartial<Swagger>): Swagger {
    return Swagger.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Swagger>): Swagger {
    const message = createBaseSwagger();
    message.swagger = object.swagger ?? "";
    message.info = (object.info !== undefined && object.info !== null) ? Info.fromPartial(object.info) : undefined;
    message.host = object.host ?? "";
    message.basePath = object.basePath ?? "";
    message.schemes = object.schemes?.map((e) => e) || [];
    message.consumes = object.consumes?.map((e) => e) || [];
    message.produces = object.produces?.map((e) => e) || [];
    message.responses = Object.entries(object.responses ?? {}).reduce<{ [key: string]: Response }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[key] = Response.fromPartial(value);
        }
        return acc;
      },
      {},
    );
    message.securityDefinitions = (object.securityDefinitions !== undefined && object.securityDefinitions !== null)
      ? SecurityDefinitions.fromPartial(object.securityDefinitions)
      : undefined;
    message.security = object.security?.map((e) => SecurityRequirement.fromPartial(e)) || [];
    message.tags = object.tags?.map((e) => Tag.fromPartial(e)) || [];
    message.externalDocs = (object.externalDocs !== undefined && object.externalDocs !== null)
      ? ExternalDocumentation.fromPartial(object.externalDocs)
      : undefined;
    message.extensions = Object.entries(object.extensions ?? {}).reduce<{ [key: string]: any | undefined }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[key] = value;
        }
        return acc;
      },
      {},
    );
    return message;
  },
};

function createBaseSwagger_ResponsesEntry(): Swagger_ResponsesEntry {
  return { key: "", value: undefined };
}

export const Swagger_ResponsesEntry = {
  encode(message: Swagger_ResponsesEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      Response.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Swagger_ResponsesEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSwagger_ResponsesEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.value = Response.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Swagger_ResponsesEntry {
    return {
      key: isSet(object.key) ? String(object.key) : "",
      value: isSet(object.value) ? Response.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: Swagger_ResponsesEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined && (obj.value = message.value ? Response.toJSON(message.value) : undefined);
    return obj;
  },

  create(base?: DeepPartial<Swagger_ResponsesEntry>): Swagger_ResponsesEntry {
    return Swagger_ResponsesEntry.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Swagger_ResponsesEntry>): Swagger_ResponsesEntry {
    const message = createBaseSwagger_ResponsesEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null)
      ? Response.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseSwagger_ExtensionsEntry(): Swagger_ExtensionsEntry {
  return { key: "", value: undefined };
}

export const Swagger_ExtensionsEntry = {
  encode(message: Swagger_ExtensionsEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      Value.encode(Value.wrap(message.value), writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Swagger_ExtensionsEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSwagger_ExtensionsEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.value = Value.unwrap(Value.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Swagger_ExtensionsEntry {
    return { key: isSet(object.key) ? String(object.key) : "", value: isSet(object?.value) ? object.value : undefined };
  },

  toJSON(message: Swagger_ExtensionsEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined && (obj.value = message.value);
    return obj;
  },

  create(base?: DeepPartial<Swagger_ExtensionsEntry>): Swagger_ExtensionsEntry {
    return Swagger_ExtensionsEntry.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Swagger_ExtensionsEntry>): Swagger_ExtensionsEntry {
    const message = createBaseSwagger_ExtensionsEntry();
    message.key = object.key ?? "";
    message.value = object.value ?? undefined;
    return message;
  },
};

function createBaseOperation(): Operation {
  return {
    tags: [],
    summary: "",
    description: "",
    externalDocs: undefined,
    operationId: "",
    consumes: [],
    produces: [],
    responses: {},
    schemes: [],
    deprecated: false,
    security: [],
    extensions: {},
    parameters: undefined,
  };
}

export const Operation = {
  encode(message: Operation, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.tags) {
      writer.uint32(10).string(v!);
    }
    if (message.summary !== "") {
      writer.uint32(18).string(message.summary);
    }
    if (message.description !== "") {
      writer.uint32(26).string(message.description);
    }
    if (message.externalDocs !== undefined) {
      ExternalDocumentation.encode(message.externalDocs, writer.uint32(34).fork()).ldelim();
    }
    if (message.operationId !== "") {
      writer.uint32(42).string(message.operationId);
    }
    for (const v of message.consumes) {
      writer.uint32(50).string(v!);
    }
    for (const v of message.produces) {
      writer.uint32(58).string(v!);
    }
    Object.entries(message.responses).forEach(([key, value]) => {
      Operation_ResponsesEntry.encode({ key: key as any, value }, writer.uint32(74).fork()).ldelim();
    });
    writer.uint32(82).fork();
    for (const v of message.schemes) {
      writer.int32(v);
    }
    writer.ldelim();
    if (message.deprecated === true) {
      writer.uint32(88).bool(message.deprecated);
    }
    for (const v of message.security) {
      SecurityRequirement.encode(v!, writer.uint32(98).fork()).ldelim();
    }
    Object.entries(message.extensions).forEach(([key, value]) => {
      if (value !== undefined) {
        Operation_ExtensionsEntry.encode({ key: key as any, value }, writer.uint32(106).fork()).ldelim();
      }
    });
    if (message.parameters !== undefined) {
      Parameters.encode(message.parameters, writer.uint32(114).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Operation {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOperation();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.tags.push(reader.string());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.summary = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.description = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.externalDocs = ExternalDocumentation.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.operationId = reader.string();
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.consumes.push(reader.string());
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.produces.push(reader.string());
          continue;
        case 9:
          if (tag != 74) {
            break;
          }

          const entry9 = Operation_ResponsesEntry.decode(reader, reader.uint32());
          if (entry9.value !== undefined) {
            message.responses[entry9.key] = entry9.value;
          }
          continue;
        case 10:
          if (tag == 80) {
            message.schemes.push(reader.int32() as any);
            continue;
          }

          if (tag == 82) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.schemes.push(reader.int32() as any);
            }

            continue;
          }

          break;
        case 11:
          if (tag != 88) {
            break;
          }

          message.deprecated = reader.bool();
          continue;
        case 12:
          if (tag != 98) {
            break;
          }

          message.security.push(SecurityRequirement.decode(reader, reader.uint32()));
          continue;
        case 13:
          if (tag != 106) {
            break;
          }

          const entry13 = Operation_ExtensionsEntry.decode(reader, reader.uint32());
          if (entry13.value !== undefined) {
            message.extensions[entry13.key] = entry13.value;
          }
          continue;
        case 14:
          if (tag != 114) {
            break;
          }

          message.parameters = Parameters.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Operation {
    return {
      tags: Array.isArray(object?.tags) ? object.tags.map((e: any) => String(e)) : [],
      summary: isSet(object.summary) ? String(object.summary) : "",
      description: isSet(object.description) ? String(object.description) : "",
      externalDocs: isSet(object.externalDocs) ? ExternalDocumentation.fromJSON(object.externalDocs) : undefined,
      operationId: isSet(object.operationId) ? String(object.operationId) : "",
      consumes: Array.isArray(object?.consumes) ? object.consumes.map((e: any) => String(e)) : [],
      produces: Array.isArray(object?.produces) ? object.produces.map((e: any) => String(e)) : [],
      responses: isObject(object.responses)
        ? Object.entries(object.responses).reduce<{ [key: string]: Response }>((acc, [key, value]) => {
          acc[key] = Response.fromJSON(value);
          return acc;
        }, {})
        : {},
      schemes: Array.isArray(object?.schemes) ? object.schemes.map((e: any) => schemeFromJSON(e)) : [],
      deprecated: isSet(object.deprecated) ? Boolean(object.deprecated) : false,
      security: Array.isArray(object?.security) ? object.security.map((e: any) => SecurityRequirement.fromJSON(e)) : [],
      extensions: isObject(object.extensions)
        ? Object.entries(object.extensions).reduce<{ [key: string]: any | undefined }>((acc, [key, value]) => {
          acc[key] = value as any | undefined;
          return acc;
        }, {})
        : {},
      parameters: isSet(object.parameters) ? Parameters.fromJSON(object.parameters) : undefined,
    };
  },

  toJSON(message: Operation): unknown {
    const obj: any = {};
    if (message.tags) {
      obj.tags = message.tags.map((e) => e);
    } else {
      obj.tags = [];
    }
    message.summary !== undefined && (obj.summary = message.summary);
    message.description !== undefined && (obj.description = message.description);
    message.externalDocs !== undefined &&
      (obj.externalDocs = message.externalDocs ? ExternalDocumentation.toJSON(message.externalDocs) : undefined);
    message.operationId !== undefined && (obj.operationId = message.operationId);
    if (message.consumes) {
      obj.consumes = message.consumes.map((e) => e);
    } else {
      obj.consumes = [];
    }
    if (message.produces) {
      obj.produces = message.produces.map((e) => e);
    } else {
      obj.produces = [];
    }
    obj.responses = {};
    if (message.responses) {
      Object.entries(message.responses).forEach(([k, v]) => {
        obj.responses[k] = Response.toJSON(v);
      });
    }
    if (message.schemes) {
      obj.schemes = message.schemes.map((e) => schemeToJSON(e));
    } else {
      obj.schemes = [];
    }
    message.deprecated !== undefined && (obj.deprecated = message.deprecated);
    if (message.security) {
      obj.security = message.security.map((e) => e ? SecurityRequirement.toJSON(e) : undefined);
    } else {
      obj.security = [];
    }
    obj.extensions = {};
    if (message.extensions) {
      Object.entries(message.extensions).forEach(([k, v]) => {
        obj.extensions[k] = v;
      });
    }
    message.parameters !== undefined &&
      (obj.parameters = message.parameters ? Parameters.toJSON(message.parameters) : undefined);
    return obj;
  },

  create(base?: DeepPartial<Operation>): Operation {
    return Operation.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Operation>): Operation {
    const message = createBaseOperation();
    message.tags = object.tags?.map((e) => e) || [];
    message.summary = object.summary ?? "";
    message.description = object.description ?? "";
    message.externalDocs = (object.externalDocs !== undefined && object.externalDocs !== null)
      ? ExternalDocumentation.fromPartial(object.externalDocs)
      : undefined;
    message.operationId = object.operationId ?? "";
    message.consumes = object.consumes?.map((e) => e) || [];
    message.produces = object.produces?.map((e) => e) || [];
    message.responses = Object.entries(object.responses ?? {}).reduce<{ [key: string]: Response }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[key] = Response.fromPartial(value);
        }
        return acc;
      },
      {},
    );
    message.schemes = object.schemes?.map((e) => e) || [];
    message.deprecated = object.deprecated ?? false;
    message.security = object.security?.map((e) => SecurityRequirement.fromPartial(e)) || [];
    message.extensions = Object.entries(object.extensions ?? {}).reduce<{ [key: string]: any | undefined }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[key] = value;
        }
        return acc;
      },
      {},
    );
    message.parameters = (object.parameters !== undefined && object.parameters !== null)
      ? Parameters.fromPartial(object.parameters)
      : undefined;
    return message;
  },
};

function createBaseOperation_ResponsesEntry(): Operation_ResponsesEntry {
  return { key: "", value: undefined };
}

export const Operation_ResponsesEntry = {
  encode(message: Operation_ResponsesEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      Response.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Operation_ResponsesEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOperation_ResponsesEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.value = Response.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Operation_ResponsesEntry {
    return {
      key: isSet(object.key) ? String(object.key) : "",
      value: isSet(object.value) ? Response.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: Operation_ResponsesEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined && (obj.value = message.value ? Response.toJSON(message.value) : undefined);
    return obj;
  },

  create(base?: DeepPartial<Operation_ResponsesEntry>): Operation_ResponsesEntry {
    return Operation_ResponsesEntry.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Operation_ResponsesEntry>): Operation_ResponsesEntry {
    const message = createBaseOperation_ResponsesEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null)
      ? Response.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseOperation_ExtensionsEntry(): Operation_ExtensionsEntry {
  return { key: "", value: undefined };
}

export const Operation_ExtensionsEntry = {
  encode(message: Operation_ExtensionsEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      Value.encode(Value.wrap(message.value), writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Operation_ExtensionsEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOperation_ExtensionsEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.value = Value.unwrap(Value.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Operation_ExtensionsEntry {
    return { key: isSet(object.key) ? String(object.key) : "", value: isSet(object?.value) ? object.value : undefined };
  },

  toJSON(message: Operation_ExtensionsEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined && (obj.value = message.value);
    return obj;
  },

  create(base?: DeepPartial<Operation_ExtensionsEntry>): Operation_ExtensionsEntry {
    return Operation_ExtensionsEntry.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Operation_ExtensionsEntry>): Operation_ExtensionsEntry {
    const message = createBaseOperation_ExtensionsEntry();
    message.key = object.key ?? "";
    message.value = object.value ?? undefined;
    return message;
  },
};

function createBaseParameters(): Parameters {
  return { headers: [] };
}

export const Parameters = {
  encode(message: Parameters, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.headers) {
      HeaderParameter.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Parameters {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseParameters();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.headers.push(HeaderParameter.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Parameters {
    return {
      headers: Array.isArray(object?.headers) ? object.headers.map((e: any) => HeaderParameter.fromJSON(e)) : [],
    };
  },

  toJSON(message: Parameters): unknown {
    const obj: any = {};
    if (message.headers) {
      obj.headers = message.headers.map((e) => e ? HeaderParameter.toJSON(e) : undefined);
    } else {
      obj.headers = [];
    }
    return obj;
  },

  create(base?: DeepPartial<Parameters>): Parameters {
    return Parameters.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Parameters>): Parameters {
    const message = createBaseParameters();
    message.headers = object.headers?.map((e) => HeaderParameter.fromPartial(e)) || [];
    return message;
  },
};

function createBaseHeaderParameter(): HeaderParameter {
  return { name: "", description: "", type: 0, format: "", required: false };
}

export const HeaderParameter = {
  encode(message: HeaderParameter, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.type !== 0) {
      writer.uint32(24).int32(message.type);
    }
    if (message.format !== "") {
      writer.uint32(34).string(message.format);
    }
    if (message.required === true) {
      writer.uint32(40).bool(message.required);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): HeaderParameter {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseHeaderParameter();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.description = reader.string();
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.type = reader.int32() as any;
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.format = reader.string();
          continue;
        case 5:
          if (tag != 40) {
            break;
          }

          message.required = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): HeaderParameter {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      description: isSet(object.description) ? String(object.description) : "",
      type: isSet(object.type) ? headerParameter_TypeFromJSON(object.type) : 0,
      format: isSet(object.format) ? String(object.format) : "",
      required: isSet(object.required) ? Boolean(object.required) : false,
    };
  },

  toJSON(message: HeaderParameter): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.description !== undefined && (obj.description = message.description);
    message.type !== undefined && (obj.type = headerParameter_TypeToJSON(message.type));
    message.format !== undefined && (obj.format = message.format);
    message.required !== undefined && (obj.required = message.required);
    return obj;
  },

  create(base?: DeepPartial<HeaderParameter>): HeaderParameter {
    return HeaderParameter.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<HeaderParameter>): HeaderParameter {
    const message = createBaseHeaderParameter();
    message.name = object.name ?? "";
    message.description = object.description ?? "";
    message.type = object.type ?? 0;
    message.format = object.format ?? "";
    message.required = object.required ?? false;
    return message;
  },
};

function createBaseHeader(): Header {
  return { description: "", type: "", format: "", default: "", pattern: "" };
}

export const Header = {
  encode(message: Header, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.description !== "") {
      writer.uint32(10).string(message.description);
    }
    if (message.type !== "") {
      writer.uint32(18).string(message.type);
    }
    if (message.format !== "") {
      writer.uint32(26).string(message.format);
    }
    if (message.default !== "") {
      writer.uint32(50).string(message.default);
    }
    if (message.pattern !== "") {
      writer.uint32(106).string(message.pattern);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Header {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseHeader();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.description = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.type = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.format = reader.string();
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.default = reader.string();
          continue;
        case 13:
          if (tag != 106) {
            break;
          }

          message.pattern = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Header {
    return {
      description: isSet(object.description) ? String(object.description) : "",
      type: isSet(object.type) ? String(object.type) : "",
      format: isSet(object.format) ? String(object.format) : "",
      default: isSet(object.default) ? String(object.default) : "",
      pattern: isSet(object.pattern) ? String(object.pattern) : "",
    };
  },

  toJSON(message: Header): unknown {
    const obj: any = {};
    message.description !== undefined && (obj.description = message.description);
    message.type !== undefined && (obj.type = message.type);
    message.format !== undefined && (obj.format = message.format);
    message.default !== undefined && (obj.default = message.default);
    message.pattern !== undefined && (obj.pattern = message.pattern);
    return obj;
  },

  create(base?: DeepPartial<Header>): Header {
    return Header.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Header>): Header {
    const message = createBaseHeader();
    message.description = object.description ?? "";
    message.type = object.type ?? "";
    message.format = object.format ?? "";
    message.default = object.default ?? "";
    message.pattern = object.pattern ?? "";
    return message;
  },
};

function createBaseResponse(): Response {
  return { description: "", schema: undefined, headers: {}, examples: {}, extensions: {} };
}

export const Response = {
  encode(message: Response, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.description !== "") {
      writer.uint32(10).string(message.description);
    }
    if (message.schema !== undefined) {
      Schema.encode(message.schema, writer.uint32(18).fork()).ldelim();
    }
    Object.entries(message.headers).forEach(([key, value]) => {
      Response_HeadersEntry.encode({ key: key as any, value }, writer.uint32(26).fork()).ldelim();
    });
    Object.entries(message.examples).forEach(([key, value]) => {
      Response_ExamplesEntry.encode({ key: key as any, value }, writer.uint32(34).fork()).ldelim();
    });
    Object.entries(message.extensions).forEach(([key, value]) => {
      if (value !== undefined) {
        Response_ExtensionsEntry.encode({ key: key as any, value }, writer.uint32(42).fork()).ldelim();
      }
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Response {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.description = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.schema = Schema.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          const entry3 = Response_HeadersEntry.decode(reader, reader.uint32());
          if (entry3.value !== undefined) {
            message.headers[entry3.key] = entry3.value;
          }
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          const entry4 = Response_ExamplesEntry.decode(reader, reader.uint32());
          if (entry4.value !== undefined) {
            message.examples[entry4.key] = entry4.value;
          }
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          const entry5 = Response_ExtensionsEntry.decode(reader, reader.uint32());
          if (entry5.value !== undefined) {
            message.extensions[entry5.key] = entry5.value;
          }
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Response {
    return {
      description: isSet(object.description) ? String(object.description) : "",
      schema: isSet(object.schema) ? Schema.fromJSON(object.schema) : undefined,
      headers: isObject(object.headers)
        ? Object.entries(object.headers).reduce<{ [key: string]: Header }>((acc, [key, value]) => {
          acc[key] = Header.fromJSON(value);
          return acc;
        }, {})
        : {},
      examples: isObject(object.examples)
        ? Object.entries(object.examples).reduce<{ [key: string]: string }>((acc, [key, value]) => {
          acc[key] = String(value);
          return acc;
        }, {})
        : {},
      extensions: isObject(object.extensions)
        ? Object.entries(object.extensions).reduce<{ [key: string]: any | undefined }>((acc, [key, value]) => {
          acc[key] = value as any | undefined;
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: Response): unknown {
    const obj: any = {};
    message.description !== undefined && (obj.description = message.description);
    message.schema !== undefined && (obj.schema = message.schema ? Schema.toJSON(message.schema) : undefined);
    obj.headers = {};
    if (message.headers) {
      Object.entries(message.headers).forEach(([k, v]) => {
        obj.headers[k] = Header.toJSON(v);
      });
    }
    obj.examples = {};
    if (message.examples) {
      Object.entries(message.examples).forEach(([k, v]) => {
        obj.examples[k] = v;
      });
    }
    obj.extensions = {};
    if (message.extensions) {
      Object.entries(message.extensions).forEach(([k, v]) => {
        obj.extensions[k] = v;
      });
    }
    return obj;
  },

  create(base?: DeepPartial<Response>): Response {
    return Response.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Response>): Response {
    const message = createBaseResponse();
    message.description = object.description ?? "";
    message.schema = (object.schema !== undefined && object.schema !== null)
      ? Schema.fromPartial(object.schema)
      : undefined;
    message.headers = Object.entries(object.headers ?? {}).reduce<{ [key: string]: Header }>((acc, [key, value]) => {
      if (value !== undefined) {
        acc[key] = Header.fromPartial(value);
      }
      return acc;
    }, {});
    message.examples = Object.entries(object.examples ?? {}).reduce<{ [key: string]: string }>((acc, [key, value]) => {
      if (value !== undefined) {
        acc[key] = String(value);
      }
      return acc;
    }, {});
    message.extensions = Object.entries(object.extensions ?? {}).reduce<{ [key: string]: any | undefined }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[key] = value;
        }
        return acc;
      },
      {},
    );
    return message;
  },
};

function createBaseResponse_HeadersEntry(): Response_HeadersEntry {
  return { key: "", value: undefined };
}

export const Response_HeadersEntry = {
  encode(message: Response_HeadersEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      Header.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Response_HeadersEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResponse_HeadersEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.value = Header.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Response_HeadersEntry {
    return {
      key: isSet(object.key) ? String(object.key) : "",
      value: isSet(object.value) ? Header.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: Response_HeadersEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined && (obj.value = message.value ? Header.toJSON(message.value) : undefined);
    return obj;
  },

  create(base?: DeepPartial<Response_HeadersEntry>): Response_HeadersEntry {
    return Response_HeadersEntry.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Response_HeadersEntry>): Response_HeadersEntry {
    const message = createBaseResponse_HeadersEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null)
      ? Header.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseResponse_ExamplesEntry(): Response_ExamplesEntry {
  return { key: "", value: "" };
}

export const Response_ExamplesEntry = {
  encode(message: Response_ExamplesEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== "") {
      writer.uint32(18).string(message.value);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Response_ExamplesEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResponse_ExamplesEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.value = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Response_ExamplesEntry {
    return { key: isSet(object.key) ? String(object.key) : "", value: isSet(object.value) ? String(object.value) : "" };
  },

  toJSON(message: Response_ExamplesEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined && (obj.value = message.value);
    return obj;
  },

  create(base?: DeepPartial<Response_ExamplesEntry>): Response_ExamplesEntry {
    return Response_ExamplesEntry.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Response_ExamplesEntry>): Response_ExamplesEntry {
    const message = createBaseResponse_ExamplesEntry();
    message.key = object.key ?? "";
    message.value = object.value ?? "";
    return message;
  },
};

function createBaseResponse_ExtensionsEntry(): Response_ExtensionsEntry {
  return { key: "", value: undefined };
}

export const Response_ExtensionsEntry = {
  encode(message: Response_ExtensionsEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      Value.encode(Value.wrap(message.value), writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Response_ExtensionsEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResponse_ExtensionsEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.value = Value.unwrap(Value.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Response_ExtensionsEntry {
    return { key: isSet(object.key) ? String(object.key) : "", value: isSet(object?.value) ? object.value : undefined };
  },

  toJSON(message: Response_ExtensionsEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined && (obj.value = message.value);
    return obj;
  },

  create(base?: DeepPartial<Response_ExtensionsEntry>): Response_ExtensionsEntry {
    return Response_ExtensionsEntry.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Response_ExtensionsEntry>): Response_ExtensionsEntry {
    const message = createBaseResponse_ExtensionsEntry();
    message.key = object.key ?? "";
    message.value = object.value ?? undefined;
    return message;
  },
};

function createBaseInfo(): Info {
  return {
    title: "",
    description: "",
    termsOfService: "",
    contact: undefined,
    license: undefined,
    version: "",
    extensions: {},
  };
}

export const Info = {
  encode(message: Info, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.termsOfService !== "") {
      writer.uint32(26).string(message.termsOfService);
    }
    if (message.contact !== undefined) {
      Contact.encode(message.contact, writer.uint32(34).fork()).ldelim();
    }
    if (message.license !== undefined) {
      License.encode(message.license, writer.uint32(42).fork()).ldelim();
    }
    if (message.version !== "") {
      writer.uint32(50).string(message.version);
    }
    Object.entries(message.extensions).forEach(([key, value]) => {
      if (value !== undefined) {
        Info_ExtensionsEntry.encode({ key: key as any, value }, writer.uint32(58).fork()).ldelim();
      }
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Info {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInfo();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.title = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.description = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.termsOfService = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.contact = Contact.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.license = License.decode(reader, reader.uint32());
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.version = reader.string();
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          const entry7 = Info_ExtensionsEntry.decode(reader, reader.uint32());
          if (entry7.value !== undefined) {
            message.extensions[entry7.key] = entry7.value;
          }
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Info {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      termsOfService: isSet(object.termsOfService) ? String(object.termsOfService) : "",
      contact: isSet(object.contact) ? Contact.fromJSON(object.contact) : undefined,
      license: isSet(object.license) ? License.fromJSON(object.license) : undefined,
      version: isSet(object.version) ? String(object.version) : "",
      extensions: isObject(object.extensions)
        ? Object.entries(object.extensions).reduce<{ [key: string]: any | undefined }>((acc, [key, value]) => {
          acc[key] = value as any | undefined;
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: Info): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.termsOfService !== undefined && (obj.termsOfService = message.termsOfService);
    message.contact !== undefined && (obj.contact = message.contact ? Contact.toJSON(message.contact) : undefined);
    message.license !== undefined && (obj.license = message.license ? License.toJSON(message.license) : undefined);
    message.version !== undefined && (obj.version = message.version);
    obj.extensions = {};
    if (message.extensions) {
      Object.entries(message.extensions).forEach(([k, v]) => {
        obj.extensions[k] = v;
      });
    }
    return obj;
  },

  create(base?: DeepPartial<Info>): Info {
    return Info.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Info>): Info {
    const message = createBaseInfo();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.termsOfService = object.termsOfService ?? "";
    message.contact = (object.contact !== undefined && object.contact !== null)
      ? Contact.fromPartial(object.contact)
      : undefined;
    message.license = (object.license !== undefined && object.license !== null)
      ? License.fromPartial(object.license)
      : undefined;
    message.version = object.version ?? "";
    message.extensions = Object.entries(object.extensions ?? {}).reduce<{ [key: string]: any | undefined }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[key] = value;
        }
        return acc;
      },
      {},
    );
    return message;
  },
};

function createBaseInfo_ExtensionsEntry(): Info_ExtensionsEntry {
  return { key: "", value: undefined };
}

export const Info_ExtensionsEntry = {
  encode(message: Info_ExtensionsEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      Value.encode(Value.wrap(message.value), writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Info_ExtensionsEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInfo_ExtensionsEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.value = Value.unwrap(Value.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Info_ExtensionsEntry {
    return { key: isSet(object.key) ? String(object.key) : "", value: isSet(object?.value) ? object.value : undefined };
  },

  toJSON(message: Info_ExtensionsEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined && (obj.value = message.value);
    return obj;
  },

  create(base?: DeepPartial<Info_ExtensionsEntry>): Info_ExtensionsEntry {
    return Info_ExtensionsEntry.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Info_ExtensionsEntry>): Info_ExtensionsEntry {
    const message = createBaseInfo_ExtensionsEntry();
    message.key = object.key ?? "";
    message.value = object.value ?? undefined;
    return message;
  },
};

function createBaseContact(): Contact {
  return { name: "", url: "", email: "" };
}

export const Contact = {
  encode(message: Contact, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.url !== "") {
      writer.uint32(18).string(message.url);
    }
    if (message.email !== "") {
      writer.uint32(26).string(message.email);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Contact {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseContact();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.url = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.email = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Contact {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      url: isSet(object.url) ? String(object.url) : "",
      email: isSet(object.email) ? String(object.email) : "",
    };
  },

  toJSON(message: Contact): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.url !== undefined && (obj.url = message.url);
    message.email !== undefined && (obj.email = message.email);
    return obj;
  },

  create(base?: DeepPartial<Contact>): Contact {
    return Contact.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Contact>): Contact {
    const message = createBaseContact();
    message.name = object.name ?? "";
    message.url = object.url ?? "";
    message.email = object.email ?? "";
    return message;
  },
};

function createBaseLicense(): License {
  return { name: "", url: "" };
}

export const License = {
  encode(message: License, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.url !== "") {
      writer.uint32(18).string(message.url);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): License {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLicense();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.url = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): License {
    return { name: isSet(object.name) ? String(object.name) : "", url: isSet(object.url) ? String(object.url) : "" };
  },

  toJSON(message: License): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.url !== undefined && (obj.url = message.url);
    return obj;
  },

  create(base?: DeepPartial<License>): License {
    return License.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<License>): License {
    const message = createBaseLicense();
    message.name = object.name ?? "";
    message.url = object.url ?? "";
    return message;
  },
};

function createBaseExternalDocumentation(): ExternalDocumentation {
  return { description: "", url: "" };
}

export const ExternalDocumentation = {
  encode(message: ExternalDocumentation, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.description !== "") {
      writer.uint32(10).string(message.description);
    }
    if (message.url !== "") {
      writer.uint32(18).string(message.url);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ExternalDocumentation {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExternalDocumentation();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.description = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.url = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ExternalDocumentation {
    return {
      description: isSet(object.description) ? String(object.description) : "",
      url: isSet(object.url) ? String(object.url) : "",
    };
  },

  toJSON(message: ExternalDocumentation): unknown {
    const obj: any = {};
    message.description !== undefined && (obj.description = message.description);
    message.url !== undefined && (obj.url = message.url);
    return obj;
  },

  create(base?: DeepPartial<ExternalDocumentation>): ExternalDocumentation {
    return ExternalDocumentation.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ExternalDocumentation>): ExternalDocumentation {
    const message = createBaseExternalDocumentation();
    message.description = object.description ?? "";
    message.url = object.url ?? "";
    return message;
  },
};

function createBaseSchema(): Schema {
  return { jsonSchema: undefined, discriminator: "", readOnly: false, externalDocs: undefined, example: "" };
}

export const Schema = {
  encode(message: Schema, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.jsonSchema !== undefined) {
      JSONSchema.encode(message.jsonSchema, writer.uint32(10).fork()).ldelim();
    }
    if (message.discriminator !== "") {
      writer.uint32(18).string(message.discriminator);
    }
    if (message.readOnly === true) {
      writer.uint32(24).bool(message.readOnly);
    }
    if (message.externalDocs !== undefined) {
      ExternalDocumentation.encode(message.externalDocs, writer.uint32(42).fork()).ldelim();
    }
    if (message.example !== "") {
      writer.uint32(50).string(message.example);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Schema {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSchema();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.jsonSchema = JSONSchema.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.discriminator = reader.string();
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.readOnly = reader.bool();
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.externalDocs = ExternalDocumentation.decode(reader, reader.uint32());
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.example = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Schema {
    return {
      jsonSchema: isSet(object.jsonSchema) ? JSONSchema.fromJSON(object.jsonSchema) : undefined,
      discriminator: isSet(object.discriminator) ? String(object.discriminator) : "",
      readOnly: isSet(object.readOnly) ? Boolean(object.readOnly) : false,
      externalDocs: isSet(object.externalDocs) ? ExternalDocumentation.fromJSON(object.externalDocs) : undefined,
      example: isSet(object.example) ? String(object.example) : "",
    };
  },

  toJSON(message: Schema): unknown {
    const obj: any = {};
    message.jsonSchema !== undefined &&
      (obj.jsonSchema = message.jsonSchema ? JSONSchema.toJSON(message.jsonSchema) : undefined);
    message.discriminator !== undefined && (obj.discriminator = message.discriminator);
    message.readOnly !== undefined && (obj.readOnly = message.readOnly);
    message.externalDocs !== undefined &&
      (obj.externalDocs = message.externalDocs ? ExternalDocumentation.toJSON(message.externalDocs) : undefined);
    message.example !== undefined && (obj.example = message.example);
    return obj;
  },

  create(base?: DeepPartial<Schema>): Schema {
    return Schema.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Schema>): Schema {
    const message = createBaseSchema();
    message.jsonSchema = (object.jsonSchema !== undefined && object.jsonSchema !== null)
      ? JSONSchema.fromPartial(object.jsonSchema)
      : undefined;
    message.discriminator = object.discriminator ?? "";
    message.readOnly = object.readOnly ?? false;
    message.externalDocs = (object.externalDocs !== undefined && object.externalDocs !== null)
      ? ExternalDocumentation.fromPartial(object.externalDocs)
      : undefined;
    message.example = object.example ?? "";
    return message;
  },
};

function createBaseJSONSchema(): JSONSchema {
  return {
    ref: "",
    title: "",
    description: "",
    default: "",
    readOnly: false,
    example: "",
    multipleOf: 0,
    maximum: 0,
    exclusiveMaximum: false,
    minimum: 0,
    exclusiveMinimum: false,
    maxLength: 0,
    minLength: 0,
    pattern: "",
    maxItems: 0,
    minItems: 0,
    uniqueItems: false,
    maxProperties: 0,
    minProperties: 0,
    required: [],
    array: [],
    type: [],
    format: "",
    enum: [],
    fieldConfiguration: undefined,
    extensions: {},
  };
}

export const JSONSchema = {
  encode(message: JSONSchema, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.ref !== "") {
      writer.uint32(26).string(message.ref);
    }
    if (message.title !== "") {
      writer.uint32(42).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(50).string(message.description);
    }
    if (message.default !== "") {
      writer.uint32(58).string(message.default);
    }
    if (message.readOnly === true) {
      writer.uint32(64).bool(message.readOnly);
    }
    if (message.example !== "") {
      writer.uint32(74).string(message.example);
    }
    if (message.multipleOf !== 0) {
      writer.uint32(81).double(message.multipleOf);
    }
    if (message.maximum !== 0) {
      writer.uint32(89).double(message.maximum);
    }
    if (message.exclusiveMaximum === true) {
      writer.uint32(96).bool(message.exclusiveMaximum);
    }
    if (message.minimum !== 0) {
      writer.uint32(105).double(message.minimum);
    }
    if (message.exclusiveMinimum === true) {
      writer.uint32(112).bool(message.exclusiveMinimum);
    }
    if (message.maxLength !== 0) {
      writer.uint32(120).uint64(message.maxLength);
    }
    if (message.minLength !== 0) {
      writer.uint32(128).uint64(message.minLength);
    }
    if (message.pattern !== "") {
      writer.uint32(138).string(message.pattern);
    }
    if (message.maxItems !== 0) {
      writer.uint32(160).uint64(message.maxItems);
    }
    if (message.minItems !== 0) {
      writer.uint32(168).uint64(message.minItems);
    }
    if (message.uniqueItems === true) {
      writer.uint32(176).bool(message.uniqueItems);
    }
    if (message.maxProperties !== 0) {
      writer.uint32(192).uint64(message.maxProperties);
    }
    if (message.minProperties !== 0) {
      writer.uint32(200).uint64(message.minProperties);
    }
    for (const v of message.required) {
      writer.uint32(210).string(v!);
    }
    for (const v of message.array) {
      writer.uint32(274).string(v!);
    }
    writer.uint32(282).fork();
    for (const v of message.type) {
      writer.int32(v);
    }
    writer.ldelim();
    if (message.format !== "") {
      writer.uint32(290).string(message.format);
    }
    for (const v of message.enum) {
      writer.uint32(370).string(v!);
    }
    if (message.fieldConfiguration !== undefined) {
      JSONSchema_FieldConfiguration.encode(message.fieldConfiguration, writer.uint32(8010).fork()).ldelim();
    }
    Object.entries(message.extensions).forEach(([key, value]) => {
      if (value !== undefined) {
        JSONSchema_ExtensionsEntry.encode({ key: key as any, value }, writer.uint32(386).fork()).ldelim();
      }
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): JSONSchema {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseJSONSchema();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 3:
          if (tag != 26) {
            break;
          }

          message.ref = reader.string();
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.title = reader.string();
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.description = reader.string();
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.default = reader.string();
          continue;
        case 8:
          if (tag != 64) {
            break;
          }

          message.readOnly = reader.bool();
          continue;
        case 9:
          if (tag != 74) {
            break;
          }

          message.example = reader.string();
          continue;
        case 10:
          if (tag != 81) {
            break;
          }

          message.multipleOf = reader.double();
          continue;
        case 11:
          if (tag != 89) {
            break;
          }

          message.maximum = reader.double();
          continue;
        case 12:
          if (tag != 96) {
            break;
          }

          message.exclusiveMaximum = reader.bool();
          continue;
        case 13:
          if (tag != 105) {
            break;
          }

          message.minimum = reader.double();
          continue;
        case 14:
          if (tag != 112) {
            break;
          }

          message.exclusiveMinimum = reader.bool();
          continue;
        case 15:
          if (tag != 120) {
            break;
          }

          message.maxLength = longToNumber(reader.uint64() as Long);
          continue;
        case 16:
          if (tag != 128) {
            break;
          }

          message.minLength = longToNumber(reader.uint64() as Long);
          continue;
        case 17:
          if (tag != 138) {
            break;
          }

          message.pattern = reader.string();
          continue;
        case 20:
          if (tag != 160) {
            break;
          }

          message.maxItems = longToNumber(reader.uint64() as Long);
          continue;
        case 21:
          if (tag != 168) {
            break;
          }

          message.minItems = longToNumber(reader.uint64() as Long);
          continue;
        case 22:
          if (tag != 176) {
            break;
          }

          message.uniqueItems = reader.bool();
          continue;
        case 24:
          if (tag != 192) {
            break;
          }

          message.maxProperties = longToNumber(reader.uint64() as Long);
          continue;
        case 25:
          if (tag != 200) {
            break;
          }

          message.minProperties = longToNumber(reader.uint64() as Long);
          continue;
        case 26:
          if (tag != 210) {
            break;
          }

          message.required.push(reader.string());
          continue;
        case 34:
          if (tag != 274) {
            break;
          }

          message.array.push(reader.string());
          continue;
        case 35:
          if (tag == 280) {
            message.type.push(reader.int32() as any);
            continue;
          }

          if (tag == 282) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.type.push(reader.int32() as any);
            }

            continue;
          }

          break;
        case 36:
          if (tag != 290) {
            break;
          }

          message.format = reader.string();
          continue;
        case 46:
          if (tag != 370) {
            break;
          }

          message.enum.push(reader.string());
          continue;
        case 1001:
          if (tag != 8010) {
            break;
          }

          message.fieldConfiguration = JSONSchema_FieldConfiguration.decode(reader, reader.uint32());
          continue;
        case 48:
          if (tag != 386) {
            break;
          }

          const entry48 = JSONSchema_ExtensionsEntry.decode(reader, reader.uint32());
          if (entry48.value !== undefined) {
            message.extensions[entry48.key] = entry48.value;
          }
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): JSONSchema {
    return {
      ref: isSet(object.ref) ? String(object.ref) : "",
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      default: isSet(object.default) ? String(object.default) : "",
      readOnly: isSet(object.readOnly) ? Boolean(object.readOnly) : false,
      example: isSet(object.example) ? String(object.example) : "",
      multipleOf: isSet(object.multipleOf) ? Number(object.multipleOf) : 0,
      maximum: isSet(object.maximum) ? Number(object.maximum) : 0,
      exclusiveMaximum: isSet(object.exclusiveMaximum) ? Boolean(object.exclusiveMaximum) : false,
      minimum: isSet(object.minimum) ? Number(object.minimum) : 0,
      exclusiveMinimum: isSet(object.exclusiveMinimum) ? Boolean(object.exclusiveMinimum) : false,
      maxLength: isSet(object.maxLength) ? Number(object.maxLength) : 0,
      minLength: isSet(object.minLength) ? Number(object.minLength) : 0,
      pattern: isSet(object.pattern) ? String(object.pattern) : "",
      maxItems: isSet(object.maxItems) ? Number(object.maxItems) : 0,
      minItems: isSet(object.minItems) ? Number(object.minItems) : 0,
      uniqueItems: isSet(object.uniqueItems) ? Boolean(object.uniqueItems) : false,
      maxProperties: isSet(object.maxProperties) ? Number(object.maxProperties) : 0,
      minProperties: isSet(object.minProperties) ? Number(object.minProperties) : 0,
      required: Array.isArray(object?.required) ? object.required.map((e: any) => String(e)) : [],
      array: Array.isArray(object?.array) ? object.array.map((e: any) => String(e)) : [],
      type: Array.isArray(object?.type) ? object.type.map((e: any) => jSONSchema_JSONSchemaSimpleTypesFromJSON(e)) : [],
      format: isSet(object.format) ? String(object.format) : "",
      enum: Array.isArray(object?.enum) ? object.enum.map((e: any) => String(e)) : [],
      fieldConfiguration: isSet(object.fieldConfiguration)
        ? JSONSchema_FieldConfiguration.fromJSON(object.fieldConfiguration)
        : undefined,
      extensions: isObject(object.extensions)
        ? Object.entries(object.extensions).reduce<{ [key: string]: any | undefined }>((acc, [key, value]) => {
          acc[key] = value as any | undefined;
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: JSONSchema): unknown {
    const obj: any = {};
    message.ref !== undefined && (obj.ref = message.ref);
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.default !== undefined && (obj.default = message.default);
    message.readOnly !== undefined && (obj.readOnly = message.readOnly);
    message.example !== undefined && (obj.example = message.example);
    message.multipleOf !== undefined && (obj.multipleOf = message.multipleOf);
    message.maximum !== undefined && (obj.maximum = message.maximum);
    message.exclusiveMaximum !== undefined && (obj.exclusiveMaximum = message.exclusiveMaximum);
    message.minimum !== undefined && (obj.minimum = message.minimum);
    message.exclusiveMinimum !== undefined && (obj.exclusiveMinimum = message.exclusiveMinimum);
    message.maxLength !== undefined && (obj.maxLength = Math.round(message.maxLength));
    message.minLength !== undefined && (obj.minLength = Math.round(message.minLength));
    message.pattern !== undefined && (obj.pattern = message.pattern);
    message.maxItems !== undefined && (obj.maxItems = Math.round(message.maxItems));
    message.minItems !== undefined && (obj.minItems = Math.round(message.minItems));
    message.uniqueItems !== undefined && (obj.uniqueItems = message.uniqueItems);
    message.maxProperties !== undefined && (obj.maxProperties = Math.round(message.maxProperties));
    message.minProperties !== undefined && (obj.minProperties = Math.round(message.minProperties));
    if (message.required) {
      obj.required = message.required.map((e) => e);
    } else {
      obj.required = [];
    }
    if (message.array) {
      obj.array = message.array.map((e) => e);
    } else {
      obj.array = [];
    }
    if (message.type) {
      obj.type = message.type.map((e) => jSONSchema_JSONSchemaSimpleTypesToJSON(e));
    } else {
      obj.type = [];
    }
    message.format !== undefined && (obj.format = message.format);
    if (message.enum) {
      obj.enum = message.enum.map((e) => e);
    } else {
      obj.enum = [];
    }
    message.fieldConfiguration !== undefined && (obj.fieldConfiguration = message.fieldConfiguration
      ? JSONSchema_FieldConfiguration.toJSON(message.fieldConfiguration)
      : undefined);
    obj.extensions = {};
    if (message.extensions) {
      Object.entries(message.extensions).forEach(([k, v]) => {
        obj.extensions[k] = v;
      });
    }
    return obj;
  },

  create(base?: DeepPartial<JSONSchema>): JSONSchema {
    return JSONSchema.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<JSONSchema>): JSONSchema {
    const message = createBaseJSONSchema();
    message.ref = object.ref ?? "";
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.default = object.default ?? "";
    message.readOnly = object.readOnly ?? false;
    message.example = object.example ?? "";
    message.multipleOf = object.multipleOf ?? 0;
    message.maximum = object.maximum ?? 0;
    message.exclusiveMaximum = object.exclusiveMaximum ?? false;
    message.minimum = object.minimum ?? 0;
    message.exclusiveMinimum = object.exclusiveMinimum ?? false;
    message.maxLength = object.maxLength ?? 0;
    message.minLength = object.minLength ?? 0;
    message.pattern = object.pattern ?? "";
    message.maxItems = object.maxItems ?? 0;
    message.minItems = object.minItems ?? 0;
    message.uniqueItems = object.uniqueItems ?? false;
    message.maxProperties = object.maxProperties ?? 0;
    message.minProperties = object.minProperties ?? 0;
    message.required = object.required?.map((e) => e) || [];
    message.array = object.array?.map((e) => e) || [];
    message.type = object.type?.map((e) => e) || [];
    message.format = object.format ?? "";
    message.enum = object.enum?.map((e) => e) || [];
    message.fieldConfiguration = (object.fieldConfiguration !== undefined && object.fieldConfiguration !== null)
      ? JSONSchema_FieldConfiguration.fromPartial(object.fieldConfiguration)
      : undefined;
    message.extensions = Object.entries(object.extensions ?? {}).reduce<{ [key: string]: any | undefined }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[key] = value;
        }
        return acc;
      },
      {},
    );
    return message;
  },
};

function createBaseJSONSchema_FieldConfiguration(): JSONSchema_FieldConfiguration {
  return { pathParamName: "" };
}

export const JSONSchema_FieldConfiguration = {
  encode(message: JSONSchema_FieldConfiguration, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pathParamName !== "") {
      writer.uint32(378).string(message.pathParamName);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): JSONSchema_FieldConfiguration {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseJSONSchema_FieldConfiguration();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 47:
          if (tag != 378) {
            break;
          }

          message.pathParamName = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): JSONSchema_FieldConfiguration {
    return { pathParamName: isSet(object.pathParamName) ? String(object.pathParamName) : "" };
  },

  toJSON(message: JSONSchema_FieldConfiguration): unknown {
    const obj: any = {};
    message.pathParamName !== undefined && (obj.pathParamName = message.pathParamName);
    return obj;
  },

  create(base?: DeepPartial<JSONSchema_FieldConfiguration>): JSONSchema_FieldConfiguration {
    return JSONSchema_FieldConfiguration.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<JSONSchema_FieldConfiguration>): JSONSchema_FieldConfiguration {
    const message = createBaseJSONSchema_FieldConfiguration();
    message.pathParamName = object.pathParamName ?? "";
    return message;
  },
};

function createBaseJSONSchema_ExtensionsEntry(): JSONSchema_ExtensionsEntry {
  return { key: "", value: undefined };
}

export const JSONSchema_ExtensionsEntry = {
  encode(message: JSONSchema_ExtensionsEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      Value.encode(Value.wrap(message.value), writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): JSONSchema_ExtensionsEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseJSONSchema_ExtensionsEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.value = Value.unwrap(Value.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): JSONSchema_ExtensionsEntry {
    return { key: isSet(object.key) ? String(object.key) : "", value: isSet(object?.value) ? object.value : undefined };
  },

  toJSON(message: JSONSchema_ExtensionsEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined && (obj.value = message.value);
    return obj;
  },

  create(base?: DeepPartial<JSONSchema_ExtensionsEntry>): JSONSchema_ExtensionsEntry {
    return JSONSchema_ExtensionsEntry.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<JSONSchema_ExtensionsEntry>): JSONSchema_ExtensionsEntry {
    const message = createBaseJSONSchema_ExtensionsEntry();
    message.key = object.key ?? "";
    message.value = object.value ?? undefined;
    return message;
  },
};

function createBaseTag(): Tag {
  return { name: "", description: "", externalDocs: undefined, extensions: {} };
}

export const Tag = {
  encode(message: Tag, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.externalDocs !== undefined) {
      ExternalDocumentation.encode(message.externalDocs, writer.uint32(26).fork()).ldelim();
    }
    Object.entries(message.extensions).forEach(([key, value]) => {
      if (value !== undefined) {
        Tag_ExtensionsEntry.encode({ key: key as any, value }, writer.uint32(34).fork()).ldelim();
      }
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Tag {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTag();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.description = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.externalDocs = ExternalDocumentation.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          const entry4 = Tag_ExtensionsEntry.decode(reader, reader.uint32());
          if (entry4.value !== undefined) {
            message.extensions[entry4.key] = entry4.value;
          }
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Tag {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      description: isSet(object.description) ? String(object.description) : "",
      externalDocs: isSet(object.externalDocs) ? ExternalDocumentation.fromJSON(object.externalDocs) : undefined,
      extensions: isObject(object.extensions)
        ? Object.entries(object.extensions).reduce<{ [key: string]: any | undefined }>((acc, [key, value]) => {
          acc[key] = value as any | undefined;
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: Tag): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.description !== undefined && (obj.description = message.description);
    message.externalDocs !== undefined &&
      (obj.externalDocs = message.externalDocs ? ExternalDocumentation.toJSON(message.externalDocs) : undefined);
    obj.extensions = {};
    if (message.extensions) {
      Object.entries(message.extensions).forEach(([k, v]) => {
        obj.extensions[k] = v;
      });
    }
    return obj;
  },

  create(base?: DeepPartial<Tag>): Tag {
    return Tag.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Tag>): Tag {
    const message = createBaseTag();
    message.name = object.name ?? "";
    message.description = object.description ?? "";
    message.externalDocs = (object.externalDocs !== undefined && object.externalDocs !== null)
      ? ExternalDocumentation.fromPartial(object.externalDocs)
      : undefined;
    message.extensions = Object.entries(object.extensions ?? {}).reduce<{ [key: string]: any | undefined }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[key] = value;
        }
        return acc;
      },
      {},
    );
    return message;
  },
};

function createBaseTag_ExtensionsEntry(): Tag_ExtensionsEntry {
  return { key: "", value: undefined };
}

export const Tag_ExtensionsEntry = {
  encode(message: Tag_ExtensionsEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      Value.encode(Value.wrap(message.value), writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Tag_ExtensionsEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTag_ExtensionsEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.value = Value.unwrap(Value.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Tag_ExtensionsEntry {
    return { key: isSet(object.key) ? String(object.key) : "", value: isSet(object?.value) ? object.value : undefined };
  },

  toJSON(message: Tag_ExtensionsEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined && (obj.value = message.value);
    return obj;
  },

  create(base?: DeepPartial<Tag_ExtensionsEntry>): Tag_ExtensionsEntry {
    return Tag_ExtensionsEntry.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Tag_ExtensionsEntry>): Tag_ExtensionsEntry {
    const message = createBaseTag_ExtensionsEntry();
    message.key = object.key ?? "";
    message.value = object.value ?? undefined;
    return message;
  },
};

function createBaseSecurityDefinitions(): SecurityDefinitions {
  return { security: {} };
}

export const SecurityDefinitions = {
  encode(message: SecurityDefinitions, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    Object.entries(message.security).forEach(([key, value]) => {
      SecurityDefinitions_SecurityEntry.encode({ key: key as any, value }, writer.uint32(10).fork()).ldelim();
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SecurityDefinitions {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSecurityDefinitions();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          const entry1 = SecurityDefinitions_SecurityEntry.decode(reader, reader.uint32());
          if (entry1.value !== undefined) {
            message.security[entry1.key] = entry1.value;
          }
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SecurityDefinitions {
    return {
      security: isObject(object.security)
        ? Object.entries(object.security).reduce<{ [key: string]: SecurityScheme }>((acc, [key, value]) => {
          acc[key] = SecurityScheme.fromJSON(value);
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: SecurityDefinitions): unknown {
    const obj: any = {};
    obj.security = {};
    if (message.security) {
      Object.entries(message.security).forEach(([k, v]) => {
        obj.security[k] = SecurityScheme.toJSON(v);
      });
    }
    return obj;
  },

  create(base?: DeepPartial<SecurityDefinitions>): SecurityDefinitions {
    return SecurityDefinitions.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SecurityDefinitions>): SecurityDefinitions {
    const message = createBaseSecurityDefinitions();
    message.security = Object.entries(object.security ?? {}).reduce<{ [key: string]: SecurityScheme }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[key] = SecurityScheme.fromPartial(value);
        }
        return acc;
      },
      {},
    );
    return message;
  },
};

function createBaseSecurityDefinitions_SecurityEntry(): SecurityDefinitions_SecurityEntry {
  return { key: "", value: undefined };
}

export const SecurityDefinitions_SecurityEntry = {
  encode(message: SecurityDefinitions_SecurityEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      SecurityScheme.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SecurityDefinitions_SecurityEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSecurityDefinitions_SecurityEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.value = SecurityScheme.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SecurityDefinitions_SecurityEntry {
    return {
      key: isSet(object.key) ? String(object.key) : "",
      value: isSet(object.value) ? SecurityScheme.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: SecurityDefinitions_SecurityEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined && (obj.value = message.value ? SecurityScheme.toJSON(message.value) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SecurityDefinitions_SecurityEntry>): SecurityDefinitions_SecurityEntry {
    return SecurityDefinitions_SecurityEntry.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SecurityDefinitions_SecurityEntry>): SecurityDefinitions_SecurityEntry {
    const message = createBaseSecurityDefinitions_SecurityEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null)
      ? SecurityScheme.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseSecurityScheme(): SecurityScheme {
  return {
    type: 0,
    description: "",
    name: "",
    in: 0,
    flow: 0,
    authorizationUrl: "",
    tokenUrl: "",
    scopes: undefined,
    extensions: {},
  };
}

export const SecurityScheme = {
  encode(message: SecurityScheme, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.type !== 0) {
      writer.uint32(8).int32(message.type);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.name !== "") {
      writer.uint32(26).string(message.name);
    }
    if (message.in !== 0) {
      writer.uint32(32).int32(message.in);
    }
    if (message.flow !== 0) {
      writer.uint32(40).int32(message.flow);
    }
    if (message.authorizationUrl !== "") {
      writer.uint32(50).string(message.authorizationUrl);
    }
    if (message.tokenUrl !== "") {
      writer.uint32(58).string(message.tokenUrl);
    }
    if (message.scopes !== undefined) {
      Scopes.encode(message.scopes, writer.uint32(66).fork()).ldelim();
    }
    Object.entries(message.extensions).forEach(([key, value]) => {
      if (value !== undefined) {
        SecurityScheme_ExtensionsEntry.encode({ key: key as any, value }, writer.uint32(74).fork()).ldelim();
      }
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SecurityScheme {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSecurityScheme();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 8) {
            break;
          }

          message.type = reader.int32() as any;
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.description = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.name = reader.string();
          continue;
        case 4:
          if (tag != 32) {
            break;
          }

          message.in = reader.int32() as any;
          continue;
        case 5:
          if (tag != 40) {
            break;
          }

          message.flow = reader.int32() as any;
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.authorizationUrl = reader.string();
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.tokenUrl = reader.string();
          continue;
        case 8:
          if (tag != 66) {
            break;
          }

          message.scopes = Scopes.decode(reader, reader.uint32());
          continue;
        case 9:
          if (tag != 74) {
            break;
          }

          const entry9 = SecurityScheme_ExtensionsEntry.decode(reader, reader.uint32());
          if (entry9.value !== undefined) {
            message.extensions[entry9.key] = entry9.value;
          }
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SecurityScheme {
    return {
      type: isSet(object.type) ? securityScheme_TypeFromJSON(object.type) : 0,
      description: isSet(object.description) ? String(object.description) : "",
      name: isSet(object.name) ? String(object.name) : "",
      in: isSet(object.in) ? securityScheme_InFromJSON(object.in) : 0,
      flow: isSet(object.flow) ? securityScheme_FlowFromJSON(object.flow) : 0,
      authorizationUrl: isSet(object.authorizationUrl) ? String(object.authorizationUrl) : "",
      tokenUrl: isSet(object.tokenUrl) ? String(object.tokenUrl) : "",
      scopes: isSet(object.scopes) ? Scopes.fromJSON(object.scopes) : undefined,
      extensions: isObject(object.extensions)
        ? Object.entries(object.extensions).reduce<{ [key: string]: any | undefined }>((acc, [key, value]) => {
          acc[key] = value as any | undefined;
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: SecurityScheme): unknown {
    const obj: any = {};
    message.type !== undefined && (obj.type = securityScheme_TypeToJSON(message.type));
    message.description !== undefined && (obj.description = message.description);
    message.name !== undefined && (obj.name = message.name);
    message.in !== undefined && (obj.in = securityScheme_InToJSON(message.in));
    message.flow !== undefined && (obj.flow = securityScheme_FlowToJSON(message.flow));
    message.authorizationUrl !== undefined && (obj.authorizationUrl = message.authorizationUrl);
    message.tokenUrl !== undefined && (obj.tokenUrl = message.tokenUrl);
    message.scopes !== undefined && (obj.scopes = message.scopes ? Scopes.toJSON(message.scopes) : undefined);
    obj.extensions = {};
    if (message.extensions) {
      Object.entries(message.extensions).forEach(([k, v]) => {
        obj.extensions[k] = v;
      });
    }
    return obj;
  },

  create(base?: DeepPartial<SecurityScheme>): SecurityScheme {
    return SecurityScheme.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SecurityScheme>): SecurityScheme {
    const message = createBaseSecurityScheme();
    message.type = object.type ?? 0;
    message.description = object.description ?? "";
    message.name = object.name ?? "";
    message.in = object.in ?? 0;
    message.flow = object.flow ?? 0;
    message.authorizationUrl = object.authorizationUrl ?? "";
    message.tokenUrl = object.tokenUrl ?? "";
    message.scopes = (object.scopes !== undefined && object.scopes !== null)
      ? Scopes.fromPartial(object.scopes)
      : undefined;
    message.extensions = Object.entries(object.extensions ?? {}).reduce<{ [key: string]: any | undefined }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[key] = value;
        }
        return acc;
      },
      {},
    );
    return message;
  },
};

function createBaseSecurityScheme_ExtensionsEntry(): SecurityScheme_ExtensionsEntry {
  return { key: "", value: undefined };
}

export const SecurityScheme_ExtensionsEntry = {
  encode(message: SecurityScheme_ExtensionsEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      Value.encode(Value.wrap(message.value), writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SecurityScheme_ExtensionsEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSecurityScheme_ExtensionsEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.value = Value.unwrap(Value.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SecurityScheme_ExtensionsEntry {
    return { key: isSet(object.key) ? String(object.key) : "", value: isSet(object?.value) ? object.value : undefined };
  },

  toJSON(message: SecurityScheme_ExtensionsEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined && (obj.value = message.value);
    return obj;
  },

  create(base?: DeepPartial<SecurityScheme_ExtensionsEntry>): SecurityScheme_ExtensionsEntry {
    return SecurityScheme_ExtensionsEntry.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SecurityScheme_ExtensionsEntry>): SecurityScheme_ExtensionsEntry {
    const message = createBaseSecurityScheme_ExtensionsEntry();
    message.key = object.key ?? "";
    message.value = object.value ?? undefined;
    return message;
  },
};

function createBaseSecurityRequirement(): SecurityRequirement {
  return { securityRequirement: {} };
}

export const SecurityRequirement = {
  encode(message: SecurityRequirement, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    Object.entries(message.securityRequirement).forEach(([key, value]) => {
      SecurityRequirement_SecurityRequirementEntry.encode({ key: key as any, value }, writer.uint32(10).fork())
        .ldelim();
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SecurityRequirement {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSecurityRequirement();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          const entry1 = SecurityRequirement_SecurityRequirementEntry.decode(reader, reader.uint32());
          if (entry1.value !== undefined) {
            message.securityRequirement[entry1.key] = entry1.value;
          }
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SecurityRequirement {
    return {
      securityRequirement: isObject(object.securityRequirement)
        ? Object.entries(object.securityRequirement).reduce<
          { [key: string]: SecurityRequirement_SecurityRequirementValue }
        >((acc, [key, value]) => {
          acc[key] = SecurityRequirement_SecurityRequirementValue.fromJSON(value);
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: SecurityRequirement): unknown {
    const obj: any = {};
    obj.securityRequirement = {};
    if (message.securityRequirement) {
      Object.entries(message.securityRequirement).forEach(([k, v]) => {
        obj.securityRequirement[k] = SecurityRequirement_SecurityRequirementValue.toJSON(v);
      });
    }
    return obj;
  },

  create(base?: DeepPartial<SecurityRequirement>): SecurityRequirement {
    return SecurityRequirement.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SecurityRequirement>): SecurityRequirement {
    const message = createBaseSecurityRequirement();
    message.securityRequirement = Object.entries(object.securityRequirement ?? {}).reduce<
      { [key: string]: SecurityRequirement_SecurityRequirementValue }
    >((acc, [key, value]) => {
      if (value !== undefined) {
        acc[key] = SecurityRequirement_SecurityRequirementValue.fromPartial(value);
      }
      return acc;
    }, {});
    return message;
  },
};

function createBaseSecurityRequirement_SecurityRequirementValue(): SecurityRequirement_SecurityRequirementValue {
  return { scope: [] };
}

export const SecurityRequirement_SecurityRequirementValue = {
  encode(message: SecurityRequirement_SecurityRequirementValue, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.scope) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SecurityRequirement_SecurityRequirementValue {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSecurityRequirement_SecurityRequirementValue();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.scope.push(reader.string());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SecurityRequirement_SecurityRequirementValue {
    return { scope: Array.isArray(object?.scope) ? object.scope.map((e: any) => String(e)) : [] };
  },

  toJSON(message: SecurityRequirement_SecurityRequirementValue): unknown {
    const obj: any = {};
    if (message.scope) {
      obj.scope = message.scope.map((e) => e);
    } else {
      obj.scope = [];
    }
    return obj;
  },

  create(
    base?: DeepPartial<SecurityRequirement_SecurityRequirementValue>,
  ): SecurityRequirement_SecurityRequirementValue {
    return SecurityRequirement_SecurityRequirementValue.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<SecurityRequirement_SecurityRequirementValue>,
  ): SecurityRequirement_SecurityRequirementValue {
    const message = createBaseSecurityRequirement_SecurityRequirementValue();
    message.scope = object.scope?.map((e) => e) || [];
    return message;
  },
};

function createBaseSecurityRequirement_SecurityRequirementEntry(): SecurityRequirement_SecurityRequirementEntry {
  return { key: "", value: undefined };
}

export const SecurityRequirement_SecurityRequirementEntry = {
  encode(message: SecurityRequirement_SecurityRequirementEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      SecurityRequirement_SecurityRequirementValue.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SecurityRequirement_SecurityRequirementEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSecurityRequirement_SecurityRequirementEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.value = SecurityRequirement_SecurityRequirementValue.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SecurityRequirement_SecurityRequirementEntry {
    return {
      key: isSet(object.key) ? String(object.key) : "",
      value: isSet(object.value) ? SecurityRequirement_SecurityRequirementValue.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: SecurityRequirement_SecurityRequirementEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined &&
      (obj.value = message.value ? SecurityRequirement_SecurityRequirementValue.toJSON(message.value) : undefined);
    return obj;
  },

  create(
    base?: DeepPartial<SecurityRequirement_SecurityRequirementEntry>,
  ): SecurityRequirement_SecurityRequirementEntry {
    return SecurityRequirement_SecurityRequirementEntry.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<SecurityRequirement_SecurityRequirementEntry>,
  ): SecurityRequirement_SecurityRequirementEntry {
    const message = createBaseSecurityRequirement_SecurityRequirementEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null)
      ? SecurityRequirement_SecurityRequirementValue.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseScopes(): Scopes {
  return { scope: {} };
}

export const Scopes = {
  encode(message: Scopes, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    Object.entries(message.scope).forEach(([key, value]) => {
      Scopes_ScopeEntry.encode({ key: key as any, value }, writer.uint32(10).fork()).ldelim();
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Scopes {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseScopes();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          const entry1 = Scopes_ScopeEntry.decode(reader, reader.uint32());
          if (entry1.value !== undefined) {
            message.scope[entry1.key] = entry1.value;
          }
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Scopes {
    return {
      scope: isObject(object.scope)
        ? Object.entries(object.scope).reduce<{ [key: string]: string }>((acc, [key, value]) => {
          acc[key] = String(value);
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: Scopes): unknown {
    const obj: any = {};
    obj.scope = {};
    if (message.scope) {
      Object.entries(message.scope).forEach(([k, v]) => {
        obj.scope[k] = v;
      });
    }
    return obj;
  },

  create(base?: DeepPartial<Scopes>): Scopes {
    return Scopes.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Scopes>): Scopes {
    const message = createBaseScopes();
    message.scope = Object.entries(object.scope ?? {}).reduce<{ [key: string]: string }>((acc, [key, value]) => {
      if (value !== undefined) {
        acc[key] = String(value);
      }
      return acc;
    }, {});
    return message;
  },
};

function createBaseScopes_ScopeEntry(): Scopes_ScopeEntry {
  return { key: "", value: "" };
}

export const Scopes_ScopeEntry = {
  encode(message: Scopes_ScopeEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== "") {
      writer.uint32(18).string(message.value);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Scopes_ScopeEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseScopes_ScopeEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.value = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Scopes_ScopeEntry {
    return { key: isSet(object.key) ? String(object.key) : "", value: isSet(object.value) ? String(object.value) : "" };
  },

  toJSON(message: Scopes_ScopeEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined && (obj.value = message.value);
    return obj;
  },

  create(base?: DeepPartial<Scopes_ScopeEntry>): Scopes_ScopeEntry {
    return Scopes_ScopeEntry.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Scopes_ScopeEntry>): Scopes_ScopeEntry {
    const message = createBaseScopes_ScopeEntry();
    message.key = object.key ?? "";
    message.value = object.value ?? "";
    return message;
  },
};

declare var self: any | undefined;
declare var window: any | undefined;
declare var global: any | undefined;
var tsProtoGlobalThis: any = (() => {
  if (typeof globalThis !== "undefined") {
    return globalThis;
  }
  if (typeof self !== "undefined") {
    return self;
  }
  if (typeof window !== "undefined") {
    return window;
  }
  if (typeof global !== "undefined") {
    return global;
  }
  throw "Unable to locate global object";
})();

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new tsProtoGlobalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isObject(value: any): boolean {
  return typeof value === "object" && value !== null;
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
