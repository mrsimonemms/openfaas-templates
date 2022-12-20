module.exports = class Routes {
  constructor(config, controller, handlerObj) {
    this.config = config;
    this.controller = controller;
    this.handlerObj = handlerObj;
  }

  get routesEnabled() {
    return this.config.routes;
  }

  // eslint-disable-next-line class-methods-use-this
  customRoutes() {
    return [];
  }

  handlerFactory(fn) {
    return (req, res) => {
      return fn({
        ...this.handlerObj,
        req,
        res,
      });
    };
  }

  async register(app, opts, done) {
    const tasks = this.routes().map(async (item) => {
      const route = {
        ...item,
        method: (item?.method ?? 'get').toUpperCase(),
        /* Replace any duplicate slashes with a single slash and remove trailing slashes */
        url: item.url.replace(/\/+/g, '/').replace(/\w\/+$/, ''),
      };

      /* Register the middleware */
      try {
        const { default: middleware } = await import('../../function/middleware.js'); // eslint-disable-line import/extensions, import/no-unresolved

        app.log.trace('Loading middleware');

        route.preHandler = middleware;
      } catch (err) {
        if (err.code === 'ERR_MODULE_NOT_FOUND') {
          app.log.trace({ err }, 'No middleware loaded');
        } else {
          app.log.error({ err }, 'Error loading middleware');
          throw err;
        }
      }

      /* Use !== false to allow custom routes to default to including */
      if (route.enabled !== false) {
        app.log.debug({ route }, 'New route registered');
        app.route(route);
      } else {
        app.log.debug({ route }, 'Route not enabled');
      }
    });

    await Promise.all(tasks);

    done();
  }

  routes() {
    const customRoutes = this.customRoutes();

    if (!Array.isArray(customRoutes)) {
      throw new Error('Custom routes must be an array');
    }

    const defaultRoutes = [
      {
        enabled: this.routesEnabled?.getMany ?? true,
        method: 'GET',
        url: this.config.server.routePrefix,
        schema: {
          querystring: {
            type: 'object',
            properties: {
              limit: {
                type: 'integer',
                maximum: this.config.pagination.maxLimit,
              },
              page: {
                type: 'integer',
              },
            },
          },
        },
        handler: this.handlerFactory((...args) => this.controller.getMany(...args)),
      },
      {
        enabled: this.routesEnabled?.create ?? true,
        method: 'POST',
        url: this.config.server.routePrefix,
        handler: this.handlerFactory((...args) => this.controller.createOne(...args)),
      },
      {
        enabled: this.routesEnabled?.getOne ?? true,
        method: 'GET',
        url: `${this.config.server.routePrefix}/:id`,
        handler: this.handlerFactory((...args) => this.controller.getOne(...args)),
      },
      {
        enabled: this.routesEnabled?.update ?? true,
        method: 'PATCH',
        url: `${this.config.server.routePrefix}/:id`,
        handler: this.handlerFactory((...args) => this.controller.updateOne(...args)),
      },
      {
        enabled: this.routesEnabled?.delete ?? true,
        method: 'DELETE',
        url: `${this.config.server.routePrefix}/:id`,
        handler: this.handlerFactory((...args) => this.controller.deleteOne(...args)),
      },
    ];

    customRoutes.push(...defaultRoutes);

    return customRoutes;
  }
};
