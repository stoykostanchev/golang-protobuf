/**
 * @fileoverview gRPC-Web generated client stub for 
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = require('./joke_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.JokerClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.JokerPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.JokesRequest,
 *   !proto.JokesReply>}
 */
const methodDescriptor_Joker_GetJokes = new grpc.web.MethodDescriptor(
  '/Joker/GetJokes',
  grpc.web.MethodType.UNARY,
  proto.JokesRequest,
  proto.JokesReply,
  /**
   * @param {!proto.JokesRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.JokesReply.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.JokesRequest,
 *   !proto.JokesReply>}
 */
const methodInfo_Joker_GetJokes = new grpc.web.AbstractClientBase.MethodInfo(
  proto.JokesReply,
  /**
   * @param {!proto.JokesRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.JokesReply.deserializeBinary
);


/**
 * @param {!proto.JokesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.JokesReply)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.JokesReply>|undefined}
 *     The XHR Node Readable Stream
 */
proto.JokerClient.prototype.getJokes =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/Joker/GetJokes',
      request,
      metadata || {},
      methodDescriptor_Joker_GetJokes,
      callback);
};


/**
 * @param {!proto.JokesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.JokesReply>}
 *     A native promise that resolves to the response
 */
proto.JokerPromiseClient.prototype.getJokes =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/Joker/GetJokes',
      request,
      metadata || {},
      methodDescriptor_Joker_GetJokes);
};


module.exports = proto;

