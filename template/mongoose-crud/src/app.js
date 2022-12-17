const fastify = require('fastify');
const metricsPlugin = require('fastify-metrics');
const crudRequest = require('@nestjsx/crud-request');
const fastifySensible = require('fastify-sensible');
const uuid = require('uuid');

const config = require('./config');
const mongooseBootstrap = require('./lib/mongoose');
const BaseController = require('./lib/controller');
const BaseRoute = require('./lib/routes');
const BaseSchema = require('./lib/schema');
const schema = require('../function/schema');
const utils = require('./lib/utils');

async function loadClass(logger, BaseClass, filePath) {
  try {
    const { default: factory } = await import(filePath);

    logger.trace({ filePath }, 'Loading custom and base class');

    return factory(BaseClass);
  } catch (err) {
    logger.trace({ err }, 'Loading base class only');

    return BaseClass;
  }
}

async function main(app) {
  const mongoose = await mongooseBootstrap(app.log);

  const Schema = await schema(mongoose, BaseSchema);

  const Routes = await loadClass(app.log, BaseRoute, '../function/routes.js');
  const Controller = await loadClass(app.log, BaseController, '../function/controller.js');

  const controller = new Controller(config, crudRequest, utils);
  const routes = new Routes(config, controller, { mongoose, Schema });

  app.register(fastifySensible);

  /* Health is a special route */
  app.get(
    '/health',
    routes.handlerFactory((...args) => controller.health(...args)),
  );
  if (config.metrics.enabled) {
    app.log.debug('Metrics enabled');
    app.register(metricsPlugin, {
      enableDefaultMetrics: config.metrics.defaultMetricsEnabled,
      enableRouteMetrics: config.metrics.routeMetricsEnabled,
      endpoint: '/metrics',
      interval: config.metrics.interval,
    });
  } else {
    app.log.debug('Metrics not enabled');
  }

  app.register((...args) => routes.register(...args));

  await app.listen(config.server.port, config.server.host);

  app.log.info({ config }, 'Server started');
}

const app = fastify({
  logger: {
    level: config.logger.level,
    redact: config.logger.redact,
  },
  requestIdHeader: config.logger.requestIdHeader,
  genReqId: () => uuid.v4(),
});

main(app).catch((err) => {
  app.log.error({ err }, 'Failed to start server');
  process.exit(1);
});
