const loggerRedact = (process.env.LOGGER_REDACT ?? '')
  .split(',')
  .map((item) => item.trim())
  .filter((item) => item);
loggerRedact.push('config.mongodb.url');

module.exports = {
  logger: {
    level: process.env.LOGGER_LEVEL ?? 'info',
    redact: loggerRedact,
    requestIdHeader: process.env.REQUEST_ID_HEADER ?? 'x-correlation-id',
  },
  metrics: {
    enabled: process.env.METRICS_ENABLED !== 'false',
    defaultMetricsEnabled: process.env.METRICS_DEFAULT_ENABLED !== 'false',
    routeMetricsEnabled: process.env.METRICS_ROUTE_ENABLED !== 'false',
    interval: Number(process.env.METRICS_INTERVAL ?? 5000),
  },
  mongodb: {
    url: process.env.MONGODB_URL,
    opts: {
      autoIndex: process.env.MONGODB_AUTO_INDEX !== 'false',
      dbName: process.env.MONGODB_DB_NAME,
      minPoolSize: Number(process.env.MONGODB_MIN_POOL_SIZE ?? 0),
      maxPoolSize: Number(process.env.MONGODB_MAX_POOL_SIZE ?? 100),
    },
  },
  pagination: {
    limit: Number(process.env.PAGINATION_LIMIT ?? 25),
    maxLimit: Number(process.env.PAGINATION_MAX_LIMIT ?? 100),
  },
  routes: {
    create: process.env.CREATE_ONE_ROUTE_ENABLED !== 'false',
    delete: process.env.DELETE_ONE_ROUTE_ENABLED !== 'false',
    getMany: process.env.GET_MANY_ROUTE_ENABLED !== 'false',
    getOne: process.env.GET_ONE_ROUTE_ENABLED !== 'false',
    update: process.env.UPDATE_ONE_ROUTE_ENABLED !== 'false',
  },
  server: {
    healthcheckTimeout: Number(process.env.HEALTHCHECK_TIMEOUT ?? 10000),
    host: '0.0.0.0',
    routePrefix: process.env.ROUTE_PREFIX ?? '/crud',
    port: Number(process.env.http_port ?? 3000), // Use http_port for compatibility with OpenFaaS watchdog
  },
};
