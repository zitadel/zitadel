export type Severity = '' | 'error' | 'warning';
export type LogType = 'cons:log' |
  'cons:info' |
  'cons:warn' |
  'cons:error' |
  'cons:debug' |
  'cy:log' |
  'cy:xhr' |
  'cy:fetch' |
  'cy:request' |
  'cy:route' |
  'cy:intercept' |
  'cy:command' |
  'ctr:info';

export interface SupportOptions {
  /**
   * What types of logs to collect and print.
   * By default all types are enabled.
   * The 'cy:command' is the general type that contain all types of commands that are not specially treated.
   * @default ['cons:log','cons:info', 'cons:warn', 'cons:error', 'cy:log', 'cy:xhr', 'cy:fetch', 'cy:request', 'cy:route', 'cy:command']
   */
  collectTypes?: readonly string[];

  /**
   * Callback to filter logs manually. The type is from the same list as for the collectTypes option.
   * Severity can be of ['', 'error', 'warning'].
   * @default undefined
   */
  filterLog?:
    | null
    | NonNullable<SupportOptions['collectTypes']>[number]
    | ((args: [/* type: */ LogType, /* message: */ string, /* severity: */ Severity]) => boolean);

  /**
   * Callback to process logs manually. The type is from the same list as for the collectTypes option.
   * Severity can be of ['', 'error', 'warning'].
   * @default undefined
   */
  processLog?:
      | null
      | NonNullable<SupportOptions['collectTypes']>[number]
      | ((args: [/* type: */ LogType, /* message: */ string, /* severity: */ Severity]) => [LogType, string, Severity]);

  /**
   * Callback to collect each test case's logs after its run.
   * @default undefined
   */
  collectTestLogs?: (
    context: {mochaRunnable: any, testState: string, testTitle: string, testLevel: number},
    messages: [/* type: */ LogType, /* message: */ string, /* severity: */ Severity][]
  ) => void;

  xhr?: {
    /**
     * Whether to print header data for XHR requests.
     * @default false
     */
    printHeaderData?: boolean;

    /**
     * Whether to print request data for XHR requests besides response data.
     * @default false
     */
    printRequestData?: boolean;
  };

  /**
   * Enables extended log collection: including after all and before all hooks.
   * @unstable
   * @default false
   */
  enableExtendedCollector: boolean;

  /**
   * Enables continuous logging of logs to terminal one by one, as they get registerd or modified.
   * @unstable
   * @default false
   */
  enableContinuousLogging: boolean;

  /**
   * Enabled debug logging.
   * @default false
   */
  debug: boolean;
}

declare function installLogsCollector(config?: SupportOptions): void;
export default installLogsCollector;
