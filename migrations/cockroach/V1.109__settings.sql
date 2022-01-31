ALTER TABLE zitadel.projections.iam ADD COLUMN default_language TEXT DEFAuLT '';

CREATE TABLE zitadel.projections.secret_generators (
    generator_type STRING NOT NULL
    , creation_date TIMESTAMPTZ NOT NULL
    , change_date TIMESTAMPTZ NOT NULL
    , resource_owner STRING NOT NULL
    , sequence INT8 NOT NULL

    , length STRING BIGINT NOT NULL
    , expiry STRING BIGINT NOT NULL
    , include_lower_letters BOOLEAN NOT NULL
    , include_upper_letters BOOLEAN NOT NULL
    , include_digits BOOLEAN NOT NULL
    , include_symbols BOOLEAN NOT NULL

    , PRIMARY KEY (generator_type)
);

	SecretGeneratorColumnLength              = "length"
	SecretGeneratorColumnExpiry              = "expiry"
	SecretGeneratorColumnIncludeLowerLetters = "include_lower_letters"
	SecretGeneratorColumnIncludeUpperLetters = "include_upper_letters"
	SecretGeneratorColumnIncludeDigits       = "include_digits"
	SecretGeneratorColumnIncludeSymbols      = "include_symbols"