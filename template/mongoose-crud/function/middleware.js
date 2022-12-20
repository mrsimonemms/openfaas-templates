/**
 * Middleware
 *
 * An array of Fastify preHandler hooks that are invoked before the routes are
 * triggered.
 *
 * This will be triggered on all CRUD routes. If you want to target specific endpoints/
 * methods then the middleware function(s) will need to apply this conditionally.
 *
 * The function receives "request", "reply" and "done".
 *
 * @link https://www.fastify.io/docs/latest/Reference/Hooks/#prehandler
 */
module.exports = [];
