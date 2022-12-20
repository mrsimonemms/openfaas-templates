/**
 * Schema
 *
 * Return a Mongoose schema to define our model. This will where the majority of
 * the work will be done.
 *
 * @link https://mongoosejs.com/docs/guide.html
 * @param {*} param0
 * @param {*} BaseSchema
 * @returns {schema}
 */
module.exports = ({ model }, BaseSchema) => {
  const modelName = 'User';

  const UserSchema = new BaseSchema({
    name: {
      type: String,
      required: true,
    },
    emailAddress: {
      type: String,
      required: true,
      unique: true,
      validate: {
        async validator(emailAddress) {
          if (!/@/.test(emailAddress)) {
            throw new Error('Value is an invalid email address');
          }

          const user = await this.model(modelName)
            .findOne({
              emailAddress,
            })
            .where('_id')
            .nin([this.get('id')]);

          if (user !== null) {
            throw new Error('Email address already registered');
          }

          return true;
        },
      },
    },
  });

  return model(modelName, UserSchema);
};
