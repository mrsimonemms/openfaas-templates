// eslint-disable-next-line import/no-extraneous-dependencies
const pino = require('pino');
const supertest = require('supertest');

const mongoose = require('../src/lib/mongoose');

describe('App test', () => {
  let app;
  let request;
  beforeEach(async () => {
    request = supertest(process.env.APP_URL);
    const { connection } = await mongoose(pino({ level: 'silent' }));

    /* This allows local development with live reload */
    await new Promise((resolve) => setTimeout(resolve, Number(process.env.START_TIMEOUT)));

    /* Get all collections */
    const collections = await connection.db.listCollections().toArray();

    /* And drop */
    await Promise.all(collections.map(({ name }) => connection.db.dropCollection(name)));
  });

  describe('/crud', () => {
    describe('GET', () => {
      beforeEach(async () => {
        const data = [
          {
            name: 'Test Testington',
            emailAddress: 'test@test.com',
          },
          {
            name: 'Test2 Testington',
            emailAddress: 'test2@test.com',
          },
          {
            name: 'Smith Smithson',
            emailAddress: 'smith@smithson.eu',
          },
        ];

        await Promise.all(data.map((input) => request.post('/crud').send(input).expect(201)));
      });

      it('should return the raw record', async () => {
        const { body } = await request.get('/crud').expect(200);

        expect(body.data).toEqual(expect.arrayContaining([]));
        expect(body.page).toBe(1);
        expect(body.pages).toBe(1);
        expect(body.count).toBe(3);
        expect(body.total).toBe(3);
      });

      it('should return the raw record with paging', async () => {
        const { body } = await request.get('/crud?limit=2&page=2').expect(200);

        expect(body.data).toEqual(expect.arrayContaining([]));
        expect(body.page).toBe(2);
        expect(body.pages).toBe(2);
        expect(body.count).toBe(1);
        expect(body.total).toBe(3);
      });

      it('should return the raw record filter', async () => {
        const { body } = await request
          .get('/crud?filter=emailAddress||$eq||test@test.com')
          .expect(200);

        expect(body.data).toEqual(expect.arrayContaining([]));
        expect(body.page).toBe(1);
        expect(body.pages).toBe(1);
        expect(body.count).toBe(1);
        expect(body.total).toBe(1);
      });
    });

    describe('POST', () => {
      beforeEach(() => {
        app = request.post('/crud');
      });

      it('should create a record', () =>
        app
          .send({
            name: 'Test Testington',
            emailAddress: 'test@test.com',
          })
          .expect(201));

      it('should fail basic validation', () =>
        app.send({}).expect(400, {
          emailAddress: {
            name: 'ValidatorError',
            message: 'Path `emailAddress` is required.',
            properties: {
              message: 'Path `emailAddress` is required.',
              type: 'required',
              path: 'emailAddress',
            },
            kind: 'required',
            path: 'emailAddress',
          },
          name: {
            name: 'ValidatorError',
            message: 'Path `name` is required.',
            properties: {
              message: 'Path `name` is required.',
              type: 'required',
              path: 'name',
            },
            kind: 'required',
            path: 'name',
          },
        }));
    });

    describe('/:id', () => {
      describe('DELETE', () => {
        it('should delete an existing record', async () => {
          const {
            body: { id },
          } = await request
            .post('/crud')
            .send({
              name: 'Test Testington',
              emailAddress: 'test@test.com',
            })
            .expect(201);

          await request.delete(`/crud/${id}`).expect(200, { id });
        });
      });

      describe('GET', () => {
        it('should delete an existing record', async () => {
          const data = {
            name: 'Test Testington',
            emailAddress: 'test@test.com',
          };

          const {
            body: { id },
          } = await request.post('/crud').send(data).expect(201);

          const { body } = await request.get(`/crud/${id}`).expect(200);

          expect(body).toEqual(expect.objectContaining(data));
        });
      });

      describe('PATCH', () => {
        it('should update an existing record', async () => {
          const {
            body: { id },
          } = await request
            .post('/crud')
            .send({
              name: 'Test Testington',
              emailAddress: 'test@test.com',
            })
            .expect(201);

          const { body } = await request.patch(`/crud/${id}`).send({ name: 'Test' }).expect(200);

          expect(body).toEqual(
            expect.objectContaining({
              id,
              name: 'Test',
              emailAddress: 'test@test.com',
            }),
          );
        });
      });
    });
  });

  describe('/health', () => {
    describe('GET', () => {
      it('should return the health check', () =>
        request.get('/health').expect(200, {
          isHealthy: true,
          checks: [
            {
              isHealthy: true,
              name: 'mongodb',
            },
          ],
        }));
    });
  });

  describe('/metrics', () => {
    describe('GET', () => {
      it('should return the metrics', () => request.get('/metrics').expect(200));
    });
  });
});
