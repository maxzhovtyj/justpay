INSERT INTO users (id,
                   created_at,
                   updated_at)
VALUES ('2c127d70-3b9b-4743-9c2e-74b9f617029f', now(), now());

INSERT INTO orders (id,
                    user_id,
                    status,
                    is_final,
                    created_at,
                    updated_at)
VALUES ('97a96c29-7631-4cbc-9559-f8866fb03392', '2c127d70-3b9b-4743-9c2e-74b9f617029f', 0, false, now(), now());