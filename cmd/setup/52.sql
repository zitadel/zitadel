DO
$$
    BEGIN
        IF NOT EXISTS
            (SELECT 1
             FROM information_schema.tables
             WHERE table_schema = 'projections'
               AND table_name = 'idp_templates6_ldap2')
        THEN
            ALTER TABLE IF EXISTS projections.idp_templates6_ldap3 RENAME COLUMN rootCA TO root_ca;
            ALTER TABLE IF EXISTS projections.idp_templates6_ldap3 RENAME TO idp_templates6_ldap2;
        END IF;
    END
$$;
