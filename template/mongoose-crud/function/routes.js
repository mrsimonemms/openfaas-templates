/**
 * Routes
 *
 * Used to configure the custom routes on the CRUD application
 *
 * @param {*} BaseRoutes
 * @returns
 */

module.exports = (BaseRoutes) => {
  return class Routes extends BaseRoutes {
    /**
     * Custom Routes
     *
     * Define Fastify route(s)
     *
     * @link https://www.fastify.io/docs/latest/Reference/Routes/
     * @returns []{fastify.Route}
     */
    // eslint-disable-next-line class-methods-use-this
    customRoutes() {
      return [];
    }
  };
};
