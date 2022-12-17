const { Schema } = require('mongoose');
const mongoosePaginate = require('mongoose-paginate-v2');

const uuid = require('uuid');

module.exports = class BaseSchema extends Schema {
  constructor(schema, opts = {}) {
    super(schema, {
      timestamps: true,
      ...opts,
    });

    this.add({
      _id: {
        type: String,
        default: () => uuid.v4(),
        immutable: true,
      },
    });

    this.set('toJSON', { virtuals: true }).set('toObject', { virtuals: true });

    if (opts.timestamps !== false) {
      this.index({
        createdAt: 1,
      });
    }

    this.plugin(mongoosePaginate);
  }
};
