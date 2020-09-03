import * as gax from 'google-gax';
import { Callback, ClientOptions } from 'google-gax';
import * as protosTypes from '../../protos/protos';
/**
 *  Each RPC normalizes the partition IDs of the keys in its input entities,
 *  and always returns entities with keys with normalized partition IDs.
 *  This applies to all keys and entities, including those in values, except keys
 *  with both an empty path and an empty or unset partition ID. Normalization of
 *  input keys sets the project ID (if not already set) to the project ID from
 *  the request.
 *
 * @class
 * @memberof v1
 */
export declare class DatastoreClient {
    private _descriptors;
    private _innerApiCalls;
    private _terminated;
    private _opts;
    private _gaxModule;
    private _gaxGrpc;
    private _protos;
    private _defaults;
    auth: gax.GoogleAuth;
    datastoreStub?: Promise<{
        [name: string]: Function;
    }>;
    /**
     * Construct an instance of DatastoreClient.
     *
     * @param {object} [options] - The configuration object. See the subsequent
     *   parameters for more details.
     * @param {object} [options.credentials] - Credentials object.
     * @param {string} [options.credentials.client_email]
     * @param {string} [options.credentials.private_key]
     * @param {string} [options.email] - Account email address. Required when
     *     using a .pem or .p12 keyFilename.
     * @param {string} [options.keyFilename] - Full path to the a .json, .pem, or
     *     .p12 key downloaded from the Google Developers Console. If you provide
     *     a path to a JSON file, the projectId option below is not necessary.
     *     NOTE: .pem and .p12 require you to specify options.email as well.
     * @param {number} [options.port] - The port on which to connect to
     *     the remote host.
     * @param {string} [options.projectId] - The project ID from the Google
     *     Developer's Console, e.g. 'grape-spaceship-123'. We will also check
     *     the environment variable GCLOUD_PROJECT for your project ID. If your
     *     app is running in an environment which supports
     *     {@link https://developers.google.com/identity/protocols/application-default-credentials Application Default Credentials},
     *     your project ID will be detected automatically.
     * @param {string} [options.apiEndpoint] - The domain name of the
     *     API remote host.
     */
    constructor(opts?: ClientOptions);
    /**
     * Initialize the client.
     * Performs asynchronous operations (such as authentication) and prepares the client.
     * This function will be called automatically when any class method is called for the
     * first time, but if you need to initialize it before calling an actual method,
     * feel free to call initialize() directly.
     *
     * You can await on this method if you want to make sure the client is initialized.
     *
     * @returns {Promise} A promise that resolves to an authenticated service stub.
     */
    initialize(): Promise<{
        [name: string]: Function;
    }>;
    /**
     * The DNS address for this API service.
     */
    static get servicePath(): string;
    /**
     * The DNS address for this API service - same as servicePath(),
     * exists for compatibility reasons.
     */
    static get apiEndpoint(): string;
    /**
     * The port for this API service.
     */
    static get port(): number;
    /**
     * The scopes needed to make gRPC calls for every method defined
     * in this service.
     */
    static get scopes(): string[];
    getProjectId(): Promise<string>;
    getProjectId(callback: Callback<string, undefined, undefined>): void;
    lookup(request: protosTypes.google.datastore.v1.ILookupRequest, options?: gax.CallOptions): Promise<[protosTypes.google.datastore.v1.ILookupResponse, protosTypes.google.datastore.v1.ILookupRequest | undefined, {} | undefined]>;
    lookup(request: protosTypes.google.datastore.v1.ILookupRequest, options: gax.CallOptions, callback: Callback<protosTypes.google.datastore.v1.ILookupResponse, protosTypes.google.datastore.v1.ILookupRequest | undefined, {} | undefined>): void;
    runQuery(request: protosTypes.google.datastore.v1.IRunQueryRequest, options?: gax.CallOptions): Promise<[protosTypes.google.datastore.v1.IRunQueryResponse, protosTypes.google.datastore.v1.IRunQueryRequest | undefined, {} | undefined]>;
    runQuery(request: protosTypes.google.datastore.v1.IRunQueryRequest, options: gax.CallOptions, callback: Callback<protosTypes.google.datastore.v1.IRunQueryResponse, protosTypes.google.datastore.v1.IRunQueryRequest | undefined, {} | undefined>): void;
    beginTransaction(request: protosTypes.google.datastore.v1.IBeginTransactionRequest, options?: gax.CallOptions): Promise<[protosTypes.google.datastore.v1.IBeginTransactionResponse, protosTypes.google.datastore.v1.IBeginTransactionRequest | undefined, {} | undefined]>;
    beginTransaction(request: protosTypes.google.datastore.v1.IBeginTransactionRequest, options: gax.CallOptions, callback: Callback<protosTypes.google.datastore.v1.IBeginTransactionResponse, protosTypes.google.datastore.v1.IBeginTransactionRequest | undefined, {} | undefined>): void;
    commit(request: protosTypes.google.datastore.v1.ICommitRequest, options?: gax.CallOptions): Promise<[protosTypes.google.datastore.v1.ICommitResponse, protosTypes.google.datastore.v1.ICommitRequest | undefined, {} | undefined]>;
    commit(request: protosTypes.google.datastore.v1.ICommitRequest, options: gax.CallOptions, callback: Callback<protosTypes.google.datastore.v1.ICommitResponse, protosTypes.google.datastore.v1.ICommitRequest | undefined, {} | undefined>): void;
    rollback(request: protosTypes.google.datastore.v1.IRollbackRequest, options?: gax.CallOptions): Promise<[protosTypes.google.datastore.v1.IRollbackResponse, protosTypes.google.datastore.v1.IRollbackRequest | undefined, {} | undefined]>;
    rollback(request: protosTypes.google.datastore.v1.IRollbackRequest, options: gax.CallOptions, callback: Callback<protosTypes.google.datastore.v1.IRollbackResponse, protosTypes.google.datastore.v1.IRollbackRequest | undefined, {} | undefined>): void;
    allocateIds(request: protosTypes.google.datastore.v1.IAllocateIdsRequest, options?: gax.CallOptions): Promise<[protosTypes.google.datastore.v1.IAllocateIdsResponse, protosTypes.google.datastore.v1.IAllocateIdsRequest | undefined, {} | undefined]>;
    allocateIds(request: protosTypes.google.datastore.v1.IAllocateIdsRequest, options: gax.CallOptions, callback: Callback<protosTypes.google.datastore.v1.IAllocateIdsResponse, protosTypes.google.datastore.v1.IAllocateIdsRequest | undefined, {} | undefined>): void;
    reserveIds(request: protosTypes.google.datastore.v1.IReserveIdsRequest, options?: gax.CallOptions): Promise<[protosTypes.google.datastore.v1.IReserveIdsResponse, protosTypes.google.datastore.v1.IReserveIdsRequest | undefined, {} | undefined]>;
    reserveIds(request: protosTypes.google.datastore.v1.IReserveIdsRequest, options: gax.CallOptions, callback: Callback<protosTypes.google.datastore.v1.IReserveIdsResponse, protosTypes.google.datastore.v1.IReserveIdsRequest | undefined, {} | undefined>): void;
    /**
     * Terminate the GRPC channel and close the client.
     *
     * The client will no longer be usable and all future behavior is undefined.
     */
    close(): Promise<void>;
}
