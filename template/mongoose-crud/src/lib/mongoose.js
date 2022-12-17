const mongoose = require('mongoose');
const config = require('../config');

module.exports = async (logger) => {
  mongoose.set('debug', (collection, method, query, ...args) => {
    /* Query should be redacted to hide personally-identifiable information */
    logger.debug({ collection, method, query, ...args }, 'MongoDB query');
  });

  try {
    /* Mongoose is a singleton internally, so connect before accepting connections */
    logger.info('Attempting to create Mongoose connection');
    await mongoose.connect(config.mongodb.url, config.mongodb.opts);
    logger.info('Connected to Mongoose');

    return mongoose;
  } catch (err) {
    logger.fatal({ err }, 'Unable to connect to Mongoose');

    /* Connection may be open */
    if (mongoose.connection.readyState === 1) {
      logger.debug('Killing Mongo connection before terminating');
      await mongoose.connection.close();
    }

    /* Preserve the error */
    throw err;
  }
};
