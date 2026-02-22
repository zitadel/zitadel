DROP TABLE IF EXISTS eventstore.events2;
DROP FUNCTION IF EXISTS eventstore.commands_to_events(commands eventstore.command[]), eventstore.push(commands eventstore.command[]);
DROP TYPE IF EXISTS eventstore.command;