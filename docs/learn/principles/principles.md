# Principles

> This is just an outline!

- Stateless Application design
- System of records is the eventstore
- Everything else need to be able to be regenerated
- Try not so solve complex problems outside of the IAM Domain
- Use scalable storage for the eventstore and querymodels
- Try to be idempotent when ever possible
- Create All-in-one binary
  - Reduce necessaty of system or external dependencies as much as possible
- Embrace testing
- Design API first
- Optimize all components for day-two operations
- Use only Opensource projects with permissive licenses
- Don't do crypto
- Embrace standard as much as possible
