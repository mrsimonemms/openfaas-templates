/* eslint-disable class-methods-use-this */

function queryMapper() {
  return [
    {
      query: '$between',
      factory: ({ value }) => {
        const [from, to] = value;
        if (!from) {
          throw new Error('Between query must have "from" defined');
        }

        if (!to) {
          throw new Error('Between query must have "to" defined');
        }

        return {
          $gte: from,
          $lte: to,
        };
      },
    },
    {
      query: '$cont',
      factory: ({ value }) => ({
        $regex: new RegExp(value),
      }),
    },
    {
      query: '$excl',
      factory: ({ value }) => ({
        $not: {
          $regex: new RegExp(value),
        },
      }),
    },
    {
      query: '$ends',
      factory: ({ value }) => ({
        $regex: new RegExp(`${value}$`),
      }),
    },
    {
      query: '$isnull',
      factory: () => ({
        $eq: null,
      }),
    },
    {
      operator: '$nin',
      query: '$notin',
    },
    {
      query: '$notnull',
      factory: () => ({
        $ne: null,
      }),
    },
    {
      query: '$starts',
      factory: ({ value }) => ({
        $regex: new RegExp(`^${value}`),
      }),
    },
  ];
}

function queryParser(factory) {
  return (item) => {
    const mapped = queryMapper().find((mapper) => mapper.query === item.operator);

    const value = {};
    if (mapped) {
      /* Operator defined - value needs to be cast */
      if (mapped.factory) {
        Object.assign(value, mapped.factory(item));
      } else {
        value[mapped.operator] = item.value;
      }
    } else {
      /* No operator mapped - use as-is */
      value[item.operator] = item.value;
    }

    factory(item, value);
  };
}

function queryBuilder({ filter, or }) {
  const query = {};

  filter.forEach(
    queryParser((item, value) => {
      query[item.field] = value;
    }),
  );
  or.forEach(
    queryParser((item, value) => {
      if (!Array.isArray(query.$or)) {
        query.$or = [];
      }

      query.$or.push({
        [item.field]: value,
      });
    }),
  );

  return query;
}

module.exports = class BaseController {
  constructor(config, crudRequest, utils) {
    this.config = config;
    this.crudRequest = crudRequest;
    this.utils = utils;
  }

  healthChecks({ mongoose }) {
    return [
      {
        name: 'mongodb',
        test: async () => {
          const { ok } = await mongoose.connection.db.admin().ping();

          return ok === 1;
        },
      },
    ];
  }

  async createOne({ req, res, Schema }) {
    req.log.debug('Creating new document');
    const document = new Schema(req.body);

    try {
      await document.validate();
    } catch (err) {
      req.log.warn({ err, errors: err.errors }, 'Create schema invalid');

      return res.status(400).send(err.errors);
    }

    await document.save({
      validateBeforeSave: false,
    });

    req.log.info({ id: document.get('id') }, 'Document saved to database');

    res.code(201);

    return document;
  }

  async deleteOne({ Schema, req, res }) {
    const { id } = req.params;

    req.log.debug({ id }, 'Deleting document');

    const data = await Schema.findByIdAndDelete(id);

    if (data === null) {
      req.log.warn({ id }, 'Document not deleted');

      res.code(404);
      return {
        message: `Item ${id} not found`,
        error: 'Not found',
        status: 404,
      };
    }

    req.log.info({ id }, 'Document deleted');

    return { id };
  }

  async getMany({ Schema, req, res }) {
    req.log.debug('Get many documents');

    const { page = 1, limit = this.config.pagination.limit } = req.query;
    let { sort = 'createdAt,ASC' } = req.query;
    let parser;
    try {
      parser = this.crudRequest.RequestQueryParser.create();
      parser.parseQuery(req.query);
    } catch (err) {
      req.log.error({ err }, 'Query parser error');
      res.code(400);
      return {
        message: err.message,
        error: 'Bad request',
        status: 400,
      };
    }

    if (!Array.isArray(sort)) {
      sort = [sort];
    }

    const sortBy = sort.reduce((result, item) => {
      const validSort = /^(\w+),(ASC|DESC)$/i.test(item);
      req.log.debug({ validSort, item }, 'Sort info');

      if (validSort) {
        const [column, type] = item.split(',');
        const direction = /^ASC$/i.test(type) ? 1 : -1;

        if (column) {
          Object.defineProperty(result, column, {
            enumerable: true,
            value: direction,
          });
        }
      }

      return result;
    }, {});

    let query;
    try {
      query = queryBuilder(parser);
    } catch (err) {
      req.log.error({ err }, 'Query builder error');
      res.code(400);
      return {
        message: err.message,
        error: 'Bad request',
        status: 400,
      };
    }

    req.log.info({ page, query, limit, sortBy }, 'getMany pagination');

    const { docs, totalPages, totalDocs } = await Schema.paginate(query, {
      limit,
      page,
      sort: sortBy,
    });

    return {
      data: docs,
      page, // Current page
      count: docs.length, // Records displayed
      pages: totalPages, // Total pages
      total: totalDocs, // Total records
    };
  }

  async getOne({ Schema, req, res }) {
    const { id } = req.params;
    req.log.debug({ id }, 'Get one document');

    const item = await Schema.findById(id);

    if (!item) {
      req.log.warn({ id }, 'Document not found');

      res.code(404);
      return {
        message: `Item ${id} not found`,
        error: 'Not found',
        status: 404,
      };
    }
    req.log.info({ id }, 'Document found');

    return item;
  }

  /**
   * Health
   *
   * To extend this, add to the healthChecks function
   *
   * @param input
   * @returns {Promise<*&{isHealthy: *}>}
   */
  async health(input) {
    const checks = (
      await Promise.allSettled(
        this.healthChecks(input).map(async (task) => {
          let isHealthy = false;

          try {
            isHealthy = await this.utils.promiseTimeout(
              this.config.server.healthcheckTimeout,
              task.test(),
            );
          } catch (err) {
            input.req.log.debug({ err, task }, 'Health check failed');
          }

          return {
            ...task,
            isHealthy,
          };
        }),
      )
    )
      .map(({ value }) => value)
      .map(({ isHealthy, name }) => ({
        isHealthy,
        name,
      }));

    const isHealthy = checks.every(({ isHealthy: healthy }) => healthy);

    const output = {
      isHealthy,
      checks,
    };

    if (!isHealthy) {
      input.req.log.error(output, 'Service unhealthy');
      input.res.code(503);
    } else {
      input.req.log.debug(output, 'Service healthy');
    }

    return output;
  }

  async updateOne({ Schema, req, res }) {
    const { id } = req.params;
    req.log.debug({ id }, 'Update existing document');

    const document = await Schema.findById(id);

    if (!document) {
      res.code(404);
      return {
        message: `Item ${id} not found`,
        error: 'Not found',
        status: 404,
      };
    }

    const { body } = req;
    Object.keys(body).forEach((key) => {
      document.set(key, body[key]);
    });

    try {
      await document.validate();
    } catch (err) {
      req.log.warn({ id, err, errors: err.errors }, 'Updated schema invalid');

      return res.status(400).send(err.errors);
    }

    await document.save({
      validateBeforeSave: false,
    });

    req.log.info({ id }, 'Updated document saved to database');

    return document;
  }
};
