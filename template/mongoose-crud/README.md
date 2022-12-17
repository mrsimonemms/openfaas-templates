# Mongoose CRUD

An OpenFaaS template to provide a CRUD application backed by Mongoose

## Getting Started

See [OpemFaaS Docs](https://github.com/openfaas/faas-cli/blob/master/guide/TEMPLATE.md) for
full information.

```shell
export OPENFAAS_TEMPLATE_URL=https://gitlab.com/MrSimonEmms/openfaas-templates

faas-cli new --lang mongoose-crud crud
```

This will create a new instance of this function in `./crud`

## Schema

This is the only required file to change. This receives an object with the
`mongoose` library injected and a `BaseSchema`

```javascript
module.exports = (mongoose, BaseSchema) => {
  const UserSchema = new BaseSchema({
    name: {
      type: String,
      required: true,
    },
  })

  return mongoose.model('User', UserSchema);
}
```

Using the `BaseSchema` is not required, but advised - it changes the `_id` to be
a UUID (which makes it easier to use in HTTP requests) and adds the `timestamps`
option by default. It also loads the [Mongoose Paginate V2](https://www.npmjs.com/package/mongoose-paginate-v2)
plugin, which is required for the `Get Many` endpoint. Aside from that, it is a
normal [Mongoose Schema](https://mongoosejs.com/docs/guide.html).

## Endpoints

This will expose the following endpoints:
- GET:/crud - Get Many
- POST:/crud - Create One
- GET:/crud/:id - Get One
- PATCH:/crud/:id - Update One
- DELETE:/crud:/id - Delete One

In all cases, the `:id` is UUID V4 specification. These can be disabled by setting
the following environment variables to `false`:
- `CREATE_ONE_ROUTE_ENABLED`
- `DELETE_ONE_ROUTE_ENABLED`
- `GET_MANY_ROUTE_ENABLED`
- `GET_ONE_ROUTE_ENABLED`
- `UPDATE_ONE_ROUTE_ENABLED`

The endpoints are prefixed with `/crud`. This can be controlled by setting the
`ROUTE_PREFIX` environment variable.

There are also some utility endpoints:
- GET:/health - health checks
- GET:/metrics - Prometheus metrics

### Get Many pagination and filtering

> Functionally, this is based upon [nestjsx/crud](https://github.com/nestjsx/crud/wiki/Requests#filter-conditions)
> and offers a compatible interface.

Of all the default routes, the Get Many is the most complex. It offers pagination
and filtering

#### Query Params

- `filter`:
- `limit`: number. Defaults to `PAGINATION_LIMIT` (or `25`) and max `PAGINATION_MAX_LIMIT` (or `100`)
- `or`:
- `page`: number. Defaults to `1`
- `sort`: string. Defaults to `createdAt,ASC`

#### Filter conditions

- **$eq** (=, equal)
- **$ne** (!=, not equal)
- **$gt** (>, greater than)
- **$lt** (<, lower that)
- **$gte** (>=, greater than or equal)
- **$lte** (<=, lower than or equal)
- **$starts** (LIKE val%, starts with)
- **$ends** (LIKE %val, ends with)
- **$cont** (LIKE %val%, contains)
- **$excl** (NOT LIKE %val%, not contains)
- **$in** (IN, in range, accepts multiple values)
- **$notin** (NOT IN, not in range, accepts multiple values)
- **$isnull** (IS NULL, is NULL, doesn't accept value)
- **$notnull** (IS NOT NULL, not NULL, doesn't accept value)
- **$between** (BETWEEN, between, accepts two values)

Response:
```json
{
  "data": [], // Array of models
  "page":  1, // Current page
  "count": 0, // Records displayed
  "pages": 1, // Total pages
  "total": 0 // Total records
}
```

#### filter

Adds fields request condition (multiple conditions) to your request.

Syntax:

> ?filter=field||$condition||value

Examples:

> ?filter=name||$eq||batman

> ?filter=isVillain||$eq||false&filter=city||$eq||Arkham (multiple filters are treated as a combination of AND type of conditions)

> ?filter=shots||$in||12,26 (some conditions accept multiple values separated by commas)

> ?filter=power||$isnull (some conditions don't accept value)

#### limit

Receive `N` amount of entities.

Syntax:

> ?limit=number

Example:

> ?limit=10

#### or

Adds OR conditions to the request.

Syntax:

> ?or=field||$condition||value

It uses the same filter conditions.

Rules and examples:

- If there is only one or present (without filter) then it will be interpreted as simple filter:

> ?or=name||$eq||batman

- If there are multiple or present (without filter) then it will be interpreted as a compination of OR conditions, as follows:
```
WHERE {or} OR {or} OR ...
```

> ?or=name||$eq||batman&or=name||$eq||joker

- If there are one or and one filter then it will be interpreted as OR condition, as follows:
```
WHERE {filter} OR {or}
```

> ?filter=name||$eq||batman&or=name||$eq||joker

- If present both or and filter in any amount (one or miltiple each) then both interpreted as a combitation of AND conditions and compared with each other by OR condition, as follows:
```
WHERE ({filter} AND {filter} AND ...) OR ({or} AND {or} AND ...)
```

> ?filter=type||$eq||hero&filter=status||$eq||alive&or=type||$eq||villain&or=status||$eq||dead

#### page

Receive a portion of limited amount of resources.

Syntax:

> ?page=number

Example:

> ?page=2

#### sort

Adds sort by field (by multiple fields) and order to query result.

Syntax:

> ?sort=field,ASC|DESC

Examples:

> ?sort=name,ASC

> ?sort=name,ASC&sort=id,DESC

### Middleware

Custom middleware can be registered by providing a `middleware.js` file.

### Custom Routes

Custom routes can be added by adding a `routes.js` file to your function. This
needs to export a factory function, which will receive the `BaseRoute` which you
can extend from.

> **NB**. The `url` here is the fully qualified URL and doesn't include the
> `ROUTE_PREFIX` environment variable.

This will use a normal [Fastify route](https://www.fastify.io/docs/latest/Routes)

```javascript
module.exports = (BaseRoutes) =>
  class Routes extends BaseRoutes {
    customRoutes() {
      return [
        {
          url: '/',
          method: 'get',
          handler: () => ({ hello: 2 }),
        },
      ];
    }
  };
```

If you want access to mongoose or the defined schema, you can wrap the `handler`
in the `handlerFactory` function:

```javascript
module.exports = (BaseRoutes) =>
  class Routes extends BaseRoutes {
    customRoutes() {
      return [
        {
          url: '/',
          method: 'get',
          handler: this.handlerFactory(({ mongoose, Schema, req, res }) => {
            console.log({
              mongoose, // The result of require('mongoose')
              Schema, // The schema constructor
              req, // The fastify request object
              res, // The fastify response object
            });
            return new Schema();
          }),
        },
      ];
    }
  };
```

## Controller

If you wish to extend the behaviour, you can add a `controller.js` file to the
function. This  needs to export a factory function, which will receive the `BaseController`
which you can extend from. This can then be used to extend or replace methods.

```javascript
module.exports = (BaseController) =>
  class Controller extends BaseController {
    // Extend the getOne method
    getOne(input) {
      console.log(`Retrieving model: ${input.req.params.id}`);

      return super.getOne(input);
    }

    // Replace the updateOne method
    async updateOne({ Schema, req }) {
      const document = await Schema.findById(req.params.id);

      ///////////////////////////////
      // Do things with the Schema //
      ///////////////////////////////

      return document;
    }
  };
```
