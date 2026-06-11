CREATE OR REPLACE FUNCTION generate_uuid_v7()
    RETURNS uuid AS $$
    BEGIN
        RETURN (
        -- timestamp
        lpad(to_hex(((extract(epoch from now()) * 1000)::bigint >> 16)), 8, '0') || '-' ||
        lpad(to_hex(((extract(epoch from now()) * 1000
            + (date_part('milliseconds', now())::bigint % 1000))::bigint & 0xffff)), 4, '0') || '-' ||
        -- version / rand_a
        lpad(to_hex((0x7000 + (random() * 0x0fff)::int)), 4, '0') || '-' ||
        -- variant / rand_b
        lpad(to_hex((0x8000 + (random() * 0x3fff)::int)), 4, '0') || '-' ||
        -- rand_b
        lpad(to_hex((floor(random() * (2^48))::bigint >> 16)), 12, '0')
        )::uuid;  -- Cast the result to uuid
    END;
    $$ LANGUAGE plpgsql;